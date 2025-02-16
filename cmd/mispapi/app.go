// Пакет mispapi реализует методы для взаимодействия с API MISP
package mispapi

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/av-belyakov/placeholder_misp/commoninterfaces"
	"github.com/av-belyakov/placeholder_misp/internal/confighandler"
	"github.com/av-belyakov/placeholder_misp/internal/supportingfunctions"
)

type RespMISP struct {
	Event map[string]interface{} `json:"event"`
}

// NewClientMISP конструктор API MISP
func NewClientMISP(host, authKey string, verify bool) (*ClientMISP, error) {
	urlBase, err := url.Parse("http://" + host)
	if err != nil {
		return &ClientMISP{}, err
	}

	return &ClientMISP{
		BaseURL:  urlBase,
		Host:     host,
		AuthHash: authKey,
		Verify:   verify,
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
	return &AuthorizationDataMISP{c, s}
}

// HandlerMISP взаимодействие с API MISP
func HandlerMISP(
	conf confighandler.AppConfigMISP,
	organistions []confighandler.Organization,
	logger commoninterfaces.Logger) (*ModuleMISP, error) {
	mmisp := ModuleMISP{
		ChanInput:  make(chan InputSettings),
		ChanOutput: make(chan OutputSetting),
	}

	client, err := NewClientMISP(conf.Host, conf.Auth, false)
	if err != nil {
		logger.Send("error", supportingfunctions.CustomError(err).Error())
	}

	connHandler := NewHandlerAuthorizationMISP(client, NewStorageAuthorizationDataMISP())
	err = connHandler.GetListAllOrganisation(organistions)
	if err != nil {
		logger.Send("error", supportingfunctions.CustomError(err).Error())
	}

	if countUser, err := connHandler.GetListAllUsers(); err != nil {
		logger.Send("error", supportingfunctions.CustomError(err).Error())
	} else {
		logger.Send("info", fmt.Sprintf("at the start of the application, all user information was downloaded %d", countUser))
	}

	//здесь обрабатываем данные из входного канала модуля MISP
	go func() {
		for msg := range mmisp.GetInputChannel() {
			authKey := conf.Auth

			// ***********************************
			// Это логирование только для теста!!!
			// ***********************************
			logger.Send("testing", fmt.Sprintf("TEST_INFO func 'HandlerMISP', reseived command '%s', case Id '%d'", msg.Command, int(msg.CaseId)))
			//
			//

			go func(data InputSettings) {
				// получаем авторизационный ключ пользователя по его email
				if data.UserEmail != "" {
					if us, err := connHandler.GetUserData(data.UserEmail); err == nil {
						authKey = us.AuthKey
					} else {
						logger.Send("error", supportingfunctions.CustomError(err).Error())

						if us, err = connHandler.CreateNewUser(data.UserEmail, data.CaseSource); err != nil {
							logger.Send("error", supportingfunctions.CustomError(fmt.Errorf("%w, case id:'%d'", err, int(data.CaseId))).Error())
						} else {
							authKey = us.AuthKey
							logger.Send("info", fmt.Sprintf("a new user '%s' has been successfully created", data.UserEmail))
						}
					}
				}

				switch data.Command {
				case "add event":
					addEvent(conf.Host, authKey, conf.Auth, data, &mmisp, logger)

				case "del event by id":
					delEventById(conf, data.EventId, logger)
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
	data InputSettings,
	mmisp *ModuleMISP,
	logger commoninterfaces.Logger) {

	// ***********************************
	// Это логирование только для теста!!!
	// ***********************************
	logger.Send("testing", fmt.Sprintf("TEST_INFO func 'HandlerMISP', send EVENTS to ----> MISP	USER EMAIL: %s, CaseId: %v", data.UserEmail, data.CaseId))
	//
	//

	//обработка только для события типа 'events'
	_, resBodyByte, err := sendEventsMispFormat(host, authKey, data)
	if err != nil {
		logger.Send("error", supportingfunctions.CustomError(err).Error())

		return
	}

	resMisp := RespMISP{}
	if err := json.Unmarshal(resBodyByte, &resMisp); err != nil {
		logger.Send("error", supportingfunctions.CustomError(err).Error())

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
		logger.Send("error", supportingfunctions.CustomError(fmt.Errorf("the formation of events of the 'Attributes' type was not performed because the EventID is empty")).Error())

		return
	}

	// ***********************************
	// Это логирование только для теста!!!
	// ***********************************
	logger.Send("testing", fmt.Sprintf("TEST_INFO func 'HandlerMISP', send EVENTS REPORTS to ----> MISP	USER EMAIL: %s, CaseId: %v", data.UserEmail, data.CaseId))
	//
	//

	// добавляем event_reports
	if err := sendEventReportsMispFormat(host, authKey, eventId, data.CaseId); err != nil {
		logger.Send("error", supportingfunctions.CustomError(err).Error())
	}

	// ***********************************
	// Это логирование только для теста!!!
	// ***********************************
	logger.Send("testing", fmt.Sprintf("TEST_INFO func 'HandlerMISP', send data to ----> RedisDB USER EMAIL: %s, CaseId: %v", data.UserEmail, data.CaseId))
	//
	//

	//отправляем запрос для добавления в БД Redis, id кейса и нового события
	mmisp.SendingDataOutput(OutputSetting{
		Command: "set new event id",
		CaseId:  fmt.Sprint(data.CaseId),
		EventId: eventId,
	})

	// ***********************************
	// Это логирование только для теста!!!
	// ***********************************
	logger.Send("testing", fmt.Sprintf("TEST_INFO func 'HandlerMISP', send ATTRIBYTES to ----> MISP	USER EMAIL: %s, CaseId: %v", data.UserEmail, data.CaseId))
	//
	//

	//добавляем атрибуты
	_, _ = sendAttribytesMispFormat(host, authKey, eventId, data, logger)

	// ***********************************
	// Это логирование только для теста!!!
	// ***********************************
	logger.Send("testing", fmt.Sprintf("TEST_INFO func 'HandlerMISP', send OBJECTS to ----> MISP	USER EMAIL: %s, CaseId: %v", data.UserEmail, data.CaseId))
	//
	//

	// добавляем объекты
	_, _ = sendObjectsMispFormat(host, authKey, eventId, data, logger)

	//берем небольшой таймаут
	time.Sleep(4 * time.Second)

	// ***********************************
	// Это логирование только для теста!!!
	// ***********************************
	logger.Send("testing", fmt.Sprintf("TEST_INFO func 'HandlerMISP', send EVENT_TAGS to ----> MISP	USER EMAIL: %s, CaseId: %v", data.UserEmail, data.CaseId))
	//
	//

	// добавляем event_tags
	if err := sendEventTagsMispFormat(host, masterKey, eventId, data, logger); err != nil {
		logger.Send("error", supportingfunctions.CustomError(err).Error())
	}

	// ***********************************
	// Это логирование только для теста!!!
	// ***********************************
	logger.Send("testing", fmt.Sprintf("TEST_INFO func 'HandlerMISP', send PUBLISH to ----> MISP	USER EMAIL: %s, CaseId: %v", data.UserEmail, data.CaseId))
	//
	//

	//публикуем добавленное событие
	//masterKey нужен для публикации события так как пользователь
	//должен иметь более расшириные права чем могут иметь некоторые
	//обычные пользователи
	resMsg, err := sendRequestPublishEvent(host, masterKey, eventId)
	if err != nil {
		logger.Send("error", supportingfunctions.CustomError(err).Error())
	}
	if resMsg != "" {
		logger.Send("warning", resMsg)
	}

	// ***********************************
	// Это логирование только для теста!!!
	// ***********************************
	logger.Send("testing", fmt.Sprintf("TEST_INFO func 'HandlerMISP', send EventId to ----> CORE	USER EMAIL: %s, CaseId: %v", data.UserEmail, data.CaseId))
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
	mmisp.SendingDataOutput(OutputSetting{
		Command: "send event id",
		EventId: eventId,
		CaseId:  fmt.Sprint(data.CaseId),
		RootId:  data.RootId,
		TaskId:  data.TaskId,
	})
}

func delEventById(conf confighandler.AppConfigMISP, eventId string, logger commoninterfaces.Logger) {

	// ***********************************
	// Это логирование только для теста!!!
	// ***********************************
	logger.Send("testing", fmt.Sprintf("TEST_INFO func 'HandlerMISP', удаление события типа event, где event id: %s", eventId))
	//
	//

	// удаление события типа event
	_, err := delEventsMispFormat(conf.Host, conf.Auth, eventId)
	if err != nil {
		logger.Send("error", supportingfunctions.CustomError(err).Error())
	}

	// ***********************************
	// Это логирование только для теста!!!
	// ***********************************
	logger.Send("testing", fmt.Sprintf("TEST_INFO func 'HandlerMISP', должно было быть успешно выполненно удаление события event id: %s", eventId))
	//
	//

	// ***********************************
	// Это логирование только для теста!!!
	// ***********************************
	logger.Send("testing", "TEST_INFO STOP")
	//
	//

	//
	//только для теста, для ОСТАНОВА
	//
	//mmisp.SendingDataOutput(SettingChanOutputMISP{
	//	Command: "TEST STOP",
	//})
}
