package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/av-belyakov/simplelogger"

	"placeholder_misp/confighandler"
	"placeholder_misp/coremodule"
	"placeholder_misp/datamodels"
	"placeholder_misp/memorytemporarystorage"
	"placeholder_misp/mispinteractions"
	"placeholder_misp/natsinteractions"
	"placeholder_misp/redisinteractions"
	rules "placeholder_misp/rulesinteraction"
	"placeholder_misp/supportingfunctions"
	"placeholder_misp/zabbixinteractions"
)

const (
	ROOT_DIR = "placeholder_misp"

	ansiWhite               = "\033[97m"
	ansiDarkGreenBackground = "\033[42m"
	boldFont                = "\033[1m"
	ansiReset               = "\033[0m"
)

func getLoggerSettings(cls []confighandler.LogSet) []simplelogger.Options {
	loggerConf := make([]simplelogger.Options, 0, len(cls))

	for _, v := range cls {
		loggerConf = append(loggerConf, simplelogger.Options{
			MsgTypeName:     v.MsgTypeName,
			WritingToFile:   v.WritingFile,
			PathDirectory:   v.PathDirectory,
			WritingToStdout: v.WritingStdout,
			MaxFileSize:     v.MaxFileSize,
		})
	}

	return loggerConf
}

// loggingHandler обработчик логов
func loggingHandler(
	channelZabbix chan<- zabbixinteractions.MessageSettings,
	sl *simplelogger.SimpleLoggerSettings,
	logging <-chan datamodels.MessageLogging) {
	for msg := range logging {
		_ = sl.WriteLoggingData(msg.MsgData, msg.MsgType)

		if msg.MsgType == "error" || msg.MsgType == "warning" {
			channelZabbix <- zabbixinteractions.MessageSettings{
				EventType: "error",
				Message:   fmt.Sprintf("%s: %s", msg.MsgType, msg.MsgData),
			}
		}

		if msg.MsgType == "info" {
			channelZabbix <- zabbixinteractions.MessageSettings{
				EventType: "info",
				Message:   msg.MsgData,
			}
		}
	}
}

// counterHandler обработчик счетчиков
func counterHandler(
	channelZabbix chan<- zabbixinteractions.MessageSettings,
	storageApp *memorytemporarystorage.CommonStorageTemporary,
	sl *simplelogger.SimpleLoggerSettings,
	counting <-chan datamodels.DataCounterSettings) {
	for data := range counting {
		d, h, m, s := supportingfunctions.GetDifference(storageApp.GetStartTimeDataCounter(), time.Now())
		patternTime := fmt.Sprintf("со старта приложения: дней %d, часов %d, минут %d, секунд %d", d, h, m, s)
		var msg string

		switch data.DataType {
		case "update accepted events":
			storageApp.SetAcceptedEventsDataCounter(data.Count)
			msg = fmt.Sprintf("принято: %d, %s", storageApp.GetAcceptedEventsDataCounter(), patternTime)
		case "update processed events":
			storageApp.SetProcessedEventsDataCounter(data.Count)
			msg = fmt.Sprintf("обработано: %d, %s", storageApp.GetProcessedEventsDataCounter(), patternTime)
		case "update events meet rules":
			storageApp.SetEventsMeetRulesDataCounter(data.Count)
			msg = fmt.Sprintf("соответствует правилам: %d, %s", storageApp.GetEventsMeetRulesDataCounter(), patternTime)
		}

		_ = sl.WriteLoggingData(msg, "debug")

		channelZabbix <- zabbixinteractions.MessageSettings{
			EventType: "info",
			Message:   msg,
		}
	}
}

// interactionZabbix осуществляет взаимодействие с Zabbix
func interactionZabbix(
	ctx context.Context,
	confApp *confighandler.ConfigApp,
	sl *simplelogger.SimpleLoggerSettings,
	channelZabbix <-chan zabbixinteractions.MessageSettings) error {

	connTimeout := time.Duration(7 * time.Second)
	hz, err := zabbixinteractions.NewZabbixConnection(
		ctx,
		zabbixinteractions.SettingsZabbixConnection{
			Port:              confApp.Zabbix.NetworkPort,
			Host:              confApp.Zabbix.NetworkHost,
			NetProto:          "tcp",
			ZabbixHost:        confApp.Zabbix.ZabbixHost,
			ConnectionTimeout: &connTimeout,
		})
	if err != nil {
		return err
	}

	et := make([]zabbixinteractions.EventType, len(confApp.Zabbix.EventTypes))
	for _, v := range confApp.Zabbix.EventTypes {
		et = append(et, zabbixinteractions.EventType{
			IsTransmit: v.IsTransmit,
			EventType:  v.EventType,
			ZabbixKey:  v.ZabbixKey,
			Handshake:  zabbixinteractions.Handshake(v.Handshake),
		})
	}

	if err = hz.Handler(et, channelZabbix); err != nil {
		return err
	}

	go func() {
		for err := range hz.GetChanErr() {
			_, f, l, _ := runtime.Caller(0)
			_ = sl.WriteLoggingData(fmt.Sprintf("zabbix module: '%s' %s:%d", err.Error(), f, l-1), "error")
		}
	}()

	return nil
}

func main() {
	ctx, ctxCancel := signal.NotifyContext(context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	//инициализируем модуль чтения конфигурационного файла
	confApp, err := confighandler.NewConfig(ROOT_DIR)
	if err != nil {
		log.Fatalf("error module 'confighandler': %v", err)
	}

	//инициализируем модуль логирования
	sl, err := simplelogger.NewSimpleLogger(ctx, ROOT_DIR, getLoggerSettings(confApp.GetListLogs()))
	if err != nil {
		log.Fatalf("error module 'simplelogger': %v", err)
	}

	defer func() {
		if r := recover(); r != nil {
			_ = sl.WriteLoggingData(fmt.Sprint(r), "error")
		}
	}()

	//инициализируем модуль чтения правил обработки MISP сообщений
	lr, warnings, err := rules.NewListRule(ROOT_DIR, confApp.RulesProcMSGMISP.Directory, confApp.RulesProcMSGMISP.File)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		_ = sl.WriteLoggingData(fmt.Sprintf("%v %s:%d", err, f, l-2), "error")

		log.Fatalf("error module 'rulesinteraction': %v\n", err)
	}

	//если есть какие либо логические ошибки в файле с YAML правилами для обработки сообщений поступающих от NATS
	if len(warnings) > 0 {
		var warningStr string

		for _, v := range warnings {
			warningStr += fmt.Sprintln(v)
		}

		_, f, l, _ := runtime.Caller(0)
		_ = sl.WriteLoggingData(fmt.Sprintf("%s:%d\n%v", f, l, warningStr), "warning")
	}

	// проверяем наличие правил Pass или Passany
	if len(lr.GetRulePass()) == 0 && !lr.GetRulePassany() {
		msg := "there are no rules for handling messages received from NATS or all rules have failed validation"
		_, f, l, _ := runtime.Caller(0)
		_ = sl.WriteLoggingData(fmt.Sprintf(" '%s' %s:%d", msg, f, l-3), "error")

		log.Fatalln(msg)
	}

	//взаимодействие с Zabbix
	channelZabbix := make(chan zabbixinteractions.MessageSettings)
	if err := interactionZabbix(ctx, confApp, sl, channelZabbix); err != nil {
		_, f, l, _ := runtime.Caller(0)
		_ = sl.WriteLoggingData(fmt.Sprintf(" '%s' %s:%d", err.Error(), f, l-3), "error")

		log.Fatalln(err.Error())
	}

	var appName string
	appStatus := "production"
	if an, err := supportingfunctions.GetAppName("README.md", 1); err != nil {
		_, f, l, _ := runtime.Caller(0)
		_ = sl.WriteLoggingData(fmt.Sprintf(" '%s' %s:%d", err, f, l-2), "warning")
	} else {
		appName = an
	}

	envValue, ok := os.LookupEnv("GO_PHMISP_MAIN")
	if ok && envValue == "development" {
		appStatus = envValue
	}

	appVersion := supportingfunctions.GetAppVersion(appName)

	msg := fmt.Sprintf("Placeholder_misp application, version %s is running. Application status is '%s'", appVersion, appStatus)
	log.Printf("%v%v%v%s%v\n", ansiDarkGreenBackground, boldFont, ansiWhite, msg, ansiReset)

	//инициализируем модуль временного хранения информации
	storageApp := memorytemporarystorage.NewTemporaryStorage()

	//добавляем время инициализации счетчика хранения
	storageApp.SetStartTimeDataCounter(time.Now())

	//вывод данных счетчика
	counting := make(chan datamodels.DataCounterSettings)
	go counterHandler(channelZabbix, storageApp, sl, counting)

	//логирование данных
	logging := make(chan datamodels.MessageLogging)
	go loggingHandler(channelZabbix, sl, logging)

	//инициализация модуля для взаимодействия с NATS (Данный модуль обязателен для взаимодействия)
	natsModule, err := natsinteractions.NewClientNATS(confApp.AppConfigNATS, confApp.AppConfigTheHive, storageApp, logging, counting)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		_ = sl.WriteLoggingData(fmt.Sprintf(" '%s' %s:%d", err, f, l-2), "error")

		log.Fatal(err)
	}

	//инициализация модуля для взаимодействия с СУБД Redis
	ctxRedis := context.Background()
	redisModule := redisinteractions.HandlerRedis(ctxRedis, *confApp.GetAppRedis(), storageApp, logging)

	//инициалиация модуля для взаимодействия с MISP
	mispModule, err := mispinteractions.HandlerMISP(*confApp.GetAppMISP(), confApp.GetListOrganization(), logging)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		_ = sl.WriteLoggingData(fmt.Sprintf(" '%s' %s:%d", err, f, l-2), "error")
	}

	//выполняется для записи в файл и отправки в Zabbix
	logging <- datamodels.MessageLogging{
		MsgData: "application '" + appName + "' is started",
		MsgType: "info",
	}

	go func() {
		sigChan := make(chan os.Signal, 1)
		osCall := <-sigChan
		log.Printf("system call:%+v", osCall)

		close(counting)
		close(logging)
		close(channelZabbix)

		ctxCancel()
	}()

	core := coremodule.NewCoreHandler(storageApp, logging, counting)
	core.CoreHandler(ctx, natsModule, mispModule, redisModule, lr)
}
