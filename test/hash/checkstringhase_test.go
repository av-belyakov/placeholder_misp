package testhash_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/av-belyakov/placeholder_misp/coremodule"
	"github.com/av-belyakov/placeholder_misp/internal/datamodels"
	"github.com/av-belyakov/placeholder_misp/internal/supportingfunctions"
)

var _ = Describe("Checkstringhase", func() {
	testList := []coremodule.ChanInputCreateMispFormat{
		{
			FieldBranch: "observables.attachment.name",
			FieldName:   "mytextfile.txt",
		},
		{
			FieldBranch: "observables.attachment.size",
			FieldName:   "817",
		},
		{
			FieldBranch: "observables.attachment.hashes.0",
			FieldName:   "c29438b04791184d3eba39bdb7cf99560ab62068fee9509d50cf59723c398ac1",
		},
		{
			FieldBranch: "observables.attachment.hashes.1",
			FieldName:   "58861ef4c118cc3270b9871734ee54852a1374e5",
		},
		{
			FieldBranch: "observables.attachment.hashes.2",
			FieldName:   "7c531394dc2f483bc6c6c628c02e0788",
		},
	}

	testListStrings := []string{
		"294593fcb93a6d6694c9670e86e649bf",                                 //md5
		"fd861b0d33cc076ded2987c94fa9860e0c4aadd0",                         //sha1
		"6b3383ad0a767b008e8a41db84efea8847de86796aefd3703dcecb7ec3203e27", //sha256
		"c3f167e719aa944af2e80941ac629d39cec22308",                         //sha1
		"78cf6611f6928a64b03a57fe218c3cd4",                                 //md5
		"2c0961c22dc6caad6210759787fb149a837ee2db",                         //sha1
		"mytextfile.txt",
	}

	Context("Тест 1. Проверяем тип хеша", func() {
		It("Должна быть найдена одна ошибка, остальные типы хеша должны быть успешно определены", func() {
			var countSuccess, countError int

			for _, v := range testListStrings {
				hashType, stringSize, err := supportingfunctions.CheckStringHash(v)
				if err != nil {
					countError++
				} else {
					countSuccess++
				}

				fmt.Println("test value:", v, "hashType:", hashType, " stringSize:", stringSize)
			}

			Expect(countError).Should(Equal(1))
			Expect(countSuccess).Should(Equal(6))
		})
	})

	Context("Тест 2. Проверяем обработку attachment", func() {
		It("Attachment должно быть успешно обработанно", func() {
			nla := datamodels.NewListAttributeTmp()

			for _, tmf := range testList {
				nla.AddAttribute(tmf.FieldBranch, tmf.FieldName, 0)
			}

			la := nla.GetListAttribute()

			fmt.Println("List attribute:")
			for k, v := range la {
				fmt.Printf("%d.\n\t%s\n", k, v)
			}

			Expect(len(la)).Should(Equal(1))
		})
	})

	/*Context("Test 3", func() {
		It("test time", func() {
			firstSeen := fmt.Sprint(time.Now().UnixMilli()) //13
			timestamp := fmt.Sprint(time.Now().UnixMicro()) //16

			fslen := len(firstSeen)
			if fslen < 16 {
				firstSeen = firstSeen + strings.Repeat("0", 16-fslen)
			} else if fslen > 16 {
				firstSeen = firstSeen[:16]
			}

			if len(timestamp) > 10 {
				timestamp = timestamp[:10]
			}

			fmt.Println("First_seen:", firstSeen)
			fmt.Println("Timestamp:", timestamp)

			fs := float64(time.Now().UnixMilli())
			fmt.Printf("%13.f", fs)

			Expect(true).Should(BeTrue())
		})
	})

	Context("Тест 4. Проверяем наличие переменных окружения", func() {
		It("Должна быть найдена переменная окружения GO_PHMISP_MAIN", func() {
			v, ok := os.LookupEnv("GO_PHMISP_MAIN")
			Expect(ok).Should(BeFalse())
			Expect(v).Should(Equal(""))

			if !ok || v != "development" {
				fmt.Println("Is production")
			} else {
				fmt.Println("Is development")
			}

			os.Setenv("GO_PHMISP_MAIN", "development")
			v, ok = os.LookupEnv("GO_PHMISP_MAIN")
			Expect(ok).Should(BeTrue())
			Expect(v).Should(Equal("development"))
		})
	})*/
})
