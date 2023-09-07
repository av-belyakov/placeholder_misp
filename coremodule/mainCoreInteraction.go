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
			uuidCase := uuid.New().String()
			fmt.Printf("\n\tfunc 'CoreHandler', RESEIVED data from NATS, UUID = %s\n", uuidCase)

			storageApp.SetRawDataHiveFormatMessage(uuidCase, data.Data)

			//формирование итоговых документов в формате MISP
			chanCreateMispFormat, chanDone := NewMispFormat(mispmodule, loging)

			//обработчик сообщений из TheHive (выполняется разбор сообщения и его разбор на основе правил)
			go HandlerMessageFromHive(data.Data, uuidCase, storageApp, listRule, chanCreateMispFormat, chanDone, loging)

			// отправка сообщения в Elasticshearch
			esmodule.SendingData(elasticsearchinteractions.SettingsInputChan{UUID: uuidCase})

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
