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
			m.logger.Send("error", supportingfunctions.CustomError(err).Error())

			return false
		}

		m.logger.Send("info", fmt.Sprintf("starting adding the case id:'%d'", int(data.CaseId)))

		//отправляет в API MISP событие в виде типа Event и возвращает результат который содержит
		//id события в MISP, у MISP свой уникальный id для событий
		//только с использованием этого id в MISP добавляются все остальные объекты
		_, resBodyByte, err := rmisp.sendEvent(ctx, data.Data.GetEvent())
		if err != nil {
			m.logger.Send("error", supportingfunctions.CustomError(err).Error())

			return false
		}

		//Все ошибки которые могут возникнуть при дальнейшем взаимодействии с MISP
		//будут попрежнему логироватся.
		//Однако, статус выполнения для функции будет ставится в TRUE, что бы не досить
		//MISP, так как все последующие попытки будут начинатся с добавления 'event', а
		//добавить 'event' с таким id нельзя.
		//Необходимо удалить предыдущий.

		resMisp := MispResponse{}
		if err := json.Unmarshal(resBodyByte, &resMisp); err != nil {
			m.logger.Send("error", supportingfunctions.CustomError(err).Error())

			return true
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

			return true
		}

		m.logger.Send("info", fmt.Sprintf("new element 'event' with id:'%s' successfully created (case id:'%d')", eventId, int(data.CaseId)))

		// добавляем event_reports
		if err := rmisp.sendEventReports(ctx, eventId, data.Data.GetReports()); err != nil {
			m.logger.Send("error", supportingfunctions.CustomError(err).Error())

			return true
		}

		m.logger.Send("info", fmt.Sprintf("element 'event_reports' successfully added to event with id:'%s' (case id:'%d')", eventId, int(data.CaseId)))

		//отправляем запрос для добавления в БД Redis, id кейса и нового события
		m.SendDataOutput(OutputSetting{
			Command: "set new event id",
			CaseId:  fmt.Sprint(data.CaseId),
			EventId: eventId,
		})

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

		m.logger.Send("info", fmt.Sprintf("some elements 'attribytes' successfully added to event with id:'%s' (case id:'%d')", eventId, int(data.CaseId)))

		// добавляем объекты
		if _, _, err = rmisp.sendObjects(ctx, eventId, data.Data.GetObjects()); err != nil {
			//тут такая же ситуация что и с ошибками при выполнении метода rmisp.sendAttribytes
			m.logger.Send("error", supportingfunctions.CustomError(err).Error())
		}

		m.logger.Send("info", fmt.Sprintf("elements 'objects' successfully added to event with id:'%s' (case id:'%d')", eventId, int(data.CaseId)))

		// берем небольшой таймаут, нужен для того что бы MISP успел обработать и добавить в БД
		// всё ранее ему переданное, если обработка переданных объектов не была завершена
		// возможны накладки или сбои при добавлении данных
		// это недостаток MISP, с этим я ничего не могу поделать
		time.Sleep(4 * time.Second)

		// добавляем event_tags
		if err := rmisp.sendEventTags(ctx, eventId, data.Data.GetObjectTags()); err != nil {
			m.logger.Send("error", supportingfunctions.CustomError(err).Error())
		}

		m.logger.Send("info", fmt.Sprintf("elements 'tags' successfully added to event with id:'%s' (case id:'%d')", eventId, int(data.CaseId)))

		time.Sleep(1 * time.Second)

		//публикуем добавленное событие
		//masterKey нужен для публикации события так как пользователь
		//должен иметь более расшириные права чем могут иметь некоторые
		//обычные пользователи
		resMsg, err := rmisp.sendRequestPublishEvent(ctx, eventId)
		if err != nil {
			m.logger.Send("error", supportingfunctions.CustomError(err).Error())
		}
		if resMsg != "" {
			m.logger.Send("info", fmt.Sprintf("event with id:'%s' (case id:'%d') %s", eventId, int(data.CaseId), resMsg))
		}

		//
		// !!!!!!!!!!!!!!!!!!!!!!
		//
		// надо разобратся что за ошибка:
		//2025/03/06 17:49:47 methods.go:67: error processing command 'send event id' for
		// caseId '39100' (rootId '~88678416456') 'nats: no responders available for
		// request' /home/artemij/go/src/placeholder_misp/cmd/natsapi/handlers.go:18
		// /home/artemij/go/src/placeholder_misp/cmd/natsapi/app.go:114
		//
		//
		// !!!!!!!!!!!!!!!!!!!!!!
		//

		// отправляем в ядро информацию по event Id
		m.SendDataOutput(OutputSetting{
			Command: "send event id",
			EventId: eventId,
			CaseId:  fmt.Sprint(data.CaseId),
			RootId:  data.RootId,
			TaskId:  data.TaskId,
		})

		//выполнено 'успешно'
		return true
	})

	//добавляем вспомогательный тип specialObject в очередь хранилища
	m.cache.PushObjectToQueue(specialObject)
}

// delEventById удаляет событие типа event
func delEventById(ctx context.Context, host, authKey, eventId string, logger commoninterfaces.Logger) {
	_, err := delEventsMispFormat(ctx, host, authKey, eventId)
	if err != nil {
		logger.Send("error", supportingfunctions.CustomError(err).Error())
	}
}
