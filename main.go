package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/av-belyakov/simplelogger"

	"placeholder_misp/confighandler"
	"placeholder_misp/coremodule"
	"placeholder_misp/mispinteractions"
	"placeholder_misp/natsinteractions"
)

var (
	err     error
	sl      simplelogger.SimpleLoggerSettings
	confApp confighandler.ConfigApp
	//ctxNATS, ctxMISP           context.Context
	//ctxCloseNATS, ctxCloseMISP context.CancelFunc
)

func init() {

	// + 1. Прочитать переменные окружения, пока одну
	// + 2. Инициировать модуль для чтения конфигурационных файлов. При этом сначало читается общий конфиг, а затем
	// тот конфиг, выбор которого зависит от переменной окружения 'GO_PH_MISP_MAIN'
	// + 3. Инициализировать обработчик ошибок (запись логов) или отправка их на stdout
	//4. Инициализировать модуль соединения с NATS
	//5. Инициализировать модуль соединения с MISP
	//6. Инициализировать модуль обработчик
	//7. Модуль взаимодействия с Забикс? или он в модуле логах должен быть?

	confApp, err = confighandler.NewConfig()
	if err != nil {
		log.Fatalf("error module 'confighandler': %v\n", err)
	}

	sl, err = simplelogger.NewSimpleLogger("simplelogger", []simplelogger.MessageTypeSettings{
		{
			MsgTypeName:   "error",
			WritingFile:   true,
			PathDirectory: "logs",
			WritingStdout: true,
			MaxFileSize:   1024,
		},
		{
			MsgTypeName:   "info",
			WritingFile:   true,
			PathDirectory: "logs",
			//WritingStdout: false,
			WritingStdout: true,
			MaxFileSize:   1024,
		},
	})
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

	natsModule, err := natsinteractions.NewClientNATS(ctxNATS, confApp.AppConfigNATS)
	if err != nil {
		_ = sl.WriteLoggingData(fmt.Sprintln(err), "error")
	}

	ctxMISP, ctxCloseMISP := context.WithTimeout(context.Background(), 2*time.Second)
	defer ctxCloseMISP()

	mispModule, err := mispinteractions.NewClientMISP(ctxMISP, confApp.AppConfigMISP)
	if err != nil {
		_ = sl.WriteLoggingData(fmt.Sprintln(err), "error")
	}

	coremodule.NewCore(*natsModule, *mispModule, sl)
}
