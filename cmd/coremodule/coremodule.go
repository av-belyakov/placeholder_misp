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
	"context"
	"fmt"
	"runtime"

	"github.com/av-belyakov/placeholder_misp/cmd/commoninterfaces"
	"github.com/av-belyakov/placeholder_misp/cmd/mispapi"
	"github.com/av-belyakov/placeholder_misp/cmd/natsapi"
	"github.com/av-belyakov/placeholder_misp/cmd/redisapi"
	"github.com/av-belyakov/placeholder_misp/internal/datamodels"
	"github.com/av-belyakov/placeholder_misp/internal/supportingfunctions"
	"github.com/av-belyakov/placeholder_misp/memorytemporarystorage"
	rules "github.com/av-belyakov/placeholder_misp/rulesinteraction"
)

type CoreHandlerSettings struct {
	storageApp *memorytemporarystorage.CommonStorageTemporary
	logging    chan<- datamodels.MessageLogging
	counting   chan<- datamodels.DataCounterSettings
}

func NewCoreHandler(
	storage *memorytemporarystorage.CommonStorageTemporary,
	log chan<- datamodels.MessageLogging,
	count chan<- datamodels.DataCounterSettings) *CoreHandlerSettings {
	return &CoreHandlerSettings{
		storageApp: storage,
		logging:    log,
		counting:   count,
	}
}

func (settings *CoreHandlerSettings) CoreHandler(
	ctx context.Context,
	natsModule commoninterfaces.ModuleNatsHandler,
	mispModule commoninterfaces.ModuleMispHandler,
	redisModule *redisapi.ModuleRedis,
	listRule *rules.ListRule) error {

	chanNatsReception := natsModule.GetDataReceptionChannel()
	chanMispReception := mispModule.GetDataReceptionChannel()
	chanRedisReception := redisModule.GetDataReceptionChannel()
	hjm := NewHandlerJsonMessage(settings.storageApp, settings.logging, settings.counting)

	for {
		select {
		case <-ctx.Done():
			settings.logging <- datamodels.MessageLogging{
				MsgData: "TEST_INFO func 'CoreHandler', reseived ctx.Done()!!!!",
				MsgType: "testing",
			}

			return ctx.Err()

		case data := <-chanNatsReception:
			// ***********************************
			// Это логирование только для теста!!!
			// ***********************************
			settings.logging <- datamodels.MessageLogging{
				MsgData: "TEST_INFO func 'CoreHandler', reseived new object",
				MsgType: "testing",
			}
			//
			//

			//нужно только для хранения событий в RedisDB для последующей обработки
			//объектов которые не были добавлены в MISP из-за отсутствия доступа
			//к MISP (пока эта часть не реализована)
			//settings.storageApp.SetRawDataHiveFormatMessage(data.MsgId, data.Data)
			//добавляем raw данные по кейсу из thehive в Redis
			//redisModule.SendingDataInput(redisinteractions.SettingsChanInputRedis{
			//	Command: "set raw case",
			//	RawData: data.Data,
			//})

			//для записи необработанных событий в лог-файл events
			go func() {
				str, err := supportingfunctions.NewReadReflectJSONSprint(data.Data)
				if err != nil {
					// ***********************************
					// Это логирование только для теста!!!
					// ***********************************
					settings.logging <- datamodels.MessageLogging{
						MsgData: "TEST_INFO func 'CoreHandler', reseived new data",
						MsgType: "testing",
					}
					//
					//
				}

				if err == nil {
					settings.logging <- datamodels.MessageLogging{
						MsgData: fmt.Sprintf("\t---------------\n\tEVENTS:\n%s\n", str),
						MsgType: "events",
					}
				}
			}()

			// обработчик JSON документа
			chanOutputDecodeJson := hjm.HandlerJsonMessage(data.Data, data.MsgId)

			//формирование итоговых документов в формате MISP
			go NewMispFormat(chanOutputDecodeJson, data.MsgId, mispModule, listRule, settings.logging, settings.counting)

		case data := <-chanMispReception:
			switch data.Command {
			//отправка eventId в NATS
			case "send event id":

				// ***********************************
				// Это логирование только для теста!!!
				// ***********************************
				settings.logging <- datamodels.MessageLogging{
					MsgData: "TEST_INFO func 'CoreHandler', send EventId to ----> NATS",
					MsgType: "testing",
				}
				//
				//

				natsModule.SendingDataInput(natsapi.InputSettings{
					Command: data.Command,
					EventId: data.EventId,
					TaskId:  data.TaskId,
					RootId:  data.RootId,
					CaseId:  data.CaseId,
				})

			//отправка данных в Redis
			case "set new event id":
				//обработка запроса на добавления новой связки caseId:eventId в Redis
				redisModule.SendingDataInput(redisapi.SettingsChanInputRedis{
					Command: "set case id",
					Data:    fmt.Sprintf("%s:%s", data.CaseId, data.EventId),
				})
			}

		case data := <-chanRedisReception:
			switch data.CommandResult {
			//получаем eventId из Redis для удаления события в MISP
			case "found event id":
				eventId, ok := data.Result.(string)
				if !ok {
					_, f, l, _ := runtime.Caller(0)
					settings.logging <- datamodels.MessageLogging{
						MsgData: fmt.Sprintf("'it is not possible to convert a value to a string' %s:%d", f, l-1),
						MsgType: "error",
					}

					break
				}

				mispModule.SendingDataInput(mispapi.InputSettings{
					Command: "del event by id",
					EventId: eventId,
				})
			}
		}
	}
}
