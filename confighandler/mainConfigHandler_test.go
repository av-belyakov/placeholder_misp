package confighandler_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"placeholder_misp/confighandler"
)

var _ = Describe("MainConfigHandler", Ordered, func() {
	BeforeAll(func() {
		os.Setenv("GO_PHMISP_MAIN", "")
		os.Setenv("GO_PHMISP_MHOST", "")
		os.Setenv("GO_PHMISP_MPORT", "")
		os.Setenv("GO_PHMISP_NHOST", "")
		os.Setenv("GO_PHMISP_NPORT", "")
	})

	Context("Тест 1. Проверяем работу функции NewConfig с разными значениями переменной окружения GO_PHMISP_MAIN", func() {
		It("Должно быть получено содержимое общего файла 'config.yaml'", func() {
			conf, err := confighandler.NewConfig()

			Expect(err).ShouldNot(HaveOccurred())
			Expect(conf.Logging.Stdout).Should(BeTrue())
			Expect(conf.Logging.FilePath).Should(Equal("/logs"))
			Expect(len(conf.Logging.FileType)).Should(Equal(2))
			Expect(conf.Logging.MaxSize).Should(Equal(1024))
		})

		It("Должно быть получено содержимое файла 'config_prod.yaml' при пустом значении переменной GO_PHMISP_MAIN", func() {
			os.Setenv("GO_PHMISP_MAIN", "")

			conf, err := confighandler.NewConfig()

			Expect(err).ShouldNot(HaveOccurred())
			Expect(conf.GetAppNATS().Host).Should(Equal("192.168.9.201"))
			Expect(conf.GetAppNATS().Port).Should(Equal(4222))
			Expect(conf.GetAppMISP().Host).Should(Equal("192.168.9.202"))
			Expect(conf.GetAppMISP().Port).Should(Equal(8080))
		})

		It("Должно быть получено содержимое файла 'config_dev.yaml' при значении переменной GO_PHMISP_MAIN=development", func() {
			os.Setenv("GO_PHMISP_MAIN", "development")

			conf, err := confighandler.NewConfig()

			Expect(err).ShouldNot(HaveOccurred())
			Expect(conf.GetAppNATS().Host).Should(Equal("127.0.0.1"))
			Expect(conf.GetAppNATS().Port).Should(Equal(4222))
			Expect(conf.GetAppMISP().Host).Should(Equal("192.168.13.3"))
			Expect(conf.GetAppMISP().Port).Should(Equal(8181))
		})

		It("Должно быть перезаписанно содержимое некоторых полей полученных с переменных окружения GO_PHMISP_*", func() {
			os.Setenv("GO_PHMISP_MAIN", "development")
			os.Setenv("GO_PHMISP_MHOST", "78.87.78.87")
			os.Setenv("GO_PHMISP_NPORT", "11111")

			conf, err := confighandler.NewConfig()

			Expect(err).ShouldNot(HaveOccurred())
			Expect(conf.GetAppMISP().Host).Should(Equal("78.87.78.87"))
			Expect(conf.GetAppNATS().Port).Should(Equal(11111))
		})
	})
})
