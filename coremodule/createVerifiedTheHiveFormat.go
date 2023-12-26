package coremodule

import (
	"placeholder_misp/datamodels"
	"placeholder_misp/elasticsearchinteractions"
)

func NewVerifiedTheHiveFormat(
	input <-chan datamodels.ChanOutputDecodeJSON,
	done <-chan bool,
	esm *elasticsearchinteractions.ModuleElasticSearch,
	logging chan<- datamodels.MessageLogging,
) {

}
