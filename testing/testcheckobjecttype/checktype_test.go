package testcheckobjecttype_test

import (
	"fmt"
	"reflect"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Checktype", Ordered, func() {
	var (
		testMap   map[string]interface{}
		testSlice []interface{}
	)

	type listTestElem struct {
		Name    string
		Age     int
		Size    int
		Weight  int
		Address string
	}

	BeforeAll(func() {
		testMap = map[string]interface{}{
			"one": listTestElem{
				Name:    "Elena",
				Age:     23,
				Size:    167,
				Weight:  58,
				Address: "country, city",
			},
			"two": listTestElem{
				Name:    "Olga",
				Age:     25,
				Size:    172,
				Weight:  63,
				Address: "country, city",
			},
			"free": listTestElem{
				Name:    "Nady",
				Age:     20,
				Size:    171,
				Weight:  66,
				Address: "country, city",
			},
		}

		testSlice = []interface{}{
			listTestElem{
				Name:    "Elena",
				Age:     23,
				Size:    167,
				Weight:  58,
				Address: "country, city",
			},
			listTestElem{
				Name:    "Olga",
				Age:     25,
				Size:    172,
				Weight:  63,
				Address: "country, city",
			},
			listTestElem{
				Name:    "Nady",
				Age:     20,
				Size:    171,
				Weight:  66,
				Address: "country, city",
			},
		}
	})

	Context("Тест 1. Проверяем определение типа map", func() {
		It("При определении типа map не должно быть ошибки", func() {
			r := reflect.TypeOf(testMap)

			fmt.Println("kind: ", r.Kind().String())
			fmt.Println("is Map", r.Kind() == reflect.Map)

			Expect(r.Kind()).Should(Equal(reflect.Map))
		})
	})

	Context("Тест 2. Проверяем определение типа slice", func() {
		It("При определении типа map не должно быть ошибки", func() {
			r := reflect.TypeOf(testSlice)

			fmt.Println("kind: ", r.Kind().String())
			fmt.Println("is Slice", r.Kind() == reflect.Slice)

			Expect(r.Kind()).Should(Equal(reflect.Slice))
		})
	})
})
