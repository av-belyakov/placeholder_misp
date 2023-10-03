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
			switch data.Command {
			case "send event id":
				// ***********************************
				// Это логирование только для теста!!!
				// ***********************************
				loging <- datamodels.MessageLoging{
					MsgData: fmt.Sprintf("TEST_INFO func 'CoreHandler', отправляем полученный event id: %s в модуль NATS", data.EventId),
					MsgType: "info",
				}
				//
				//

				//отправка eventId в NATS
				natsmodule.SendingDataInput(natsinteractions.SettingsInputChan{
					Command: data.Command,
					EventId: data.EventId,
				})

			case "set new event id":
				// ***********************************
				// Это логирование только для теста!!!
				// ***********************************
				loging <- datamodels.MessageLoging{
					MsgData: fmt.Sprintf("TEST_INFO func 'CoreHandler', надо отправить инфу CaseID '%s' и EventId '%s' to REDIS DB", data.CaseId, data.EventId),
					MsgType: "info",
				}
				//
				//

				//обработка запроса на добавления новой связки caseId:eventId в Redis
				redismodule.SendingDataInput(redisinteractions.SettingsChanInputRedis{
					Command: "set case id",
					Data:    fmt.Sprintf("%s:%s", data.CaseId, data.EventId),
				})
			}

		case data := <-redisChanReception:
			switch data.CommandResult {
			case "found event id":
				// ***********************************
				// Это логирование только для теста!!!
				// ***********************************
				loging <- datamodels.MessageLoging{
					MsgData: fmt.Sprintf("TEST_INFO func 'CoreHandler', здесь, получаем event id: '%v' из Redis для удаления события в MISP", data.Result),
					MsgType: "info",
				}
				//
				//

				// Здесь, получаем eventId из Redis для удаления события в MISP
				eventId, ok := data.Result.(string)
				if !ok {
					_, f, l, _ := runtime.Caller(0)

					loging <- datamodels.MessageLoging{
						MsgData: fmt.Sprintf("'it is not possible to convert a value to a string' %s:%d", f, l-1),
						MsgType: "error",
					}

					break
				}

				// ***********************************
				// Это логирование только для теста!!!
				// ***********************************
				loging <- datamodels.MessageLoging{
					MsgData: fmt.Sprintf("TEST_INFO func 'CoreHandler', отправляем event id: '%s' в MISP для удаления события", eventId),
					MsgType: "info",
				}
				//
				//

				mispmodule.SendingDataInput(mispinteractions.SettingsChanInputMISP{
					Command: "del event by id",
					EventId: eventId,
				})
			}
		}
	}
}
