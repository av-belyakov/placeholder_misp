package mispinteractions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"

	"placeholder_misp/datamodels"
)

// sendEventsMispFormat отправляет в API MISP событие в виде типа Event и возвращает полученный ответ
func sendEventsMispFormat(host, authKey string, d SettingsChanInputMISP) (*http.Response, []byte, error) {
	var (
		res         *http.Response
		resBodyByte = make([]byte, 0)
	)

	c, err := NewClientMISP(host, authKey, false)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return nil, resBodyByte, fmt.Errorf("'events add, %s' %s:%d", err.Error(), f, l-2)
	}

	ed, ok := d.MajorData["events"]
	if !ok {
		_, f, l, _ := runtime.Caller(0)

		return nil, resBodyByte, fmt.Errorf("'the properties of \"events\" were not found in the received data' %s:%d", f, l-2)
	}

	b, err := json.Marshal(ed)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)

		return nil, resBodyByte, fmt.Errorf("'events add, %s' %s:%d", err.Error(), f, l-2)
	}

	res, resBodyByte, err = c.Post("/events/add", b)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)

		return nil, resBodyByte, fmt.Errorf("'events add, %s' %s:%d", err.Error(), f, l-2)
	}

	if res.StatusCode != http.StatusOK {
		_, f, l, _ := runtime.Caller(0)

		return nil, resBodyByte, fmt.Errorf("'events add, %s' %s:%d", res.Status, f, l-1)
	}

	return res, resBodyByte, nil
}

// sendAttribytesMispFormat отправляет в API MISP список атрибутов в виде среза типов Attribytes и возвращает полученный ответ
func sendAttribytesMispFormat(host, authKey, eventId string, d SettingsChanInputMISP, logging chan<- datamodels.MessageLogging) (*http.Response, []byte) {
	var (
		res         *http.Response
		resBodyByte = make([]byte, 0)
	)

	c, err := NewClientMISP(host, authKey, false)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		logging <- datamodels.MessageLogging{
			MsgData: fmt.Sprintf("'attributes for event id%s add, %s' %s:%d", eventId, err.Error(), f, l-2),
			MsgType: "error",
		}

		return nil, resBodyByte
	}

	ad, ok := d.MajorData["attributes"]
	if !ok {
		_, f, l, _ := runtime.Caller(0)

		logging <- datamodels.MessageLogging{
			MsgData: fmt.Sprintf("'the properties of \"attributes\" were not found in the received data' %s:%d", f, l-2),
			MsgType: "error",
		}

		return nil, resBodyByte
	}

	lamf, ok := ad.([]datamodels.AttributesMispFormat)
	if !ok {
		_, f, l, _ := runtime.Caller(0)

		logging <- datamodels.MessageLogging{
			MsgData: fmt.Sprintf("'the received data does not match the type \"attributes\"' %s:%d", f, l-2),
			MsgType: "error",
		}

		return nil, resBodyByte
	}

	for k := range lamf {
		lamf[k].EventId = eventId

		if lamf[k].Value == "" {
			_, f, l, _ := runtime.Caller(0)

			logging <- datamodels.MessageLogging{
				MsgData: fmt.Sprintf("'attributes for event id №%s is not added, the \"Value\" type property should not be empty' %s:%d", eventId, f, l-1),
				MsgType: "warning",
			}

			continue
		}

		b, err := json.Marshal(lamf[k])
		if err != nil {
			_, f, l, _ := runtime.Caller(0)

			logging <- datamodels.MessageLogging{
				MsgData: fmt.Sprintf("'attributes №%s add, %s' %s:%d", eventId, err.Error(), f, l-2),
				MsgType: "warning",
			}

			continue
		}

		res, resBodyByte, err = c.Post("/attributes/add/"+eventId, b)
		if err != nil {
			_, f, l, _ := runtime.Caller(0)

			logging <- datamodels.MessageLogging{
				MsgData: fmt.Sprintf("'attributes №%s add, %s' %s:%d", eventId, err.Error(), f, l-2),
				MsgType: "warning",
			}

			continue
		}

		if res.StatusCode != http.StatusOK {
			_, f, l, _ := runtime.Caller(0)

			logging <- datamodels.MessageLogging{
				MsgData: fmt.Sprintf("'attributes №%s add, %s' %s:%d", eventId, res.Status, f, l-1),
				MsgType: "warning",
			}
		}
	}

	return res, resBodyByte
}

// sendObjectsMispFormat отправляет в API MISP список объектов содержащихся в свойстве observables.attachment (как правило это описание вложеного файла)
func sendObjectsMispFormat(host, authKey, eventId string, d SettingsChanInputMISP, logging chan<- datamodels.MessageLogging) (*http.Response, []byte) {
	var (
		res         *http.Response
		resBodyByte = make([]byte, 0)
	)

	c, err := NewClientMISP(host, authKey, false)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		logging <- datamodels.MessageLogging{
			MsgData: fmt.Sprintf("'objects for event id:%s add, %s' %s:%d", eventId, err.Error(), f, l-2),
			MsgType: "error",
		}

		return nil, resBodyByte
	}

	od, ok := d.MajorData["objects"]
	if !ok {
		_, f, l, _ := runtime.Caller(0)

		logging <- datamodels.MessageLogging{
			MsgData: fmt.Sprintf("'the properties of \"objects\" were not found in the received data' %s:%d", f, l-2),
			MsgType: "error",
		}

		return nil, resBodyByte
	}

	lomf, ok := od.(map[int]datamodels.ObjectsMispFormat)
	if !ok {
		_, f, l, _ := runtime.Caller(0)

		logging <- datamodels.MessageLogging{
			MsgData: fmt.Sprintf("'the received data does not match the type \"objects\"' %s:%d", f, l-2),
			MsgType: "error",
		}

		return nil, resBodyByte
	}

	for _, v := range lomf {
		v.EventId = eventId

		b, err := json.Marshal(v)
		if err != nil {
			_, f, l, _ := runtime.Caller(0)

			logging <- datamodels.MessageLogging{
				MsgData: fmt.Sprintf("'objects №%s add, %s' %s:%d", eventId, err.Error(), f, l-2),
				MsgType: "warning",
			}

			continue
		}

		res, resBodyByte, err = c.Post("/objects/add/"+eventId, b)
		if err != nil {
			_, f, l, _ := runtime.Caller(0)

			logging <- datamodels.MessageLogging{
				MsgData: fmt.Sprintf("'objects №%s add, %s' %s:%d", eventId, err.Error(), f, l-2),
				MsgType: "warning",
			}

			continue
		}

		if res.StatusCode != http.StatusOK {
			_, f, l, _ := runtime.Caller(0)

			logging <- datamodels.MessageLogging{
				MsgData: fmt.Sprintf("'objects №%s add, %s' %s:%d", eventId, res.Status, f, l-1),
				MsgType: "warning",
			}
		}
	}

	return res, resBodyByte
}

func sendRequestPublishEvent(host, authKey, eventId string) (string, error) {
	var resultMsg string

	c, err := NewClientMISP(host, authKey, false)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)

		return resultMsg, fmt.Errorf("'event publish add, %s' %s:%d", err.Error(), f, l-2)
	}

	res, b, err := c.Post("/events/publish/"+eventId, []byte{})
	if err != nil {
		_, f, l, _ := runtime.Caller(0)

		return resultMsg, fmt.Errorf("'event publish add, %s' %s:%d", err.Error(), f, l-2)
	}

	resData := decodeResponseMIspMessage(b)
	resultMsg = fmt.Sprintf("result published event with id '%s' - %s '%s' %s", eventId, resData.name, resData.message, resData.success)

	if res.StatusCode != http.StatusOK {
		_, f, l, _ := runtime.Caller(0)

		return resultMsg, fmt.Errorf("'event publish add, %s' %s:%d", res.Status, f, l-1)
	}

	return resultMsg, nil
}

func sendEventReportsMispFormat(host, authKey, eventId string, caseId float64) error {
	c, err := NewClientMISP(host, authKey, false)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return fmt.Errorf("'event report add, %s' %s:%d", err.Error(), f, l-2)
	}

	b, err := json.Marshal(datamodels.EventReports{
		Name:         fmt.Sprint(caseId),
		Distribution: "1",
	})
	if err != nil {
		_, f, l, _ := runtime.Caller(0)

		return fmt.Errorf("'event report add, %s' %s:%d", err.Error(), f, l-2)
	}

	res, _, err := c.Post("/event_reports/add/"+eventId, b)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)

		return fmt.Errorf("'event report add, %s' %s:%d", err.Error(), f, l-2)
	}

	if res.StatusCode != http.StatusOK {
		_, f, l, _ := runtime.Caller(0)

		return fmt.Errorf("'event report add, %s' %s:%d", res.Status, f, l-1)
	}

	return nil
}

func sendEventTagsMispFormat(host, authKey, eventId string, d SettingsChanInputMISP, logging chan<- datamodels.MessageLogging) error {
	c, err := NewClientMISP(host, authKey, false)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return fmt.Errorf("'event tags add, %s' %s:%d", err.Error(), f, l-2)
	}

	eot, ok := d.MajorData["event.object.tags"]
	if !ok {
		_, f, l, _ := runtime.Caller(0)

		logging <- datamodels.MessageLogging{
			MsgData: fmt.Sprintf("'the properties of \"objects\" were not found in the received data' %s:%d", f, l-2),
			MsgType: "error",
		}

		return nil
	}

	leot, ok := eot.(datamodels.ListEventObjectTags)
	if !ok {
		_, f, l, _ := runtime.Caller(0)

		logging <- datamodels.MessageLogging{
			MsgData: fmt.Sprintf("'the received data does not match the type \"objects\"' %s:%d", f, l-2),
			MsgType: "error",
		}

		return nil
	}

	eotmf := datamodels.EventObjectTagsMispFormat{}

	// ***********************************
	// Это логирование только для теста!!!
	// ***********************************
	logging <- datamodels.MessageLogging{
		MsgData: fmt.Sprintf("TEST_INFO func 'sendEventTagsMispFormat', готовимся добавлять event tags - %v", leot.GetListTags()),
		MsgType: "testing",
	}
	//
	//

	for _, v := range leot {
		eotmf.Event = eventId
		eotmf.Tag = v

		b, err := json.Marshal(eotmf)
		if err != nil {
			_, f, l, _ := runtime.Caller(0)

			logging <- datamodels.MessageLogging{
				MsgData: fmt.Sprintf("'event tags №%s add, %s' %s:%d", eventId, err.Error(), f, l-2),
				MsgType: "warning",
			}

			continue
		}

		// ***********************************
		// Это логирование только для теста!!!
		// ***********************************
		logging <- datamodels.MessageLogging{
			MsgData: fmt.Sprintf("TEST_INFO func 'sendEventTagsMispFormat', готовимся отправить POST запрос для добавления тега %s", v),
			MsgType: "testing",
		}
		//
		//

		res, b, err := c.Post("/events/addTag", b)
		if err != nil {
			_, f, l, _ := runtime.Caller(0)

			logging <- datamodels.MessageLogging{
				MsgData: fmt.Sprintf("'event tags №%s add, %s' %s:%d", eventId, err.Error(), f, l-2),
				MsgType: "warning",
			}

			continue
		}

		resData := decodeResponseMIspMessage(b)
		resultMsg := fmt.Sprintf("tag: '%s' %s '%s' %s errors:'%s'", v, resData.name, resData.message, resData.success, resData.errors)
		// ***********************************
		// Это логирование только для теста!!!
		// ***********************************
		logging <- datamodels.MessageLogging{
			MsgData: fmt.Sprintf("TEST_INFO func 'sendEventTagsMispFormat', результат выполнения POST запроса - '%s'", resultMsg),
			MsgType: "testing",
		}
		//
		//

		if res.StatusCode != http.StatusOK {
			_, f, l, _ := runtime.Caller(0)

			logging <- datamodels.MessageLogging{
				MsgData: fmt.Sprintf("'event tags №%s add, %s' %s:%d", eventId, res.Status, f, l-1),
				MsgType: "warning",
			}
		}
	}

	return nil
}

// удаляем дублирующиеся события из MISP
func delEventsMispFormat(host, authKey, eventId string) (*http.Response, error) {
	c, err := NewClientMISP(host, authKey, false)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return nil, fmt.Errorf("'events delete, %s' %s:%d", err.Error(), f, l-2)
	}

	res, _, err := c.Delete("/events/delete/" + eventId)
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
