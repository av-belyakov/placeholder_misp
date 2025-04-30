package commoninterfaces

type Counter interface {
	SendMessage(string, int)
}

//************** логирование ***************

// Logger интерфейс для передачи логов через вспомогательный модуль, при этом
// принятые логи не только завписаваются в лог-файлы и выводяться в консоль (работа модуля
// simplelogger), но и передаются в систему мониторинга, такую как Zabbix
type Logger interface {
	GetChan() <-chan Messager
	Send(msgType, msgData string)
}

type Messager interface {
	GetType() string
	SetType(v string)
	GetMessage() string
	SetMessage(v string)
}

// WriterLoggingData интерфейс для записи логов напрямую через simplelogger
type WriterLoggingData interface {
	Write(typeLogFile, str string) bool
}
