package testing_test

import (
	"fmt"
	"placeholder_misp/datamodels"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("RegexpExample", Ordered, func() {
	var listAttributes *datamodels.ListAttributesMispFormat

	getNewListAttributes := func(al []datamodels.AttributesMispFormat, lat map[int][][2]string) []datamodels.AttributesMispFormat {
		countAttr := len(al)

		nal := make([]datamodels.AttributesMispFormat, 0, countAttr)

		for k, v := range al {
			if elem, ok := lat[k]; ok {
				for _, value := range elem {
					v.Category = value[0]
					v.Type = value[1]

					nal = append(nal, v)
				}

				continue
			}

			nal = append(nal, v)
		}

		return nal
	}

	BeforeAll(func() {
		listAttributes = datamodels.NewListAttributesMispFormat()
		listAttributes.SetValueValueAttributesMisp("test value 1", true)
		listAttributes.SetValueValueAttributesMisp("test value 2", true)
		listAttributes.HandlingValueEventIdAttributesMisp(88545, 1)
	})

	Context("Тест 2. Разбираем строку из observable.tags", func() {
		It("Функция обработки списка тегов с помощью регулярного выражения должна отработать успешно", func() {
			result := datamodels.HandlingListTags([]string{
				"misp:Network activity=\"email-src\"",
				"misp:Persistence mechanism=\"email-body\"",
				"misp;Bjkfo-ffkkk=\"nvii fefeni\"",
			})

			fmt.Println(result)

			Expect(len(result)).Should(Equal(2))
			Expect(result[0][0]).Should(Equal("Network activity"))
		})

		It("Должно быть пустое значение при обработке ошибочных тегов", func() {
			result := datamodels.HandlingListTags([]string{
				"misp:Network activity-\"email-src\"",
				"misp;Bjkfo-ffkkk=\"nvii fefeni\"",
			})

			fmt.Println(result)

			Expect(len(result)).Should(Equal(0))
		})
	})

	Context("Тест 2. Проверка заполнения списка attributeTags", func() {
		It("Должно быть добавлено определенной кол-во елементов в attributeTags", func() {
			listAttributes.HandlingValueTagsAttributesMisp([]string{
				"misp:Network activity=\"email-src\"",
				"misp:Persistence mechanism=\"email-body\"",
			}, false)

			Expect(listAttributes.GetCountListAttributesMisp()).Should(Equal(2))

			listAttributes.SetValueValueAttributesMisp("test value 3", true)
			listAttributes.SetValueCommentAttributesMisp("yytttttt", false)
			listAttributes.SetValueValueAttributesMisp("test value 4", true)
			listAttributes.HandlingValueEventIdAttributesMisp(1234, 4)
			listAttributes.HandlingValueTagsAttributesMisp([]string{
				"misp:Person=\"target-email\"",
				"misp:Antivirus detection=\"uri\"",
			}, false)

			Expect(listAttributes.GetCountListAttributesMisp()).Should(Equal(4))

			//fmt.Println("List attributes misp =", listAttributes.GetListAttributesMisp())

			lat := listAttributes.GetListAttributeTags()
			for k, v := range lat {
				fmt.Printf("Key: %d, Value: %v\n", k, v)
			}

			Expect(len(lat)).Should(Equal(2))

			ndamf := getNewListAttributes(listAttributes.GetListAttributesMisp(), listAttributes.GetListAttributeTags())

			for k, v := range ndamf {
				fmt.Printf("%d. %v\n", k, v)
			}

			Expect(len(ndamf)).Should(Equal(6))
		})
	})
})
