package datamodels

import (
	"fmt"
)

func (cf CustomFields) ToStringBeautiful() string {
	var str string

	for key, value := range cf {
		str += fmt.Sprintf("\t'%s':\n", key)
		for k, v := range value {
			str += fmt.Sprintf("\t\t'%s': '%v'\n", k, v)
		}
	}

	return str
}

// CustomFields настраиваемые поля
//type CustomFields map[string]map[string]interface{}
/*
str := istix.CommonPropertiesObjectSTIX.ToStringBeautiful()
	str += istix.CommonPropertiesDomainObjectSTIX.ToStringBeautiful()
	str += fmt.Sprintf("name: '%s'\n", istix.Name)
	str += fmt.Sprintf("description: '%s'\n", istix.Description)
	str += fmt.Sprintf("infrastructure_types: \n%v", func(l []OpenVocabTypeSTIX) string {
		var str string
		for k, v := range l {
			str += fmt.Sprintf("\tinfrastructure_type '%d': '%v'\n", k, v)
		}
		return str
	}(istix.InfrastructureTypes))
	str += fmt.Sprintf("aliases: \n%v", func(l []string) string {
		var str string
		for k, v := range l {
			str += fmt.Sprintf("\taliase '%d': '%s'\n", k, v)
		}
		return str
	}(istix.Aliases))
"customFields": {
                "class-attack": {
                    "order": 1,
                    "string": "Exploit"
                },
                "first-time": {
                    "date": 1630543560000,
                    "order": 0
                },
                "last-time": {
                    "date": 1630543560000,
                    "order": 0
                },
                "misp-event-id": {
                    "order": 2,
                    "string": "7481"
                },
                "ncircc-bulletine-id": {
                    "order": 0,
                    "string": "21-09-45"
                },
                "ncircc-class-attack": {
                    "order": 1,
                    "string": "Попытки эксплуатации уязвимости;attack"
                }
            },
*/
