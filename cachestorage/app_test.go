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

	t.Run("Тест 2. Найти и удалить самую старую запись", func(t *testing.T) {
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
