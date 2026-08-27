package natsapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"golang.org/x/sync/errgroup"

	"github.com/av-belyakov/placeholder_misp/v2/constants"
	"github.com/av-belyakov/placeholder_misp/v2/internal/supportingfunctions"
)

// subscriptionCaseHandler обработчик подписок для получения кейсов
func (api *ApiNatsModule) subscriptionCaseHandler() {
	//event.object.caseId
	eventStruct := struct {
		Event struct {
			Object struct {
				CaseId int `json:"caseId"`
			} `json:"object"`
		} `json:"event"`
	}{}

	//приём кейсов
	api.natsConn.Subscribe(api.subscriptions.listenerCase, func(m *nats.Msg) {
		err := json.Unmarshal(m.Data, &eventStruct)
		if err != nil {
			fmt.Println("Error:", err)
		}

		api.logger.Send("info", fmt.Sprintf("a new case with id '%d' has been accepted", eventStruct.Event.Object.CaseId))

		api.SendingDataOutput(OutputSettings{
			MsgType: "case",
			MsgId:   uuid.NewString(),
			Data:    m.Data,
		})

		//счетчик принятых кейсов
		api.counting.SendMessage("update accepted events", 1)

	})

	lisSub := fmt.Sprintf("%v, listening to a subscription:%v'%s'%v", constants.Ansi_Bright_Green, constants.Ansi_Dark_Gray, api.subscriptions.listenerCase, constants.Ansi_Reset)
	log.Printf("%vconnect to NATS with address %v%s:%d%v%s\n", constants.Ansi_Bright_Green, constants.Ansi_Dark_Gray, api.host, api.port, constants.Ansi_Reset, lisSub)
}

// incomingInformationHandler обработчик информации полученной изнутри приложения
func (api *ApiNatsModule) incomingInformationHandler(ctx context.Context) {
	//обработка данных приходящих в модуль от ядра приложения фактически это команды на добавления
	//тега - 'add_case_tag' и команда на добавление MISP id в поле customField
	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case incomingData := <-api.GetChannelToModule():
				switch incomingData.Command {
				case "get sensor information":
					//
					// получение информации о сенсоре
					go func() {
						ctx, cancel := context.WithTimeout(ctx, constants.Default_Timeout_Get_Sensor_Information*time.Second)
						defer cancel()

						api.logger.Send("info", fmt.Sprintf("a request has been sent to get sensor information for an object with rootId:'%s', caseId:'%s'", incomingData.RootId, incomingData.CaseId))

						res, err := api.natsConn.RequestWithContext(ctx, api.subscriptions.getSensorInfo, incomingData.Data)
						if err != nil {
							api.logger.Send("error", supportingfunctions.CustomError(err).Error())
						}

						if res == nil {
							api.logger.Send("warning", fmt.Sprintf("an empty response to a request for additional sensor information was returned for an object with rootId:'%s', caseId:'%s'", incomingData.RootId, incomingData.CaseId))

							return
						}

						api.logger.Send("info", fmt.Sprintf("a response was received to a request for sensor information for an object with rootId:'%s', caseId:'%s'", incomingData.RootId, incomingData.CaseId))

						api.SendingDataOutput(OutputSettings{
							MsgType:       "sensor information",
							MsgId:         incomingData.CaseId,     // id кейса TheHive
							MsgAdditional: incomingData.CaseSource, // дополнительно отправляется источник для корректного поиска в БД
							Data:          res.Data,
						})
					}()

				case "send event id":
					//
					// установка тегов и customFields в TheHive

					//не отправляем eventId в TheHive
					if !api.sendCommand {
						continue
					}

					rootId := incomingData.RootId
					regionalObject := incomingData.CaseSource

					g := errgroup.Group{}
					g.Go(func() error {
						//команда на установку тега
						if err := api.natsConn.Publish(api.subscriptions.senderCommand,
							fmt.Appendf(
								nil,
								`{
					          "service": "MISP",
					          "command": "add_case_tag",
					  		  "for_regional_object": "%s",
					          "root_id": "%s",
					          "case_id": "%s",
					          "value": "Webhook: send=\"MISP\""
					        }`,
								regionalObject,
								rootId,
								incomingData.CaseId,
							)); err != nil {
							return err
						}

						return nil
					})
					g.Go(func() error {
						//команда на добавление значения поля customFields
						if err := api.natsConn.Publish(api.subscriptions.senderCommand,
							fmt.Appendf(
								nil,
								`{
						      "service": "MISP",
					          "command": "set_case_custom_field",
     					  	  "for_regional_object": "%s", 
							  "root_id": "%s",
					          "case_id": "%s",
					          "field_name": "misp-event-id.string",
					          "value": "%s"
						    }`,
								regionalObject,
								rootId,
								incomingData.CaseId,
								incomingData.EventId,
							)); err != nil {
							return err
						}

						return nil
					})

					if err := g.Wait(); err != nil {
						api.logger.Send("error", supportingfunctions.CustomError(err).Error())

						continue
					}

					api.logger.Send("info", fmt.Sprintf("comand:'%s' for case id:'%s' (root id:'%s') was successfully sent", incomingData.Command, incomingData.CaseId, incomingData.RootId))
				}
			}
		}
	}()
}
