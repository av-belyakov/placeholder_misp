package confighandler_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"

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

			fmt.Println("conf = ", conf)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(conf.LogList)).Should(Equal(2))
			Expect(conf.LogList[0].PathDirectory).Should(Equal("logs"))
			Expect(conf.LogList[0].MsgTypeName).Should(Equal("error"))
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

	Context("Тест 2. Тестируем функцию возвращающую путь до текущего приложения", func() {
		It("Должен построен определенный путь до текущего приложения", func() {
			getRootPath := func(rootDir string) (string, error) {
				currentDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
				if err != nil {
					return "", err
				}

				tmp := strings.Split(currentDir, "/")

				if tmp[len(tmp)-1] == rootDir {
					return currentDir, nil
				}

				var path string = ""
				for _, v := range tmp {
					path += v + "/"

					if v == rootDir {
						return path, nil
					}
				}

				return path, nil
			}

			rootPath, err := getRootPath("placeholder_misp")

			Expect(err).ShouldNot(HaveOccurred())
			Expect(rootPath).Should(Equal("/home/artemij/go/src/placeholder_misp/"))
		})
	})

	Context("Тест 3. Проверка выполнения анмаршалинга конфигурационного файла", func() {
		It("Должен быть выполнен анмаршалинг части основного конфигурационного файла", func() {
			type Logs struct {
				Log []confighandler.LogSet
			}

			ls := Logs{}

			viper.SetConfigFile("../configs/config.yaml")
			viper.SetConfigType("yaml")
			errRC := viper.ReadInConfig()

			Expect(errRC).ShouldNot(HaveOccurred())

			ok := viper.IsSet("LOGGING")
			err := viper.GetViper().Unmarshal(&ls)

			fmt.Println("result unmarshal config file LOG: ", ls)

			Expect(ok).Should(BeTrue())
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
