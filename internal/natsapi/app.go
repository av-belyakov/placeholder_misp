// Пакет natsapi реализует методы для взаимодействия с NATS
package natsapi

import (
	"errors"

	"github.com/av-belyakov/placeholder_misp/v2/commoninterfaces"
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
