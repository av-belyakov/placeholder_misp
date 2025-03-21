package sendcasetonats_test

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

const (
	Host string = "192.168.9.208"
	Port int    = 4222
)

var (
	conf  *confighandler.ConfigApp
	nc    *nats.Conn
	rbyte []byte

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

	rbyte, err = os.ReadFile("../test_json/event_39100.json")
	if err != nil {
		log.Panicln(err)
	}

	nc, err = nats.Connect(fmt.Sprintf("%s:%d", Host, Port),
		//nats.RetryOnFailedConnect(true),
		//имя клиента
		nats.Name("sender_case"),
		//неограниченное количество попыток переподключения
		nats.MaxReconnects(-1),
		//время ожидания до следующей попытки переподключения (по умолчанию 2 сек.)
		nats.ReconnectWait(3*time.Second))
	if err != nil {
		log.Panicln(err)
	}

	os.Exit(m.Run())
}

func TestSendToNats(t *testing.T) {
	nc.Subscribe(conf.Subscriptions.ListenerCase, func(msg *nats.Msg) {
		t.Log("received ", string(msg.Data))
	})

	fmt.Println("Sending 1 request with case ->")
	err = nc.Publish(conf.Subscriptions.ListenerCase, rbyte)
	assert.NoError(t, err)

	time.Sleep(time.Second * 3)

	fmt.Println("Sending 2 request with case ->")
	err = nc.Publish(conf.Subscriptions.ListenerCase, rbyte)
	assert.NoError(t, err)

	nc.Close()
}
