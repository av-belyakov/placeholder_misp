package coremodule_test

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"placeholder_misp/confighandler"
	"placeholder_misp/coremodule"
	"placeholder_misp/datamodels"
	"placeholder_misp/memorytemporarystorage"
	"placeholder_misp/mispinteractions"
	rules "placeholder_misp/rulesinteraction"
	"placeholder_misp/tmpdata"
)

/*
type Sizer interface {
	Size() uintptr
}

type Producer[T any] struct {
	Val T
}

// Producer implements the Sizer interface

func (p Producer[T]) Size(T) uintptr {
	return unsafe.Sizeof(p.Val)
}

func (p Producer[T]) Produce() T {
	return p.Val
}

func Execute[T any](s Sizer[T]) {
	switch p := s.(type) {
	case Producer[T]:
		fmt.Println(p.Produce())
	default:
		panic("This should not happen")
	}
}

coremodule.Execute[string](coremodule.Producer[string]{"23"})
coremodule.Execute[string](coremodule.Producer[int64]{27})


func PrintStringOrInt[T string | int](v T) {
	switch any(v).(type) {
	case string:
		fmt.Printf("String: %v\n", v)
	case int:
		fmt.Printf("Int: %v\n", v)
	default:
		panic("Impossible")
	}
}
*/

var _ = Describe("CreateMispFormat", Ordered, func() {
	var (
		listRule             rules.ListRulesProcessingMsgMISP
		storageApp           *memorytemporarystorage.CommonStorageTemporary
		logging              chan datamodels.MessageLogging
		mispOutput           <-chan mispinteractions.SettingChanOutputMISP
		moduleMisp           *mispinteractions.ModuleMISP
		exampleByte          []byte
		errGetRule, errHMisp error
	)

	getExampleByte := func() []byte {
		byteTmp := []byte{}

		for _, v := range strings.Split(tmpdata.GetExampleDataThree(), " ") {
			i, err := strconv.Atoi(v)
			if err != nil {
				continue
			}

			byteTmp = append(byteTmp, uint8(i))
		}

		return byteTmp
	}

	//getAnalysis := func() string {
	//	return "2"
	//}

	BeforeAll(func() {
		//инициализация модуля конфига
		ca, _ := confighandler.NewConfig()

		//канал для логирования
		logging = make(chan datamodels.MessageLogging)

		//пример кейса в виде байтов
		exampleByte = getExampleByte()

		//инициализация списка правил
		listRule, _, errGetRule = rules.GetRuleProcessingMsgForMISP("rules", "procmispmsg_test.yaml")

		//fmt.Println("list RULES = ", listRule)

		//инициализируем модуль временного хранения информации
		storageApp = memorytemporarystorage.NewTemporaryStorage()

		//инициализация модуля для взаимодействия с MISP
		moduleMisp, errHMisp = mispinteractions.HandlerMISP(*ca.GetAppMISP(), storageApp, logging)

		mispOutput = moduleMisp.GetDataReceptionChannel()
	})

	Context("Тест 1. Проверка формирования правил фильтрации", func() {
		It("При конвертировании времени в формат ISO 8601 ошибок быть не должно", func() {
			dateTest := "1692322715497"
			dt, err := strconv.ParseInt(dateTest, 10, 64)

			Expect(err).ShouldNot(HaveOccurred())

			t := time.UnixMilli(dt)
			tstr := t.Format(time.RFC3339)

			fmt.Println("		Time format RFC3339 as string: ", tstr)

			Expect(true).Should(BeTrue())
		})

		It("При формировании правил фильтрации ошибки быть не должно", func() {
			Expect(errGetRule).ShouldNot(HaveOccurred())
		})

		It("При инициализации модуля взаимодействия с MISP ошибки быть не должно", func() {
			Expect(errHMisp).ShouldNot(HaveOccurred())
		})
	})

	/*
		//
		//
		// Для тестов незабывать устанавливать переменную окружения
		// export GO_PHMISP_MAIN="development"
		// что бы взаимодействовать с misp-world.cloud.gcm
		//
		//
	*/

	Context("Тест 2. Проверяем правельность обработки сообщения по правилам и корректность формирования MISP форматов", func() {
		It("Сообщение должно быть успешно обработано, должены быть сформированы типы Events и Attributes", func(ctx SpecContext) {
			cd := make(chan struct{})

			uuidTask := uuid.NewString()
			storageApp.SetRawDataHiveFormatMessage(uuidTask, exampleByte)

			//формирование итоговых документов в формате MISP
			chanCreateMispFormat, chanDone := coremodule.NewMispFormat(moduleMisp, logging)

			go coremodule.HandlerMessageFromHive(exampleByte, uuidTask, storageApp, listRule, chanCreateMispFormat, chanDone, logging)

			go func() {
				for {
					select {
					case log := <-logging:
						fmt.Println("___ Log = ", log, " ____")
					case <-cd:
						fmt.Println("-== STOP TEST ==-")

						return
					}
				}
			}()

			fmt.Println("reseived data: ", <-mispOutput)
			cd <- struct{}{}

			Expect(true).Should(BeTrue())
		}, SpecTimeout(time.Second*15))
	})

	Context("Тест 3. Тестируем удаление выбранных событий", func() {
		It("Выбранное событие должно быть успешно удалено", func() {
			res, err := mispinteractions.DelEventsMispFormat("misp-world.cloud.gcm", "TvHkjH8jVQEIdvAxjxnL4H6wDoKyV7jobDjndvAo", "6029")

			Expect(err).ShouldNot(HaveOccurred())
			Expect(res.StatusCode).Should(Equal(http.StatusOK))
		})
	})

	/*Context("Тест 2. Проверяем правельность обработки правил для формирования MISP форматов", func() {
		It("Обработка правил для формирования MISP форматов должна выполнятся успешно", func() {
			chanOutMispFormat := make(chan coremodule.ChanInputCreateMispFormat)

			procMsgHive, err := coremodule.NewHandleMessageFromHive(exampleByte, listRules)
			Expect(err).ShouldNot(HaveOccurred())

			go func() {
				for v := range chanOutMispFormat {

					//fmt.Println("___ v.FieldName = ", v.FieldName, " v.Value = ", v.Value, " ____")

					if v.FieldBranch == "event.object.tlp" {
						fmt.Println("|||| 'event.object.tlp' |||| FieldBranch: ", v.FieldBranch, " v.Value: ", v.Value, " |||||||||")
					}

					if lf, ok := listHandlerMisp[v.FieldBranch]; ok {
						for _, f := range lf {
							_ = f(v.Value)
						}
					}
				}
			}()

			ok, warningMsg := procMsgHive.HandleMessage(chanOutMispFormat)

			fmt.Println("EventsType: ", eventsMisp, " AttributesType: ", attributesMisp)

			Expect(len(warningMsg)).Should(Equal(0))
			Expect(ok).Should(BeTrue())
		})
	})*/
})
