package testupdatestdout_test

import (
	"bufio"
	"fmt"
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Stdoutupdate", Ordered, func() {
	getAppName := func(pf string, nl int) (string, error) {
		var line string

		f, err := os.OpenFile(pf, os.O_RDONLY, os.ModePerm)
		if err != nil {
			return line, err
		}
		defer f.Close()

		num := 1
		sc := bufio.NewScanner(f)
		for sc.Scan() {
			if num == nl {
				return sc.Text(), nil
			}

			num++
		}

		return line, nil
	}

	sendingInt := func(numBegin int) chan int {
		sendInt := make(chan int)

		go func() {
			defer close(sendInt)

			for i := numBegin; i < 100; i = i + numBegin {
				sendInt <- i

				time.Sleep(1 * time.Second)
			}
		}()

		return sendInt
	}

	sendInt := sendingInt(10)

	Context("Тест 1. Вывод информационного дачборда", func() {
		It("Должен быть успешно выведен информационный дачборд", func(ctx SpecContext) {
			appN := "placeholder_misp"
			an, err := getAppName("../../README.md", 1)

			Expect(err).ShouldNot(HaveOccurred())
			appN = an
			/*msg := `The application %s is running

			Всего событий полученно: %d
			Соответствуют правилам: %d
			--------------------------------
			`*/

			for d := range sendInt {
				fmt.Printf("The application %s is running", appN)
				fmt.Printf("Всего событий полученно: %d\r", d)
				fmt.Printf("Соответствуют правилам: %d\r", d-2)
				//fmt.Printf("%s\r", fmt.Sprintf(msg, appN, d, d-2))
			}

			fmt.Println("")

			Expect(true).Should(BeTrue())
		}, SpecTimeout(time.Second*10))
	})

	Context("Тест 2. Обновляем вывод информации в консоли без перевода сторки", func() {
		It("Должно быть успешно обновлены данные в одной и тойже строке", func() {
			var i int

			for i = 3; i >= 0; i-- {
				fmt.Printf("\033[2K\r%d", i)
				time.Sleep(1 * time.Second)
			}
			fmt.Println()

			Expect(i).Should(Equal(-1))
		})
	})
})
