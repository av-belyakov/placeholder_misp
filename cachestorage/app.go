package cachestorage

import (
	"context"
	"errors"
	"fmt"
	"time"
)

/*
То что ниже надо переделать с учетом следующих особенностей:

1. Первичное добавление сформированых MISP структур должно происходить в
стек очередей, где стек формируется по принципу, первый вошел, первый вышел.
Далее все данные будут братся из этого стека и передаватся обработчику.

2. Автоматический обработчик по сигналу tick проверяет есть ли объект, который
нужно обработать. Обработчик работает только с одним единственным объектом.

3. По истечении времени жизни TTL для объекта (время жизни должно быть достачно
продолжительное) объект должен быть удален в независимости от успешности выполнения
действия "isCompletedSuccessfully" по передачи данных в MISP. Объект не удаляется
до тех пор пока статус выполнения задачи "isFunctionExecution" находится в TRUE

4. Для хранения выполняемого объекта можно использовать элемент ввиде массив с 1 элементом.

5. Пока обработка объекта выполняется флаг "isFunctionExecution" ставится в TRUE,
после выполнения ставится в FALSE.

6. После обработки, объект остается в массиве с 1 элементом и не удаляется до тех пор
пока не истечет время TTL или не появится новый объект полученный из стека объектов.
При этом старый объект (находящийся в массиве обработчике) срвынивается с новый объектом,
полученном из стека объектов, сначала по таким параметрам как caseId и rootId.
a. Если caseId и rootId совпадают, то выполняется сравнение структур и содержимого двух
объектов (для этого надо добавить в структуры методы сравнения есть в placeholder_elasticsearch).
Если все параметры и значения совпадают, то считается что это тот же объект и его дольнейшая
обработка НЕ ВЫПОЛНЯЕТСЯ.
b. Если caseId и rootId не совпадают, то выполняется замена объекта в массиве для обработке,
объектом полученным из стека объектов.

	??? как отслеживать выплнялся ли объект в массиве 1 элемента или нет, по типам isCompletedSuccessfully
и isFunctionExecution. Если isFunctionExecution в FALSE и isCompletedSuccessfully в FALSE то НЕ ВЫПОЛНЯЛСЯ или
БЫЛ ВЫПОЛНЕН НЕ УСПЕШНО, соответственно надо попробовать выполнить его ЕЩЕ РАЗ???
Если isFunctionExecution в TRUE то выполняется.
Если isFunctionExecution в FALSE и isCompletedSuccessfully в TRUE то выполнен успешно, его можно заменить.
*/

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
