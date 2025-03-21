package mispapi

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/av-belyakov/placeholder_misp/internal/supportingfunctions"
)

type ClientMISP struct {
	BaseURL  *url.URL
	Host     string
	AuthHash string
	Verify   bool
}

func (client *ClientMISP) SetAuthData(ah string) {
	client.AuthHash = ah
}

func (client *ClientMISP) GetAuthData() string {
	return client.AuthHash
}

func (client *ClientMISP) Get(ctx context.Context, path string, data []byte) (*http.Response, []byte, error) {
	return client.Do(ctx, "GET", path, data)
}

func (client *ClientMISP) Post(ctx context.Context, path string, data []byte) (*http.Response, []byte, error) {
	return client.Do(ctx, "POST", path, data)
}

func (client *ClientMISP) Delete(ctx context.Context, path string) (*http.Response, []byte, error) {
	return client.Do(ctx, "DELETE", path, []byte{})
}

// Do выполняет запрос к API MISP и возвращает заголовок ответа и и тело ответа в виде среза байт
func (client *ClientMISP) Do(ctx context.Context, method, path string, data []byte) (*http.Response, []byte, error) {
	ctxTimeout, CancelFunc := context.WithTimeout(ctx, time.Second*15)
	defer CancelFunc()

	dataLen := 0
	resBodyByte := []byte{}

	reader := bytes.NewReader(data)
	httpReq, err := http.NewRequestWithContext(ctxTimeout, method, path, reader)
	if err != nil {
		return nil, resBodyByte, supportingfunctions.CustomError(err)
	}

	dataLen = reader.Len()
	if dataLen > 0 && method == "POST" {
		httpReq.ContentLength = int64(dataLen)
		httpReq.Body = io.NopCloser(reader)
	}

	httpTrp := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !client.Verify},
	}

	httpReq.URL = client.BaseURL
	httpReq.URL.Path = path

	httpReq.Header = http.Header{}
	httpReq.Header.Set("Authorization", client.AuthHash)
	httpReq.Header.Set("Content-type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	httpClient := http.Client{
		Transport: httpTrp,
	}

	res, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, resBodyByte, supportingfunctions.CustomError(err)
	}
	defer res.Body.Close()

	resBodyByte, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, resBodyByte, supportingfunctions.CustomError(err)
	}

	if res.StatusCode != http.StatusOK {
		return res, resBodyByte, supportingfunctions.CustomError(fmt.Errorf("message from MISP: status '%s', error - %v", res.Status, string(resBodyByte)))
	}

	return res, resBodyByte, err
}
