package rules_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	rules "placeholder_misp/rulesinteraction"
)

var _ = Describe("Readrule", func() {
	printRuleResult := func(lr []rules.RuleProcMISPMessageField) string {
		var resultPrint string

		for k, v := range lr {
			resultPrint += fmt.Sprintln(k, ".")
			resultPrint += fmt.Sprintf("  actionType: %s\n", v.ActionType)

			for key, value := range v.ListRequiredValues {
				resultPrint += fmt.Sprintf("   %d.\n", key)
				resultPrint += fmt.Sprintf("    fieldName: %s.\n", value.FieldName)
				resultPrint += fmt.Sprintf("    typeValue: %s.\n", value.TypeValue)
				resultPrint += fmt.Sprintf("    value: %s.\n", value.Value)
				resultPrint += fmt.Sprintf("    replaceValue: %s.\n", value.ReplaceValue)
			}
		}

		return resultPrint
	}

	printVerificationWarning := func(lvw []string) string {
		var resultPrint string

		for _, v := range lvw {
			resultPrint += fmt.Sprintln(v)
		}

		return resultPrint
	}

	Context("Тест 1. Чтение основного файла с правилами", func() {
		It("При чтении файла с правилами ошибок быть не должно, файл должен быть успешно прочитан", func() {
			lrp, lvw, err := rules.GetRuleProcessedMISPMsg("rules", "procmispmsg.yaml")

			fmt.Println("1. _________ RULE procmispmsg.yaml.")
			fmt.Println("new result:")
			fmt.Println(printRuleResult(lrp.Rulles))

			fmt.Println("list verification warning:")
			fmt.Println(printVerificationWarning(lvw))

			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("Тест 2. Чтение тестового файла procmispmsg_test_1.yaml с правилами", func() {
		It("При чтении файла с правилами ошибок быть не должно, файл должен быть успешно прочитан", func() {
			lrp, lvw, err := rules.GetRuleProcessedMISPMsg("rules", "procmispmsg_test_1.yaml")

			fmt.Println("2. _________ RULE procmispmsg_test_1.yaml.")
			fmt.Println("new result:")
			fmt.Println(printRuleResult(lrp.Rulles))

			fmt.Println("list verification warning:")
			fmt.Println(printVerificationWarning(lvw))

			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("Тест 3. Чтение тестового файла procmispmsg_test_2.yaml с правилами", func() {
		It("При чтении файла с правилами ошибок быть не должно, файл должен быть успешно прочитан", func() {
			lrp, lvw, err := rules.GetRuleProcessedMISPMsg("rules", "procmispmsg_test_2.yaml")

			fmt.Println("3. _________ RULE procmispmsg_test_2.yaml.")
			fmt.Println("new result:")
			fmt.Println(printRuleResult(lrp.Rulles))

			fmt.Println("list verification warning:")
			fmt.Println(printVerificationWarning(lvw))

			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("Тест 4. Чтение тестового файла procmispmsg_test_error.yaml с ошибочным построение правил", func() {
		It("При чтении файла с правилами ошибок быть не должно, файл должен быть успешно прочитан", func() {
			lrp, lvw, err := rules.GetRuleProcessedMISPMsg("rules", "procmispmsg_test_error.yaml")

			fmt.Println("4. _________ RULE procmispmsg_test_error.yaml.")
			fmt.Println("new result:")
			fmt.Println(printRuleResult(lrp.Rulles))

			fmt.Println("list verification warning:")
			fmt.Println(printVerificationWarning(lvw))

			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
