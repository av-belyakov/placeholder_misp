package mispapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/av-belyakov/placeholder_misp/internal/supportingfunctions"
)

// удаляем дублирующиеся события из MISP
func delEventsMispFormat(ctx context.Context, host, authKey, eventId string) (*http.Response, error) {
	ctxTimeout, CancelFunc := context.WithTimeout(ctx, time.Second*15)
	defer CancelFunc()

	c, err := NewClientMISP(host, authKey, false)
	if err != nil {
		return nil, supportingfunctions.CustomError(fmt.Errorf("events delete, %w", err))
	}

	res, _, err := c.Delete(ctxTimeout, "/events/delete/"+eventId)
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
