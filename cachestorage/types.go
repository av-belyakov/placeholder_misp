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

//1. Положить в очередь.
//2. В цикле проверка кеша. Если кеш пустой или меньше максимального значения, максимальное
//значение кеша, предположим, будет 10 объектов, то следующее значение будет взято из очереди.
//Если кеш равен максимальному значению, поиск и удаление не выполняемых в настоящее время
//и выполненных успешно объектов. Удаляется самый старый по времени объект.
//3. В кеше ищется по ключу (rootId) схожий объект, если он есть, выполнетсся сравнение
//двух объектов, того что пришел и объекта из кеша.
//4. Если объект совпадает с объектом в кеш, то даный объект не обрабатывается.
//5. В кеш кладется объект, ставится статус 'выполняется'. Значения из объекта добавляется
//в MISP. При успешном добавлении ставится тригер 'успешно обработано' и статус 'не выполняется'.
//При не успешном выполнении статус 'успешно обработано' не ставится.

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
