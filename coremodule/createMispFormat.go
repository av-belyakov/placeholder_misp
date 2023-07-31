package coremodule

import (
	"fmt"
	"placeholder_misp/datamodels"
)

type ChanInputCreateMispFormat struct {
	FieldName string
	ValueType string
	Value     interface{}
}

type FieldsNameMapping struct {
	inputFieldName, outputFieldName string
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
		Analysis: getAnalysis(),
	}
	attributesMisp = datamodels.AttributesMispFormat{
		Category: "Other",
		Type:     "other",
	}
	galaxyClustersMisp = datamodels.GalaxyClustersMispFormat{}
	galaxyElementMisp = datamodels.GalaxyElementMispFormat{}
	usersMisp = datamodels.UsersMispFormat{}
	organizationsMisp = datamodels.OrganisationsMispFormat{}
	serversMisp = datamodels.ServersMispFormat{}
	feedsMisp = datamodels.FeedsMispFormat{}
	tagsMisp = datamodels.TagsMispFormat{}

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
			{inputFieldName: "title", outputFieldName: "info"},
			{inputFieldName: "startDate", outputFieldName: "timestamp"},
			{inputFieldName: "tlp", outputFieldName: "distribution"},
			{inputFieldName: "severity", outputFieldName: "threat_level_id"},
			{inputFieldName: "organisationId", outputFieldName: "org_id"},
			//{ inputFieldName: "", outputFieldName: "" },
		},
		"attributes": {
			{inputFieldName: "tlp", outputFieldName: "tags"},
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
