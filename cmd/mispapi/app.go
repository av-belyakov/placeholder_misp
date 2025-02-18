// Пакет mispapi реализует методы для взаимодействия с API MISP
package mispapi

import (
	"context"
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

				// ***********************************
				// Это логирование только для теста!!!
				// ***********************************
				m.logger.Send("testing", fmt.Sprintf("TEST_INFO func 'HandlerMISP', reseived command '%s', case Id '%d'", msg.Command, int(msg.CaseId)))
				//
				//

				go func(data InputSettings) {
					userAuthKey := m.authKey
					// получаем авторизационный ключ пользователя по его email
					if data.UserEmail != "" {
						if us, err := connHandler.GetUserData(data.UserEmail); err == nil {
							userAuthKey = us.AuthKey
						} else {
							m.logger.Send("error", supportingfunctions.CustomError(err).Error())

							if us, err = connHandler.CreateNewUser(data.UserEmail, data.CaseSource); err != nil {
								m.logger.Send("error", supportingfunctions.CustomError(fmt.Errorf("%w, case id:'%d'", err, int(data.CaseId))).Error())
							} else {
								userAuthKey = us.AuthKey
								m.logger.Send("info", fmt.Sprintf("a new user '%s' has been successfully created", data.UserEmail))
							}
						}
					}

					switch data.Command {
					case "add event":
						//addEvent(m.host, authKey, m.authKey, data, m, m.logger)

						m.addNewObject(ctx, userAuthKey, data)

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

func (m *ModuleMISP) addNewObject(ctx context.Context, userAuthKey string, data InputSettings) {
	specialObject := NewCacheSpecialObject[*objectsmispformat.ListFormatsMISP]()
	specialObject.SetID(data.RootId)
	specialObject.SetObject(&data.Data)
	specialObject.SetFunc(func(i int) bool {

		/*

			сюда нужно перенести все функции необходимые для добавления данных в MISP

		*/

		//по умолчанию 'не успешно'
		return false
	})
}
