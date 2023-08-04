package memorytemporarystorage

import (
	"sync"
)

// CommonStorageTemporary содержит информацию предназначенную для временного хранения
type CommonStorageTemporary struct {
	HiveFormatMessage HiveFormatMessages
}

// HiveFormatMessages содержит временные данные относящиеся к обработки событий из TheHive
type HiveFormatMessages struct {
	Storages map[string]StorageHiveFormatMessages
	mutex    sync.Mutex
}

// StorageHiveFormatMessages хранит сообщение формата TheHave
// RawMessage сырые данные
// ProcessedMessage данные частично разобранные Unmarshal JSON
// AllowedTransfer указывает, можно ли пропустить сообщение на дальнейшую обработку
// isDelete указывает, можно ли удалить данные
type StorageHiveFormatMessages struct {
	RawMessage       []byte
	ProcessedMessage map[string]interface{}
	AllowedTransfer  bool
	isDelete         bool
}
