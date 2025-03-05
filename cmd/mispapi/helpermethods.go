package mispapi

import "github.com/av-belyakov/placeholder_misp/commoninterfaces"

// NewLogWrite создаёт вспомогательный тип для логирования
func NewLogWrite(l commoninterfaces.Logger) *LogWrite {
	return &LogWrite{logger: l}
}

// Write логирование данных
func (lw *LogWrite) Write(msgType, msg string) bool {
	lw.logger.Send(msgType, msg)

	return true
}
