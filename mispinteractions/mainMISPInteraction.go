package mispinteractions

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime"

	"placeholder_misp/confighandler"
	"placeholder_misp/datamodels"
)

var mmisp ModuleMISP

type ClientMISP struct {
	BaseURL  *url.URL
	Host     string
	AuthHash string
	Verify   bool
}

type RespMISP struct {
	Event map[string]interface{} `json:"event"`
}

func init() {
	mmisp = ModuleMISP{
		chanInputMISP:  make(chan map[string]interface{}),
		chanOutputMISP: make(chan interface{}),
	}
}

func NewClientMISP(h, a string, v bool) (ClientMISP, error) {
	urlBase, err := url.Parse("http://" + h)
	if err != nil {
		return ClientMISP{}, err
	}

	return ClientMISP{
		BaseURL:  urlBase,
		Host:     h,
		AuthHash: a,
		Verify:   v,
	}, nil
}

// Get это обертка для функции Do()
func (client *ClientMISP) Get(path string, data []byte) (*http.Response, error) {
	return client.Do("GET", path, data)
}

// Post это обертка для функции Do()
func (client *ClientMISP) Post(path string, data []byte) (*http.Response, error) {
	return client.Do("POST", path, data)
}

func (client *ClientMISP) Do(method, path string, data []byte) (*http.Response, error) {
	dataLen := 0

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

	//fmt.Println("func 'Do', Method: ", method, " client.BaseURL: ", client.BaseURL, " path: ", path)
	//fmt.Println("REGUEST HEADER: ", httpReq, "httpReq.ContentLength: %d\n\n", httpReq.ContentLength)

	httpClient := http.Client{
		Transport: httpTrp,
	}
	resp, err := httpClient.Do(httpReq)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("MISP server replied status=%d", resp.StatusCode)
	}

	return resp, nil
}

func HandlerMISP(
	ctx context.Context,
	conf confighandler.AppConfigMISP,
	/*testChan chan<- struct {
		Status     string
		StatusCode int
		Body       []byte
	},*/
	loging chan<- datamodels.MessageLoging) (*ModuleMISP, error) {

	client, err := NewClientMISP(conf.Host, conf.Auth, false)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)

		loging <- datamodels.MessageLoging{
			MsgData: fmt.Sprintf("%s %s:%d", fmt.Sprint(err), f, l-2),
			MsgType: "error",
		}
	}

	//здесь обрабатываем входной канал
	go func() {
		for data := range mmisp.chanInputMISP {
			//обработка только для события типа 'events'
			if ed, ok := data["events"]; ok {
				b, err := json.Marshal(ed)
				if err != nil {
					_, f, l, _ := runtime.Caller(0)

					loging <- datamodels.MessageLoging{
						MsgData: fmt.Sprintf("%s %s:%d", fmt.Sprint(err), f, l-2),
						MsgType: "error",
					}

					continue
				}

				/*
						ТОЛЬКО ДЛЯ ТЕСТОВ
					str, err := supportingfunctions.NewReadReflectJSONSprint(b)
					if err != nil {
						fmt.Println("ERROR NewReadReflectJSONSprint:", err)
					}

					fmt.Printf("JSON string:\n%v\n", str)
				*/

				res, err := client.Post("/events/add", b)
				if err != nil {
					_, f, l, _ := runtime.Caller(0)

					loging <- datamodels.MessageLoging{
						MsgData: fmt.Sprintf("%s %s:%d", fmt.Sprint(err), f, l-2),
						MsgType: "error",
					}

					continue
				}

				//fmt.Println("___ RESPONSE HEADER:", res.Header)

				resByte, err := io.ReadAll(res.Body)
				if err != nil {
					_, f, l, _ := runtime.Caller(0)

					loging <- datamodels.MessageLoging{
						MsgData: fmt.Sprintf("%s %s:%d", fmt.Sprint(err), f, l-2),
						MsgType: "error",
					}

					continue
				}
				res.Body.Close()

				//тут надо получить id сообщения типа 'events'
				var eventId string
				resMisp := RespMISP{}
				if err := json.Unmarshal(resByte, &resMisp); err != nil {
					fmt.Println("Error create ResponseMISP: ", err)
				}

				fmt.Println("resMIsp.Event _______________ ")
				for key, value := range resMisp.Event {
					//fmt.Printf("Key: '%s' - Value: '%v'\n", key, value)

					if key == "id" {
						if str, ok := value.(string); ok {
							eventId = str
						}
					}
				}

				fmt.Println("EventId '", eventId, "' send to NATS")

				//отправляем данные в coremodule
				//тут надо подумать!!!
				//mmisp.SendingDataOutputMisp(eventId)

				//Это тоже только для теста
				testChan <- struct {
					Status     string
					StatusCode int
					Body       []byte
				}{
					Status:     res.Status,
					StatusCode: res.StatusCode,
					Body:       resByte,
				}

			}
		}
	}()

	return &mmisp, nil
}
