package redisapi_test

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/av-belyakov/placeholder_misp/cmd/redisapi"
	"github.com/av-belyakov/placeholder_misp/commoninterfaces"
	"github.com/av-belyakov/placeholder_misp/constants"
	"github.com/av-belyakov/placeholder_misp/internal/confighandler"
	"github.com/av-belyakov/placeholder_misp/internal/logginghandler"
	"github.com/av-belyakov/placeholder_misp/internal/supportingfunctions"
	"github.com/av-belyakov/simplelogger"
)

var _ = Describe("Interactionredisdb", Ordered, func() {
	var (
		errReadFile error
		exampleByte []byte
		module      *redisapi.ModuleRedis
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
		ca, _ := confighandler.New(constants.Root_Dir)

		simpleLogger, err := simplelogger.NewSimpleLogger(context.Background(), constants.Root_Dir, []simplelogger.Options{})
		if err != nil {
			log.Fatalf("error module 'simplelogger': %v", err)
		}

		chZabbix := make(chan commoninterfaces.Messager)

		// канал для логирования
		logging := logginghandler.New(simpleLogger, chZabbix)
		//logging.Start(ctx)
		exampleByte, errReadFile = readFileJson("testing/test_json", "example_caseId_33705_1.json")

		ctxRedis, _ := context.WithTimeout(context.Background(), 2*time.Second)
		module = redisapi.HandlerRedis(ctxRedis, ca.AppConfigRedis, logging)

		go func() {
			for log := range logging.GetChan() {
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
			module.SendingDataInput(redisapi.SettingsChanInputRedis{
				Command: "set caseId",
				//caseId:eventId
				Data: "12003:789",
			})

			module.SendingDataInput(redisapi.SettingsChanInputRedis{
				Command: "set caseId",
				//caseId:eventId
				Data: "14012:893",
			})

			Expect(true).Should(BeTrue())
		})

		It("Должно быть успешно найдено ранее добавленое значение", func(ctx SpecContext) {
			module.SendingDataInput(redisapi.SettingsChanInputRedis{
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
			module.SendingDataInput(redisapi.SettingsChanInputRedis{
				Command: "set caseId",
				//caseId:eventId
				Data: "12003:112340",
			})

			module.SendingDataInput(redisapi.SettingsChanInputRedis{
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
				module.SendingDataInput(redisapi.SettingsChanInputRedis{
					Command: "set raw case",
					//caseId:eventId
					//Data: "14012:893",
					RawData: exampleByte,
				})

			})
			It("Из БД должен быть успешно получен кейс в RAW формате", func() {
				module.SendingDataInput(redisapi.SettingsChanInputRedis{
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
