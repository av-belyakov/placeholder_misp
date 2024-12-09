package cachestorage

import (
	"sync"
	"time"
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
