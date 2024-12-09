package cachestorage_test

import (
	"context"
	"log"
	"os"
	"testing"

	"placeholder_misp/cachestorage"
	"placeholder_misp/datamodels"
)

var (
	cache *cachestorage.CacheExecutedObjects

	err error
)

func TestMain(m *testing.M) {
	cache, err = cachestorage.NewCacheStorage(context.Background(), 30)
	if err != nil {
		log.Fatalln(err)
	}

	os.Exit(m.Run())
}

func TestQueueHandler(t *testing.T) {
	event := datamodels.NewEventMisp()
	//event.GetAnalysis()

	cache.AddObjectToQueue(datamodels.ListFormatsMISP{
		Event: datamodels.NewEventMisp(),
	})
}
