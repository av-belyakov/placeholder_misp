package dicontainer

import (
	"github.com/av-belyakov/placeholder_misp/v2/commoninterfaces"
)

// DiContainer DI контейнер
type DiContainer struct {
	logger       Logger
	counter      Counter
	configer     Configer
	simpleLogger SimpleLogger
	rules        RulesHandler

	db       DB
	dbLogger DbLogger
	nats     NatsConnecter
	misp     MispConnecter

	ch      chan commoninterfaces.Messager
	rootDir string
}
