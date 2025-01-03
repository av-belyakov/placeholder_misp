package mispinteractions

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime"
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

func (client *ClientMISP) Get(path string, data []byte) (*http.Response, []byte, error) {
	return client.Do("GET", path, data)
}

func (client *ClientMISP) Post(path string, data []byte) (*http.Response, []byte, error) {
	return client.Do("POST", path, data)
}

func (client *ClientMISP) Delete(path string) (*http.Response, []byte, error) {
	return client.Do("DELETE", path, []byte{})
}

// Do выполняет запрос к API MISP и возвращает заголовок ответа и и тело ответа в виде среза байт
func (client *ClientMISP) Do(method, path string, data []byte) (*http.Response, []byte, error) {
	dataLen := 0
	resBodyByte := []byte{}

	httpReq := &http.Request{}
	reader := bytes.NewReader(data)
	dataLen = reader.Len()
	if dataLen > 0 && method == "POST" {
		httpReq.ContentLength = int64(dataLen)
		httpReq.Body = io.NopCloser(reader)
	}

	httpTrp := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !client.Verify},
	}

	httpReq.Method = method
	httpReq.URL = client.BaseURL
	httpReq.URL.Path = path

	httpReq.Header = http.Header{}
	httpReq.Header.Set("Authorization", client.AuthHash)
	httpReq.Header.Set("Content-type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	httpClient := http.Client{
		Transport: httpTrp,
	}

	resp, err := httpClient.Do(httpReq)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)

		return nil, resBodyByte, fmt.Errorf(" '%s' %s:%d", err.Error(), f, l-2)
	}
	defer resp.Body.Close()

	resBodyByte, err = io.ReadAll(resp.Body)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)

		return resp, resBodyByte, fmt.Errorf(" '%s' %s:%d", err.Error(), f, l-2)
	}

	if resp.StatusCode != http.StatusOK {
		_, f, l, _ := runtime.Caller(0)

		mferr := map[string]any{}
		if err := json.Unmarshal(resBodyByte, &mferr); err != nil {
			lerr := []interface{}{}
			if err := json.Unmarshal(resBodyByte, &lerr); err == nil {
				return resp, resBodyByte, fmt.Errorf("message from MISP: '%s: err - %v' %s:%d", resp.Status, lerr, f, l-1)
			}

			var serr string
			if err := json.Unmarshal(resBodyByte, &serr); err == nil {
				return resp, resBodyByte, fmt.Errorf("message from MISP: '%s: err - %s' %s:%d", resp.Status, serr, f, l-1)
			}
		}

		return resp, resBodyByte, fmt.Errorf("message from MISP: '%s: msg - %+v %s:%d", resp.Status, mferr, f, l-1)
	}

	return resp, resBodyByte, err
}
