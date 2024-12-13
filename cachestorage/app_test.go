package cachestorage_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"placeholder_misp/cachestorage"
	"placeholder_misp/datamodels"
)

var (
	cache *cachestorage.CacheExecutedObjects[*datamodels.ListFormatsMISP]

	err error
)

type SpecialObjectComparator interface {
	ComparisonID(string) bool
	ComparisonEvent(*datamodels.EventsMispFormat) bool
	ComparisonReports(*datamodels.EventReports) bool
	ComparisonAttributes(*datamodels.AttributesMispFormat) bool
	ComparisonObjects(map[int]*datamodels.ObjectsMispFormat) bool
	ComparisonObjectTags(*datamodels.ListEventObjectTags) bool
	SpecialObjectGetter
}

type SpecialObjectGetter interface {
	GetID() string
	GetEvent() *datamodels.EventsMispFormat
	GetReports() *datamodels.EventReports
	GetAttributes() []*datamodels.AttributesMispFormat
	GetObjects() map[int]*datamodels.ObjectsMispFormat
	GetObjectTags() *datamodels.ListEventObjectTags
}

type SpecialObjectForCache[T SpecialObjectComparator] struct {
	object      T
	handlerFunc func(int) bool
}

func NewSpecialObjectForCache[T SpecialObjectComparator]() *SpecialObjectForCache[T] {
	return &SpecialObjectForCache[T]{}
}

func (o *SpecialObjectForCache[T]) SetObject(v T) {
	o.object = v
}

func (o *SpecialObjectForCache[T]) GetObject() T {
	return o.object
}

func (o *SpecialObjectForCache[T]) SetFunc(f func(int) bool) {
	o.handlerFunc = f
}

func (o *SpecialObjectForCache[T]) GetFunc() func(int) bool {
	return o.handlerFunc
}

func (o *SpecialObjectForCache[T]) Comparison(objFromCache T) bool {
	//выполнить сравнение
	//o.object и objFromCache

	if !o.object.ComparisonID(objFromCache.GetID()) {
		return false
	}

	if !o.object.ComparisonEvent(objFromCache.GetEvent()) {
		return false
	}

	if !o.object.ComparisonReports(objFromCache.GetReports()) {
		return false
	}

	if !o.object.ComparisonAttributes(objFromCache.GetAttributes()) {
		return false

	}

	return true
}

func TestMain(m *testing.M) {
	cache, err = cachestorage.NewCacheStorage[*datamodels.ListFormatsMISP](context.Background(), 30, 10)
	if err != nil {
		log.Fatalln(err)
	}

	os.Exit(m.Run())
}

func TestQueueHandler(t *testing.T) {
	t.Run("Тест 1. Работа с очередью", func(t *testing.T) {
		cache.PushObjectToQueue(datamodels.NewListFormatsMISP())
		cache.PushObjectToQueue(datamodels.NewListFormatsMISP())
		cache.PushObjectToQueue(datamodels.NewListFormatsMISP())

		assert.Equal(t, cache.SizeObjectToQueue(), 3)

		_, ok := cache.PullObjectToQueue()
		assert.True(t, ok)
		assert.Equal(t, cache.SizeObjectToQueue(), 2)

		_, _ = cache.PullObjectToQueue()
		_, _ = cache.PullObjectToQueue()
		assert.Equal(t, cache.SizeObjectToQueue(), 0)

		_, ok = cache.PullObjectToQueue()
		assert.False(t, ok)
	})

	t.Run("Тест 1.1. Добавить в очередь некоторое количество объектов", func(t *testing.T) {
		cache.CleanQueue()

		objectTemplate := datamodels.NewListFormatsMISP()

		objectTemplate.ID = "3255-46673"
		cache.PushObjectToQueue(objectTemplate)
		cache.PushObjectToQueue(objectTemplate) //дублирующийся объект
		objectTemplate.ID = "8483-78578"
		cache.PushObjectToQueue(objectTemplate)
		objectTemplate.ID = "3132-11223"
		cache.PushObjectToQueue(objectTemplate)
		objectTemplate.ID = "6553-13323"
		cache.PushObjectToQueue(objectTemplate)
		objectTemplate.ID = "8474-37722"
		cache.PushObjectToQueue(objectTemplate)
		objectTemplate.ID = "9123-84885"
		cache.PushObjectToQueue(objectTemplate)
		objectTemplate.ID = "1200-04993"
		cache.PushObjectToQueue(objectTemplate)
		objectTemplate.ID = "4323-29909"
		cache.PushObjectToQueue(objectTemplate)
		objectTemplate.ID = "7605-89493"
		cache.PushObjectToQueue(objectTemplate)

		assert.Equal(t, cache.SizeObjectToQueue(), 10)
	})

	t.Run("Тест 2. Добавить в кэш хранилищя некоторое количество объектов находящихся в очереди", func(t *testing.T) {
		obj, isEmpty := cache.PullObjectToQueue()
		assert.False(t, isEmpty)

		specialObject := NewSpecialObjectForCache[*datamodels.ListFormatsMISP]()
		specialObject.SetObject(obj)
		specialObject.SetFunc(func(int) bool {
			//здесь некий обработчик...
			//в контексе работы с MISP здесь должен быть код отвечающий
			//за REST запросы к серверу MISP

			return true
		})

		err := cache.AddObjectToCache(specialObject.object.ID, specialObject)
		assert.NoError(t, err)
	})

	t.Run("Тест 3. Найти и удалить самую старую запись", func(t *testing.T) {
		var (
			index      string
			timeExpiry time.Time

			timeNow time.Time            = time.Now()
			listObj map[string]time.Time = map[string]time.Time{
				"1": timeNow.Add(time.Second * 7),
				"2": timeNow.Add(time.Second * 3),
				"3": timeNow.Add(time.Second * 8),
				"4": timeNow.Add(time.Second * 1),
				"5": timeNow.Add(time.Second * 4),
				"6": timeNow.Add(time.Second * 3),
			}
		)

		for k, v := range listObj {
			if index == "" {
				index = k
				timeExpiry = v

				continue
			}

			if v.Before(timeExpiry) {
				index = k
				timeExpiry = v
			}
		}

		assert.Equal(t, index, "4")
	})
}
