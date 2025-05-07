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

	"github.com/av-belyakov/placeholder_misp/cmd/mispapi"
	"github.com/av-belyakov/placeholder_misp/cmd/natsapi"
	"github.com/av-belyakov/placeholder_misp/cmd/sqlite3api"
	"github.com/av-belyakov/placeholder_misp/commoninterfaces"
	rules "github.com/av-belyakov/placeholder_misp/internal/ruleshandler"
	"github.com/av-belyakov/placeholder_misp/internal/supportingfunctions"
)

type CoreHandlerSettings struct {
	logger    commoninterfaces.Logger
	counter   commoninterfaces.Counter
	listRules *rules.ListRule
}

func NewCoreHandler(
	counter commoninterfaces.Counter,
	listRules *rules.ListRule,
	logger commoninterfaces.Logger) *CoreHandlerSettings {
	return &CoreHandlerSettings{
		logger:    logger,
		listRules: listRules,
		counter:   counter,
	}
}

func (settings *CoreHandlerSettings) Start(
	ctx context.Context,
	natsModule *natsapi.ApiNatsModule,
	mispModule mispapi.ModuleMispHandler,
	sqlite3Module *sqlite3api.ApiSqlite3Module) {

	chanNatsReception := natsModule.GetChannelFromModule()
	chanMispReception := mispModule.GetReceptionChannel()

	hjson := NewHandlerJSON(settings.counter, settings.logger)

	generatorFormatMISP := NewGenerateObjectsFormatMISP(
		SettingsGenerateObjectsFormatMISP{
			MispModule:    mispModule,
			Sqlite3Module: sqlite3Module,
			ListRule:      settings.listRules,
			Counter:       settings.counter,
			Logger:        settings.logger,
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
					Command: data.Command,
					EventId: data.EventId,
					TaskId:  data.TaskId,
					RootId:  data.RootId,
					CaseId:  data.CaseId,
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
