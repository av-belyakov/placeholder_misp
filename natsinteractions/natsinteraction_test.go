package natsinteractions_test

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"placeholder_misp/confighandler"
	"placeholder_misp/datamodels"
	"placeholder_misp/natsinteractions"
	"placeholder_misp/tmpdata"
)

func readReflectAnyType(name, anyType interface{}) string {
	var str, nameStr string
	r := reflect.TypeOf(anyType)

	if n, ok := name.(int); ok {
		nameStr = fmt.Sprintf("%v.", n)
	} else if n, ok := name.(string); ok {
		name = fmt.Sprintln(n, ": ")
	}

	if r == nil {
		return str
	}

	switch r.Kind() {
	case reflect.String:
		str += fmt.Sprintf("\t - %s %s\n", nameStr, reflect.ValueOf(anyType).String())

	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		str += fmt.Sprintf("\t - %s %d\n", nameStr, reflect.ValueOf(anyType).Int())

	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		str += fmt.Sprintf("\t - %s %d\n", nameStr, reflect.ValueOf(anyType).Uint())

	case reflect.Float32, reflect.Float64:
		str += fmt.Sprintf("\t - %s %v\n", nameStr, reflect.ValueOf(anyType).Float())

	case reflect.Bool:
		str += fmt.Sprintf("\t - %s %v\n", nameStr, reflect.ValueOf(anyType).Bool())
	}

	return str
}

func readReflectMap(list map[string]interface{}) string {
	var str string

	for k, v := range list {
		r := reflect.TypeOf(v)

		if r == nil {
			return str
		}

		str += readReflectAnyType(k, v)

		switch r.Kind() {
		case reflect.Map:
			if v, ok := v.(map[string]interface{}); ok {
				str += readReflectMap(v)
			}

		case reflect.Slice:
			if v, ok := v.([]interface{}); ok {
				str += redReflectSlice(v)
			}

		case reflect.Array:
			str += fmt.Sprintf("%s: %s (it is array)\n", k, reflect.ValueOf(v).String())
		}
	}

	return str
}

func redReflectSlice(list []interface{}) string {
	var str string

	for k, v := range list {
		r := reflect.TypeOf(v)

		if r == nil {
			return str
		}

		str += readReflectAnyType(k, v)

		switch r.Kind() {
		case reflect.Map:
			if v, ok := v.(map[string]interface{}); ok {
				str += readReflectMap(v)
			}

		case reflect.Slice:
			if v, ok := v.([]interface{}); ok {
				str += redReflectSlice(v)
			}

		case reflect.Array:
			str += fmt.Sprintf("%d. %s (it is array)\n", k, reflect.ValueOf(v).String())
		}
	}

	return str
}

var _ = Describe("Natsinteraction", Ordered, func() {
	var (
		ctx          context.Context
		errConn      error
		closeCtx     context.CancelFunc
		enumChannels *natsinteractions.EnumChannelsNATS
	)

	BeforeAll(func() {
		ctx, closeCtx = context.WithTimeout(context.Background(), 2*time.Second)

		enumChannels, errConn = natsinteractions.NewClientNATS(ctx, confighandler.AppConfigNATS{
			//Host: "nats.cloud.gcm",
			Host: "127.0.0.1",
			Port: 4222,
		})
	})

	Context("Тест 1.1. Проверка декадирования тестовых данных из файла 'binaryDataOne'", func() {
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

	Context("Тест 1.2. Проверка декадирования тестовых данных из файла 'binaryDataTwo'", func() {
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

	Context("Тест 1.3. Проверка декадирования тестовых данных из файла 'binaryDataThree'", func() {
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

			fmt.Println("---- ExampleDataThree ----")
			fmt.Println(mm.ToStringBeautiful(0))
			fmt.Println("--------------------------")

			Expect(err).ShouldNot(HaveOccurred())
		})

		It("При анмаршалинге в НЕИЗВЕСТНЫЙ тип ошибки быть не должно", func() {
			result := map[string]interface{}{}

			err := json.Unmarshal(exampleByte, &result)

			/*
				Можно попробовать описать все типы для объектов сообщений из примера,
				потом путем определения наименования поля подставлять типы формируя итоговый тип.
				Все наименования полей для которых нет типа записывать в отдельный файл для их дольнейшго
				анализа.
				Вопрос, как формировать итоговый тип который должен состоять из переменного числа объектов?
			*/

			fmt.Println("---- REFLECTION MAPPING ExampleDataThree ----")
			fmt.Println(fmt.Println(readReflectMap(result)))
			fmt.Println("---------------------------------------------")

			/* сделал но надо облагородить с пробелами и нет наименование полей почему то */

			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("Тест 2. Проверка инициализации соединения с NATS", func() {
		It("При инициализации соединения с NATS не должно быть ошибки", func() {
			Expect(errConn).ShouldNot(HaveOccurred())

			fmt.Println("Resevid message = ", <-enumChannels.GetDataReceptionChannel())
		})
	})

	AfterAll(func() {
		closeCtx()
	})
})
