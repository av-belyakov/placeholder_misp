package updateeventsmisp

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	"github.com/av-belyakov/objectsmispformat"
	"github.com/av-belyakov/placeholder_misp/cmd/mispapi"
)

func TestEditElementMisp(t *testing.T) {
	var (
		eventId string = "43940" // = case id 39100
	)

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
		events := objectsmispformat.EventsMispFormat{
			Info:              "Test Case for Belyakov AV (modified event!!!) :::TheHive caseId:'ANY-ID':::",
			Uuid:              "cde5088d-4b30-4365-bdfd-82226d2fd30f", // обязательный параметр
			Analysis:          "2",                                    // должно быть любое число из списка 0,1,2
			Timestamp:         "1778428082",                           // обязательный параметр (10 символов), должно быть старше timestamp уже загруженного события
			Date:              "2026-04-17",                           // обязательный параметр (YYYY-MM-DD)
			Distribution:      "1",
			ThreatLevelId:     "3",
			OrgId:             "~4192222",
			SightingTimestamp: "1776428092112",
			EventCreatorEmail: "a.belyakov-modified-event@cloud.gcm",
		}

		// для того что бы получить ряд основных полей нужно сделать запрос на /events/view/eventId
		// res, raw, err := client.Get(t.Context(), fmt.Sprintf("/events/view/%s", eventId), []byte{})

		b, err := json.Marshal(events)
		assert.NoError(t, err)

		res, b, err := client.Post(t.Context(), fmt.Sprintf("/events/edit/%s", eventId), b)
		assert.NoError(t, err)
		assert.Equal(t, res.StatusCode, 200)

		fmt.Println("Response:", string(b))
	})

	t.Run("Тест 2. Редактируем объект типа 'event_reports' события case", func(t *testing.T) {

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
