package rules_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	rules "placeholder_misp/rulesinteraction"
)

var _ = Describe("Readrule", Ordered, func() {
	const fileName = "mispmsgrule.yaml"
	var (
		lr         *rules.ListRule
		lw         []string
		errGetRule error
	)

	printRuleResult := func(r *rules.ListRule) string {
		resultPrint := fmt.Sprintln("RULES:")

		resultPrint += fmt.Sprintln("  REPLACE:")
		for k, v := range r.Rules.Replace {
			resultPrint += fmt.Sprintln("  ", k+1, ".")
			resultPrint += fmt.Sprintf("    searchField: '%s'\n", v.SearchField)
			resultPrint += fmt.Sprintf("    searchValue: '%s'\n", v.SearchValue)
			resultPrint += fmt.Sprintf("    replaceValue: '%s'\n", v.ReplaceValue)
		}

		resultPrint += fmt.Sprintln("  PASS:")
		for key, value := range r.Rules.Pass {
			resultPrint += fmt.Sprintln("  ", key+1, ".")
			for k, v := range value.ListAnd {
				resultPrint += fmt.Sprintln("    ", k+1, ".")
				resultPrint += fmt.Sprintf("      searchField: '%s'\n", v.SearchField)
				resultPrint += fmt.Sprintf("      searchValue: '%s'\n", v.SearchValue)
			}
		}

		resultPrint += fmt.Sprintln("  EXCLUDE:")
		for k, v := range r.Rules.Exclude {
			resultPrint += fmt.Sprintln("    ", k+1, ".")
			resultPrint += fmt.Sprintf("      searchField: '%s'\n", v.SearchField)
			resultPrint += fmt.Sprintf("      searchValue: '%s'\n", v.SearchValue)
			resultPrint += fmt.Sprintf("      accurateComparison: '%v'\n", v.AccurateComparison)
		}

		resultPrint += fmt.Sprintf("  PASSANY: '%v'\n", r.Rules.Passany)

		return resultPrint
	}

	BeforeAll(func() {
		//инициализация списка правил
		lr, lw, errGetRule = rules.NewListRule("placeholder_misp", "rules", "mispmsgrule.yaml")
	})

	Context("Тест 1. Чтение нового тестового файла, в другом формате", func() {
		It("Новый тестовый файл должен быть успешно прочитан", func() {
			//инициализация списка правил
			fmt.Println("NEW RULES FILE", fileName, ":")
			for k, v := range lw {
				fmt.Printf("%d. %s\n", k, v)
			}
			fmt.Println("new rule result:")
			fmt.Println(printRuleResult(lr))

			Expect(errGetRule).ShouldNot(HaveOccurred())
		})
	})
})
