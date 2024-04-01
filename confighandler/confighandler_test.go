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

			//fmt.Println("conf = ", conf)

			//for k, v := range conf.GetListOrganization() {
			//	fmt.Printf("%d. OrgName: %s, SourceName: %s\n", k, v.OrgName, v.SourceName)
			//}

			commonApp := conf.GetCommonApp()

			//fmt.Println("------------------------ ZABBIX -------------------------")
			//fmt.Println("NetworkHost:", conf.GetCommonApp().Zabbix.NetworkHost)
			//fmt.Println("NetworkPort:", conf.GetCommonApp().Zabbix.NetworkPort)

			Expect(commonApp.Zabbix.NetworkHost).Should(Equal("192.168.9.45"))
			Expect(commonApp.Zabbix.NetworkPort).Should(Equal(10051))
			Expect(commonApp.Zabbix.ZabbixHost).Should(Equal("test-uchet-db.cloud.gcm"))
			Expect(len(commonApp.Zabbix.EventTypes)).Should(Equal(3))
			Expect(commonApp.Zabbix.EventTypes[0].EventType).Should(Equal("error"))
			Expect(commonApp.Zabbix.EventTypes[0].ZabbixKey).Should(Equal("placeholder_misp.error"))
			Expect(commonApp.Zabbix.EventTypes[0].IsTransmit).Should(BeTrue())
			Expect(commonApp.Zabbix.EventTypes[0].Handshake.TimeInterval).Should(Equal(0))
			Expect(commonApp.Zabbix.EventTypes[0].Handshake.Message).Should(Equal(""))

			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(conf.GetListLogs())).Should(Equal(5))
			Expect(len(conf.GetListOrganization())).Should(Equal(12))
			Expect(conf.LogList[0].PathDirectory).Should(Equal("logs"))
			Expect(conf.LogList[0].MsgTypeName).Should(Equal("error"))
			Expect(conf.GetAppTheHive().Send).Should(BeTrue())
		})

		It("Должно быть получено содержимое файла 'config_prod.yaml' при пустом значении переменной GO_PHMISP_MAIN", func() {
			os.Setenv("GO_PHMISP_MAIN", "")

			conf, err := confighandler.NewConfig()

			Expect(err).ShouldNot(HaveOccurred())
			Expect(conf.GetAppNATS().Host).Should(Equal("nats.cloud.gcm"))
			Expect(conf.GetAppNATS().Port).Should(Equal(4222))
			Expect(conf.GetAppMISP().Host).Should(Equal("misp-center.cloud.gcm"))
		})

		It("Должно быть получено содержимое файла 'config_dev.yaml' при значении переменной GO_PHMISP_MAIN=development", func() {
			os.Setenv("GO_PHMISP_MAIN", "development")

			conf, err := confighandler.NewConfig()

			Expect(err).ShouldNot(HaveOccurred())
			Expect(conf.GetAppNATS().Host).Should(Equal("nats.cloud.gcm"))
			Expect(conf.GetAppNATS().Port).Should(Equal(4222))
			Expect(conf.GetAppMISP().Host).Should(Equal("misp-world.cloud.gcm"))

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
			ls := confighandler.Logs{}

			viper.SetConfigFile("../configs/config.yaml")
			viper.SetConfigType("yaml")
			errRC := viper.ReadInConfig()

			Expect(errRC).ShouldNot(HaveOccurred())

			ok := viper.IsSet("LOGGING")
			err := viper.GetViper().Unmarshal(&ls)

			Expect(len(ls.Logging)).ShouldNot(Equal(0))
			fmt.Println("result unmarshal config file LOG: ", ls)

			Expect(ok).Should(BeTrue())
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
