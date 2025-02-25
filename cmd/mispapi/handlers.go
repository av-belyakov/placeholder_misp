package mispapi

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/av-belyakov/objectsmispformat"
	"github.com/av-belyakov/placeholder_misp/commoninterfaces"
	"github.com/av-belyakov/placeholder_misp/internal/supportingfunctions"
)

func (m *ModuleMISP) addNewObject(ctx context.Context, userAuthKey string, data InputSettings) {
	specialObject := NewCacheSpecialObject[*objectsmispformat.ListFormatsMISP]()
	specialObject.SetID(data.RootId)
	specialObject.SetObject(&data.Data)
	specialObject.SetFunc(func(i int) bool {
		rmisp, err := NewMispRequest(
			WithHost(m.host),
			WithUserAuthKey(userAuthKey),
			WithMasterAuthKey(m.authKey))
		if err != nil {
			return false
		}

		// ***********************************
		// Это логирование только для теста!!!
		// ***********************************
		m.logger.Send("testing", fmt.Sprintf("TEST_INFO func 'HandlerMISP', send EVENTS to ----> MISP	USER EMAIL: %s, CaseId: %v", data.UserEmail, data.CaseId))
		//
		//

		//отправляет в API MISP событие в виде типа Event и возвращает результат который содержит
		//id события в MISP, у MISP свой уникальный id для событий
		//только с использованием этого id в MISP добавляются все остальные объекты
		_, resBodyByte, err := rmisp.sendEvent(ctx, data.Data.GetEvent())
		if err != nil {
			m.logger.Send("error", supportingfunctions.CustomError(err).Error())

			return false
		}

		resMisp := MispResponse{}
		if err := json.Unmarshal(resBodyByte, &resMisp); err != nil {
			m.logger.Send("error", supportingfunctions.CustomError(err).Error())

			return false
		}

		//получаем уникальный id MISP
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

			return false
		}

		// ***********************************
		// Это логирование только для теста!!!
		// ***********************************
		m.logger.Send("testing", fmt.Sprintf("TEST_INFO func 'HandlerMISP', send EVENTS REPORTS to ----> MISP	USER EMAIL: %s, CaseId: %v", data.UserEmail, data.CaseId))
		//
		//

		// добавляем event_reports
		if err := rmisp.sendEventReports(ctx, eventId, data.Data.GetReports()); err != nil {
			m.logger.Send("error", supportingfunctions.CustomError(err).Error())

			return false
		}

		// ***********************************
		// Это логирование только для теста!!!
		// ***********************************
		m.logger.Send("testing", fmt.Sprintf("TEST_INFO func 'HandlerMISP', send data to ----> RedisDB USER EMAIL: %s, CaseId: %v", data.UserEmail, data.CaseId))
		//
		//

		//отправляем запрос для добавления в БД Redis, id кейса и нового события
		m.SendingDataOutput(OutputSetting{
			Command: "set new event id",
			CaseId:  fmt.Sprint(data.CaseId),
			EventId: eventId,
		})

		// ***********************************
		// Это логирование только для теста!!!
		// ***********************************
		m.logger.Send("testing", fmt.Sprintf("TEST_INFO func 'HandlerMISP', send ATTRIBYTES to ----> MISP	USER EMAIL: %s, CaseId: %v", data.UserEmail, data.CaseId))
		//
		//

		//добавляем атрибуты
		_, _, warning, err := rmisp.sendAttribytes(ctx, eventId, data.Data.GetAttributes())
		if err != nil {
			// тут ошибка может быть при добавлении только одного из многих объектов
			// соответственно тормозить весь процесс только из-за того что была ошибка
			// при добавлении одного или нескольких объектов не стоит
			// если же не был добавлен ни один из объектов, то это возможно глабальная
			// ошибка доступа, следовательно, при добавлении следующих объектов она также
			// может вылезти, тогда там и будет выполнен останов всей цепочки
			m.logger.Send("error", supportingfunctions.CustomError(err).Error())
		}
		if warning != "" {
			m.logger.Send("warning", warning)
		}

		// ***********************************
		// Это логирование только для теста!!!
		// ***********************************
		m.logger.Send("testing", fmt.Sprintf("TEST_INFO func 'HandlerMISP', send OBJECTS to ----> MISP	USER EMAIL: %s, CaseId: %v", data.UserEmail, data.CaseId))
		//
		//

		// добавляем объекты
		if _, _, err = rmisp.sendObjects(ctx, eventId, data.Data.GetObjects()); err != nil {
			//тут такая же ситуация что и с ошибками при выполнении метода rmisp.sendAttribytes
			m.logger.Send("error", supportingfunctions.CustomError(err).Error())
		}

		// берем небольшой таймаут, нужен для того что бы MISP успел обработать и добавить в БД
		// всё ранее ему переданное, если обработка переданных объектов не была завершена
		// возможны накладки или сбои при добавлении данных
		// это недостаток MISP, с этим я ничего не могу поделать
		time.Sleep(4 * time.Second)

		// ***********************************
		// Это логирование только для теста!!!
		// ***********************************
		m.logger.Send("testing", fmt.Sprintf("TEST_INFO func 'HandlerMISP', send EVENT_TAGS to ----> MISP	USER EMAIL: %s, CaseId: %v", data.UserEmail, data.CaseId))
		//
		//

		// добавляем event_tags
		if err := rmisp.sendEventTags(ctx, eventId, data.Data.GetObjectTags()); err != nil {
			m.logger.Send("error", supportingfunctions.CustomError(err).Error())
		}

		// ***********************************
		// Это логирование только для теста!!!
		// ***********************************
		m.logger.Send("testing", fmt.Sprintf("TEST_INFO func 'HandlerMISP', send PUBLISH to ----> MISP	USER EMAIL: %s, CaseId: %v", data.UserEmail, data.CaseId))
		//
		//

		//публикуем добавленное событие
		//masterKey нужен для публикации события так как пользователь
		//должен иметь более расшириные права чем могут иметь некоторые
		//обычные пользователи
		resMsg, err := rmisp.sendRequestPublishEvent(ctx, eventId)
		if err != nil {
			m.logger.Send("error", supportingfunctions.CustomError(err).Error())
		}
		if resMsg != "" {
			m.logger.Send("warning", resMsg)
		}

		// ***********************************
		// Это логирование только для теста!!!
		// ***********************************
		m.logger.Send("testing", fmt.Sprintf("TEST_INFO func 'HandlerMISP', send EventId to ----> CORE	USER EMAIL: %s, CaseId: %v", data.UserEmail, data.CaseId))
		//
		//

		// отправляем в ядро информацию по event Id
		m.SendingDataOutput(OutputSetting{
			Command: "send event id",
			EventId: eventId,
			CaseId:  fmt.Sprint(data.CaseId),
			RootId:  data.RootId,
			TaskId:  data.TaskId,
		})

		//по умолчанию 'не успешно'
		return false
	})
}

func delEventById(ctx context.Context, host, authKey, eventId string, logger commoninterfaces.Logger) {
	// ***********************************
	// Это логирование только для теста!!!
	// ***********************************
	logger.Send("testing", fmt.Sprintf("TEST_INFO func 'HandlerMISP', удаление события типа event, где event id: %s", eventId))
	//
	//

	// удаление события типа event
	_, err := delEventsMispFormat(ctx, host, authKey, eventId)
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
