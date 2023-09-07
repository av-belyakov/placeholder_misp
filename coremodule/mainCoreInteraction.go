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
)

func CoreHandler(
	natsmodule *natsinteractions.ModuleNATS,
	mispmodule *mispinteractions.ModuleMISP,
	esmodule *elasticsearchinteractions.ModuleElasticSearch,
	nkckimodule *nkckiinteractions.ModuleNKCKI,
	listRule rules.ListRulesProcessingMsgMISP,
	storageApp *memorytemporarystorage.CommonStorageTemporary,
	loging chan<- datamodels.MessageLoging) {

	natsChanReception := natsmodule.GetDataReceptionChannel()
	mispChanReception := mispmodule.GetDataReceptionChannel()

	for {
		select {
		case data := <-natsChanReception:
			fmt.Printf("\n\tfunc 'CoreHandler', RESEIVED data from NATS, UUID = %s\n", data.UUID)

			storageApp.SetRawDataHiveFormatMessage(data.UUID, data.Data)

			//формирование итоговых документов в формате MISP
			chanCreateMispFormat, chanDone := NewMispFormat(data.UUID, mispmodule, loging)

			//обработчик сообщений из TheHive (выполняется разбор сообщения и его разбор на основе правил)
			go HandlerMessageFromHive(data.UUID, storageApp, listRule, chanCreateMispFormat, chanDone, loging)

			// отправка сообщения в Elasticshearch
			esmodule.SendingData(elasticsearchinteractions.SettingsInputChan{UUID: data.UUID})

			// отправка сообщения в НКЦКИ (пока заглушка)
			//nkckimodule.SendingData(procMsg.Message)

		case data := <-mispChanReception:
			fmt.Println("func 'NewCore', MISP reseived message from chanOutMISP: ", data)

			if data.Command == "send eventId" {
				fmt.Println("func 'NewCore', надо отправить инфу в NATS")

				natsmodule.SendingDataInput(natsinteractions.SettingsInputChan{
					Command: data.Command,
					EventId: data.EventId,
				})
			}
		}
	}
}
