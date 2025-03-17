// Пакет natsapi реализует методы для взаимодействия с NATS
package natsapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
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

	//event.object.caseId
	eventStruct := struct {
		Event struct {
			Object struct {
				CaseId int `json:"caseId"`
			} `json:"object"`
		} `json:"event"`
	}{}

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

	//приём кейсов
	nc.Subscribe(api.subscriptions.listenerCase, func(m *nats.Msg) {
		err := json.Unmarshal(m.Data, &eventStruct)
		if err != nil {
			fmt.Println("Error:", err)
		}

		api.logger.Send("info", fmt.Sprintf("a new case with id '%d' has been accepted", eventStruct.Event.Object.CaseId))

		api.SendingDataOutput(OutputSettings{
			MsgId: uuid.NewString(),
			Data:  m.Data,
		})

		//счетчик принятых кейсов
		api.counting.SendMessage("update accepted events", 1)

	})

	lisSub := fmt.Sprintf("%v, listening to a subscription:%v'%s'%v", constants.Ansi_Bright_Green, constants.Ansi_Dark_Gray, api.subscriptions.listenerCase, constants.Ansi_Reset)
	log.Printf("%vconnect to NATS with address %v%s:%d%v%s\n", constants.Ansi_Bright_Green, constants.Ansi_Dark_Gray, api.host, api.port, constants.Ansi_Reset, lisSub)

	go func(ctx context.Context, nc *nats.Conn) {
		<-ctx.Done()
		nc.Drain()
	}(ctx, nc)

	// обработка данных приходящих в модуль от ядра приложения фактически это команды на добавления
	//тега - 'add_case_tag' и команда на добавление MISP id в поле customField
	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case incomingData := <-api.GetChannelToModule():
				//не отправляем eventId в TheHive
				if !api.sendCommand {
					continue
				}

				//отправляем команды на установку тега и значения поля customFields
				err := SendRequestCommandExecute(nc, api.subscriptions.senderCommand, incomingData)
				if err != nil {
					api.logger.Send("error", supportingfunctions.CustomError(err).Error())

					continue
				}

				api.logger.Send("info", fmt.Sprintf("comand:'%s' for case id:'%s' (root id:'%s') was successfully sent", incomingData.Command, incomingData.CaseId, incomingData.RootId))

			}
		}
	}()

	return nil
}
