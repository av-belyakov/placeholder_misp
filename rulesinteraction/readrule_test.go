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
			resultPrint += fmt.Sprintln("  listRequiredValues:")

			for key, value := range v.ListRequiredValues {
				resultPrint += fmt.Sprintf("   %d.\n", key)
				resultPrint += fmt.Sprintf("    fieldSearchName: %s.\n", value.FieldSearchName)
				resultPrint += fmt.Sprintf("    typeValue: %s.\n", value.TypeValue)
				//resultPrint += fmt.Sprintf("    fieldValue: %s.\n", value.FindValue)
				resultPrint += fmt.Sprintf("    replaceValue: %s.\n", value.ReplaceValue)
			}
		}

		return resultPrint
	}

	printSearchName := func(name string, list map[string][][2]int) string {
		strRes := fmt.Sprintf("%s:\n", name)

		for key, value := range list {
			strRes += fmt.Sprintf("  '%s': %v\n", key, value)
		}

		return strRes
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
			fmt.Println("new rule result:")
			fmt.Println(printRuleResult(lrp.Rules))
			fmt.Println(printSearchName("SearchFieldsName", lrp.SearchFieldsName))
			fmt.Println(printSearchName("SearchValuesName", lrp.SearchValuesName))

			fmt.Println("list verification warning:")
			fmt.Println(printVerificationWarning(lvw))

			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	/*Context("Тест 2. Чтение тестового файла procmispmsg_test_1.yaml с правилами", func() {
		It("При чтении файла с правилами ошибок быть не должно, файл должен быть успешно прочитан", func() {
			lrp, lvw, err := rules.GetRuleProcessedMISPMsg("rules", "procmispmsg_test_1.yaml")

			fmt.Println("2. _________ RULE procmispmsg_test_1.yaml.")
			fmt.Println("new result:")
			fmt.Println(printRuleResult(lrp.Rules))

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
			fmt.Println(printRuleResult(lrp.Rules))

			fmt.Println("list verification warning:")
			fmt.Println(printVerificationWarning(lvw))

			Expect(err).ShouldNot(HaveOccurred())
		})
	})*/

	Context("Тест 4. Чтение тестового файла procmispmsg_test_error.yaml с ошибочным построение правил", func() {
		It("При чтении файла с правилами ошибок быть не должно, файл должен быть успешно прочитан", func() {
			lrp, lvw, err := rules.GetRuleProcessedMISPMsg("rules", "procmispmsg_test_error.yaml")

			/*
				Если сделать два типа списков lf (список названий свойств) b lv (список названий искомых значений). Так как оба
				списка будут в виде map[string][][2]int{} где свойство map будет наименование свойства или искомого значения, то
				выполняем поиск наименования свойства (из сообщения) и значения (из сообщения) по спискам lf, lv. Для actionType
				pass, reject должны совпасть оба, для replace можно только lv. Дальше можно найти исходное правило по полученным
				номерам и выполнить дальнейшую обработку.


						тольк для правли типа pass и reject
				   1. создать отдельные списки полей и значений в которые занести порядковый номер правила где они есть по типу 0.1. список map[string] где
				      свойство map будет имя поля или само искомое значение
				   2. сравнивать поле и значение по мере обработке элементов списка
				   2.1 ищем имя поля в списке lf
				   2.2 при совпадении сравниваем значение из общего правила с значением поля в списке

				   			для правила типа replace
			*/

			fmt.Println("list verification warning:")
			fmt.Println(printVerificationWarning(lvw))

			fmt.Println("4. _________ RULE procmispmsg_test_error.yaml.")
			fmt.Println("new rule result:")
			fmt.Println(printRuleResult(lrp.Rules))
			fmt.Println(printSearchName("SearchFieldsName", lrp.SearchFieldsName))
			fmt.Println(printSearchName("SearchValuesName", lrp.SearchValuesName))

			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("Тест 5. Чтение нового тестового файла, в другом формате", func() {
		It("Новый тестовый файл должен быть успешно прочитан", func() {
			r, err := rules.GetRuleProcessingMsgForMISP("rules", "procmispmsg_test.yaml")

			fmt.Println("NEW RULES FILE:", r)

			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
