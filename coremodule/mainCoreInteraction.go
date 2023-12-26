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

	natsChanReception := natsModule.GetDataReceptionChannel()
	mispChanReception := mispModule.GetDataReceptionChannel()
	redisChanReception := redisModule.GetDataReceptionChannel()

	//СТАРЫЙ обработчик
	//hmfh := NewHandlerMessageFromHive(storageApp, listRule, logging, counting)
	//НОВЫЙ обработчик
	decodeJson := NewDecodeJsonMessageSettings(listRule, logging, counting)

	for {
		select {
		case data := <-natsChanReception:
			storageApp.SetRawDataHiveFormatMessage(data.MsgId, data.Data)

			//добавляем raw данные по кейсу из thehive в Redis
			redisModule.SendingDataInput(redisinteractions.SettingsChanInputRedis{
				Command: "set raw case",
				RawData: data.Data,
			})

			//формирование итоговых документов в формате MISP
			//chanOutputDecodeJSON, chanDone := NewMispFormat(data.MsgId, mispModule, logging)

			//СТАРЫЙ обработчик сообщений из TheHive (выполняется разбор сообщения и его разбор на основе правил)
			//go hmfh.HandlerMessageFromHive(chanOutputDecodeJSON, data.Data, data.MsgId, chanDone)

			//******************** Внимание *******************
			//Надо проверить декодирование JSON документов и формирование
			//новых объектов в MISP формате. Внес некоторые изменения в обработчики
			//JSON сообщения и формирования нового MISP формата. Восновном изменения
			//касаются зеркальной, по отношению к старым, заменой формирования и
			//использования каналов. Это сделано из-за того что появился новый обработчик
			//формирующий верефицированный JSON объект (TheHive JSON) для его отправки
			//в Elasticsearch

			//НОВЫЙ обработчик сообщений из TheHive (выполняется разбор сообщения и его разбор на основе правил)
			chanOutputJsonDecode, chanDecodeDone := decodeJson.HandlerJsonMessage(data.Data, data.MsgId)
			//НОВЫЙ формирование итоговых документов в формате MISP
			go NewVerifiedMispFormat(chanOutputJsonDecode, chanDecodeDone, data.MsgId, mispModule, logging)

			//!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
			// отправка СЫРОГО сообщения в Elasticshearch
			//ПОКА ОСТАВЛЯЮ, ПОТОМ НАДО БУДЕТ УДАЛИТЬ
			esModule.HandlerData(elasticsearchinteractions.SettingsInputChan{
				UUID: data.MsgId,
				Data: data.Data,
			})

			// отправка сообщения в НКЦКИ (пока заглушка)
			//nkckiModule.SendingData(procMsg.Message)

		case data := <-mispChanReception:
			switch data.Command {
			//отправка eventId в NATS
			case "send event id":
				natsModule.SendingDataInput(natsinteractions.SettingsInputChan{
					Command: data.Command,
					EventId: data.EventId,
					TaskId:  data.TaskId,
				})

			//отправка данных в Redis
			case "set new event id":
				//обработка запроса на добавления новой связки caseId:eventId в Redis
				redisModule.SendingDataInput(redisinteractions.SettingsChanInputRedis{
					Command: "set case id",
					Data:    fmt.Sprintf("%s:%s", data.CaseId, data.EventId),
				})
			}

		case data := <-redisChanReception:
			switch data.CommandResult {
			//получаем eventId из Redis для удаления события в MISP
			case "found event id":
				eventId, ok := data.Result.(string)
				if !ok {
					_, f, l, _ := runtime.Caller(0)

					logging <- datamodels.MessageLogging{
						MsgData: fmt.Sprintf("'it is not possible to convert a value to a string' %s:%d", f, l-1),
						MsgType: "error",
					}

					break
				}

				mispModule.SendingDataInput(mispinteractions.SettingsChanInputMISP{
					Command: "del event by id",
					EventId: eventId,
				})
			}
		}
	}
}
