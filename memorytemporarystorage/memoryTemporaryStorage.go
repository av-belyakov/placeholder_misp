// Пакет memorytemporarystorage реализует временное хранилище информации
package memorytemporarystorage

import (
	"sync"
	"time"
)

var once sync.Once
var cst CommonStorageTemporary

func NewTemporaryStorage() *CommonStorageTemporary {
	once.Do(func() {
		cst = CommonStorageTemporary{
			temporaryInputCase: TemporaryInputCases{
				Cases: make(map[int]SettingsInputCase),
			},
			HiveFormatMessage: HiveFormatMessages{
				Storages: make(map[string]StorageHiveFormatMessages),
			},
			ListUserSettingsMISP: make([]UserSettingsMISP, 0),
			dataCounter:          DataCounterStorage{},
		}

		go checkTimeDelete(&cst)
	})

	return &cst
}

// checkTimeDeleteTemporaryStorageSearchQueries очистка информации о задаче по истечении определенного времени или неактуальности данных
func checkTimeDelete(cst *CommonStorageTemporary) {
	c := time.Tick(3 * time.Second)

	for range c {
		go func() {
			listDel := []string(nil)

			//так как чтение Map тоже вызывает конкурентный доступ
			cst.HiveFormatMessage.mu.RLock()
			for k, v := range cst.HiveFormatMessage.Storages {
				if v.isProcessedMisp && v.isProcessedElasticsearsh && v.isProcessedNKCKI {
					listDel = append(listDel, k)
				}
			}
			cst.HiveFormatMessage.mu.RUnlock()

			if len(listDel) == 0 {
				return
			}

			cst.HiveFormatMessage.mu.Lock()
			defer cst.HiveFormatMessage.mu.Unlock()

			for _, v := range listDel {
				delete(cst.HiveFormatMessage.Storages, v)
			}
		}()

		go func() {
			listDel := []int(nil)

			// так как чтение Map тоже вызывает конкурентный доступ
			cst.temporaryInputCase.mu.RLock()
			for k, v := range cst.temporaryInputCase.Cases {
				if time.Now().Unix() > (v.TimeCreate + 54000) {
					listDel = append(listDel, k)
				}
			}
			cst.temporaryInputCase.mu.Unlock()

			if len(listDel) == 0 {
				return
			}

			cst.temporaryInputCase.mu.Lock()
			defer cst.temporaryInputCase.mu.Unlock()

			for _, v := range listDel {
				delete(cst.temporaryInputCase.Cases, v)
			}
		}()
	}
}
