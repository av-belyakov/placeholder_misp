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
	c := time.Tick(5 * time.Second)

	var wg sync.WaitGroup

	for range c {
		wg.Add(2)

		go func() {
			for k, v := range cst.HiveFormatMessage.Storages {
				if v.isProcessedMisp && v.isProcessedElasticsearsh && v.isProcessedNKCKI {
					delete(cst.HiveFormatMessage.Storages, k)
				}
			}

			wg.Done()
		}()

		go func() {
			for k, v := range cst.temporaryInputCase.Cases {
				if time.Now().Unix() > (v.TimeCreate + 54000) {
					delete(cst.temporaryInputCase.Cases, k)
				}
			}

			wg.Done()
		}()

		wg.Wait()
	}
}
