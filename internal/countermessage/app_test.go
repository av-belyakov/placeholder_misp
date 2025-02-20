package countermessage_test

import (
	"context"
	"os"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/av-belyakov/placeholder_misp/commoninterfaces"
	"github.com/av-belyakov/placeholder_misp/internal/countermessage"
)

var (
	cm         *countermessage.CounterMessage
	chMessager chan commoninterfaces.Messager = make(chan commoninterfaces.Messager)
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
	ctx, close := context.WithCancel(context.Background())

	reg, err := regexp.Compile(`\d`)
	assert.NoError(t, err)

	cm.Start(ctx)
	go func() {
		var num int = 1
		for {
			select {
			case <-ctx.Done():
				return

			case data, isOpen := <-chMessager:
				if isOpen {
					if num == 5 {
						num = 1
					}

					t.Log(data.GetMessage())

					assert.Equal(t, data.GetType(), "info")
					str := reg.FindString(data.GetMessage())
					assert.Equal(t, str, strconv.Itoa(num))

					num++
				}
			}
		}
	}()

	cm.SendMessage("update accepted events", 1)
	cm.SendMessage("update accepted events", 1)
	cm.SendMessage("update accepted events", 1)
	cm.SendMessage("update accepted events", 1)

	time.Sleep(1 * time.Second)

	cm.SendMessage("update processed events", 1)
	cm.SendMessage("update processed events", 1)

	time.Sleep(1 * time.Second)

	close()

	assert.True(t, true)
	// "update accepted events"
	// "update processed events"
	// "update events meet rules"
}
