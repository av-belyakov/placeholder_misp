package mispinteractions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"runtime"

	"placeholder_misp/confighandler"
	"placeholder_misp/datamodels"
	"placeholder_misp/memorytemporarystorage"
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
		chanInputMISP:  make(chan SettingsChanInputMISP),
		chanOutputMISP: make(chan interface{}),
	}
}

func HandlerMISP(
	ctx context.Context,
	conf confighandler.AppConfigMISP,
	storageApp *memorytemporarystorage.CommonStorageTemporary,
	testChan chan<- struct {
		Status     string
		StatusCode int
		Body       []byte
	},
	loging chan<- datamodels.MessageLoging) (*ModuleMISP, error) {

	//выполнеяем запрос для получения настроек пользователей через API MISP
	//и сохраняем полученные параметры во временном хранилище
	err := getUserMisp(conf, storageApp)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)

		loging <- datamodels.MessageLoging{
			MsgData: fmt.Sprintf("%s %s:%d", err.Error(), f, l-2),
			MsgType: "error",
		}
	}

	//здесь обрабатываем данные из входного канала модуля MISP
	go func() {
		for data := range mmisp.chanInputMISP {
			authKey := conf.Auth

			// получаем авторизационный ключ пользователя по его email
			if us, ok := storageApp.GetUserSettingsMISP(data.UserEmail); ok {
				authKey = us.AuthKey
			}

			//обработка только для события типа 'events'
			_, resBodyByte, err := sendEventsMispFormat(conf.Host, authKey, data)
			if err != nil {
				_, f, l, _ := runtime.Caller(0)

				fmt.Println("sendEventsMispFormat ERROR: ", err)

				loging <- datamodels.MessageLoging{
					MsgData: fmt.Sprintf("%s %s:%d", err.Error(), f, l-2),
					MsgType: "error",
				}

				continue
			}

			resMisp := RespMISP{}
			if err := json.Unmarshal(resBodyByte, &resMisp); err != nil {
				_, f, l, _ := runtime.Caller(0)

				loging <- datamodels.MessageLoging{
					MsgData: fmt.Sprintf("%s %s:%d", err.Error(), f, l-2),
					MsgType: "error",
				}

				continue
			}

			var eventId string
			for key, value := range resMisp.Event {
				if key == "id" {
					if str, ok := value.(string); ok {
						eventId = str
					}
				}
			}

			fmt.Println("EventId '", eventId, "' send to NATS")

			if eventId == "" {
				_, f, l, _ := runtime.Caller(0)

				loging <- datamodels.MessageLoging{
					MsgData: fmt.Sprintf("the formation of events of the 'Attributes' type was not performed because the EventID is empty %s:%d", f, l-1),
					MsgType: "error",
				}

				continue
			}

			res, _ := sendAttribytesMispFormat(conf.Host, authKey, eventId, data, loging)

			/*
				/ !!!!!! тут нужно отправить EventId через канал в модуль NATS !!!
				/ отправляем id добавленного в MISP события модулю NATS для передачи в TheHive

				сделал получения списка пользователей через API MISP и
				отправку json сообщений в формате Events, Attributes
				от имени пользователя сгенерирувавшего case в The Hive

				НО НАДо выполнить отладку и ПОТЕСТРОВАТЬ

			*/

			//Это тоже только для теста
			testChan <- struct {
				Status     string
				StatusCode int
				Body       []byte
			}{
				Status:     res.Status,
				StatusCode: res.StatusCode,
				//Body:       resByte,
			}
		}
	}()

	return &mmisp, nil
}

// NewClientMISP возвращает ытруктуру типа ClientMISP с предустановленными значениями
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

// getUserMisp выполнеяет запрос для получения настроек пользователей через API MISP
// и сохраняет полученные параметры во временном хранилище
func getUserMisp(conf confighandler.AppConfigMISP, storageApp *memorytemporarystorage.CommonStorageTemporary) error {
	client, err := NewClientMISP(conf.Host, conf.Auth, false)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return fmt.Errorf("%s %s:%d", err.Error(), f, l-2)
	}

	_, resByte, err := client.Get("/admin/users", nil)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return fmt.Errorf("%s %s:%d", err.Error(), f, l-2)
	}

	usmispf := []datamodels.UsersSettingsMispFormat{}
	err = json.Unmarshal(resByte, &usmispf)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return fmt.Errorf("%s %s:%d", err.Error(), f, l-2)
	}

	for _, v := range usmispf {
		storageApp.AddUserSettingsMISP(memorytemporarystorage.UserSettingsMISP{
			UserId:  v.User.Id,
			OrgId:   v.Organisation.Id,
			OrgName: v.Organisation.Name,
			Email:   v.User.Email,
			AuthKey: v.User.Authkey,
			Role:    v.Role.Name,
		})
	}

	return nil
}

// sendEventsMispFormat отправляет в API MISP событие в виде типа Event и возвращает полученный ответ
func sendEventsMispFormat(host, authKey string, d SettingsChanInputMISP) (*http.Response, []byte, error) {
	var (
		res         *http.Response
		resBodyByte = make([]byte, 0)
	)

	c, err := NewClientMISP(host, authKey, false)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return nil, resBodyByte, fmt.Errorf("%s %s:%d", err.Error(), f, l-2)
	}

	ed, ok := d.MajorData["events"]
	if !ok {
		_, f, l, _ := runtime.Caller(0)

		return nil, resBodyByte, fmt.Errorf("the properties of 'events' were not found in the received data %s:%d", f, l-2)
	}

	fmt.Println("func 'sendEventsMispFormat', EVENTS: ", ed)

	b, err := json.Marshal(ed)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)

		return nil, resBodyByte, fmt.Errorf("%s %s:%d", err.Error(), f, l-2)
	}

	res, resBodyByte, err = c.Post("/events/add", b)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)

		return nil, resBodyByte, fmt.Errorf("%s %s:%d", err.Error(), f, l-2)
	}

	return res, resBodyByte, nil
}

// sendAttribytesMispFormat отправляет в API MISP список атрибутов в виде среза типов Attribytes и возвращает полученный ответ
func sendAttribytesMispFormat(host, authKey string, eventId string, d SettingsChanInputMISP, loging chan<- datamodels.MessageLoging) (*http.Response, []byte) {
	var (
		res         *http.Response
		resBodyByte = make([]byte, 0)
	)

	c, err := NewClientMISP(host, authKey, false)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		loging <- datamodels.MessageLoging{
			MsgData: fmt.Sprintf("%s %s:%d", err.Error(), f, l-2),
			MsgType: "error",
		}

		return res, resBodyByte
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

		return nil, resBodyByte
	}

	for k := range lamf {
		lamf[k].EventId = eventId

		fmt.Println("func 'sendAttribytesMispFormat', lamf[k] = ", lamf[k])

		b, err := json.Marshal(lamf[k])
		if err != nil {
			_, f, l, _ := runtime.Caller(0)

			loging <- datamodels.MessageLoging{
				MsgData: fmt.Sprintf("%s %s:%d", err.Error(), f, l-2),
				MsgType: "error",
			}

			continue
		}

		fmt.Printf("func 'sendAttribytesMispFormat', AttributesMisp: %v\n", string(b))

		res, resBodyByte, err = c.Post("/attributes/add/"+eventId, b)
		if err != nil {
			_, f, l, _ := runtime.Caller(0)

			loging <- datamodels.MessageLoging{
				MsgData: fmt.Sprintf("%s %s:%d", err.Error(), f, l-2),
				MsgType: "error",
			}

			continue
		}
	}

	/*
				{
					"event_id":"3592",
					"object_id":"~207683696",
					"object_relation":"",
					"category":"Other",
					"type":"other",
					"value":"176.192.244.107",
					"to_ids":true,
					"uuid":"",
					"timestamp":"0",
					"distribution":"3",
					"sharing_group_id":"",
					"comment":"Download a piece of traffic",
					"first_seen":"0", //похоже дело в нуле для этого типа
					//надо в таком формате 1581984000000000 кол-во символов
					//должно быть не 13 а больше (16)
					"last_seen":"0", //похоже дело в нуле для этого типа
					"deleted":false,
					"disable_correlation":false
				}

				Возвращает ошибку 403
			{
		    	"saved": false,
		    	"name": "Could not add Attribute",
		    	"message": "Could not add Attribute",
		    	"url": "\/attributes\/add",
		    	"errors": {
		        	"first_seen": [
		            	"Invalid ISO 8601 format"
		        	],
		        	"last_seen": [
		            	"Invalid ISO 8601 format"
		        	]
		    	}
			}
	*/

	return res, resBodyByte
}
