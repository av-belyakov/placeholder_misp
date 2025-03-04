package supportingfunctions

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
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

// NewStorageNATS конструктор storageNATS
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
