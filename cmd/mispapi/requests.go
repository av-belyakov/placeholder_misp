package mispapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/av-belyakov/objectsmispformat"
	"github.com/av-belyakov/placeholder_misp/internal/supportingfunctions"
)

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

// sendEvent отправляет в API MISP событие в виде типа Event и возвращает полученный ответ
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

	if res.StatusCode != http.StatusOK {
		return nil, resBodyByte, supportingfunctions.CustomError(fmt.Errorf("events add, status '%s'", res.Status))
	}

	return res, resBodyByte, nil
}
