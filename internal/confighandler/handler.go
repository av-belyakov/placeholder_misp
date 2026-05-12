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
	conf := &ConfigApp{}

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

			conf.CommonCfg.LogList = ls.Logging
		}

		orgs := Orgs{}
		if ok := viper.IsSet("ORGANIZATIONS"); ok {
			if err := viper.GetViper().Unmarshal(&orgs); err != nil {
				return err
			}

			conf.CommonCfg.Organizations = orgs.Organizations
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

			conf.CommonCfg.Zabbix = ZabbixOptions{
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
			conf.CfgNATS.Host = viper.GetString("NATS.host")
		}
		if viper.IsSet("NATS.port") {
			conf.CfgNATS.Port = viper.GetInt("NATS.port")
		}
		if viper.IsSet("NATS.cache_ttl") {
			conf.CfgNATS.CacheTTL = viper.GetInt("NATS.cache_ttl")
		}
		if viper.IsSet("NATS.subscriptions.listener_case") {
			conf.CfgNATS.Subscriptions.ListenerCase = viper.GetString("NATS.subscriptions.listener_case")
		}
		if viper.IsSet("NATS.subscriptions.sender_command") {
			conf.CfgNATS.Subscriptions.SenderCommand = viper.GetString("NATS.subscriptions.sender_command")
		}

		//Настройки для модуля подключения к MISP
		if viper.IsSet("MISP.host") {
			conf.CfgMISP.Host = viper.GetString("MISP.host")
		}

		//Настройки для модуля подключения к Sqlite3
		if viper.IsSet("SQLITE3.path_file_db") {
			conf.CfgSqlite3.PathFileDb = viper.GetString("SQLITE3.path_file_db")
		}

		//Настройки для взаимодействия с TheHive
		if viper.IsSet("THEHIVE.send") {
			conf.CfgTheHive.Send = viper.GetBool("THEHIVE.send")
		}

		//Настройки для модуля правил обработки сообщений
		if viper.IsSet("RULES_PROC_MSG_FOR_MISP.directory") {
			conf.RulesProcMSGMISP.Directory = viper.GetString("RULES_PROC_MSG_FOR_MISP.directory")
		}
		if viper.IsSet("RULES_PROC_MSG_FOR_MISP.file") {
			conf.RulesProcMSGMISP.File = viper.GetString("RULES_PROC_MSG_FOR_MISP.file")
		}

		// Настройки доступа к БД в которую будут записыватся логи
		if viper.IsSet("DATABASEWRITELOG.host") {
			conf.CfgWriteLogDB.Host = viper.GetString("DATABASEWRITELOG.host")
		}
		if viper.IsSet("DATABASEWRITELOG.port") {
			conf.CfgWriteLogDB.Port = viper.GetInt("DATABASEWRITELOG.port")
		}
		if viper.IsSet("DATABASEWRITELOG.user") {
			conf.CfgWriteLogDB.User = viper.GetString("DATABASEWRITELOG.user")
		}
		if viper.IsSet("DATABASEWRITELOG.namedb") {
			conf.CfgWriteLogDB.NameDB = viper.GetString("DATABASEWRITELOG.namedb")
		}
		if viper.IsSet("DATABASEWRITELOG.storage_name_db") {
			conf.CfgWriteLogDB.StorageNameDB = viper.GetString("DATABASEWRITELOG.storage_name_db")
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
		return conf, err
	}

	confPath := filepath.Join(rootPath, "config")
	list, err := os.ReadDir(confPath)
	if err != nil {
		return conf, err
	}

	fileNameCommon, err := getFileName("config.yml", confPath, list)
	if err != nil {
		return conf, err
	}

	//читаем общий конфигурационный файл
	if err := setCommonSettings(fileNameCommon); err != nil {
		return conf, err
	}

	var fn string
	if envList["GO_PHMISP_MAIN"] == "development" {
		fn, err = getFileName("config_dev.yml", confPath, list)
		if err != nil {
			return conf, err
		}
	} else {
		fn, err = getFileName("config_prod.yml", confPath, list)
		if err != nil {
			return conf, err
		}
	}

	if err := setSpecial(fn); err != nil {
		return conf, err
	}

	//Настройки для модуля подключения к NATS
	if envList["GO_PHMISP_NHOST"] != "" {
		conf.CfgNATS.Host = envList["GO_PHMISP_NHOST"]
	}
	if envList["GO_PHMISP_NPORT"] != "" {
		if p, err := strconv.Atoi(envList["GO_PHMISP_NPORT"]); err == nil {
			conf.CfgNATS.Port = p
		}
	}
	if envList["GO_PHMISP_NCACHETTL"] != "" {
		if ttl, err := strconv.Atoi(envList["GO_PHMISP_NCACHETTL"]); err == nil {
			conf.CfgNATS.CacheTTL = ttl
		}
	}
	if envList["GO_PHMISP_NSUBLISTENERCASE"] != "" {
		conf.CfgNATS.Subscriptions.ListenerCase = envList["GO_PHMISP_NSUBLISTENERCASE"]
	}
	if envList["GO_PHMISP_NSUBSENDERCOMMAND"] != "" {
		conf.CfgNATS.Subscriptions.SenderCommand = envList["GO_PHMISP_NSUBSENDERCOMMAND"]
	}

	//Настройки для модуля подключения к MISP
	if envList["GO_PHMISP_MHOST"] != "" {
		conf.CfgMISP.Host = envList["GO_PHMISP_MHOST"]
	}
	if envList["GO_PHMISP_MAUTH"] != "" {
		conf.CfgMISP.Auth = envList["GO_PHMISP_MAUTH"]
	}

	//Настройки для модуля подключения к Sqlite3
	if envList["GO_PHMISP_SQLITE3PATH"] != "" {
		conf.CfgSqlite3.PathFileDb = envList["GO_PHMISP_SQLITE3PATH"]
	}

	//Настройки для модуля правил обработки сообщений
	if envList["GO_PHMISP_RULES_DIR"] != "" {
		conf.RulesProcMSGMISP.Directory = envList["GO_PHMISP_RULES_DIR"]
	}
	if envList["GO_PHMISP_RULES_FILE"] != "" {
		conf.RulesProcMSGMISP.File = envList["GO_PHMISP_RULES_FILE"]
	}

	//Настройки доступа к БД в которую будут записыватся логи
	if envList["GO_PHMISP_DBWLOGHOST"] != "" {
		conf.CfgWriteLogDB.Host = envList["GO_PHMISP_DBWLOGHOST"]
	}
	if envList["GO_PHMISP_DBWLOGPORT"] != "" {
		if p, err := strconv.Atoi(envList["GO_PHMISP_DBWLOGPORT"]); err == nil {
			conf.CfgWriteLogDB.Port = p
		}
	}
	if envList["GO_PHMISP_DBWLOGNAME"] != "" {
		conf.CfgWriteLogDB.NameDB = envList["GO_PHMISP_DBWLOGNAME"]
	}
	if envList["GO_PHMISP_DBWLOGUSER"] != "" {
		conf.CfgWriteLogDB.User = envList["GO_PHMISP_DBWLOGUSER"]
	}
	if envList["GO_PHMISP_DBWLOGPASSWD"] != "" {
		conf.CfgWriteLogDB.Passwd = envList["GO_PHMISP_DBWLOGPASSWD"]
	}
	if envList["GO_PHMISP_DBWLOGSTORAGENAME"] != "" {
		conf.CfgWriteLogDB.StorageNameDB = envList["GO_PHMISP_DBWLOGSTORAGENAME"]
	}

	//выполняем проверку заполненой структуры
	if err = validate.Struct(conf); err != nil {
		return conf, err
	}

	return conf, nil
}
