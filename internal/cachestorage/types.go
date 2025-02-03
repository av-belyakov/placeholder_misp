package cachestorage

import (
	"sync"
	"time"
)

// CacheExecutedObjects кеш выполненных объектов
type CacheExecutedObjects[T any] struct {
	maxTtl time.Duration       //максимальное время, по истечении которого запись в cacheStorages будет удалена
	queue  listQueueObjects[T] //очередь объектов предназначенных для выполнения
	cache  cacheStorages[T]    //кеш хранилища обработанных объектов
}

type listQueueObjects[T any] struct {
	mutex    sync.Mutex
	storages []T
}

type cacheStorages[T any] struct {
	mutex sync.RWMutex
	//максимальный размер кеша при привышении которого выполняется удаление самой старой записи
	maxSize int
	//основное хранилище
	storages map[string]storageParameters[T]
}

type storageParameters[T any] struct {
	//статус выполнения
	isExecution bool
	//результат выполнения
	isCompletedSuccessfully bool
	//основное время, по данному времени можно найти самый старый объект в кеше
	timeMain time.Time
	//общее время истечения жизни, время по истечению которого объект удаляется в любом
	//случае в независимости от того, был ли он выполнен или нет, формируется time.Now().Add(c.maxTTL)
	timeExpiry time.Time
	//исходный объект над которым выполняются действия
	originalObject T
	//фунция-обертка выполнения
	cacheFunc func(int) bool
}
