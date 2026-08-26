package natsapi

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/av-belyakov/placeholder_misp/v2/internal/supportingfunctions"
)

// Start инициализирует новый модуль взаимодействия с API NATS
func (api *ApiNatsModule) Start(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	nc, err := nats.Connect(
		fmt.Sprintf("%s:%d", api.host, api.port),
		//nats.RetryOnFailedConnect(true),
		//имя клиента
		nats.Name("placeholder_misp"),
		//неограниченное количество попыток переподключения
		nats.MaxReconnects(-1),
		//время ожидания до следующей попытки переподключения (по умолчанию 2 сек.)
		nats.ReconnectWait(3*time.Second),
		//максимальное количество запросов ping, которые могут остаться без ответа от сервера,
		// прежде чем закрыть соединение
		nats.MaxPingsOutstanding(3),
		//обработка разрыва соединения с NATS
		nats.DisconnectErrHandler(func(c *nats.Conn, err error) {
			api.logger.Send("error", supportingfunctions.CustomError(fmt.Errorf("the connection with NATS has been disconnected (%w)", err)).Error())
		}),
		//обработка переподключения к NATS
		nats.ReconnectHandler(func(c *nats.Conn) {
			api.logger.Send("info", "the connection to NATS has been re-established")
		}))
	if err != nil {
		return supportingfunctions.CustomError(err)
	}
	api.natsConn = nc

	// обработчик подписок для получения кейсов
	go api.subscriptionCaseHandler()

	// обработчик информации полученной изнутри приложения
	go api.incomingInformationHandler(ctx)

	context.AfterFunc(ctx, func() {
		nc.Drain()
	})

	return nil
}

// GetChannelFromModule возвращает канал приема данных из модуля
func (api *ApiNatsModule) GetChannelFromModule() <-chan OutputSettings {
	return api.chanOutput
}

// GetChannelToModule возвращает канал приема данных в модуль
func (api *ApiNatsModule) GetChannelToModule() chan InputSettings {
	return api.chanInput
}

// SendingDataInput отправка данных в модуль
func (api *ApiNatsModule) SendingDataInput(data InputSettings) {
	api.chanInput <- data
}

// SendingDataOutput отправка данных из модуля
func (api *ApiNatsModule) SendingDataOutput(data OutputSettings) {
	api.chanOutput <- data
}
