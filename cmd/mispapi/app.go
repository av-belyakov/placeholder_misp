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

	l := NewLogWrite(logger)
	cache, err := cachingstoragewithqueue.NewCacheStorage(
		cachingstoragewithqueue.WithMaxTtl[*objectsmispformat.ListFormatsMISP](3600),
		cachingstoragewithqueue.WithMaxSize[*objectsmispformat.ListFormatsMISP](20),
		cachingstoragewithqueue.WithLogging[*objectsmispformat.ListFormatsMISP](l),
		cachingstoragewithqueue.WithTimeTick[*objectsmispformat.ListFormatsMISP](7),
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
		logger:       logger,
	}

	return module, nil
}

func (m *ModuleMISP) Start(ctx context.Context) error {
	client, err := NewClientMISP(m.host, m.authKey, false)
	if err != nil {
		return err
	}

	connHandler := NewHandlerAuthorizationMISP(client, NewStorageAuthorizationDataMISP())
	err = connHandler.GetListAllOrganisation(ctx, m.organistions)
	if err != nil {
		return err
	}

	countUser, err := connHandler.GetListAllUsers(ctx)
	if err != nil {
		return err
	}

	//инициализируем автоматическую обработку объектов попадающих в кеш
	m.cache.StartAutomaticExecution(ctx)

	m.logger.Send("info", fmt.Sprintf("at the start of the application, all user information was downloaded %d", countUser))

	//здесь обрабатываем данные из входного канала модуля MISP
	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case msg := <-m.GetInputChannel():
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
								m.logger.Send("info", fmt.Sprintf("a new user '%s' from %s has been successfully created", data.UserEmail, data.CaseSource))
							}
						}
					}

					switch data.Command {
					case "add event":
						m.addNewObject(ctx, userAuthKey, data)

					case "del event":
						m.delObject(ctx, data.EventId)

					}
				}(msg)
			}
		}
	}()

	return nil
}
