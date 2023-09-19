package testaddnewelements_test

import (
	"bufio"
	"fmt"
	"os"
	"path"

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
		loging                         chan<- datamodels.MessageLoging
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
		loging = make(chan<- datamodels.MessageLoging)

		confApp.AppConfigMISP.Host = "misp-world.cloud.gcm"
		confApp.AppConfigMISP.Auth = "TvHkjH8jVQEIdvAxjxnL4H6wDoKyV7jobDjndvAo"

		//инициализируем модуль временного хранения информации
		storageApp = memorytemporarystorage.NewTemporaryStorage()

		// инициализируем модуль чтения правил обработки MISP сообщений
		listRules, _, errRules = rules.GetRuleProcessingMsgForMISP("rules", "mispmsgrule.yaml")

		//читаем тестовый файл
		exampleByte, errReadFile = readFileJson("natsinteractions/test_json", "example_3.json")

		//инициалиация модуля для взаимодействия с MISP
		mispModule, errMisp = mispinteractions.HandlerMISP(confApp.AppConfigMISP, storageApp, loging)

		//формирование итоговых документов в формате MISP
		chanCreateMispFormat, chanDone = coremodule.NewMispFormat(mispModule, loging)
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
			//обработчик сообщений из TheHive (выполняется разбор сообщения и его разбор на основе правил)
			coremodule.HandlerMessageFromHive(exampleByte, uuid.New().String(), storageApp, listRules, chanCreateMispFormat, chanDone, loging)

			Expect(true).Should(BeTrue())
		})
	})
})
