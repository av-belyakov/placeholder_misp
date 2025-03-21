package natsapi

import "errors"

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

//******************* функции настройки опций natsapi ***********************

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

// WithListenerCase устанавливает имя канала NATS который необходимо прослушивать для
// получения объектов типа 'case'
func WithListenerCase(v string) NatsApiOptions {
	return func(n *ApiNatsModule) error {
		if v == "" {
			return errors.New("the value of 'sender_case' cannot be empty")
		}

		n.subscriptions.listenerCase = v

		return nil
	}
}

// WithSenderCommand устанавливает имя канала NATS через которые будут передаваться
// команды для выполнения определенных действий в TheHive
func WithSenderCommand(v string) NatsApiOptions {
	return func(n *ApiNatsModule) error {
		if v == "" {
			return errors.New("the value of 'listener_command' cannot be empty")
		}

		n.subscriptions.senderCommand = v

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
