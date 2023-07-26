package natsinteractions_test

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"placeholder_misp/confighandler"
	"placeholder_misp/coremodule"
	"placeholder_misp/datamodels"
	"placeholder_misp/natsinteractions"
	rules "placeholder_misp/rulesinteraction"
	"placeholder_misp/supportingfunctions"
	"placeholder_misp/tmpdata"
)

var _ = Describe("Natsinteraction", Ordered, func() {
	var (
		ctx      context.Context
		errConn  error
		closeCtx context.CancelFunc
		mnats    *natsinteractions.ModuleNATS
		chanLog  chan<- datamodels.MessageLoging
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
		chanLog = make(chan<- datamodels.MessageLoging)

		ctx, closeCtx = context.WithTimeout(context.Background(), 2*time.Second)

		mnats, errConn = natsinteractions.NewClientNATS(ctx, confighandler.AppConfigNATS{
			//Host: "nats.cloud.gcm",
			Host: "127.0.0.1",
			Port: 4222,
		}, chanLog)
	})

	Context("Тест 1.1. Проверка декодирования тестовых данных из файла 'binaryDataOne'", func() {
		var exampleByte []byte

		for _, v := range strings.Split(tmpdata.GetExampleDataOne(), " ") {
			i, err := strconv.Atoi(v)
			if err != nil {
				continue
			}

			exampleByte = append(exampleByte, uint8(i))
		}

		/*It("Должно нормально отрабатывать функция  GetWhitespace", func() {
			fmt.Printf("%s none Whitespace\n", datamodels.GetWhitespace(0))
			fmt.Printf("%s one Whitespace\n", datamodels.GetWhitespace(1))
			fmt.Printf("%s two Whitespace\n", datamodels.GetWhitespace(2))
			fmt.Printf("%s three Whitespace\n", datamodels.GetWhitespace(3))

			Expect(true).Should(BeTrue())
		})*/

		It("При анмаршалинге данных в ИЗВЕСТНЫЙ ТИП ошибки быть не должно", func() {
			mm := datamodels.MainMessage{}
			err := json.Unmarshal(exampleByte, &mm)

			//fmt.Println("---- ExampleDataOne ----")
			//fmt.Println(mm.ToStringBeautiful(0))
			//fmt.Println("------------------------")

			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("Тест 1.2. Проверка декодирования тестовых данных из файла 'binaryDataTwo'", func() {
		var exampleByte []byte

		for _, v := range strings.Split(tmpdata.GetExampleDataTwo(), " ") {
			i, err := strconv.Atoi(v)
			if err != nil {
				continue
			}

			exampleByte = append(exampleByte, uint8(i))
		}

		It("При анмаршалинге данных в ИЗВЕСТНЫЙ ТИП ошибки быть не должно", func() {
			mm := datamodels.MainMessage{}
			err := json.Unmarshal(exampleByte, &mm)

			//fmt.Println("---- ExampleDataTwo ----")
			//fmt.Println(mm.ToStringBeautiful(0))
			//fmt.Println("------------------------")

			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("Тест 1.3. Проверка декодирования тестовых данных из файла 'binaryDataThree'", func() {
		var exampleByte []byte

		for _, v := range strings.Split(tmpdata.GetExampleDataThree(), " ") {
			i, err := strconv.Atoi(v)
			if err != nil {
				continue
			}

			exampleByte = append(exampleByte, uint8(i))
		}

		It("При анмаршалинге данных в ИЗВЕСТНЫЙ ТИП ошибки быть не должно", func() {
			mm := datamodels.MainMessage{}
			err := json.Unmarshal(exampleByte, &mm)

			Expect(err).ShouldNot(HaveOccurred())
		})

		It("При анмаршалинге в НЕИЗВЕСТНЫЙ тип ошибки быть не должно", func() {

			/*
				Теперь для вывода данных из JSON сообщения используется функция
			*/

			_, err := supportingfunctions.NewReadReflectJSONSprint(exampleByte)

			/*
				!!!!!!
				В coremodule.decodeMessageReflect тестово реализованна возможность замены
				строковых данных в определенных полях, подробнее смотреть в функции readReflectAnyTypeSprint
				надо написать для других типов, таких как int, bool. И сделать замену на основе правил из
				пакета rules
				!!!!!!

				result := map[string]interface{}{}
				err := json.Unmarshal(exampleByte, &result)
				Expect(err).ShouldNot(HaveOccurred())

				strData := coremodule.ReadReflectMapSprint(result, rules.ListRulesProcessedMISPMessage{}, 0)
				Expect(err).ShouldNot(HaveOccurred())
			*/

			//fmt.Println("---- REFLECTION MAPPING ExampleDataThree ----")
			//fmt.Println(strData)
			//fmt.Println("---------------------------------------------")

			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("Тест 2. Проверка замены некоторых значений в json файле", func() {
		var eb []byte

		for _, v := range strings.Split(tmpdata.GetExampleDataThree(), " ") {
			i, err := strconv.Atoi(v)
			if err != nil {
				continue
			}

			eb = append(eb, uint8(i))
		}

		It("Должны быть заменины некоторые значения на основе правил из файла 'procmispmsg.yaml'", func() {
			//lrp, lvw, err := rules.GetRuleProcessedMISPMsg("rules", "procmispmsg.yaml")
			var (
				strData     string
				procMsgHive coremodule.ProcessMessageFromHive
			)

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

			ok, warningMsg := procMsgHive.HandleMessage()
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

/*
	Извесный тип
ttps:
  ttp:
    1.
      _createdAt: '1686652917436'
      _createdBy: 'p.delyukin@cloud.rcm'
      _id: '~248369200'
        _createdAt: '1679910385644'
        _createdBy: 'admin@thehive.local'
        _id: '~282632'
        _type: 'Pattern'
        dataSources:
          1. 'User Account: User Account Authentication'
          2. 'Command: Command Execution'
          3. 'Application Log: Application Log Content'
        defenseBypassed:
        description: 'Adversaries may use brute force techniques to gain access to accounts when passwords are unknown or when password hashes are obtained. Without knowledge of the password for an account or set of accounts, an adversary may systematically guess the password using a repetitive or iterative mechanism. Brute forcing passwords can take place via interaction with a service that will check the validity of those credentials or offline against previously acquired credential data, such as password hashes.

Brute forcing credentials may take place at various points during a breach. For example, adversaries may attempt to brute force access to [Valid Accounts](https://attack.mitre.org/techniques/T1078) within a victim environment leveraging knowledge gathered from other post-compromise behaviors such as [OS Credential Dumping](https://attack.mitre.org/techniques/T1003), [Account Discovery](https://attack.mitre.org/techniques/T1087), or [Password Policy Discovery](https://attack.mitre.org/techniques/T1201). Adversaries may also combine brute forcing activity with behaviors such as [External Remote Services](https://attack.mitre.org/techniques/T1133) as part of Initial Access.'
        extraData:
        name: 'Brute Force'
        patternId: 'T1110'
        patternType: 'attack-pattern'
        permissionsRequired:
        platforms:
          1. 'Windows'
          2. 'Azure AD'
          3. 'Office 365'
          4. 'SaaS'
          5. 'IaaS'
          6. 'Linux'
          7. 'macOS'
          8. 'Google Workspace'
          9. 'Containers'
          10. 'Network'
        remoteSupport: 'false'
        revoked: 'false'
        systemRequirements:
        tactics:
          1. 'credential-access'
        URL: 'https://attack.mitre.org/techniques/T1110'
        version: '2.4'
        _createdAt: '0'
        _createdBy: ''
        _id: ''
        _type: ''
        dataSources:
        defenseBypassed:
        description: ''
        extraData:
        name: ''
        patternId: ''
        patternType: ''
        permissionsRequired:
        platforms:
        remoteSupport: 'false'
        revoked: 'false'
        systemRequirements:
        tactics:
        URL: ''
        version: ''
      occurDate: '1686652860000'
      patternId: 'T1110'
      tactic: 'credential-access'

	REFLECTION MAPPING
	ttp:
    occurDate: 1686652860000
    patternId: 'T1110'
    tactic: 'credential-access'
    _createdAt: 1686652917436
    _createdBy: 'p.delyukin@cloud.rcm'
    _id: '~248369200'
    extraData:
      pattern:
        platforms:
          1. 'Windows'
          2. 'Azure AD'
          3. 'Office 365'
          4. 'SaaS'
          5. 'IaaS'
          6. 'Linux'
          7. 'macOS'
          8. 'Google Workspace'
          9. 'Containers'
          10. 'Network'
        remoteSupport: false
        url: 'https://attack.mitre.org/techniques/T1110'
        version: '2.4'
        capecId: 'CAPEC-49'
        _createdBy: 'admin@thehive.local'
        _id: '~282632'
        dataSources:
          1. 'User Account: User Account Authentication'
          2. 'Command: Command Execution'
          3. 'Application Log: Application Log Content'
        description: 'Adversaries may use brute force techniques to gain access to accounts when passwords are unknown or when password hashes are obtained. Without knowledge of the password for an account or set of accounts, an adversary may systematically guess the password using a repetitive or iterative mechanism. Brute forcing passwords can take place via interaction with a service that will check the validity of those credentials or offline against previously acquired credential data, such as password hashes.

Brute forcing credentials may take place at various points during a breach. For example, adversaries may attempt to brute force access to [Valid Accounts](https://attack.mitre.org/techniques/T1078) within a victim environment leveraging knowledge gathered from other post-compromise behaviors such as [OS Credential Dumping](https://attack.mitre.org/techniques/T1003), [Account Discovery](https://attack.mitre.org/techniques/T1087), or [Password Policy Discovery](https://attack.mitre.org/techniques/T1201). Adversaries may also combine brute forcing activity with behaviors such as [External Remote Services](https://attack.mitre.org/techniques/T1133) as part of Initial Access.'
        extraData:
        name: 'Brute Force'
        patternType: 'attack-pattern'
        _createdAt: 1679910385644
        revoked: false
        capecUrl: 'https://capec.mitre.org/data/definitions/49.html'
        defenseBypassed:
        detection: 'Monitor authentication logs for system and application login failures of [Valid Accounts](https://attack.mitre.org/techniques/T1078). If authentication failures are high, then there may be a brute force attempt to gain access to a system using legitimate credentials. Also monitor for many failed authentication attempts across various accounts that may result from password spraying attempts. It is difficult to detect when hashes are cracked, since this is generally done outside the scope of the target network.'
        patternId: 'T1110'
        systemRequirements:
        _type: 'Pattern'
        tactics:
          1. 'credential-access'
        permissionsRequired:
*/
