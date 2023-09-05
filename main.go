package main

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/av-belyakov/simplelogger"

	"placeholder_misp/confighandler"
	"placeholder_misp/coremodule"
	"placeholder_misp/datamodels"
	"placeholder_misp/elasticsearchinteractions"
	"placeholder_misp/memorytemporarystorage"
	"placeholder_misp/mispinteractions"
	"placeholder_misp/natsinteractions"
	"placeholder_misp/nkckiinteractions"
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

func init() {
	loging = make(chan datamodels.MessageLoging)

	//инициализируем модуль логирования
	sl, err = simplelogger.NewSimpleLogger("placeholder_misp", []simplelogger.MessageTypeSettings{
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
	})
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		log.Fatalf("error module 'simplelogger': %v %s:%d\n", err, f, l+18)
	}

	//инициализируем модуль чтения конфигурационного файла
	confApp, err = confighandler.NewConfig()
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		_ = sl.WriteLoggingData(fmt.Sprintf("%v %s:%d", err, f, l-2), "error")

		log.Fatalf("error module 'confighandler': %v\n", err)
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
		_ = sl.WriteLoggingData(fmt.Sprintf("%s %s:%d", msg, f, l-3), "error")

		log.Fatalf("%s\n", msg)
	}

	//инициализируем модуль временного хранения информации
	storageApp = memorytemporarystorage.NewTemporaryStorage()
}

func main() {
	log.Println("Application 'placeholder_misp' is start")

	//инициализация модуля для взаимодействия с NATS (Данный модуль обязателен для взаимодействия)
	ctxNATS, ctxCloseNATS := context.WithTimeout(context.Background(), 2*time.Second)
	defer ctxCloseNATS()
	natsModule, err := natsinteractions.NewClientNATS(ctxNATS, confApp.AppConfigNATS, storageApp, loging)
	if err != nil {
		_ = sl.WriteLoggingData(fmt.Sprintln(err), "error")

		log.Fatal(err)
	}

	//инициалиация модуля для взаимодействия с MISP
	ctxMISP, ctxCloseMISP := context.WithTimeout(context.Background(), 2*time.Second)
	defer ctxCloseMISP()
	mispModule, err := mispinteractions.HandlerMISP(ctxMISP, confApp.AppConfigMISP, storageApp, loging)
	if err != nil {
		_ = sl.WriteLoggingData(fmt.Sprintln(err), "error")
	}

	//инициализация модуля для взаимодействия с ElasticSearch
	ctxES, ctxCloseES := context.WithTimeout(context.Background(), 2*time.Second)
	defer ctxCloseES()
	esModule, err := elasticsearchinteractions.HandlerElasticSearch(ctxES, confApp.AppConfigElasticSearch, storageApp, loging)
	if err != nil {
		_ = sl.WriteLoggingData(fmt.Sprintln(err), "error")
	}

	// инициализация модуля для взаимодействия с NKCKI
	ctxNKCKI, ctxCloseNKCKI := context.WithTimeout(context.Background(), 2*time.Second)
	defer ctxCloseNKCKI()
	nkckiModule, err := nkckiinteractions.NewClientNKCKI(ctxNKCKI, confApp.AppConfigNKCKI, loging)
	if err != nil {
		_ = sl.WriteLoggingData(fmt.Sprintln(err), "error")
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

	coremodule.CoreHandler(natsModule, mispModule, esModule, nkckiModule, listRulesProcMISPMsg, storageApp, loging)
}
