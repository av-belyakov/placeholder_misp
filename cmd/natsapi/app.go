// Пакет natsapi реализует методы для взаимодействия с NATS
package natsapi

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/av-belyakov/placeholder_misp/commoninterfaces"
	"github.com/av-belyakov/placeholder_misp/constants"
	"github.com/av-belyakov/placeholder_misp/internal/countermessage"
	"github.com/av-belyakov/placeholder_misp/internal/supportingfunctions"
)

// New настраивает новый модуль взаимодействия с API NATS
func New(logger commoninterfaces.Logger, counting *countermessage.CounterMessage, opts ...NatsApiOptions) (*ApiNatsModule, error) {
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

// Start инициализирует новый модуль взаимодействия с API NATS при инициализации
// возращается канал для взаимодействия с модулем, все запросы к модулю выполняются
// через данный канал
func (api *ApiNatsModule) Start(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	nc, err := nats.Connect(
		fmt.Sprintf("%s:%d", api.host, api.port),
		//nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(3*time.Second))
	if err != nil {
		return supportingfunctions.CustomError(err)
	}
	api.natsConn = nc

	//обработка разрыва соединения с NATS
	nc.SetDisconnectErrHandler(func(c *nats.Conn, err error) {
		api.logger.Send("error", supportingfunctions.CustomError(fmt.Errorf("the connection with NATS has been disconnected (%w)", err)).Error())
	})

	//обработка переподключения к NATS
	nc.SetReconnectHandler(func(c *nats.Conn) {
		api.logger.Send("info", "the connection to NATS has been re-established")
	})

	//приём кейсов
	nc.Subscribe(api.subscriptions.listenerCase, func(m *nats.Msg) {
		fmt.Printf("func 'NewClientNATS', reseived new case")

		// ***********************************
		// Это логирование только для теста!!!
		// ***********************************
		api.logger.Send("testing", "------------|||||| TEST_INFO func 'NewClientNATS', reseived new object ||||||------------")
		//
		//

		api.SendingDataOutput(OutputSettings{
			MsgId: ns.setElement(m),
			Data:  m.Data,
		})

		//счетчик принятых кейсов
		api.counting.SendMessage("update accepted events", 1)
	})

	nc.Flush()

	log.Printf("%vconnect to NATS with address %v%s:%d%v\n", constants.Ansi_Bright_Green, constants.Ansi_Dark_Gray, api.host, api.port, constants.Ansi_Reset)

	go func(ctx context.Context, nc *nats.Conn) {
		<-ctx.Done()

		fmt.Println("ApiNatsModule, stop 1")

		nc.Drain()
		nc.Close()
	}(ctx, nc)

	// обработка данных приходящих в модуль от ядра приложения фактически это команды на добавления
	//тега - 'add_case_tag' и команда на добавление MISP id в поле customField
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("ApiNatsModule, stop 2")

				return

			case incomingData := <-api.GetChannelToModule():
				//не отправляем eventId в TheHive
				if !api.sendCommand {
					continue
				}

				//отправляем команды на установку тега и значения поля customFields
				go func() {
					info, err := SendRequestCommandExecute(nc, api.subscriptions.senderCommand, incomingData)
					if err != nil {
						api.logger.Send("error", supportingfunctions.CustomError(err).Error())

						return
					}

					api.logger.Send("info", info)
				}()

			}
		}
	}()

	return nil
}
