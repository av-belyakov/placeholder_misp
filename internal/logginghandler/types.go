package logginghandler

import "github.com/av-belyakov/placeholder_misp/commoninterfaces"

type LoggingChan struct {
	logChan chan commoninterfaces.Messager
}

// MessageLogging содержит информацию используемую при логировании
type MessageLogging struct {
	Message, Type string
}
