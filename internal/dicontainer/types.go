package dicontainer

import (
	"github.com/av-belyakov/placeholder_misp/commoninterfaces"
)

type diContainer struct {
	logger       Logger
	counter      Counter
	configer     Configer
	simpleLogger SimpleLogger

	db       DB
	dbLogger DbLogger
	nats     NatsConnecter
	misp     MispConnecter

	ch      chan commoninterfaces.Messager
	rootDir string
}
