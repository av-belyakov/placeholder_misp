package mispapi

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/av-belyakov/placeholder_misp/internal/datamodels"
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
		return nil, resBodyByte, supportingfunctions.CustomError(err)
	}
	defer resp.Body.Close()

	resBodyByte, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, resBodyByte, supportingfunctions.CustomError(err)
	}

	if resp.StatusCode != http.StatusOK {
		mferr := datamodels.MispFormatError{Errors: map[string]interface{}{}}
		if err := json.Unmarshal(resBodyByte, &mferr); err != nil {
			lerr := []interface{}{}
			if err := json.Unmarshal(resBodyByte, &lerr); err == nil {
				return resp, resBodyByte, supportingfunctions.CustomError(fmt.Errorf("message from MISP: status '%s', error - %s", resp.Status, lerr))
			}

			var serr string
			if err := json.Unmarshal(resBodyByte, &serr); err == nil {
				return resp, resBodyByte, supportingfunctions.CustomError(fmt.Errorf("message from MISP: status '%s' error - %s", resp.Status, serr))
			}
		}

		return resp, resBodyByte, supportingfunctions.CustomError(fmt.Errorf("message from MISP: staus '%s', message '%s', url '%s', error - %s", resp.Status, mferr.Message, mferr.URL, mferr.Errors))
	}

	return resp, resBodyByte, err
}
