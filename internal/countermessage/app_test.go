package countermessage_test

import (
	"context"
	"os"
	"testing"

	"github.com/av-belyakov/placeholder_misp/commoninterfaces"
	"github.com/av-belyakov/placeholder_misp/internal/countermessage"
	"github.com/stretchr/testify/assert"
)

var (
	cm         *countermessage.CounterMessage
	chMessager chan commoninterfaces.Messager = make(chan commoninterfaces.Messager, 10)
)

type MessageSettings struct {
	msgType, msg string
}

func NewMessageSettinngs() *MessageSettings {
	return &MessageSettings{}
}

func (ms *MessageSettings) GetType() string {
	return ms.msgType
}

func (ms *MessageSettings) SetType(v string) {
	ms.msgType = v
}

func (ms *MessageSettings) GetMessage() string {
	return ms.msg
}

func (ms *MessageSettings) SetMessage(v string) {
	ms.msg = v
}

func TestMain(m *testing.M) {
	cm = countermessage.New(chMessager)

	os.Exit(m.Run())
}

func TestCounterMessage(t *testing.T) {
	go cm.Handler(context.Background())

	cm.SendMessage("update accepted events", 1)
	cm.SendMessage("update accepted events", 1)

	data := <-chMessager
	t.Log(data.GetMessage())
	assert.Equal(t, data.GetType(), "update accepted events")
	data = <-chMessager
	t.Log(data.GetMessage())
	assert.Equal(t, data.GetType(), "update accepted events")

	// "update accepted events"
	// "update processed events"
	// "update events meet rules"
}
