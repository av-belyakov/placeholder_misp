package memorytemporarystorage

import (
	"sync"
)

// CommonStorageTemporary содержит информацию предназначенную для временного хранения
// HiveFormatMessage временное хранилище для сообщений MISP
// ListUserSettingsMISP список настроек пользователей MISP
type CommonStorageTemporary struct {
	temporaryInputCase   TemporaryInputCases
	HiveFormatMessage    HiveFormatMessages
	ListUserSettingsMISP []UserSettingsMISP
	dataCounter          DataCounter
}

type TemporaryInputCases struct {
	Cases map[int]SettingsInputCase
	mutex sync.Mutex
}

type DataCounter struct {
	acceptedEvents       int
	processedEvents      int
	eventsDoNotMeetRules int
	eventsMeetRules      int
	mutex                sync.Mutex
}

type SettingsInputCase struct {
	TimeCreate int64
	EventId    string
}

type UserSettingsMISP struct {
	UserId  string
	OrgId   string
	Email   string
	AuthKey string
	OrgName string
	Role    string
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
