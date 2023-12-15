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
)

type InputDataSender interface {
	TestSend()
}

func CoreHandler(
	natsModule *natsinteractions.ModuleNATS,
	mispModule *mispinteractions.ModuleMISP,
	redisModule *redisinteractions.ModuleRedis,
	esModule *elasticsearchinteractions.ModuleElasticSearch,
	nkckiModule *nkckiinteractions.ModuleNKCKI,
	listRule *rules.ListRule,
	storageApp *memorytemporarystorage.CommonStorageTemporary,
	logging chan<- datamodels.MessageLogging,
	counting chan<- datamodels.DataCounterSettings) {

	listRedisChanReception := map[string]InputDataSender{
		"found event id": mispModule.SendingDataInput,
	}

	natsChanReception := natsModule.GetDataReceptionChannel()
	mispChanReception := mispModule.GetDataReceptionChannel()
	redisChanReception := redisModule.GetDataReceptionChannel()
	hmfh := NewHandlerMessageFromHive(storageApp, listRule, logging, counting)

	for {
		select {
		case data := <-natsChanReception:
			//формирование итоговых документов в формате MISP
			chanCreateMispFormat, chanDone := NewMispFormat(data.MsgId, mispModule, logging)

			//обработчик сообщений из TheHive (выполняется разбор сообщения и его разбор на основе правил)
			go hmfh.HandlerMessageFromHive(chanCreateMispFormat, data.Data, data.MsgId, chanDone)

			// отправка сообщения в Elasticshearch
			esModule.HandlerData(elasticsearchinteractions.SettingsInputChan{
				UUID: data.MsgId,
				Data: data.Data,
			})

			// отправка сообщения в НКЦКИ (пока заглушка)
			//nkckiModule.SendingData(procMsg.Message)

		case data := <-mispChanReception:
			switch data.Command {
			//отправка данных в NATS
			case "send event id":
				// ***********************************
				// Это логирование только для теста!!!
				// ***********************************
				logging <- datamodels.MessageLogging{
					MsgData: fmt.Sprintf("TEST_INFO func 'CoreHandler', отправляем полученный event id: %s в модуль NATS", data.EventId),
					MsgType: "testing",
				}
				//
				//

				//отправка eventId в NATS
				natsModule.SendingDataInput(natsinteractions.SettingsInputChan{
					Command: data.Command,
					EventId: data.EventId,
					TaskId:  data.TaskId,
				})

				//отправка данных в Redis
			case "set new event id":
				// ***********************************
				// Это логирование только для теста!!!
				// ***********************************
				logging <- datamodels.MessageLogging{
					MsgData: fmt.Sprintf("TEST_INFO func 'CoreHandler', надо отправить инфу CaseID '%s' и EventId '%s' to REDIS DB", data.CaseId, data.EventId),
					MsgType: "testing",
				}
				//
				//

				//обработка запроса на добавления новой связки caseId:eventId в Redis
				redisModule.SendingDataInput(redisinteractions.SettingsChanInputRedis{
					Command: "set case id",
					Data:    fmt.Sprintf("%s:%s", data.CaseId, data.EventId),
				})
			}

		case data := <-redisChanReception:
			if f, ok := listRedisChanReception[data.CommandResult]; ok {
				f()
			}

			switch data.CommandResult {
			case "found event id":
				// ***********************************
				// Это логирование только для теста!!!
				// ***********************************
				logging <- datamodels.MessageLogging{
					MsgData: fmt.Sprintf("TEST_INFO func 'CoreHandler', здесь, получаем event id: '%v' из Redis для удаления события в MISP", data.Result),
					MsgType: "testing",
				}
				//
				//

				// Здесь, получаем eventId из Redis для удаления события в MISP
				eventId, ok := data.Result.(string)
				if !ok {
					_, f, l, _ := runtime.Caller(0)

					logging <- datamodels.MessageLogging{
						MsgData: fmt.Sprintf("'it is not possible to convert a value to a string' %s:%d", f, l-1),
						MsgType: "error",
					}

					break
				}

				// ***********************************
				// Это логирование только для теста!!!
				// ***********************************
				logging <- datamodels.MessageLogging{
					MsgData: fmt.Sprintf("TEST_INFO func 'CoreHandler', отправляем event id: '%s' в MISP для удаления события", eventId),
					MsgType: "testing",
				}
				//
				//

				mispModule.SendingDataInput(mispinteractions.SettingsChanInputMISP{
					Command: "del event by id",
					EventId: eventId,
				})
			}
		}
	}
}
