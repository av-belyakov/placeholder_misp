package testaddnewelements_test

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/av-belyakov/placeholder_misp/cmd/coremodule"
	"github.com/av-belyakov/placeholder_misp/cmd/mispapi"
	"github.com/av-belyakov/placeholder_misp/commoninterfaces"
	"github.com/av-belyakov/placeholder_misp/constants"
	"github.com/av-belyakov/placeholder_misp/internal/confighandler"
	"github.com/av-belyakov/placeholder_misp/internal/countermessage"
	"github.com/av-belyakov/placeholder_misp/internal/logginghandler"
	rules "github.com/av-belyakov/placeholder_misp/internal/ruleshandler"
	"github.com/av-belyakov/placeholder_misp/internal/supportingfunctions"
	"github.com/av-belyakov/simplelogger"
)

var _ = Describe("Addneweventandattributes", Ordered, func() {
	var (
		logging                        *logginghandler.LoggingChan
		counting                       *countermessage.CounterMessage
		confApp                        confighandler.ConfigApp
		listRules                      *rules.ListRule
		mispModule                     *mispapi.ModuleMISP
		exampleByte                    []byte
		errReadFile, errMisp, errRules error
	)

	readFileJson := func(fpath, fname string) ([]byte, error) {
		var newResult []byte

		rootPath, err := supportingfunctions.GetRootPath("placeholder_misp")
		if err != nil {
			return newResult, err
		}

		fmt.Println("func 'readFileJson', path = ", path.Join(rootPath, fpath, fname))

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
		simpleLogger, err := simplelogger.NewSimpleLogger(context.Background(), constants.Root_Dir, []simplelogger.Options{})
		if err != nil {
			log.Fatalf("error module 'simplelogger': %v", err)
		}

		chZabbix := make(chan commoninterfaces.Messager)
		counting = countermessage.New(chZabbix)
		logging = logginghandler.New(simpleLogger, chZabbix)

		// NATS
		confApp.AppConfigNATS.Host = "nats.cloud.gcm"
		confApp.AppConfigNATS.Port = 4222

		//misp-world
		//confApp.AppConfigMISP.Host = "misp-world.cloud.gcm"
		//confApp.AppConfigMISP.Auth = "TvHkjH8jVQEIdvAxjxnL4H6wDoKyV7jobDjndvAo"

		//misp-center
		confApp.AppConfigMISP.Host = "misp-center.cloud.gcm"
		confApp.AppConfigMISP.Auth = "Z2PwRBdP5lFP7rdDJBzxmSahaLEwIvJoeOuwhRYQ"

		go func() {
			fmt.Println("___ Logging START")
			defer fmt.Println("___ Logging STOP")

			for log := range logging.GetChan() {
				fmt.Println("----", log, "----")
			}
		}()

		//вывод данных для Zabbix
		go func() {
			for msg := range chZabbix {
				fmt.Println("message for Zabbix:", msg)
			}
		}()

		taskId := uuid.New().String()

		//инициализация списка правил
		listRules, _, errRules = rules.NewListRule("placeholder_misp", "rules", "mispmsgrule.yaml")

		//читаем тестовый файл
		//"example_caseId_33705.json" совпадает с правилами
		//"example_caseId_33807.json" НЕ совпадает с правилами
		exampleByte, errReadFile = readFileJson("test/test_json", "examplenew.json")

		//инициалиация модуля для взаимодействия с MISP
		mispModule, errMisp = mispapi.NewModuleMISP(confApp.GetAppMISP().Host, confApp.GetAppMISP().Auth, confApp.GetListOrganization(), logging)

		hjson := coremodule.NewHandlerJSON(counting, logging)
		// обработчик JSON документа
		chanOutputDecodeJson := hjson.Start(exampleByte, taskId)
		//формирование итоговых документов в формате MISP
		go coremodule.CreateObjectsFormatMISP(chanOutputDecodeJson, taskId, mispModule, listRules, counting, logging)
	})

	Context("Тест 1. Проверка инициализации модулей", func() {
		It("При инициализации модуля чтения правил обработки не должно быть ошибки", func() {
			Expect(errRules).ShouldNot(HaveOccurred())
		})

		It("При инициализации модуля чтения файла примера не должно быть ошибки", func() {
			Expect(errReadFile).ShouldNot(HaveOccurred())
		})

		It("При инициализации модуля обработки MISP не должно быть ошибки", func() {
			Expect(errMisp).ShouldNot(HaveOccurred())
		})
	})

	Context("Тест 2. Проверяем обработчик кейсов", func() {
		It("Должны прийти два события от модуля misp", func() {
			mispOutput := mispModule.GetReceptionChannel()

			var err error

			for count := 0; count < 2; count++ {
				mop := <-mispOutput

				fmt.Println(count, ". TEST , mispOutput = ", mop)

				if mop.Command == "send event id" {
					fmt.Println("__#######$$$$$$$$$$$%%%%%%%%%%%%%%______________________________")
					fmt.Println("Sending event id", mop.EventId, " to NATS, taskId:", mop.TaskId)

					//natsSendData := func(conf confighandler.AppConfigNATS, eventId string)
					/*if err = natsSendData(confApp.AppConfigNATS, mop.EventId); err != nil {
						break
					}*/
				}
			}

			Expect(err).ShouldNot(HaveOccurred())
			Expect(true).Should(BeTrue())
		})
	})
})
