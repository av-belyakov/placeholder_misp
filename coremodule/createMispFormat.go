package coremodule

import (
	"fmt"
	"regexp"
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
		//"observables.tags":       {listAttributesMisp.HandlingValueTagsAttributesMisp},
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
+ 2. Заменить objectId на caseId



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
			currentCount, maxCountObservables, seqNum int
			userEmail, observableId                   string
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
		listTags := make(map[int][2]string)

		for {
			select {
			case tmf := <-chanInput:
				if tmf.FieldBranch == "event.object.caseId" {
					if cid, ok := tmf.Value.(float64); ok {
						caseId = cid
					}
				}

				if tmf.FieldBranch == "observables._id" {
					if id, ok := tmf.Value.(string); ok {
						if observableId == "" {
							observableId = id
						}

						if observableId != id {
							seqNum++
							fmt.Printf(":::::::: observableId '%s' == '%s' id, NUM: %d\n", observableId, id, seqNum)

							observableId = id
						}
					}
				}

				//обрабатываем свойство observables.tags
				if tmf.FieldBranch == "observables.tags" {
					fmt.Println("========== func 'NewMispFormat' =========== observables.tags, tags:", tmf.Value)

					if tag, ok := tmf.Value.(string); ok {
						result, err := HandlingListTags(tag)
						if err == nil {
							listTags[seqNum] = result

							fmt.Printf("|||||| %d. add new tag: '%s'\n", seqNum, result)
						}
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

				//основной обработчик путей из tmf.FieldBranch
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

					fmt.Println("+++++++++++++++___________________")

					isNew = true
				}

			case isAllowed := <-chanDone:

				fmt.Printf("\n\tfunc 'NewMispFormat', RESEIVED chanDone, eventsMisp: %v, isAllowed: %v\nLIST TAGS: '%v'\n", eventsMisp, isAllowed, listTags)

				if !isAllowed {
					_, f, l, _ := runtime.Caller(0)

					loging <- datamodels.MessageLoging{
						MsgData: fmt.Sprintf(" 'the message with %d was not sent to MISP because it does not comply with the rules' %s:%d", int(caseId), f, l-1),
						MsgType: "warning",
					}
				} else {
					attr := getNewListAttributes(
						listAttributesMisp.GetListAttributesMisp(),
						listTags,
						/*listAttributesMisp.GetListAttributeTags()*/)

					fmt.Println("----------- NEW ATTRIBUTES -----------")
					for k, v := range attr {
						fmt.Printf("%d. \n\tData: '%s'\n\tCategory: '%s'\n\tType: '%s'\n", k, v.Value, v.Category, v.Type)
					}

					//fmt.Printf("\n\tfunc 'NewMispFormat', NEW ATTRIBUTESSSSSSSS: %v\n", attr)

					//тут отправляем сформированные по формату MISP пользовательские структуры
					mispmodule.SendingDataInput(mispinteractions.SettingsChanInputMISP{
						CaseId:    caseId,
						UserEmail: userEmail,
						//AttributeTags: listAttributesMisp.GetListAttributeTags(),
						MajorData: map[string]interface{}{
							"events": eventsMisp,
							//"attributes": listAttributesMisp.GetListAttributesMisp(),
							"attributes": attr,
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

func getNewListAttributes(al []datamodels.AttributesMispFormat, lat map[int][2]string) []datamodels.AttributesMispFormat {
	countAttr := len(al)

	nal := make([]datamodels.AttributesMispFormat, 0, countAttr)

	for k, v := range al {
		if elem, ok := lat[k]; ok {
			v.Category = elem[0]
			v.Type = elem[1]

			nal = append(nal, v)

			continue
		}

		nal = append(nal, v)
	}

	return nal
}

/*func getNewListAttributes(al []datamodels.AttributesMispFormat, lat map[int][][2]string) []datamodels.AttributesMispFormat {
	countAttr := len(al)

	nal := make([]datamodels.AttributesMispFormat, 0, countAttr)

	for k, v := range al {
		if elem, ok := lat[k]; ok {
			for _, value := range elem {
				v.Category = value[0]
				v.Type = value[1]

				nal = append(nal, v)
			}

			continue
		}

		nal = append(nal, v)
	}

	fmt.Println("________ ___________ func 'getNewListAttributes', NEW LIST ATTRIBUTES = ", nal)

	return nal
}*/

func HandlingListTags(tag string) ([2]string, error) {
	nl := [2]string{}
	patter := regexp.MustCompile(`^misp:([\w\-].*)=\"([\w\-].*)\"$`)

	if !patter.MatchString(tag) {
		return nl, fmt.Errorf("the accepted value does not match the regular expression")
	}

	result := patter.FindAllStringSubmatch(tag, -1)

	fmt.Println("func 'HandlingListTags', ======= RESULT REGEXP = ", result)

	if len(result) > 0 && len(result[0]) == 3 {
		nl = [2]string{result[0][1], result[0][2]}
	}

	return nl, nil
}

/*func HandlingListTags(l []string) [][2]string {
	nl := make([][2]string, 0, len(l))
	patter := regexp.MustCompile(`^misp:([\w\-].*)=\"([\w\-].*)\"$`)

	for _, v := range l {
		if !patter.MatchString(v) {
			continue
		}

		result := patter.FindAllStringSubmatch(v, -1)
		if len(result) > 0 && len(result[0]) == 3 {
			nl = append(nl, [2]string{result[0][1], result[0][2]})
		}
	}

	return nl
}*/
