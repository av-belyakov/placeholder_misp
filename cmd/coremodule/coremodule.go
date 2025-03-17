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
	sqlite3Module *sqlite3api.ApiSqlite3Module) {

	chanNatsReception := natsModule.GetChannelFromModule()
	chanMispReception := mispModule.GetReceptionChannel()

	hjson := NewHandlerJSON(settings.counting, settings.logger)

	generatorFormatMISP := NewGenerateObjectsFormatMISP(
		SettingsGenerateObjectsFormatMISP{
			MispModule:    mispModule,
			Sqlite3Module: sqlite3Module,
			ListRule:      settings.listRules,
			Counter:       settings.counting,
			Logger:        settings.logger,
		})

	for {
		select {
		case <-ctx.Done():
			return

		case data := <-chanNatsReception:
			//для записи необработанных событий в лог-файл events
			go func() {
				str, err := supportingfunctions.NewReadReflectJSONSprint(data.Data)
				if err == nil {
					settings.logger.Send("events", fmt.Sprintf("\t---------------\n\tEVENTS:\n%s\n", str))
				}
			}()

			go func() {
				// обработчик JSON документа
				chanOutputDecodeJson := hjson.Start(data.Data, data.MsgId)

				//формирование итоговых документов в формате MISP
				generatorFormatMISP.Start(chanOutputDecodeJson, data.MsgId)
			}()

		case data := <-chanMispReception:
			switch data.Command {
			case "send event id":
				//отправка eventId в NATS
				natsModule.SendingDataInput(natsapi.InputSettings{
					Command: data.Command,
					EventId: data.EventId,
					TaskId:  data.TaskId,
					RootId:  data.RootId,
					CaseId:  data.CaseId,
				})

				//отправка eventId в Sqlite3
				sqlite3Module.SendDataToModule(sqlite3api.Request{
					Command: "set case id",
					Payload: fmt.Append(nil, fmt.Sprintf("%s:%s", data.CaseId, data.EventId)),
				})

			}
		}
	}
}
