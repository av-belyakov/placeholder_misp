package rules_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"placeholder_misp/rules"
)

var _ = Describe("Readrule", func() {
	Context("Тест 1. Чтение файла с правилами", func() {
		It("При чтении файла с правилами ошибок быть не должно, файл должен быть успешно прочитан", func() {
			lrp, err := rules.GetRuleProcessedMISPMsg("rules", "processedmispmsg.yaml")

			fmt.Println("RULE processedmispmsg.yaml: ", lrp)

			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
