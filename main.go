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
	"placeholder_misp/mispinteractions"
	"placeholder_misp/natsinteractions"
	rules "placeholder_misp/rulesinteraction"
)

var (
	err                  error
	sl                   simplelogger.SimpleLoggerSettings
	confApp              confighandler.ConfigApp
	listRulesProcMISPMsg rules.ListRulesProcMISPMessage
	listWarning          []string
	msgOutChan           chan datamodels.MessageLoging
)

func init() {

	// + 1. Прочитать переменные окружения, пока одну
	// + 2. Инициировать модуль для чтения конфигурационных файлов. При этом сначало читается общий конфиг, а затем
	// тот конфиг, выбор которого зависит от переменной окружения 'GO_PH_MISP_MAIN'
	// + 3. Инициализировать обработчик ошибок (запись логов) или отправка их на stdout
	// + 4. Инициализировать модуль соединения с NATS
	//5. Инициализировать модуль соединения с MISP
	//6. Инициализировать модуль обработчик
	//7. Модуль взаимодействия с Забикс? или он в модуле логах должен быть?

	msgOutChan = make(chan datamodels.MessageLoging)

	sl, err = simplelogger.NewSimpleLogger("simplelogger", []simplelogger.MessageTypeSettings{
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

	listRulesProcMISPMsg, listWarning, err = rules.GetRuleProcessedMISPMsg(confApp.RulesProcMsg.Directory, confApp.RulesProcMsg.File)
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

	// проверяем наличие вообще каких либо правил
	if len(listRulesProcMISPMsg.Rulles) == 0 {
		msg := "there are no rules for processing messages received from NATS or have not been verified"
		_, f, l, _ := runtime.Caller(0)
		_ = sl.WriteLoggingData(fmt.Sprintf("%s %s:%d", msg, f, l-3), "error")

		log.Fatalf("%s\n", msg)
	}
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

	natsModule, err := natsinteractions.NewClientNATS(ctxNATS, confApp.AppConfigNATS, msgOutChan)
	if err != nil {
		_ = sl.WriteLoggingData(fmt.Sprintln(err), "error")

		log.Fatal(err)
	}

	ctxMISP, ctxCloseMISP := context.WithTimeout(context.Background(), 2*time.Second)
	defer ctxCloseMISP()

	mispModule, err := mispinteractions.NewClientMISP(ctxMISP, confApp.AppConfigMISP, msgOutChan)
	if err != nil {
		_ = sl.WriteLoggingData(fmt.Sprintln(err), "error")

		log.Fatal()
	}

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

	coremodule.NewCore(*natsModule, *mispModule, msgOutChan)
}
