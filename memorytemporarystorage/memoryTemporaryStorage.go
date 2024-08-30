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
	var wg sync.WaitGroup
	c := time.Tick(5 * time.Second)

	for range c {
		wg.Add(2)

		go func() {
			cst.HiveFormatMessage.mu.Lock()
			for k, v := range cst.HiveFormatMessage.Storages {
				if v.isProcessedMisp && v.isProcessedElasticsearsh && v.isProcessedNKCKI {
					delete(cst.HiveFormatMessage.Storages, k)
				}
			}
			cst.HiveFormatMessage.mu.Unlock()

			wg.Done()
		}()

		go func() {
			cst.temporaryInputCase.mu.Lock()
			for k, v := range cst.temporaryInputCase.Cases {
				if time.Now().Unix() > (v.TimeCreate + 54000) {
					delete(cst.temporaryInputCase.Cases, k)
				}
			}
			cst.temporaryInputCase.mu.Unlock()

			wg.Done()
		}()

		wg.Wait()
	}
}
