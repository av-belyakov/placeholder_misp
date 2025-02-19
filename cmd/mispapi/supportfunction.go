package mispapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/av-belyakov/placeholder_misp/internal/supportingfunctions"
)

/*
// sendEventsMispFormat отправляет в API MISP событие в виде типа Event и возвращает полученный ответ
func sendEventsMispFormat(host, authKey string, d InputSettings) (*http.Response, []byte, error) {
	var (
		res         *http.Response
		resBodyByte = make([]byte, 0)
	)

	c, err := NewClientMISP(host, authKey, false)
	if err != nil {
		return nil, resBodyByte, supportingfunctions.CustomError(fmt.Errorf("events add, %w", err))
	}

	ed, ok := d.MajorData["events"]
	if !ok {
		return nil, resBodyByte, supportingfunctions.CustomError(errors.New("the properties of 'events' were not found in the received data"))
	}

	b, err := json.Marshal(ed)
	if err != nil {
		return nil, resBodyByte, supportingfunctions.CustomError(fmt.Errorf("events add, %w", err))
	}

	res, resBodyByte, err = c.Post("/events/add", b)
	if err != nil {
		return nil, resBodyByte, supportingfunctions.CustomError(fmt.Errorf("events add, %w", err))
	}

	if res.StatusCode != http.StatusOK {
		return nil, resBodyByte, supportingfunctions.CustomError(fmt.Errorf("events add, status '%s'", res.Status))
	}

	return res, resBodyByte, nil
}

// sendAttribytesMispFormat отправляет в API MISP список атрибутов в виде среза типов Attribytes и возвращает полученный ответ
func sendAttribytesMispFormat(host, authKey, eventId string, d InputSettings, logger commoninterfaces.Logger) (*http.Response, []byte) {
	var (
		res         *http.Response
		resBodyByte = make([]byte, 0)
	)

	c, err := NewClientMISP(host, authKey, false)
	if err != nil {
		logger.Send("error", supportingfunctions.CustomError(fmt.Errorf("'attributes' for event id:'%s' add, %w", eventId, err)).Error())

		return nil, resBodyByte
	}

	ad, ok := d.MajorData["attributes"]
	if !ok {
		logger.Send("error", supportingfunctions.CustomError(errors.New("the properties of 'attributes' were not found in the received data")).Error())

		return nil, resBodyByte
	}

	lamf, ok := ad.([]datamodels.AttributesMispFormat)
	if !ok {
		logger.Send("error", supportingfunctions.CustomError(errors.New("the received data does not match the type 'attributes'")).Error())

		return nil, resBodyByte
	}

	for k := range lamf {
		lamf[k].EventId = eventId

		if lamf[k].Value == "" {
			logger.Send("warning", fmt.Sprintf("'attributes' for event id:'%s' is not added, the 'Value' type property should not be empty", eventId))

			continue
		}

		b, err := json.Marshal(lamf[k])
		if err != nil {
			logger.Send("error", supportingfunctions.CustomError(fmt.Errorf("'attributes' id:'%s' add, %w", eventId, err)).Error())

			continue
		}

		res, resBodyByte, err = c.Post("/attributes/add/"+eventId, b)
		if err != nil {
			attrObject, errMarshal := json.MarshalIndent(lamf[k], "", "  ")
			if errMarshal != nil {
				logger.Send("error", supportingfunctions.CustomError(errors.New("the received data does not match the type 'attributes'")).Error())
			}

			logger.Send("warning", fmt.Sprintf("'attributes' with id:'%s' add, object:\n%s\n%s", eventId, string(attrObject), err.Error()))

			continue
		}

		if res.StatusCode != http.StatusOK {
			logger.Send("error", supportingfunctions.CustomError(fmt.Errorf("'attributes' with id:'%s' add, status '%s'", eventId, res.Status)).Error())
		}
	}

	return res, resBodyByte
}

// sendObjectsMispFormat отправляет в API MISP список объектов содержащихся в свойстве observables.attachment (как правило это описание вложеного файла)
func sendObjectsMispFormat(host, authKey, eventId string, d InputSettings, logger commoninterfaces.Logger) (*http.Response, []byte) {
	var (
		res         *http.Response
		resBodyByte = make([]byte, 0)
	)

	c, err := NewClientMISP(host, authKey, false)
	if err != nil {
		logger.Send("error", supportingfunctions.CustomError(fmt.Errorf("objects for event id:'%s' add, %w", eventId, err)).Error())

		return nil, resBodyByte
	}

	od, ok := d.MajorData["objects"]
	if !ok {
		logger.Send("error", supportingfunctions.CustomError(errors.New("the properties of 'objects' were not found in the received data")).Error())

		return nil, resBodyByte
	}

	lomf, ok := od.(map[int]datamodels.ObjectsMispFormat)
	if !ok {
		logger.Send("error", supportingfunctions.CustomError(errors.New("the received data does not match the type 'objects'")).Error())

		return nil, resBodyByte
	}

	for _, v := range lomf {
		v.EventId = eventId

		b, err := json.Marshal(v)
		if err != nil {
			logger.Send("error", supportingfunctions.CustomError(fmt.Errorf("objects with id:'%s' add, %w", eventId, err)).Error())

			continue
		}

		res, resBodyByte, err = c.Post("/objects/add/"+eventId, b)
		if err != nil {
			logger.Send("error", supportingfunctions.CustomError(fmt.Errorf("objects with id:'%s' add, %w", eventId, err)).Error())

			continue
		}

		if res.StatusCode != http.StatusOK {
			logger.Send("error", supportingfunctions.CustomError(fmt.Errorf("objects with id:'%s' add, status '%s'", eventId, res.Status)).Error())
		}
	}

	return res, resBodyByte
}

func sendRequestPublishEvent(host, authKey, eventId string) (string, error) {
	var resultMsg string

	c, err := NewClientMISP(host, authKey, false)
	if err != nil {
		return resultMsg, supportingfunctions.CustomError(fmt.Errorf("event publish add, %w", err))
	}

	res, b, err := c.Post("/events/publish/"+eventId, []byte{})
	if err != nil {
		return resultMsg, supportingfunctions.CustomError(fmt.Errorf("event publish add, %w", err))
	}

	resData := decodeResponseMIspMessage(b)
	resultMsg = fmt.Sprintf("result published event with id '%s' - %s '%s' %s", eventId, resData.name, resData.message, resData.success)

	if res.StatusCode != http.StatusOK {
		return resultMsg, supportingfunctions.CustomError(fmt.Errorf("event publish add, status '%s'", res.Status))
	}

	return resultMsg, nil
}

func sendEventReportsMispFormat(host, authKey, eventId string, caseId float64) error {
	c, err := NewClientMISP(host, authKey, false)
	if err != nil {
		return supportingfunctions.CustomError(fmt.Errorf("event report add, %w", err))
	}

	b, err := json.Marshal(datamodels.EventReports{
		Name:         fmt.Sprint(caseId),
		Distribution: "1",
	})
	if err != nil {
		return supportingfunctions.CustomError(fmt.Errorf("event report add, %w", err))
	}

	res, _, err := c.Post("/event_reports/add/"+eventId, b)
	if err != nil {
		return supportingfunctions.CustomError(fmt.Errorf("event report add, %w", err))
	}

	if res.StatusCode != http.StatusOK {
		return supportingfunctions.CustomError(fmt.Errorf("event report add, status '%s'", res.Status))
	}

	return nil
}

func sendEventTagsMispFormat(host, authKey, eventId string, d InputSettings, logger commoninterfaces.Logger) error {
	c, err := NewClientMISP(host, authKey, false)
	if err != nil {
		return supportingfunctions.CustomError(fmt.Errorf("event tags add, %w", err))
	}

	eot, ok := d.MajorData["event.object.tags"]
	if !ok {
		logger.Send("error", supportingfunctions.CustomError(errors.New("the properties of 'objects' were not found in the received data")).Error())

		return nil
	}

	leot, ok := eot.(datamodels.ListEventObjectTags)
	if !ok {
		logger.Send("error", supportingfunctions.CustomError(errors.New("the received data does not match the type 'objects'")).Error())

		return nil
	}

	eotmf := datamodels.EventObjectTagsMispFormat{}

	// ***********************************
	// Это логирование только для теста!!!
	// ***********************************
	logger.Send("testing", fmt.Sprintf("TEST_INFO func 'sendEventTagsMispFormat', готовимся добавлять event tags - %v", leot.GetListTags()))
	//
	//

	for _, v := range leot {
		eotmf.Event = eventId
		eotmf.Tag = v

		b, err := json.Marshal(eotmf)
		if err != nil {
			logger.Send("error", supportingfunctions.CustomError(fmt.Errorf("'event tags with id:'%s' add, %w", eventId, err)).Error())

			continue
		}

		// ***********************************
		// Это логирование только для теста!!!
		// ***********************************
		logger.Send("testing", fmt.Sprintf("TEST_INFO func 'sendEventTagsMispFormat', готовимся отправить POST запрос для добавления тега %s", v))
		//
		//

		res, b, err := c.Post("/events/addTag", b)
		if err != nil {
			logger.Send("error", supportingfunctions.CustomError(fmt.Errorf("'event tags with id:'%s' add, %w", eventId, err)).Error())

			continue
		}

		resData := decodeResponseMIspMessage(b)
		resultMsg := fmt.Sprintf("tag: '%s' %s '%s' %s errors:'%s'", v, resData.name, resData.message, resData.success, resData.errors)
		logger.Send("warning", fmt.Sprintf("event tags with id:'%s' the result of executing the POST query - '%s'", eventId, resultMsg))

		if res.StatusCode != http.StatusOK {
			logger.Send("error", supportingfunctions.CustomError(fmt.Errorf("'event tags with id:'%s' add, status '%s'", eventId, res.Status)).Error())
		}
	}

	return nil
}*/

// удаляем дублирующиеся события из MISP
func delEventsMispFormat(ctx context.Context, host, authKey, eventId string) (*http.Response, error) {
	ctxTimeout, CancelFunc := context.WithTimeout(ctx, time.Second*15)
	defer CancelFunc()

	c, err := NewClientMISP(host, authKey, false)
	if err != nil {
		return nil, supportingfunctions.CustomError(fmt.Errorf("events delete, %w", err))
	}

	res, _, err := c.Delete(ctxTimeout, "/events/delete/"+eventId)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func decodeResponseMIspMessage(b []byte) struct {
	name    string
	errors  string
	message string
	success string
} {
	msg := struct {
		name    string
		errors  string
		message string
		success string
	}{}
	resTmp := map[string]interface{}{}
	if err := json.Unmarshal(b, &resTmp); err == nil {
		for k, v := range resTmp {
			switch k {
			case "name":
				msg.name = fmt.Sprint(v)

			case "errors":
				msg.errors = fmt.Sprint(v)

			case "message":
				msg.message = fmt.Sprint(v)

			case "success":
				msg.success = fmt.Sprint(v)

			}
		}
	}

	return msg
}
