package sqlite3api_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/av-belyakov/simplelogger"
	"github.com/stretchr/testify/assert"

	"github.com/av-belyakov/placeholder_misp/cmd/sqlite3api"
	"github.com/av-belyakov/placeholder_misp/commoninterfaces"
	"github.com/av-belyakov/placeholder_misp/constants"
	"github.com/av-belyakov/placeholder_misp/internal/logginghandler"
)

var (
	module *sqlite3api.ApiSqlite3Module
)

func TestMain(m *testing.M) {
	simpleLogger, err := simplelogger.NewSimpleLogger(context.Background(), constants.Root_Dir, []simplelogger.Options{})
	if err != nil {
		log.Fatalln(err)
	}

	chZabbix := make(chan commoninterfaces.Messager)
	logging := logginghandler.New(simpleLogger, chZabbix)

	module, err = sqlite3api.New(context.Background(), "../../internal/backupdb/sqlite3_backup.db", logging)
	if err != nil {
		log.Fatalln(err)
	}

	os.Exit(m.Run())
}

func TestSqlite3Api(t *testing.T) {
	defer module.ConnectionClose()

	t.Run("Тест 1. Получаем существующую запись", func(t *testing.T) {
		chRes := make(chan sqlite3api.Response)

		module.SendDataToModule(sqlite3api.Request{
			Command:    "search caseId",
			ChResponse: chRes,
			Payload:    []byte("852"),
		})

		msg := <-chRes
		eventId := string(msg.Payload)
		assert.Equal(t, eventId, "9824")
	})

	t.Run("Тест 2. Добавляем информацию если её нет", func(t *testing.T) {
		//добавляем информацию
		module.SendDataToModule(sqlite3api.Request{
			Command:    "set case id",
			ChResponse: make(chan sqlite3api.Response),
			Payload:    []byte("999999:989898"),
		})

		//проверяем её наличие
		chRes := make(chan sqlite3api.Response)
		module.SendDataToModule(sqlite3api.Request{
			Command:    "search caseId",
			ChResponse: chRes,
			Payload:    []byte("999999"),
		})

		msg := <-chRes
		eventId := string(msg.Payload)
		assert.Equal(t, eventId, "989898")
	})

	t.Run("Тест 3. Обновляем существующую информацию", func(t *testing.T) {
		// обновляем информацию
		module.SendDataToModule(sqlite3api.Request{
			Command:    "set case id",
			ChResponse: make(chan sqlite3api.Response),
			Payload:    []byte("999999:898989"),
		})

		//проверяем результат
		chRes := make(chan sqlite3api.Response)
		module.SendDataToModule(sqlite3api.Request{
			Command:    "search caseId",
			ChResponse: chRes,
			Payload:    []byte("999999"),
		})

		msg := <-chRes
		eventId := string(msg.Payload)
		assert.Equal(t, eventId, "898989")
	})

	t.Run("Тест 4. Удаляем существующую информацию", func(t *testing.T) {
		// удаляем информацию
		module.SendDataToModule(sqlite3api.Request{
			Command:    "delete case id",
			ChResponse: make(chan sqlite3api.Response),
			Payload:    []byte("999999"),
		})

		//проверяем наличие
		chRes := make(chan sqlite3api.Response)
		module.SendDataToModule(sqlite3api.Request{
			Command:    "search caseId",
			ChResponse: chRes,
			Payload:    []byte("999999"),
		})

		msg := <-chRes
		eventId := string(msg.Payload)
		assert.Equal(t, eventId, "0")
	})
}
