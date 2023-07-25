package rules_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	rules "placeholder_misp/rulesinteraction"
)

var _ = Describe("Readrule", func() {
	printRuleResult := func(r rules.ListRulesProcessingMsgMISP) string {
		resultPrint := fmt.Sprintln("RULES:")

		resultPrint += fmt.Sprintln("  REPLACE:")
		for k, v := range r.Rules.Replace {
			resultPrint += fmt.Sprintln("  ", k+1, ".")
			resultPrint += fmt.Sprintf("    searchField: '%s'\n", v.SearchField)
			resultPrint += fmt.Sprintf("    searchValue: '%s'\n", v.SearchValue)
			resultPrint += fmt.Sprintf("    replaceValue: '%s'\n", v.ReplaceValue)
		}

		resultPrint += fmt.Sprintln("  PASS:")
		for k, v := range r.Rules.Pass {
			resultPrint += fmt.Sprintln("  ", k+1, ".")
			resultPrint += fmt.Sprintf("    searchField: '%s'\n", v.SearchField)
			resultPrint += fmt.Sprintf("    searchValue: '%s'\n", v.SearchValue)
		}

		resultPrint += fmt.Sprintf("  PASSANY: '%v'\n", r.Rules.Passany)

		return resultPrint
	}

	printRuleIndex := func(r rules.ListRulesProcessingMsgMISP) string {
		resultPrint := fmt.Sprintln("RULE INDEX:")

		for k, v := range r.RulesIndex {
			resultPrint += fmt.Sprintf("  %s:\n", k)
			for num, value := range v {
				resultPrint += fmt.Sprintln("    ", num+1, ".")
				resultPrint += fmt.Sprintf("      ruleType: '%s'\n", value.RuleType)
				resultPrint += fmt.Sprintf("      searchField: '%s'\n", value.SearchField)
				resultPrint += fmt.Sprintf("      replaceValue: '%s'\n", value.ReplaceValue)
			}
		}

		return resultPrint
	}

	printRulePasstest := func(r rules.ListRulesProcessingMsgMISP) string {
		resultPasstest := fmt.Sprintln("RULE PASSTEST:")

		for key, value := range r.Rules.Passtest {
			resultPasstest += fmt.Sprintln("  ", key+1, ".")
			for k, v := range value.ListAnd {
				resultPasstest += fmt.Sprintln("    ", k+1, ".")
				resultPasstest += fmt.Sprintf("      searchField: '%s'\n", v.SearchField)
				resultPasstest += fmt.Sprintf("      searchValue: '%s'\n", v.SearchValue)
			}
		}

		return resultPasstest
	}

	printVerificationWarning := func(lvw []string) string {
		var resultPrint string

		for _, v := range lvw {
			resultPrint += fmt.Sprintln(v)
		}

		return resultPrint
	}

	Context("Тест 1. Чтение тестового файла procmispmsg_test_error.yaml с ошибочным построение правил", func() {
		It("При чтении файла с правилами ошибок быть не должно, файл должен быть успешно прочитан", func() {
			r, lw, err := rules.GetRuleProcessingMsgForMISP("rules", "procmispmsg_test_error.yaml")

			fmt.Println("-------- FILE 'procmispmsg_test_error.yaml' ----------")
			fmt.Println("NEW RULES FILE: ", r)
			fmt.Println("LIST WARNING: ", lw)

			fmt.Println("list verification warning:")
			fmt.Println(printVerificationWarning(lw))

			fmt.Println("1. _________ RULE procmispmsg_test_error.yaml.")
			fmt.Println("new rule result:")
			fmt.Println(printRuleResult(r))
			fmt.Println(printRuleIndex(r))

			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("Тест 2. Чтение нового тестового файла, в другом формате", func() {
		It("Новый тестовый файл должен быть успешно прочитан", func() {
			r, lw, err := rules.GetRuleProcessingMsgForMISP("rules", "procmispmsg_test.yaml")

			fmt.Println("NEW RULES FILE: ", r.Rules.Passtest)
			fmt.Println("LIST WARNING: ", lw)

			fmt.Println("2. _________ RULE procmispmsg_test.yaml.")
			fmt.Println("new rule result:")
			fmt.Println(printRuleResult(r))
			fmt.Println(printRuleIndex(r))
			fmt.Println(printRulePasstest(r))

			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
