package coremodule

import (
	"fmt"
	"placeholder_misp/datamodels"
)

type ChanInputCreateMispFormat struct {
	FieldName   string
	ValueType   string
	Value       interface{}
	FieldBranch string
}

type FieldsNameMapping struct {
	InputFieldName, MispFieldName string
}

var (
	eventsMisp         datamodels.EventsMispFormat
	attributesMisp     datamodels.AttributesMispFormat
	galaxyClustersMisp datamodels.GalaxyClustersMispFormat
	galaxyElementMisp  datamodels.GalaxyElementMispFormat
	usersMisp          datamodels.UsersMispFormat
	organizationsMisp  datamodels.OrganisationsMispFormat
	serversMisp        datamodels.ServersMispFormat
	feedsMisp          datamodels.FeedsMispFormat
	tagsMisp           datamodels.TagsMispFormat
	listFieldsMispType map[string][]FieldsNameMapping
)

func init() {
	eventsMisp = datamodels.EventsMispFormat{
		Analysis:          getAnalysis(),
		Timestamp:         "0",
		ThreatLevelId:     "4",
		PublishTimestamp:  "0",
		SightingTimestamp: "0",
	}
	attributesMisp = datamodels.AttributesMispFormat{
		Category:     "Other",
		Type:         "other",
		Timestamp:    "0",
		Distribution: "3",
		FirstSeen:    "0",
		LastSeen:     "0",
	}
	galaxyClustersMisp = datamodels.GalaxyClustersMispFormat{
		Description:   "3",
		GalaxyElement: []datamodels.GalaxyElementMispFormat{},
	}
	usersMisp = datamodels.UsersMispFormat{
		Newsread:     "0",
		ChangePw:     "0",
		CurrentLogin: "0",
		LastLogin:    "0",
		DateCreated:  "0",
		DateModified: "0",
	}
	organizationsMisp = datamodels.OrganisationsMispFormat{
		DateCreated:  "0",
		DateModified: "0",
	}
	serversMisp = datamodels.ServersMispFormat{}
	feedsMisp = datamodels.FeedsMispFormat{
		Distribution: "3",
		SourceFormat: "misp",
		InputSource:  "network",
	}
	tagsMisp = datamodels.TagsMispFormat{
		Exportable:     true,
		IsGalaxy:       true,
		IsCustomGalaxy: true,
		Inherited:      1,
	}

	/*
	   Список полей которые мне не удалось размапить.
	   	Events:
	   	- in 'непонятно' out 'tags' //однако в спецификации MISP нет такого
	   	// поля, а по коду на питон оно должно быть срезом
	   	- in 'непонятно' out 'event_reports' // однако в спецификации MISP
	   	// нет такого поля, а по коду на питон оно должно быть срезом из некоего
	   	// объекта типа EventReport

	*/

	listFieldsMispType = map[string][]FieldsNameMapping{
		"events": {
			{InputFieldName: "event.object.title", MispFieldName: "info"},
			{InputFieldName: "event.object.startDate", MispFieldName: "timestamp"},
			{InputFieldName: "event.object.tlp", MispFieldName: "distribution"},
			{InputFieldName: "event.object.severity", MispFieldName: "threat_level_id"},
			{InputFieldName: "event.organisationId", MispFieldName: "org_id"},
			{InputFieldName: "event.object.updatedAt", MispFieldName: "sighting_timestamp"},
			{InputFieldName: "event.object.owner", MispFieldName: "event_creator_email"},
		},
		"attributes": {
			{InputFieldName: "event.object.tlp", MispFieldName: "tags"},
			{InputFieldName: "observables._id", MispFieldName: "object_id"},
			{InputFieldName: "observables.data", MispFieldName: "value"},
			{InputFieldName: "observables._createdAt", MispFieldName: "timestamp"},
			{InputFieldName: "observables.message", MispFieldName: "comment"},
			{InputFieldName: "observables.startDate", MispFieldName: "first_seen"},
		},
	}
}

func NewMispFormat() (chan ChanInputCreateMispFormat, chan struct{}) {
	fmt.Println("func 'NewMispFormat', START...")

	chanInput := make(chan ChanInputCreateMispFormat)
	chanDone := make(chan struct{})

	go func(ci <-chan ChanInputCreateMispFormat, cd <-chan struct{}) {
		for {
			select {
			case tmf := <-ci:
				fmt.Println("INPUT ChanInputCreateMispFormat: ", tmf)
				//Здесь нужно сделать обработкик приема значений от decoderMessage
				// и формирование типов формата MISP

			case <-cd:
				return
			}
		}
	}(chanInput, chanDone)

	return chanInput, chanDone
}

func getAnalysis() string {
	return "2"
}

func getTagTLP(tlp int) datamodels.TagsMispFormat {
	tag := datamodels.TagsMispFormat{Name: "tlp:red", Colour: "#cc0033"}

	switch tlp {
	case 0:
		tag.Name = "tlp:white"
		tag.Colour = "#ffffff"
	case 1:
		tag.Name = "tlp:green"
		tag.Colour = "#339900"
	case 2:
		tag.Name = "tlp:amber"
		tag.Colour = "#ffc000"
	}

	return tag
}
