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
		listRules                      rules.ListRulesProcessingMsgMISP
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

		confApp.AppConfigMISP.Host = "misp-world.cloud.gcm"
		confApp.AppConfigMISP.Auth = "TvHkjH8jVQEIdvAxjxnL4H6wDoKyV7jobDjndvAo"

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

		//инициализируем модуль временного хранения информации
		storageApp = memorytemporarystorage.NewTemporaryStorage()

		// инициализируем модуль чтения правил обработки MISP сообщений
		listRules, _, errRules = rules.GetRuleProcessingMsgForMISP("rules", "mispmsgrule.yaml")

		//читаем тестовый файл
		//exampleByte, errReadFile = readFileJson("natsinteractions/test_json", "example_3.json")
		exampleByte, errReadFile = readFileJson("natsinteractions/test_json", "example_caseId_33705_3.json")

		/*
			по тестам не дублируется, а в боевом дублируется 3 раза
			однако и кейсов 33705 из хайва приходили 3 раза
			может это связанно

			а еще в других событиях появляется 3 объекта, таких же как
			и в кейсе 33705
		*/

		//инициалиация модуля для взаимодействия с MISP
		mispModule, errMisp = mispinteractions.HandlerMISP(confApp.AppConfigMISP, storageApp, logging)

		//формирование итоговых документов в формате MISP
		chanCreateMispFormat, chanDone = coremodule.NewMispFormat(mispModule, logging)

		//обработчик сообщений из TheHive (выполняется разбор сообщения и его разбор на основе правил)
		coremodule.HandlerMessageFromHive(exampleByte, uuid.New().String(), storageApp, listRules, chanCreateMispFormat, chanDone, logging, counting)
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
		It("", func() {

			mispOutput := mispModule.GetDataReceptionChannel()

			fmt.Println("TEST 2.1, mispOutput = ", <-mispOutput)
			fmt.Println("TEST 2.2, mispOutput = ", <-mispOutput)

			Expect(true).Should(BeTrue())
		})
	})
})
