package redisdbtosqlite3_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

//для запуска docker image с redis
//docker run -d -p 6379:6379 --volume /home/artemij/go/src/placeholder_misp/test/redis_file/dump.rdb:/data/dump.rdb --name redisdb redis:latest

const (
	HostRDb string = "127.0.0.1"
	PortRDb int    = 6379

	pathSqlite3 string = "../sqlite3_file/sqlite3.db"
)

var (
	ctx    context.Context
	cancel context.CancelFunc

	redisClient   *ConnectionRedis
	sqlite3Client *ConnectionSqlite3

	response  string
	container map[string]int = map[string]int{}

	err error
)

// ConnectionRedis
type ConnectionRedis struct {
	client *redis.Client
}

// NewConnectionRedis устанавливает новое соединение с БД
func NewConnectionRedis(host string, port int) *ConnectionRedis {
	return &ConnectionRedis{
		client: redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%d", host, port),
		})}
}

// Ping проверка соединения с БД
func (db *ConnectionRedis) Ping(ctx context.Context) (string, error) {
	status := db.client.Ping(ctx)

	return status.Result()
}

// ConnectionClose закрывает соединение с БД
func (db *ConnectionRedis) ConnectionClose() {
	db.client.Close()
}

// ConnectionSqlite3
type ConnectionSqlite3 struct {
	client *sql.DB
}

// NewConnectionSqlite3 устанавливает новое соединение с БД
func NewConnectionSqlite3(dbPath string) (*ConnectionSqlite3, error) {
	sqlite3Client, err := sql.Open("sqlite3", dbPath)

	return &ConnectionSqlite3{
		client: sqlite3Client,
	}, err
}

// Ping проверка соединения с БД
func (db *ConnectionSqlite3) Ping(ctx context.Context) error {
	return db.client.PingContext(ctx)
}

// ConnectionClose закрывает соединение с БД
func (db *ConnectionSqlite3) ConnectionClose() {
	db.client.Close()
}

func TestMain(m *testing.M) {
	ctx, cancel = context.WithCancel(context.Background())

	redisClient = NewConnectionRedis(HostRDb, PortRDb)

	sqlite3Client, err = NewConnectionSqlite3(pathSqlite3)
	if err != nil {
		log.Fatalln(err)
	}

	os.Exit(m.Run())
}

func TestGetDumpdb(t *testing.T) {
	defer func() {
		//закрываем соединение с Redis БД
		redisClient.ConnectionClose()

		//закрываем соединение с Sqlite3 БД
		sqlite3Client.ConnectionClose()

		cancel()
	}()

	t.Run("Тест 1. Получаем эхо ответ от Redis database", func(t *testing.T) {
		response, err = redisClient.Ping(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, response)
		assert.Equal(t, response, "PONG")
	})

	t.Run("Тест 2. Получаем эхо ответ от Sqlite3 database", func(t *testing.T) {
		err = sqlite3Client.Ping(ctx)
		assert.NoError(t, err)
	})

	t.Run("Тест 3. Получаем все caseId:eventId", func(t *testing.T) {
		scan := redisClient.client.Keys(ctx, "*")
		keys := []string{}
		err = scan.ScanSlice(&keys)
		assert.NoError(t, err)

		t.Log("count keys =", len(keys))

		for _, caseId := range keys {
			strCmd := redisClient.client.Get(ctx, caseId)
			eventId, err := strCmd.Int()
			assert.NoError(t, err)

			container[caseId] = eventId
		}

		assert.NotEmpty(t, container)
	})

	t.Run("Тест 4. Создаем новую таблицу", func(t *testing.T) {
		/*queryRes*/ _, err := sqlite3Client.client.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS placeholder_misp(caseId INT, eventId INT)")
		assert.NoError(t, err)
	})

	t.Run("Тест 5. Добавляем записи", func(t *testing.T) {
		record, err := sqlite3Client.client.PrepareContext(ctx, "INSERT INTO placeholder_misp (caseId, eventId) VALUES (?,?)")
		assert.NoError(t, err)

		for caseId, eventId := range container {
			cid, err := strconv.Atoi(caseId)
			assert.NoError(t, err)

			_, err = record.ExecContext(ctx, cid, eventId)
			assert.NoError(t, err)
		}
	})
}
