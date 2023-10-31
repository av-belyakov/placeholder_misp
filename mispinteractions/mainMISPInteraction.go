package mispinteractions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"runtime"

	"placeholder_misp/confighandler"
	"placeholder_misp/datamodels"
	"placeholder_misp/memorytemporarystorage"
	"placeholder_misp/supportingfunctions"
)

var mmisp ModuleMISP

type ClientMISP struct {
	BaseURL  *url.URL
	Host     string
	AuthHash string
	Verify   bool
}

type RespMISP struct {
	Event map[string]interface{} `json:"event"`
}

func init() {
	mmisp = ModuleMISP{
		ChanInputMISP:  make(chan SettingsChanInputMISP),
		ChanOutputMISP: make(chan SettingChanOutputMISP),
	}
}

func HandlerMISP(
	conf confighandler.AppConfigMISP,
	storageApp *memorytemporarystorage.CommonStorageTemporary,
	logging chan<- datamodels.MessageLogging) (*ModuleMISP, error) {

	//выполнеяем запрос для получения настроек пользователей через API MISP
	//и сохраняем полученные параметры во временном хранилище
	err := getUserMisp(conf, storageApp)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)

		logging <- datamodels.MessageLogging{
			MsgData: fmt.Sprintf("'%s' %s:%d", err.Error(), f, l-2),
			MsgType: "error",
		}
	}

	//здесь обрабатываем данные из входного канала модуля MISP
	go func() {
		for data := range mmisp.GetInputChannel() {

			// ***********************************
			// Это логирование только для теста!!!
			// ***********************************
			logging <- datamodels.MessageLogging{
				MsgData: fmt.Sprintf("TEST_INFO func 'HandlerMISP', reseived command '%s', data object '%v'", data.Command, data),
				MsgType: "testing",
			}
			//
			//

			switch data.Command {
			case "add event":
				go addEvent(conf, storageApp, data, logging)

			case "del event by id":
				go delEventById(conf, data.EventId, logging)
			}
		}
	}()

	return &mmisp, nil
}

// NewClientMISP возвращает структуру типа ClientMISP с предустановленными значениями
func NewClientMISP(h, a string, v bool) (ClientMISP, error) {
	urlBase, err := url.Parse("http://" + h)
	if err != nil {
		return ClientMISP{}, err
	}

	return ClientMISP{
		BaseURL:  urlBase,
		Host:     h,
		AuthHash: a,
		Verify:   v,
	}, nil
}

func addEvent(conf confighandler.AppConfigMISP,
	storageApp *memorytemporarystorage.CommonStorageTemporary,
	data SettingsChanInputMISP,
	logging chan<- datamodels.MessageLogging) {
	authKey := conf.Auth

	// ***********************************
	// Это логирование только для теста!!!
	// ***********************************
	logging <- datamodels.MessageLogging{
		MsgData: fmt.Sprintf("TEST_INFO func 'HandlerMISP', --=== RESEIVED DATA ===--	USER EMAIL: %s, ObjectId: %v", data.UserEmail, data.CaseId),
		MsgType: "testing",
	}
	//
	//

	// получаем авторизационный ключ пользователя по его email
	if us, ok := storageApp.GetUserSettingsMISP(data.UserEmail); ok {
		authKey = us.AuthKey
	}

	//обработка только для события типа 'events'
	_, resBodyByte, err := sendEventsMispFormat(conf.Host, authKey, data)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)

		logging <- datamodels.MessageLogging{
			MsgData: fmt.Sprintf("'%s' %s:%d", err.Error(), f, l-2),
			MsgType: "error",
		}

		return
	}

	resMisp := RespMISP{}
	if err := json.Unmarshal(resBodyByte, &resMisp); err != nil {
		_, f, l, _ := runtime.Caller(0)

		logging <- datamodels.MessageLogging{
			MsgData: fmt.Sprintf("'%s' %s:%d", err.Error(), f, l-2),
			MsgType: "error",
		}

		return
	}

	var eventId string
	for key, value := range resMisp.Event {
		if key == "id" {
			if str, ok := value.(string); ok {
				eventId = str

				break
			}
		}
	}

	if eventId == "" {
		_, f, l, _ := runtime.Caller(0)

		logging <- datamodels.MessageLogging{
			MsgData: fmt.Sprintf("'the formation of events of the 'Attributes' type was not performed because the EventID is empty' %s:%d", f, l-1),
			MsgType: "error",
		}

		return
	}

	// ***********************************
	// Это логирование только для теста!!!
	// ***********************************
	logging <- datamodels.MessageLogging{
		MsgData: fmt.Sprintf("TEST_INFO func 'HandlerMISP', отправляем запрос для добавления в БД Redis, id кейса и нового события, где case id: %v, event id: %s", data.CaseId, eventId),
		MsgType: "testing",
	}
	//
	//

	// добавляем event_reports
	if err := sendEventReportsMispFormat(conf.Host, authKey, eventId, data.CaseId, logging); err != nil {
		_, f, l, _ := runtime.Caller(0)

		logging <- datamodels.MessageLogging{
			MsgData: fmt.Sprintf("'%s' %s:%d", err.Error(), f, l-2),
			MsgType: "error",
		}
	}

	//отправляем запрос для добавления в БД Redis, id кейса и нового события
	mmisp.SendingDataOutput(SettingChanOutputMISP{
		Command: "set new event id",
		CaseId:  fmt.Sprint(data.CaseId),
		EventId: eventId,
	})

	//добавляем атрибуты
	_, _ = sendAttribytesMispFormat(conf.Host, authKey, eventId, data, logging)

	// добавляем объекты
	_, _ = sendObjectsMispFormat(conf.Host, authKey, eventId, data, logging)

	//отправляем в ядро информацию по event Id
	mmisp.SendingDataOutput(SettingChanOutputMISP{
		Command: "send event id",
		EventId: eventId,
	})
}

func delEventById(conf confighandler.AppConfigMISP,
	eventId string,
	logging chan<- datamodels.MessageLogging) {

	// ***********************************
	// Это логирование только для теста!!!
	// ***********************************
	logging <- datamodels.MessageLogging{
		MsgData: fmt.Sprintf("TEST_INFO func 'HandlerMISP', удаление события типа event, где event id: %s", eventId),
		MsgType: "testing",
	}
	//
	//

	// удаление события типа event
	_, err := delEventsMispFormat(conf.Host, conf.Auth, eventId)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)

		logging <- datamodels.MessageLogging{
			MsgData: fmt.Sprintf("'%s' %s:%d", err.Error(), f, l-2),
			MsgType: "error",
		}
	}

	// ***********************************
	// Это логирование только для теста!!!
	// ***********************************
	logging <- datamodels.MessageLogging{
		MsgData: fmt.Sprintf("TEST_INFO func 'HandlerMISP', должно было быть успешно выполненно удаление события event id: %s", eventId),
		MsgType: "testing",
	}
	//
	//

	// ***********************************
	// Это логирование только для теста!!!
	// ***********************************
	logging <- datamodels.MessageLogging{
		MsgData: "TEST_INFO STOP",
		MsgType: "testing",
	}
	//
	//

	//
	//только для теста, для ОСТАНОВА
	//
	mmisp.SendingDataOutput(SettingChanOutputMISP{
		Command: "TEST STOP",
	})
}

// getUserMisp выполнеяет запрос для получения настроек пользователей через API MISP
// и сохраняет полученные параметры во временном хранилище
func getUserMisp(conf confighandler.AppConfigMISP, storageApp *memorytemporarystorage.CommonStorageTemporary) error {
	client, err := NewClientMISP(conf.Host, conf.Auth, false)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return fmt.Errorf("'%s' %s:%d", err.Error(), f, l-2)
	}

	_, resByte, err := client.Get("/admin/users", nil)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return fmt.Errorf("'%s' %s:%d", err.Error(), f, l-2)
	}

	usmispf := []datamodels.UsersSettingsMispFormat{}
	err = json.Unmarshal(resByte, &usmispf)
	if err != nil {

		_, f, l, _ := runtime.Caller(0)
		return fmt.Errorf("'%s' %s:%d", err.Error(), f, l-2)
	}

	for _, v := range usmispf {
		storageApp.AddUserSettingsMISP(memorytemporarystorage.UserSettingsMISP{
			UserId:  v.User.Id,
			OrgId:   v.Organisation.Id,
			OrgName: v.Organisation.Name,
			Email:   v.User.Email,
			AuthKey: v.User.Authkey,
			Role:    v.Role.Name,
		})
	}

	return nil
}

// sendEventsMispFormat отправляет в API MISP событие в виде типа Event и возвращает полученный ответ
func sendEventsMispFormat(host, authKey string, d SettingsChanInputMISP) (*http.Response, []byte, error) {
	var (
		res         *http.Response
		resBodyByte = make([]byte, 0)
	)

	c, err := NewClientMISP(host, authKey, false)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return nil, resBodyByte, fmt.Errorf("'events add, %s' %s:%d", err.Error(), f, l-2)
	}

	ed, ok := d.MajorData["events"]
	if !ok {
		_, f, l, _ := runtime.Caller(0)

		return nil, resBodyByte, fmt.Errorf("'the properties of \"events\" were not found in the received data' %s:%d", f, l-2)
	}

	b, err := json.Marshal(ed)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)

		return nil, resBodyByte, fmt.Errorf("'events add, %s' %s:%d", err.Error(), f, l-2)
	}

	res, resBodyByte, err = c.Post("/events/add", b)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)

		return nil, resBodyByte, fmt.Errorf("'events add, %s' %s:%d", err.Error(), f, l-2)
	}

	if res.StatusCode != http.StatusOK {
		_, f, l, _ := runtime.Caller(0)

		return nil, resBodyByte, fmt.Errorf("'events add, %s' %s:%d", res.Status, f, l-1)
	}

	return res, resBodyByte, nil
}

// sendAttribytesMispFormat отправляет в API MISP список атрибутов в виде среза типов Attribytes и возвращает полученный ответ
func sendAttribytesMispFormat(host, authKey, eventId string, d SettingsChanInputMISP, logging chan<- datamodels.MessageLogging) (*http.Response, []byte) {
	var (
		res         *http.Response
		resBodyByte = make([]byte, 0)
	)

	c, err := NewClientMISP(host, authKey, false)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		logging <- datamodels.MessageLogging{
			MsgData: fmt.Sprintf("'attributes for event id%s add, %s' %s:%d", eventId, err.Error(), f, l-2),
			MsgType: "error",
		}

		return nil, resBodyByte
	}

	ad, ok := d.MajorData["attributes"]
	if !ok {
		_, f, l, _ := runtime.Caller(0)

		logging <- datamodels.MessageLogging{
			MsgData: fmt.Sprintf("'the properties of \"attributes\" were not found in the received data' %s:%d", f, l-2),
			MsgType: "error",
		}

		return nil, resBodyByte
	}

	lamf, ok := ad.([]datamodels.AttributesMispFormat)
	if !ok {
		_, f, l, _ := runtime.Caller(0)

		logging <- datamodels.MessageLogging{
			MsgData: fmt.Sprintf("'the received data does not match the type \"attributes\"' %s:%d", f, l-2),
			MsgType: "error",
		}

		return nil, resBodyByte
	}

	for k := range lamf {
		lamf[k].EventId = eventId

		if lamf[k].Value == "" {
			_, f, l, _ := runtime.Caller(0)

			logging <- datamodels.MessageLogging{
				MsgData: fmt.Sprintf("'attributes №%s add, the \"Value\" type property should not be empty' %s:%d", eventId, f, l-1),
				MsgType: "warning",
			}

			continue
		}

		b, err := json.Marshal(lamf[k])
		if err != nil {
			_, f, l, _ := runtime.Caller(0)

			logging <- datamodels.MessageLogging{
				MsgData: fmt.Sprintf("'attributes №%s add, %s' %s:%d", eventId, err.Error(), f, l-2),
				MsgType: "warning",
			}

			continue
		}

		fmt.Println("|||||||||||||||||||||||||||||||||||||_______________________ START _____________________------")
		fmt.Println(supportingfunctions.NewReadReflectJSONSprint(b))
		fmt.Println("|||||||||||||||||||||||||||||||||||||________________________ END ______________________------")

		res, resBodyByte, err = c.Post("/attributes/add/"+eventId, b)
		if err != nil {
			_, f, l, _ := runtime.Caller(0)

			logging <- datamodels.MessageLogging{
				MsgData: fmt.Sprintf("'attributes №%s add, %s' %s:%d", eventId, err.Error(), f, l-2),
				MsgType: "warning",
			}

			continue
		}

		if res.StatusCode != http.StatusOK {
			_, f, l, _ := runtime.Caller(0)

			logging <- datamodels.MessageLogging{
				MsgData: fmt.Sprintf("'attributes №%s add, %s' %s:%d", eventId, res.Status, f, l-1),
				MsgType: "warning",
			}
		}
	}

	return res, resBodyByte
}

// sendObjectsMispFormat отправляет в API MISP список объектов содержащихся в свойстве observables.attachment (как правило это описание вложеного файла)
func sendObjectsMispFormat(host, authKey, eventId string, d SettingsChanInputMISP, logging chan<- datamodels.MessageLogging) (*http.Response, []byte) {
	var (
		res         *http.Response
		resBodyByte = make([]byte, 0)
	)

	c, err := NewClientMISP(host, authKey, false)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		logging <- datamodels.MessageLogging{
			MsgData: fmt.Sprintf("'objects for event id:%s add, %s' %s:%d", eventId, err.Error(), f, l-2),
			MsgType: "error",
		}

		return nil, resBodyByte
	}

	od, ok := d.MajorData["objects"]
	if !ok {
		_, f, l, _ := runtime.Caller(0)

		logging <- datamodels.MessageLogging{
			MsgData: fmt.Sprintf("'the properties of \"objects\" were not found in the received data' %s:%d", f, l-2),
			MsgType: "error",
		}

		return nil, resBodyByte
	}

	lomf, ok := od.(map[int]datamodels.ObjectsMispFormat)
	if !ok {
		_, f, l, _ := runtime.Caller(0)

		logging <- datamodels.MessageLogging{
			MsgData: fmt.Sprintf("'the received data does not match the type \"objects\"' %s:%d", f, l-2),
			MsgType: "error",
		}

		return nil, resBodyByte
	}

	for _, v := range lomf {
		v.EventId = eventId

		b, err := json.Marshal(v)
		if err != nil {
			_, f, l, _ := runtime.Caller(0)

			logging <- datamodels.MessageLogging{
				MsgData: fmt.Sprintf("'objects №%s add, %s' %s:%d", eventId, err.Error(), f, l-2),
				MsgType: "warning",
			}

			continue
		}

		res, resBodyByte, err = c.Post("/objects/add/"+eventId, b)
		if err != nil {
			_, f, l, _ := runtime.Caller(0)

			logging <- datamodels.MessageLogging{
				MsgData: fmt.Sprintf("'objects №%s add, %s' %s:%d", eventId, err.Error(), f, l-2),
				MsgType: "warning",
			}

			continue
		}

		if res.StatusCode != http.StatusOK {
			_, f, l, _ := runtime.Caller(0)

			logging <- datamodels.MessageLogging{
				MsgData: fmt.Sprintf("'objects №%s add, %s' %s:%d", eventId, res.Status, f, l-1),
				MsgType: "warning",
			}
		}
	}

	return res, resBodyByte
}

func sendEventReportsMispFormat(host, authKey, eventId string, caseId float64, logging chan<- datamodels.MessageLogging) error {
	c, err := NewClientMISP(host, authKey, false)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return fmt.Errorf("'events add, %s' %s:%d", err.Error(), f, l-2)
	}

	b, err := json.Marshal(datamodels.EventReports{
		Name:         fmt.Sprint(caseId),
		Distribution: "0",
	})
	if err != nil {
		_, f, l, _ := runtime.Caller(0)

		return fmt.Errorf("'events add, %s' %s:%d", err.Error(), f, l-2)
	}

	res, _, err := c.Post("/event_reports/add/"+eventId, b)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)

		return fmt.Errorf("'events add, %s' %s:%d", err.Error(), f, l-2)
	}

	if res.StatusCode != http.StatusOK {
		_, f, l, _ := runtime.Caller(0)

		return fmt.Errorf("'events add, %s' %s:%d", res.Status, f, l-1)
	}

	return nil
}

// удаляем дублирующиеся события из MISP
func delEventsMispFormat(host, authKey, eventId string) (*http.Response, error) {
	fmt.Println("func 'delEventsMispFormat', удаляем дублирующиеся события из MISP START...")

	c, err := NewClientMISP(host, authKey, false)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return nil, fmt.Errorf("'events delete, %s' %s:%d", err.Error(), f, l-2)
	}

	res, _, err := c.Delete("/events/delete/" + eventId)
	if err != nil {
		return nil, err
	}

	fmt.Println("func 'delEventsMispFormat', res = ", res.Status)

	return res, nil
}
