package supportingfunctions

import (
	"sync"
	"time"
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

func (ns *natsStorage) deleteElement(id string) {
	delete(ns.storage, id)
}
