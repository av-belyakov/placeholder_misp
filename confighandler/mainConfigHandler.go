package confighandler

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type ConfigApp struct {
	CommonAppConfig
	AppConfigNATS
	AppConfigMISP
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

func NewConfig() (ConfigApp, error) {
	conf := ConfigApp{}
	var envList map[string]string = map[string]string{
		"GO_PHMISP_MAIN":  "",
		"GO_PHMISP_MHOST": "",
		"GO_PHMISP_MPORT": "",
		"GO_PHMISP_NHOST": "",
		"GO_PHMISP_NPORT": "",
	}

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

		if viper.IsSet("NATS.host") {
			conf.AppConfigNATS.Host = viper.GetString("NATS.host")
		}

		if viper.IsSet("NATS.port") {
			conf.AppConfigNATS.Port = viper.GetInt("NATS.port")
		}

		if viper.IsSet("MISP.host") {
			conf.AppConfigMISP.Host = viper.GetString("MISP.host")
		}

		if viper.IsSet("MISP.port") {
			conf.AppConfigMISP.Port = viper.GetInt("MISP.port")
		}

		return nil
	}

	for v := range envList {
		if env, ok := os.LookupEnv(v); ok {
			envList[v] = env
		}
	}

	rootPath, err := getRootPath("placeholder_misp")
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

	if envList["GO_PHMISP_MHOST"] != "" {
		conf.AppConfigMISP.Host = envList["GO_PHMISP_MHOST"]
	}

	if envList["GO_PHMISP_MPORT"] != "" {
		if p, err := strconv.Atoi(envList["GO_PHMISP_MPORT"]); err == nil {
			conf.AppConfigMISP.Port = p
		}
	}

	if envList["GO_PHMISP_NHOST"] != "" {
		conf.AppConfigNATS.Host = envList["GO_PHMISP_NHOST"]
	}

	if envList["GO_PHMISP_NPORT"] != "" {
		if p, err := strconv.Atoi(envList["GO_PHMISP_NPORT"]); err == nil {
			conf.AppConfigNATS.Port = p
		}
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
