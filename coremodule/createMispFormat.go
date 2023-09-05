package coremodule

import (
	"fmt"
	"placeholder_misp/datamodels"
	"placeholder_misp/mispinteractions"
	"runtime"
	"strings"
)

type ChanInputCreateMispFormat struct {
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

	//пока не нужны, временно отключаем
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
	eventsMisp = datamodels.EventsMispFormat{
		Analysis:          getAnalysis(),
		Distribution:      getDistributionEvent(),
		Timestamp:         "0",
		ThreatLevelId:     getThreatLevelId(),
		PublishTimestamp:  "0",
		SightingTimestamp: "0",
		SharingGroupId:    getSharingGroupId(),
	}
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

func NewMispFormat(
	uuidTask string,
	mispmodule *mispinteractions.ModuleMISP,
	loging chan<- datamodels.MessageLoging) (chan ChanInputCreateMispFormat, chan bool) {

	//канал принимающий данные необходимые для заполнения MISP форматов
	chanInput := make(chan ChanInputCreateMispFormat)
	//останавливает обработчик канала chanInput (при завершении декодировании сообщения)
	chanDone := make(chan bool)

	go func() {
		var (
			currentCount, maxCountObservables int
			userEmail                         string
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

		fmt.Println("maxCountObservables = ", maxCountObservables)

		isNew := true

		for {
			select {
			case tmf := <-chanInput:
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
				if !isAllowed {
					_, f, l, _ := runtime.Caller(0)

					loging <- datamodels.MessageLoging{
						MsgData: fmt.Sprintf("the message with %s was not sent to MISP because it does not comply with the rules %s:%d", uuidTask, f, l-2),
						MsgType: "warning",
					}
				} else {
					//тут отправляем сформированные по формату MISP пользовательские структуры
					mispmodule.SendingDataInput(mispinteractions.SettingsChanInputMISP{
						UserEmail: userEmail,
						MajorData: map[string]interface{}{
							"events":     eventsMisp,
							"attributes": listAttributesMisp.GetListAttributesMisp(),
						}})

					userEmail = ""
				}

				return
			}
		}
	}()

	return chanInput, chanDone
}

func getAnalysis() string {
	return "2"
}

func getDistributionEvent() string {
	return "3"
}

func getThreatLevelId() string {
	return "4"
}

func getSharingGroupId() string {
	return "1"
}

func getToIds() bool {
	return true
}

/*func getTagTLP(tlp int) datamodels.TagsMispFormat {
	tag := datamodels.TagsMispFormat{Name: "tlp:red", Colour: "#cc0033"}

	switch tlp {
	case 0:
		tag.Name = "tlp:white"
		tag.Colour = "#ffffff"
	case 1:
		tag.Name = "tlp:green"
		tag.Colour = "#339900"
	case 2:
		tag.Name = "tlp:amber"
		tag.Colour = "#ffc000"
	}

	return tag
}*/
