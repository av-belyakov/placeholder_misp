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
	ns    *natsStorage
	once  sync.Once
	mnats ModuleNATS
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

	//инициируем хранилище для дескрипторов сообщений NATS
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

	nc.Subscribe("main_caseupdate", func(m *nats.Msg) {
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

	log.Printf("Connect to NATS with address %s:%d\n", conf.Host, conf.Port)

	// обработка данных приходящих в модуль от ядра приложения
	go func() {
		for data := range mnats.chanInputNATS {
			//получаем дескриптор соединения с NATS для отправки eventId
			ncd, ok := ns.getElement(data.TaskId)
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

			tmpPackage := nrm.GetResponseMessageFromMispToTheHave()
			fmt.Printf("========= data.Command: '%s', EventId: '%s'\n", data.Command, data.EventId)
			fmt.Println("========= CREATE package for sending to TheHive, Package:")
			fmt.Printf("tmpPackage.Success: %v\n", tmpPackage.Success)
			fmt.Printf("tmpPackage.Service: %v\n", tmpPackage.Service)
			fmt.Println("tmpPackage.Commands:")
			for k, v := range tmpPackage.Commands {
				fmt.Printf("\t%d. \n\tCommand:%s\n\tName:%s\n\tString:%s\n", k, v.Command, v.Name, v.String)
			}

			/*
			   tmpPackage.Success: true
			   tmpPackage.Service: MISP
			   tmpPackage.Commands:
			   	0.
			   	Command:addtag
			   	Name:
			   	String:Webhook: send="MISP"
			   	1.
			   	Command:setcustomfield
			   	Name:misp-event-id.string
			   	String:19288
			*/

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
			}
		}
	}()

	return &mnats, nil
}
