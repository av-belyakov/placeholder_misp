package cachestorage

import (
	"sync"
	"time"

	"placeholder_misp/datamodels"
)

type CacheExecutedObjects struct {
	maxTTL time.Duration    //максимальное время, через которое запись в cacheStorages будет удалена
	queue  listQueueObjects //очередь объектов предназначенных для выполнения
	cache  cacheStorages    //кеш хранилища
}

type listQueueObjects struct {
	mutex    sync.Mutex
	storages []FormatImplementer //listFormatsMISP
}

// listFormatsMISP содержит описание типов добавляемых в MISP и их порядок добавления.
// По результатам добавления Event, MISP возвращает id котрый необходим как для добавления
// следующих типов объектов MISP, так и для добавления этого значения в поле
// 'customFields' TheHive. Не все из этих объектов могут сразу добавлятся в MISP 'как есть',
// некоторые из них подлежат дополнительной обработке, см. обработчик для каждого из объектов.
// После добавления всех объектов, событие MISP необходимо опобликовать, как это сделать
// см. обработчик публикации.
type ListFormatsMISP struct {
	Event      datamodels.EventsMispFormat
	Reports    datamodels.EventReports
	Attributes []datamodels.AttributesMispFormat
	Objects    map[int]datamodels.ObjectsMispFormat
	ObjectTags datamodels.ListEventObjectTags
}

type cacheStorages struct {
	size     int
	mutex    sync.RWMutex
	storages map[string]storageParameters
}

type storageParameters struct {
	isExecution bool
	//статус выполнения
	isCompletedSuccessfully bool
	//результат выполнения
	timeMain time.Time
	//основное время, по данному времени можно найти самый старый объект в кеше
	timeExpiry time.Time
	//общее время истечения жизни, время по истечению которого объект удаляется в любом
	//случае в независимости от того, был ли он выполнен или нет
	cacheFunc CacheStorageFuncHandler //func(int) bool
	//фунция-обертка выполнения
}
