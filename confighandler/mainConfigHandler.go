package confighandler

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strconv"

	"github.com/spf13/viper"

	"placeholder_misp/supportingfunctions"
)

type ConfigApp struct {
	CommonAppConfig
	AppConfigNATS
	AppConfigMISP
	AppConfigElasticSearch
	AppConfigNKCKI
	RulesProcMSGMISP
}

type CommonAppConfig struct {
	LogList []LogSet
}

type Logs struct {
	Logging []LogSet
}

type LogSet struct {
	MsgTypeName   string
	WritingFile   bool
	PathDirectory string
	WritingStdout bool
	MaxFileSize   int
}

type AppConfigNATS struct {
	Host string
	Port int
}

type AppConfigMISP struct {
	Host string
	Port int
}

type AppConfigElasticSearch struct {
	Host string
	Port int
}

type AppConfigNKCKI struct {
	Host string
	Port int
}

type RulesProcMSGMISP struct {
	Directory, File string
}

func NewConfig() (ConfigApp, error) {
	conf := ConfigApp{}
	var envList map[string]string = map[string]string{
		"GO_PHMISP_MAIN":       "",
		"GO_PHMISP_MHOST":      "",
		"GO_PHMISP_MPORT":      "",
		"GO_PHMISP_NHOST":      "",
		"GO_PHMISP_NPORT":      "",
		"GO_PHMISP_ESHOST":     "",
		"GO_PHMISP_ESPORT":     "",
		"GO_PHMISP_NKCKIHOST":  "",
		"GO_PHMISP_NKCKIPORT":  "",
		"GO_PHMISP_RULES_DIR":  "",
		"GO_PHMISP_RULES_FILE": "",
	}

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

		//Настройки для модуля подключения к MISP
		if viper.IsSet("MISP.host") {
			conf.AppConfigMISP.Host = viper.GetString("MISP.host")
		}
		if viper.IsSet("MISP.port") {
			conf.AppConfigMISP.Port = viper.GetInt("MISP.port")
		}

		// ПРЕДВАРИТЕЛЬНЫЕ Настройки для модуля подключения к ElasticSearch
		if viper.IsSet("ElasticSearch.host") {
			conf.AppConfigElasticSearch.Host = viper.GetString("ElasticSearch.host")
		}
		if viper.IsSet("ElasticSearch.port") {
			conf.AppConfigElasticSearch.Port = viper.GetInt("ElasticSearch.port")
		}

		// ПРЕДВАРИТЕЛЬНЫЕ Настройки для модуля подключения к NKCKI
		if viper.IsSet("NKCKI.host") {
			conf.AppConfigNKCKI.Host = viper.GetString("NKCKI.host")
		}
		if viper.IsSet("NKCKI.port") {
			conf.AppConfigNKCKI.Port = viper.GetInt("NKCKI.port")
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

	for v := range envList {
		if env, ok := os.LookupEnv(v); ok {
			envList[v] = env
		}
	}

	rootPath, err := supportingfunctions.GetRootPath("placeholder_misp")
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

	//Настройки для модуля подключения к MISP
	if envList["GO_PHMISP_MHOST"] != "" {
		conf.AppConfigMISP.Host = envList["GO_PHMISP_MHOST"]
	}
	if envList["GO_PHMISP_MPORT"] != "" {
		if p, err := strconv.Atoi(envList["GO_PHMISP_MPORT"]); err == nil {
			conf.AppConfigMISP.Port = p
		}
	}

	// ПРЕДВАРИТЕЛЬНЫЕ Настройки для модуля подключения к ElasticSearch
	if envList["GO_PHMISP_ESHOST"] != "" {
		conf.AppConfigMISP.Host = envList["GO_PHMISP_ESHOST"]
	}
	if envList["GO_PHMISP_ESPORT"] != "" {
		if p, err := strconv.Atoi(envList["GO_PHMISP_ESPORT"]); err == nil {
			conf.AppConfigMISP.Port = p
		}
	}

	// ПРЕДВАРИТЕЛЬНЫЕ Настройки для модуля подключения к NKCKI
	if envList["GO_PHMISP_NKCKIHOST"] != "" {
		conf.AppConfigMISP.Host = envList["GO_PHMISP_NKCKIHOST"]
	}
	if envList["GO_PHMISP_NKCKIPORT"] != "" {
		if p, err := strconv.Atoi(envList["GO_PHMISP_NKCKIPORT"]); err == nil {
			conf.AppConfigMISP.Port = p
		}
	}

	//Настройки для модуля правил обработки сообщений
	if envList["GO_PHMISP_RULES_DIR"] != "" {
		conf.RulesProcMSGMISP.Directory = envList["GO_PHMISP_RULES_DIR"]
	}
	if envList["GO_PHMISP_RULES_FILE"] != "" {
		conf.RulesProcMSGMISP.File = envList["GO_PHMISP_RULES_FILE"]
	}

	return conf, nil
}

func (conf *ConfigApp) GetCommonApp() *CommonAppConfig {
	return &conf.CommonAppConfig
}

func (conf *ConfigApp) GetAppNATS() *AppConfigNATS {
	return &conf.AppConfigNATS
}

func (conf *ConfigApp) GetAppMISP() *AppConfigMISP {
	return &conf.AppConfigMISP
}

func (conf *ConfigApp) Clean() {
	conf = &ConfigApp{}
}
