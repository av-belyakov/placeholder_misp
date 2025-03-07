package sqlite3api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"runtime"

	"github.com/av-belyakov/placeholder_misp/commoninterfaces"
)

// New инициализирует новый модуль взаимодействия с API TheHive
func New(majorFilePath string, logging commoninterfaces.Logger) (*SqliteApiModule, error) {
	options := SqliteApiModule{
		logger:         logging,
		backupFilePath: "",
	}

	if majorFilePath == "" {
		return &options, errors.New("the path to the database major file database should not be empty")
	}

	options.majorFilePath = majorFilePath

	return &options, nil
}

func (opts *SqliteApiModule) Start(ctx context.Context) (chan<- ChanApiSqlite, error) {
	chanListene := make(chan ChanApiSqlite)

	sqldb, err := sql.Open("sqlite3", opts.majorFilePath)
	if err != nil {
		return chanListene, err
	}

	if sqldb.Ping() != nil {
		return chanListene, err
	}

	opts.db = sqldb

	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case msg := <-chanListene:
				opts.route(routeSettings{
					command:      msg.Command,
					taskId:       msg.TaskID,
					service:      msg.Service,
					data:         msg.Data,
					chanResponse: msg.ChanResponse,
				})
			}
		}
	}()

	return chanListene, nil
}

// Route маршрутизатор обработки запросов
func (opts *SqliteApiModule) route(settings routeSettings) {
	if settings.taskId == "" || settings.command == "" {
		_, f, l, _ := runtime.Caller(0)
		opts.logger.Send("error", fmt.Sprintf(" 'the sql query cannot be processed, the command and the task ID must not be empty' %s:%d", f, l-1))

		return
	}

	switch settings.command {
	case "insert section tags":
		go opts.handlerSectionInsertTags(settings.taskId, settings.service, settings.data, settings.chanResponse)
	case "insert section creater":
		go opts.handlerSectionInsertCreater(settings.taskId, settings.service, settings.data, settings.chanResponse)
	case "select section tags":
		go opts.handlerSectionSelectTags(settings.taskId, settings.service, settings.data, settings.chanResponse)
	case "select section creater":
		go opts.handlerSectionSelectCreater(settings.taskId, settings.service, settings.data, settings.chanResponse)
	}
}
