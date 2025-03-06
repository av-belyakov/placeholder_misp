package coremodule

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/av-belyakov/objectsmispformat"
	"github.com/av-belyakov/placeholder_misp/cmd/mispapi"
	"github.com/av-belyakov/placeholder_misp/commoninterfaces"
	"github.com/av-belyakov/placeholder_misp/internal/countermessage"
	rules "github.com/av-belyakov/placeholder_misp/internal/ruleshandler"
)

// CreateObjectsFormatMISP создаёт набор объектов в формате MISP
func CreateObjectsFormatMISP(
	chDecodeJSON <-chan ChanInputCreateMispFormat,
	taskId string,
	mispModule mispapi.ModuleMispHandler,
	listRule *rules.ListRule,
	counting *countermessage.CounterMessage,
	logger commoninterfaces.Logger) {
	var (
		maxCountObservables int
		seqNumObservable    int
		caseId              float64
		rootId              string
		userEmail           string
		caseSource          string
		patterIsNum         *regexp.Regexp = regexp.MustCompile(`^\d+$`)
	)

	//формируем шаблоны для заполнения
	eventsMisp := objectsmispformat.NewEventMisp()
	listObjectsMisp := objectsmispformat.NewListObjectsMispFormat()
	listAttributeTmp := objectsmispformat.NewListAttributeTmp()
	listAttributesMisp := objectsmispformat.NewListAttributesMispFormat()
	leot := objectsmispformat.NewListEventObjectTags()

	supportiveListExcludeRule := NewSupportiveListExcludeRuleTmp(listRule.GetRuleExclude())

	//это основной обработчик параметров входящего объекта
	listHandlerMisp := map[string][]func(interface{}, int){
		//event -> events
		"event.object.title":     {eventsMisp.SetAnyInfo},
		"event.object.startDate": {eventsMisp.SetAnyTimestamp},
		"event.details.endDate":  {eventsMisp.SetAnyDate},
		"event.object.tlp":       {eventsMisp.SetAnyDistribution},
		"event.object.severity":  {eventsMisp.SetAnyThreatLevelId},
		"event.organisationId":   {eventsMisp.SetAnyOrgId},
		"event.object.updatedAt": {eventsMisp.SetAnySightingTimestamp},
		"event.object.owner":     {eventsMisp.SetAnyEventCreatorEmail},
		"event.object.customFields.class-attack.string": {func(i interface{}, num int) {
			leot.SetTag(fmt.Sprintf("class-attack=\"%v\"", i))
		}},
		//observables -> attributes
		"observables._id": {
			listAttributesMisp.SetValueObjectId,
			listObjectsMisp.SetValueId,
		},
		"observables.data": {listAttributesMisp.SetValueValue},
		"observables.dataType": {
			listObjectsMisp.SetValueName,
			//здесь выполняем автоподстановку значений для полей Type и Category
			//объекта AttributesMisp на основе определенной логике и уже предустановленных
			//значений, при этом значения, заданные пользователем для этих полей, обрабатываются
			//отдельно и хранятся в listTags, а после закрытия канала совмещаются с
			//объектами AttributesMisp и следовательно перезаписывают значения выполненные
			//через автоподстановку
			listAttributesMisp.HandlingValueDataType,
		},
		"observables._createdAt": {listAttributesMisp.SetValueTimestamp},
		"observables.message":    {listAttributesMisp.SetValueComment},
		"observables.startDate": {
			listAttributesMisp.SetValueFirstSeen,
			listObjectsMisp.SetValueFirstSeen,
			listObjectsMisp.SetValueTimestamp,
		},
		//observables.attachment -> objects
		"observables.attachment.size": {listObjectsMisp.SetValueSize},
	}

	svn := NewStorageValueName()
	exclusionRules := NewExclusionRules()
	listGalaxyTags := NewMispGalaxyTags()
	addFuncGalaxyTags := addListGalaxyTags(listGalaxyTags)

	for key := range listHandlerMisp {
		if strings.Contains(key, "observables") {
			maxCountObservables++
		}
	}

	listTags := make(map[int][2]string)

	for msg := range chDecodeJSON {
		//ищем источник события
		if source, ok := searchEventSource(msg); ok {
			caseSource = source
		}

		//ищем id события
		if cid, ok := searchCaseId(msg); ok {
			caseId = cid
		}

		// ищем email владельца события
		if uemail, ok := searchOwnerEmail(msg); ok {
			userEmail = uemail
		}

		//*************** Обработка правил ***************
		//обработка правил REPLACEMENT (замена)
		newValue, _, err := listRule.ReplacementRuleHandler(msg.ValueType, msg.FieldBranch, msg.Value)
		if err != nil {
			logger.Send("warning", fmt.Sprintf("search value '%s' from rule of section 'REPLACE' is not fulfilled", msg.Value))
		}
		//обработка правил PASS (пропуск)
		listRule.PassRuleHandler(msg.FieldBranch, newValue)
		//**********************************************

		//для observables которые содержат свойства, являющиеся картами,
		//такими как свойства 'attachment', 'reports' и т.д. не
		//осуществлять подсчет свойств
		isObservables := strings.Contains(msg.FieldBranch, "observables")
		countOne := strings.Count(msg.FieldBranch, ".") <= 1

		if isObservables && countOne {
			var newFieldName = msg.FieldName
			//сделал проверку на число что бы исключить повторение
			//имен для свойств являющихся срезами, так как в данной ситуации
			//имя содержащееся в msg.FieldName представляет собой числовой
			//индекс, соответственное, если будет еще одно свойство являющееся
			//срезом, то может быть совпадение имен и изменение seqNumObservable, а как
			//результат будет переход на другой объект 'observables'
			if patterIsNum.MatchString(msg.FieldName) {
				tmp := strings.Split(msg.FieldBranch, ".")
				var nameTmp string
				if len(tmp) > 0 {
					nameTmp = tmp[len(tmp)-1] + "_"
				}

				newFieldName = nameTmp + msg.FieldName
			}

			//подсчет свойств для объектов типа 'observables' выполняется для
			//того что бы отделить один объект 'observables' от другого
			if svn.GetValueName(newFieldName) {
				svn.CleanValueName()
				seqNumObservable++
			}

			svn.SetValueName(newFieldName)

			// ************* обработка правил EXCLUSION (исключения) *************
			if addrRule, isEqual := listRule.ExcludeRuleHandler(msg.FieldBranch, newValue); isEqual {
				supportiveListExcludeRule.Add(seqNumObservable, msg.FieldBranch, newValue, addrRule, isEqual)
			}
		}

		//обрабатываем свойство observables.attachment
		if strings.Contains(msg.FieldBranch, "attachment") {
			listAttributeTmp.AddAttribute(msg.FieldBranch, newValue, seqNumObservable)
		}

		//получаем rootId
		if msg.FieldBranch == "event.rootId" || msg.FieldBranch == "event.object.rootId" {
			rid, ok := newValue.(string)
			if ok && rootId == "" {
				rootId = rid
			}
		}

		//обрабатываем свойство event.object.tags, оно ответственно за
		//наполнение поля "Теги" MISP
		if msg.FieldBranch == "event.object.tags" {
			if tag, ok := newValue.(string); ok {
				leot.SetTag(tag)
			}
		}

		//заполняем временный объект listGalaxyTags данными, предназначенными
		//для формирования специализированных тегов на основе которых в MISP
		//будут формироватся галактики
		addFuncGalaxyTags(msg.FieldBranch, newValue)

		//обрабатываем свойство observables.tags
		if msg.FieldBranch == "observables.tags" {
			listTags = handlerObservablesTags(newValue, listTags, listAttributesMisp, seqNumObservable)
		}

		//проверяем есть ли путь до обрабатываемого свойства в списке обработчиков
		//это главный обработчик который выполняет всю работу, особенно если выше идущие
		//обработчики не выполнялись так как принимаемые значения не соответствовали
		//параметрам для их вызова
		lf, ok := listHandlerMisp[msg.FieldBranch]
		if ok {
			//основной обработчик путей из msg.FieldBranch
			for _, f := range lf {
				f(newValue, seqNumObservable)
			}
		}
	}

	var isAllowed bool
	//проверяем что бы хотя бы одно правило разрешало пропуск кейса
	if listRule.GetRulePassany() || listRule.SomePassRuleIsTrue() {
		isAllowed = true

		//сетчик кейсов соответствующих или не соответствующих правилам
		counting.SendMessage("update events meet rules", 1)
	}

	//удаляем те объекты Attributes которые соответствуют правилам EXCLUDE
	delElementAttributes(exclusionRules, listAttributesMisp, logger)

	//выполняет очистку значения StatementExpression что равно отсутствию совпадений в правилах Pass
	// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	// ВРЕМЕННО КОМЕНТИРУЕТСЯ для проведения тестов
	//
	// в production не забывать убирать коментарий
	//
	listRule.CleanStatementExpressionRulePass()

	if !isAllowed {
		logger.Send("warning", fmt.Sprintf("the message with case id %d was not sent to MISP because it does not comply with the rules", int(caseId)))
	} else {
		//добавляем case id в поле Info
		//		eventsMisp.Info += fmt.Sprintf(" :::TheHive case id '%d':::", int(caseId))
		eventsMisp.Info += fmt.Sprintf(" TesT:_TheHive case id '%d'_:TesT", int(caseId))

		//добавляем в datemodels.ListObjectEventTags дополнительные теги
		//ответственные за формирование галактик в MISP
		galaxyTags := createGalaxyTags(listGalaxyTags)
		joinEventTags(leot, galaxyTags)

		for k := range listAttributesMisp.GetList() {
			if supportiveListExcludeRule.CheckRuleTrue(k) {
				//удаляем элементы подходящие под правила группы EXCLUDE
				listAttributesMisp.DelElementList(k)
			}
		}

		mispFormat := objectsmispformat.NewListFormatsMISP()
		mispFormat.ID = rootId
		mispFormat.Event = eventsMisp
		mispFormat.Objects = getNewListObjects(listObjectsMisp.GetList(), listAttributeTmp.GetListAttribute())
		tmpListTags := leot.GetListTags()
		mispFormat.ObjectTags = &tmpListTags
		mispFormat.Attributes = getNewListAttributes(listAttributesMisp.GetList(), listTags)
		reports := objectsmispformat.NewEventReports()
		reports.SetName(fmt.Sprint(caseId))
		reports.SetDistribution("1")
		mispFormat.Reports = reports

		logger.Send("info", fmt.Sprintf("the case with id:'%d' complies with the specified rules and has been submitted for further processing", int(caseId)))

		//тут отправляем сформированные по формату MISP пользовательские структуры
		mispModule.SendDataInput(mispapi.InputSettings{
			Command:    "add event",
			TaskId:     taskId,
			CaseId:     caseId,
			RootId:     rootId,
			CaseSource: caseSource,
			UserEmail:  userEmail,
			Data:       *mispFormat,
		})
	}

	//очищаем события, список аттрибутов и текущий email пользователя
	userEmail, caseSource = "", ""
	leot.CleanListTags()
	listObjectsMisp.CleanList()
	listAttributeTmp.CleanAttribute()
	listAttributesMisp.CleanList()
}
