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

	"github.com/av-belyakov/placeholder_misp/commoninterfaces"
	"github.com/av-belyakov/placeholder_misp/internal/dicontainer"
	"github.com/av-belyakov/placeholder_misp/internal/mispapi"
	"github.com/av-belyakov/placeholder_misp/internal/natsapi"
	"github.com/av-belyakov/placeholder_misp/internal/sqlite3api"
	"github.com/av-belyakov/placeholder_misp/internal/supportingfunctions"
)

type CoreHandler struct {
	logger  commoninterfaces.Logger
	counter commoninterfaces.Counter
	rules   dicontainer.RulesHandler
}

func NewCoreHandler(
	logger commoninterfaces.Logger,
	counter commoninterfaces.Counter,
	rules dicontainer.RulesHandler,
) *CoreHandler {
	return &CoreHandler{
		rules:   rules,
		logger:  logger,
		counter: counter,
	}
}

func (settings *CoreHandler) Start(
	ctx context.Context,
	natsModule dicontainer.NatsConnecter,
	mispModule dicontainer.MispConnecter,
	sqlite3Module dicontainer.DB,
) {

	chanNatsReception := natsModule.GetChannelFromModule()
	chanMispReception := mispModule.GetReceptionChannel()

	hjson := NewHandlerJSON(settings.counter, settings.logger)

	generatorFormatMISP := NewGenerateObjectsFormatMISP(
		SettingsGenerateObjectsFormatMISP{
			MispModule:    mispModule,
			Sqlite3Module: sqlite3Module,
			ListRule:      settings.rules,
			Logger:        settings.logger,
			Counter:       settings.counter,
		})

	for {
		select {
		case <-ctx.Done():
			return

		case data := <-chanNatsReception:
			go func() {
				//----------------------------------------------------------------
				//----------- запись в файл необработанных объектов --------------
				//----------------------------------------------------------------
				str, err := supportingfunctions.NewReadReflectJSONSprint(data.Data)
				if err == nil {
					settings.logger.Send("events", fmt.Sprintf("\t---------------\n\tEVENTS:\n%s\n", str))
				}
				//----------------------------------------------------------------

				// обработчик JSON документа
				chanOutputDecodeJson := hjson.Start(data.Data, data.MsgId)

				//формирование итоговых документов в формате MISP
				generatorFormatMISP.Start(chanOutputDecodeJson, data.MsgId)
			}()

		case data := <-chanMispReception:
			switch data.Command {
			case "get event id":
				//отправка eventId в Sqlite3
				sqlite3Module.SendDataToModule(sqlite3api.Request{
					Command: "set case id",
					Payload: fmt.Append(nil, fmt.Sprintf("%s:%s", data.CaseId, data.EventId)),
				})

			case "send event id":
				//отправка eventId в NATS
				natsModule.SendingDataInput(natsapi.InputSettings{
					Command:    data.Command,
					EventId:    data.EventId,
					TaskId:     data.TaskId,
					RootId:     data.RootId,
					CaseId:     data.CaseId,
					CaseSource: data.CaseSource,
				})

				go func() {
					//поиск старого eventId в Sqlite3
					chRes := make(chan sqlite3api.Response)
					sqlite3Module.SendDataToModule(sqlite3api.Request{
						Command:    "search caseId",
						ChResponse: chRes,
						Payload:    fmt.Append(nil, data.CaseId),
					})
					res := <-chRes
					oldEeventId := string(res.Payload)

					if res.Error == nil && oldEeventId != "" {
						//запрос на удаление старого event в MISP
						mispModule.SendDataInput(mispapi.InputSettings{
							Command: "del event",
							EventId: oldEeventId,
						})
					}

					//передача нового eventId в Sqlite3
					sqlite3Module.SendDataToModule(sqlite3api.Request{
						Command: "set case id",
						Payload: fmt.Append(nil, fmt.Sprintf("%s:%s", data.CaseId, data.EventId)),
					})
				}()
			}
		}
	}
}
