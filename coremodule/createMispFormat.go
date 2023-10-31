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

// storageValueName временное хранилище свойств элементов observables
type storageValueName []string

func NewStorageValueName() *storageValueName {
	return &storageValueName{}
}

func (svn *storageValueName) SetValueName(value string) {
	*svn = append(*svn, value)
}

func (svn *storageValueName) GetValueName(value string) bool {
	for _, v := range *svn {
		if v == value {
			return true
		}
	}

	return false
}

func (svn *storageValueName) CleanValueName() {
	*svn = storageValueName{}
}

var (
	eventsMisp         datamodels.EventsMispFormat
	listObjectsMisp    *datamodels.ListObjectsMispFormat
	listAttributeTmp   *datamodels.ListAttributeTmp
	listAttributesMisp *datamodels.ListAttributesMispFormat

	//		пока не нужны, временно отключаем
	//galaxyClustersMisp datamodels.GalaxyClustersMispFormat
	//galaxyElementMisp  datamodels.GalaxyElementMispFormat
	//usersMisp          datamodels.UsersMispFormat
	//organizationsMisp  datamodels.OrganisationsMispFormat
	//serversMisp        datamodels.ServersMispFormat
	//feedsMisp          datamodels.FeedsMispFormat
	//tagsMisp           datamodels.TagsMispFormat

	listHandlerMisp map[string][]func(interface{}, int)
)

func init() {
	eventsMisp = datamodels.NewEventMisp()
	listObjectsMisp = datamodels.NewListObjectsMispFormat()
	listAttributeTmp = datamodels.NewListAttributeTmp()
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

	listHandlerMisp = map[string][]func(interface{}, int){
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
		"observables._id":        {listAttributesMisp.SetValueObjectIdAttributesMisp},
		"observables.data":       {listAttributesMisp.SetValueValueAttributesMisp},
		"observables.dataType":   {listObjectsMisp.SetValueNameObjectsMisp},
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
}

func NewMispFormat(
	mispmodule *mispinteractions.ModuleMISP,
	logging chan<- datamodels.MessageLogging) (chan ChanInputCreateMispFormat, chan bool) {

	//канал принимающий данные необходимые для заполнения MISP форматов
	chanInput := make(chan ChanInputCreateMispFormat)
	//останавливает обработчик канала chanInput (при завершении декодировании сообщения)
	chanDone := make(chan bool)

	go func() {
		var (
			maxCountObservables, seqNum int
			userEmail                   string
			caseId                      float64
		)
		defer func() {
			close(chanInput)
			close(chanDone)
		}()

		svn := NewStorageValueName()
		leot := datamodels.NewListEventObjectTags()

		for key := range listHandlerMisp {
			if strings.Contains(key, "observables") {
				maxCountObservables++
			}
		}

		listTags := make(map[int][2]string)

		for {
			select {
			case tmf := <-chanInput:
				//ищем id события
				if tmf.FieldBranch == "event.object.caseId" {
					if cid, ok := tmf.Value.(float64); ok {
						caseId = cid
					}
				}

				// ищем email владельца события
				if tmf.FieldBranch == "event.object.owner" {
					if email, ok := tmf.Value.(string); ok {
						userEmail = email
					}
				}

				if strings.Contains(tmf.FieldBranch, "observables") && !strings.Contains(tmf.FieldBranch, "attachment") {
					if svn.GetValueName(tmf.FieldName) {
						svn.CleanValueName()
						seqNum++
					}

					svn.SetValueName(tmf.FieldName)
				}

				//обрабатываем свойство observables.attachment
				if strings.Contains(tmf.FieldBranch, "attachment") {
					listAttributeTmp.AddAttribute(tmf.FieldBranch, tmf.Value, seqNum)
				}

				//обрабатываем свойство event.object.tags, оно ответственно за
				//наполнение поля "Теги" MISP
				if tmf.FieldBranch == "event.object.tags" {
					if tag, ok := tmf.Value.(string); ok {
						leot.SetTag(tag)
					}
				}

				//обрабатываем свойство observables.tags
				if tmf.FieldBranch == "observables.tags" {
					if tag, ok := tmf.Value.(string); ok {
						result, err := HandlingObservablesTag(tag)
						if err == nil {
							listTags[seqNum] = result
						}
					}
				}

				//проверяем есть ли путь до обрабатываемого свойства в списке обработчиков
				lf, ok := listHandlerMisp[tmf.FieldBranch]
				if ok {
					//основной обработчик путей из tmf.FieldBranch
					for _, f := range lf {
						f(tmf.Value, seqNum)
					}
				}

			case isAllowed := <-chanDone:
				if !isAllowed {
					_, f, l, _ := runtime.Caller(0)

					logging <- datamodels.MessageLogging{
						MsgData: fmt.Sprintf("'the message with case id %d was not sent to MISP because it does not comply with the rules' %s:%d", int(caseId), f, l-1),
						MsgType: "warning",
					}
				} else {
					//добавляем case id в поле Info
					eventsMisp.Info += fmt.Sprintf(" :::TheHive case id '%d':::", int(caseId))

					//тут отправляем сформированные по формату MISP пользовательские структуры
					mispmodule.SendingDataInput(mispinteractions.SettingsChanInputMISP{
						Command:   "add event",
						CaseId:    caseId,
						UserEmail: userEmail,
						MajorData: map[string]interface{}{
							"events": eventsMisp,
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
				userEmail = ""
				leot.CleanListTags()
				eventsMisp.CleanEventsMispFormat()
				listObjectsMisp.CleanListObjectsMisp()
				listAttributeTmp.CleanAttribute()
				listAttributesMisp.CleanListAttributesMisp()

				return
			}
		}
	}()

	return chanInput, chanDone
}

func getNewListAttributes(al map[int]datamodels.AttributesMispFormat, lat map[int][2]string) []datamodels.AttributesMispFormat {
	nal := make([]datamodels.AttributesMispFormat, 0, len(al))

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

func getNewListObjects(
	listObjects map[int]datamodels.ObjectsMispFormat,
	attachment map[int][]datamodels.AttributeMispFormat,
) map[int]datamodels.ObjectsMispFormat {
	nlo := make(map[int]datamodels.ObjectsMispFormat, len(attachment))
	for k, v := range attachment {
		if obj, ok := listObjects[k]; ok {
			obj.Attribute = v
			nlo[k] = obj
		}
	}

	return nlo
}

func HandlingObservablesTag(tag string) ([2]string, error) {
	nl := [2]string{}
	patter := regexp.MustCompile(`^misp:([\w\-].*)=\"([\w\-].*)\"$`)

	if !patter.MatchString(tag) {
		return nl, fmt.Errorf("the accepted value does not match the regular expression")
	}

	result := patter.FindAllStringSubmatch(tag, -1)

	if len(result) > 0 && len(result[0]) == 3 {
		nl = [2]string{result[0][1], result[0][2]}
	}

	return nl, nil
}
