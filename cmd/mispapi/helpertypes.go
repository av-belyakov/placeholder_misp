package mispapi

import "github.com/av-belyakov/placeholder_misp/commoninterfaces"

// LogWrite вспомогательный тип применяемый для логирования
type LogWrite struct {
	logger commoninterfaces.Logger
}
