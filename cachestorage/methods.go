package cachestorage

// AddObjectToQueue добавляет в очередь объектов новый объект
func (cache *CacheExecutedObjects) AddObjectToQueue(v listFormatsMISP) {
	cache.queueObjects.mutex.Lock()
	defer cache.queueObjects.mutex.Unlock()

	cache.queueObjects.storages = append(cache.queueObjects.storages, v)
}

// GetObjectToQueue забирает с начала очереди новый объект или возвращает
// FALSE если очередь пуста
func (cache *CacheExecutedObjects) GetObjectToQueue() (listFormatsMISP, bool) {
	cache.queueObjects.mutex.Lock()
	defer cache.queueObjects.mutex.Unlock()

	size := len(cache.queueObjects.storages)
	if size == 0 {
		return listFormatsMISP{}, false
	}

	obj := cache.queueObjects.storages[0]
	if size == 1 {
		cache.queueObjects.storages = make([]listFormatsMISP, 0, 0)
	}

	cache.queueObjects.storages = cache.queueObjects.storages[1:]

	return obj, true
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
