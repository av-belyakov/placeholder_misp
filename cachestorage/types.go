package cachestorage

import (
	"placeholder_misp/datamodels"
	"sync"
	"time"
)

// CacheRunningFunctions хранилище функций
type CacheRunningFunctions struct {
	ttl          time.Duration
	cacheStorage cacheStorageParameters
	stackObjects stackObjectsParameters
}

type cacheStorageParameters struct {
	mutex    sync.RWMutex
	storages map[string]storageParameters
}

type storageParameters struct {
	isFunctionExecution     bool
	isCompletedSuccessfully bool
	numberAttempts          int
	timeExpiry              time.Time
	cacheFunc               func(int) bool
}

type stackObjectsParameters struct {
	mutex    sync.Mutex
	storages []listFormatsMISP
}

// listFormatsMISP содержит описание типов добавляемых в MISP
// и их порядок добавления. По результатам добавления Event,
// MISP возвращает id котрый необходим как для добавления следующих
// типов объектов MISP, так и для добавления этого значения в поле
// 'customFields' TheHive.
// Не все из этих объектов могут сразу добавлятся в MISP 'как есть',
// некоторые из них подлежат дополнительной обработке, см. обработчик
// для каждого из объектов.
// После добавления всех объектов, событие MISP необходимо опобликовать,
// как это сделать см. обработчик публикации.
type listFormatsMISP struct {
	Event      datamodels.EventsMispFormat
	Reports    datamodels.EventReports
	Attributes []datamodels.AttributesMispFormat
	Objects    map[int]datamodels.ObjectsMispFormat
	ObjectTags datamodels.ListEventObjectTags
}
