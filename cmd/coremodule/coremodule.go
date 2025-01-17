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

	"github.com/av-belyakov/placeholder_misp/cmd/mispapi"
	"github.com/av-belyakov/placeholder_misp/cmd/natsapi"
	"github.com/av-belyakov/placeholder_misp/cmd/redisapi"
	"github.com/av-belyakov/placeholder_misp/commoninterfaces"
	"github.com/av-belyakov/placeholder_misp/internal/countermessage"
	"github.com/av-belyakov/placeholder_misp/internal/supportingfunctions"
	rules "github.com/av-belyakov/placeholder_misp/rulesinteraction"
)

type CoreHandlerSettings struct {
	logger   commoninterfaces.Logger
	counting *countermessage.CounterMessage
}

func NewCoreHandler(counting *countermessage.CounterMessage, logger commoninterfaces.Logger) *CoreHandlerSettings {
	return &CoreHandlerSettings{
		logger:   logger,
		counting: counting,
	}
}

func (settings *CoreHandlerSettings) CoreHandler(
	ctx context.Context,
	natsModule natsapi.ModuleNatsHandler,
	mispModule mispapi.ModuleMispHandler,
	redisModule *redisapi.ModuleRedis,
	listRule *rules.ListRule) error {

	chanNatsReception := natsModule.GetDataReceptionChannel()
	chanMispReception := mispModule.GetDataReceptionChannel()
	chanRedisReception := redisModule.GetDataReceptionChannel()
	hjm := NewHandlerJsonMessage(settings.counting, settings.logger)

	for {
		select {
		case <-ctx.Done():
			settings.logger.Send("testing", "TEST_INFO func 'CoreHandler', reseived ctx.Done()!!!!")

			return ctx.Err()

		case data := <-chanNatsReception:
			// ***********************************
			// Это логирование только для теста!!!
			// ***********************************
			settings.logger.Send("testing", "TEST_INFO func 'CoreHandler', reseived new object")
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
					settings.logger.Send("testing", "TEST_INFO func 'CoreHandler', reseived new data")
					//
					//
				}

				if err == nil {
					settings.logger.Send("events", fmt.Sprintf("\t---------------\n\tEVENTS:\n%s\n", str))
				}
			}()

			// обработчик JSON документа
			chanOutputDecodeJson := hjm.HandlerJsonMessage(data.Data, data.MsgId)

			//формирование итоговых документов в формате MISP
			go NewMispFormat(chanOutputDecodeJson, data.MsgId, mispModule, listRule, settings.counting, settings.logger)

		case data := <-chanMispReception:
			switch data.Command {
			//отправка eventId в NATS
			case "send event id":

				// ***********************************
				// Это логирование только для теста!!!
				// ***********************************
				settings.logger.Send("testing", "TEST_INFO func 'CoreHandler', send EventId to ----> NATS")
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
					settings.logger.Send("error", fmt.Sprintf("'it is not possible to convert a value to a string' %s:%d", f, l-2))

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
