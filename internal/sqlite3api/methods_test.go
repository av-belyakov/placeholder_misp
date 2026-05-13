package sqlite3api_test

import (
	"context"
	"log"
	"testing"

	"github.com/av-belyakov/placeholder_misp/cmd/sqlite3api"
	"github.com/av-belyakov/placeholder_misp/commoninterfaces"
	"github.com/av-belyakov/placeholder_misp/constants"
	"github.com/av-belyakov/placeholder_misp/internal/logginghandler"
	"github.com/av-belyakov/simplelogger"
	"github.com/stretchr/testify/assert"
)

const Case_Id int = 711711

func TestMethods(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	simpleLogger, err := simplelogger.NewSimpleLogger(context.Background(), constants.Root_Dir, []simplelogger.Options{})
	if err != nil {
		log.Fatalln(err)
	}

	chZabbix := make(chan commoninterfaces.Messager)
	logging := logginghandler.New(simpleLogger, chZabbix)

	go func(ctx context.Context, log *logginghandler.LoggingChan) {
		for {
			select {
			case <-ctx.Done():
				return

			case msg := <-log.GetChan():
				t.Log("Log message:", msg)

			}
		}
	}(ctx, logging)

	module, err := sqlite3api.New(ctx, "../../backupdb/sqlite3_backup.db", logging)
	assert.NoError(t, err)
	defer module.ConnectionClose()

	t.Run("Тест 1. Добавляем запись", func(t *testing.T) {
		err := module.UpdateCaseId(ctx, Case_Id, 700001)
		assert.NoError(t, err)

		eventId, err := module.SearchCaseId(ctx, Case_Id)
		assert.NoError(t, err)
		assert.Equal(t, eventId, 700001)
	})

	t.Run("Тест 2. Обновляем запись", func(t *testing.T) {
		err := module.UpdateCaseId(ctx, Case_Id, 7000221)
		assert.NoError(t, err)

		eventId, err := module.SearchCaseId(ctx, Case_Id)
		assert.NoError(t, err)
		assert.Equal(t, eventId, 7000221)
	})

	t.Run("Тест 3. Удаляем запись", func(t *testing.T) {
		err := module.DeleteCaseId(ctx, Case_Id)
		assert.NoError(t, err)

		eventId, err := module.SearchCaseId(ctx, Case_Id)
		assert.NoError(t, err)
		assert.Equal(t, eventId, 0)
	})
}
