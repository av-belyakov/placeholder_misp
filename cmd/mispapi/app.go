// Пакет mispapi реализует методы для взаимодействия с API MISP
package mispapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/av-belyakov/cachingstoragewithqueue"
	"github.com/av-belyakov/objectsmispformat"
	"github.com/av-belyakov/placeholder_misp/commoninterfaces"
	"github.com/av-belyakov/placeholder_misp/internal/confighandler"
	"github.com/av-belyakov/placeholder_misp/internal/supportingfunctions"
)

// NewStorageAuthorizationDataMISP хранилище с данными пользователей MISP
func NewStorageAuthorizationDataMISP() *StorageAuthorizationData {
	return &StorageAuthorizationData{
		AuthList:         []UserSettings{},
		OrganisationList: map[string]OrganisationOptions{},
	}
}

// NewHandlerAuthorizationMISP обработчик соединений с MISP
func NewHandlerAuthorizationMISP(c ConnectMISPHandler, s *StorageAuthorizationData) *AuthorizationDataMISP {
	return &AuthorizationDataMISP{c, s}
}

// NewClientMISP клиент API MISP
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

// NewModuleMISP модуль взаимодействия с API MISP
func NewModuleMISP(
	host string,
	authKey string,
	org []confighandler.Organization,
	logger commoninterfaces.Logger) (*ModuleMISP, error) {
	cache, err := cachingstoragewithqueue.NewCacheStorage[*objectsmispformat.ListFormatsMISP](
		cachingstoragewithqueue.WithMaxTtl[*objectsmispformat.ListFormatsMISP](300),
		cachingstoragewithqueue.WithTimeTick[*objectsmispformat.ListFormatsMISP](3),
		cachingstoragewithqueue.WithMaxSize[*objectsmispformat.ListFormatsMISP](10),
		cachingstoragewithqueue.WithEnableAsyncProcessing[*objectsmispformat.ListFormatsMISP](1))
	if err != nil {
		return &ModuleMISP{}, err
	}

	module := &ModuleMISP{
		cache:        cache,
		organistions: org,
		host:         host,
		authKey:      authKey,
		chInput:      make(chan InputSettings),
		chOutput:     make(chan OutputSetting),
	}

	return module, nil
}

func (m *ModuleMISP) Start(ctx context.Context) error {
	client, err := NewClientMISP(m.host, m.authKey, false)
	if err != nil {
		//m.logger.Send("error", supportingfunctions.CustomError(err).Error())
		return err
	}

	connHandler := NewHandlerAuthorizationMISP(client, NewStorageAuthorizationDataMISP())
	err = connHandler.GetListAllOrganisation(m.organistions)
	if err != nil {
		//m.logger.Send("error", supportingfunctions.CustomError(err).Error())
		return err
	}

	countUser, err := connHandler.GetListAllUsers()
	if err != nil {
		//m.logger.Send("error", supportingfunctions.CustomError(err).Error())
		return err
	}

	m.logger.Send("info", fmt.Sprintf("at the start of the application, all user information was downloaded %d", countUser))

	//здесь обрабатываем данные из входного канала модуля MISP
	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case msg := <-m.GetInputChannel():
				authKey := m.authKey

				// ***********************************
				// Это логирование только для теста!!!
				// ***********************************
				m.logger.Send("testing", fmt.Sprintf("TEST_INFO func 'HandlerMISP', reseived command '%s', case Id '%d'", msg.Command, int(msg.CaseId)))
				//
				//

				go func(data InputSettings) {
					// получаем авторизационный ключ пользователя по его email
					if data.UserEmail != "" {
						if us, err := connHandler.GetUserData(data.UserEmail); err == nil {
							authKey = us.AuthKey
						} else {
							m.logger.Send("error", supportingfunctions.CustomError(err).Error())

							if us, err = connHandler.CreateNewUser(data.UserEmail, data.CaseSource); err != nil {
								m.logger.Send("error", supportingfunctions.CustomError(fmt.Errorf("%w, case id:'%d'", err, int(data.CaseId))).Error())
							} else {
								authKey = us.AuthKey
								m.logger.Send("info", fmt.Sprintf("a new user '%s' has been successfully created", data.UserEmail))
							}
						}
					}

					switch data.Command {
					case "add event":
						addEvent(m.host, authKey, m.authKey, data, m, m.logger)

					case "del event by id":
						delEventById(m.host, m.authKey, data.EventId, m.logger)
					}

					//m.authKey отличается от authKey тем что m.authKey является ключем с повышенными привилегиями
					//которые необходимы для большинства 'важных' операций, например таких как удаление объектов,
					//а authKey это авторизационный ключ пользователя от имени которого выполняются операции по
					//добавлению объектов MISP
				}(msg)

			}
		}
	}()

	return nil
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

func delEventById(host, authKey, eventId string, logger commoninterfaces.Logger) {

	// ***********************************
	// Это логирование только для теста!!!
	// ***********************************
	logger.Send("testing", fmt.Sprintf("TEST_INFO func 'HandlerMISP', удаление события типа event, где event id: %s", eventId))
	//
	//

	// удаление события типа event
	_, err := delEventsMispFormat(host, authKey, eventId)
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
