package mispapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/av-belyakov/objectsmispformat"
	"github.com/av-belyakov/placeholder_misp/internal/supportingfunctions"
)

// NewMispRequest конструктор запроса к MISP
func NewMispRequest(opts ...RequestMISPOptions) (*requestMISP, error) {
	mispReq := &requestMISP{}

	for _, opt := range opts {
		if err := opt(mispReq); err != nil {
			return mispReq, err
		}
	}

	return mispReq, nil
}

// WithHost имя или ip адрес хоста API
func WithHost(v string) RequestMISPOptions {
	return func(n *requestMISP) error {
		if v == "" {
			return errors.New("the value of 'host' cannot be empty")
		}

		n.host = v

		return nil
	}
}

// WithUserAuthKey пользовательский ключ авторизации
func WithUserAuthKey(v string) RequestMISPOptions {
	return func(n *requestMISP) error {
		if v == "" {
			return errors.New("the value of 'userAuthKey' cannot be empty")
		}

		n.userAuthKey = v

		return nil
	}
}

// WithMasterAuthKey привилегированный ключ авторизации
func WithMasterAuthKey(v string) RequestMISPOptions {
	return func(n *requestMISP) error {
		if v == "" {
			return errors.New("the value of 'masterAuthKey' cannot be empty")
		}

		n.masterAuthKey = v

		return nil
	}
}

func (rmisp *requestMISP) SendEvent_ForTest(ctx context.Context, data *objectsmispformat.EventsMispFormat) (*http.Response, []byte, error) {
	return rmisp.sendEvent(ctx, data)
}

// sendEvent добавляет в MISP объект типа 'event'
func (rmisp *requestMISP) sendEvent(ctx context.Context, data *objectsmispformat.EventsMispFormat) (*http.Response, []byte, error) {
	var (
		res         *http.Response
		resBodyByte = make([]byte, 0)
	)

	c, err := NewClientMISP(rmisp.host, rmisp.userAuthKey, false)
	if err != nil {
		return nil, resBodyByte, supportingfunctions.CustomError(fmt.Errorf("events add, %w", err))
	}

	b, err := json.Marshal(data)
	if err != nil {
		return nil, resBodyByte, supportingfunctions.CustomError(fmt.Errorf("events add, %w", err))
	}

	res, resBodyByte, err = c.Post(ctx, "/events/add", b)
	if err != nil {
		return nil, resBodyByte, supportingfunctions.CustomError(fmt.Errorf("events add, %w", err))
	}

	return res, resBodyByte, nil
}

func (rmisp *requestMISP) SendEventReports_ForTest(ctx context.Context, eventId string, data *objectsmispformat.EventReports) error {
	return rmisp.sendEventReports(ctx, eventId, data)
}

// sendEventReports добавляет в MISP объект типа 'event_reports'
func (rmisp *requestMISP) sendEventReports(ctx context.Context, eventId string, data *objectsmispformat.EventReports) error {
	c, err := NewClientMISP(rmisp.host, rmisp.userAuthKey, false)
	if err != nil {
		return supportingfunctions.CustomError(fmt.Errorf("event report add, %w", err))
	}

	b, err := json.Marshal(data)
	if err != nil {
		return supportingfunctions.CustomError(fmt.Errorf("event report add, %w", err))
	}

	_, _, err = c.Post(ctx, "/event_reports/add/"+eventId, b)
	if err != nil {
		return supportingfunctions.CustomError(fmt.Errorf("event report add, %w", err))
	}

	return nil
}

func (rmisp *requestMISP) SendAttribytes_ForTest(ctx context.Context, eventId string, data []*objectsmispformat.AttributesMispFormat) (*http.Response, []byte, string, error) {
	return rmisp.sendAttribytes(ctx, eventId, data)
}

// sendAttribytes отправляет в MISP список атрибутов в виде среза объектов типа 'attribytes'
func (rmisp *requestMISP) sendAttribytes(ctx context.Context, eventId string, data []*objectsmispformat.AttributesMispFormat) (*http.Response, []byte, string, error) {
	var (
		res         *http.Response
		resBodyByte = make([]byte, 0)
	)

	warning := strings.Builder{}
	defer warning.Reset()

	c, err := NewClientMISP(rmisp.host, rmisp.userAuthKey, false)
	if err != nil {
		return nil, resBodyByte, warning.String(), supportingfunctions.CustomError(fmt.Errorf("'attributes' for event id:'%s' add, %w", eventId, err))
	}

	for k := range data {
		data[k].EventId = eventId

		if data[k].Value == "" {
			warning.WriteString(fmt.Sprintf("'attributes' for event id:'%s' is not added, the 'Value' type property should not be empty", eventId))

			continue
		}

		b, errTmp := json.Marshal(data[k])
		if errTmp != nil {
			err = errors.Join(err, supportingfunctions.CustomError(fmt.Errorf("'attributes' id:'%s' add, %w", eventId, errTmp)))

			continue
		}

		res, resBodyByte, errTmp = c.Post(ctx, "/attributes/add/"+eventId, b)
		if errTmp != nil {
			err = errors.Join(err, supportingfunctions.CustomError(fmt.Errorf("'attributes' id:'%s' add, %w", eventId, errTmp)))

			attrObject, errTmp := json.MarshalIndent(data[k], "", "  ")
			if errTmp != nil {
				err = errors.Join(err, supportingfunctions.CustomError(fmt.Errorf("'attributes' id:'%s' add, %w", eventId, errTmp)))
			}

			warning.WriteString(fmt.Sprintf("'attributes' with id:'%s' add, object:%s\n", eventId, string(attrObject)))

			continue
		}
	}

	return res, resBodyByte, warning.String(), err
}

func (rmisp *requestMISP) SendObjects_ForTest(ctx context.Context, eventId string, data map[int]*objectsmispformat.ObjectsMispFormat) (*http.Response, []byte, error) {
	return rmisp.sendObjects(ctx, eventId, data)
}

// sendObjects отправляет в MISP список объектов содержащихся в свойстве observables.attachment
// (как правило это описание вложеного файла)
func (rmisp *requestMISP) sendObjects(ctx context.Context, eventId string, data map[int]*objectsmispformat.ObjectsMispFormat) (*http.Response, []byte, error) {
	var (
		res         *http.Response
		resBodyByte = make([]byte, 0)
	)

	c, err := NewClientMISP(rmisp.host, rmisp.userAuthKey, false)
	if err != nil {
		return nil, resBodyByte, supportingfunctions.CustomError(fmt.Errorf("objects for event id:'%s' add, %w", eventId, err))
	}

	for _, v := range data {
		v.EventId = eventId

		b, errTmp := json.Marshal(v)
		if errTmp != nil {
			err = errors.Join(err, supportingfunctions.CustomError(fmt.Errorf("objects with id:'%s' add, %w", eventId, errTmp)))

			continue
		}

		res, resBodyByte, errTmp = c.Post(ctx, "/objects/add/"+eventId, b)
		if errTmp != nil {
			err = errors.Join(err, supportingfunctions.CustomError(fmt.Errorf("objects with id:'%s' add, %w", eventId, errTmp)))

			continue
		}
	}

	return res, resBodyByte, err
}

func (rmisp *requestMISP) SendEventTags_ForTest(ctx context.Context, eventId string, data *objectsmispformat.ListEventObjectTags) error {
	return rmisp.sendEventTags(ctx, eventId, data)
}

// sendEventTags отправляет в MISP объекты типа 'tags'
func (rmisp *requestMISP) sendEventTags(ctx context.Context, eventId string, data *objectsmispformat.ListEventObjectTags) error {
	var (
		client *ClientMISP
		err    error
	)

	client, err = NewClientMISP(rmisp.host, rmisp.userAuthKey, false)
	if err != nil {
		return supportingfunctions.CustomError(fmt.Errorf("event tags add, %w", err))
	}

	eotmf := objectsmispformat.EventObjectTagsMispFormat{}
	for _, v := range *data {
		var tagColor string = "#98bb1a"

		eotmf.Event = eventId
		eotmf.Tag = v

		if strings.Contains(strings.ToLower(v), "sensor:id") {
			tagColor = "#a70a92"
		}

		if strings.Contains(strings.ToLower(v), "sensor:object") {
			tagColor = "#c56415"
		}

		if strings.Contains(strings.ToLower(v), "class-attack") {
			tagColor = "#1535c5"
		}

		// проверка наличия тега который нужно добавить
		res, errTmp := rmisp.searchTag(ctx, v)
		if errTmp != nil {
			err = errors.Join(err, supportingfunctions.CustomError(fmt.Errorf("'event tags with id:'%s' add, %w", eventId, errTmp)))
		} else {
			if res != nil || len(res) == 0 {
				if strings.Contains(strings.ToLower(v), "sensor:id") {
					tagColor = "#a70a92"
				}

				if strings.Contains(strings.ToLower(v), "sensor:name") {
					tagColor = "#229696"
				}

				if strings.Contains(strings.ToLower(v), "sensor:object") {
					tagColor = "#c56415"
				}

				if strings.Contains(strings.ToLower(v), "class-attack") {
					tagColor = "#1535c5"
				}

				rmisp.addTag(ctx, v, tagColor)
			}
		}

		//fmt.Printf("method 'requestMISP.sendEventTags' %d. eotmf:'%+v'\n", k, eotmf)

		b, errTmp := json.Marshal(eotmf)
		if errTmp != nil {
			err = errors.Join(err, supportingfunctions.CustomError(fmt.Errorf("'event tags with id:'%s' add, %w", eventId, errTmp)))

			continue
		}

		_, _, errTmp = client.Post(ctx, "/events/addTag", b)
		if errTmp != nil {
			//fmt.Println("method 'requestMISP.sendEventTags', error:", errTmp)

			err = errors.Join(err, supportingfunctions.CustomError(fmt.Errorf("'event tags with id:'%s' add, %w", eventId, err)))

			continue
		}
	}

	return err
}

func (rmisp *requestMISP) AddTag_ForTest(ctx context.Context, tag, color string) (ResponseTagMISPFormat, error) {
	return rmisp.addTag(ctx, tag, color)
}

// addTag добавляет в MISP новый тег
func (rmisp *requestMISP) addTag(ctx context.Context, tag, color string) (ResponseTagMISPFormat, error) {
	response := ResponseTagMISPFormat{}

	c, err := NewClientMISP(rmisp.host, rmisp.masterAuthKey, false)
	if err != nil {
		return response, supportingfunctions.CustomError(fmt.Errorf("tag add, %w", err))
	}

	req, err := json.Marshal(struct {
		Name       string `json:"name"`
		Colour     string `json:"colour"`
		Exportable bool   `json:"exportable"`
	}{tag, color, true})
	if err != nil {
		return response, supportingfunctions.CustomError(fmt.Errorf("add tag, %w", err))
	}

	_, b, err := c.Post(ctx, "/tags/add/", req)
	if err != nil {
		return response, supportingfunctions.CustomError(fmt.Errorf("add tag, %w", err))
	}

	if err := json.Unmarshal(b, &response); err != nil {
		return response, err
	}

	return response, nil
}

func (rmisp *requestMISP) SearchTag_ForTest(ctx context.Context, tag string) (ResponseTagsMISPFormat, error) {
	return rmisp.searchTag(ctx, tag)
}

// searchTag поиск тегов
func (rmisp *requestMISP) searchTag(ctx context.Context, tag string) (ResponseTagsMISPFormat, error) {
	foundTags := ResponseTagsMISPFormat{}

	c, err := NewClientMISP(rmisp.host, rmisp.masterAuthKey, false)
	if err != nil {
		return foundTags, supportingfunctions.CustomError(fmt.Errorf("event publish add, %w", err))
	}

	req, err := json.Marshal(struct {
		Name string `json:"name"`
	}{tag})
	if err != nil {
		return foundTags, supportingfunctions.CustomError(fmt.Errorf("search tag, %w", err))
	}

	_, b, err := c.Post(ctx, "/tags/search/"+tag, req)
	if err != nil {
		return foundTags, supportingfunctions.CustomError(fmt.Errorf("search tag, %w", err))
	}

	if err := json.Unmarshal(b, &foundTags); err != nil {
		return foundTags, err
	}

	return foundTags, nil
}

func (rmisp *requestMISP) SendRequestPublishEvent_ForTest(ctx context.Context, eventId string) (string, error) {
	return rmisp.sendRequestPublishEvent(ctx, eventId)
}

// sendRequestPublishEvent запрос на публикацию события
func (rmisp *requestMISP) sendRequestPublishEvent(ctx context.Context, eventId string) (string, error) {
	var resultMsg string

	c, err := NewClientMISP(rmisp.host, rmisp.masterAuthKey, false)
	if err != nil {
		return resultMsg, supportingfunctions.CustomError(fmt.Errorf("event publish add, %w", err))
	}

	_, b, err := c.Post(ctx, "/events/publish/"+eventId, []byte{})
	if err != nil {
		return resultMsg, supportingfunctions.CustomError(fmt.Errorf("event publish add, %w", err))
	}

	resData := decodeResponseMispMessage(b)
	resultMsg = fmt.Sprintf("result published event with id '%s' - %s '%s' %s", eventId, resData.name, resData.message, resData.success)

	return resultMsg, nil
}

func (rmisp *requestMISP) DeleteEvent_ForTest(ctx context.Context, eventId string) error {
	return rmisp.deleteEvent(ctx, eventId)
}

// deleteEvent удаляет событие по его eventId
func (rmisp *requestMISP) deleteEvent(ctx context.Context, eventId string) error {
	ctxTimeout, CancelFunc := context.WithTimeout(ctx, time.Second*15)
	defer CancelFunc()

	c, err := NewClientMISP(rmisp.host, rmisp.masterAuthKey, false)
	if err != nil {
		return supportingfunctions.CustomError(fmt.Errorf("events delete, %w", err))
	}

	_, _, err = c.Delete(ctxTimeout, "/events/delete/"+eventId)
	if err != nil {
		return err
	}

	return nil
}
