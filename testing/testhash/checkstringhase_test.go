package testhash_test

import (
	"fmt"
	"os"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"placeholder_misp/coremodule"
	"placeholder_misp/datamodels"
)

func CheckHash(h []string) map[string]string {
	result := make(map[string]string, len(h))

	for _, v := range h {
		switch len(v) {
		case 32:
			result[v] = "md5"
		case 40:
			result[v] = "sha1"
		case 64:
			result[v] = "sha256"
		default:
			result[v] = "other"
		}
	}

	return result
}

/*
    {
      "template_uuid": "c8cc27a6-4bd31-1f72-afa5-7b9bb4ac3b3b",
      "template_version": "1",
      "first_seen": "1581984000000000",
      "timestamp": "1617875568",
      "name": "file",
      "description": "size 817 byte",
      "event_id": "12660",
      "meta-category": "file",
      "distribution": "5",
      "Attribute": [
        {
          "category": "Payload delivery",
          "type": "other",
          "value": "n[n.txt",
          "distribution": "0",
          "object_relation": "filename"
        },
 	  ],
	}
*/

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

	Context("Тест 1. Проверяем тип хеша", func() {
		It("Должно быть определен тип хеша", func() {
			hashes := []string{
				"c29438b04791184d3eba39bdb7cf99560ab62068fee9509d50cf59723c398ac1",
				"58861ef4c118cc3270b9871734ee54852a1374e5",
				"7c531394dc2f483bc6c6c628c02e0788",
			}

			r := CheckHash(hashes)

			Expect(len(r)).Should(Equal(3))
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

	Context("Test 3", func() {
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
	})
})
