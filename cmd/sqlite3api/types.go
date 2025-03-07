package sqlite3api

import (
	"database/sql"

	"github.com/av-belyakov/placeholder_misp/commoninterfaces"
)

// SqliteApiModule модуль взаимодействия с БД Sqlite3
type SqliteApiModule struct {
	db             *sql.DB                 //дескриптор соединения с БД
	logger         commoninterfaces.Logger //логгер
	majorFilePath  string                  //основной файл с базой
	backupFilePath string                  //файл с базой, используется как сонова для основного
}

// routeSettings настройки маршрутизатора
type routeSettings struct {
	data         []byte
	command      string
	taskId       string
	service      string
	chanResponse chan<- ChanOutputApiSqlite
}

// ChanApiSqlite канал для взаимодействия с API SQLite
type ChanApiSqlite struct {
	Data         []byte                     //данные передаваемые в API SQLite
	Command      string                     //команда которую должен выполнить API SQLite
	TaskID       string                     //id задачи
	Service      string                     //имя сервиса, за пределами NATS, от имени которого происходит запрос (например MISP, ES)
	ChanResponse chan<- ChanOutputApiSqlite //канал для ответа
}

// ChanOutputApiSqlite
type ChanOutputApiSqlite struct {
	Data   []byte //передаваемые данные
	TaskID string //id задачи
	Status bool   //статус выполнения
}
