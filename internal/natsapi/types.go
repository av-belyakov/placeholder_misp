package natsapi

import (
	"github.com/nats-io/nats.go"

	"github.com/av-belyakov/placeholder_misp/commoninterfaces"
)

// ApiNatsModule настройки для API NATS
type ApiNatsModule struct {
	natsConn      *nats.Conn
	subscriptions subscription
	host          string
	cachettl      int
	port          int
	logger        commoninterfaces.Logger
	counting      commoninterfaces.Counter
	chanOutput    chan OutputSettings //канал для отправки данных ИЗ модуля
	chanInput     chan InputSettings  //канал для приема данных В модуль
	sendCommand   bool
}

type subscription struct {
	listenerCase  string
	senderCommand string
	getSensorInfo string
}

// NatsApiOptions функциональные опции
type NatsApiOptions func(*ApiNatsModule) error

// OutputSettings параметры для канала отправки данных из модуля
type OutputSettings struct {
	Data    []byte
	MsgId   string
	MsgType string
}

// InputSettings параметры для канала приема данных в модуль
type InputSettings struct {
	Data       []byte //данные
	Command    string //команда
	TaskId     string //внутренний id задачи
	EventId    string //id события
	CaseId     string //id кейса
	RootId     string //rootId
	CaseSource string //источник кейса
}

type ResponseToCommand struct {
	Data       any    `json:"data"`
	ID         string `json:"id"`
	Error      string `json:"error"`
	Command    string `json:"command"`
	StatusCode int    `json:"status_code"`
}
