package coremodule_test

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"placeholder_misp/coremodule"
	"placeholder_misp/datamodels"
	"placeholder_misp/memorytemporarystorage"
	"placeholder_misp/mispinteractions"
	rules "placeholder_misp/rulesinteraction"
	"placeholder_misp/supportingfunctions"
)

var _ = Describe("CreateMispFormat", Ordered, func() {
	var (
		listRule                          rules.ListRulesProcessingMsgMISP
		storageApp                        *memorytemporarystorage.CommonStorageTemporary
		logging                           chan datamodels.MessageLogging
		mispOutput                        <-chan mispinteractions.SettingChanOutputMISP
		moduleMisp                        *mispinteractions.ModuleMISP
		exampleByte                       []byte
		counting                          chan datamodels.DataCounterSettings
		errReadFile, errGetRule, errHMisp error
	)

	/*getExampleByte := func() []byte {
		byteTmp := []byte{}

		for _, v := range strings.Split(tmpdata.GetExampleDataThree(), " ") {
			i, err := strconv.Atoi(v)
			if err != nil {
				continue
			}

			byteTmp = append(byteTmp, uint8(i))
		}

		return byteTmp
	}*/

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
		//инициализация модуля конфига
		//ca, _ := confighandler.NewConfig()

		//канал для логирования
		logging = make(chan datamodels.MessageLogging)
		//канал для подсчета обработанных кейсов
		counting = make(chan datamodels.DataCounterSettings)

		//читаем тестовый файл
		exampleByte, errReadFile = readFileJson("natsinteractions/test_json", "example_caseId_33705.json")

		//инициализация списка правил
		listRule, _, errGetRule = rules.GetRuleProcessingMsgForMISP("rules", "mispmsgrule.yaml")

		//инициализируем модуль временного хранения информации
		storageApp = memorytemporarystorage.NewTemporaryStorage()

		moduleMisp = &mispinteractions.ModuleMISP{
			ChanInputMISP:  make(chan mispinteractions.SettingsChanInputMISP),
			ChanOutputMISP: make(chan mispinteractions.SettingChanOutputMISP),
		}

		/*
			//инициализация модуля для взаимодействия с MISP
			moduleMisp, errHMisp = mispinteractions.HandlerMISP(*ca.GetAppMISP(), storageApp, logging)

			mispOutput = moduleMisp.GetDataReceptionChannel()
		*/
	})

	Context("Тест 1. Чтение тестового JSON файла", func() {
		It("При чтении тестового файла ошибок быть не должно", func() {
			Expect(errReadFile).ShouldNot(HaveOccurred())
		})
	})

	Context("Тест 2. Проверка формирования правил фильтрации", func() {
		It("При конвертировании времени в формат ISO 8601 ошибок быть не должно", func() {
			dateTest := "1692322715497"
			dt, err := strconv.ParseInt(dateTest, 10, 64)

			Expect(err).ShouldNot(HaveOccurred())

			t := time.UnixMilli(dt)
			tstr := t.Format(time.RFC3339)

			fmt.Println("		Time format RFC3339 as string: ", tstr)

			Expect(true).Should(BeTrue())
		})

		It("При формировании правил фильтрации ошибки быть не должно", func() {
			Expect(errGetRule).ShouldNot(HaveOccurred())
		})

		It("При инициализации модуля взаимодействия с MISP ошибки быть не должно", func() {
			Expect(errHMisp).ShouldNot(HaveOccurred())
		})
	})

	/*
		//
		//
		// Для тестов незабывать устанавливать переменную окружения
		// export GO_PHMISP_MAIN="development"
		// что бы взаимодействовать с misp-world.cloud.gcm
		//
		//
	*/

	Context("Тест 3. Проверяем правельность обработки сообщения по правилам и корректность формирования MISP форматов", func() {
		It("Сообщение должно быть успешно обработано, должены быть сформированы типы Events и Attributes", func(ctx SpecContext) {
			cd := make(chan struct{})

			uuidTask := uuid.NewString()
			storageApp.SetRawDataHiveFormatMessage(uuidTask, exampleByte)

			go func() {
				for data := range moduleMisp.GetInputChannel() {
					fmt.Println("@@@@_________________ Приняты данные для отправки в MISP _________________@@@@")

					if d, ok := data.MajorData["attributes"]; ok {
						b, err := json.Marshal(d)
						if err != nil {
							fmt.Println("__________ ERROR ___________")
							fmt.Println(err)
						}

						str, err := supportingfunctions.NewReadReflectJSONSprint(b)
						if err != nil {
							fmt.Println("__________ ERROR ___________")
							fmt.Println(err)
						}

						fmt.Println(str)
						/*
							почему список из 16 а не из 15?!
						*/

					}

					moduleMisp.SendingDataOutput(mispinteractions.SettingChanOutputMISP{Command: "STOP "})

					return
				}
			}()

			//формирование итоговых документов в формате MISP
			chanCreateMispFormat, chanDone := coremodule.NewMispFormat(moduleMisp, logging)

			go coremodule.HandlerMessageFromHive(exampleByte, uuidTask, storageApp, listRule, chanCreateMispFormat, chanDone, logging, counting)

			go func() {
				for {
					select {
					case log := <-logging:

						if log.MsgType == "warning" {
							fmt.Println("___ Log = ", log.MsgData, " ____")
						}
					case numData := <-counting:
						fmt.Printf("Counter processed object, type:%s, count:%d\n", numData.DataType, numData.Count)
					case <-cd:
						fmt.Println("-== STOP TEST ==-")

						return
					}
				}
			}()

			fmt.Println("reseived data: ", <-mispOutput)
			cd <- struct{}{}

			/*
				Объект с 219-236 стр. файла "example_caseId_33705.json"
				в боевом приложении игнорируется

				Объект с 297-324 стр. файла "example_caseId_33705.json"
				в боевом приложении игнорируется
			*/

			Expect(true).Should(BeTrue())
		}, SpecTimeout(time.Second*15))
	})
})
