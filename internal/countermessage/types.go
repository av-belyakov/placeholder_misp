package countermessage

import (
	"github.com/av-belyakov/placeholder_misp/v2/commoninterfaces"
	"github.com/av-belyakov/placeholder_misp/v2/internal/informationcountingstorage"
)

// CounterMessage счетчик сообщений
type CounterMessage struct {
	storage  *informationcountingstorage.InformationCountingStorage
	chInput  chan DataCounterMessage
	chOutput chan<- commoninterfaces.Messager
}

// SomeMessage некоторое сообщение
type SomeMessage struct {
	Type, Message string
}

// DataCounterMessage содержит информацию для подсчета
type DataCounterMessage struct {
	msgType string
	count   int
}
