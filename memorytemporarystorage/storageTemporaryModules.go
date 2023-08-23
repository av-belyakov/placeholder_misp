package memorytemporarystorage

import (
	"sync"
	"time"
)

func NewTemporaryStorage() *CommonStorageTemporary {
	cst := CommonStorageTemporary{
		HiveFormatMessage: HiveFormatMessages{
			Storages: make(map[string]StorageHiveFormatMessages),
			mutex:    sync.Mutex{},
		},
	}

	go checkTimeDelete(&cst)

	return &cst
}

// checkTimeDeleteTemporaryStorageSearchQueries очистка информации о задаче по истечении определенного времени или неактуальности данных
func checkTimeDelete(cst *CommonStorageTemporary) {
	c := time.Tick(3 * time.Second)

	for range c {
		for k, v := range cst.HiveFormatMessage.Storages {
			if v.isProcessedMisp && v.isProcessedElasticsearsh && v.isProcessedNKCKI {
				deleteHiveFormatMessageElement(k, cst)
			}
		}
	}
}

// SetOriginalHaveFormatMessage добавляет сырые данные полученные от TheHive
func (cst *CommonStorageTemporary) SetRawDataHiveFormatMessage(uuid string, data []byte) {
	cst.HiveFormatMessage.mutex.Lock()

	if _, ok := cst.HiveFormatMessage.Storages[uuid]; !ok {
		cst.HiveFormatMessage.Storages[uuid] = StorageHiveFormatMessages{}
	}

	storage := cst.HiveFormatMessage.Storages[uuid]

	cst.HiveFormatMessage.Storages[uuid] = StorageHiveFormatMessages{
		rawMessage:               data,
		processedMessage:         storage.processedMessage,
		allowedTransfer:          storage.allowedTransfer,
		isProcessedMisp:          storage.isProcessedMisp,
		isProcessedElasticsearsh: storage.isProcessedElasticsearsh,
		isProcessedNKCKI:         storage.isProcessedNKCKI,
	}

	cst.HiveFormatMessage.mutex.Unlock()
}

// GetRawDataHiveFormatMessage возвращает из хранилища сырые данные полученные от TheHive
func (cst *CommonStorageTemporary) GetRawDataHiveFormatMessage(uuid string) ([]byte, bool) {
	if s, ok := cst.HiveFormatMessage.Storages[uuid]; ok {
		return s.rawMessage, true
	}

	return []byte{}, false
}

// SetProcessedDataHiveFormatMessage добавляет данные частично разобранные Unmarshal JSON
func (cst *CommonStorageTemporary) SetProcessedDataHiveFormatMessage(uuid string, data map[string]interface{}) {
	cst.HiveFormatMessage.mutex.Lock()

	if _, ok := cst.HiveFormatMessage.Storages[uuid]; !ok {
		cst.HiveFormatMessage.Storages[uuid] = StorageHiveFormatMessages{}
	}

	storage := cst.HiveFormatMessage.Storages[uuid]

	cst.HiveFormatMessage.Storages[uuid] = StorageHiveFormatMessages{
		rawMessage:               storage.rawMessage,
		processedMessage:         data,
		allowedTransfer:          storage.allowedTransfer,
		isProcessedMisp:          storage.isProcessedMisp,
		isProcessedElasticsearsh: storage.isProcessedElasticsearsh,
		isProcessedNKCKI:         storage.isProcessedNKCKI,
	}

	cst.HiveFormatMessage.mutex.Unlock()
}

// GetProcessedDataHiveFormatMessage возвращает из хранилища сырые данные полученные от TheHive
func (cst *CommonStorageTemporary) GetProcessedDataHiveFormatMessage(uuid string) (map[string]interface{}, bool) {
	if s, ok := cst.HiveFormatMessage.Storages[uuid]; ok {
		return s.processedMessage, true
	}

	return map[string]interface{}{}, false
}

// SetAllowedTransferTrueHiveFormatMessage устанавливает поле информирующее о разрешении пропуска события в TRUE
func (cst *CommonStorageTemporary) SetAllowedTransferTrueHiveFormatMessage(uuid string) {
	cst.HiveFormatMessage.mutex.Lock()

	if _, ok := cst.HiveFormatMessage.Storages[uuid]; !ok {
		cst.HiveFormatMessage.Storages[uuid] = StorageHiveFormatMessages{}
	}

	storage := cst.HiveFormatMessage.Storages[uuid]

	cst.HiveFormatMessage.Storages[uuid] = StorageHiveFormatMessages{
		rawMessage:               storage.rawMessage,
		processedMessage:         storage.processedMessage,
		allowedTransfer:          true,
		isProcessedMisp:          storage.isProcessedMisp,
		isProcessedElasticsearsh: storage.isProcessedElasticsearsh,
		isProcessedNKCKI:         storage.isProcessedNKCKI,
	}

	cst.HiveFormatMessage.mutex.Unlock()
}

// SetAllowedTransferFalseHiveFormatMessage устанавливает поле информирующее о разрешении пропуска события в FALSE
func (cst *CommonStorageTemporary) SetAllowedTransferFalseHiveFormatMessage(uuid string) {
	if _, ok := cst.HiveFormatMessage.Storages[uuid]; !ok {
		cst.HiveFormatMessage.Storages[uuid] = StorageHiveFormatMessages{}
	}

	storage := cst.HiveFormatMessage.Storages[uuid]

	cst.HiveFormatMessage.Storages[uuid] = StorageHiveFormatMessages{
		rawMessage:               storage.rawMessage,
		processedMessage:         storage.processedMessage,
		allowedTransfer:          false,
		isProcessedMisp:          storage.isProcessedMisp,
		isProcessedElasticsearsh: storage.isProcessedElasticsearsh,
		isProcessedNKCKI:         storage.isProcessedNKCKI,
	}
}

// GetAllowedTransferHiveFormatMessage возвращает значение поля информирующее о разрешении пропуска события
func (cst *CommonStorageTemporary) GetAllowedTransferHiveFormatMessage(uuid string) (bool, bool) {
	if s, ok := cst.HiveFormatMessage.Storages[uuid]; ok {
		return s.allowedTransfer, true
	}

	return false, false
}

func (cst *CommonStorageTemporary) SetIsProcessedMispHiveFormatMessage(uuid string) bool {
	if _, ok := cst.HiveFormatMessage.Storages[uuid]; !ok {
		return false
	}

	storage := cst.HiveFormatMessage.Storages[uuid]

	cst.HiveFormatMessage.Storages[uuid] = StorageHiveFormatMessages{
		rawMessage:               storage.rawMessage,
		processedMessage:         storage.processedMessage,
		allowedTransfer:          storage.allowedTransfer,
		isProcessedMisp:          true,
		isProcessedElasticsearsh: storage.isProcessedElasticsearsh,
		isProcessedNKCKI:         storage.isProcessedNKCKI,
	}

	return true
}

func (cst *CommonStorageTemporary) SetIsProcessedElasticsearshHiveFormatMessage(uuid string) bool {
	if _, ok := cst.HiveFormatMessage.Storages[uuid]; !ok {
		return false
	}

	storage := cst.HiveFormatMessage.Storages[uuid]

	cst.HiveFormatMessage.Storages[uuid] = StorageHiveFormatMessages{
		rawMessage:               storage.rawMessage,
		processedMessage:         storage.processedMessage,
		allowedTransfer:          storage.allowedTransfer,
		isProcessedMisp:          storage.isProcessedMisp,
		isProcessedElasticsearsh: true,
		isProcessedNKCKI:         storage.isProcessedNKCKI,
	}

	return true
}

func (cst *CommonStorageTemporary) SetIsProcessedNKCKIHiveFormatMessage(uuid string) bool {
	if _, ok := cst.HiveFormatMessage.Storages[uuid]; !ok {
		return false
	}

	storage := cst.HiveFormatMessage.Storages[uuid]

	cst.HiveFormatMessage.Storages[uuid] = StorageHiveFormatMessages{
		rawMessage:               storage.rawMessage,
		processedMessage:         storage.processedMessage,
		allowedTransfer:          storage.allowedTransfer,
		isProcessedMisp:          storage.isProcessedMisp,
		isProcessedElasticsearsh: storage.isProcessedElasticsearsh,
		isProcessedNKCKI:         true,
	}

	return true
}

func deleteHiveFormatMessageElement(uuid string, cst *CommonStorageTemporary) {
	cst.HiveFormatMessage.mutex.Lock()

	delete(cst.HiveFormatMessage.Storages, uuid)

	cst.HiveFormatMessage.mutex.Unlock()
}
