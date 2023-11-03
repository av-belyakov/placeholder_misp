package natsinteractions

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/nats-io/nats.go"

	"placeholder_misp/confighandler"
	"placeholder_misp/datamodels"
	"placeholder_misp/memorytemporarystorage"
)

var mnats ModuleNATS

func init() {
	mnats.chanOutputNATS = make(chan SettingsOutputChan)
	mnats.chanInputNATS = make(chan SettingsInputChan)
}

func NewClientNATS(
	conf confighandler.AppConfigNATS,
	storageApp *memorytemporarystorage.CommonStorageTemporary,
	logging chan<- datamodels.MessageLogging,
	counting chan<- datamodels.DataCounterSettings) (*ModuleNATS, error) {

	nc, err := nats.Connect(
		fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(3*time.Second))
	_, f, l, _ := runtime.Caller(0)

	if err != nil {
		return &mnats, fmt.Errorf("'%s' %s:%d", err.Error(), f, l-4)
	}

	//обработка разрыва соединения с NATS
	nc.SetDisconnectErrHandler(func(c *nats.Conn, err error) {
		logging <- datamodels.MessageLogging{
			MsgData: fmt.Sprintf("the connection with NATS has been disconnected %s:%d", f, l-4),
			MsgType: "error",
		}
	})

	nc.SetReconnectHandler(func(c *nats.Conn) {
		logging <- datamodels.MessageLogging{
			MsgData: fmt.Sprintf("the connection to NATS has been re-established %s:%d", f, l-4),
			MsgType: "info",
		}
	})

	log.Printf("Connect to NATS with address %s:%d\n", conf.Host, conf.Port)

	// обработка данных приходящих в модуль от ядра приложения
	go func() {
		for data := range mnats.chanInputNATS {
			nrm := datamodels.NewResponseMessage()

			if data.Command == "send eventId" {
				nrm.ResponseMessageAddNewCommand(datamodels.ResponseCommandForTheHive{
					Command: "setcustomfield",
					Name:    "misp-event-id.string",
					String:  data.EventId,
				})
			}

			// Далее нужно сделать Marchal для ResponseMessageFromMispToTheHave и отправить
			// в TheHive через Webhook и NATS
			res, err := json.Marshal(nrm.GetResponseMessageFromMispToTheHave())
			if err != nil {
				_, f, l, _ := runtime.Caller(0)

				logging <- datamodels.MessageLogging{
					MsgData: fmt.Sprintf("%s %s:%d", err.Error(), f, l-2),
					MsgType: "error",
				}

				continue
			}

			//err = nc.Publish("setcustomfield", res)
			err = nc.Publish("main_caseupdate", res)
			if err != nil {
				_, f, l, _ := runtime.Caller(0)

				logging <- datamodels.MessageLogging{
					MsgData: fmt.Sprintf("%s %s:%d", err.Error(), f, l-2),
					MsgType: "error",
				}
			}
		}
	}()

	nc.Subscribe("main_caseupdate", func(msg *nats.Msg) {
		mnats.chanOutputNATS <- SettingsOutputChan{
			Data: msg.Data,
		}

		//сетчик принятых кейсов
		counting <- datamodels.DataCounterSettings{
			DataType: "update accepted events",
			Count:    1,
		}

		//msg.Respond()
	})

	return &mnats, nil
}
