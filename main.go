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
	msgOutChan           chan datamodels.MessageLoging
)

func init() {
	msgOutChan = make(chan datamodels.MessageLoging)

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

	confApp, err = confighandler.NewConfig()
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		_ = sl.WriteLoggingData(fmt.Sprintf("%v %s:%d", err, f, l-2), "error")

		log.Fatalf("error module 'confighandler': %v\n", err)
	}

	fmt.Println("func 'main', config application = ", confApp)

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

	storageApp = memorytemporarystorage.NewTemporaryStorage()
}

func main() {
	fmt.Println("func 'main', START...")
	fmt.Println("config application:", confApp)

	log.Println("Application 'placeholder_misp' is start")

	// i, err := os.Stdout.Write([]byte("test writing to stdout"))
	// fmt.Println("os.Stdout.Write i = ", i, " error = ", err)
	//_ = sl.WriteLoggingData("my test message about trable", "error")
	//_ = sl.WriteLoggingData("my test information message", "info")

	ctxNATS, ctxCloseNATS := context.WithTimeout(context.Background(), 2*time.Second)
	defer ctxCloseNATS()

	//инициализация модуля для взаимодействия с NATS (Данный модуль обязателен для взаимодействия)
	natsModule, err := natsinteractions.NewClientNATS(ctxNATS, confApp.AppConfigNATS, msgOutChan)
	if err != nil {
		_ = sl.WriteLoggingData(fmt.Sprintln(err), "error")

		log.Fatal(err)
	}

	//инициалиация модуля для взаимодействия с MISP
	ctxMISP, ctxCloseMISP := context.WithTimeout(context.Background(), 2*time.Second)
	defer ctxCloseMISP()

	mispModule, err := mispinteractions.NewClientMISP(ctxMISP, confApp.AppConfigMISP, msgOutChan)
	if err != nil {
		_ = sl.WriteLoggingData(fmt.Sprintln(err), "error")
	}

	//инициализация модуля для взаимодействия с ElasticSearch
	ctxES, ctxCloseES := context.WithTimeout(context.Background(), 2*time.Second)
	defer ctxCloseES()

	esModule, err := elasticsearchinteractions.NewClientElasticSearch(ctxES, confApp.AppConfigElasticSearch, msgOutChan)
	if err != nil {
		_ = sl.WriteLoggingData(fmt.Sprintln(err), "error")
	}

	// инициализация модуля для взаимодействия с NKCKI
	ctxNKCKI, ctxCloseNKCKI := context.WithTimeout(context.Background(), 2*time.Second)
	defer ctxCloseNKCKI()

	nkckiModule, err := nkckiinteractions.NewClientNKCKI(ctxNKCKI, confApp.AppConfigNKCKI, msgOutChan)
	if err != nil {
		_ = sl.WriteLoggingData(fmt.Sprintln(err), "error")
	}

	//Если подключений к API MISP, NKCKI, ElasticSearch нет следует продолжить выполнения программы
	//а в дальнейшем пытаться установить соединение с данными API

	/*
			ОБЯЗАТЕЛЬНО К ПРОЧТЕНИЮ
		Появилась задача осуществлять подключения к Elasticsearch и к НКЦКИ. И отправлять туда данные
		из Hive. В ответ с MISP, Elasticsearch и НКЦКИ должны приходить ID принятого сообщения, который
		нужно отправлять назад в Hive.

		здесь сделать инициализацию подключения к zabbix
		вообще правильнее основное логирование сделать здесь, а все ошибки кидать
		из модулей через канал сюда. Тогда можно выполнять логирование и отправлять
		логи ошибок через модуль подключения к zabbix
	*/

	go func() {
		for msg := range msgOutChan {
			_ = sl.WriteLoggingData(msg.MsgData, msg.MsgType)

			/*
				отправлять логи в zabbix
			*/
		}
	}()

	coremodule.NewCore(*natsModule, *mispModule, *esModule, *nkckiModule, listRulesProcMISPMsg, storageApp, msgOutChan)
}
