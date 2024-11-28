package redisinteractions_test

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"placeholder_misp/confighandler"
	"placeholder_misp/datamodels"
	"placeholder_misp/memorytemporarystorage"
	"placeholder_misp/redisinteractions"
	"placeholder_misp/supportingfunctions"
)

var _ = Describe("Interactionredisdb", Ordered, func() {
	var (
		errReadFile error
		exampleByte []byte
		module      *redisinteractions.ModuleRedis
	)

	readFileJson := func(fpath, fname string) ([]byte, error) {
		var newResult []byte

		rootPath, err := supportingfunctions.GetRootPath("placeholder_misp")
		if err != nil {
			return newResult, err
		}

		//fmt.Println("func 'readFileJson', path = ", path.Join(rootPath, fpath, fname))

		f, err := os.OpenFile(path.Join(rootPath, fpath, fname), os.O_RDONLY, os.ModePerm)
		if err != nil {
			return newResult, err
		}
		defer f.Close()

		sc := bufio.NewScanner(f)
		for sc.Scan() {
			newResult = append(newResult, sc.Bytes()...)
		}

		return newResult, nil
	}

	BeforeAll(func() {
		// инициализация модуля конфига
		ca, _ := confighandler.NewConfig("placeholder_misp")

		// канал для логирования
		logging := make(chan datamodels.MessageLogging)

		exampleByte, errReadFile = readFileJson("testing/test_json", "example_caseId_33705_1.json")

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

	Context("Тест 0. Чтение тестового файла с кейсом", func() {
		It("При чтении файла ошибок быть не должно", func() {
			Expect(errReadFile).ShouldNot(HaveOccurred())
		})
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

		Context("Тест 2. Тестируем возможность добавления кейса (формат RAW) из TheHive в List", func() {
			It("При добавлении RAW формата не должно быть ошибок", func() {
				module.SendingDataInput(redisinteractions.SettingsChanInputRedis{
					Command: "set raw case",
					//caseId:eventId
					//Data: "14012:893",
					RawData: exampleByte,
				})

			})
			It("Из БД должен быть успешно получен кейс в RAW формате", func() {
				module.SendingDataInput(redisinteractions.SettingsChanInputRedis{
					Command: "get next raw case",
				})

				chanOutput := module.GetDataReceptionChannel()
				data := <-chanOutput

				fmt.Println("RECEIVED DATA = ", data.Result)

				Expect(data.CommandResult).Should(Equal("sending next raw case"))
			})
		})
	})
})
