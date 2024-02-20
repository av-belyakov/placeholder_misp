package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
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

const ROOT_DIR = "placeholder_misp"

var (
	err        error
	warnings   []string
	confApp    confighandler.ConfigApp
	sl         simplelogger.SimpleLoggerSettings
	lr         *rules.ListRule
	hz         *zabbixinteractions.HandlerZabbixConnection
	storageApp *memorytemporarystorage.CommonStorageTemporary

	iz       chan string
	logging  chan datamodels.MessageLogging
	counting chan datamodels.DataCounterSettings
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

// loggingHandler обработчик логов
func loggingHandler(
	iz chan<- string,
	sl simplelogger.SimpleLoggerSettings,
	logging <-chan datamodels.MessageLogging) {
	for msg := range logging {
		_ = sl.WriteLoggingData(msg.MsgData, msg.MsgType)

		if msg.MsgType == "error" || msg.MsgType == "info" {
			iz <- msg.MsgData
		}
	}
}

// counterHandler обработчик счетчиков
func counterHandler(
	iz chan<- string,
	storageApp *memorytemporarystorage.CommonStorageTemporary,
	counting <-chan datamodels.DataCounterSettings) {
	var ae, emr int

	for d := range counting {
		switch d.DataType {
		case "update accepted events":
			storageApp.SetAcceptedEventsDataCounter(d.Count)
		case "update processed events":
			storageApp.SetProcessedEventsDataCounter(d.Count)
		case "update events meet rules":
			storageApp.SetEventsMeetRulesDataCounter(d.Count)
		case "events do not meet rules":
			storageApp.SetEventsDoNotMeetRulesDataCounter(d.Count)
		}

		dc := storageApp.GetDataCounter()
		d, h, m, s := supportingfunctions.GetDifference(dc.StartTime, time.Now())

		log.Printf("\tсобытий принятых/обработанных: %d/%d, соответствие/не соответствие правилам: %d/%d, время со старта приложения: дней %d, часов %d, минут %d, секунд %d\n", dc.AcceptedEvents, dc.ProcessedEvents, dc.EventsMeetRules, dc.EventsDoNotMeetRules, d, h, m, s)

		if ae != dc.AcceptedEvents || emr != dc.EventsMeetRules {
			iz <- fmt.Sprintf("событий принятых: %d, соответствие правилам: %d, время со старта приложения: дней %d, часов %d, минут %d, секунд %d\n", dc.AcceptedEvents, dc.EventsMeetRules, d, h, m, s)

			ae = dc.AcceptedEvents
			emr = dc.EventsMeetRules
		}
	}
}

// interactionZabbix осуществляет взаимодействие с Zabbix
func interactionZabbix(
	confApp confighandler.ConfigApp,
	hz *zabbixinteractions.HandlerZabbixConnection,
	sl simplelogger.SimpleLoggerSettings,
	iz <-chan string) {
	co := confApp.GetCommonApp()
	t := time.Tick(time.Duration(co.Zabbix.TimeInterval) * time.Minute)

	for {
		select {
		case <-t:
			if !co.Zabbix.IsTransmit {
				continue
			}

			if _, err := hz.SendData([]string{co.Zabbix.Handshake}); err != nil {
				_, f, l, _ := runtime.Caller(0)
				_ = sl.WriteLoggingData(fmt.Sprintf(" '%v' %s:%d", err, f, l-1), "error")
			}

		case msg := <-iz:
			if !co.Zabbix.IsTransmit {
				continue
			}

			if _, err := hz.SendData([]string{msg}); err != nil {
				_, f, l, _ := runtime.Caller(0)
				_ = sl.WriteLoggingData(fmt.Sprintf(" '%v' %s:%d", err, f, l-1), "error")
			}
		}
	}
}

func init() {
	iz = make(chan string)
	logging = make(chan datamodels.MessageLogging)
	counting = make(chan datamodels.DataCounterSettings)

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
	lr, warnings, err = rules.NewListRule(ROOT_DIR, confApp.RulesProcMSGMISP.Directory, confApp.RulesProcMSGMISP.File)
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

		log.Fatalf("%s\n", msg)
	}

	//инициализируем модуль временного хранения информации
	storageApp = memorytemporarystorage.NewTemporaryStorage()

	//добавляем время инициализации счетчика хранения
	storageApp.SetStartTimeDataCounter(time.Now())

	commOpt := confApp.GetCommonApp()
	host := fmt.Sprintf("%s:%d", commOpt.Zabbix.NetworkHost, commOpt.Zabbix.NetworkPort)

	//инициализируем модуль связи с Zabbix
	hz = zabbixinteractions.NewHandlerZabbixConnection(host, commOpt.Zabbix.ZabbixHost, commOpt.Zabbix.ZabbixKey)
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			_ = sl.WriteLoggingData(fmt.Sprintf("stop 'main' function, %v", err), "error")
		}
	}()

	var appName string
	appStatus := "production"
	if an, err := getAppName("README.md", 1); err != nil {
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
	log.Printf("Placeholder_misp application, version %s is running. Application status is '%s'\n", appVersion, appStatus)

	//взаимодействие с Zabbix
	go interactionZabbix(confApp, hz, sl, iz)

	//вывод данных счетчика
	go counterHandler(iz, storageApp, counting)

	// логирование данных
	go loggingHandler(iz, sl, logging)

	//инициализация модуля для взаимодействия с NATS (Данный модуль обязателен для взаимодействия)
	natsModule, err := natsinteractions.NewClientNATS(confApp.AppConfigNATS, confApp.AppConfigTheHive, storageApp, logging, counting)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		_ = sl.WriteLoggingData(fmt.Sprintf(" '%s' %s:%d", err, f, l-2), "error")

		log.Fatal(err)
	}

	// инициализация модуля для взаимодействия с СУБД Redis
	ctxRedis := context.Background()
	redisModule := redisinteractions.HandlerRedis(ctxRedis, *confApp.GetAppRedis(), storageApp, logging)

	//инициалиация модуля для взаимодействия с MISP
	mispModule, err := mispinteractions.HandlerMISP(*confApp.GetAppMISP(), confApp.GetListOrganization(), logging)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		_ = sl.WriteLoggingData(fmt.Sprintf(" '%s' %s:%d", err, f, l-2), "error")
	}

	logging <- datamodels.MessageLogging{
		MsgData: "application '" + appName + "' is started",
		MsgType: "info",
	}

	coremodule.CoreHandler(natsModule, mispModule, redisModule, lr, storageApp, logging, counting)
}
