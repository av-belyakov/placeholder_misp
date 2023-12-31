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

func NewConfig() (ConfigApp, error) {
	conf := ConfigApp{}
	var envList map[string]string = map[string]string{
		"GO_PHMISP_MAIN":       "",
		"GO_PHMISP_MHOST":      "",
		"GO_PHMISP_MAUTH":      "",
		"GO_PHMISP_NHOST":      "",
		"GO_PHMISP_NPORT":      "",
		"GO_PHMISP_REDISHOST":  "",
		"GO_PHMISP_REDISPORT":  "",
		"GO_PHMISP_ESSEND":     "",
		"GO_PHMISP_ESHOST":     "",
		"GO_PHMISP_ESPORT":     "",
		"GO_PHMISP_ESPREFIX":   "",
		"GO_PHMISP_ESINDEX":    "",
		"GO_PHMISP_ESUSER":     "",
		"GO_PHMISP_ESPASSWD":   "",
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

			ti := 10
			if z.Zabbix.TimeInterval > 0 && z.Zabbix.TimeInterval <= 30 {
				ti = z.Zabbix.TimeInterval
			}

			hs := "0"
			if z.Zabbix.Handshake != "" && len(z.Zabbix.Handshake) <= 60 {
				hs = z.Zabbix.Handshake
			}

			np := 10051
			if z.Zabbix.NetworkPort != 0 && z.Zabbix.NetworkPort < 65536 {
				np = z.Zabbix.NetworkPort
			}

			conf.CommonAppConfig.Zabbix = ZabbixOptions{
				IsTransmit:   z.Zabbix.IsTransmit,
				TimeInterval: ti,
				NetworkPort:  np,
				NetworkHost:  z.Zabbix.NetworkHost,
				ZabbixHost:   z.Zabbix.ZabbixHost,
				ZabbixKey:    z.Zabbix.ZabbixKey,
				Handshake:    hs,
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

		// Настройки для модуля подключения к ElasticSearch
		if viper.IsSet("ElasticSearch.send") {
			conf.AppConfigElasticSearch.Send = viper.GetBool("ElasticSearch.send")
		}
		if viper.IsSet("ElasticSearch.host") {
			conf.AppConfigElasticSearch.Host = viper.GetString("ElasticSearch.host")
		}
		if viper.IsSet("ElasticSearch.port") {
			conf.AppConfigElasticSearch.Port = viper.GetInt("ElasticSearch.port")
		}
		if viper.IsSet("ElasticSearch.prefix") {
			conf.AppConfigElasticSearch.Prefix = viper.GetString("ElasticSearch.prefix")
		}
		if viper.IsSet("ElasticSearch.index") {
			conf.AppConfigElasticSearch.Index = viper.GetString("ElasticSearch.index")
		}
		if viper.IsSet("ElasticSearch.user") {
			conf.AppConfigElasticSearch.User = viper.GetString("ElasticSearch.user")
		}
		if viper.IsSet("ElasticSearch.passwd") {
			conf.AppConfigElasticSearch.Passwd = viper.GetString("ElasticSearch.passwd")
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

	//Настройки для модуля подключения к ElasticSearch
	if envList["GO_PHMISP_ESSEND"] != "" {
		if envList["GO_PHMISP_ESSEND"] == "true" {
			conf.AppConfigElasticSearch.Send = true
		}

		if envList["GO_PHMISP_ESSEND"] == "false" {
			conf.AppConfigElasticSearch.Send = false
		}
	}
	if envList["GO_PHMISP_ESHOST"] != "" {
		conf.AppConfigElasticSearch.Host = envList["GO_PHMISP_ESHOST"]
	}
	if envList["GO_PHMISP_ESPORT"] != "" {
		if p, err := strconv.Atoi(envList["GO_PHMISP_ESPORT"]); err == nil {
			conf.AppConfigElasticSearch.Port = p
		}
	}
	if envList["GO_PHMISP_ESPREFIX"] != "" {
		conf.AppConfigElasticSearch.Prefix = envList["GO_PHMISP_ESPREFIX"]
	}
	if envList["GO_PHMISP_ESINDEX"] != "" {
		conf.AppConfigElasticSearch.Index = envList["GO_PHMISP_ESINDEX"]
	}
	if envList["GO_PHMISP_ESUSER"] != "" {
		conf.AppConfigElasticSearch.User = envList["GO_PHMISP_ESUSER"]
	}
	if envList["GO_PHMISP_ESPASSWD"] != "" {
		conf.AppConfigElasticSearch.Passwd = envList["GO_PHMISP_ESPASSWD"]
	}

	// ПРЕДВАРИТЕЛЬНЫЕ Настройки для модуля подключения к NKCKI
	if envList["GO_PHMISP_NKCKIHOST"] != "" {
		conf.AppConfigNKCKI.Host = envList["GO_PHMISP_NKCKIHOST"]
	}
	if envList["GO_PHMISP_NKCKIPORT"] != "" {
		if p, err := strconv.Atoi(envList["GO_PHMISP_NKCKIPORT"]); err == nil {
			conf.AppConfigNKCKI.Port = p
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
