package mispapi

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

func decodeResponseMispMessage(b []byte) struct {
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

// создает список сенсоров получаемых из тегов
func createListSensors(listTags []string) []string {
	list := []string{}

	for _, v := range listTags {
		if !strings.Contains(strings.ToLower(v), "sensor:id") {
			continue
		}

		rgx := regexp.MustCompile(`(\w+):id=\"(\d+)\"`)
		tmp := rgx.FindStringSubmatch(v)

		if len(tmp) < 2 {
			continue
		}

		list = append(list, tmp[2])
	}

	return list
}
