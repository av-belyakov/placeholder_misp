// Пакет redisinteractions содержит маршрутизатор для обработки запросов к СУБД Redis
package redisinteractions

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"strings"

	redis "github.com/redis/go-redis/v9"

	"placeholder_misp/confighandler"
	"placeholder_misp/datamodels"
	"placeholder_misp/memorytemporarystorage"
)

func HandlerRedis(
	ctx context.Context,
	conf confighandler.AppConfigRedis,
	storageApp *memorytemporarystorage.CommonStorageTemporary,
	logging chan<- datamodels.MessageLogging) *ModuleRedis {

	mredis := ModuleRedis{
		chanInputRedis:  make(chan SettingsChanInputRedis),
		chanOutputRedis: make(chan SettingChanOutputRedis),
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", conf.Host, conf.Port),
	})

	log.Printf("Connect to Redis DB with address %s:%d", conf.Host, conf.Port)

	go func() {
		for data := range mredis.chanInputRedis {
			switch data.Command {
			case "search caseId":
				strCmd := rdb.Get(ctx, data.Data)
				if strResult, err := strCmd.Result(); err == nil {
					mredis.SendingDataOutput(SettingChanOutputRedis{
						CommandResult: "found caseId",
						//возвращает eventId MISP
						Result: strResult,
					})
				}

			case "set case id":
				// ***********************************
				// Это логирование только для теста!!!
				// ***********************************
				logging <- datamodels.MessageLogging{
					MsgData: fmt.Sprintf("TEST_INFO func 'HandlerRedis', обрабатываем добавление CaseID и EventId '%s' to REDIS DB", data.Data),
					MsgType: "testing",
				}
				//
				//

				tmp := strings.Split(data.Data, ":")
				if len(tmp) == 0 {
					_, f, l, _ := runtime.Caller(0)

					logging <- datamodels.MessageLogging{
						MsgData: fmt.Sprintf("'it is not possible to split a string '%s' to add case and event information to the Redis DB' %s:%d", data.Data, f, l-1),
						MsgType: "warning",
					}

					continue
				}

				//получаем старое значение eventId по текущему caseId (если оно есть)
				strCmd := rdb.Get(ctx, tmp[0])
				eventId, err := strCmd.Result()
				if err == nil {
					// ***********************************
					// Это логирование только для теста!!!
					// ***********************************
					logging <- datamodels.MessageLogging{
						MsgData: fmt.Sprintf("TEST_INFO func 'HandlerRedis', НАЙДЕНО СТАРОЕ значение CaseID '%s' отправляем в ядро найденное событие с event id '%s'", tmp[0], eventId),
						MsgType: "testing",
					}
					//
					//

					//отправляем eventId для удаления события в MISP
					mredis.SendingDataOutput(SettingChanOutputRedis{
						CommandResult: "found event id",
						Result:        eventId,
					})
				}

				//заменяем старое значение (если есть) или создаем новое
				//tmp[0] - caseId и tmp[1] - eventId
				if err := rdb.Set(ctx, tmp[0], tmp[1], 0).Err(); err != nil {
					_, f, l, _ := runtime.Caller(0)

					logging <- datamodels.MessageLogging{
						MsgData: fmt.Sprintf("'%s' %s:%d", fmt.Sprint(err), f, l-1),
						MsgType: "error",
					}

					continue
				}

				// ***********************************
				// Это логирование только для теста!!!
				// ***********************************
				logging <- datamodels.MessageLogging{
					MsgData: fmt.Sprintf("TEST_INFO func 'HandlerRedis', выполнили замену старого значения event id: %s новым значением event id: %s, для case id: %s", eventId, tmp[1], tmp[0]),
					MsgType: "testing",
				}
				//
				//

			case "set raw case":
				/*

					Тут нужно добавить RAW данные кейса из TheHive в List



				*/

			case "get next raw case":
				/*

					запрос на получение из БД следующего кейса в формате RAW
					ответ должен содержать CommandResult = "sending next raw case"

				*/
			}
		}
	}()

	return &mredis
}
