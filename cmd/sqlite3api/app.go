package sqlite3api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"github.com/av-belyakov/placeholder_misp/commoninterfaces"
	"github.com/av-belyakov/placeholder_misp/internal/supportingfunctions"
)

// New инициализирует новый модуль взаимодействия с Sqlite3 API
func New(ctx context.Context, pathDb string, logging commoninterfaces.Logger) (*ApiSqlite3Module, error) {
	module := &ApiSqlite3Module{
		logger:    logging,
		chRequest: make(chan Request)}

	if pathDb == "" {
		return module, supportingfunctions.CustomError(errors.New("the pathDb parameter must not be empty"))
	}

	sqlite3Client, err := sql.Open("sqlite3", pathDb)
	if err != nil {
		return module, supportingfunctions.CustomError(fmt.Errorf("module Sqlite3 API, %w", err))
	}

	module.pathSqlite3Db = pathDb
	module.db = sqlite3Client

	if err = module.Ping(ctx); err != nil {
		return module, err
	}

	module.route(ctx)

	go func(ctx context.Context, m *ApiSqlite3Module) {
		<-ctx.Done()
		m.ConnectionClose()
	}(ctx, module)

	return module, nil
}

func (module *ApiSqlite3Module) route(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case data := <-module.GetChRequest():
				switch data.Command {
				case "search caseId":
					str := string(data.Payload)
					caseId, err := strconv.Atoi(str)
					if err != nil {
						data.ChResponse <- Response{Error: err}
						module.logger.Send("error", supportingfunctions.CustomError(err).Error())

						continue
					}

					res, err := module.SearchCaseId(ctx, caseId)
					if err != nil {
						data.ChResponse <- Response{Error: err}
						module.logger.Send("error", supportingfunctions.CustomError(err).Error())

						continue
					}

					data.ChResponse <- Response{Payload: fmt.Append(nil, res)}

				case "set case id":
					tmp := strings.Split(string(data.Payload), ":")
					if len(tmp) == 0 {
						module.logger.Send("warning", supportingfunctions.CustomError(errors.New("it is not possible to split a string")).Error())

						continue
					}

					caseId, err := strconv.Atoi(tmp[0])
					if err != nil {
						module.logger.Send("warning", supportingfunctions.CustomError(err).Error())

						continue
					}

					eventId, err := strconv.Atoi(tmp[1])
					if err != nil {
						module.logger.Send("warning", supportingfunctions.CustomError(err).Error())

						continue
					}

					fmt.Printf("func 'ApiSqlite3Module.route' data.Command:'%s' updates information about case, caseId:%d, new eventId:%d\n", data.Command, caseId, eventId)

					if err = module.UpdateCaseId(ctx, caseId, eventId); err != nil {
						fmt.Println("func 'ApiSqlite3Module.route' ERROR:", err)

						module.logger.Send("warning", supportingfunctions.CustomError(err).Error())
					}

				case "delete case id":
					str := string(data.Payload)
					caseId, err := strconv.Atoi(str)
					if err != nil {
						module.logger.Send("error", supportingfunctions.CustomError(err).Error())

						continue
					}

					if err := module.DeleteCaseId(ctx, caseId); err != nil {
						module.logger.Send("error", supportingfunctions.CustomError(err).Error())
					}

				}
			}
		}
	}()
}
