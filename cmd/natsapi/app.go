// Пакет natsapi реализует методы для взаимодействия с NATS
package natsapi

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/av-belyakov/placeholder_misp/commoninterfaces"
	"github.com/av-belyakov/placeholder_misp/constants"
	"github.com/av-belyakov/placeholder_misp/internal/confighandler"
	"github.com/av-belyakov/placeholder_misp/internal/countermessage"
	"github.com/av-belyakov/placeholder_misp/internal/supportingfunctions"
)

// NewClientNATS конструктор API NATS
func NewClientNATS(
	confNats confighandler.AppConfigNATS,
	confTheHive confighandler.AppConfigTheHive,
	counting *countermessage.CounterMessage,
	logger commoninterfaces.Logger) (*ModuleNATS, error) {

	var mnats ModuleNATS = ModuleNATS{
		chanOutput: make(chan OutputSettings),
		chanInput:  make(chan InputSettings),
	}

	//инициируем хранилище для дескрипторов сообщений NATS
	ns := NewStorageNATS()

	nc, err := nats.Connect(
		fmt.Sprintf("%s:%d", confNats.Host, confNats.Port),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(3*time.Second))
	if err != nil {
		return &mnats, supportingfunctions.CustomError(err)
	}

	//обработка разрыва соединения с NATS
	nc.SetDisconnectErrHandler(func(c *nats.Conn, err error) {
		logger.Send("error", supportingfunctions.CustomError(fmt.Errorf("the connection with NATS has been disconnected %w", err)).Error())
	})

	//обработка переподключения к NATS
	nc.SetReconnectHandler(func(c *nats.Conn) {
		logger.Send("info", "the connection to NATS has been re-established")
	})

	nc.Subscribe(confNats.Subscriptions.SenderCase, func(m *nats.Msg) {
		// ***********************************
		// Это логирование только для теста!!!
		// ***********************************
		logger.Send("testing", "------------|||||| TEST_INFO func 'NewClientNATS', reseived new object ||||||------------")
		//
		//

		mnats.chanOutput <- OutputSettings{
			MsgId: ns.setElement(m),
			Data:  m.Data,
		}

		//счетчик принятых кейсов
		counting.SendMessage("update accepted events", 1)
	})

	nc.Flush()

	log.Printf("%vconnect to NATS with address %v%s:%d%v\n", constants.Ansi_Bright_Green, constants.Ansi_Dark_Gray, confNats.Host, confNats.Port, constants.Ansi_Reset)

	// обработка данных приходящих в модуль от ядра приложения фактически это команды на добавления
	//тега - 'add_case_tag' и команда на добавление MISP id в поле customField
	go func() {
		for incomingData := range mnats.chanInput {
			//не отправляем eventId в TheHive
			if !confTheHive.Send {
				continue
			}

			//отправляем команды на установку тега и значения поля customFields
			go func() {
				info, err := SendRequestCommandExecute(nc, confNats.Subscriptions.ListenerCommand, incomingData)
				if err != nil {
					logger.Send("error", supportingfunctions.CustomError(err).Error())

					return
				}

				logger.Send("info", info)
			}()
		}
	}()

	return &mnats, nil
}
