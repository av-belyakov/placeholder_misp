package natsinteractions

import (
	"context"
	"fmt"
	"runtime"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"

	"placeholder_misp/confighandler"
	"placeholder_misp/datamodels"
	"placeholder_misp/memorytemporarystorage"
)

var mnats ModuleNATS

type SettingsOutputChan struct {
	UUID string
}

type SettingsInputChan struct {
	EventId string
}

// ModuleNATS инициализированный модуль
// ChanOutputMISP - канал для отправки полученных данных из модуля
// chanInputNATS - канал для принятия данных в модуль
// ChanLogging - канал для отправки логов
type ModuleNATS struct {
	chanOutputNATS chan SettingsOutputChan
	chanInputNATS  chan SettingsInputChan
	ChanLogging    chan<- datamodels.MessageLoging
}

func init() {
	mnats.chanOutputNATS = make(chan SettingsOutputChan)
	mnats.chanInputNATS = make(chan SettingsInputChan)
}

func NewClientNATS(
	ctx context.Context,
	conf confighandler.AppConfigNATS,
	storageApp *memorytemporarystorage.CommonStorageTemporary,
	chanLog chan<- datamodels.MessageLoging) (*ModuleNATS, error) {
	mnats.ChanLogging = chanLog

	nc, err := nats.Connect(fmt.Sprintf("%s:%d", conf.Host, conf.Port))
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return &mnats, fmt.Errorf("%s %s:%d", err.Error(), f, l-2)
	}

	fmt.Println("func 'NewClientNATS', STATUS:", nc.Stats())

	// обработка данных приходящих в модуль от ядра приложения
	go func() {
		for data := range mnats.chanInputNATS {
			nrm := datamodels.NewResponseMessage()
			nrm.ResponseMessageAddNewCommand(datamodels.ResponseCommandForTheHive{
				Command: "setcustomfield",
				Name:    "misp-event-id.string",
				String:  data.EventId,
			})

			fmt.Println("func 'NewClientNATS', ResponseMessageFromMispToTheHave: ", nrm.GetResponseMessageFromMispToTheHave())

			//
			// Далее нужно сделать Unmarchal для ResponseMessageFromMispToTheHave и отправить
			// в TheHive через Webhook и NATS
			//

		}
	}()

	nc.Subscribe("main_caseupdate", func(msg *nats.Msg) {
		uuidTask := uuid.NewString()
		storageApp.SetRawDataHiveFormatMessage(uuidTask, msg.Data)

		//mnats.chanOutputNATS <- msg.Data
		mnats.chanOutputNATS <- SettingsOutputChan{
			UUID: uuidTask,
		}
	})

	return &mnats, nil
}

func (mnats ModuleNATS) GetDataReceptionChannel() <-chan SettingsOutputChan /*[]byte*/ {
	return mnats.chanOutputNATS
}
