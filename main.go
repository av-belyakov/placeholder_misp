package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/av-belyakov/simplelogger"

	"placeholder_misp/confighandler"
	"placeholder_misp/coremodule"
	"placeholder_misp/datamodels"
	"placeholder_misp/elasticsearchinteractions"
	"placeholder_misp/memorytemporarystorage"
	"placeholder_misp/mispinteractions"
	"placeholder_misp/natsinteractions"
	"placeholder_misp/nkckiinteractions"
	"placeholder_misp/redisinteractions"
	rules "placeholder_misp/rulesinteraction"
)

var (
	err                  error
	sl                   simplelogger.SimpleLoggerSettings
	confApp              confighandler.ConfigApp
	listRulesProcMISPMsg rules.ListRulesProcessingMsgMISP
	listWarning          []string
	storageApp           *memorytemporarystorage.CommonStorageTemporary
	loging               chan datamodels.MessageLoging
)

func getAppName(pf string, nl int) (string, error) {
	var line string

	f, err := os.OpenFile(pf, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return line, err
	}
	defer f.Close()

	num := 1
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		if num == nl {
			return sc.Text(), nil
		}

		num++
	}

	return line, nil
}

func getLoggerSettings(cls []confighandler.LogSet) []simplelogger.MessageTypeSettings {
	loggerConf := make([]simplelogger.MessageTypeSettings, 0, len(cls))

	for _, v := range cls {
		loggerConf = append(loggerConf, simplelogger.MessageTypeSettings{
			MsgTypeName:   v.MsgTypeName,
			WritingFile:   v.WritingFile,
			PathDirectory: v.PathDirectory,
			WritingStdout: v.WritingStdout,
			MaxFileSize:   v.MaxFileSize,
		})
	}

	return loggerConf
}

/*
[]simplelogger.MessageTypeSettings{
		{
			MsgTypeName:   "error",
			WritingFile:   true,
			PathDirectory: "logs",
			WritingStdout: true,
			MaxFileSize:   1024,
		},
		{
			MsgTypeName:   "warning",
			WritingFile:   true,
			PathDirectory: "logs",
			WritingStdout: false,
			MaxFileSize:   1024,
		},
		{
			MsgTypeName:   "info",
			WritingFile:   true,
			PathDirectory: "logs",
			WritingStdout: true,
			MaxFileSize:   1024,
		},
	}
*/

func init() {
	loging = make(chan datamodels.MessageLoging)

	//инициализируем модуль чтения конфигурационного файла
	confApp, err = confighandler.NewConfig()
	if err != nil {
		log.Fatalf("error module 'confighandler': %v", err)
	}

	//инициализируем модуль логирования
	sl, err = simplelogger.NewSimpleLogger("placeholder_misp", getLoggerSettings(confApp.GetListLogs()))
	if err != nil {
		log.Fatalf("error module 'simplelogger': %v", err)
	}

	//инициализируем модуль чтения правил обработки MISP сообщений
	listRulesProcMISPMsg, listWarning, err = rules.GetRuleProcessingMsgForMISP(confApp.RulesProcMSGMISP.Directory, confApp.RulesProcMSGMISP.File)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		_ = sl.WriteLoggingData(fmt.Sprintf("%v %s:%d", err, f, l-2), "error")

		log.Fatalf("error module 'rulesinteraction': %v\n", err)
	}

	//если есть какие либо логические ошибки в файле с YAML правилами для обработки сообщений поступающих от NATS
	if len(listWarning) > 0 {
		var warningStr string

		for _, v := range listWarning {
			warningStr += fmt.Sprintln(v)
		}

		_, f, l, _ := runtime.Caller(0)
		_ = sl.WriteLoggingData(fmt.Sprintf("%s:%d\n%v", f, l, warningStr), "warning")
	}

	// проверяем наличие правил Pass или Passany
	if len(listRulesProcMISPMsg.Rules.Pass) == 0 && !listRulesProcMISPMsg.Rules.Passany {
		msg := "there are no rules for handling messages received from NATS or all rules have failed validation"
		_, f, l, _ := runtime.Caller(0)
		_ = sl.WriteLoggingData(fmt.Sprintf(" '%s' %s:%d", msg, f, l-3), "error")

		log.Fatalf("%s\n", msg)
	}

	//инициализируем модуль временного хранения информации
	storageApp = memorytemporarystorage.NewTemporaryStorage()
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			_ = sl.WriteLoggingData(fmt.Sprintf("stop 'main' function, %v", err), "error")
		}
	}()

	var appName = "placeholder_misp"
	if an, err := getAppName("README.md", 1); err != nil {
		_, f, l, _ := runtime.Caller(0)
		_ = sl.WriteLoggingData(fmt.Sprintf(" '%s' %s:%d", err, f, l-2), "warning")
	} else {
		appName = an
	}

	log.Printf("Application '%s' is start", appName)

	//инициализация модуля для взаимодействия с NATS (Данный модуль обязателен для взаимодействия)
	natsModule, err := natsinteractions.NewClientNATS(confApp.AppConfigNATS, storageApp, loging)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		_ = sl.WriteLoggingData(fmt.Sprintf(" '%s' %s:%d", err, f, l-2), "error")

		log.Fatal(err)
	}

	// инициализация модуля для взаимодействия с СУБД Redis
	ctxRedis := context.Background()
	redisModule := redisinteractions.HandlerRedis(ctxRedis, confApp.AppConfigRedis, storageApp, loging)

	//инициалиация модуля для взаимодействия с MISP
	mispModule, err := mispinteractions.HandlerMISP(confApp.AppConfigMISP, storageApp, loging)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		_ = sl.WriteLoggingData(fmt.Sprintf(" '%s' %s:%d", err, f, l-2), "error")
	}

	//инициализация модуля для взаимодействия с ElasticSearch
	esModule, err := elasticsearchinteractions.HandlerElasticSearch(confApp.AppConfigElasticSearch, storageApp, loging)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		_ = sl.WriteLoggingData(fmt.Sprintf(" '%s' %s:%d", err, f, l-2), "error")
	}

	// инициализация модуля для взаимодействия с NKCKI
	nkckiModule, err := nkckiinteractions.NewClientNKCKI(confApp.AppConfigNKCKI, loging)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		_ = sl.WriteLoggingData(fmt.Sprintf(" '%s' %s:%d", err, f, l-2), "error")
	}

	go func() {
		for msg := range loging {
			_ = sl.WriteLoggingData(msg.MsgData, msg.MsgType)

			/*
				здесь нужно дополнительно сделат
					отправку логов в zabbix
			*/
		}
	}()

	_ = sl.WriteLoggingData("application '"+appName+"' is started", "info")

	coremodule.CoreHandler(natsModule, mispModule, redisModule, esModule, nkckiModule, listRulesProcMISPMsg, storageApp, loging)
}
