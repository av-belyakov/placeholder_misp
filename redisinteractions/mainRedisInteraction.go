package redisinteractions

import (
	"context"
	"fmt"
	"log"
	"placeholder_misp/confighandler"
	"placeholder_misp/datamodels"
	"placeholder_misp/memorytemporarystorage"
	"runtime"
	"strings"

	redis "github.com/redis/go-redis/v9"
)

var mredis ModuleRedis

func init() {
	mredis = ModuleRedis{
		chanInputRedis:  make(chan SettingsChanInputRedis),
		chanOutputRedis: make(chan SettingChanOutputRedis),
	}
}

func HandlerRedis(
	ctx context.Context,
	conf confighandler.AppConfigRedis,
	storageApp *memorytemporarystorage.CommonStorageTemporary,
	loging chan<- datamodels.MessageLoging) *ModuleRedis {
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
				loging <- datamodels.MessageLoging{
					MsgData: fmt.Sprintf("TEST_INFO func 'HandlerRedis', обрабатываем добавление CaseID и EventId '%s' to REDIS DB", data.Data),
					MsgType: "info",
				}
				//
				//

				tmp := strings.Split(data.Data, ":")
				if len(tmp) == 0 {
					_, f, l, _ := runtime.Caller(0)

					loging <- datamodels.MessageLoging{
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
					loging <- datamodels.MessageLoging{
						MsgData: fmt.Sprintf("TEST_INFO func 'HandlerRedis', НАЙДЕНО СТАРОЕ значение CaseID '%s' отправляем в ядро найденное событие с event id '%s'", tmp[0], eventId),
						MsgType: "info",
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

					loging <- datamodels.MessageLoging{
						MsgData: fmt.Sprintf("'%s' %s:%d", fmt.Sprint(err), f, l-1),
						MsgType: "error",
					}

					continue
				}

				// ***********************************
				// Это логирование только для теста!!!
				// ***********************************
				loging <- datamodels.MessageLoging{
					MsgData: fmt.Sprintf("TEST_INFO func 'HandlerRedis', выполнили замену старого значения event id: %s новым значением event id: %s, для case id: %s", eventId, tmp[1], tmp[0]),
					MsgType: "info",
				}
				//
				//
			}
		}
	}()

	return &mredis
}
