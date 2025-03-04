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
	"errors"
	"fmt"

	"github.com/av-belyakov/placeholder_misp/cmd/mispapi"
	"github.com/av-belyakov/placeholder_misp/cmd/natsapi"
	"github.com/av-belyakov/placeholder_misp/cmd/redisapi"
	"github.com/av-belyakov/placeholder_misp/commoninterfaces"
	"github.com/av-belyakov/placeholder_misp/internal/countermessage"
	rules "github.com/av-belyakov/placeholder_misp/internal/ruleshandler"
	"github.com/av-belyakov/placeholder_misp/internal/supportingfunctions"
)

type CoreHandlerSettings struct {
	logger    commoninterfaces.Logger
	listRules *rules.ListRule
	counting  *countermessage.CounterMessage
}

func NewCoreHandler(
	counting *countermessage.CounterMessage,
	listRules *rules.ListRule,
	logger commoninterfaces.Logger) *CoreHandlerSettings {
	return &CoreHandlerSettings{
		logger:    logger,
		listRules: listRules,
		counting:  counting,
	}
}

func (settings *CoreHandlerSettings) Start(
	ctx context.Context,
	natsModule *natsapi.ApiNatsModule,
	mispModule mispapi.ModuleMispHandler,
	redisModule *redisapi.ModuleRedis) {

	chanNatsReception := natsModule.GetChannelFromModule()
	chanMispReception := mispModule.GetReceptionChannel()
	chanRedisReception := redisModule.GetReceptionChannel()
	hjson := NewHandlerJSON(settings.counting, settings.logger)

	for {
		select {
		case <-ctx.Done():
			return

		case data := <-chanNatsReception:
			fmt.Println("func 'CoreHandlerSettings.Start', reseived new case")

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

			go func() {
				// обработчик JSON документа
				chanOutputDecodeJson := hjson.Start(data.Data, data.MsgId)

				//формирование итоговых документов в формате MISP
				go CreateObjectsFormatMISP(chanOutputDecodeJson, data.MsgId, mispModule, settings.listRules, settings.counting, settings.logger)
			}()

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
				redisModule.SendDataInput(redisapi.SettingsInput{
					Command: "set case id",
					Data:    fmt.Sprintf("%s:%s", data.CaseId, data.EventId),
				})
			}

		case data := <-chanRedisReception:
			//получаем eventId из Redis для удаления события в MISP
			if data.CommandResult != "found event id" {
				continue
			}

			eventId, ok := data.Result.(string)
			if !ok {
				settings.logger.Send("error", supportingfunctions.CustomError(errors.New("it is not possible to convert a value to a string")).Error())

				continue
			}

			mispModule.SendDataInput(mispapi.InputSettings{
				Command: "del event by id",
				EventId: eventId,
			})
		}
	}
}
