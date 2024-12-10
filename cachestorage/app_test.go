package cachestorage_test

import (
	"context"
	"log"
	"os"
	"testing"

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
}
