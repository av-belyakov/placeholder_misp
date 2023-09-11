package coremodule

import (
	"fmt"
	"runtime"
	"strings"

	"placeholder_misp/datamodels"
	"placeholder_misp/mispinteractions"
)

type ChanInputCreateMispFormat struct {
	UUID        string
	FieldName   string
	ValueType   string
	Value       interface{}
	FieldBranch string
}

type FieldsNameMapping struct {
	InputFieldName, MispFieldName string
}

var (
	eventsMisp         datamodels.EventsMispFormat
	listAttributesMisp *datamodels.ListAttributesMispFormat

	//		пока не нужны, временно отключаем
	//galaxyClustersMisp datamodels.GalaxyClustersMispFormat
	//galaxyElementMisp  datamodels.GalaxyElementMispFormat
	//usersMisp          datamodels.UsersMispFormat
	//organizationsMisp  datamodels.OrganisationsMispFormat
	//serversMisp        datamodels.ServersMispFormat
	//feedsMisp          datamodels.FeedsMispFormat
	//tagsMisp           datamodels.TagsMispFormat

	listHandlerMisp map[string][]func(interface{}, bool)
)

func init() {
	eventsMisp = datamodels.NewEventMisp()
	listAttributesMisp = datamodels.NewListAttributesMispFormat()

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

	listHandlerMisp = map[string][]func(interface{}, bool){
		//events
		"event.object.title":     {eventsMisp.SetValueInfoEventsMisp},
		"event.object.startDate": {eventsMisp.SetValueTimestampEventsMisp},
		"event.details.endDate":  {eventsMisp.SetValueDateEventsMisp},
		"event.object.tlp":       {eventsMisp.SetValueDistributionEventsMisp},
		"event.object.severity":  {eventsMisp.SetValueThreatLevelIdEventsMisp},
		"event.organisationId":   {eventsMisp.SetValueOrgIdEventsMisp},
		"event.object.updatedAt": {eventsMisp.SetValueSightingTimestampEventsMisp},
		"event.object.owner":     {eventsMisp.SetValueEventCreatorEmailEventsMisp},
		//attributes
		"observables._id":        {listAttributesMisp.SetValueObjectIdAttributesMisp},
		"observables.data":       {listAttributesMisp.SetValueValueAttributesMisp},
		"observables._createdAt": {listAttributesMisp.SetValueTimestampAttributesMisp},
		"observables.message":    {listAttributesMisp.SetValueCommentAttributesMisp},
		"observables.startDate":  {listAttributesMisp.SetValueFirstSeenAttributesMisp},
	}
}

/*
Для типов Attributes нужно сделать следующее:
1. свойства MISP формата такие как Category и Type заполнять на основе значений
 поля observables[*].tags из формата Hive, где из строки "misp:Network activity="email-src""
 Network activity идет в Category, а email-src в Type, если observables[*].tags пустое,
 то ставим "other" для Category и Type. Однако если observables[*].tags имеет
 болше чем одно значение то тогда нужно создать то количество Attributes с одними
 и теме же данными (поле data, время создание), но с разными Category и Type.
2. Заменить objectId на caseId



*/

func NewMispFormat(
	mispmodule *mispinteractions.ModuleMISP,
	loging chan<- datamodels.MessageLoging) (chan ChanInputCreateMispFormat, chan bool) {

	//канал принимающий данные необходимые для заполнения MISP форматов
	chanInput := make(chan ChanInputCreateMispFormat)
	//останавливает обработчик канала chanInput (при завершении декодировании сообщения)
	chanDone := make(chan bool)

	fmt.Printf("\n\n\tfunc 'NewMispFormat', START...")

	go func() {
		var (
			currentCount, maxCountObservables int
			userEmail                         string
			//objectId               string
			caseId float64
		)
		defer func() {
			close(chanInput)
			close(chanDone)
		}()

		for key := range listHandlerMisp {
			if strings.Contains(key, "observables") {
				maxCountObservables++
			}
		}

		isNew := true

		for {
			select {
			case tmf := <-chanInput:
				//ищем id обрабатываемого события
				/*if tmf.FieldBranch == "event.objectId" || tmf.FieldBranch == "event.object.id" {

					fmt.Println("===================== event.objectId 1111")

					if oid, ok := tmf.Value.(string); ok {
						fmt.Println("===================== tmf.Value.(string) 222 oid = ", oid)

						objectId = oid
					}
				}*/
				if tmf.FieldBranch == "event.object.caseId" {

					fmt.Println("===================== event.objectId 1111")
					if cid, ok := tmf.Value.(float64); ok {
						fmt.Println("===================== tmf.Value.(float64) 222 cid = ", cid)

						caseId = cid
					}
				}

				//проверяем есть ли путь до обрабатываемого свойства в списке обработчиков
				lf, ok := listHandlerMisp[tmf.FieldBranch]
				if !ok {
					continue
				}

				if currentCount > 0 {
					isNew = false
				}

				for _, f := range lf {
					f(tmf.Value, isNew)
				}

				// ищем email владельца события
				if tmf.FieldBranch == "event.object.owner" {
					if email, ok := tmf.Value.(string); ok {
						userEmail = email
					}
				}

				if strings.Contains(tmf.FieldBranch, "observables") {
					currentCount++
				}

				if currentCount == maxCountObservables {
					currentCount = 0
					isNew = true
				}

			case isAllowed := <-chanDone:

				fmt.Printf("\n\tfunc 'NewMispFormat', RESEIVED chanDone, eventsMisp: %v, isAllowed: %v\n", eventsMisp, isAllowed)

				if !isAllowed {
					_, f, l, _ := runtime.Caller(0)

					loging <- datamodels.MessageLoging{
						MsgData: fmt.Sprintf(" 'the message with %d was not sent to MISP because it does not comply with the rules' %s:%d", int(caseId), f, l-1),
						MsgType: "warning",
					}
				} else {
					//тут отправляем сформированные по формату MISP пользовательские структуры
					mispmodule.SendingDataInput(mispinteractions.SettingsChanInputMISP{
						//ObjectId:  objectId,
						CaseId:    caseId,
						UserEmail: userEmail,
						MajorData: map[string]interface{}{
							"events":     eventsMisp,
							"attributes": listAttributesMisp.GetListAttributesMisp(),
						}})
				}

				//очищаем события, список аттрибутов и текущий email пользователя
				userEmail = ""
				eventsMisp.CleanEventsMispFormat()
				listAttributesMisp.CleanListAttributesMisp()

				return
			}
		}
	}()

	return chanInput, chanDone
}
