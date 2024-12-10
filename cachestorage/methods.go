package cachestorage

// SizeObjectToQueue размер очереди
func (c *CacheExecutedObjects[T]) SizeObjectToQueue() int {
	return len(c.queue.storages)

	// PushObjectToQueue добавляет в очередь объектов новый объект
}
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
func (c *CacheExecutedObjects[T]) AddObjectToCache(v CacheStorageFuncHandler[T]) error {
	c.cache.mutex.Lock()
	defer c.cache.mutex.Unlock()

	if len(c.cache.storages) == c.cache.maxSize {
		//запустить удаление самого старого объекта
		//при успешном удалении не должно быть ошибки
		//если есть ошибка то добавление не выполняется
	}

	return nil
}

// получить объект из хранилища
//func (c *CacheExecutedObjects[T]) GetObjectToCache(v string) CacheStorageFuncHandler[T] {
// вот здесь возвращает объект или только исполняемую функцию?
// какой объект возвращать, по старее или по новее?
//}

// deleteOldObjectFromCache удаляет самый старый объект по timeMain
func (c *CacheExecutedObjects[T]) deleteOldObjectFromCache() error {
	//как то наод еще удалять объекты у которых истекло timeExpiry

	return nil
}

//isExecution установить в TRUE
//isExecution установить в FALSE

//isCompletedSuccessfully установить в TRUE
//isCompletedSuccessfully установить в FALSE

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
