package mispapi

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/av-belyakov/placeholder_misp/commoninterfaces"
	"github.com/av-belyakov/placeholder_misp/internal/supportingfunctions"
)

func addEvent(
	host string,
	authKey string,
	masterKey string,
	data InputSettings,
	mmisp *ModuleMISP,
	logger commoninterfaces.Logger) {

	// ***********************************
	// Это логирование только для теста!!!
	// ***********************************
	logger.Send("testing", fmt.Sprintf("TEST_INFO func 'HandlerMISP', send EVENTS to ----> MISP	USER EMAIL: %s, CaseId: %v", data.UserEmail, data.CaseId))
	//
	//

	//обработка только для события типа 'events'
	_, resBodyByte, err := sendEventsMispFormat(host, authKey, data)
	if err != nil {
		logger.Send("error", supportingfunctions.CustomError(err).Error())

		return
	}

	resMisp := RespMISP{}
	if err := json.Unmarshal(resBodyByte, &resMisp); err != nil {
		logger.Send("error", supportingfunctions.CustomError(err).Error())

		return
	}

	var eventId string
	for key, value := range resMisp.Event {
		if key == "id" {
			if str, ok := value.(string); ok {
				eventId = str

				break
			}
		}
	}

	if eventId == "" {
		logger.Send("error", supportingfunctions.CustomError(fmt.Errorf("the formation of events of the 'Attributes' type was not performed because the EventID is empty")).Error())

		return
	}

	// ***********************************
	// Это логирование только для теста!!!
	// ***********************************
	logger.Send("testing", fmt.Sprintf("TEST_INFO func 'HandlerMISP', send EVENTS REPORTS to ----> MISP	USER EMAIL: %s, CaseId: %v", data.UserEmail, data.CaseId))
	//
	//

	// добавляем event_reports
	if err := sendEventReportsMispFormat(host, authKey, eventId, data.CaseId); err != nil {
		logger.Send("error", supportingfunctions.CustomError(err).Error())
	}

	// ***********************************
	// Это логирование только для теста!!!
	// ***********************************
	logger.Send("testing", fmt.Sprintf("TEST_INFO func 'HandlerMISP', send data to ----> RedisDB USER EMAIL: %s, CaseId: %v", data.UserEmail, data.CaseId))
	//
	//

	//отправляем запрос для добавления в БД Redis, id кейса и нового события
	mmisp.SendingDataOutput(OutputSetting{
		Command: "set new event id",
		CaseId:  fmt.Sprint(data.CaseId),
		EventId: eventId,
	})

	// ***********************************
	// Это логирование только для теста!!!
	// ***********************************
	logger.Send("testing", fmt.Sprintf("TEST_INFO func 'HandlerMISP', send ATTRIBYTES to ----> MISP	USER EMAIL: %s, CaseId: %v", data.UserEmail, data.CaseId))
	//
	//

	//добавляем атрибуты
	_, _ = sendAttribytesMispFormat(host, authKey, eventId, data, logger)

	// ***********************************
	// Это логирование только для теста!!!
	// ***********************************
	logger.Send("testing", fmt.Sprintf("TEST_INFO func 'HandlerMISP', send OBJECTS to ----> MISP	USER EMAIL: %s, CaseId: %v", data.UserEmail, data.CaseId))
	//
	//

	// добавляем объекты
	_, _ = sendObjectsMispFormat(host, authKey, eventId, data, logger)

	//берем небольшой таймаут
	time.Sleep(4 * time.Second)

	// ***********************************
	// Это логирование только для теста!!!
	// ***********************************
	logger.Send("testing", fmt.Sprintf("TEST_INFO func 'HandlerMISP', send EVENT_TAGS to ----> MISP	USER EMAIL: %s, CaseId: %v", data.UserEmail, data.CaseId))
	//
	//

	// добавляем event_tags
	if err := sendEventTagsMispFormat(host, masterKey, eventId, data, logger); err != nil {
		logger.Send("error", supportingfunctions.CustomError(err).Error())
	}

	// ***********************************
	// Это логирование только для теста!!!
	// ***********************************
	logger.Send("testing", fmt.Sprintf("TEST_INFO func 'HandlerMISP', send PUBLISH to ----> MISP	USER EMAIL: %s, CaseId: %v", data.UserEmail, data.CaseId))
	//
	//

	//публикуем добавленное событие
	//masterKey нужен для публикации события так как пользователь
	//должен иметь более расшириные права чем могут иметь некоторые
	//обычные пользователи
	resMsg, err := sendRequestPublishEvent(host, masterKey, eventId)
	if err != nil {
		logger.Send("error", supportingfunctions.CustomError(err).Error())
	}
	if resMsg != "" {
		logger.Send("warning", resMsg)
	}

	// ***********************************
	// Это логирование только для теста!!!
	// ***********************************
	logger.Send("testing", fmt.Sprintf("TEST_INFO func 'HandlerMISP', send EventId to ----> CORE	USER EMAIL: %s, CaseId: %v", data.UserEmail, data.CaseId))
	//
	//

	/*
	   типы команд для передачи в NATS

	   requests := map[string][]byte{
	   					"add_case_tag": []byte(
	   						fmt.Sprintf(`{
	   							"service": "MISP",
	   							"command": "add_case_tag",
	   							"root_id": "%s",
	   							"case_id": "%s",
	   							"value": "Webhook: send=\"MISP\""}`,
	   							data.RootId,
	   							data.CaseId)),
	   					"set_case_custom_field": []byte(
	   						fmt.Sprintf(`{
	   							"service": "MISP",
	   							"command": "set_case_custom_field",
	   							"root_id": "%s",
	   	  						"field_name": "misp-event-id.string",
	   							"value": "%s"}`,
	   							data.RootId,
	   							data.EventId)),
	   				}

	*/

	//отправляем в ядро информацию по event Id
	mmisp.SendingDataOutput(OutputSetting{
		Command: "send event id",
		EventId: eventId,
		CaseId:  fmt.Sprint(data.CaseId),
		RootId:  data.RootId,
		TaskId:  data.TaskId,
	})
}

func delEventById(host, authKey, eventId string, logger commoninterfaces.Logger) {

	// ***********************************
	// Это логирование только для теста!!!
	// ***********************************
	logger.Send("testing", fmt.Sprintf("TEST_INFO func 'HandlerMISP', удаление события типа event, где event id: %s", eventId))
	//
	//

	// удаление события типа event
	_, err := delEventsMispFormat(host, authKey, eventId)
	if err != nil {
		logger.Send("error", supportingfunctions.CustomError(err).Error())
	}

	// ***********************************
	// Это логирование только для теста!!!
	// ***********************************
	logger.Send("testing", fmt.Sprintf("TEST_INFO func 'HandlerMISP', должно было быть успешно выполненно удаление события event id: %s", eventId))
	//
	//

	// ***********************************
	// Это логирование только для теста!!!
	// ***********************************
	logger.Send("testing", "TEST_INFO STOP")
	//
	//

	//
	//только для теста, для ОСТАНОВА
	//
	//mmisp.SendingDataOutput(SettingChanOutputMISP{
	//	Command: "TEST STOP",
	//})
}
