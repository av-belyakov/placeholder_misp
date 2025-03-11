package mispapi

import (
	"encoding/json"
	"fmt"
)

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
