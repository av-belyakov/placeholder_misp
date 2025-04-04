package nats_test

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/av-belyakov/placeholder_misp/internal/confighandler"
	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

var (
	conf *confighandler.ConfigApp
	nc   *nats.Conn
	b    []byte

	err error
)

func TestMain(m *testing.M) {
	if err = godotenv.Load("../../.env"); err != nil {
		log.Fatalln(err)
	}

	conf, err = confighandler.New("placeholder_misp")
	if err != nil {
		log.Fatalln(err)
	}

	nc, err = nats.Connect("192.168.9.208:4222",
		nats.MaxReconnects(-1),
		nats.ReconnectWait(3*time.Second))
	if err != nil {
		log.Fatalln(err)
	}

	// обработка разрыва соединения с NATS
	nc.SetDisconnectErrHandler(func(c *nats.Conn, err error) {
		log.Println(err)
	})

	// обработка переподключения к NATS
	nc.SetReconnectHandler(func(c *nats.Conn) {
		log.Println(err)
	})

	os.Exit(m.Run())
}

func TestSendCase(t *testing.T) {
	chDone := make(chan struct{})

	go func(nc *nats.Conn) {
		nc.Subscribe(conf.Subscriptions.SenderCommand, func(msg *nats.Msg) {
			fmt.Printf("Received command: %s\n", string(msg.Data))

			chDone <- struct{}{}
		})
	}(nc)

	b, err = os.ReadFile("../test_json/event_39100.json")
	assert.NoError(t, err)

	t.Log("GO_PHMISP_MAIN =", os.Getenv("GO_PHMISP_MAIN"))
	t.Log("conf.Subscriptions.ListenerCase:", conf.Subscriptions.ListenerCase)

	err = nc.Publish(conf.Subscriptions.ListenerCase, b)
	assert.NoError(t, err)

	fmt.Println("Before")

	<-chDone

	fmt.Println("After")

	nc.Close()
}
