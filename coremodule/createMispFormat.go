package coremodule

import (
	"fmt"
	"regexp"
	"runtime"
	"strings"

	"placeholder_misp/datamodels"
	"placeholder_misp/mispinteractions"
	rules "placeholder_misp/rulesinteraction"
)

func NewMispFormat(
	chanOutputDecodeJson <-chan ChanInputCreateMispFormat,
	taskId string,
	mispmodule *mispinteractions.ModuleMISP,
	listRule *rules.ListRule,
	logging chan<- datamodels.MessageLogging,
	counting chan<- datamodels.DataCounterSettings) {

	eventsMisp := datamodels.NewEventMisp()
	listObjectsMisp := datamodels.NewListObjectsMispFormat()
	listAttributeTmp := datamodels.NewListAttributeTmp()
	listAttributesMisp := datamodels.NewListAttributesMispFormat()

	listHandlerMisp := map[string][]func(interface{}, int){
		//event -> events
		"event.object.title":     {eventsMisp.SetValueInfoEventsMisp},
		"event.object.startDate": {eventsMisp.SetValueTimestampEventsMisp},
		"event.details.endDate":  {eventsMisp.SetValueDateEventsMisp},
		"event.object.tlp":       {eventsMisp.SetValueDistributionEventsMisp},
		"event.object.severity":  {eventsMisp.SetValueThreatLevelIdEventsMisp},
		"event.organisationId":   {eventsMisp.SetValueOrgIdEventsMisp},
		"event.object.updatedAt": {eventsMisp.SetValueSightingTimestampEventsMisp},
		"event.object.owner":     {eventsMisp.SetValueEventCreatorEmailEventsMisp},
		//observables -> attributes
		"observables._id":  {listAttributesMisp.SetValueObjectIdAttributesMisp},
		"observables.data": {listAttributesMisp.SetValueValueAttributesMisp},
		"observables.dataType": {
			listObjectsMisp.SetValueNameObjectsMisp,
			//здесь выполняем автоподстановку значений для полей Type и Category
			//объекта AttributesMisp на основе определенной логике и уже предустановленных
			//значений, при этом значения, заданные пользователем для этих полей, обрабатываются
			//отдельно и хранятся в listTags, а после закрытия канала совмещаются с
			//объектами AttributesMisp и следовательно перезаписывают значения выполненные
			//через автоподстановку
			listAttributesMisp.HandlingValueDataTypeAttributesMisp,
		},
		"observables._createdAt": {listAttributesMisp.SetValueTimestampAttributesMisp},
		"observables.message":    {listAttributesMisp.SetValueCommentAttributesMisp},
		"observables.startDate": {
			listAttributesMisp.SetValueFirstSeenAttributesMisp,
			listObjectsMisp.SetValueFirstSeenObjectsMisp,
			listObjectsMisp.SetValueTimestampObjectsMisp,
		},
		//observables.attachment -> objects
		"observables.attachment.size": {listObjectsMisp.SetValueSizeObjectsMisp},
	}

	var (
		maxCountObservables int
		seqNumObservable    int
		userEmail           string
		caseSource          string
		caseId              float64
		patterIsNum         *regexp.Regexp = regexp.MustCompile(`^\d+$`)
	)

	svn := NewStorageValueName()
	leot := datamodels.NewListEventObjectTags()
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

		//************* Обработка правил ***************
		//обработка правил Replacement
		newValue, _, err := listRule.ReplacementRuleHandler(tmf.ValueType, tmf.FieldBranch, tmf.Value)
		if err != nil {
			_, f, l, _ := runtime.Caller(0)

			logging <- datamodels.MessageLogging{
				MsgData: fmt.Sprintf("'search value \"%s\" from rule of section \"REPLACE\" is not fulfilled' %s:%d", tmf.Value, f, l-1),
				MsgType: "warning",
			}
		}
		//обработка правил Pass
		listRule.PassRuleHandler(tmf.FieldBranch, newValue)

		//fmt.Printf("func 'NewMispFormat', FieldBranch: '%s', Value: '%v' (newValue: '%v')\n", tmf.FieldBranch, tmf.Value, newValue)

		/*
			!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
				здесь и далее в других типах надо подумать как
				обрабатывать группу условий тип "И"

			!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		*/
		//обработка правил EXCLUSION (исключения)
		//ExclusionRuleWorked: lr.ExcludeRuleHandler(fieldBranch, ncv),
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

			/*
				!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
					здесь и далее в других типах надо подумать как
					обрабатывать группу условий тип "И"

				!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
			*/
			if tmf.ExclusionRuleWorked {
				exclusionRules.Add(seqNumObservable, tmf.FieldBranch)
			}
		}

		//обрабатываем свойство observables.attachment
		if strings.Contains(tmf.FieldBranch, "attachment") {
			listAttributeTmp.AddAttribute(tmf.FieldBranch, newValue, seqNumObservable)
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
		lf, ok := listHandlerMisp[tmf.FieldBranch]
		if ok {
			//основной обработчик путей из tmf.FieldBranch
			for _, f := range lf {
				f(newValue, seqNumObservable)
			}
		}
	}

	var isAllowed bool
	dt := "events do not meet rules"
	//проверяем что бы хотя бы одно правило разрешало пропуск кейса
	if listRule.GetRulePassany() || listRule.SomePassRuleIsTrue() {
		isAllowed = true
		dt = "update events meet rules"
	}

	//сетчик кейсов соответствующих или не соответствующих правилам
	counting <- datamodels.DataCounterSettings{
		DataType: dt,
		Count:    1,
	}

	//удаляем те объекты Attributes которые соответствуют правилам EXCLUDE
	delElementAttributes(exclusionRules, listAttributesMisp, logging)

	//выполняет очистку значения StatementExpression что равно отсутствию совпадений в правилах Pass
	// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	// ВРЕМЕННО ЗАКОМЕНТИРОВАЛ
	// не забыть снять коментарий
	//
	listRule.CleanStatementExpressionRulePass()

	if !isAllowed {
		_, f, l, _ := runtime.Caller(0)
		logging <- datamodels.MessageLogging{
			MsgData: fmt.Sprintf("'the message with case id %d was not sent to MISP because it does not comply with the rules' %s:%d", int(caseId), f, l-1),
			MsgType: "warning",
		}
	} else {
		//добавляем case id в поле Info
		eventsMisp.Info += fmt.Sprintf(" :::TheHive case id '%d':::", int(caseId))

		//добавляем в datemodels.ListObjectEventTags дополнительные теги
		//ответственные за формирование галактик в MISP
		joinEventTags(leot, createGalaxyTags(listGalaxyTags))

		//тут отправляем сформированные по формату MISP пользовательские структуры
		mispmodule.SendingDataInput(mispinteractions.SettingsChanInputMISP{
			Command:    "add event",
			TaskId:     taskId,
			CaseId:     caseId,
			CaseSource: caseSource,
			UserEmail:  userEmail,
			MajorData: map[string]interface{}{
				"events": eventsMisp,
				//getNewListAttributes влияет на поля Category и Type
				//типа Attributes
				"attributes": getNewListAttributes(
					listAttributesMisp.GetListAttributesMisp(),
					listTags),
				"objects": getNewListObjects(
					listObjectsMisp.GetListObjectsMisp(),
					listAttributeTmp.GetListAttribute()),
				"event.object.tags": leot.GetListTags(),
			}})
	}

	//очищаем события, список аттрибутов и текущий email пользователя
	userEmail, caseSource = "", ""
	leot.CleanListTags()
	eventsMisp.CleanEventsMispFormat()
	//exclusionRules.Clean()
	listObjectsMisp.CleanListObjectsMisp()
	listAttributeTmp.CleanAttribute()
	listAttributesMisp.CleanListAttributesMisp()

	// ТОЛЬКО ДЛЯ ТЕСТОВ, что бы завершить гроутину вывода информации и логирования
	//при выполнения тестирования
	logging <- datamodels.MessageLogging{
		MsgData: "",
		MsgType: "STOP TEST",
	}
}
