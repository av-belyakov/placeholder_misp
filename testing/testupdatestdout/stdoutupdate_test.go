package testupdatestdout_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Stdoutupdate", Ordered, func() {
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
			/*msg := `The application %s is running

			Всего событий полученно: %d
			Соответствуют правилам: %d
			--------------------------------
			`*/

			for d := range sendInt {
				fmt.Printf("\tСОБЫТИЙ получено/обработано - %d/%d\r", d, d-2)
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
