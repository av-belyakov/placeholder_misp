package confighandler

import (
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	"github.com/av-belyakov/placeholder_misp/constants"
	"github.com/av-belyakov/placeholder_misp/internal/confighandler"
)

var (
	cfg *confighandler.ConfigApp

	err error
)

func TestMain(m *testing.M) {
	unsetEnvVariables()

	//загружаем ключи и пароли
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalln(err)
	}

	os.Setenv("GO_PHMISP_MAIN", "development")

	cfg, err = confighandler.New(constants.Root_Dir)
	if err != nil {
		log.Fatalln(err)
	}

	os.Exit(m.Run())
}

func TestConfigHandler(t *testing.T) {
	t.Run("Тест чтения конфигурационного файла config_dev.yml", func(t *testing.T) {
		t.Run("Тест 1. Проверка настройки параметров доступа к правилам", func(t *testing.T) {
			assert.Equal(t, cfg.GetRules().Directory, "rules")
			assert.Equal(t, cfg.GetRules().File, "mispmsgrule.yml")
		})

		t.Run("Тест 2. Проверка настройки параметра доступа к SQLite3", func(t *testing.T) {
			assert.Equal(t, cfg.GetSqlite3().PathFileDb, "/sqlite3/sqlite3.db")
		})

		t.Run("Тест 3. Проверка настройки NATS", func(t *testing.T) {
			assert.Equal(t, cfg.GetNATS().Host, "192.168.9.208")
			assert.Equal(t, cfg.GetNATS().Port, 4222)
			assert.Equal(t, cfg.GetNATS().CacheTTL, 3600)
			assert.Equal(t, cfg.GetNATS().Subscriptions.ListenerCase, "object.casetype.test")
			assert.Equal(t, cfg.GetNATS().Subscriptions.SenderCommand, "object.commandstype.local")
			assert.Equal(t, cfg.GetNATS().Subscriptions.GetSensorInfo, "object.sensor-info-request")
		})

		t.Run("Тест 4. Проверка настройки MISP", func(t *testing.T) {
			assert.Equal(t, cfg.GetMISP().Host, "misp-world.cloud.gcm")
			assert.NotEmpty(t, cfg.GetMISP().Auth)
		})

		t.Run("Тест 5. Проверка настройки TheHive", func(t *testing.T) {
			assert.True(t, cfg.GetTheHive().Send)
		})

		t.Run("Тест 5. Проверка настройки WriteLogDataBase", func(t *testing.T) {
			assert.Equal(t, cfg.GetLogDB().Host, "datahook.cloud.gcm")
			assert.Equal(t, cfg.GetLogDB().Port, 9200)
			assert.Equal(t, cfg.GetLogDB().User, "log_writer")
			assert.Equal(t, cfg.GetLogDB().NameDB, "")
			assert.Equal(t, cfg.GetLogDB().StorageNameDB, "placeholder_misp")
			assert.NotEmpty(t, cfg.GetLogDB().Passwd)
		})

		t.Run("Тест 6. Проверка настройки сервера отладки", func(t *testing.T) {
			assert.True(t, cfg.GetDebugServer().Enable)
			assert.Equal(t, cfg.GetDebugServer().Host, "localhost")
			assert.Equal(t, cfg.GetDebugServer().Port, 6464)
		})
	})

	t.Run("Тест чтения переменных окружения", func(t *testing.T) {
		t.Run("Тест 1. Проверка настройки MISP", func(t *testing.T) {
			host := "misp-example.gcm"
			auth := "1234567890"

			os.Setenv("GO_PHMISP_MHOST", host)
			os.Setenv("GO_PHMISP_MAUTH", auth)

			cfg, err := confighandler.New(constants.Root_Dir)
			assert.NoError(t, err)

			assert.Equal(t, cfg.GetMISP().Host, host)
			assert.Equal(t, cfg.GetMISP().Auth, auth)
		})
		t.Run("Тест 2. Проверка настройки NATS", func(t *testing.T) {
			host := "127.0.0.1"
			port := "4242"
			cacheTTL := "650"
			subListenerCase := "obj.subscript.any_test_request"
			subSenderCommand := "obj.subscript.any_test_command"
			subGetSensorInfo := "obj.subscript.any_test_sensor_info"

			os.Setenv("GO_PHMISP_NHOST", host)
			os.Setenv("GO_PHMISP_NPORT", port)
			os.Setenv("GO_PHMISP_NCACHETTL", cacheTTL)
			os.Setenv("GO_PHMISP_NSUBLISTENERCASE", subListenerCase)
			os.Setenv("GO_PHMISP_NSUBSENDERCOMMAND", subSenderCommand)
			os.Setenv("GO_PHMISP_NSUBGETSENSORINFO", subGetSensorInfo)

			cfg, err := confighandler.New(constants.Root_Dir)
			assert.NoError(t, err)

			var intPort int
			intPort, err = strconv.Atoi(port)
			assert.NoError(t, err)

			var intCacheTTL int
			intCacheTTL, err = strconv.Atoi(cacheTTL)
			assert.NoError(t, err)

			assert.Equal(t, cfg.GetNATS().Host, host)
			assert.Equal(t, cfg.GetNATS().Port, intPort)
			assert.Equal(t, cfg.GetNATS().CacheTTL, intCacheTTL)
			assert.Equal(t, cfg.GetNATS().Subscriptions.ListenerCase, subListenerCase)
			assert.Equal(t, cfg.GetNATS().Subscriptions.SenderCommand, subSenderCommand)
			assert.Equal(t, cfg.GetNATS().Subscriptions.GetSensorInfo, subGetSensorInfo)
		})

		t.Run("Тест 3. Проверка настройки доступа к Sqlite3", func(t *testing.T) {
			dbPath := "sqlite3.db"

			//настройки доступа к Sqlite3
			os.Setenv("GO_PHMISP_SQLITE3PATH", dbPath)

			cfg, err := confighandler.New(constants.Root_Dir)
			assert.NoError(t, err)

			assert.Equal(t, cfg.GetSqlite3().PathFileDb, dbPath)
		})

		t.Run("Тест 4. Проверка настройки доступа к параметрам правил", func(t *testing.T) {
			rDir := "any-folder/rules"
			rFile := "rules.json"

			//настройки доступа к параметрам правил
			os.Setenv("GO_PHMISP_RULES_DIR", rDir)
			os.Setenv("GO_PHMISP_RULES_FILE", rFile)

			cfg, err := confighandler.New(constants.Root_Dir)
			assert.NoError(t, err)

			assert.Equal(t, cfg.GetRules().Directory, rDir)
			assert.Equal(t, cfg.GetRules().File, rFile)

		})

		t.Run("Тест 3. Проверка настройки WriteLogDataBase", func(t *testing.T) {
			host := "domaniname.database.cm"
			port := "8989"
			user := "somebody_user"
			name := "any_name"
			passwd := "your_passwd"
			nameDB := "any_name_db"

			os.Setenv("GO_PHMISP_DBWLOGHOST", host)
			os.Setenv("GO_PHMISP_DBWLOGPORT", port)
			os.Setenv("GO_PHMISP_DBWLOGUSER", user)
			os.Setenv("GO_PHMISP_DBWLOGNAME", name)
			os.Setenv("GO_PHMISP_DBWLOGPASSWD", passwd)
			os.Setenv("GO_PHMISP_DBWLOGSTORAGENAME", nameDB)

			cfg, err := confighandler.New(constants.Root_Dir)
			assert.NoError(t, err)

			var intPort int
			intPort, err = strconv.Atoi(port)
			assert.NoError(t, err)

			assert.Equal(t, cfg.GetLogDB().Host, host)
			assert.Equal(t, cfg.GetLogDB().Port, intPort)
			assert.Equal(t, cfg.GetLogDB().User, user)
			assert.Equal(t, cfg.GetLogDB().Passwd, passwd)
			assert.Equal(t, cfg.GetLogDB().NameDB, name)
			assert.Equal(t, cfg.GetLogDB().StorageNameDB, nameDB)
		})
	})

	t.Cleanup(func() {
		unsetEnvVariables()
	})
}

func unsetEnvVariables() {
	os.Unsetenv("GO_PHMISP_MAIN")

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
	os.Unsetenv("GO_PHMISP_NSUBGETSENSORINFO")

	//настройки доступа к Sqlite3
	os.Unsetenv("GO_PHMISP_SQLITE3PATH")

	//настройки доступа к параметрам правил
	os.Unsetenv("GO_PHMISP_RULES_DIR")
	os.Unsetenv("GO_PHMISP_RULES_FILE")

	//настройки доступа к БД в которую будут записыватся логи
	os.Unsetenv("GO_PHMISP_DBWLOGHOST")
	os.Unsetenv("GO_PHMISP_DBWLOGPORT")
	os.Unsetenv("GO_PHMISP_DBWLOGNAME")
	os.Unsetenv("GO_PHMISP_DBWLOGUSER")
	os.Unsetenv("GO_PHMISP_DBWLOGSTORAGENAME")
}
