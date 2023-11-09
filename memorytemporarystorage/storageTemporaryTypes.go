package memorytemporarystorage

import (
	"sync"
	"time"
)

// CommonStorageTemporary содержит информацию предназначенную для временного хранения
// HiveFormatMessage временное хранилище для сообщений MISP
// ListUserSettingsMISP список настроек пользователей MISP
type CommonStorageTemporary struct {
	temporaryInputCase   TemporaryInputCases
	HiveFormatMessage    HiveFormatMessages
	ListUserSettingsMISP []UserSettingsMISP
	dataCounter          DataCounterStorage
}

type TemporaryInputCases struct {
	Cases map[int]SettingsInputCase
	sync.Mutex
}

// DataCounterStorage
type DataCounterStorage struct {
	acceptedEvents       int
	processedEvents      int
	eventsDoNotMeetRules int
	eventsMeetRules      int
	startTime            time.Time
	sync.Mutex
}

// DataCounter счетчик данных
// AcceptedEvents       количество принятых событий
// ProcessedEvents      количество обработанных событий
// EventsDoNotMeetRules количество событий не соответствующих правилам
// EventsMeetRules количество событий соответствующих правилам
// StartTime время инициализации счетчика
type DataCounter struct {
	AcceptedEvents       int
	ProcessedEvents      int
	EventsDoNotMeetRules int
	EventsMeetRules      int
	StartTime            time.Time
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
	sync.Mutex
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
