package cachestorage

import (
	"sync"
	"time"
)

type CacheExecutedObjects[T any] struct {
	maxTTL time.Duration       //максимальное время, через которое запись в cacheStorages будет удалена
	queue  listQueueObjects[T] //очередь объектов предназначенных для выполнения
	cache  cacheStorages[T]    //кеш хранилища
}

type listQueueObjects[T any] struct {
	mutex    sync.Mutex
	storages []T
}

type cacheStorages[T any] struct {
	mutex   sync.RWMutex
	maxSize int
	//максимальный размер кеша при привышении которого выполняется удаление
	//самой старой записи
	storages map[string]storageParameters[T]
}

type storageParameters[T any] struct {
	isExecution bool
	//статус выполнения
	isCompletedSuccessfully bool
	//результат выполнения
	timeMain time.Time
	//основное время, по данному времени можно найти самый старый объект в кеше
	timeExpiry time.Time
	//общее время истечения жизни, время по истечению которого объект удаляется в любом
	//случае в независимости от того, был ли он выполнен или нет
	cacheFunc CacheStorageFuncHandler[T] //func(int) bool
	//фунция-обертка выполнения
}
