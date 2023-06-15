package natsinteractions

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"

	"placeholder_misp/confighandler"
)

var enumChannels EnumChannelsNATS

// EnumChannelsMISP перечень каналов для взаимодействия с модулем
// chanOutput - канал для принятия данных из модуля
type EnumChannelsNATS struct {
	chanOutput chan interface{}
}

func init() {
	enumChannels = EnumChannelsNATS{
		chanOutput: make(chan interface{}),
	}
}

func NewClientNATS(ctx context.Context, conf confighandler.AppConfigNATS) (*EnumChannelsNATS, error) {
	fmt.Println("func 'NewClientNATS', START...")

	/*
			Хотя context  может и не надо использовать?


			CLIENT_NATS_HOST=nats.cloud.gcm
		CLIENT_NATS_PORT=4222
		CLIENT_NATS_SUBJECT_PUB=nkcki.notification.feedback
		CLIENT_NATS_SUBJECT_SUB=nkcki.notification.feedback

	*/

	nc, err := nats.Connect(fmt.Sprintf("%s:%d", conf.Host, conf.Port))
	if err != nil {
		return &enumChannels, err
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

	nc.Subscribe("main_caseupdate" /*"nkcki.notification.feedback"*/, func(msg *nats.Msg) {
		//fmt.Printf("Received a message: %s\n", string(msg.Data))
		enumChannels.chanOutput <- msg.Data
	})

	return &enumChannels, nil
}

func (ecmisp EnumChannelsNATS) GetDataReceptionChannel() <-chan interface{} {
	return ecmisp.chanOutput
}

func (ecnats EnumChannelsNATS) GettingData() interface{} {
	return <-ecnats.chanOutput
}

func (ecnats EnumChannelsNATS) SendingData(data interface{}) {
	//ecmisp.chanInput <- data
}
