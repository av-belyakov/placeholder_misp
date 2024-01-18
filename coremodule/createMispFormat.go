package coremodule

import (
	"fmt"
	"regexp"
	"runtime"
	"strings"

	"placeholder_misp/datamodels"
	"placeholder_misp/mispinteractions"
)

func NewMispFormat(
	taskId string,
	mispmodule *mispinteractions.ModuleMISP,
	logging chan<- datamodels.MessageLogging) (chan ChanInputCreateMispFormat, chan bool) {

	//канал принимающий данные необходимые для заполнения MISP форматов
	chanInput := make(chan ChanInputCreateMispFormat)
	//останавливает обработчик канала chanInput (при завершении декодировании сообщения)
	chanDone := make(chan bool)

	eventsMisp := datamodels.NewEventMisp()
	listObjectsMisp := datamodels.NewListObjectsMispFormat()
	listAttributeTmp := datamodels.NewListAttributeTmp()
	listAttributesMisp := datamodels.NewListAttributesMispFormat()

	/*galaxyClustersMisp = datamodels.GalaxyClustersMispFormat{
		Description:   "3",
		GalaxyElement: []datamodels.GalaxyElementMispFormat{},
	}
	usersMisp = datamodels.UsersMispFormat{
		Newsread:     "0",
		ChangePw:     "0",
		CurrentLogin: "0",
		LastLogin:    "0",
		DateCreated:  "0",
		DateModified: "0",
	}
	organizationsMisp = datamodels.OrganisationsMispFormat{
		DateCreated:  "0",
		DateModified: "0",
	}
	serversMisp = datamodels.ServersMispFormat{}
	feedsMisp = datamodels.FeedsMispFormat{
		Distribution: "3",
		SourceFormat: "misp",
		InputSource:  "network",
	}
	tagsMisp = datamodels.TagsMispFormat{
		Exportable:     true,
		IsGalaxy:       true,
		IsCustomGalaxy: true,
		Inherited:      1,
	}*/

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

	go func() {
		var (
			maxCountObservables int
			seqNumObservable    int
			userEmail           string
			caseSource          string
			caseId              float64
			patterIsNum         *regexp.Regexp = regexp.MustCompile(`^\d+$`)
		)
		defer func() {
			close(chanInput)
			close(chanDone)
		}()

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

		for {
			select {
			case tmf := <-chanInput:
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

					if tmf.ExclusionRuleWorked {
						exclusionRules.Add(seqNumObservable, tmf.FieldBranch)
					}
				}

				//обрабатываем свойство observables.attachment
				if strings.Contains(tmf.FieldBranch, "attachment") {
					listAttributeTmp.AddAttribute(tmf.FieldBranch, tmf.Value, seqNumObservable)
				}

				//обрабатываем свойство event.object.tags, оно ответственно за
				//наполнение поля "Теги" MISP
				if tmf.FieldBranch == "event.object.tags" {

					/*
						ДЛЯ ПОДКЛЮЧЕНИЯ ГАЛАКТИК нужно сформировать, на основе данных из TTP,
						тег со следующем форматом. Пример ниже

						"misp-galaxy:mitre-attack-pattern=\"Match Legitimate Name or Location - T1036.005\""

						ttp.extraData.pattern.patternId - содержит что то похожее на T1036.005
						ttp.extraData.pattern.patternType - содержит что то похожее на attack-pattern
						ttp.extraData.pattern.name - вроде описание
					*/

					if tag, ok := tmf.Value.(string); ok {
						leot.SetTag(tag)
					}
				}

				//заполняем временный объект listGalaxyTags данными, предназначенными
				//для формирования специализированных тегов на основе которых в MISP
				//будут формироватся галактики
				addFuncGalaxyTags(tmf.FieldBranch, tmf.Value)

				//обрабатываем свойство observables.tags
				if tmf.FieldBranch == "observables.tags" {
					listTags = handlerObservablesTags(tmf.Value, listTags, listAttributesMisp, seqNumObservable)
				}

				//проверяем есть ли путь до обрабатываемого свойства в списке обработчиков
				lf, ok := listHandlerMisp[tmf.FieldBranch]
				if ok {
					//основной обработчик путей из tmf.FieldBranch
					for _, f := range lf {
						f(tmf.Value, seqNumObservable)
					}
				}

			case isAllowed := <-chanDone:
				//удаляем те объекты Attributes которые соответствуют правилам EXCLUDE
				delElementAttributes(exclusionRules, listAttributesMisp, logging)

				//fmt.Println("==================== exclusionRules =====================")
				//for k, v := range exclusionRules.SearchObjectName("observables") {
				//	fmt.Printf("%d. %d\n", k+1, v.SequenceNumber)
				//}

				//fmt.Println("____________________ Attributes _______________________")
				//i := 1
				//for k, v := range getNewListAttributes(
				//	listAttributesMisp.GetListAttributesMisp(),
				//	listTags) {
				//	fmt.Printf("%d. num: %d, Value: %s\n", i, k, v.Value)
				//	i++
				//}
				//fmt.Println("_______________________________________________________")

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

					//fmt.Println("==================== The Galaxise =======================")
					//fmt.Println(listGalaxyTags.Get())
					//fmt.Println(leot.GetListTags())
					//fmt.Println(createGalaxyTags(listGalaxyTags))
					//fmt.Println("===========================================")

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

				// только для тестов, что бы завершить гроутину вывода информации и логирования
				logging <- datamodels.MessageLogging{
					MsgData: "",
					MsgType: "STOP TEST",
				}

				return
			}
		}
	}()

	return chanInput, chanDone
}
