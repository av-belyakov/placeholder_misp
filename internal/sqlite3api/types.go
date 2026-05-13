package sqlite3api

import (
	"database/sql"

	"github.com/av-belyakov/placeholder_misp/commoninterfaces"
)

// ApiSqlite3Module модуль взаимодействия с БД
type ApiSqlite3Module struct {
	db            *sql.DB                 //дескриптор соединения с БД
	logger        commoninterfaces.Logger //логирование событий
	pathSqlite3Db string                  //путь к файлу с БД
	chRequest     chan Request            //канал для запросов к БД
}

// Request запрос к модулю
type Request struct {
	Payload    []byte
	Command    string
	ChResponse chan Response
}

// Response ответ от модуля
type Response struct {
	Payload []byte
	Error   error
}
