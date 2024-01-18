// Пакет mispinteractions реализует методы для взаимодействия с API MISP
package mispinteractions

import (
	"encoding/json"
	"fmt"
	"net/url"
	"runtime"

	"placeholder_misp/confighandler"
	"placeholder_misp/datamodels"
)

type RespMISP struct {
	Event map[string]interface{} `json:"event"`
}

// NewClientMISP возвращает структуру типа ClientMISP с предустановленными значениями
func NewClientMISP(h, a string, v bool) (*ClientMISP, error) {
	urlBase, err := url.Parse("http://" + h)
	if err != nil {
		return &ClientMISP{}, err
	}

	return &ClientMISP{
		BaseURL:  urlBase,
		Host:     h,
		AuthHash: a,
		Verify:   v,
	}, nil
}

// NewStorageAuthorizationDataMISP создает новое хранилище с данными пользователей MISP
func NewStorageAuthorizationDataMISP() *StorageAuthorizationData {
	return &StorageAuthorizationData{
		AuthList:         []UserSettings{},
		OrganisationList: map[string]OrganisationOptions{},
	}
}

// NewHandlerAuthorizationMISP создает новый обработчик соединений с MISP
func NewHandlerAuthorizationMISP(c ConnectMISPHandler, s *StorageAuthorizationData) *AuthorizationDataMISP {
	return &AuthorizationDataMISP{
		c,
		s,
	}
}

func HandlerMISP(
	conf confighandler.AppConfigMISP,
	organistions []confighandler.Organization,
	logging chan<- datamodels.MessageLogging) (*ModuleMISP, error) {

	mmisp := ModuleMISP{
		ChanInputMISP:  make(chan SettingsChanInputMISP),
		ChanOutputMISP: make(chan SettingChanOutputMISP),
	}

	client, err := NewClientMISP(conf.Host, conf.Auth, false)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		logging <- datamodels.MessageLogging{
			MsgData: fmt.Sprintf("'%s' %s:%d", err.Error(), f, l-2),
			MsgType: "error",
		}
	}

	connHandler := NewHandlerAuthorizationMISP(client, NewStorageAuthorizationDataMISP())
	err = connHandler.GetListAllOrganisation(organistions)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		logging <- datamodels.MessageLogging{
			MsgData: fmt.Sprintf("'%s' %s:%d", err.Error(), f, l-2),
			MsgType: "error",
		}
	}

	if countUser, err := connHandler.GetListAllUsers(); err != nil {
		_, f, l, _ := runtime.Caller(0)
		logging <- datamodels.MessageLogging{
			MsgData: fmt.Sprintf("'%s' %s:%d", err.Error(), f, l-2),
			MsgType: "error",
		}
	} else {
		logging <- datamodels.MessageLogging{
			MsgData: fmt.Sprintf("at the start of the application, all user information was downloaded %d", countUser),
			MsgType: "info",
		}
	}

	//здесь обрабатываем данные из входного канала модуля MISP
	go func() {
		for data := range mmisp.GetInputChannel() {
			authKey := conf.Auth

			// ***********************************
			// Это логирование только для теста!!!
			// ***********************************
			logging <- datamodels.MessageLogging{
				MsgData: fmt.Sprintf("TEST_INFO func 'HandlerMISP', reseived command '%s', data object '%v'", data.Command, data),
				MsgType: "testing",
			}
			//
			//

			// получаем авторизационный ключ пользователя по его email
			if data.UserEmail != "" {
				if us, err := connHandler.GetUserData(data.UserEmail); err == nil {
					authKey = us.AuthKey
				} else {
					_, f, l, _ := runtime.Caller(0)
					logging <- datamodels.MessageLogging{
						MsgData: fmt.Sprintf("'%s' %s:%d", err.Error(), f, l-3),
						MsgType: "error",
					}

					if us, err = connHandler.CreateNewUser(data.UserEmail, data.CaseSource); err != nil {
						_, f, l, _ = runtime.Caller(0)
						logging <- datamodels.MessageLogging{
							MsgData: fmt.Sprintf("'%s, case id %d' %s:%d", err.Error(), int(data.CaseId), f, l-3),
							MsgType: "error",
						}
					} else {
						authKey = us.AuthKey

						logging <- datamodels.MessageLogging{
							MsgData: fmt.Sprintf("a new user %s has been successfully created", data.UserEmail),
							MsgType: "info",
						}
					}
				}
			}

			//
			// --------------------- ТОЛЬКО ДЛЯ ОТЛАДКИ ----------------------
			//
			//if data.CaseId == 0 || data.UserEmail == "" {
			//	if ed, ok := data.MajorData["events"]; !ok {
			//		logging <- datamodels.MessageLogging{
			//			MsgData: fmt.Sprintf("TEST_ERROR func 'HandlerMISP', reseived command '%s', the properties of \"events\" were not found in the received data. DATA: %v", data.Command, data.MajorData),
			//			MsgType: "error",
			//		}
			//	} else {
			//		logging <- datamodels.MessageLogging{
			//			MsgData: fmt.Sprintf("TEST_ERROR func 'HandlerMISP', reseived command '%s', data object '%v'", data.Command, ed),
			//			MsgType: "error",
			//		}
			//	}
			//}
			//
			// ----------------------------------------------------------------
			//

			switch data.Command {
			case "add event":
				go addEvent(conf.Host, authKey, conf.Auth, data, &mmisp, logging)

			case "del event by id":
				go delEventById(conf, data.EventId, logging)
			}
		}
	}()

	return &mmisp, nil
}

func addEvent(
	host string,
	authKey string,
	masterKey string,
	data SettingsChanInputMISP,
	mmisp *ModuleMISP,
	logging chan<- datamodels.MessageLogging) {

	// ***********************************
	// Это логирование только для теста!!!
	// ***********************************
	logging <- datamodels.MessageLogging{
		MsgData: fmt.Sprintf("TEST_INFO func 'HandlerMISP', --=== RESEIVED DATA ===--	USER EMAIL: %s, ObjectId: %v", data.UserEmail, data.CaseId),
		MsgType: "testing",
	}
	//
	//

	//обработка только для события типа 'events'
	_, resBodyByte, err := sendEventsMispFormat(host, authKey, data)
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
	if err := sendEventReportsMispFormat(host, authKey, eventId, data.CaseId); err != nil {
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
	_, _ = sendAttribytesMispFormat(host, authKey, eventId, data, logging)

	// добавляем объекты
	_, _ = sendObjectsMispFormat(host, authKey, eventId, data, logging)

	// добавляем event_tags
	if err := sendEventTagsMispFormat(host, authKey, eventId, data, logging); err != nil {
		_, f, l, _ := runtime.Caller(0)

		logging <- datamodels.MessageLogging{
			MsgData: fmt.Sprintf("'%s' %s:%d", err.Error(), f, l-2),
			MsgType: "error",
		}
	}

	//публикуем добавленное событие
	//masterKey нужен для публикации события так как пользователь
	//должен иметь более расшириные права чем могут иметь некоторые
	//обычные пользователи
	resMsg, err := sendRequestPublishEvent(host, masterKey, eventId)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)

		logging <- datamodels.MessageLogging{
			MsgData: fmt.Sprintf("'%s' %s:%d", err.Error(), f, l-2),
			MsgType: "error",
		}
	}
	if resMsg != "" {
		logging <- datamodels.MessageLogging{
			MsgData: resMsg,
			MsgType: "warning",
		}
	}

	//отправляем в ядро информацию по event Id
	mmisp.SendingDataOutput(SettingChanOutputMISP{
		Command: "send event id",
		EventId: eventId,
		TaskId:  data.TaskId,
	})
}

func delEventById(
	conf confighandler.AppConfigMISP,
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
	//mmisp.SendingDataOutput(SettingChanOutputMISP{
	//	Command: "TEST STOP",
	//})
}
