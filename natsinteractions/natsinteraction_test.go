package natsinteractions_test

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"placeholder_misp/confighandler"
	"placeholder_misp/coremodule"
	"placeholder_misp/datamodels"
	"placeholder_misp/memorytemporarystorage"
	"placeholder_misp/natsinteractions"
	rules "placeholder_misp/rulesinteraction"
	"placeholder_misp/supportingfunctions"
	"placeholder_misp/tmpdata"
)

var _ = Describe("Natsinteraction", Ordered, func() {
	var (
		errConn    error
		closeCtx   context.CancelFunc
		mnats      *natsinteractions.ModuleNATS
		chanLog    chan<- datamodels.MessageLogging
		storageApp *memorytemporarystorage.CommonStorageTemporary
	)

	/*
		Для отправки логов в zabbix см. https://habr.com/ru/companies/nixys/news/503104/
	*/

	printVerificationWarning := func(lvw []string) string {
		var resultPrint string

		for _, v := range lvw {
			resultPrint += fmt.Sprintln(v)
		}

		return resultPrint
	}

	BeforeAll(func() {
		chanLog = make(chan<- datamodels.MessageLogging)
		//инициализируем модуль временного хранения информации
		storageApp = memorytemporarystorage.NewTemporaryStorage()

		mnats, errConn = natsinteractions.NewClientNATS(confighandler.AppConfigNATS{
			//Host: "nats.cloud.gcm",
			Host: "127.0.0.1",
			Port: 4222,
		},
			storageApp,
			chanLog)
	})

	Context("Тест 1. Проверка замены некоторых значений в json файле", func() {
		var eb []byte

		for _, v := range strings.Split(tmpdata.GetExampleDataThree(), " ") {
			i, err := strconv.Atoi(v)
			if err != nil {
				continue
			}

			eb = append(eb, uint8(i))
		}

		It("Должны быть заменены некоторые значения на основе правил из файла 'procmispmsg.yaml'", func() {
			//lrp, lvw, err := rules.GetRuleProcessedMISPMsg("rules", "procmispmsg.yaml")
			var (
				strData     string
				procMsgHive coremodule.FieldsNameMapping
			)

			chanOutMispFormat := make(chan coremodule.ChanInputCreateMispFormat)

			lr, lw, err := rules.GetRuleProcessingMsgForMISP("rules", "procmispmsg_test.yaml")

			fmt.Println("list verification warning:")
			fmt.Println(printVerificationWarning(lw))

			Expect(err).ShouldNot(HaveOccurred())

			reg := regexp.MustCompile(`_createdBy: \'[a-zA-Z_.@0-9]+\'`)
			regCaseId := regexp.MustCompile(`caseId: [0-9]+`)
			regRevoked := regexp.MustCompile(`revoked: [a-zA-Z]+`)
			regStartDate := regexp.MustCompile(`startDate: [a-zA-Z_.@0-9]+`)
			regPatternId := regexp.MustCompile(`capecId: \'[a-zA-Z0-9\-]+\'`)

			fmt.Println("____BEFORE reflect modify____")
			strData, err = supportingfunctions.NewReadReflectJSONSprint(eb)
			bl := reg.FindAllString(strData, 10)
			for k, v := range bl {
				fmt.Printf("%d. %s\n", k+1, v)
			}
			cibl := regCaseId.FindAllString(strData, 10)
			for k, v := range cibl {
				fmt.Printf("%d. %s\n", k+1, v)
			}
			rbl := regRevoked.FindAllString(strData, 10)
			for k, v := range rbl {
				fmt.Printf("%d. %s\n", k+1, v)
			}
			pidbl := regPatternId.FindAllString(strData, 10)
			for k, v := range pidbl {
				fmt.Printf("%d. %s\n", k+1, v)
			}

			Expect(err).ShouldNot(HaveOccurred())

			procMsgHive, err = coremodule.NewHandleMessageFromHive(eb, lr)

			Expect(err).ShouldNot(HaveOccurred())

			go func() {
				for v := range chanOutMispFormat {
					//что бы не выполнялось
					if v.Value == "111" {
						fmt.Printf("\n RESEIVED MESSAGE:\n - FieldName: %s\n - ValueType: %s\n - Value: %v\n - FieldBranch: %s\n", v.FieldName, v.ValueType, v.Value, v.FieldBranch)
					}
				}
			}()

			ok, warningMsg := procMsgHive.HandleMessage(chanOutMispFormat)
			neb, err := procMsgHive.GetMessage()

			Expect(err).ShouldNot(HaveOccurred())

			fmt.Println("____AFTER reflect modify____:")
			strData, err = supportingfunctions.NewReadReflectJSONSprint(neb)
			bl = reg.FindAllString(strData, 10)
			for k, v := range bl {
				fmt.Printf("%d. %s\n", k+1, v)
			}
			cibl = regCaseId.FindAllString(strData, 10)
			for k, v := range cibl {
				fmt.Printf("%d. %s\n", k+1, v)
			}
			rbl = regRevoked.FindAllString(strData, 10)
			for k, v := range rbl {
				fmt.Printf("%d. %s\n", k+1, v)
			}
			pidbl = regPatternId.FindAllString(strData, 10)
			for k, v := range pidbl {
				fmt.Printf("%d. %s\n", k+1, v)
			}
			sdbl := regStartDate.FindAllString(strData, 10)
			for k, v := range sdbl {
				fmt.Printf("%d. %s\n", k+1, v)
			}

			fmt.Println("procMsgHive.ProcessMessage() is true: ", ok)
			fmt.Println("warningMsg: ", warningMsg)
			fmt.Println("")

			Expect(err).ShouldNot(HaveOccurred())
		})

		It("Должен быть получен полный текстовый вывод объекта из хайва", func() {
			str, err := supportingfunctions.NewReadReflectJSONSprint(eb)

			//fmt.Println(str)

			Expect(str).ShouldNot(Equal(""))
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("Тест 3. Проверка инициализации соединения с NATS", func() {
		It("При инициализации соединения с NATS не должно быть ошибки", func() {
			Expect(errConn).ShouldNot(HaveOccurred())

			fmt.Println("Resevid message = ", <-mnats.GetDataReceptionChannel())
		})
	})

	AfterAll(func() {
		closeCtx()
	})
})
