package testprocessingrules_test

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"placeholder_misp/coremodule"
	"placeholder_misp/datamodels"
	"placeholder_misp/memorytemporarystorage"
	"placeholder_misp/mispinteractions"
	rules "placeholder_misp/rulesinteraction"
	"placeholder_misp/supportingfunctions"
	"sync"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Processingrules", Ordered, func() {
	var (
		lr *rules.ListRule
		//storageApp                        *memorytemporarystorage.CommonStorageTemporary
		logging chan datamodels.MessageLogging
		//mispOutput                        <-chan mispinteractions.SettingChanOutputMISP
		moduleMisp  *mispinteractions.ModuleMISP
		exampleByte []byte
		counting    chan datamodels.DataCounterSettings
		errReadFile, errGetRule/*, errHMisp*/ error
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
		//канал для логирования
		logging = make(chan datamodels.MessageLogging)
		//канал для подсчета обработанных кейсов
		counting = make(chan datamodels.DataCounterSettings)

		//читаем тестовый файл
		exampleByte, errReadFile = readFileJson("testing/test_json", "example_caseId_33705.json")

		//инициализация списка правил
		lr, _, errGetRule = rules.NewListRule("placeholder_misp", "rules", "mispmsgrule.yaml")

		//эмулируем результат инициализации модуля MISP
		moduleMisp = &mispinteractions.ModuleMISP{
			ChanInputMISP:  make(chan mispinteractions.SettingsChanInputMISP),
			ChanOutputMISP: make(chan mispinteractions.SettingChanOutputMISP),
		}
	})

	BeforeEach(func() {
		//выполняет очистку значения StatementExpression что равно отсутствию совпадений в правилах Pass
		lr.CleanStatementExpressionRulePass()
	})

	Context("Тест 0. Проверка функции PassRuleHandler", func() {
		It("Должны быть успешно найдены все элементы из правила", func() {
			list := map[string]interface{}{
				"event.object.resolutionStatus": "TruePositive",
				"event.object.impactStatus":     "WithImpact",
				"event.object.tlp":              "not:3",
			}

			for fieldBranch, v := range list {
				lr.PassRuleHandler(fieldBranch, v)
				//coremodule.PassRuleHandler(listRule.Rules.Pass, fieldBranch, v)
			}

			//fmt.Println("----------------------- Equal rules --------------------------")
			var count int
			for _, v := range lr.GetRulePass() {
				for _, value := range v.ListAnd {
					//		fmt.Printf("field '%s' is exist '%v'\n", value.SearchField, value.StatementExpression)

					if value.StatementExpression {
						count++
					}
				}
			}
			//fmt.Println("--------------------------------------------------------------")

			Expect(count).Should(Equal(3))
		})
	})

	Context("Тест 1. Чтение тестового JSON файла", func() {
		It("При чтении тестового файла ошибок быть не должно", func() {
			//fmt.Println("count example byte", len(exampleByte))

			Expect(errReadFile).ShouldNot(HaveOccurred())
		})
	})

	Context("Тест 2.1. Проверка формирования правил фильтрации", func() {
		It("При формировании правил фильтрации ошибки быть не должно", func() {
			//fmt.Println("Rules Pass:", lr.GetRulePass())

			Expect(errGetRule).ShouldNot(HaveOccurred())
		})

		//It("При инициализации модуля взаимодействия с MISP ошибки быть не должно", func() {
		//	Expect(errHMisp).ShouldNot(HaveOccurred())
		//})
	})

	Context("Тест 2.2. Проверка формирования правил фильтрации на основе НОВОГО конструктора", func() {
		It("Не должно быть ошибок при формировании правил фильтрации", func() {
			_, warnings, err := rules.NewListRule("placeholder_misp", "rules", "mispmsgrule.yaml")

			//fmt.Println()
			//fmt.Println("Rules warnings 1111 START")
			//for k, v := range warnings {
			//	fmt.Printf("%d. %s\n", k, v)
			//}
			//fmt.Println("Rules warnings 1111 END")

			//fmt.Println("-------- LIST RULE ---------", r, "------------------")

			Expect(len(warnings)).Should(Equal(0))
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("Должны быть ошибки при формировании правил на основе невалидного файла", func() {
			_, warnings, err := rules.NewListRule("placeholder_misp", "rules", "procmispmsg_test_error.yaml")

			//fmt.Println()
			//fmt.Println("Rules warnings 2222 START")
			//for k, v := range warnings {
			//	fmt.Printf("%d. %s\n", k, v)
			//}
			//fmt.Println("Rules warnings 2222 END")

			Expect(len(warnings)).ShouldNot(Equal(0))

			//нет ошибок так как они возможны только при чтении или парсинге файла
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("Тест 3. Формирование итоговых документов и проверка их соответствия правилам", func() {
		It("При выполнении формирования документов, соответствующих формату MISP ошибок быть не должно", func(ctx SpecContext) {
			done := make(chan struct{})

			var wg sync.WaitGroup
			wg.Add(1)

			//вывод логов и счетчиков
			go func() {
				fmt.Println("function SHOW Logs and Counts is START")

				for {
					select {
					case log := <-logging:
						if log.MsgType == "warning" {
							fmt.Println("___ Log = ", log.MsgData, " ____")
						}

						if log.MsgType == "STOP TEST" {
							fmt.Println("-== STOP function SHOW Logs and Counts ==-")

							done <- struct{}{}

							return
						}
					case numData := <-counting:
						fmt.Printf("Counter processed object, type:%s, count:%d\n", numData.DataType, numData.Count)
					}
				}
			}()

			go func() {
				fmt.Println("function SHOW Major Data is START")

				for {
					select {
					case data := <-moduleMisp.GetInputChannel():
						fmt.Println("@@@@_________________ Приняты данные для отправки в MISP _________________@@@@")

						if d, ok := data.MajorData["events"]; ok {
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

							fmt.Println("________________ events _________________")
							fmt.Println(str)
						}

						if d, ok := data.MajorData["attributes"]; ok {
							if l, ok := d.([]datamodels.AttributesMispFormat); ok {
								fmt.Println("___________ MISP type: ATTRIBUTES ___________")
								fmt.Println("+++++++++++++++ Length: ", len(l))
								fmt.Println("_____________________________________________")
							}
						}
						/*if d, ok := data.MajorData["attributes"]; ok {
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

							fmt.Println("________________ attributes _________________")
							fmt.Println(str)
						}*/

						/*if d, ok := data.MajorData["objects"]; ok {
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

							fmt.Println("________________ objects _________________")
							fmt.Println(str)
						}*/
					case <-done:
						fmt.Println("-== STOP function SHOW Major Data ==-")

						wg.Done()

						return
					}
				}
			}()

			// инициализируем модуль временного хранения информации
			storageApp := memorytemporarystorage.NewTemporaryStorage()

			hmfh := coremodule.NewHandlerMessageFromHive(storageApp, lr, logging, counting)

			msgId := uuid.New().String()

			// формирование итоговых документов в формате MISP
			chanCreateMispFormat, chanDone := coremodule.NewMispFormat(msgId, moduleMisp, logging)
			// обработчик сообщений из TheHive (выполняется разбор сообщения и его разбор на основе правил)
			go hmfh.HandlerMessageFromHive(chanCreateMispFormat, exampleByte, msgId, chanDone)

			wg.Wait()

			// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
			//
			//тест может нормально не проходить потому что в hmfh.HandlerMessageFromHive
			//есть еще один метод lr.CleanStatementExpressionRulePass()
			//который выполняет очистку значения StatementExpression
			//по этому перед тестом его рекомендуется отключить
			//
			// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

			Expect(lr.SomePassRuleIsTrue()).Should(BeTrue())
			Expect(true).Should(BeTrue())
		}, SpecTimeout(time.Second*15))
	})

	/*Context("", func(){
		It("", func ()  {

		})
	})*/
})
