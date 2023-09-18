package redisinteractions

import (
	"context"
	"fmt"
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

			case "set caseId":
				tmp := strings.Split(data.Data, ":")
				if len(tmp) == 0 {
					_, f, l, _ := runtime.Caller(0)

					loging <- datamodels.MessageLoging{
						MsgData: fmt.Sprintf(" 'it is not possible to split a string '%s' to add case and event information to the Redis DB' %s:%d", data.Data, f, l-2),
						MsgType: "warning",
					}

					continue
				}

				//получаем старое значение eventId по текущему caseId (если оно есть)
				strCmd := rdb.Get(ctx, tmp[0])
				if eventId, err := strCmd.Result(); err == nil {
					//отправляем eventId для удаления события в MISP
					mredis.SendingDataOutput(SettingChanOutputRedis{
						CommandResult: "found eventId",
						Result:        eventId,
					})
				}

				//заменяем старое значение (если есть) или создаем новое
				//tmp[0] - caseId и tmp[1] - eventId
				if err := rdb.Set(ctx, tmp[0], tmp[1], 0).Err(); err != nil {
					_, f, l, _ := runtime.Caller(0)

					loging <- datamodels.MessageLoging{
						MsgData: fmt.Sprintf(" '%s' %s:%d", fmt.Sprint(err), f, l-2),
						MsgType: "error",
					}
				}
			}
		}
	}()

	return &mredis
}
