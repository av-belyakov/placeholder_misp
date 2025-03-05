package sendcasetonats_test

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

const (
	Host string = "192.168.9.208"
	Port int    = 4222
)

var (
	rbyte []byte
	nc    *nats.Conn

	err error
)

func TestMain(m *testing.M) {
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
	nc.Subscribe("object.casetype.local", func(msg *nats.Msg) {
		t.Log("received ", string(msg.Data))
	})

	err = nc.Publish("object.casetype.local", rbyte)
	assert.NoError(t, err)

	nc.Close()
}
