package redisinteractions_test

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"placeholder_misp/confighandler"
	"placeholder_misp/datamodels"
	"placeholder_misp/memorytemporarystorage"
	"placeholder_misp/redisinteractions"
)

var _ = Describe("Interactionredisdb", Ordered, func() {
	var module *redisinteractions.ModuleRedis

	BeforeAll(func() {
		// инициализация модуля конфига
		ca, _ := confighandler.NewConfig()

		// канал для логирования
		logging := make(chan datamodels.MessageLogging)

		// инициализируем модуль временного хранения информации
		storageApp := memorytemporarystorage.NewTemporaryStorage()

		ctxRedis, _ := context.WithTimeout(context.Background(), 2*time.Second)
		module = redisinteractions.HandlerRedis(ctxRedis, ca.AppConfigRedis, storageApp, logging)

		go func() {
			for log := range logging {
				fmt.Println("LOGGING: ", log)
			}
		}()
	})

	Context("Тест 1. Тестируем взаимодействие с СУБД Redis", func() {
		It("Должно быть успешно добавлено новое значение в БД", func() {
			module.SendingDataInput(redisinteractions.SettingsChanInputRedis{
				Command: "set caseId",
				//caseId:eventId
				Data: "12003:789",
			})

			module.SendingDataInput(redisinteractions.SettingsChanInputRedis{
				Command: "set caseId",
				//caseId:eventId
				Data: "14012:893",
			})

			Expect(true).Should(BeTrue())
		})

		It("Должно быть успешно найдено ранее добавленое значение", func(ctx SpecContext) {
			module.SendingDataInput(redisinteractions.SettingsChanInputRedis{
				Command: "search caseId",
				Data:    "12003",
			})

			ch := module.GetDataReceptionChannel()

			info := <-ch

			fmt.Println("RESULT: ", info)

			Expect(info.CommandResult).Should(Equal("found caseId"))

			strRes, ok := info.Result.(string)

			Expect(ok).Should(BeTrue())
			Expect(strRes).Should(Equal("789"))

			//заменяю значение
			module.SendingDataInput(redisinteractions.SettingsChanInputRedis{
				Command: "set caseId",
				//caseId:eventId
				Data: "12003:112340",
			})

			module.SendingDataInput(redisinteractions.SettingsChanInputRedis{
				Command: "search caseId",
				Data:    "12003",
			})

			ch = module.GetDataReceptionChannel()

			info = <-ch

			fmt.Println("RESULT: ", info)

			strRes, ok = info.Result.(string)

			Expect(ok).Should(BeTrue())
			Expect(strRes).Should(Equal("112340"))
		}, SpecTimeout(time.Second*5))
	})
})
