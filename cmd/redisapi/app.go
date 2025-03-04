// Пакет redisapi содержит маршрутизатор для обработки запросов к СУБД Redis
package redisapi

import (
	"context"
	"fmt"
	"log"
	"strings"

	redis "github.com/redis/go-redis/v9"

	"github.com/av-belyakov/placeholder_misp/commoninterfaces"
	"github.com/av-belyakov/placeholder_misp/constants"
	"github.com/av-belyakov/placeholder_misp/internal/supportingfunctions"
)

// NewModuleRedis модуль для взаимодействия с Redis API
func NewModuleRedis(host string, port int, logger commoninterfaces.Logger) *ModuleRedis {
	return &ModuleRedis{
		client: redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%d", host, port),
		}),
		chInput:  make(chan SettingsInput),
		chOutput: make(chan SettingsOutput),
		logger:   logger,
		host:     host,
		port:     port,
	}
}

// Start запуск модуля
func (r *ModuleRedis) Start(ctx context.Context) error {
	defer func() {
		close(r.chInput)
		close(r.chOutput)
	}()

	if ctx.Err() != nil {
		return ctx.Err()
	}

	if err := r.client.Ping(ctx).Err(); err != nil {
		return supportingfunctions.CustomError(fmt.Errorf("module Redis API, %w", err))
	}

	log.Printf("%vconnect to Redis database with address %v%s:%d%v\n", constants.Ansi_Bright_Green, constants.Ansi_Dark_Gray, r.host, r.port, constants.Ansi_Reset)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case msg := <-r.chInput:
				switch msg.Command {
				case "search caseId":
					strCmd := r.client.Get(ctx, msg.Data)
					if strResult, err := strCmd.Result(); err == nil {
						r.SendDataOutput(SettingsOutput{
							CommandResult: "found caseId",
							//возвращает eventId MISP
							Result: strResult,
						})
					}

				case "set case id":
					// ***********************************
					// Это логирование только для теста!!!
					// ***********************************
					r.logger.Send("testing", fmt.Sprintf("TEST_INFO func 'HandlerRedis', обрабатываем добавление CaseID и EventId '%s' to REDIS DB", msg.Data))
					//
					//

					tmp := strings.Split(msg.Data, ":")
					if len(tmp) == 0 {
						r.logger.Send("warning", fmt.Sprintf("it is not possible to split a string '%s' to add case and event information to the Redis DB", msg.Data))

						continue
					}

					//получаем старое значение eventId по текущему caseId (если оно есть)
					strCmd := r.client.Get(ctx, tmp[0])
					eventId, err := strCmd.Result()
					if err == nil {
						// ***********************************
						// Это логирование только для теста!!!
						// ***********************************
						r.logger.Send("testing", fmt.Sprintf("TEST_INFO func 'HandlerRedis', НАЙДЕНО СТАРОЕ значение CaseID '%s' отправляем в ядро найденное событие с event id '%s'", tmp[0], eventId))
						//
						//

						//отправляем eventId для удаления события в MISP
						r.SendDataOutput(SettingsOutput{
							CommandResult: "found event id",
							Result:        eventId,
						})
					}

					//заменяем старое значение (если есть) или создаем новое
					//tmp[0] - caseId и tmp[1] - eventId
					if err := r.client.Set(ctx, tmp[0], tmp[1], 0).Err(); err != nil {
						r.logger.Send("error", supportingfunctions.CustomError(err).Error())

						continue
					}

					// ***********************************
					// Это логирование только для теста!!!
					// ***********************************
					r.logger.Send("testing", fmt.Sprintf("TEST_INFO func 'HandlerRedis', выполнили замену старого значения event id: %s новым значением event id: %s, для case id: %s", eventId, tmp[1], tmp[0]))
					//
					//
				}
			}
		}
	}()

	return nil
}
