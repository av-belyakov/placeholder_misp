package coremodule

import (
	"fmt"
	"runtime"

	"placeholder_misp/datamodels"
	"placeholder_misp/elasticsearchinteractions"
	"placeholder_misp/memorytemporarystorage"
	"placeholder_misp/mispinteractions"
	"placeholder_misp/natsinteractions"
	"placeholder_misp/nkckiinteractions"
	"placeholder_misp/redisinteractions"
	rules "placeholder_misp/rulesinteraction"

	"github.com/google/uuid"
)

func CoreHandler(
	natsmodule *natsinteractions.ModuleNATS,
	mispmodule *mispinteractions.ModuleMISP,
	redismodule *redisinteractions.ModuleRedis,
	esmodule *elasticsearchinteractions.ModuleElasticSearch,
	nkckimodule *nkckiinteractions.ModuleNKCKI,
	listRule rules.ListRulesProcessingMsgMISP,
	storageApp *memorytemporarystorage.CommonStorageTemporary,
	loging chan<- datamodels.MessageLoging) {

	natsChanReception := natsmodule.GetDataReceptionChannel()
	mispChanReception := mispmodule.GetDataReceptionChannel()
	redisChanReception := redismodule.GetDataReceptionChannel()

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

			switch data.Command {
			case "send eventId":
				fmt.Println("func 'NewCore', надо отправить инфу в NATS")

				//отправка eventId в NATS
				natsmodule.SendingDataInput(natsinteractions.SettingsInputChan{
					Command: data.Command,
					EventId: data.EventId,
				})

			case "set new event id":
				//обработка запроса на добавления новой связки caseId:eventId в Redis
				redismodule.SendingDataInput(redisinteractions.SettingsChanInputRedis{
					Command: "set caseId",
					Data:    fmt.Sprintf("%s:%s", data.CaseId, data.EventId),
				})
			}

		case data := <-redisChanReception:
			fmt.Println("RESEIVED DATA FROM REDIS: ", data)

			switch data.CommandResult {
			case "found eventId":
				// Здесь, получаем eventId из Redis для удаления события в MISP
				eventId, ok := data.Result.(string)
				if !ok {
					_, f, l, _ := runtime.Caller(0)

					loging <- datamodels.MessageLoging{
						MsgData: fmt.Sprintf(" 'it is not possible to convert a value to a string' %s:%d", f, l-1),
						MsgType: "warning",
					}

					break
				}

				mispmodule.SendingDataInput(mispinteractions.SettingsChanInputMISP{
					Command: "del event by id",
					EventId: eventId,
				})
			}
		}
	}
}
