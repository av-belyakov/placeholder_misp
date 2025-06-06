package coremodule

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/av-belyakov/objectsmispformat"
	"github.com/av-belyakov/placeholder_misp/cmd/mispapi"
)

// Start создаёт набор объектов в формате MISP
func (gen *GenerateObjectsFormatMISP) Start(chDecodeJSON <-chan ChanInputCreateMispFormat, taskId string) {
	go func() {
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
		defer listObjectsMisp.CleanList()

		listAttributeTmp := objectsmispformat.NewListAttributeTmp()
		defer listAttributeTmp.CleanAttribute()

		listAttributesMisp := objectsmispformat.NewListAttributesMispFormat()
		defer listAttributesMisp.CleanList()

		leot := objectsmispformat.NewListEventObjectTags()
		defer leot.CleanListTags()

		supportiveListExcludeRule := NewSupportiveListExcludeRuleTmp(gen.listRule.GetRuleExclude())

		//это основной обработчик параметров входящего объекта
		listHandlerMisp := map[string][]func(any, int){
			//event -> events
			"event.object.title":     {eventsMisp.SetAnyInfo},
			"event.object.startDate": {eventsMisp.SetAnyTimestamp},
			"event.details.endDate":  {eventsMisp.SetAnyDate},
			"event.object.tlp":       {eventsMisp.SetAnyDistribution},
			"event.object.severity":  {eventsMisp.SetAnyThreatLevelId},
			"event.organisationId":   {eventsMisp.SetAnyOrgId},
			"event.object.updatedAt": {eventsMisp.SetAnySightingTimestamp},
			"event.object.owner":     {eventsMisp.SetAnyEventCreatorEmail},
			"event.object.customFields.class-attack.string": {func(i any, num int) {
				leot.SetTag(fmt.Sprintf("class-attack=\"%v\"", i))
			}},
			//observables -> attributes
			"observables._id": {
				listAttributesMisp.SetAnyObjectId,
				listObjectsMisp.SetValueId,
			},
			"observables.data": {listAttributesMisp.SetAnyValue},
			"observables.dataType": {
				listObjectsMisp.SetValueName,
				//здесь выполняем автоподстановку значений для полей Type и Category объекта
				//AttributesMisp на основе некоторой логикиб а так же уже предустановленных
				//значений, при этом значения, заданные пользователем для этих полей, обрабатываются
				//отдельно и хранятся в listTags, а после закрытия канала совмещаются с
				//объектами AttributesMisp и следовательно перезаписывают значения выполненные
				//через автоподстановку
				listAttributesMisp.HandlingAnyDataType,
			},
			"observables._createdAt": {listAttributesMisp.SetAnyTimestamp},
			"observables.message":    {listAttributesMisp.SetAnyComment},
			"observables.startDate": {
				listAttributesMisp.SetAnyFirstSeen,
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
			newValue, _, err := gen.listRule.ReplacementRuleHandler(msg.ValueType, msg.FieldBranch, msg.Value)
			if err != nil {
				gen.logger.Send("warning", fmt.Sprintf("search value '%s' from rule of section 'REPLACE' is not fulfilled", msg.Value))
			}
			//обработка правил PASS (пропуск)
			gen.listRule.PassRuleHandler(msg.FieldBranch, newValue)
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
				if addrRule, isEqual := gen.listRule.ExcludeRuleHandler(msg.FieldBranch, newValue); isEqual {
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
		if gen.listRule.GetRulePassany() || gen.listRule.SomePassRuleIsTrue() {
			isAllowed = true

			//сетчик кейсов соответствующих или не соответствующих правилам
			gen.counter.SendMessage("update events meet rules", 1)
		}

		//удаляем те объекты Attributes которые соответствуют правилам EXCLUDE
		delElementAttributes(exclusionRules, listAttributesMisp, gen.logger)

		gen.listRule.CleanStatementExpressionRulePass()

		if caseId == 0 {
			if len(eventsMisp.EventCreatorEmail) != 0 {
				gen.logger.Send("error", fmt.Sprintf("the caseId in the event cannot be equal to 0 (event creator email '%s')", eventsMisp.EventCreatorEmail))
			}

			return
		}

		if !isAllowed {
			gen.logger.Send("warning", fmt.Sprintf("the message with case id %d was not sent to MISP because it does not comply with the rules", int(caseId)))

			return
		}

		//добавляем case id в поле Info
		eventsMisp.Info += fmt.Sprintf(" :::TheHive caseId:'%d':::", int(caseId))

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

		gen.logger.Send("info", fmt.Sprintf("the case with id:'%d' complies with the specified rules and has been submitted for further processing", int(caseId)))

		//тут отправляем сформированные по формату MISP пользовательские структуры
		gen.mispModule.SendDataInput(mispapi.InputSettings{
			Command:    "add event",
			TaskId:     taskId,
			CaseId:     caseId,
			RootId:     rootId,
			UserEmail:  userEmail,
			CaseSource: caseSource,
			Data:       *mispFormat,
		})
	}()
}
