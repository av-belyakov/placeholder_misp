// Пакет natsinteractions реализует методы для взаимодействия с NATS
package natsinteractions

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"

	"placeholder_misp/confighandler"
	"placeholder_misp/datamodels"
	"placeholder_misp/memorytemporarystorage"
)

const (
	ansiReset    = "\033[0m"
	ansiDarkGray = "\033[90m"
)

var (
	ns   *natsStorage
	once sync.Once
)

type natsStorage struct {
	storage map[string]messageDescriptors
	mutex   sync.Mutex
}

type messageDescriptors struct {
	timeCreate int64
	msgNats    *nats.Msg
}

func NewStorageNATS() *natsStorage {
	once.Do(func() {
		ns = &natsStorage{storage: make(map[string]messageDescriptors)}

		go checkLiveTime(ns)
	})

	return ns
}

func checkLiveTime(ns *natsStorage) {
	for range time.Tick(5 * time.Second) {
		go func() {
			ns.mutex.Lock()
			defer ns.mutex.Unlock()

			for k, v := range ns.storage {
				if time.Now().Unix() > (v.timeCreate + 360) {
					ns.deleteElement(k)
				}
			}
		}()
	}
}

func (ns *natsStorage) setElement(m *nats.Msg) string {
	id := uuid.New().String()

	ns.mutex.Lock()
	defer ns.mutex.Unlock()

	ns.storage[id] = messageDescriptors{
		timeCreate: time.Now().Unix(),
		msgNats:    m,
	}

	return id
}

func (ns *natsStorage) getElement(id string) (*nats.Msg, bool) {
	if elem, ok := ns.storage[id]; ok {
		return elem.msgNats, ok
	}

	return nil, false
}

func (ns *natsStorage) deleteElement(id string) {
	delete(ns.storage, id)
}

func NewClientNATS(
	conf confighandler.AppConfigNATS,
	confTheHive confighandler.AppConfigTheHive,
	storageApp *memorytemporarystorage.CommonStorageTemporary,
	logging chan<- datamodels.MessageLogging,
	counting chan<- datamodels.DataCounterSettings) (*ModuleNATS, error) {

	var mnats ModuleNATS = ModuleNATS{
		chanOutputNATS: make(chan SettingsOutputChan),
		chanInputNATS:  make(chan SettingsInputChan),
	}

	//инициируем хранилище для дескрипторов сообщений NATS
	ns := NewStorageNATS()

	nc, err := nats.Connect(
		fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(3*time.Second))
	_, f, l, _ := runtime.Caller(0)
	if err != nil {
		return &mnats, fmt.Errorf("'%s' %s:%d", err.Error(), f, l-4)
	}

	//обработка разрыва соединения с NATS
	nc.SetDisconnectErrHandler(func(c *nats.Conn, err error) {
		logging <- datamodels.MessageLogging{
			MsgData: fmt.Sprintf("the connection with NATS has been disconnected %s:%d", f, l-4),
			MsgType: "error",
		}
	})

	//обработка переподключения к NATS
	nc.SetReconnectHandler(func(c *nats.Conn) {
		logging <- datamodels.MessageLogging{
			MsgData: fmt.Sprintf("the connection to NATS has been re-established %s:%d", f, l-4),
			MsgType: "info",
		}
	})

	nc.Subscribe(conf.Subscriptions.SenderCase, func(m *nats.Msg) {
		// ***********************************
		// Это логирование только для теста!!!
		// ***********************************
		logging <- datamodels.MessageLogging{
			MsgData: "------------|||||| TEST_INFO func 'NewClientNATS', reseived new object ||||||------------",
			MsgType: "testing",
		}
		//
		//

		mnats.chanOutputNATS <- SettingsOutputChan{
			MsgId: ns.setElement(m),
			Data:  m.Data,
		}

		//счетчик принятых кейсов
		counting <- datamodels.DataCounterSettings{
			DataType: "update accepted events",
			Count:    1,
		}
	})

	log.Printf("%vConnect to NATS with address %s:%d%v\n", ansiDarkGray, conf.Host, conf.Port, ansiReset)

	// обработка данных приходящих в модуль от ядра приложения
	go func(chInput <-chan SettingsInputChan, toSend bool, log chan<- datamodels.MessageLogging) {
		for incomingData := range chInput {
			//не отправляем eventId в TheHive
			if !toSend {
				continue
			}

			//
			//Здесь нужно сделать обработку отправки тега и custom field
			//по новому
			//КРОМЕ ТО нужно предотвратить повторную отправку команд на добавление тегов!!!
			//

			//отправляем команды на установку тега и значения поля customFields
			go func(data SettingsInputChan, log chan<- datamodels.MessageLogging) {
				requests := map[string][]byte{
					"add_case_tag": []byte(
						fmt.Sprintf(`{
							"service": "MISP",
							"command": "add_case_tag",
							"root_id": "%s",
							"case_id": "%s",
							"value": "Webhook: send=\"MISP\""}`,
							data.RootId,
							data.CaseId)),
					"set_case_custom_field": []byte(
						fmt.Sprintf(`{
							"service": "MISP",
							"command": "set_case_custom_field",
							"root_id": "%s",
	  						"field_name": "misp-event-id.string",
							"value": "%s"}`,
							data.RootId,
							data.EventId)),
				}

				reqHandler := func(command string, req []byte) (string, error) {
					ctx, ctxCancel := context.WithTimeout(context.Background(), 300*time.Second)
					defer ctxCancel()

					res, err := nc.RequestWithContext(ctx, conf.Subscriptions.ListenerCommand, req)
					if err != nil {
						_, f, l, _ := runtime.Caller(0)
						return "", fmt.Errorf("error processing command '%s' for caseId '%s' (rootId '%s') '%s' %s:%d", command, data.CaseId, data.RootId, err.Error(), f, l-2)
					}

					resToComm := ResponseToCommand{}
					if err = json.Unmarshal(res.Data, &resToComm); err != nil {
						_, f, l, _ := runtime.Caller(0)
						return "", fmt.Errorf("%s %s:%d", err.Error(), f, l-2)
					}

					if resToComm.Error != "" {
						return "", fmt.Errorf("command '%s' for caseId '%s' (rootId '%s') return status code '%d' with error message '%s'", command, data.CaseId, data.RootId, resToComm.StatusCode, resToComm.Error)
					}

					return fmt.Sprintf("command '%s' for caseId '%s' (rootId '%s') return status code '%d'", command, data.CaseId, data.RootId, resToComm.StatusCode), nil
				}

				for command, req := range requests {
					info, err := reqHandler(command, req)
					if err != nil {
						log <- datamodels.MessageLogging{MsgType: "error", MsgData: err.Error()}

						continue
					}

					log <- datamodels.MessageLogging{MsgType: "info", MsgData: info}
				}
			}(incomingData, log)

			//*****************************************************
			//все что ниже не подходит для новой реализации
			//получаем дескриптор соединения с NATS для отправки eventId
			/*ncd, ok := ns.getElement(data.TaskId)
			if !ok {
				_, f, l, _ := runtime.Caller(0)

				logging <- datamodels.MessageLogging{
					MsgData: fmt.Sprintf("connection descriptor for task id '%s' not found %s:%d", data.TaskId, f, l-2),
					MsgType: "error",
				}

				continue
			}

			nrm := datamodels.NewResponseMessage()

			if data.Command == "send event id" {
				nrm.ResponseMessageAddNewCommand(datamodels.ResponseCommandForTheHive{
					Command: "setcustomfield",
					Name:    "misp-event-id.string",
					String:  data.EventId,
				})
			}

			res, err := json.Marshal(nrm.GetResponseMessageFromMispToTheHave())
			if err != nil {
				_, f, l, _ := runtime.Caller(0)

				logging <- datamodels.MessageLogging{
					MsgData: fmt.Sprintf("%s %s:%d", err.Error(), f, l-2),
					MsgType: "error",
				}

				continue
			}

			//отправляем в NATS пакет с eventId для добавления его в TheHive
			if err := ncd.Respond(res); err != nil {
				_, f, l, _ := runtime.Caller(0)

				logging <- datamodels.MessageLogging{
					MsgData: fmt.Sprintf("%s %s:%d", err.Error(), f, l-2),
					MsgType: "error",
				}
			}*/
		}
	}(mnats.chanInputNATS, confTheHive.Send, logging)

	return &mnats, nil
}
