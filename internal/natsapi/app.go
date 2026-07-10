// Пакет natsapi реализует методы для взаимодействия с NATS
package natsapi

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/av-belyakov/placeholder_misp/v2/commoninterfaces"
	"github.com/av-belyakov/placeholder_misp/v2/internal/supportingfunctions"
)

// New настраивает новый модуль взаимодействия с API NATS
func New(logger commoninterfaces.Logger, counting commoninterfaces.Counter, opts ...NatsApiOptions) (*ApiNatsModule, error) {
	api := &ApiNatsModule{
		cachettl:    60,
		sendCommand: true,
		logger:      logger,
		counting:    counting,
		//прием запросов в NATS
		chanInput: make(chan InputSettings),
		//передача запросов из NATS
		chanOutput: make(chan OutputSettings),
	}

	for _, opt := range opts {
		if err := opt(api); err != nil {
			return api, err
		}
	}

	return api, nil
}

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

// WithHost метод устанавливает имя или ip адрес хоста API
func WithHost(v string) NatsApiOptions {
	return func(n *ApiNatsModule) error {
		if v == "" {
			return errors.New("the value of 'host' cannot be empty")
		}

		n.host = v

		return nil
	}
}

// WithPort метод устанавливает порт API
func WithPort(v int) NatsApiOptions {
	return func(n *ApiNatsModule) error {
		if v <= 0 || v > 65535 {
			return errors.New("an incorrect network port value was received")
		}

		n.port = v

		return nil
	}
}

// WithCacheTTL устанавливает время жизни для кэша хранящего функции-обработчики
// запросов к модулю
func WithCacheTTL(v int) NatsApiOptions {
	return func(th *ApiNatsModule) error {
		if v <= 10 || v > 86400 {
			return errors.New("the lifetime of a cache entry should be between 10 and 86400 seconds")
		}

		th.cachettl = v

		return nil
	}
}

// WithSubcriptionListenerCase устанавливает имя канала NATS который необходимо прослушивать для
// получения объектов типа 'case'
func WithSubcriptionListenerCase(v string) NatsApiOptions {
	return func(n *ApiNatsModule) error {
		if v == "" {
			return errors.New("the value of 'sender_case' cannot be empty")
		}

		n.subscriptions.listenerCase = v

		return nil
	}
}

// WithSubcriptionSenderCommand устанавливает имя канала NATS через которые будут передаваться
// команды для выполнения определенных действий в TheHive
func WithSubcriptionSenderCommand(v string) NatsApiOptions {
	return func(n *ApiNatsModule) error {
		if v == "" {
			return errors.New("the value of 'listener_command' cannot be empty")
		}

		n.subscriptions.senderCommand = v

		return nil
	}
}

// WithSubcriptionGetSensorInfo устанавливает имя канала NATS через которые будут отправляться
// запросы на получение информации о сенсорах
func WithSubcriptionGetSensorInfo(v string) NatsApiOptions {
	return func(n *ApiNatsModule) error {
		if v == "" {
			return errors.New("the value of 'get_sensor_info' cannot be empty")
		}

		n.subscriptions.getSensorInfo = v

		return nil
	}
}

// WithNotSendCommand запрещает отправлять команды в ответ на данные полученные через NATS
// по умолчанию отправка команд всегда разрешена
func WithNotSendCommand() NatsApiOptions {
	return func(n *ApiNatsModule) error {
		n.sendCommand = false

		return nil
	}
}
