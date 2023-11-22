package testaddnewelements_test

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"placeholder_misp/confighandler"
	"placeholder_misp/coremodule"
	"placeholder_misp/datamodels"
	"placeholder_misp/memorytemporarystorage"
	"placeholder_misp/mispinteractions"
	rules "placeholder_misp/rulesinteraction"
	"placeholder_misp/supportingfunctions"
)

var _ = Describe("Addneweventandattributes", Ordered, func() {
	var (
		logging                        chan datamodels.MessageLogging
		counting                       chan datamodels.DataCounterSettings
		confApp                        confighandler.ConfigApp
		listRules                      *rules.ListRule
		mispModule                     *mispinteractions.ModuleMISP
		storageApp                     *memorytemporarystorage.CommonStorageTemporary
		chanCreateMispFormat           chan coremodule.ChanInputCreateMispFormat
		chanDone                       chan bool
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
		logging = make(chan datamodels.MessageLogging)
		counting = make(chan datamodels.DataCounterSettings)

		// NATS
		confApp.AppConfigNATS.Host = "nats.cloud.gcm"
		confApp.AppConfigNATS.Port = 4222

		//misp-world
		confApp.AppConfigMISP.Host = "misp-world.cloud.gcm"
		confApp.AppConfigMISP.Auth = "TvHkjH8jVQEIdvAxjxnL4H6wDoKyV7jobDjndvAo"

		//misp-center
		//confApp.AppConfigMISP.Host = "misp-center.cloud.gcm"
		//confApp.AppConfigMISP.Auth = "Z2PwRBdP5lFP7rdDJBzxmSahaLEwIvJoeOuwhRYQ"

		go func() {
			fmt.Println("___ Logging START")
			defer fmt.Println("___ Logging STOP")

			for log := range logging {
				fmt.Println("----", log, "----")
			}
		}()

		//вывод данных счетчика
		go func() {
			dc := storageApp.GetDataCounter()
			d, h, m, s := supportingfunctions.GetDifference(dc.StartTime, time.Now())

			fmt.Printf("\tСОБЫТИЙ принятых/обработанных: %d/%d, соответствие/не соответствие правилам: %d/%d, время со старта приложения: дней %d, часов %d, минут %d, секунд %d\n", dc.AcceptedEvents, dc.ProcessedEvents, dc.EventsMeetRules, dc.EventsDoNotMeetRules, d, h, m, s)

			for d := range counting {
				switch d.DataType {
				case "update accepted events":
					storageApp.SetAcceptedEventsDataCounter(d.Count)
				case "update processed events":
					storageApp.SetProcessedEventsDataCounter(d.Count)
				case "update events meet rules":
					storageApp.SetEventsMeetRulesDataCounter(d.Count)
				case "events do not meet rules":
					storageApp.SetEventsDoNotMeetRulesDataCounter(d.Count)
				}

				dc := storageApp.GetDataCounter()
				d, h, m, s := supportingfunctions.GetDifference(dc.StartTime, time.Now())

				fmt.Printf("\tСОБЫТИЙ принятых/обработанных: %d/%d, соответствие/не соответствие правилам: %d/%d, время со старта приложения: дней %d, часов %d, минут %d, секунд %d\n", dc.AcceptedEvents, dc.ProcessedEvents, dc.EventsMeetRules, dc.EventsDoNotMeetRules, d, h, m, s)
			}
		}()

		taskId := uuid.New().String()

		//инициализируем модуль временного хранения информации
		storageApp = memorytemporarystorage.NewTemporaryStorage()

		//инициализация списка правил
		listRules, _, errRules = rules.NewListRule("placeholder_misp", "rules", "mispmsgrule.yaml")

		//читаем тестовый файл
		//"example_caseId_33705.json" совпадает с правилами
		//"example_caseId_33807.json" НЕ совпадает с правилами
		exampleByte, errReadFile = readFileJson("natsinteractions/test_json", "example_caseId_33705.json")

		//инициалиация модуля для взаимодействия с MISP
		mispModule, errMisp = mispinteractions.HandlerMISP(confApp.AppConfigMISP, confApp.Organizations, logging)

		//формирование итоговых документов в формате MISP
		chanCreateMispFormat, chanDone = coremodule.NewMispFormat(taskId, mispModule, logging)

		//обработчик сообщений из TheHive (выполняется разбор сообщения и его разбор на основе правил)
		hmfh := coremodule.NewHandlerMessageFromHive(storageApp, listRules, logging, counting)
		go hmfh.HandlerMessageFromHive(chanCreateMispFormat, exampleByte, taskId, chanDone)
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
			mispOutput := mispModule.GetDataReceptionChannel()

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
