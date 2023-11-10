package mispinteractions

import (
	"encoding/json"
	"fmt"
	"net/url"
	"runtime"

	"placeholder_misp/confighandler"
	"placeholder_misp/datamodels"
	"placeholder_misp/memorytemporarystorage"
)

var mmisp ModuleMISP

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

	// добавляем event_tags
	if err := sendEventTagsMispFormat(conf.Host, authKey, eventId, data, logging); err != nil {
		_, f, l, _ := runtime.Caller(0)

		logging <- datamodels.MessageLogging{
			MsgData: fmt.Sprintf("'%s' %s:%d", err.Error(), f, l-2),
			MsgType: "error",
		}
	}

	//отправляем в ядро информацию по event Id
	mmisp.SendingDataOutput(SettingChanOutputMISP{
		Command: "send event id",
		EventId: eventId,
		TaskId:  data.TaskId,
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
