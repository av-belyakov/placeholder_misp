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
// rawMessage сырые данные
// processedMessage данные частично разобранные Unmarshal JSON
// // parsingCompleted была ли завершена обработка сообщения
// allowedTransfer указывает, можно ли пропустить сообщение на дальнейшую обработку
// isProcessedMisp указывает обработан ли модулем MISP
// isProcessedElasticsearsh указывает обработан ли модулем Elasticsearch
// isProcessedNKCKI указывает обработан ли модулем NKCKI
type StorageHiveFormatMessages struct {
	rawMessage               []byte
	processedMessage         map[string]interface{}
	allowedTransfer          bool
	isProcessedMisp          bool
	isProcessedElasticsearsh bool
	isProcessedNKCKI         bool
}
