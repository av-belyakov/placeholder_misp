// Пакет coremodule является ядром/маршрутизатором приложения
//
// Пакет осуществляет декодирование JSON сообщения получаемого от TheHive,
// проверку значений на соответствие списку правил и формирование объекта
// соответствующего формату MISP.
// Кроме того пакет обеспечивает связь и взаимодействие между различными
// специализированными модулями, такими как например, модули взаимодействия
// с NATS, MISP и т.д.
package coremodule

import (
	"fmt"
	"runtime"

	"placeholder_misp/datamodels"
	"placeholder_misp/memorytemporarystorage"
	"placeholder_misp/mispinteractions"
	"placeholder_misp/natsinteractions"
	"placeholder_misp/redisinteractions"
	rules "placeholder_misp/rulesinteraction"
	"placeholder_misp/supportingfunctions"
)

func CoreHandler(
	natsModule *natsinteractions.ModuleNATS,
	mispModule *mispinteractions.ModuleMISP,
	redisModule *redisinteractions.ModuleRedis,
	listRule *rules.ListRule,
	storageApp *memorytemporarystorage.CommonStorageTemporary,
	logging chan<- datamodels.MessageLogging,
	counting chan<- datamodels.DataCounterSettings) {

	natsChanReception := natsModule.GetDataReceptionChannel()
	mispChanReception := mispModule.GetDataReceptionChannel()
	redisChanReception := redisModule.GetDataReceptionChannel()
	hjm := NewHandlerJsonMessage(storageApp, logging, counting)

	for {
		select {
		case data := <-natsChanReception:
			storageApp.SetRawDataHiveFormatMessage(data.MsgId, data.Data)

			//добавляем raw данные по кейсу из thehive в Redis
			redisModule.SendingDataInput(redisinteractions.SettingsChanInputRedis{
				Command: "set raw case",
				RawData: data.Data,
			})

			//для записи необработанных событий в лог-файл events
			if str, err := supportingfunctions.NewReadReflectJSONSprint(data.Data); err != nil {
				logging <- datamodels.MessageLogging{
					MsgData: fmt.Sprintf("\t---------------\n\tEVENTS:\n%s\n", str),
					MsgType: "events",
				}
			}

			// обработчик JSON документа
			chanOutputDecodeJson := hjm.HandlerJsonMessage(data.Data, data.MsgId)

			//формирование итоговых документов в формате MISP
			go NewMispFormat(chanOutputDecodeJson, data.MsgId, mispModule, listRule, logging, counting)

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
