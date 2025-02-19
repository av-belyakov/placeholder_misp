// Пакет mispapi реализует методы для взаимодействия с API MISP
package mispapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

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
	cache, err := cachingstoragewithqueue.NewCacheStorage[objectsmispformat.ListFormatsMISP](
		cachingstoragewithqueue.WithMaxTtl[objectsmispformat.ListFormatsMISP](300),
		cachingstoragewithqueue.WithTimeTick[objectsmispformat.ListFormatsMISP](3),
		cachingstoragewithqueue.WithMaxSize[objectsmispformat.ListFormatsMISP](10),
		cachingstoragewithqueue.WithEnableAsyncProcessing[objectsmispformat.ListFormatsMISP](1))
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
	err = connHandler.GetListAllOrganisation(ctx, m.organistions)
	if err != nil {
		//m.logger.Send("error", supportingfunctions.CustomError(err).Error())
		return err
	}

	countUser, err := connHandler.GetListAllUsers(ctx)
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

				// ***********************************
				// Это логирование только для теста!!!
				// ***********************************
				m.logger.Send("testing", fmt.Sprintf("TEST_INFO func 'HandlerMISP', reseived command '%s', case Id '%d'", msg.Command, int(msg.CaseId)))
				//
				//

				go func(data InputSettings) {
					userAuthKey := m.authKey

					// получаем авторизационный ключ пользователя по его email
					// m.authKey отличается от authKey тем что m.authKey является ключем с
					// повышенными привилегиями которые необходимы для большинства 'важных'
					// операций, например таких как удаление объектов, а authKey это авторизационный
					// ключ пользователя от имени которого выполняются операции по добавлению объектов MISP
					if data.UserEmail != "" {
						if us, err := connHandler.GetUserData(ctx, data.UserEmail); err == nil {
							userAuthKey = us.AuthKey
						} else {
							m.logger.Send("error", supportingfunctions.CustomError(err).Error())

							if us, err = connHandler.CreateNewUser(ctx, data.UserEmail, data.CaseSource); err != nil {
								m.logger.Send("error", supportingfunctions.CustomError(fmt.Errorf("%w, case id:'%d'", err, int(data.CaseId))).Error())
							} else {
								userAuthKey = us.AuthKey
								m.logger.Send("info", fmt.Sprintf("a new user '%s' has been successfully created", data.UserEmail))
							}
						}
					}

					switch data.Command {
					case "add event":
						m.addNewObject(ctx, userAuthKey, data)

					case "del event by id":
						delEventById(ctx, m.host, m.authKey, data.EventId, m.logger)

					}
				}(msg)

			}
		}
	}()

	return nil
}

func (m *ModuleMISP) addNewObject(ctx context.Context, userAuthKey string, data InputSettings) {
	specialObject := NewCacheSpecialObject[*objectsmispformat.ListFormatsMISP]()
	specialObject.SetID(data.RootId)
	specialObject.SetObject(&data.Data)
	specialObject.SetFunc(func(i int) bool {

		rmisp, err := NewMispRequest(
			WithHost(m.host),
			WithUserAuthKey(userAuthKey),
			WithMasterAuthKey(m.authKey))
		if err != nil {
			return false
		}

		//отправляет в API MISP событие в виде типа Event и возвращает полученный ответ
		_, resBodyByte, err := rmisp.sendEvent(ctx, data.Data.GetEvent())
		if err != nil {
			m.logger.Send("error", supportingfunctions.CustomError(err).Error())

			return false
		}

		resMisp := MispResponse{}
		if err := json.Unmarshal(resBodyByte, &resMisp); err != nil {
			m.logger.Send("error", supportingfunctions.CustomError(err).Error())

			return false
		}

		/*

			сюда нужно перенести все функции необходимые для добавления данных в MISP

		*/

		//по умолчанию 'не успешно'
		return false
	})
}

/*
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



	//отправляем в ядро информацию по event Id
	mmisp.SendingDataOutput(SettingChanOutputMISP{
		Command: "send event id",
		EventId: eventId,
		CaseId:  fmt.Sprint(data.CaseId),
		RootId:  data.RootId,
		TaskId:  data.TaskId,
	})
}
*/
