// Пакет confighandler формирует конфигурационные настройки приложения
package confighandler

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"

	"placeholder_misp/supportingfunctions"
)

func NewConfig(rootDir string) (*ConfigApp, error) {
	conf := &ConfigApp{}

	var (
		validate *validator.Validate
		envList  map[string]string = map[string]string{
			"GO_PHMISP_MAIN": "",

			//Подключение к MISP
			"GO_PHMISP_MHOST": "",
			"GO_PHMISP_MAUTH": "",

			//Подключение к NATS
			"GO_PHMISP_NHOST":               "",
			"GO_PHMISP_NPORT":               "",
			"GO_PHMISP_NCACHETTL":           "",
			"GO_PHMISP_NSUBSENDERCASE":      "",
			"GO_PHMISP_NSUBLISTENERCOMMAND": "",

			//Подключение к Redis DB
			"GO_PHMISP_REDISHOST": "",
			"GO_PHMISP_REDISPORT": "",

			//Правила обработки событий
			"GO_PHMISP_RULES_DIR":  "",
			"GO_PHMISP_RULES_FILE": "",
		}
	)

	getFileName := func(sf, confPath string, lfs []fs.DirEntry) (string, error) {
		for _, v := range lfs {
			if v.Name() == sf && !v.IsDir() {
				return path.Join(confPath, v.Name()), nil
			}
		}

		return "", fmt.Errorf("file '%s' is not found", sf)
	}

	setCommonSettings := func(fn string) error {
		viper.SetConfigFile(fn)
		viper.SetConfigType("yaml")
		if err := viper.ReadInConfig(); err != nil {
			return err
		}

		ls := Logs{}
		if ok := viper.IsSet("LOGGING"); ok {
			if err := viper.GetViper().Unmarshal(&ls); err != nil {
				return err
			}

			conf.CommonAppConfig.LogList = ls.Logging
		}

		orgs := Orgs{}
		if ok := viper.IsSet("ORGANIZATIONS"); ok {
			if err := viper.GetViper().Unmarshal(&orgs); err != nil {
				return err
			}

			conf.CommonAppConfig.Organizations = orgs.Organizations
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

			conf.CommonAppConfig.Zabbix = ZabbixOptions{
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
		viper.SetConfigType("yaml")
		if err := viper.ReadInConfig(); err != nil {
			return err
		}

		//Настройки для модуля подключения к NATS
		if viper.IsSet("NATS.host") {
			conf.AppConfigNATS.Host = viper.GetString("NATS.host")
		}
		if viper.IsSet("NATS.port") {
			conf.AppConfigNATS.Port = viper.GetInt("NATS.port")
		}
		if viper.IsSet("NATS.cacheTtl") {
			conf.AppConfigNATS.CacheTTL = viper.GetInt("NATS.cacheTtl")
		}

		if viper.IsSet("NATS.subscriptions.sender_case") {
			conf.AppConfigNATS.Subscriptions.SenderCase = viper.GetString("NATS.subscriptions.sender_case")
		}
		if viper.IsSet("NATS.subscriptions.listener_command") {
			conf.AppConfigNATS.Subscriptions.ListenerCommand = viper.GetString("NATS.subscriptions.listener_command")
		}

		//Настройки для модуля подключения к MISP
		if viper.IsSet("MISP.host") {
			conf.AppConfigMISP.Host = viper.GetString("MISP.host")
		}
		if viper.IsSet("MISP.auth") {
			conf.AppConfigMISP.Auth = viper.GetString("MISP.auth")
		}

		//Настройки для модуля подключения к СУБД Redis
		if viper.IsSet("REDIS.host") {
			conf.AppConfigRedis.Host = viper.GetString("REDIS.host")
		}
		if viper.IsSet("REDIS.port") {
			conf.AppConfigRedis.Port = viper.GetInt("REDIS.port")
		}

		//Настройки для взаимодействия с TheHive
		if viper.IsSet("THEHIVE.send") {
			conf.AppConfigTheHive.Send = viper.GetBool("THEHIVE.send")
		}

		//Настройки для модуля правил обработки сообщений
		if viper.IsSet("RULES_PROC_MSG_FOR_MISP.directory") {
			conf.RulesProcMSGMISP.Directory = viper.GetString("RULES_PROC_MSG_FOR_MISP.directory")
		}
		if viper.IsSet("RULES_PROC_MSG_FOR_MISP.file") {
			conf.RulesProcMSGMISP.File = viper.GetString("RULES_PROC_MSG_FOR_MISP.file")
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

	confPath := path.Join(rootPath, "configs")

	list, err := os.ReadDir(confPath)
	if err != nil {
		return conf, err
	}

	fileNameCommon, err := getFileName("config.yaml", confPath, list)
	if err != nil {
		return conf, err
	}

	//читаем общий конфигурационный файл
	if err := setCommonSettings(fileNameCommon); err != nil {
		return conf, err
	}

	var fn string
	if envList["GO_PHMISP_MAIN"] == "development" {
		fn, err = getFileName("config_dev.yaml", confPath, list)
		if err != nil {
			return conf, err
		}
	} else {
		fn, err = getFileName("config_prod.yaml", confPath, list)
		if err != nil {
			return conf, err
		}
	}

	if err := setSpecial(fn); err != nil {
		return conf, err
	}

	//Настройки для модуля подключения к NATS
	if envList["GO_PHMISP_NHOST"] != "" {
		conf.AppConfigNATS.Host = envList["GO_PHMISP_NHOST"]
	}
	if envList["GO_PHMISP_NPORT"] != "" {
		if p, err := strconv.Atoi(envList["GO_PHMISP_NPORT"]); err == nil {
			conf.AppConfigNATS.Port = p
		}
	}
	if envList["GO_PHMISP_NCACHETTL"] != "" {
		if v, err := strconv.Atoi(envList["GO_PHMISP_NCACHETTL"]); err == nil {
			conf.AppConfigNATS.CacheTTL = v
		}
	}
	if envList["GO_PHMISP_NSUBSENDERCASE"] != "" {
		conf.AppConfigNATS.Subscriptions.SenderCase = envList["GO_PHMISP_NSUBSENDERCASE"]
	}
	if envList["GO_PHMISP_NSUBLISTENERCOMMAND"] != "" {
		conf.AppConfigNATS.Subscriptions.ListenerCommand = envList["GO_PHMISP_NSUBLISTENERCOMMAND"]
	}

	//Настройки для модуля подключения к MISP
	if envList["GO_PHMISP_MHOST"] != "" {
		conf.AppConfigMISP.Host = envList["GO_PHMISP_MHOST"]
	}
	if envList["GO_PHMISP_MAUTH"] != "" {
		conf.AppConfigMISP.Auth = envList["GO_PHMISP_MAUTH"]
	}

	//Настройки для модуля подключения к СУБД Redis
	if envList["GO_PHMISP_REDISHOST"] != "" {
		conf.AppConfigRedis.Host = envList["GO_PHMISP_REDISHOST"]
	}
	if envList["GO_PHMISP_REDISPORT"] != "" {
		if p, err := strconv.Atoi(envList["GO_PHMISP_REDISPORT"]); err == nil {
			conf.AppConfigRedis.Port = p
		}
	}

	//Настройки для модуля правил обработки сообщений
	if envList["GO_PHMISP_RULES_DIR"] != "" {
		conf.RulesProcMSGMISP.Directory = envList["GO_PHMISP_RULES_DIR"]
	}
	if envList["GO_PHMISP_RULES_FILE"] != "" {
		conf.RulesProcMSGMISP.File = envList["GO_PHMISP_RULES_FILE"]
	}

	//выполняем проверку заполненой структуры
	if err = validate.Struct(conf); err != nil {
		return conf, err
	}

	return conf, nil
}
