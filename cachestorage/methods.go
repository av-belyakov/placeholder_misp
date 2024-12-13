package cachestorage

import (
	"fmt"
	"time"
)

// SizeObjectToQueue размер очереди
func (c *CacheExecutedObjects[T]) SizeObjectToQueue() int {
	return len(c.queue.storages)
}

// CleanQueue очистка очереди
func (c *CacheExecutedObjects[T]) CleanQueue() {
	c.queue.mutex.Lock()
	defer c.queue.mutex.Unlock()

	c.queue.storages = []T(nil)
}

// PushObjectToQueue добавляет в очередь объектов новый объект
func (c *CacheExecutedObjects[T]) PushObjectToQueue(v T) {
	c.queue.mutex.Lock()
	defer c.queue.mutex.Unlock()

	c.queue.storages = append(c.queue.storages, v)
}

// PullObjectToQueue забирает с начала очереди новый объект или возвращает
// FALSE если очередь пуста
func (c *CacheExecutedObjects[T]) PullObjectToQueue() (T, bool) {
	c.queue.mutex.Lock()
	defer c.queue.mutex.Unlock()

	var obj T
	size := len(c.queue.storages)
	if size == 0 {
		return obj, false
	}

	obj = c.queue.storages[0]
	if size == 1 {
		c.queue.storages = make([]T, 0)

		return obj, true
	}

	c.queue.storages = c.queue.storages[1:]

	return obj, true
}

// AddObjectToCache добавляет новый объект в хранилище
func (c *CacheExecutedObjects[T]) AddObjectToCache(key string, value CacheStorageFuncHandler[T]) error {
	c.cache.mutex.Lock()
	defer c.cache.mutex.Unlock()

	if len(c.cache.storages) >= c.cache.maxSize {
		//удаление самого старого объекта, осуществляется по параметру timeMain
		c.deleteOldObjectFromCache()
	}

	storage, ok := c.cache.storages[key]
	if !ok {
		c.cache.storages[key] = storageParameters[T]{
			timeMain:       time.Now(),
			timeExpiry:     time.Now().Add(c.maxTTL),
			originalObject: value.GetObject(),
		}
		c.cache.storages[key].cacheFunc.SetFunc(value.GetFunc())

		return nil
	}

	//найден объект у которого ключ совпадает с объектом принятом в обработку

	//объект в настоящее время выполняется
	if storage.isExecution {
		return fmt.Errorf("an object has been received whose key ID '%s' matches the already running object, ignore it", key)
	}

	//сравнение объектов из кеша и полученного из очереди
	if value.Comparison(storage.originalObject) {
		return fmt.Errorf("objects with key ID '%s' are completely identical, adding an object to the cache is not performed", key)
	}

	storage.timeMain = time.Now()
	storage.timeExpiry = time.Now().Add(c.maxTTL)
	storage.isExecution = false
	storage.isCompletedSuccessfully = false
	storage.originalObject = value.GetObject()
	storage.cacheFunc.SetFunc(value.GetFunc())

	c.cache.storages[key] = storage

	return nil
}

// GetObjectFromCacheByKey возвращает исполняемую функцию из кэша по ключу
func (c *CacheExecutedObjects[T]) GetObjectFromCacheByKey(key string) (CacheStorageFuncHandler[T], bool) {
	c.cache.mutex.RLock()
	defer c.cache.mutex.RUnlock()

	storage, ok := c.cache.storages[key]

	return storage.cacheFunc, ok
}

// GetObjectFromCacheMinTimeExpiry возвращает из кэша исполняемую функцию которая в настоящее время не
// выполняется, не была успешно выполнена и время истечения жизни объекта которой самое меньшее
func (c *CacheExecutedObjects[T]) GetObjectFromCacheMinTimeExpiry() (key string, f CacheStorageFuncHandler[T]) {
	c.cache.mutex.RLock()
	defer c.cache.mutex.RUnlock()

	var early time.Time
	for k, v := range c.cache.storages {
		if key == "" {
			key = k
			early = v.timeExpiry
			f = v.cacheFunc

			continue
		}

		if v.timeExpiry.Before(early) {
			key = k
			early = v.timeExpiry
			f = v.cacheFunc
		}
	}

	return key, f
}

// SetTimeExpiry устанавливает или обновляет значение параметра timeExpiry
func (c *CacheExecutedObjects[T]) SetTimeExpiry(key string) {
	c.cache.mutex.Lock()
	defer c.cache.mutex.Unlock()

	if storage, ok := c.cache.storages[key]; ok {
		storage.timeExpiry = time.Now().Add(c.maxTTL)
		c.cache.storages[key] = storage
	}
}

// SetIsExecutionTrue устанавливает значение параметра isExecution
func (c *CacheExecutedObjects[T]) SetIsExecutionTrue(key string) {
	c.cache.mutex.Lock()
	defer c.cache.mutex.Unlock()

	if storage, ok := c.cache.storages[key]; ok {
		storage.isExecution = true
		c.cache.storages[key] = storage
	}
}

// SetIsExecutionFalse устанавливает значение параметра isExecution
func (c *CacheExecutedObjects[T]) SetIsExecutionFalse(key string) {
	c.cache.mutex.Lock()
	defer c.cache.mutex.Unlock()

	if storage, ok := c.cache.storages[key]; ok {
		storage.isExecution = false
		c.cache.storages[key] = storage
	}
}

// SetIsCompletedSuccessfullyTrue устанавливает значение параметра isCompletedSuccessfully
func (c *CacheExecutedObjects[T]) SetIsCompletedSuccessfullyTrue(key string) {
	c.cache.mutex.Lock()
	defer c.cache.mutex.Unlock()

	if storage, ok := c.cache.storages[key]; ok {
		storage.isCompletedSuccessfully = true
		c.cache.storages[key] = storage
	}
}

// SetIsCompletedSuccessfullyFalse устанавливает значение параметра isCompletedSuccessfully
func (c *CacheExecutedObjects[T]) SetIsCompletedSuccessfullyFalse(key string) {
	c.cache.mutex.Lock()
	defer c.cache.mutex.Unlock()

	if storage, ok := c.cache.storages[key]; ok {
		storage.isCompletedSuccessfully = false
		c.cache.storages[key] = storage
	}
}

/*
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
*/

// DeleteForTimeExpiryObjectFromCache удаляет все объекты у которых истекло время жизни
func (c *CacheExecutedObjects[T]) DeleteForTimeExpiryObjectFromCache() {
	c.cache.mutex.Lock()
	defer c.cache.mutex.Unlock()

	for key, storage := range c.cache.storages {
		if storage.timeExpiry.Before(time.Now()) {
			delete(c.cache.storages, key)
		}
	}
}

// deleteOldObjectFromCache удаляет самый старый объект по timeMain
func (c *CacheExecutedObjects[T]) deleteOldObjectFromCache() {
	var (
		index      string
		timeExpiry time.Time
	)

	for k, v := range c.cache.storages {
		if index == "" {
			index = k
			timeExpiry = v.timeExpiry

			continue
		}

		if v.timeExpiry.Before(timeExpiry) {
			index = k
			timeExpiry = v.timeExpiry
		}
	}

	delete(c.cache.storages, index)
}

/*
// SetMethod создает новую запись, принимает значение которое нужно сохранить
// и id по которому данное значение можно будет найти
func (crm *CacheRunningFunctions) SetMethod(id string, f func(v int) bool) string {
	crm.cacheStorage.mutex.Lock()
	defer crm.cacheStorage.mutex.Unlock()

	crm.cacheStorage.storages[id] = storageParameters{
		timeExpiry: time.Now().Add(crm.ttl),
		cacheFunc:  f,
	}

	return id
}

// GetMethod возвращает данные по полученому id
func (crm *CacheRunningFunctions) GetMethod(id string) (func(int) bool, bool) {
	crm.cacheStorage.mutex.RLock()
	defer crm.cacheStorage.mutex.Unlock()

	if storage, ok := crm.cacheStorage.storages[id]; ok {
		return storage.cacheFunc, ok
	}

	return nil, false
}

// DeleteElement удаляет заданный элемент по его id
func (crm *CacheRunningFunctions) DeleteElement(id string) {
	crm.cacheStorage.mutex.Lock()
	defer crm.cacheStorage.mutex.Unlock()

	delete(crm.cacheStorage.storages, id)
}

// getNumberAttempts количество попыток вызова функции
func (crm *CacheRunningFunctions) getNumberAttempts(id string) int {
	storage, ok := crm.cacheStorage.storages[id]
	if !ok {
		return 0
	}

	return storage.numberAttempts
}

// increaseNumberAttempts количество попыток вызова функции
func (crm *CacheRunningFunctions) increaseNumberAttempts(id string) {
	crm.cacheStorage.mutex.Lock()
	defer crm.cacheStorage.mutex.Unlock()

	storage, ok := crm.cacheStorage.storages[id]
	if !ok {
		return
	}

	storage.numberAttempts++
	crm.cacheStorage.storages[id] = storage
}

// setIsCompletedSuccessfully выполняемая функция завершилась успехом
func (crm *CacheRunningFunctions) setIsCompletedSuccessfully(id string) {
	crm.cacheStorage.mutex.Lock()
	defer crm.cacheStorage.mutex.Unlock()

	storage, ok := crm.cacheStorage.storages[id]
	if !ok {
		return
	}

	storage.isCompletedSuccessfully = true
	crm.cacheStorage.storages[id] = storage
}

// setIsFunctionExecution функция находится в процессе выполнения
func (crm *CacheRunningFunctions) setIsFunctionExecution(id string) {
	crm.cacheStorage.mutex.Lock()
	defer crm.cacheStorage.mutex.Unlock()

	storage, ok := crm.cacheStorage.storages[id]
	if !ok {
		return
	}

	storage.isFunctionExecution = true
	crm.cacheStorage.storages[id] = storage
}

// setIsFunctionNotExecution функция не выполняется
func (crm *CacheRunningFunctions) setIsFunctionNotExecution(id string) {
	crm.cacheStorage.mutex.Lock()
	defer crm.cacheStorage.mutex.Unlock()

	storage, ok := crm.cacheStorage.storages[id]
	if !ok {
		return
	}

	storage.isFunctionExecution = false
	crm.cacheStorage.storages[id] = storage
}
*/
