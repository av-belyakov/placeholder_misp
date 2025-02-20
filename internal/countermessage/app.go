// Пакет countermassage подсчитывает входящие сообщения и отправляет их системе мониторинга
package countermessage

import (
	"context"
	"fmt"
	"time"

	"github.com/av-belyakov/placeholder_misp/commoninterfaces"
	"github.com/av-belyakov/placeholder_misp/internal/informationcountingstorage"
	"github.com/av-belyakov/placeholder_misp/internal/supportingfunctions"
)

// New конструктор счетчика сообщений
func New(ch chan<- commoninterfaces.Messager) *CounterMessage {
	return &CounterMessage{
		storage:  informationcountingstorage.NewInformationMessageCountingStorage(),
		chOutput: ch,
		chInput:  make(chan DataCounterMessage),
	}
}

// Handler обработчик подсчитывающий входящие данные
func (c *CounterMessage) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(c.chInput)

				return

			case data := <-c.chInput:
				d, h, m, s := supportingfunctions.GetDifference(c.storage.GetStartTime(), time.Now())
				patternTime := fmt.Sprintf("со старта приложения: дней %d, часов %d, минут %d, секунд %d", d, h, m, s)

				var count uint = uint(data.count)
				c.storage.Increase(data.msgType)
				if v, err := c.storage.GetCount(data.msgType); err == nil {
					count = v
				}

				var msg string = data.msgType
				switch data.msgType {
				case "update accepted events":
					msg = fmt.Sprintf("принято: %d, %s", count, patternTime)

				case "update processed events":
					msg = fmt.Sprintf("обработано: %d, %s", count, patternTime)

				case "update events meet rules":
					msg = fmt.Sprintf("соответствует правилам: %d, %s", count, patternTime)

				}

				message := NewSomeMessage()
				message.SetType("info")
				message.SetMessage(msg)

				c.chOutput <- message
			}
		}
	}()
}

// SendMessage отправка количественных значений счетчику сообщений
func (c *CounterMessage) SendMessage(msgType string, count int) {
	c.chInput <- DataCounterMessage{msgType: msgType, count: count}
}
