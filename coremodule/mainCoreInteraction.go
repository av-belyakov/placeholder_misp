package coremodule

import (
	"fmt"

	"placeholder_misp/datamodels"
	"placeholder_misp/elasticsearchinteractions"
	"placeholder_misp/memorytemporarystorage"
	"placeholder_misp/mispinteractions"
	"placeholder_misp/natsinteractions"
	"placeholder_misp/nkckiinteractions"
	rules "placeholder_misp/rulesinteraction"

	"github.com/google/uuid"
)

func CoreHandler(
	natsmodule *natsinteractions.ModuleNATS,
	mispmodule *mispinteractions.ModuleMISP,
	elmodule *elasticsearchinteractions.ModuleElasticSearch,
	nkckimodule *nkckiinteractions.ModuleNKCKI,
	listRule rules.ListRulesProcessingMsgMISP,
	storageApp *memorytemporarystorage.CommonStorageTemporary,
	loging chan<- datamodels.MessageLoging) {

	natsChanReception := natsmodule.GetDataReceptionChannel()
	mispChanReception := mispmodule.GetDataReceptionChannel()

	for {
		select {
		case data := <-natsChanReception:
			uuidTask := uuid.NewString()
			storageApp.SetRawDataHiveFormatMessage(uuidTask, data)

			//формирование итоговых документов в формате MISP
			chanCreateMispFormat, chanDone := NewMispFormat(uuidTask, mispmodule, loging)

			//обработчик сообщений из TheHive (выполняется разбор сообщения и его разбор на основе правил)
			go HandlerMessageFromHive(uuidTask, storageApp, listRule, chanCreateMispFormat, chanDone, loging)

			// отправка сообщения в Elasticshearch (пока заглушка)
			//elmodule.SendingData(procMsg.Message)

			// отправка сообщения в НКЦКИ (пока заглушка)
			//nkckimodule.SendingData(procMsg.Message)

		case data := <-mispChanReception:
			fmt.Println("func 'NewCore', MISP reseived message from chanOutMISP: ", data)

		}
	}
}
