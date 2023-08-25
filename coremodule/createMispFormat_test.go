package coremodule_test

import (
	"context"
	"fmt"
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
		listRule    rules.ListRulesProcessingMsgMISP
		storageApp  *memorytemporarystorage.CommonStorageTemporary
		loging      chan datamodels.MessageLoging
		moduleMisp  *mispinteractions.ModuleMISP
		exampleByte []byte
		errGetRule  error

		testChan chan struct {
			Status     string
			StatusCode int
			Body       []byte
		}
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
		//PrintStringOrInt("hello")
		//PrintStringOrInt(42)

		testChan = make(chan struct {
			Status     string
			StatusCode int
			Body       []byte
		})

		ca, _ := confighandler.NewConfig()
		mts := memorytemporarystorage.NewTemporaryStorage()

		loging = make(chan datamodels.MessageLoging)
		exampleByte = getExampleByte()
		listRule, _, errGetRule = rules.GetRuleProcessingMsgForMISP("rules", "procmispmsg_test.yaml")

		//инициализируем модуль временного хранения информации
		storageApp = memorytemporarystorage.NewTemporaryStorage()
		ctxMISP, _ := context.WithTimeout(context.Background(), 2*time.Second)
		moduleMisp, _ = mispinteractions.HandlerMISP(ctxMISP, *ca.GetAppMISP(), mts, testChan, loging)

		/*
			eventsMisp = datamodels.EventsMispFormat{
				Analysis:          getAnalysis(),
				Timestamp:         "0",
				ThreatLevelId:     "4",
				PublishTimestamp:  "0",
				SightingTimestamp: "0",
			}

			listHandlerMisp = map[string][]func(interface{}) bool{
				"event.object.title":     {eventsMisp.SetValueInfoEventsMisp},
				"event.object.startDate": {eventsMisp.SetValueTimestampEventsMisp},
				"event.details.endDate":  {eventsMisp.SetValueDateEventsMisp},
				"event.object.tlp":       {eventsMisp.SetValueDistributionEventsMisp},
				"event.object.severity":  {eventsMisp.SetValueThreatLevelIdEventsMisp},
				"event.organisationId":   {eventsMisp.SetValueOrgIdEventsMisp},
				"event.object.updatedAt": {eventsMisp.SetValueSightingTimestampEventsMisp},
				"event.object.owner":     {eventsMisp.SetValueEventCreatorEmailEventsMisp},
				"observables._id":        {attributesMisp.SetValueObjectIdAttributesMisp},
				"observables.data":       {attributesMisp.SetValueValueAttributesMisp},
				"observables._createdAt": {attributesMisp.SetValueTimestampAttributesMisp},
				"observables.message":    {attributesMisp.SetValueCommentAttributesMisp},
				"observables.startDate":  {attributesMisp.SetValueFirstSeenAttributesMisp},
			}
		*/
	})

	Context("Тест 1. Проверка формирования правил фильтрации", func() {
		It("При формировании правил фильтрации ошибки быть не должно", func() {
			Expect(errGetRule).ShouldNot(HaveOccurred())
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

			uuidTask := uuid.NewString()
			storageApp.SetRawDataHiveFormatMessage(uuidTask, exampleByte)

			//формирование итоговых документов в формате MISP
			chanCreateMispFormat, chanDone := coremodule.NewMispFormat(uuidTask, moduleMisp, loging)

			go coremodule.HandlerMessageFromHive(uuidTask, storageApp, listRule, chanCreateMispFormat, chanDone, loging)

			/*
				data := <-testChan
				fmt.Println("Resiver DATA: ", data)


					ed, ok := data["events"]
					fmt.Println("EVENTS: ", ed)
					Expect(ok).ShouldNot(Equal(BeFalse()))

					ad, ok := data["attributes"]
					fmt.Println("ATTRIBUTES: ", ad)
					Expect(ok).ShouldNot(Equal(BeFalse()))
			*/

			resp := <-testChan

			fmt.Println("RESPONSE status: ", resp.Status)
			fmt.Println("RESPONSE status code: ", resp.StatusCode)
			fmt.Println("RESPONSE count body = ", len(resp.Body))

			//str, err := supportingfunctions.NewReadReflectJSONSprint(resp.Body)
			//str := []interface{}{}
			//err := json.Unmarshal(resp.Body, &str)
			//Expect(err).ShouldNot(HaveOccurred())
			fmt.Printf("BODY:\n%s\n", string(resp.Body))

			Expect(true).Should(BeTrue())
		}, SpecTimeout(time.Second*15))
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
