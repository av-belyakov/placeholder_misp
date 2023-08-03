package coremodule_test

import (
	"fmt"
	"strconv"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"placeholder_misp/coremodule"
	"placeholder_misp/datamodels"
	rules "placeholder_misp/rulesinteraction"
	"placeholder_misp/tmpdata"
)

var _ = Describe("CreateMispFormat", Ordered, func() {
	var (
		listRules       rules.ListRulesProcessingMsgMISP
		listHandlerMisp map[string][]func(interface{}) bool
		eventsMisp      datamodels.EventsMispFormat
		attributesMisp  datamodels.AttributesMispFormat
		exampleByte     []byte
		errGetRule      error
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

	getAnalysis := func() string {
		return "2"
	}

	BeforeAll(func() {
		exampleByte = getExampleByte()
		listRules, _, errGetRule = rules.GetRuleProcessingMsgForMISP("rules", "procmispmsg_test.yaml")

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
	})

	Context("Тест 1. Проверка формирования правил фильтрации", func() {
		It("При формировании правил фильтрации ошибки быть не должно", func() {
			Expect(errGetRule).ShouldNot(HaveOccurred())
		})
	})

	Context("Тест 2. Проверяем правельность обработки правил для формирования MISP форматов", func() {
		It("Обработка правил для формирования MISP форматов должна выполнятся успешно", func() {
			chanOutMispFormat := make(chan coremodule.ChanInputCreateMispFormat)

			procMsgHive, err := coremodule.NewHandleMessageFromHive(exampleByte, listRules)
			Expect(err).ShouldNot(HaveOccurred())

			/*
				Правила отрабатывает нормально, теперь необходимо сделать
				формирование типов MISP типа event и attributes
			*/

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
	})
})
