package natsinteractions

import (
	"context"
	"fmt"
	"runtime"

	"github.com/nats-io/nats.go"

	"placeholder_misp/confighandler"
	"placeholder_misp/datamodels"
)

var mnats ModuleNATS

// ModuleNATS инициализированный модуль
// ChanOutputMISP - канал для принятия данных из модуля
// ChanLogging - канал для отправки логов
type ModuleNATS struct {
	ChanOutputNATS chan []byte
	ChanLogging    chan<- datamodels.MessageLoging
}

func init() {
	mnats.ChanOutputNATS = make(chan []byte)
}

func NewClientNATS(
	ctx context.Context,
	conf confighandler.AppConfigNATS,
	chanLog chan<- datamodels.MessageLoging) (*ModuleNATS, error) {
	fmt.Println("func 'NewClientNATS', START...")

	mnats.ChanLogging = chanLog

	nc, err := nats.Connect(fmt.Sprintf("%s:%d", conf.Host, conf.Port))
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return &mnats, fmt.Errorf("%s %s:%d", fmt.Sprint(err), f, l-2)
	}

	fmt.Println("func 'NewClientNATS', STATUS:", nc.Stats())

	// Simple Async Subscriber
	/*nc.Subscribe("foo", func(msg *nats.Msg) {
		fmt.Printf("Received a message: %s\n", string(msg.Data))
	})

	// Simple Publisher
	nc.Publish("foo", []byte("Hello World"))
	time.Sleep(900 * time.Millisecond)
	nc.Publish("foo", []byte("Send messgae after sleep 2s"))

	fmt.Println("func 'NewClientNATS', END")*/

	nc.Subscribe("main_caseupdate", func(msg *nats.Msg) {
		mnats.ChanOutputNATS <- msg.Data
	})

	return &mnats, nil
}

func (mnats ModuleNATS) GetDataReceptionChannel() <-chan []byte {
	/*
		для формирование правильного сообщения об ошибке
		if err != nil {
			_, f, l, _ := runtime.Caller(0)

			mnats.ChanLogging <- datamodels.MessageLoging{
				MsgData: fmt.Sprintf("%s %s:%d", fmt.Sprint(err), f, l-2),
				MsgType: "error",
			}
		}
	*/

	return mnats.ChanOutputNATS
}

func (mnats ModuleNATS) GettingData() interface{} {
	/*
		для формирование правильного сообщения об ошибке
		if err != nil {
			_, f, l, _ := runtime.Caller(0)

			mnats.ChanLogging <- datamodels.MessageLoging{
				MsgData: fmt.Sprintf("%s %s:%d", fmt.Sprint(err), f, l-2),
				MsgType: "error",
			}
		}
	*/

	return <-mnats.ChanOutputNATS
}

func (cnats ModuleNATS) SendingData(data interface{}) {
	//ecmisp.chanInput <- data
}
