// Пакет mispinteractions реализует методы для взаимодействия с API MISP
package mispinteractions

import (
	"encoding/json"
	"fmt"
	"net/url"
	"runtime"
	"time"

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
		for msg := range mmisp.GetInputChannel() {
			authKey := conf.Auth

			// ***********************************
			// Это логирование только для теста!!!
			// ***********************************
			logging <- datamodels.MessageLogging{
				MsgData: fmt.Sprintf("TEST_INFO func 'HandlerMISP', reseived command '%s', case Id '%d'", msg.Command, int(msg.CaseId)),
				MsgType: "testing",
			}
			//
			//

			go func(data SettingsChanInputMISP) {
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

				switch data.Command {
				case "add event":
					addEvent(conf.Host, authKey, conf.Auth, data, &mmisp, logging)

				case "del event by id":
					delEventById(conf, data.EventId, logging)
				}
			}(msg)
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
		MsgData: fmt.Sprintf("TEST_INFO func 'HandlerMISP', send EVENTS to ----> MISP	USER EMAIL: %s, CaseId: %v", data.UserEmail, data.CaseId),
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
		MsgData: fmt.Sprintf("TEST_INFO func 'HandlerMISP', send EVENTS REPORTS to ----> MISP	USER EMAIL: %s, CaseId: %v", data.UserEmail, data.CaseId),
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

	// ***********************************
	// Это логирование только для теста!!!
	// ***********************************
	logging <- datamodels.MessageLogging{
		MsgData: fmt.Sprintf("TEST_INFO func 'HandlerMISP', send data to ----> RedisDB USER EMAIL: %s, CaseId: %v", data.UserEmail, data.CaseId),
		MsgType: "testing",
	}
	//
	//

	//отправляем запрос для добавления в БД Redis, id кейса и нового события
	mmisp.SendingDataOutput(SettingChanOutputMISP{
		Command: "set new event id",
		CaseId:  fmt.Sprint(data.CaseId),
		EventId: eventId,
	})

	// ***********************************
	// Это логирование только для теста!!!
	// ***********************************
	logging <- datamodels.MessageLogging{
		MsgData: fmt.Sprintf("TEST_INFO func 'HandlerMISP', send ATTRIBYTES to ----> MISP	USER EMAIL: %s, CaseId: %v", data.UserEmail, data.CaseId),
		MsgType: "testing",
	}
	//
	//

	//добавляем атрибуты
	_, _ = sendAttribytesMispFormat(host, authKey, eventId, data, logging)

	// ***********************************
	// Это логирование только для теста!!!
	// ***********************************
	logging <- datamodels.MessageLogging{
		MsgData: fmt.Sprintf("TEST_INFO func 'HandlerMISP', send OBJECTS to ----> MISP	USER EMAIL: %s, CaseId: %v", data.UserEmail, data.CaseId),
		MsgType: "testing",
	}
	//
	//

	// добавляем объекты
	_, _ = sendObjectsMispFormat(host, authKey, eventId, data, logging)

	//берем небольшой таймаут
	time.Sleep(4 * time.Second)

	// ***********************************
	// Это логирование только для теста!!!
	// ***********************************
	logging <- datamodels.MessageLogging{
		MsgData: fmt.Sprintf("TEST_INFO func 'HandlerMISP', send EVENT_TAGS to ----> MISP	USER EMAIL: %s, CaseId: %v", data.UserEmail, data.CaseId),
		MsgType: "testing",
	}
	//
	//

	// добавляем event_tags
	if err := sendEventTagsMispFormat(host, masterKey, eventId, data, logging); err != nil {
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
		MsgData: fmt.Sprintf("TEST_INFO func 'HandlerMISP', send PUBLISH to ----> MISP	USER EMAIL: %s, CaseId: %v", data.UserEmail, data.CaseId),
		MsgType: "testing",
	}
	//
	//

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

	// ***********************************
	// Это логирование только для теста!!!
	// ***********************************
	logging <- datamodels.MessageLogging{
		MsgData: fmt.Sprintf("TEST_INFO func 'HandlerMISP', send EventId to ----> CORE	USER EMAIL: %s, CaseId: %v", data.UserEmail, data.CaseId),
		MsgType: "testing",
	}
	//
	//

	/*
	   типы команд для передачи в NATS

	   requests := map[string][]byte{
	   					"add_case_tag": []byte(
	   						fmt.Sprintf(`{
	   							"service": "MISP",
	   							"command": "add_case_tag",
	   							"root_id": "%s",
	   							"case_id": "%s",
	   							"value": "Webhook: send=\"MISP\""}`,
	   							data.RootId,
	   							data.CaseId)),
	   					"set_case_custom_field": []byte(
	   						fmt.Sprintf(`{
	   							"service": "MISP",
	   							"command": "set_case_custom_field",
	   							"root_id": "%s",
	   	  						"field_name": "misp-event-id.string",
	   							"value": "%s"}`,
	   							data.RootId,
	   							data.EventId)),
	   				}

	*/

	//отправляем в ядро информацию по event Id
	mmisp.SendingDataOutput(SettingChanOutputMISP{
		Command: "send event id",
		EventId: eventId,
		CaseId:  fmt.Sprint(data.CaseId),
		RootId:  data.RootId,
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
