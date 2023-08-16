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

// ModuleNATS инициализированный модуль
// ChanOutputMISP - канал для принятия данных из модуля
// ChanLogging - канал для отправки логов
type ModuleNATS struct {
	//chanOutputNATS chan []byte
	chanOutputNATS chan SettingsOutputChan
	chanInputNATS  chan []byte
	ChanLogging    chan<- datamodels.MessageLoging
}

func init() {
	//mnats.chanOutputNATS = make(chan []byte)
	mnats.chanOutputNATS = make(chan SettingsOutputChan)
	mnats.chanInputNATS = make(chan []byte)
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
		return &mnats, fmt.Errorf("%s %s:%d", fmt.Sprint(err), f, l-2)
	}

	fmt.Println("func 'NewClientNATS', STATUS:", nc.Stats())

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

func (cnats ModuleNATS) SendingData(data []byte) {
	cnats.chanInputNATS <- data
}
