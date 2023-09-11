package tmpdata

import (
	"bufio"
	"fmt"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ReadFileReadme", func() {
	Context("Тест 1. Читаем файл README.md", func() {
		It("При чтении файла ошибок быть не должно", func() {
			f, err := os.OpenFile("../README.md", os.O_RDONLY, os.ModePerm)
			Expect(err).ShouldNot(HaveOccurred())
			defer f.Close()

			num := 1
			sc := bufio.NewScanner(f)
			for sc.Scan() {
				str := sc.Text()

				fmt.Printf("строка %d. %s\n", num, str)
				num++
			}

			Expect(sc.Err()).ShouldNot(HaveOccurred())
		})
	})
})
