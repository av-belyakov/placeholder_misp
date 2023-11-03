package natsinteractions

import (
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

var (
	mnats ModuleNATS
	once  sync.Once
	ns    *natsStorage
)

type natsStorage struct {
	storage map[string]messageDescriptors
	sync.Mutex
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

	ns.Lock()
	defer ns.Unlock()

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
	ns.Lock()
	defer ns.Unlock()

	delete(ns.storage, id)
}

func init() {
	mnats.chanOutputNATS = make(chan SettingsOutputChan)
	mnats.chanInputNATS = make(chan SettingsInputChan)

	//инициируем хранилище
	ns = NewStorageNATS()
}

func NewClientNATS(
	conf confighandler.AppConfigNATS,
	storageApp *memorytemporarystorage.CommonStorageTemporary,
	logging chan<- datamodels.MessageLogging,
	counting chan<- datamodels.DataCounterSettings) (*ModuleNATS, error) {

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

	log.Printf("Connect to NATS with address %s:%d\n", conf.Host, conf.Port)

	// обработка данных приходящих в модуль от ядра приложения
	go func() {
		for data := range mnats.chanInputNATS {
			/*

			   здесь из канала нужно получить msgId, по нему найти
			   в хранилище natsStorage дискриптор для ответа с помощью

			   nm, ok := ns.setElement(m)

			   и отправить ответ назад в go-webhook с помощью
			   			nc.Publish(m.Reply, []byte("I will help you"))
			   				или
			   			m.Respond([]byte("answer is 42"))
			*/

			nrm := datamodels.NewResponseMessage()

			if data.Command == "send eventId" {
				nrm.ResponseMessageAddNewCommand(datamodels.ResponseCommandForTheHive{
					Command: "setcustomfield",
					Name:    "misp-event-id.string",
					String:  data.EventId,
				})
			}

			// Далее нужно сделать Marchal для ResponseMessageFromMispToTheHave и отправить
			// в TheHive через Webhook и NATS
			res, err := json.Marshal(nrm.GetResponseMessageFromMispToTheHave())
			if err != nil {
				_, f, l, _ := runtime.Caller(0)

				logging <- datamodels.MessageLogging{
					MsgData: fmt.Sprintf("%s %s:%d", err.Error(), f, l-2),
					MsgType: "error",
				}

				continue
			}

			//err = nc.Publish("setcustomfield", res)
			err = nc.Publish("main_caseupdate", res)
			if err != nil {
				_, f, l, _ := runtime.Caller(0)

				logging <- datamodels.MessageLogging{
					MsgData: fmt.Sprintf("%s %s:%d", err.Error(), f, l-2),
					MsgType: "error",
				}
			}
		}
	}()

	nc.Subscribe("main_caseupdate", func(m *nats.Msg) {
		mnats.chanOutputNATS <- SettingsOutputChan{
			MsgId: ns.setElement(m),
			Data:  m.Data,
		}

		/*
						похоже надо так

						nc.Subscribe("foo", func(m *nats.Msg) {
							nc.Publish(m.Reply, []byte("I will help you"))
						})

						или вот так

						// Responding to a request message
						nc.Subscribe("request", func(m *nats.Msg) {
			    			m.Respond([]byte("answer is 42"))
						})
		*/

		//сетчик принятых кейсов
		counting <- datamodels.DataCounterSettings{
			DataType: "update accepted events",
			Count:    1,
		}

	})

	return &mnats, nil
}
