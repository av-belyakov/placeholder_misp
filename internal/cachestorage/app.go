package cachestorage

import (
	"context"
	"errors"
	"fmt"
	"time"
)

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

// NewCacheStorage создает новое кэширующее хранилище. Время по истечение которого
// данные из кэша будут удалены, задается в секундах в диапазоне от 10 до 86400
// секунд, где 86400 секунд равны одним суткам.
func NewCacheStorage[T any](ctx context.Context, ttl, maxSize int) (*CacheExecutedObjects[T], error) {
	cacheExObj := &CacheExecutedObjects[T]{
		maxTtl: time.Duration(30 * time.Second),
		queue: listQueueObjects[T]{
			storages: []T(nil),
		},
		cache: cacheStorages[T]{
			maxSize:  maxSize,
			storages: map[string]storageParameters[T]{},
		},
	}

	if ttl < 10 || ttl > 86400 {
		return cacheExObj, errors.New("the lifetime of the temporary information should not be less than 10 seconds and more than 86400 seconds")
	}

	timeToLive, err := time.ParseDuration(fmt.Sprintf("%ds", ttl))
	if err != nil {
		return cacheExObj, err
	}

	cacheExObj.maxTtl = timeToLive

	go cacheExObj.automaticExecution(ctx, 10)

	return cacheExObj, err
}

//	isExecution             bool           //статус выполнения
//	isCompletedSuccessfully bool           //результат выполнения

// automaticExecution работает с очередями и кешем объектов
// С заданным интервалом времени выполняет следующие действия:
// 1. Проверка, есть ли в кеше объекты. Если кеш пустой, проверка
// наличия объектов в очереди, если она пуста - ожидание. Если очередь
// не пуста, взять объект из очереди и положить в кеш. Запуск вновь
// добавленного объекта. Если кеш не пуст, то выполняется пункт 2.
// 2. Проверка, есть ли в кеше объект со стаусом isExecution=true, если
// есть, то ожидание завершения выполнения объекта в результате которого
// меняется статус объекта на isExecution=true (может быть false?) и isCompletedSuccessfully=true.
// Если статус объекта меняется на isExecution=true (может быть false?), а isCompletedSuccessfully=false
// то выполняется повторный запуск объекта. Если статус объекта isExecution=true (может быть false?)
// и isCompletedSuccessfully=true, то выполняется пункт 3.
// 3. Проверка, есть ли в кеше свободное место, равное переменной maxSize, если
// свободное место есть, выполняется ПОИСК в кеше, по уникальному ключу, объекта
// соответствующего объекту принимаемому из очереди, при условии что, очередь
// не пуста. Если такой объект находится, то выполняется СРАВНЕНИЕ двух объектов.
// При их полном совпадении объект, полученный из очереди удаляется, а в обработку
// добавляется следующий в очереди объект.
// Вновь добавленный объект запускается на выполнение. Если в кеше нет свободного
// места переходим к пункту 4.
// 4. Проверка, есть ли в кеше объекты со статусом isExecution=false и
// isCompletedSuccessfully=true, если есть, то формируется список объектов подходящих
// под заданные условия из которых удаляется объект с самым старым timeMain. После
// удаления объекта с самым старым timeMain, обращение к очереди за новым объектом
// предназначенным для обработки.
func (cache *CacheExecutedObjects[Q]) automaticExecution(ctx context.Context, maxSize int) {
	// !!!!!!!!!!!!!
	//теперь нужно написать алгоритм работы с кешем на основе вышеперечисленных
	//пунктов и с использованием написанных, для работы с кешем, методов
	// !!!!!!!!!!!!!
}

/*
// CreateCache создает новое хранилище кэширующее исполняемые функции. Время по
// истечение которого кэшированная функция будет удалена, задается в секундах и
// варьируется в диапазоне от 10 до 86400 секунд, что эквивалентно одним суткам.
func CreateCache(ctx context.Context, ttl int) (*CacheRunningFunctions, error) {
	cacheRunCom := CacheRunningFunctions{
		ttl: time.Duration(30 * time.Second),
	}

	if ttl < 10 || ttl > 86400 {
		return &cacheRunCom, errors.New("the lifetime of the temporary information should not be less than 10 seconds and more than 86400 seconds")
	}

	timeToLive, err := time.ParseDuration(fmt.Sprintf("%ds", ttl))
	if err != nil {
		return &cacheRunCom, err
	}

	cacheRunCom.ttl = timeToLive
	cacheRunCom.cacheStorage = cacheStorageParameters{
		storages: make(map[string]storageParameters),
	}

	go cacheRunCom.automaticExecutionMethods(ctx)

	return &cacheRunCom, err
}

func (crm *CacheRunningFunctions) automaticExecutionMethods(ctx context.Context) {
	tick := time.NewTicker(5 * time.Second)

	go func(ctx context.Context, tick *time.Ticker) {
		<-ctx.Done()
		tick.Stop()
	}(ctx, tick)

	for range tick.C {
		crm.cacheStorage.mutex.RLock()
		for id, storage := range crm.cacheStorage.storages {
			fmt.Println("func 'automaticExecutionMethods' new tick:")

			//удаление слишком старых записей
			if storage.timeExpiry.Before(time.Now()) {
				go crm.DeleteElement(id)

				fmt.Println("func 'automaticExecutionMethods' new tick: before delete id:", id)

				continue
			}

			//удаление записей если функция в настоящее время не выполняется и вернула
			// положительный результат
			if !storage.isFunctionExecution && storage.isCompletedSuccessfully {
				go crm.DeleteElement(id)

				fmt.Println("func 'automaticExecutionMethods' new tick: delete id:", id)

				continue
			}

			if storage.isFunctionExecution {
				continue
			}

			//выполнение кешированной функции
			go func(cache *CacheRunningFunctions, id string, f func(int) bool) {
				fmt.Println("func 'automaticExecutionMethods' new tick: cacheFunc, id:", id)

				//устанавливает что функция выполняется
				cache.setIsFunctionExecution(id)
				//увеличивает количество попыток выполения функции на 1
				cache.increaseNumberAttempts(id)

				//при вызове, функция принимает количество попыток обработки
				if f(cache.getNumberAttempts(id)) {
					cache.setIsCompletedSuccessfully(id)
				}

				//отмечает что функция завершила выполнение
				cache.setIsFunctionNotExecution(id)
			}(crm, id, storage.cacheFunc)
		}
		crm.cacheStorage.mutex.RUnlock()
	}
}
*/
