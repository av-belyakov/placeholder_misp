package confighandler

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/av-belyakov/placeholder_misp/constants"
	"github.com/av-belyakov/placeholder_misp/internal/confighandler"
	"github.com/av-belyakov/placeholder_misp/internal/supportingfunctions"
)

var _ = Describe("MainConfigHandler", Ordered, func() {
	var (
		rootPath string

		err error
	)

	BeforeAll(func() {
		//настройки MISP
		os.Unsetenv("GO_PHMISP_MAIN")
		os.Unsetenv("GO_PHMISP_MHOST")
		os.Unsetenv("GO_PHMISP_MPORT")

		//настройки NATS
		os.Unsetenv("GO_PHMISP_NHOST")
		os.Unsetenv("GO_PHMISP_NPORT")
		os.Unsetenv("GO_PHMISP_NCACHETTL")
		os.Unsetenv("GO_PHMISP_NSUBLISTENERCASE")
		os.Unsetenv("GO_PHMISP_NSUBSENDERCOMMAND")

		//настройки доступа к БД в которую будут записыватся логи
		os.Unsetenv("GO_PHMISP_DBWLOGHOST")
		os.Unsetenv("GO_PHMISP_DBWLOGPORT")
		os.Unsetenv("GO_PHMISP_DBWLOGNAME")
		os.Unsetenv("GO_PHMISP_DBWLOGUSER")
		os.Unsetenv("GO_PHMISP_DBWLOGSTORAGENAME")

		rootPath, err = supportingfunctions.GetRootPath(constants.Root_Dir)
		if err != nil {
			log.Fatalf("error, it is impossible to form root path (%s)", err.Error())
		}

		if err = godotenv.Load("../../.env"); err != nil {
			log.Fatalln(err)
		}
	})

	AfterAll(func() {
		os.Unsetenv("GO_PHMISP_MAUTH")
		os.Unsetenv("GO_PHMISP_DBWLOGPASSWD")
	})

	Context("Тест 1. Проверяем работу функции NewConfig с разными значениями переменной окружения GO_PHMISP_MAIN", func() {
		It("Должно быть получено содержимое общего файла 'config.yml'", func() {
			// инициализируем модуль чтения конфигурационного файла
			conf, err := confighandler.New(rootPath)
			Expect(err).ShouldNot(HaveOccurred())

			//fmt.Println("conf = ", conf)

			//for k, v := range conf.GetListOrganization() {
			//	fmt.Printf("%d. OrgName: %s, SourceName: %s\n", k, v.OrgName, v.SourceName)
			//}

			commonApp := conf.GetCommonApp()

			fmt.Println("------------------------ application config NATS -------------------------")
			fmt.Printf("%+v\n", conf.AppConfigNATS)

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

		It("Должно быть получено содержимое файла 'config_prod.yml' при пустом значении переменной GO_PHMISP_MAIN", func() {
			os.Setenv("GO_PHMISP_MAIN", "")

			// инициализируем модуль чтения конфигурационного файла
			conf, err := confighandler.New(rootPath)
			Expect(err).ShouldNot(HaveOccurred())

			confNats := conf.GetAppNATS()
			Expect(confNats.Host).Should(Equal("192.168.9.208"))
			Expect(confNats.Port).Should(Equal(4222))
			Expect(confNats.CacheTTL).Should(Equal(3600))
			Expect(confNats.Subscriptions.ListenerCase).Should(Equal("object.casetype.local"))
			Expect(confNats.Subscriptions.SenderCommand).Should(Equal("object.commandstype"))

			//параметры подключения к MISP
			Expect(conf.GetAppMISP().Host).Should(Equal("misp-center.cloud.gcm"))

			//параметры подключения к Sqlite3
			Expect(conf.GetAppSqlite3().PathFileDb).Should(Equal("/sqlite3/sqlite3.db"))

			//параметры БД для логирования
			confLoggingDB := conf.GetApplicationWriteLogDB()
			Expect(confLoggingDB.Host).Should(Equal("datahook.cloud.gcm"))
			Expect(confLoggingDB.Port).Should(Equal(9200))
			Expect(confLoggingDB.NameDB).Should(Equal(""))
			Expect(confLoggingDB.StorageNameDB).Should(Equal("placeholder_misp"))
			Expect(confLoggingDB.User).Should(Equal("log_writer"))
			Expect(len(confLoggingDB.Passwd)).ShouldNot(Equal(0))
		})

		It("Должно быть получено содержимое файла 'config_dev.yml' при значении переменной GO_PHMISP_MAIN=development", func() {
			os.Setenv("GO_PHMISP_MAIN", "development")

			conf, err := confighandler.New(rootPath)
			Expect(err).ShouldNot(HaveOccurred())

			confNats := conf.GetAppNATS()
			Expect(confNats.Host).Should(Equal("nats.cloud.gcm"))
			Expect(confNats.Port).Should(Equal(4222))
			Expect(confNats.CacheTTL).Should(Equal(3600))
			Expect(confNats.Subscriptions.ListenerCase).Should(Equal("object.casetype.local"))
			Expect(confNats.Subscriptions.SenderCommand).Should(Equal("object.commandstype"))

			//параметры подключения к MISP
			Expect(conf.GetAppMISP().Host).Should(Equal("misp-world.cloud.gcm"))

			//параметры подключения к Sqlite3
			Expect(conf.GetAppSqlite3().PathFileDb).Should(Equal("/test/sqlite3_file/sqlite3.db"))

			//параметры БД для логирования
			confLoggingDB := conf.GetApplicationWriteLogDB()
			Expect(confLoggingDB.Host).Should(Equal("datahook.cloud.gcm"))
			Expect(confLoggingDB.Port).Should(Equal(9200))
			Expect(confLoggingDB.NameDB).Should(Equal(""))
			Expect(confLoggingDB.StorageNameDB).Should(Equal("placeholder_misp"))
			Expect(confLoggingDB.User).Should(Equal("log_writer"))
			Expect(len(confLoggingDB.Passwd)).ShouldNot(Equal(0))
		})

		It("Должно быть перезаписанно содержимое некоторых полей полученных с переменных окружения GO_PHMISP_*", func() {
			os.Setenv("GO_PHMISP_MAIN", "development")

			os.Setenv("GO_PHMISP_NHOST", "78.87.78.87")
			os.Setenv("GO_PHMISP_NPORT", "11111")
			os.Setenv("GO_PHMISP_NCACHETTL", "4500")
			os.Setenv("GO_PHMISP_NSUBLISTENERCASE", "object.casetype.test")
			os.Setenv("GO_PHMISP_NSUBSENDERCOMMAND", "object.commandstype.test")
			os.Setenv("GO_PHMISP_SQLITE3PATH", "somepath/path/file_db")

			conf, err := confighandler.New(rootPath)
			Expect(err).ShouldNot(HaveOccurred())

			confNats := conf.GetAppNATS()
			Expect(confNats.Host).Should(Equal("78.87.78.87"))
			Expect(confNats.Port).Should(Equal(11111))
			Expect(confNats.CacheTTL).Should(Equal(4500))
			Expect(confNats.Subscriptions.ListenerCase).Should(Equal("object.casetype.test"))
			Expect(confNats.Subscriptions.SenderCommand).Should(Equal("object.commandstype.test"))

			Expect(conf.AppConfigSqlite3.PathFileDb).Should(Equal("somepath/path/file_db"))
		})
	})

	Context("Тест 2. Проверяем установленные для DATABASEWRITELOG значения переменных окружения", func() {
		const (
			PHMISP_DBWLOGHOST        = "45.10.32.1"
			PHMISP_DBWLOGPORT        = 11123
			PHMISP_DBWLOGNAME        = "log_db"
			PHMISP_DBWLOGUSER        = "nreuser"
			PHMISP_DBWLOGPASSWD      = "pass123wd"
			PHMISP_DBWLOGSTORAGENAME = "thehivehookgolog"
		)

		BeforeAll(func() {
			os.Setenv("GO_PHMISP_DBWLOGHOST", PHMISP_DBWLOGHOST)
			os.Setenv("GO_PHMISP_DBWLOGPORT", strconv.Itoa(PHMISP_DBWLOGPORT))
			os.Setenv("GO_PHMISP_DBWLOGNAME", PHMISP_DBWLOGNAME)
			os.Setenv("GO_PHMISP_DBWLOGUSER", PHMISP_DBWLOGUSER)
			os.Setenv("GO_PHMISP_DBWLOGPASSWD", PHMISP_DBWLOGPASSWD)
			os.Setenv("GO_PHMISP_DBWLOGSTORAGENAME", PHMISP_DBWLOGSTORAGENAME)

		})

		It("Все параметры конфигурации для WEBHOOKSERVER должны быть успешно установлены через соответствующие переменные окружения", func() {
			conf, err := confighandler.New(rootPath)
			Expect(err).ShouldNot(HaveOccurred())

			wldb := conf.GetApplicationWriteLogDB()

			Expect(wldb.Host).Should(Equal(PHMISP_DBWLOGHOST))
			Expect(wldb.Port).Should(Equal(PHMISP_DBWLOGPORT))
			Expect(wldb.NameDB).Should(Equal(PHMISP_DBWLOGNAME))
			Expect(wldb.User).Should(Equal(PHMISP_DBWLOGUSER))
			Expect(wldb.Passwd).Should(Equal(PHMISP_DBWLOGPASSWD))
			Expect(wldb.StorageNameDB).Should(Equal(PHMISP_DBWLOGSTORAGENAME))
		})
	})

	Context("Тест 3. Тестируем функцию возвращающую путь до текущего приложения", func() {
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
})
