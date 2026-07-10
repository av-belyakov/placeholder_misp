package mispapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/av-belyakov/objectsmispformat"

	"github.com/av-belyakov/placeholder_misp/v2/internal/responses"
	"github.com/av-belyakov/placeholder_misp/v2/internal/supportingfunctions"
)

func (m *ModuleMISP) processingEvent(ctx context.Context, userAuthKey string, data InputSettings) {
	// проверяем наличие идентификатора события в MISP
	if data.EventId == "" {
		m.addNewEvent(ctx, userAuthKey, data)

		return
	}

	rmisp, err := NewMispRequest(
		WithHost(m.host),
		WithUserAuthKey(userAuthKey),
		WithMasterAuthKey(m.authKey))
	if err != nil {
		m.logger.Send("error", supportingfunctions.CustomError(err).Error())

		return
	}

	// получаем событие по его event id
	statusCode, raw, err := rmisp.getEvent(ctx, data.EventId)
	if err != nil {
		m.logger.Send("error", supportingfunctions.CustomError(err).Error())
	}

	fmt.Println("method 'processingEvent', get event status code:", statusCode)

	if statusCode == http.StatusNotFound {
		m.addNewEvent(ctx, userAuthKey, data)

		return
	}

	var oldEvent EventFromMISP
	if err = json.Unmarshal(raw, &oldEvent); err != nil {
		m.logger.Send("error", supportingfunctions.CustomError(err).Error())

		return
	}

	m.editEvent(ctx, userAuthKey, data, oldEvent)
}

// addNewEvent добавление нового события
func (m *ModuleMISP) addNewEvent(ctx context.Context, userAuthKey string, data InputSettings) {
	specialObject := NewCacheSpecialObject[*objectsmispformat.ListFormatsMISP]()
	specialObject.SetID(data.RootId)
	specialObject.SetObject(&data.Data)
	specialObject.SetFunc(func(i int) bool {
		rmisp, err := NewMispRequest(
			WithHost(m.host),
			WithUserAuthKey(userAuthKey),
			WithMasterAuthKey(m.authKey))
		if err != nil {
			m.logger.Send("error", supportingfunctions.CustomError(err).Error())

			return false
		}

		m.logger.Send("info", fmt.Sprintf("starting adding the case id:'%s'", data.CaseId))

		// отправляет в API MISP событие типа Event и возвращает результат который содержит
		// id события в MISP, у MISP свой уникальный id для событий
		// только с использованием этого id в MISP добавляются все остальные объекты
		_, resBodyByte, err := rmisp.sendEvent(ctx, data.Data.GetEvent())
		if err != nil {
			m.logger.Send("error", supportingfunctions.CustomError(err).Error())

			return false
		}

		// Все ошибки которые могут возникнуть при дальнейшем взаимодействии с MISP
		// будут попрежнему логироватся.
		// Однако, статус выполнения для функции будет ставится в TRUE, что бы не досить
		// MISP, так как все последующие попытки будут начинатся с добавления 'event', а
		// добавить 'event' с таким id нельзя.
		// Необходимо удалить предыдущий.

		resMisp := MispResponse{}
		if err := json.Unmarshal(resBodyByte, &resMisp); err != nil {
			m.logger.Send("error", supportingfunctions.CustomError(err).Error())

			return true
		}

		// получаем уникальный id MISP
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
			m.logger.Send("error", supportingfunctions.CustomError(fmt.Errorf("the formation of events of the 'Attributes' type was not performed because the EventID is empty")).Error())

			return true
		}

		m.logger.Send("info", fmt.Sprintf("new element 'event' with id:'%s' successfully created (case id:'%s')", eventId, data.CaseId))

		// добавляем event_reports
		if err := rmisp.sendEventReports(ctx, eventId, data.Data.GetReports()); err != nil {
			m.logger.Send("error", supportingfunctions.CustomError(err).Error())

			return true
		}

		m.logger.Send("info", fmt.Sprintf("element 'event_reports' successfully added to event with id:'%s' (case id:'%s')", eventId, data.CaseId))

		// добавляем атрибуты
		warning, err := rmisp.sendAttribytes(ctx, eventId, data.Data.GetAttributes())
		if err != nil {
			// Тут ошибка может быть при добавлении только одного из многих объектов,
			// соответственно тормозить весь процесс, только из-за того, что была ошибка
			// при добавлении одного или нескольких объектов не стоит.
			// Если же не был добавлен ни один из объектов, то это возможно глабальная
			// ошибка доступа, следовательно, при добавлении следующих объектов она также
			// может вылезти, тогда там и будет выполнен останов всей цепочки
			m.logger.Send("error", supportingfunctions.CustomError(err).Error())
		}
		if warning != "" {
			m.logger.Send("warning", warning)
		}

		m.logger.Send("info", fmt.Sprintf("some elements 'attribytes' successfully added to event with id:'%s' (case id:'%s')", eventId, data.CaseId))

		// добавляем объекты
		if err = rmisp.sendObjects(ctx, eventId, data.Data.GetObjects()); err != nil {
			// тут такая же ситуация что и с ошибками при выполнении метода rmisp.sendAttribytes
			m.logger.Send("error", supportingfunctions.CustomError(err).Error())
		}

		m.logger.Send("info", fmt.Sprintf("elements 'objects' successfully added to event with id:'%s' (case id:'%s')", eventId, data.CaseId))

		// берем небольшой таймаут, нужен для того что бы MISP успел обработать и добавить в БД
		// всё ранее ему переданное, если обработка переданных объектов не была завершена
		// возможны накладки или сбои при добавлении данных
		// это недостаток MISP, с этим я ничего не могу поделать
		time.Sleep(5 * time.Second)

		// добавляем event_tags
		if err := rmisp.sendEventTags(ctx, eventId, data.Data.GetObjectTags()); err != nil {
			m.logger.Send("error", supportingfunctions.CustomError(err).Error())
		}

		m.logger.Send("info", fmt.Sprintf("elements 'tags' successfully added to event with id:'%s' (case id:'%s')", eventId, data.CaseId))

		time.Sleep(5 * time.Second)

		// публикуем добавленное событие
		// masterKey нужен для публикации события так как пользователь должен иметь более расшириные права
		// чем могут иметь некоторые обычные пользователи
		resMsg, err := rmisp.sendRequestPublishEvent(ctx, eventId)
		if err != nil {
			m.logger.Send("error", supportingfunctions.CustomError(err).Error())
		}
		if resMsg != "" {
			m.logger.Send("info", fmt.Sprintf("event with id:'%s' (case id:'%s') %s", eventId, data.CaseId, resMsg))
		}

		// отправляем в ядро информацию по event Id, при этом новый eventId
		// передаётся для отправки в NATS, а так же передается в Sqlite3 для
		// обновления или создания новой связки caseId - eventId
		outMsg := OutputSetting{
			Command:    "send event id",
			EventId:    eventId,
			CaseId:     data.CaseId,
			RootId:     data.RootId,
			TaskId:     data.TaskId,
			CaseSource: data.CaseSource,
		}
		m.SendDataOutput(outMsg)

		// получаем из тегов список сенсоров по которым нужна дополнитеьльная информация
		listSensors := createListSensors(*data.Data.ObjectTags)
		if len(listSensors) == 0 {
			m.logger.Send("error", supportingfunctions.CustomError(fmt.Errorf("the sensor list  on object tags '%+v' for event with id:'%s' (case id:'%s') is empty", data.Data.ObjectTags, eventId, data.CaseId)).Error())

			return true
		}

		reqSensorId, err := json.Marshal(RequestSensorInformation{
			Source:      "placeholder_misp",
			TaskId:      data.CaseId,
			ListSensors: listSensors,
		})
		if err != nil {
			m.logger.Send("error", supportingfunctions.CustomError(err).Error())

			return true
		}

		// запрос в ядро на получение информации о сенсорах
		outMsg.Command = "get sensor information"
		outMsg.Data = reqSensorId
		m.SendDataOutput(outMsg)

		m.logger.Send("info", fmt.Sprintf("event with id:'%s' (case id:'%s') send request 'get sensor information' to core module", eventId, data.CaseId))

		return true
	})

	// добавляем вспомогательный тип specialObject в очередь хранилища
	m.cache.PushObjectToQueue(specialObject)
}

// editEvent редактирование существующего события
func (m *ModuleMISP) editEvent(ctx context.Context, userAuthKey string, data InputSettings, oldEvent EventFromMISP) {
	specialObject := NewCacheSpecialObject[*objectsmispformat.ListFormatsMISP]()
	specialObject.SetID(data.RootId)
	specialObject.SetObject(&data.Data)
	specialObject.SetFunc(func(i int) bool {
		/*oldEvent := struct {
			Event struct {
				UUID      string `json:"uuid"`
				Attribute []struct {
					Id string `json:"id"`
				} `json:"Attribute"`
				Object []struct {
					Id string `json:"id"`
				} `json:"Object"`
				Tag []struct {
					Id     string `json:"id"`
					Name   string `json:"name"`
					Colour string `json:"colour"`
				} `json:"tag"`
			} `json:"Event"`
		}{}*/

		rmisp, err := NewMispRequest(
			WithHost(m.host),
			WithUserAuthKey(userAuthKey),
			WithMasterAuthKey(m.authKey))
		if err != nil {
			m.logger.Send("error", supportingfunctions.CustomError(err).Error())

			return true
		}

		m.logger.Send("info", fmt.Sprintf("starting editing the event id:'%s', case id:'%s'", data.EventId, data.CaseId))

		// получаем событие по его event id
		/*_, raw, err := rmisp.getEvent(ctx, data.EventId)
		if err != nil {
			m.logger.Send("error", supportingfunctions.CustomError(err).Error())

			return true
		}
		if err = json.Unmarshal(raw, &oldEvent); err != nil {
			m.logger.Send("error", supportingfunctions.CustomError(err).Error())

			return true
		}*/

		updateEvent := objectsmispformat.EventsMispFormat{
			OrgId:              data.Data.Event.OrgId,
			OrgcId:             data.Data.Event.OrgcId,
			Info:               data.Data.Event.Info,
			Uuid:               data.Data.Event.Uuid,
			Date:               data.Data.Event.Date,
			Analysis:           data.Data.Event.Analysis,
			Distribution:       data.Data.Event.Distribution,
			AttributeCount:     data.Data.Event.AttributeCount,
			SharingGroupId:     data.Data.Event.SharingGroupId,
			ThreatLevelId:      data.Data.Event.ThreatLevelId,
			Timestamp:          data.Data.Event.Timestamp, // ВЫЖНЫЙ ПАРАМЕТР, если это значение меньше значения в MISP событие не будет обноавлено
			PublishTimestamp:   data.Data.Event.PublishTimestamp,
			SightingTimestamp:  data.Data.Event.SightingTimestamp,
			ExtendsUuid:        data.Data.Event.ExtendsUuid,
			EventCreatorEmail:  data.Data.Event.EventCreatorEmail,
			Published:          data.Data.Event.Published,
			ProposalEmailLock:  data.Data.Event.ProposalEmailLock,
			Locked:             data.Data.Event.Locked,
			DisableCorrelation: data.Data.Event.DisableCorrelation,
		}

		// обновляем существующее событие
		_, _, err = rmisp.editEvent(ctx, data.EventId, updateEvent)
		if err != nil {
			m.logger.Send("error", supportingfunctions.CustomError(err).Error())

			return true
		}

		// event_reports оставляем неизменным, этот объект не должен изменятся

		// удаляем все атрибуты события
		for _, v := range oldEvent.Event.Attribute {
			if err := rmisp.deleteAttributes(ctx, v.Id, data.EventId); err != nil {
				m.logger.Send("error", supportingfunctions.CustomError(err).Error())
			}
		}
		// добавляем все атрибуты повторно
		warning, err := rmisp.sendAttribytes(ctx, data.EventId, data.Data.GetAttributes())
		if err != nil {
			m.logger.Send("error", supportingfunctions.CustomError(err).Error())
		}
		if warning != "" {
			m.logger.Send("warning", warning)
		}

		m.logger.Send("info", fmt.Sprintf("some elements 'attribytes' successfully added to event with id:'%s' (case id:'%s') again", data.EventId, data.CaseId))

		// удаляем все прикрепленные к событию объекты
		for _, v := range oldEvent.Event.Object {
			if err := rmisp.deleteObject(ctx, v.Id); err != nil {
				m.logger.Send("error", supportingfunctions.CustomError(err).Error())
			}
		}
		// добавляем все объекты события повторно
		if err = rmisp.sendObjects(ctx, data.EventId, data.Data.GetObjects()); err != nil {
			m.logger.Send("error", supportingfunctions.CustomError(err).Error())
		}

		m.logger.Send("info", fmt.Sprintf("elements 'objects' successfully added to event with id:'%s' (case id:'%s') again", data.EventId, data.CaseId))

		// берем небольшой таймаут, нужен для того что бы MISP успел обработать и добавить в БД
		// всё ранее ему переданное, если обработка переданных объектов не была завершена
		// возможны накладки или сбои при добавлении данных
		// это недостаток MISP, с этим я ничего не могу поделать
		time.Sleep(5 * time.Second)

		// удаляем все теги события
		for _, v := range oldEvent.Event.Tag {
			if err := rmisp.deleteTag(ctx, v.Id, data.EventId); err != nil {
				m.logger.Send("error", supportingfunctions.CustomError(err).Error())
			}
		}
		// добавляем все теги повторно
		if err := rmisp.sendEventTags(ctx, data.EventId, data.Data.GetObjectTags()); err != nil {
			m.logger.Send("error", supportingfunctions.CustomError(err).Error())
		}

		m.logger.Send("info", fmt.Sprintf("elements 'tags' successfully added to event with id:'%s' (case id:'%s')", data.EventId, data.CaseId))

		time.Sleep(5 * time.Second)

		// публикуем добавленное событие
		// masterKey нужен для публикации события так как пользователь должен иметь более расшириные права
		// чем могут иметь некоторые обычные пользователи
		resMsg, err := rmisp.sendRequestPublishEvent(ctx, data.EventId)
		if err != nil {
			m.logger.Send("error", supportingfunctions.CustomError(err).Error())
		}
		if resMsg != "" {
			m.logger.Send("info", fmt.Sprintf("event with id:'%s' (case id:'%s') %s", data.EventId, data.CaseId, resMsg))
		}

		// получаем из тегов список сенсоров по которым нужна дополнитеьльная информация
		listSensors := createListSensors(*data.Data.ObjectTags)
		if len(listSensors) == 0 {
			m.logger.Send("error", supportingfunctions.CustomError(fmt.Errorf("the sensor list  on object tags '%+v' for event with id:'%s' (case id:'%s') is empty", data.Data.ObjectTags, data.EventId, data.CaseId)).Error())

			return true
		}

		reqSensorId, err := json.Marshal(RequestSensorInformation{
			Source:      "placeholder_misp",
			TaskId:      data.CaseId,
			ListSensors: listSensors,
		})
		if err != nil {
			m.logger.Send("error", supportingfunctions.CustomError(err).Error())

			return true
		}

		// запрос в ядро на получение информации о сенсорах
		m.SendDataOutput(OutputSetting{
			Command: "get sensor information",
			CaseId:  data.CaseId,
			RootId:  data.RootId,
			Data:    reqSensorId,
		})

		m.logger.Send("info", fmt.Sprintf("event with id:'%s' (case id:'%s') send request 'get sensor information' to core module", data.EventId, data.CaseId))

		return true
	})

	// добавляем вспомогательный тип specialObject в очередь хранилища
	m.cache.PushObjectToQueue(specialObject)
}

// delObject удаляет старое событие типа Event и все связанные с ним объекты
func (m *ModuleMISP) delObject(ctx context.Context, eventId string) {
	rmisp, err := NewMispRequest(
		WithHost(m.host),
		WithMasterAuthKey(m.authKey))
	if err != nil {
		m.logger.Send("error", supportingfunctions.CustomError(err).Error())

		return
	}

	if err = rmisp.deleteEvent(ctx, eventId); err != nil {
		m.logger.Send("error", supportingfunctions.CustomError(err).Error())
	}
}

// addSensorInformation добавляет информацию о сенсорах
func (m *ModuleMISP) addSensorInformation(ctx context.Context, userAuthKey string, eventId string, data []byte) {
	var sensorInfo responses.ResponseSensorsInformation
	if err := json.Unmarshal(data, &sensorInfo); err != nil {
		m.logger.Send("error", supportingfunctions.CustomError(err).Error())

		return
	}

	m.logger.Send("info", fmt.Sprintf("section:'information handling', command:'add sensor information', accepted object:'%#v'", sensorInfo))

	// если в принятом ответе от модуля обогащения информацией о сенсорах есть глобальная ошибка
	if sensorInfo.Error != "" {
		m.logger.Send("error", supportingfunctions.CustomError(errors.New(sensorInfo.Error)).Error())

		return
	}

	rmisp, err := NewMispRequest(
		WithHost(m.host),
		WithUserAuthKey(userAuthKey),
		WithMasterAuthKey(m.authKey))
	if err != nil {
		m.logger.Send("error", supportingfunctions.CustomError(err).Error())

		return
	}

	m.logger.Send("info", fmt.Sprintf("starting adding tags to event id:'%s'", eventId))

	for _, info := range sensorInfo.Informations {
		if err := rmisp.addTagToEvent(ctx, objectsmispformat.EventObjectTagsMispFormat{
			Event: eventId,
			Tag:   fmt.Sprintf("misp-galaxy:Sector=\"%s\"", info.NetboxTenantGroup),
		}); err != nil {
			m.logger.Send("error", supportingfunctions.CustomError(err).Error())

			continue
		}

		if _, err = rmisp.sendRequestPublishEvent(ctx, eventId); err != nil {
			m.logger.Send("error", supportingfunctions.CustomError(err).Error())
		}
	}
}
