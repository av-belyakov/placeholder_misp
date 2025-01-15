package natsapi

// ModuleNATS инициализированный модуль
type ModuleNATS struct {
	chanOutput chan OutputSettings //канал для отправки данных ИЗ модуля
	chanInput  chan InputSettings  //канал для приема данных В модуль
}

// OutputSettings параметры для канала отправки данных из модуля
type OutputSettings struct {
	MsgId string
	Data  []byte
}

// InputSettings параметры для канала приема данных в модуль
type InputSettings struct {
	Command string //команда
	TaskId  string //внутренний id задачи
	EventId string //id события
	CaseId  string //id кейса
	RootId  string //rootId
	Data    []byte //данные
}

type ResponseToCommand struct {
	StatusCode int         `json:"status_code"`
	ID         string      `json:"id"`
	Error      string      `json:"error"`
	Command    string      `json:"command"`
	Data       interface{} `json:"data"`
}
