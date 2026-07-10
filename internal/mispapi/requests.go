package mispapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/av-belyakov/objectsmispformat"
	"github.com/av-belyakov/placeholder_misp/v2/internal/supportingfunctions"
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

func (rmisp requestMISP) GetEvent_ForTest(ctx context.Context, eventId string) (int, []byte, error) {
	return rmisp.getEvent(ctx, eventId)
}

// getEvents поиск события по его идентификатору
func (rmisp requestMISP) getEvent(ctx context.Context, eventId string) (int, []byte, error) {
	if err := ctx.Err(); err != nil {
		return 0, []byte{}, supportingfunctions.CustomError(err)
	}

	client, err := NewClientMISP(rmisp.host, rmisp.userAuthKey, false)
	if err != nil {
		return 0, []byte{}, supportingfunctions.CustomError(err)
	}

	statusCode, resBodyByte, err := client.Post(ctx, fmt.Sprintf("/events/view/%s", eventId), []byte{})
	if err != nil {
		return statusCode, resBodyByte, supportingfunctions.CustomError(err)
	}

	return statusCode, resBodyByte, nil
}

func (rmisp *requestMISP) EditEvent_ForTest(ctx context.Context, eventId string, event objectsmispformat.EventsMispFormat) (int, []byte, error) {
	return rmisp.editEvent(ctx, eventId, event)
}

// editEvent редактирует событие
func (rmisp *requestMISP) editEvent(ctx context.Context, eventId string, event objectsmispformat.EventsMispFormat) (int, []byte, error) {
	if err := ctx.Err(); err != nil {
		return 0, []byte{}, supportingfunctions.CustomError(err)
	}

	b, err := json.Marshal(event)
	if err != nil {
		return 0, []byte{}, supportingfunctions.CustomError(err)
	}

	client, err := NewClientMISP(rmisp.host, rmisp.userAuthKey, false)
	if err != nil {
		return 0, []byte{}, supportingfunctions.CustomError(err)
	}

	statusCode, resBodyByte, err := client.Post(ctx, fmt.Sprintf("/events/edit/%s", eventId), b)
	if err != nil {
		return statusCode, resBodyByte, supportingfunctions.CustomError(err)
	}

	return statusCode, resBodyByte, nil
}

func (rmisp *requestMISP) SendEvent_ForTest(ctx context.Context, data *objectsmispformat.EventsMispFormat) (int, []byte, error) {
	return rmisp.sendEvent(ctx, data)
}

// sendEvent добавляет в MISP объект типа 'event'
func (rmisp *requestMISP) sendEvent(ctx context.Context, data *objectsmispformat.EventsMispFormat) (int, []byte, error) {
	if err := ctx.Err(); err != nil {
		return 0, []byte{}, supportingfunctions.CustomError(err)
	}

	b, err := json.Marshal(data)
	if err != nil {
		return 0, []byte{}, supportingfunctions.CustomError(err)
	}

	client, err := NewClientMISP(rmisp.host, rmisp.userAuthKey, false)
	if err != nil {
		return 0, []byte{}, supportingfunctions.CustomError(err)
	}

	statusCode, resBodyByte, err := client.Post(ctx, "/events/add", b)
	if err != nil {
		return statusCode, resBodyByte, supportingfunctions.CustomError(err)
	}

	return statusCode, resBodyByte, nil
}

func (rmisp *requestMISP) SendEventReports_ForTest(ctx context.Context, eventId string, data *objectsmispformat.EventReports) error {
	return rmisp.sendEventReports(ctx, eventId, data)
}

// sendEventReports добавляет в MISP объект типа 'event_reports'
func (rmisp *requestMISP) sendEventReports(ctx context.Context, eventId string, data *objectsmispformat.EventReports) error {
	if err := ctx.Err(); err != nil {
		return supportingfunctions.CustomError(err)
	}

	b, err := json.Marshal(data)
	if err != nil {
		return supportingfunctions.CustomError(err)
	}

	client, err := NewClientMISP(rmisp.host, rmisp.userAuthKey, false)
	if err != nil {
		return supportingfunctions.CustomError(err)
	}

	_, _, err = client.Post(ctx, "/event_reports/add/"+eventId, b)
	if err != nil {
		return supportingfunctions.CustomError(err)
	}

	return nil
}

func (rmisp *requestMISP) SendAttribytes_ForTest(ctx context.Context, eventId string, data []*objectsmispformat.AttributesMispFormat) (string, error) {
	return rmisp.sendAttribytes(ctx, eventId, data)
}

// sendAttribytes отправляет в MISP список атрибутов в виде среза объектов типа 'attribytes'
func (rmisp *requestMISP) sendAttribytes(ctx context.Context, eventId string, data []*objectsmispformat.AttributesMispFormat) (string, error) {
	warning := strings.Builder{}
	defer warning.Reset()

	if err := ctx.Err(); err != nil {
		return warning.String(), supportingfunctions.CustomError(err)
	}

	client, err := NewClientMISP(rmisp.host, rmisp.userAuthKey, false)
	if err != nil {
		return warning.String(), supportingfunctions.CustomError(fmt.Errorf("'attributes' for event id:'%s' add, %w", eventId, err))
	}

	for k := range data {
		data[k].EventId = eventId

		if data[k].Value == "" {
			fmt.Fprintf(&warning, "'attributes' for event id:'%s' is not added, the 'Value' type property should not be empty", eventId)

			continue
		}

		b, errTmp := json.Marshal(data[k])
		if errTmp != nil {
			err = errors.Join(err, supportingfunctions.CustomError(fmt.Errorf("'attributes' id:'%s' add, %w", eventId, errTmp)))

			continue
		}

		_, _, errTmp = client.Post(ctx, "/attributes/add/"+eventId, b)
		if errTmp != nil {
			err = errors.Join(err, supportingfunctions.CustomError(fmt.Errorf("'attributes' id:'%s' add, %w", eventId, errTmp)))

			attrObject, errTmp := json.MarshalIndent(data[k], "", "  ")
			if errTmp != nil {
				err = errors.Join(err, supportingfunctions.CustomError(fmt.Errorf("'attributes' id:'%s' add, %w", eventId, errTmp)))
			}

			fmt.Fprintf(&warning, "'attributes' with id:'%s' add, object:%s\n", eventId, string(attrObject))

			continue
		}
	}

	return warning.String(), err
}

func (rmisp *requestMISP) SendObjects_ForTest(ctx context.Context, eventId string, data map[int]*objectsmispformat.ObjectsMispFormat) error {
	return rmisp.sendObjects(ctx, eventId, data)
}

// sendObjects отправляет в MISP список объектов содержащихся в свойстве observables.attachment
// (как правило это описание вложеного файла)
func (rmisp *requestMISP) sendObjects(ctx context.Context, eventId string, data map[int]*objectsmispformat.ObjectsMispFormat) error {
	if err := ctx.Err(); err != nil {
		return supportingfunctions.CustomError(err)
	}

	client, err := NewClientMISP(rmisp.host, rmisp.userAuthKey, false)
	if err != nil {
		return supportingfunctions.CustomError(fmt.Errorf("objects for event id:'%s' add, %w", eventId, err))
	}

	for _, v := range data {
		v.EventId = eventId

		b, errTmp := json.Marshal(v)
		if errTmp != nil {
			err = errors.Join(err, supportingfunctions.CustomError(fmt.Errorf("objects with id:'%s' add, %w", eventId, errTmp)))

			continue
		}

		_, _, errTmp = client.Post(ctx, "/objects/add/"+eventId, b)
		if errTmp != nil {
			err = errors.Join(err, supportingfunctions.CustomError(fmt.Errorf("objects with id:'%s' add, %w", eventId, errTmp)))

			continue
		}
	}

	return err
}

func (rmisp *requestMISP) SendEventTags_ForTest(ctx context.Context, eventId string, data *objectsmispformat.ListEventObjectTags) error {
	return rmisp.sendEventTags(ctx, eventId, data)
}

// sendEventTags отправляет в MISP объекты типа 'tags', при этом проверяет наличие тега в MISP и если его
// не существует, то добавляет его в список тегов через rmisp.addTag() после чего добавляет тег в событие
func (rmisp *requestMISP) sendEventTags(ctx context.Context, eventId string, data *objectsmispformat.ListEventObjectTags) error {
	var (
		err        error
		objectTags objectsmispformat.EventObjectTagsMispFormat
	)

	if err := ctx.Err(); err != nil {
		return supportingfunctions.CustomError(err)
	}

	for _, v := range *data {
		var tagColor string = "#98bb1a"

		objectTags.Event = eventId
		objectTags.Tag = v

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

		if strings.Contains(strings.ToLower(v), "misp-galaxy:sector") {
			tagColor = "#e21452"
		}

		// проверка наличия тега, который нужно добавить, в списке тегов MISP (если тега нет в списке
		// MISP, то его нельзя добавить в событие)
		res, errTmp := rmisp.searchTag(ctx, v)
		if errTmp != nil {
			err = errors.Join(err, supportingfunctions.CustomError(fmt.Errorf("'event tags with event id:'%s' add, %w", eventId, errTmp)))
		} else {
			if res != nil || len(res) == 0 {
				// добавляем тег в список существующих тегов MISP
				if _, errTmp = rmisp.addTagToListTags(ctx, v, tagColor); errTmp != nil {
					err = errors.Join(err, supportingfunctions.CustomError(errTmp))
				}
			}
		}

		if err = rmisp.addTagToEvent(ctx, objectTags); err != nil {
			err = errors.Join(err, supportingfunctions.CustomError(fmt.Errorf("'event tags with id:'%s' add, %w", eventId, err)))
		}
	}

	return err
}

func (rmisp *requestMISP) AddTagToEvent_ForTest(ctx context.Context, objectTags objectsmispformat.EventObjectTagsMispFormat) error {
	return rmisp.addTagToEvent(ctx, objectTags)
}

// addTagToEvent просто добавляет тег в событие
func (rmisp *requestMISP) addTagToEvent(ctx context.Context, objectTags objectsmispformat.EventObjectTagsMispFormat) error {
	if err := ctx.Err(); err != nil {
		return supportingfunctions.CustomError(err)
	}

	b, err := json.Marshal(objectTags)
	if err != nil {
		return supportingfunctions.CustomError(err)
	}

	client, err := NewClientMISP(rmisp.host, rmisp.userAuthKey, false)
	if err != nil {
		return supportingfunctions.CustomError(err)
	}

	// добавляем тег в событие
	_, res, err := client.Post(ctx, "/events/addTag", b)
	if err != nil {
		return supportingfunctions.CustomError(err)
	}

	response := struct {
		Saved  bool   `json:"saved"`
		Errors string `json:"errors"`
	}{}
	if err := json.Unmarshal(res, &response); err != nil {
		return supportingfunctions.CustomError(err)
	}

	if !response.Saved {
		return supportingfunctions.CustomError(fmt.Errorf("response not saved: %s", response.Errors))
	}

	return nil
}

func (rmisp *requestMISP) AddTagToListTags_ForTest(ctx context.Context, tag, color string) (ResponseTagMISPFormat, error) {
	return rmisp.addTagToListTags(ctx, tag, color)
}

// addTagToListTags добавляет в список тегов новый тег (не в событие)
func (rmisp *requestMISP) addTagToListTags(ctx context.Context, tag, color string) (ResponseTagMISPFormat, error) {
	response := ResponseTagMISPFormat{}

	if err := ctx.Err(); err != nil {
		return response, supportingfunctions.CustomError(err)
	}

	req, err := json.Marshal(struct {
		Name       string `json:"name"`
		Colour     string `json:"colour"`
		Exportable bool   `json:"exportable"`
	}{tag, color, true})
	if err != nil {
		return response, supportingfunctions.CustomError(err)
	}

	client, err := NewClientMISP(rmisp.host, rmisp.masterAuthKey, false)
	if err != nil {
		return response, supportingfunctions.CustomError(err)
	}

	_, b, err := client.Post(ctx, "/tags/add/", req)
	if err != nil {
		return response, supportingfunctions.CustomError(err)
	}

	if err := json.Unmarshal(b, &response); err != nil {
		return response, supportingfunctions.CustomError(err)
	}

	return response, nil
}

func (rmisp *requestMISP) SearchTag_ForTest(ctx context.Context, tag string) (ResponseTagsMISPFormat, error) {
	return rmisp.searchTag(ctx, tag)
}

// searchTag поиск тегов
func (rmisp *requestMISP) searchTag(ctx context.Context, tag string) (ResponseTagsMISPFormat, error) {
	foundTags := ResponseTagsMISPFormat{}

	if err := ctx.Err(); err != nil {
		return foundTags, supportingfunctions.CustomError(err)
	}

	req, err := json.Marshal(struct {
		Name string `json:"name"`
	}{tag})
	if err != nil {
		return foundTags, supportingfunctions.CustomError(err)
	}

	client, err := NewClientMISP(rmisp.host, rmisp.masterAuthKey, false)
	if err != nil {
		return foundTags, supportingfunctions.CustomError(err)
	}

	_, b, err := client.Post(ctx, "/tags/search/"+tag, req)
	if err != nil {
		return foundTags, supportingfunctions.CustomError(err)
	}

	if err := json.Unmarshal(b, &foundTags); err != nil {
		return foundTags, supportingfunctions.CustomError(err)
	}

	return foundTags, nil
}

func (rmisp *requestMISP) SendRequestPublishEvent_ForTest(ctx context.Context, eventId string) (string, error) {
	return rmisp.sendRequestPublishEvent(ctx, eventId)
}

// sendRequestPublishEvent запрос на публикацию события
func (rmisp *requestMISP) sendRequestPublishEvent(ctx context.Context, eventId string) (string, error) {
	var resultMsg string

	if err := ctx.Err(); err != nil {
		return resultMsg, supportingfunctions.CustomError(err)
	}

	client, err := NewClientMISP(rmisp.host, rmisp.masterAuthKey, false)
	if err != nil {
		return resultMsg, supportingfunctions.CustomError(err)
	}

	_, b, err := client.Post(ctx, "/events/publish/"+eventId, []byte{})
	if err != nil {
		return resultMsg, supportingfunctions.CustomError(err)
	}

	resData := decodeResponseMispMessage(b)
	resultMsg = fmt.Sprintf("result published event with id '%s' - %s '%s' %s", eventId, resData.name, resData.message, resData.success)

	return resultMsg, nil
}

func (rmisp *requestMISP) SendRequestUnpublishEvent_ForTest(ctx context.Context, eventId string) (string, error) {
	return rmisp.sendRequestUnpublishEvent(ctx, eventId)
}

// sendRequestUnpublishEvent запрос на отмену публикации события
func (rmisp *requestMISP) sendRequestUnpublishEvent(ctx context.Context, eventId string) (string, error) {
	var resultMsg string

	if err := ctx.Err(); err != nil {
		return resultMsg, supportingfunctions.CustomError(err)
	}

	client, err := NewClientMISP(rmisp.host, rmisp.masterAuthKey, false)
	if err != nil {
		return resultMsg, supportingfunctions.CustomError(err)
	}

	_, b, err := client.Post(ctx, "/events/unpublish/"+eventId, []byte{})
	if err != nil {
		return resultMsg, supportingfunctions.CustomError(err)
	}

	resData := decodeResponseMispMessage(b)
	resultMsg = fmt.Sprintf("result unpublished event with id '%s' - %s '%s' %s", eventId, resData.name, resData.message, resData.success)

	return resultMsg, nil
}

func (rmisp *requestMISP) DeleteEvent_ForTest(ctx context.Context, eventId string) error {
	return rmisp.deleteEvent(ctx, eventId)
}

// deleteEvent удаляет событие по его eventId
func (rmisp *requestMISP) deleteEvent(ctx context.Context, eventId string) error {
	if err := ctx.Err(); err != nil {
		return supportingfunctions.CustomError(err)
	}

	client, err := NewClientMISP(rmisp.host, rmisp.masterAuthKey, false)
	if err != nil {
		return supportingfunctions.CustomError(err)
	}

	_, _, err = client.Delete(ctx, "/events/delete/"+eventId)
	if err != nil {
		return supportingfunctions.CustomError(err)
	}

	return nil
}

func (rmisp *requestMISP) DeleteAttributes_ForTest(ctx context.Context, attributeId, eventId string) error {
	return rmisp.deleteAttributes(ctx, attributeId, eventId)
}

// deleteAttributes удаляет атрибут привязанный к событию
func (rmisp *requestMISP) deleteAttributes(ctx context.Context, attributeId, eventId string) error {
	if err := ctx.Err(); err != nil {
		return supportingfunctions.CustomError(err)
	}

	req, err := json.Marshal(struct {
		Id      string `json:"id"`
		EventId string `json:"event_id"`
	}{
		Id:      attributeId,
		EventId: eventId,
	})
	if err != nil {
		return supportingfunctions.CustomError(err)
	}

	client, err := NewClientMISP(rmisp.host, rmisp.masterAuthKey, false)
	if err != nil {
		return supportingfunctions.CustomError(err)
	}

	if _, _, err = client.Post(ctx, fmt.Sprintf("/attributes/deleteSelected/%s", eventId), req); err != nil {
		return supportingfunctions.CustomError(err)
	}

	return nil
}

func (rmisp *requestMISP) DeleteObject_ForTest(ctx context.Context, objectId string) error {
	return rmisp.deleteObject(ctx, objectId)
}

// deleteObject удаляет объект привязанный к событию
func (rmisp *requestMISP) deleteObject(ctx context.Context, objectId string) error {
	if err := ctx.Err(); err != nil {
		return supportingfunctions.CustomError(err)
	}

	req, err := json.Marshal(struct {
		URL     string `json:"url"`
		Name    string `json:"name"`
		Message string `json:"message"`
		Saved   bool   `json:"saved"`
		Success bool   `json:"success"`
	}{
		Saved:   true,
		Success: true,
	})
	if err != nil {
		return supportingfunctions.CustomError(err)
	}

	client, err := NewClientMISP(rmisp.host, rmisp.masterAuthKey, false)
	if err != nil {
		return supportingfunctions.CustomError(err)
	}

	if _, _, err = client.Post(ctx, fmt.Sprintf("/objects/delete/%s/1", objectId), req); err != nil {
		return supportingfunctions.CustomError(err)
	}

	return nil
}

func (rmisp *requestMISP) DeleteTag_ForTest(ctx context.Context, tagId, eventId string) error {
	return rmisp.deleteTag(ctx, tagId, eventId)
}

// deleteTag удаляет тег привязанный к событию
func (rmisp *requestMISP) deleteTag(ctx context.Context, tagId, eventId string) error {
	if err := ctx.Err(); err != nil {
		return supportingfunctions.CustomError(err)
	}

	req, err := json.Marshal(struct {
		Tag   string `json:"tag"`
		Event string `json:"event"`
	}{
		Tag:   tagId,
		Event: eventId,
	})
	if err != nil {
		return supportingfunctions.CustomError(err)
	}

	client, err := NewClientMISP(rmisp.host, rmisp.masterAuthKey, false)
	if err != nil {
		return supportingfunctions.CustomError(err)
	}

	if _, _, err = client.Post(ctx, "/events/removeTag", req); err != nil {
		return supportingfunctions.CustomError(err)
	}

	return nil
}
