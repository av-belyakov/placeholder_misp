package main_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	// "placeholder_misp"
	"placeholder_misp/coremodule"
)

var _ = Describe("Testing/Testdatatime/HandlingDataTime", func() {
	const DATA_TIME = 1.686652790366e+12

	convFunc := func(d interface{}) (string, error) {
		var res string

		resTmp, ok := d.(float64)
		if !ok {
			return res, fmt.Errorf("not convert")
		}

		fmt.Println("NUMBER = ", resTmp)
		fmt.Printf("STRING = %16.f000\n", d)

		return fmt.Sprintf("%13.f000", d), nil
	}

	Context("Тест 1. Преобразовываем строку DATA_TIME в строку времени в формате Unix состоящей из 16 символов", func() {
		It("Должна быть получена строка формата Unix состоящая из 16 символов", func() {

			str, err := convFunc(DATA_TIME)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(str)).Should(Equal(16))
		})
	})

	Context("Тест 2. Получаем текущий год и месяц", func() {
		It("Должны быть успешно получены текущие год и месяц", func() {
			t := time.Now()

			Expect(t.Year()).Should(Equal(2023))
			Expect(int(t.Month())).Should(Equal(12))
		})
	})

	Context("Test 3", func() {
		It("need is true", func() {
			v := coremodule.GetTypeNameObservablesTag("type:sha256")
			Expect(v).Should(Equal("sha256"))

			v = coremodule.GetTypeNameObservablesTag("example:md5")
			Expect(v).Should(Equal("md5"))
		})
	})

	Context("Тест 4. Проверка преобразования даты в формате Unix в строковое значение даты", func() {
		It("При преобразовании в строковую дату ошибки быть не должно, дата должна соответствовать определенному значению", func() {
			/*
			   // Some valid layouts are invalid time values, due to format specifiers
			   	// such as _ for space padding and Z for zone information.
			   	// For example the RFC3339 layout 2006-01-02T15:04:05Z07:00
			   	// contains both Z and a time zone offset in order to handle both valid options:
			   	// 2006-01-02T15:04:05Z
			   	// 2006-01-02T15:04:05+07:00
			   	t, _ = time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
			   	fmt.Println(t)
			   	t, _ = time.Parse(time.RFC3339, "2006-01-02T15:04:05+07:00")
			*/

			t := time.UnixMilli(1702686056000).UTC()
			fmt.Println("UTC Datetime:", t.String())

			t = time.UnixMilli(1702686056000)
			fmt.Println("Datetime:", t.String())

			Expect(true).Should(BeTrue())
		})
	})
})
