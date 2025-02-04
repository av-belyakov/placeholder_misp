package cachestorage

import (
	"context"
	"errors"
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

/*
кэширующее хранилище произвольных объектов с очередью и обработкой этих объектов в автоматическом режиме по средствам выполнения функций переданных пользователем
*/

// NewCacheStorage создает новое кэширующее хранилище, а также очередь из которой будут, в автоматическом
// режиме, братся объекты предназначенные для обработки. Для обработки объектов будет использоватся
// пользовательская функция-обёртка, которую, как и обрабатываемый объект, добавляют с использованием
// вспомогательного пользовательского типа.
func NewCacheStorage[T any](ctx context.Context, opts ...cacheOptions[T]) (*CacheExecutedObjects[T], error) {
	cacheExObj := &CacheExecutedObjects[T]{
		//значение по умолчанию для интервала автоматической обработки
		timeTick: time.Duration(5 * time.Second),
		//значение по умолчанию для времени жизни объекта
		maxTtl: time.Duration(30 * time.Second),
		//очередь
		queue: listQueueObjects[T]{
			storages: []T(nil),
		},
		cache: cacheStorages[T]{
			//значение по умолчанию максимального размера кэша
			maxSize: 8,
			//основное хранилище
			storages: map[string]storageParameters[T]{},
		},
	}

	for _, opt := range opts {
		if err := opt(cacheExObj); err != nil {
			return cacheExObj, err
		}
	}

	go cacheExObj.automaticExecution(ctx, 10)

	return cacheExObj, nil
}

//	isExecution             bool           //статус выполнения
//	isCompletedSuccessfully bool           //результат выполнения

// automaticExecution работает с очередями и кэшем объектов
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
	//удаляем объекты время жизни которых превысело максимальное значение cache.maxTtl
	//if cache.maxTtl.storage.timeExpiry.Before(time.Now())

}

// WithMaxTtl устанавливает максимальное время, по истечении которого запись в cacheStorages будет
// удалена, допустимый интервал времени хранения записи от 10 до 86400 секунд
func WithMaxTtl[T any](v int) cacheOptions[T] {
	return func(ceo *CacheExecutedObjects[T]) error {
		if v < 10 || v > 86400 {
			return errors.New("the maximum time after which an entry in the cache will be deleted should not be less than 10 seconds or more than 24 hours (86400 seconds)")
		}

		ceo.maxTtl = time.Duration(v)

		return nil
	}
}

// WithTimeTick устанавливает интервал времени, заданное время такта, по истечении которого
// запускается новый виток автоматической обработки содержимого кэша, интервал значений должен
// быть в диапазоне от 3 до 120 секунд
func WithTimeTick[T any](v int) cacheOptions[T] {
	return func(ceo *CacheExecutedObjects[T]) error {
		if v < 3 || v > 120 {
			return errors.New("the set clock cycle time should not be less than 3 seconds or more than 120 seconds")
		}

		ceo.timeTick = time.Duration(v)

		return nil
	}
}

// WithMaxSize устанавливает максимальный размер кэша, не может быть меньше 3 и больше 1000 хранимых объектов
func WithMaxSize[T any](v int) cacheOptions[T] {
	return func(ceo *CacheExecutedObjects[T]) error {
		if v < 3 || v > 1000 {
			return errors.New("the maximum cache size cannot be less than 3 or more than 1000 objects")
		}

		ceo.cache.maxSize = v

		return nil
	}
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
