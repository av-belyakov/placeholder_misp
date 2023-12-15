package testreadfileevents_test

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Readfile", func() {
	Context("Test 1", func() {
		It("Should by is success", func() {
			f, err := os.OpenFile("../../natsinteractions/test_json/events_1696583308", os.O_RDONLY, os.ModePerm)

			Expect(err).ShouldNot(HaveOccurred())

			wr := bytes.Buffer{}
			sc := bufio.NewScanner(f)
			for sc.Scan() {
				wr.WriteString(sc.Text())
				//newResult = append(newResult, sc.Bytes()...)
			}

			list := strings.Split(wr.String(), "EVENTS:")

			//for k, v := range list {
			//	fmt.Printf("%d./n	%s/n", k, v)
			//}
			fmt.Println("000000000000000000", list[1])

			Expect(true).Should(BeTrue())
		})
	})
})
