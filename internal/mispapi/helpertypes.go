package mispapi

import "github.com/av-belyakov/placeholder_misp/v2/commoninterfaces"

// LogWrite вспомогательный тип применяемый для логирования
type LogWrite struct {
	logger commoninterfaces.Logger
}
