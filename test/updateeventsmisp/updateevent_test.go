package updateeventsmisp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	"github.com/av-belyakov/objectsmispformat"
	"github.com/av-belyakov/placeholder_misp/internal/mispapi"
)

func TestEditElementMisp(t *testing.T) {
	var (
		eventId string = "43940" // = case id 39100
		//attributeIds []string
	)

	oldEvents := struct {
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
	}{}

	/*
			{
				"id": "1542",
				"name": "Sensor:id=\"8030129\"",
				"colour": "#a70a92",
				"exportable": true,
				"user_id": "0",
				"hide_tag": false,
				"numerical_value": null,
				"is_galaxy": false,
				"is_custom_galaxy": false,
				"local": 0
		    }
	*/

	if err := godotenv.Load("../../.env"); err != nil {
		t.Fatal(err)
	}

	fmt.Printf("Host:'%s', auth token:'%s'", os.Getenv("GO_PHMISP_MHOST"), os.Getenv("GO_PHMISP_MAUTH"))

	client, err := mispapi.NewClientMISP(os.Getenv("GO_PHMISP_MHOST"), os.Getenv("GO_PHMISP_MAUTH"), false)
	if err != nil {
		t.Fatal(err)
	}

	/*
		res, raw, err := client.Get(t.Context(), fmt.Sprintf("/events/view/%s", eventId), []byte{})
		assert.NoError(t, err)
		assert.Equal(t, res.StatusCode, http.StatusOK)

		fmt.Println("Get event response:", string(raw))

		/*
				//event -> events
			"event.object.title":     {eventsMisp.SetAnyInfo},
			"event.object.startDate": {eventsMisp.SetAnyTimestamp},
			"event.details.endDate":  {eventsMisp.SetAnyDate},
			"event.object.tlp":       {eventsMisp.SetAnyDistribution},
			"event.object.severity":  {eventsMisp.SetAnyThreatLevelId},
			"event.organisationId":   {eventsMisp.SetAnyOrgId},
			"event.object.updatedAt": {eventsMisp.SetAnySightingTimestamp},
			"event.object.owner":     {eventsMisp.SetAnyEventCreatorEmail},
	*/

	t.Run("Тест 1. Редактируем основной объект 'event' события case", func(t *testing.T) {
		var eventUUID string

		t.Run("Тест 1.1. Поиск UUID редактируемого события", func(tt *testing.T) {

			res, raw, err := client.Get(t.Context(), fmt.Sprintf("/events/view/%s", eventId), []byte{})
			assert.NoError(t, err)
			assert.Equal(t, res.StatusCode, http.StatusOK)

			assert.NoError(t, json.Unmarshal(raw, &oldEvents))

			//fmt.Println("VIEW EVENT RAW:", string(raw))

			eventUUID = oldEvents.Event.UUID

			fmt.Printf("Events view: %+v", oldEvents)

			//fmt.Println("Event UUID:", eventUUID)
			//fmt.Println("Current data:", time.Now().Format("2006-01-02"), " current date unix time:", time.Now().Unix())
		})

		t.Run("Тест 1.2. Замена значения полей существующего события", func(tt *testing.T) {
			events := objectsmispformat.EventsMispFormat{
				Info:              fmt.Sprintf("Test Case for Belyakov AV (modified event!!!) :::TheHive caseId:'ANY-ID', time=%s:::", time.Now().String()),
				Uuid:              eventUUID,                       // обязательный параметр
				Analysis:          "2",                             // должно быть любое число из списка 0,1,2
				Timestamp:         fmt.Sprint(time.Now().Unix()),   //"1778428082",                         // обязательный параметр (10 символов), должно быть старше timestamp уже загруженного события
				Date:              time.Now().Format("2006-01-02"), //"2026-04-17", // обязательный параметр (YYYY-MM-DD)
				Distribution:      "1",
				ThreatLevelId:     "3",
				OrgId:             "~4192222",
				SightingTimestamp: fmt.Sprint(gofakeit.Date().UnixMicro()), //"1776428092112",
				EventCreatorEmail: "a.belyakov-modified-event@cloud.gcm",
			}

			// для того что бы получить ряд основных полей нужно сделать запрос на /events/view/eventId
			// res, raw, err := client.Get(t.Context(), fmt.Sprintf("/events/view/%s", eventId), []byte{})

			b, err := json.Marshal(events)
			assert.NoError(t, err)

			res, b, err := client.Post(t.Context(), fmt.Sprintf("/events/edit/%s", eventId), b)
			assert.NoError(t, err)
			assert.Equal(t, res.StatusCode, 200)

			//fmt.Println("Response:", string(b))
		})
	})

	t.Run("Тест 2. Удаляем список атрибутов", func(t *testing.T) {
		for _, v := range oldEvents.Event.Attribute {
			reqStruct := struct {
				Id      string `json:"id"`
				EventId string `json:"event_id"`
			}{
				Id:      fmt.Sprint(v.Id),
				EventId: eventId,
			}

			fmt.Printf("||| req:'%+v'", reqStruct)

			req, err := json.Marshal(reqStruct)
			assert.NoError(t, err)

			res, b, err := client.Post(t.Context(), fmt.Sprintf("/attributes/deleteSelected/%s", eventId), req)
			assert.NoError(t, err)

			fmt.Println("attributes deleteSelected response:", res, " byte response:", string(b))
		}
	})

	t.Run("Тест 3. Удаляем список объектов", func(t *testing.T) {
		for _, v := range oldEvents.Event.Object {
			reqStruct := struct {
				URL     string `json:"url"`
				Name    string `json:"name"`
				Message string `json:"message"`
				Saved   bool   `json:"saved"`
				Success bool   `json:"success"`
			}{
				Name:    v.Id, // это похоже не обязательно
				Saved:   true,
				Success: true,
			}

			fmt.Printf("||| req:'%+v'", reqStruct)

			req, err := json.Marshal(reqStruct)
			assert.NoError(t, err)

			res, b, err := client.Post(t.Context(), fmt.Sprintf("/objects/delete/%s/1", v.Id), req)
			assert.NoError(t, err)

			fmt.Println("objects delete response:", res, " byte response:", string(b))
		}
	})

	t.Run("Тест 4. Удаляем теги", func(t *testing.T) {
		for _, v := range oldEvents.Event.Tag {
			reqStruct := struct {
				Tag   string `json:"tag"`
				Event string `json:"event"`
			}{
				Tag:   v.Id,
				Event: eventId,
			}

			fmt.Printf("||| req:'%+v'", reqStruct)

			req, err := json.Marshal(reqStruct)
			assert.NoError(t, err)

			res, b, err := client.Post(t.Context(), "/events/removeTag", req)
			assert.NoError(t, err)

			fmt.Println("tags delete response:", res, " byte response:", string(b))
		}
	})

	// перечень обновляемых или добавляемых объектов
	//
	// добавляет в MISP объект типа 'event_reports'
	// "/event_reports/add/"+eventId
	// 	 ДЛЯ ЭТОГО ОБЪЕКТА НЕТ МЕТОДА DELETE, есть только add и edit
	//
	// отправляет в MISP список атрибутов в виде среза объектов типа 'attribytes'
	// "/attributes/add/"+eventId
	// 	ДЛЯ ЭТОГО ОБЪЕКТА ЕСТЬ МЕТОД DELETE
	//
	// отправляет в MISP список объектов содержащихся в свойстве observables.attachment
	// (как правило это описание вложеного файла)
	// "/objects/add/"+eventId
	// ДЛЯ ЭТОГО ОБЪЕКТА ВООБЩЕ НЕТ ОПИСАНИЯ API
	//
	// отправляет в MISP объекты типа 'tags'
	// "/events/addTag"
	// 	ДЛЯ ЭТОГО ОБЪЕКТА ЕСТЬ МЕТОД remoteTag
	//
	// запрос на публикацию события
	// "/events/publish/"+eventId
	// как только событие event изменяется оно считается не опубликованным

	//t.Run("", func(t *testing.T) {})

	t.Cleanup(func() {
		os.Unsetenv("GO_PHMISP_MHOST")
		os.Unsetenv("GO_PHMISP_MAUTH")
	})
}
