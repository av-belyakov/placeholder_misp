// Пакет confighandler формирует конфигурационные настройки приложения
package confighandler

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"

	"github.com/av-belyakov/placeholder_misp/internal/supportingfunctions"
)

func New(rootDir string) (*ConfigApp, error) {
	cfg := &ConfigApp{}

	var (
		validate *validator.Validate
		envList  map[string]string = map[string]string{
			"GO_PHMISP_MAIN": "",

			//Подключение к MISP
			"GO_PHMISP_MHOST": "",
			"GO_PHMISP_MAUTH": "",

			//Подключение к NATS
			"GO_PHMISP_NHOST":             "",
			"GO_PHMISP_NPORT":             "",
			"GO_PHMISP_NCACHETTL":         "",
			"GO_PHMISP_NSUBLISTENERCASE":  "",
			"GO_PHMISP_NSUBSENDERCOMMAND": "",
			"GO_PHMISP_NSUBGETSENSORINFO": "",

			//Подключение к Sqlite3
			"GO_PHMISP_SQLITE3PATH": "",

			//Правила обработки событий
			"GO_PHMISP_RULES_DIR":  "",
			"GO_PHMISP_RULES_FILE": "",

			//Настройки доступа к БД в которую будут записыватся логи
			"GO_PHMISP_DBWLOGHOST":        "",
			"GO_PHMISP_DBWLOGPORT":        "",
			"GO_PHMISP_DBWLOGNAME":        "",
			"GO_PHMISP_DBWLOGUSER":        "",
			"GO_PHMISP_DBWLOGPASSWD":      "",
			"GO_PHMISP_DBWLOGSTORAGENAME": "",
		}
	)

	getFileName := func(sf, confPath string, lfs []fs.DirEntry) (string, error) {
		for _, v := range lfs {
			if v.Name() == sf && !v.IsDir() {
				return filepath.Join(confPath, v.Name()), nil
			}
		}

		return "", fmt.Errorf("file '%s' is not found", sf)
	}

	setCommonSettings := func(fn string) error {
		viper.SetConfigFile(fn)
		viper.SetConfigType("yml")
		if err := viper.ReadInConfig(); err != nil {
			return err
		}

		ls := Logs{}
		if ok := viper.IsSet("LOGGING"); ok {
			if err := viper.GetViper().Unmarshal(&ls); err != nil {
				return err
			}

			cfg.Common.LogList = ls.Logging
		}

		orgs := Orgs{}
		if ok := viper.IsSet("ORGANIZATIONS"); ok {
			if err := viper.GetViper().Unmarshal(&orgs); err != nil {
				return err
			}

			cfg.Common.Organizations = orgs.Organizations
		}

		z := ZabbixSet{}
		if ok := viper.IsSet("ZABBIX"); ok {
			if err := viper.GetViper().Unmarshal(&z); err != nil {
				return err
			}

			np := 10051
			if z.Zabbix.NetworkPort != 0 && z.Zabbix.NetworkPort < 65536 {
				np = z.Zabbix.NetworkPort
			}

			cfg.Common.Zabbix = ZabbixOptions{
				NetworkPort: np,
				NetworkHost: z.Zabbix.NetworkHost,
				ZabbixHost:  z.Zabbix.ZabbixHost,
				EventTypes:  z.Zabbix.EventTypes,
			}
		}

		return nil
	}

	setSpecial := func(fn string) error {
		viper.SetConfigFile(fn)
		viper.SetConfigType("yml")
		if err := viper.ReadInConfig(); err != nil {
			return err
		}

		//Настройки для модуля подключения к NATS
		if viper.IsSet("NATS.host") {
			cfg.NATS.Host = viper.GetString("NATS.host")
		}
		if viper.IsSet("NATS.port") {
			cfg.NATS.Port = viper.GetInt("NATS.port")
		}
		if viper.IsSet("NATS.cache_ttl") {
			cfg.NATS.CacheTTL = viper.GetInt("NATS.cache_ttl")
		}
		if viper.IsSet("NATS.subscriptions.listener_case") {
			cfg.NATS.Subscriptions.ListenerCase = viper.GetString("NATS.subscriptions.listener_case")
		}
		if viper.IsSet("NATS.subscriptions.sender_command") {
			cfg.NATS.Subscriptions.SenderCommand = viper.GetString("NATS.subscriptions.sender_command")
		}
		if viper.IsSet("NATS.subscriptions.get_sensor_info") {
			cfg.NATS.Subscriptions.GetSensorInfo = viper.GetString("NATS.subscriptions.get_sensor_info")
		}

		//Настройки для модуля подключения к MISP
		if viper.IsSet("MISP.host") {
			cfg.MISP.Host = viper.GetString("MISP.host")
		}

		//Настройки для модуля подключения к Sqlite3
		if viper.IsSet("SQLITE3.path_file_db") {
			cfg.Sqlite3.PathFileDb = viper.GetString("SQLITE3.path_file_db")
		}

		//Настройки для взаимодействия с TheHive
		if viper.IsSet("THEHIVE.send") {
			cfg.TheHive.Send = viper.GetBool("THEHIVE.send")
		}

		//Настройки для модуля правил обработки сообщений
		if viper.IsSet("RULES_PROC_MSG_FOR_MISP.directory") {
			cfg.Rules.Directory = viper.GetString("RULES_PROC_MSG_FOR_MISP.directory")
		}
		if viper.IsSet("RULES_PROC_MSG_FOR_MISP.file") {
			cfg.Rules.File = viper.GetString("RULES_PROC_MSG_FOR_MISP.file")
		}

		// Настройки доступа к БД в которую будут записыватся логи
		if viper.IsSet("DATABASEWRITELOG.host") {
			cfg.WriteLogDB.Host = viper.GetString("DATABASEWRITELOG.host")
		}
		if viper.IsSet("DATABASEWRITELOG.port") {
			cfg.WriteLogDB.Port = viper.GetInt("DATABASEWRITELOG.port")
		}
		if viper.IsSet("DATABASEWRITELOG.user") {
			cfg.WriteLogDB.User = viper.GetString("DATABASEWRITELOG.user")
		}
		if viper.IsSet("DATABASEWRITELOG.namedb") {
			cfg.WriteLogDB.NameDB = viper.GetString("DATABASEWRITELOG.namedb")
		}
		if viper.IsSet("DATABASEWRITELOG.storage_name_db") {
			cfg.WriteLogDB.StorageNameDB = viper.GetString("DATABASEWRITELOG.storage_name_db")
		}

		// Настройки для отладочного сервера
		if viper.IsSet("DebugServer.enable") {
			cfg.DebugServer.Enable = viper.GetBool("DebugServer.enable")
		}
		if viper.IsSet("DebugServer.host") {
			cfg.DebugServer.Host = viper.GetString("DebugServer.host")
		}
		if viper.IsSet("DebugServer.port") {
			cfg.DebugServer.Port = viper.GetInt("DebugServer.port")
		}

		return nil
	}

	validate = validator.New(validator.WithRequiredStructEnabled())

	for v := range envList {
		if env, ok := os.LookupEnv(v); ok {
			envList[v] = env
		}
	}

	rootPath, err := supportingfunctions.GetRootPath(rootDir)
	if err != nil {
		return cfg, err
	}

	confPath := filepath.Join(rootPath, "config")
	list, err := os.ReadDir(confPath)
	if err != nil {
		return cfg, err
	}

	fileNameCommon, err := getFileName("config.yml", confPath, list)
	if err != nil {
		return cfg, err
	}

	//читаем общий конфигурационный файл
	if err := setCommonSettings(fileNameCommon); err != nil {
		return cfg, err
	}

	var fn string
	if envList["GO_PHMISP_MAIN"] == "development" {
		fn, err = getFileName("config_dev.yml", confPath, list)
		if err != nil {
			return cfg, err
		}
	} else {
		fn, err = getFileName("config_prod.yml", confPath, list)
		if err != nil {
			return cfg, err
		}
	}

	if err := setSpecial(fn); err != nil {
		return cfg, err
	}

	//Настройки для модуля подключения к NATS
	if envList["GO_PHMISP_NHOST"] != "" {
		cfg.NATS.Host = envList["GO_PHMISP_NHOST"]
	}
	if envList["GO_PHMISP_NPORT"] != "" {
		if p, err := strconv.Atoi(envList["GO_PHMISP_NPORT"]); err == nil {
			cfg.NATS.Port = p
		}
	}
	if envList["GO_PHMISP_NCACHETTL"] != "" {
		if ttl, err := strconv.Atoi(envList["GO_PHMISP_NCACHETTL"]); err == nil {
			cfg.NATS.CacheTTL = ttl
		}
	}
	if envList["GO_PHMISP_NSUBLISTENERCASE"] != "" {
		cfg.NATS.Subscriptions.ListenerCase = envList["GO_PHMISP_NSUBLISTENERCASE"]
	}
	if envList["GO_PHMISP_NSUBSENDERCOMMAND"] != "" {
		cfg.NATS.Subscriptions.SenderCommand = envList["GO_PHMISP_NSUBSENDERCOMMAND"]
	}
	if envList["GO_PHMISP_NSUBGETSENSORINFO"] != "" {
		cfg.NATS.Subscriptions.GetSensorInfo = envList["GO_PHMISP_NSUBGETSENSORINFO"]
	}

	//Настройки для модуля подключения к MISP
	if envList["GO_PHMISP_MHOST"] != "" {
		cfg.MISP.Host = envList["GO_PHMISP_MHOST"]
	}
	if envList["GO_PHMISP_MAUTH"] != "" {
		cfg.MISP.Auth = envList["GO_PHMISP_MAUTH"]
	}

	//Настройки для модуля подключения к Sqlite3
	if envList["GO_PHMISP_SQLITE3PATH"] != "" {
		cfg.Sqlite3.PathFileDb = envList["GO_PHMISP_SQLITE3PATH"]
	}

	//Настройки для модуля правил обработки сообщений
	if envList["GO_PHMISP_RULES_DIR"] != "" {
		cfg.Rules.Directory = envList["GO_PHMISP_RULES_DIR"]
	}
	if envList["GO_PHMISP_RULES_FILE"] != "" {
		cfg.Rules.File = envList["GO_PHMISP_RULES_FILE"]
	}

	//Настройки доступа к БД в которую будут записыватся логи
	if envList["GO_PHMISP_DBWLOGHOST"] != "" {
		cfg.WriteLogDB.Host = envList["GO_PHMISP_DBWLOGHOST"]
	}
	if envList["GO_PHMISP_DBWLOGPORT"] != "" {
		if p, err := strconv.Atoi(envList["GO_PHMISP_DBWLOGPORT"]); err == nil {
			cfg.WriteLogDB.Port = p
		}
	}
	if envList["GO_PHMISP_DBWLOGNAME"] != "" {
		cfg.WriteLogDB.NameDB = envList["GO_PHMISP_DBWLOGNAME"]
	}
	if envList["GO_PHMISP_DBWLOGUSER"] != "" {
		cfg.WriteLogDB.User = envList["GO_PHMISP_DBWLOGUSER"]
	}
	if envList["GO_PHMISP_DBWLOGPASSWD"] != "" {
		cfg.WriteLogDB.Passwd = envList["GO_PHMISP_DBWLOGPASSWD"]
	}
	if envList["GO_PHMISP_DBWLOGSTORAGENAME"] != "" {
		cfg.WriteLogDB.StorageNameDB = envList["GO_PHMISP_DBWLOGSTORAGENAME"]
	}

	//выполняем проверку заполненой структуры
	if err = validate.Struct(cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
