package memorytemporarystorage

import "time"

func NewTemporaryStorage() *CommonStorageTemporary {
	cst := CommonStorageTemporary{
		HiveFormatMessage: HiveFormatMessages{
			Storages: make(map[string]StorageHiveFormatMessages),
		},
	}

	go checkTimeDelete(&cst)

	return &cst
}

// checkTimeDeleteTemporaryStorageSearchQueries очистка информации о задаче по истечении определенного времени или неактуальности данных
func checkTimeDelete(cst *CommonStorageTemporary) {
	ticker := time.NewTicker(time.Duration(3) * time.Second)

	for range ticker.C {
		for k, v := range cst.HiveFormatMessage.Storages {
			if !v.isDelete {
				continue
			}

			deleteHiveFormatMessageElement(k, cst)
		}
	}
}

// SetOriginalHaveFormatMessage добавляет сырые данные полученные от TheHive
func (cst *CommonStorageTemporary) SetRawDataHiveFormatMessage(uuid string, data []byte) {
	cst.HiveFormatMessage.mutex.Lock()

	if _, ok := cst.HiveFormatMessage.Storages[uuid]; !ok {
		cst.HiveFormatMessage.Storages[uuid] = StorageHiveFormatMessages{}
	}

	cst.HiveFormatMessage.Storages[uuid] = StorageHiveFormatMessages{
		RawMessage: data,
	}

	cst.HiveFormatMessage.mutex.Unlock()
}

// GetRawDataHiveFormatMessage возвращает из хранилища сырые данные полученные от TheHive
func (cst *CommonStorageTemporary) GetRawDataHiveFormatMessage(uuid string) ([]byte, bool) {
	if s, ok := cst.HiveFormatMessage.Storages[uuid]; ok {
		return s.RawMessage, true
	}

	return []byte{}, false
}

// SetProcessedDataHiveFormatMessage добавляет данные частично разобранные Unmarshal JSON
func (cst *CommonStorageTemporary) SetProcessedDataHiveFormatMessage(uuid string, data map[string]interface{}) {
	cst.HiveFormatMessage.mutex.Lock()

	if _, ok := cst.HiveFormatMessage.Storages[uuid]; !ok {
		cst.HiveFormatMessage.Storages[uuid] = StorageHiveFormatMessages{}
	}

	cst.HiveFormatMessage.Storages[uuid] = StorageHiveFormatMessages{
		ProcessedMessage: data,
	}

	cst.HiveFormatMessage.mutex.Unlock()
}

// GetProcessedDataHiveFormatMessage возвращает из хранилища сырые данные полученные от TheHive
func (cst *CommonStorageTemporary) GetProcessedDataHiveFormatMessage(uuid string) (map[string]interface{}, bool) {
	if s, ok := cst.HiveFormatMessage.Storages[uuid]; ok {
		return s.ProcessedMessage, true
	}

	return map[string]interface{}{}, false
}

// SetAllowedTransferTrueHiveFormatMessage устанавливает поле информирующее о разрешении пропуска события в TRUE
func (cst *CommonStorageTemporary) SetAllowedTransferTrueHiveFormatMessage(uuid string) {
	if _, ok := cst.HiveFormatMessage.Storages[uuid]; !ok {
		cst.HiveFormatMessage.Storages[uuid] = StorageHiveFormatMessages{}
	}

	cst.HiveFormatMessage.Storages[uuid] = StorageHiveFormatMessages{
		AllowedTransfer: true,
	}
}

// SetAllowedTransferFalseHiveFormatMessage устанавливает поле информирующее о разрешении пропуска события в FALSE
func (cst *CommonStorageTemporary) SetAllowedTransferFalseHiveFormatMessage(uuid string) {
	if _, ok := cst.HiveFormatMessage.Storages[uuid]; !ok {
		cst.HiveFormatMessage.Storages[uuid] = StorageHiveFormatMessages{}
	}

	cst.HiveFormatMessage.Storages[uuid] = StorageHiveFormatMessages{
		AllowedTransfer: false,
	}
}

// GetAllowedTransferHiveFormatMessage возвращает значение поля информирующее о разрешении пропуска события
func (cst *CommonStorageTemporary) GetAllowedTransferHiveFormatMessage(uuid string) (bool, bool) {
	if s, ok := cst.HiveFormatMessage.Storages[uuid]; ok {
		return s.AllowedTransfer, true
	}

	return false, false
}

func deleteHiveFormatMessageElement(uuid string, cst *CommonStorageTemporary) {
	cst.HiveFormatMessage.mutex.Lock()

	delete(cst.HiveFormatMessage.Storages, uuid)

	cst.HiveFormatMessage.mutex.Unlock()
}
