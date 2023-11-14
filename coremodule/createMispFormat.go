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
	taskId string,
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
			caseSource                  string
			caseId                      float64
			patterIsNum                 *regexp.Regexp = regexp.MustCompile(`^\d+$`)
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

				//if isObservables && countOne {
				//	fmt.Println("00000000000000000000 tmf.FieldBranch = ", tmf.FieldBranch)
				//}

				//isAttachment := strings.Contains(tmf.FieldBranch, "attachment")
				if isObservables && countOne /*!isAttachment*/ {
					var newFieldName = tmf.FieldName
					fmt.Println("BEFORE ++++++++ tmf.FieldName:", tmf.FieldBranch, " seqNum = ", seqNum)

					//сделал проверку на число что бы исключить повторение
					//имен для свойств являющихся срезами, так как в данной ситуации
					//имя содержащееся в tmf.FieldName представляет собой числовой
					//индекс, соответственное, если будет еще одно свойство являющееся
					//срезом, то может быть совпадение имен и изменение seqNum, а как
					//результат перескакиванее на другой объект 'observables'
					//if strings.Contains(tmf.FieldBranch, ".tags") {
					if patterIsNum.MatchString(tmf.FieldName) {
						fmt.Println(")))))))))))))))))))))))))))))))))))))))")
						/*

							Никак с регуляркой не справлюсь
							не работает этот кусок кода похоже что это из-за
							перевода строки в tmf.FieldName

							кроме того имя должно формироваться из последнего
							после точки значения tmf.FieldBranch, так будет красивше

						*/

						newFieldName = tmf.FieldName + "_" + tmf.FieldName
					} else {
						fmt.Println("&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&")
					}

					//подсчет свойств для объектов типа 'observables' выполняется для
					//того что бы отделить один объект 'observables' от другого
					if svn.GetValueName(newFieldName) {
						svn.CleanValueName()
						seqNum++
					}

					fmt.Println("AFTER ++++++++ tmf.FieldName:", tmf.FieldBranch, " seqNum = ", seqNum)
					fmt.Printf("newFieldName: '%s'\n", newFieldName)

					svn.SetValueName(newFieldName)
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
				fmt.Println("_________________ listTags: __________________")
				for k, v := range listTags {
					fmt.Printf("%d. value[0]'%s', value[1]'%s'\n", k, v[0], v[1])
				}
				fmt.Println("______________________________________________")
				fmt.Println("_____________________---- Attributes ----____________________")
				for k, v := range getNewListAttributes(
					listAttributesMisp.GetListAttributesMisp(),
					listTags) {
					fmt.Printf("%d.\n%v\n", k, v)
				}
				fmt.Println("_____________________-------____________________")

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
						Command:    "add event",
						TaskId:     taskId,
						CaseId:     caseId,
						CaseSource: caseSource,
						UserEmail:  userEmail,
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
				userEmail, caseSource = "", ""
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

// searchEventSource выполняет поиск источника события
func searchEventSource(tmf ChanInputCreateMispFormat) (string, bool) {
	var (
		source string
		ok     bool
	)

	if tmf.FieldBranch == "source" {
		source, ok = tmf.Value.(string)
	}

	return source, ok
}

// searchCaseId выполняет поиск id кейса
func searchCaseId(tmf ChanInputCreateMispFormat) (float64, bool) {
	var (
		cid float64
		ok  bool
	)

	if tmf.FieldBranch == "event.object.caseId" {
		cid, ok = tmf.Value.(float64)
	}

	return cid, ok
}

// searchOwnerEmail выполняет поиск email владельца события
func searchOwnerEmail(tmf ChanInputCreateMispFormat) (string, bool) {
	var (
		email string
		ok    bool
	)

	if tmf.FieldBranch == "event.object.owner" {
		email, ok = tmf.Value.(string)
	}

	return email, ok
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

func testSearchFieldBranch(fb string) {
	patter := regexp.MustCompile(`^misp:([\w\-].*)=\"([\w\-].*)\"$`)
	if patter.MatchString(fb) {
		fmt.Println("!!!!!!! ======== В fb содержится что то больше чем observables.<любое значение>", fb)
	}
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
