package coremodule

import (
	"fmt"
	"regexp"
	"runtime"
	"strings"

	"github.com/av-belyakov/placeholder_misp/cmd/mispapi"
	"github.com/av-belyakov/placeholder_misp/commoninterfaces"
	"github.com/av-belyakov/placeholder_misp/internal/countermessage"
	"github.com/av-belyakov/placeholder_misp/internal/datamodels"
	rules "github.com/av-belyakov/placeholder_misp/rulesinteraction"
)

func NewMispFormat(
	chanOutputDecodeJson <-chan ChanInputCreateMispFormat,
	taskId string,
	mispModule mispapi.ModuleMispHandler,
	listRule *rules.ListRule,
	counting *countermessage.CounterMessage,
	logger commoninterfaces.Logger) {

	eventsMisp := datamodels.NewEventMisp()
	listObjectsMisp := datamodels.NewListObjectsMispFormat()
	listAttributeTmp := datamodels.NewListAttributeTmp()
	listAttributesMisp := datamodels.NewListAttributesMispFormat()
	leot := datamodels.NewListEventObjectTags()

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

	var (
		maxCountObservables int
		seqNumObservable    int
		caseId              float64
		rootId              string
		userEmail           string
		caseSource          string
		patterIsNum         *regexp.Regexp = regexp.MustCompile(`^\d+$`)
	)

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

	for tmf := range chanOutputDecodeJson {
		//ищем источник события
		if source, ok := searchEventSource(tmf); ok {
			caseSource = source
		}

		//ищем id события
		if cid, ok := searchCaseId(tmf); ok {
			caseId = cid
		}

		// ищем email владельца события
		if uemail, ok := searchOwnerEmail(tmf); ok {
			userEmail = uemail
		}

		//*************** Обработка правил ***************
		//обработка правил REPLACEMENT (замена)
		newValue, _, err := listRule.ReplacementRuleHandler(tmf.ValueType, tmf.FieldBranch, tmf.Value)
		if err != nil {
			_, f, l, _ := runtime.Caller(0)
			logger.Send("warning", fmt.Sprintf("'search value \"%s\" from rule of section \"REPLACE\" is not fulfilled' %s:%d", tmf.Value, f, l-2))
		}
		//обработка правил PASS (пропуск)
		listRule.PassRuleHandler(tmf.FieldBranch, newValue)
		//**********************************************

		//для observables которые содержат свойства, являющиеся картами,
		//такими как свойства 'attachment', 'reports' и т.д. не
		//осуществлять подсчет свойств
		isObservables := strings.Contains(tmf.FieldBranch, "observables")
		countOne := strings.Count(tmf.FieldBranch, ".") <= 1

		if isObservables && countOne {
			var newFieldName = tmf.FieldName
			//сделал проверку на число что бы исключить повторение
			//имен для свойств являющихся срезами, так как в данной ситуации
			//имя содержащееся в tmf.FieldName представляет собой числовой
			//индекс, соответственное, если будет еще одно свойство являющееся
			//срезом, то может быть совпадение имен и изменение seqNumObservable, а как
			//результат будет переход на другой объект 'observables'
			if patterIsNum.MatchString(tmf.FieldName) {
				tmp := strings.Split(tmf.FieldBranch, ".")
				var nameTmp string
				if len(tmp) > 0 {
					nameTmp = tmp[len(tmp)-1] + "_"
				}

				newFieldName = nameTmp + tmf.FieldName
			}

			//подсчет свойств для объектов типа 'observables' выполняется для
			//того что бы отделить один объект 'observables' от другого
			if svn.GetValueName(newFieldName) {
				svn.CleanValueName()
				seqNumObservable++
			}

			svn.SetValueName(newFieldName)

			// ************* обработка правил EXCLUSION (исключения) *************
			if addrRule, isEqual := listRule.ExcludeRuleHandler(tmf.FieldBranch, newValue); isEqual {
				supportiveListExcludeRule.Add(seqNumObservable, tmf.FieldBranch, newValue, addrRule, isEqual)
			}
		}

		//обрабатываем свойство observables.attachment
		if strings.Contains(tmf.FieldBranch, "attachment") {
			listAttributeTmp.AddAttribute(tmf.FieldBranch, newValue, seqNumObservable)
		}

		//получаем rootId
		if tmf.FieldBranch == "event.object.rootId" {
			if rid, ok := newValue.(string); ok {
				rootId = rid
			}
		}

		//обрабатываем свойство event.object.tags, оно ответственно за
		//наполнение поля "Теги" MISP
		if tmf.FieldBranch == "event.object.tags" {
			if tag, ok := newValue.(string); ok {
				leot.SetTag(tag)
			}
		}

		//заполняем временный объект listGalaxyTags данными, предназначенными
		//для формирования специализированных тегов на основе которых в MISP
		//будут формироватся галактики
		addFuncGalaxyTags(tmf.FieldBranch, newValue)

		//обрабатываем свойство observables.tags
		if tmf.FieldBranch == "observables.tags" {
			listTags = handlerObservablesTags(newValue, listTags, listAttributesMisp, seqNumObservable)
		}

		//проверяем есть ли путь до обрабатываемого свойства в списке обработчиков
		//это главный обработчик который выполняет всю работу, особенно если выше идущие
		//обработчики не выполнялись так как принимаемые значения не соответствовали
		//параметрам для их вызова
		lf, ok := listHandlerMisp[tmf.FieldBranch]
		if ok {
			//основной обработчик путей из tmf.FieldBranch
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
		// ***********************************
		// Это логирование только для теста!!!
		// ***********************************
		logger.Send("testing", fmt.Sprintf("TEST_INFO func 'NewMispFormat', case with id '%d' does not comply with the rules", int(caseId)))
		//
		//

		logger.Send("warning", fmt.Sprintf("'the message with case id %d was not sent to MISP because it does not comply with the rules'", int(caseId)))
	} else {
		//добавляем case id в поле Info
		eventsMisp.Info += fmt.Sprintf(" :::TheHive case id '%d':::", int(caseId))

		//добавляем в datemodels.ListObjectEventTags дополнительные теги
		//ответственные за формирование галактик в MISP
		joinEventTags(leot, createGalaxyTags(listGalaxyTags))

		for k := range listAttributesMisp.GetList() {
			if supportiveListExcludeRule.CheckRuleTrue(k) {
				//удаляем элементы подходящие под правила группы EXCLUDE
				listAttributesMisp.DelElementList(k)
			}
		}

		// ***********************************
		// Это логирование только для теста!!!
		// ***********************************
		logger.Send("testing", fmt.Sprintf("TEST_INFO func 'NewMispFormat', case with id '%d' equal rules, send data to MISP module", int(caseId)))
		//
		//

		//тут отправляем сформированные по формату MISP пользовательские структуры
		mispModule.SendingDataInput(mispapi.InputSettings{
			Command:    "add event",
			TaskId:     taskId,
			CaseId:     caseId,
			RootId:     rootId,
			CaseSource: caseSource,
			UserEmail:  userEmail,
			MajorData: map[string]interface{}{
				"events": eventsMisp,
				//getNewListAttributes влияет на поля Category и Type типа Attributes
				"attributes": getNewListAttributes(
					listAttributesMisp.GetList(),
					listTags),
				"objects": getNewListObjects(
					listObjectsMisp.GetList(),
					listAttributeTmp.GetListAttribute()),
				"event.object.tags": leot.GetListTags(),
			}})
	}

	//очищаем события, список аттрибутов и текущий email пользователя
	userEmail, caseSource = "", ""
	leot.CleanListTags()
	eventsMisp.CleanEventsMispFormat()
	listObjectsMisp.CleanList()
	listAttributeTmp.CleanAttribute()
	listAttributesMisp.CleanList()

	// ТОЛЬКО ДЛЯ ТЕСТОВ, что бы завершить гроутину вывода информации и логирования
	//при выполнения тестирования
	logger.Send("STOP TEST", "")
}
