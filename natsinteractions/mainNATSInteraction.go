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

		fmt.Println("func 'NewClientNATS', 111 ERROR...", err)

		return &enumChannels, err
	}

	fmt.Println("func 'NewClientNATS', STATUS:", nc.Stats())

	// Simple Publisher
	nc.Publish("foo", []byte("Hello World"))

	// Simple Async Subscriber
	nc.Subscribe("foo", func(msg *nats.Msg) {
		fmt.Printf("Received a message: %s\n", string(msg.Data))
	})

	/*nc.Subscribe("nkcki.notification.feedback", func(msg *nats.Msg) {
		fmt.Printf("Received a message: %s\n", string(msg.Data))
	})

	//msg, err := nc.RequestWithContext(ctx, "foo", []byte("bar"))
	/*go func(nc *nats.Conn) {
		for {
			nc.Subscribe("nkcki.notification.feedback", func(msg *nats.Msg) {
				fmt.Printf("Received a message: %s\n", string(msg.Data))
			})
		}
	}(nc)*/

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
