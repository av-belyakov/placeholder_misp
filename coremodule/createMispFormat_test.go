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
		listRules          rules.ListRulesProcessingMsgMISP
		listFieldsMispType map[string][]coremodule.FieldsNameMapping
		eventsMisp         datamodels.EventsMispFormat
		exampleByte        []byte
		errGetRule         error
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

		listFieldsMispType = map[string][]coremodule.FieldsNameMapping{
			"events": {
				{InputFieldName: "event.object.title", MispFieldName: "info"},
				{InputFieldName: "event.object.startDate", MispFieldName: "timestamp"},
				{InputFieldName: "event.object.tlp", MispFieldName: "distribution"},
				{InputFieldName: "event.object.severity", MispFieldName: "threat_level_id"},
				{InputFieldName: "event.organisationId", MispFieldName: "org_id"},
				{InputFieldName: "event.object.updatedAt", MispFieldName: "sighting_timestamp"},
				{InputFieldName: "event.object.owner", MispFieldName: "event_creator_email"},
			},
			"attributes": {
				{InputFieldName: "event.object.tlp", MispFieldName: "tags"},
				{InputFieldName: "observables._id", MispFieldName: "object_id"},
				{InputFieldName: "observables.data", MispFieldName: "value"},
				{InputFieldName: "observables._createdAt", MispFieldName: "timestamp"},
				{InputFieldName: "observables.message", MispFieldName: "comment"},
				{InputFieldName: "observables.startDate", MispFieldName: "first_seen"},
			},
		}

		eventsMisp = datamodels.EventsMispFormat{
			Analysis:          getAnalysis(),
			Timestamp:         "0",
			ThreatLevelId:     "4",
			PublishTimestamp:  "0",
			SightingTimestamp: "0",
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
					for _, value := range listFieldsMispType["events"] {
						if v.FieldBranch == value.InputFieldName {
							fmt.Printf("\nFieldBranch: %s\nFieldName: %s\nValue: %v\n  MispFieldName: %s\n", v.FieldBranch, v.FieldName, v.Value, value.MispFieldName)

							/*
								Здесь надо продумать возможность использования интерфейса который
								соответствовал бы методу который будет добавлять значения в поля
								event тип. Так как в v.Value имеет тип interface{} то метод
								setValue реализовать не сложно, кроме того можно будет избавится
								от switch и использовать что то типа map[string]интерфейс
							*/

							switch value.MispFieldName {
							case "info":
								eventsMisp.Info = v.Value
							case "timestamp":

							case "distribution":

							case "threat_level_id":

							case "org_id":

							case "sighting_timestamp":

							case "event_creator_email":

							}

						}
					}

					//что бы не выполнялось
					//if v.Value == "111" {
					//	fmt.Printf("\n RESEIVED MESSAGE:\n - FieldName: %s\n - ValueType: %s\n - Value: %v\n - FieldBranch: %s\n", v.FieldName, v.ValueType, v.Value, v.FieldBranch)
					//}
				}
			}()

			ok, warningMsg := procMsgHive.HandleMessage(chanOutMispFormat)

			fmt.Println("EventsType: ", eventsMisp)

			Expect(len(warningMsg)).Should(Equal(0))
			Expect(ok).Should(BeTrue())
		})
	})
})
