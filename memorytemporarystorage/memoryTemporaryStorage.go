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
			for k, v := range copyMap(cst.HiveFormatMessage.Storages) {
				if v.isProcessedMisp && v.isProcessedElasticsearsh && v.isProcessedNKCKI {
					deleteHiveFormatMessageElement(k, cst)
				}
			}
		}()

		go func() {
			for k, v := range copyMap(cst.temporaryInputCase.Cases) {
				if time.Now().Unix() > (v.TimeCreate + 54000) {
					deleteTemporaryCase(k, cst)
				}
			}
		}()
	}
}

func copyMap[K comparable, T any](oldMap map[K]T) map[K]T {
	newMap := make(map[K]T, len(oldMap))
	for key, value := range oldMap {
		newMap[key] = value
	}

	return newMap
}
