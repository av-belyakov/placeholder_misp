package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/av-belyakov/placeholder_misp/cmd/coremodule"
	"github.com/av-belyakov/placeholder_misp/cmd/mispapi"
	"github.com/av-belyakov/placeholder_misp/cmd/natsapi"
	"github.com/av-belyakov/placeholder_misp/cmd/redisapi"
	"github.com/av-belyakov/placeholder_misp/cmd/wrappers"
	"github.com/av-belyakov/placeholder_misp/commoninterfaces"
	"github.com/av-belyakov/placeholder_misp/constants"
	"github.com/av-belyakov/placeholder_misp/internal/appname"
	"github.com/av-belyakov/placeholder_misp/internal/appversion"
	"github.com/av-belyakov/placeholder_misp/internal/confighandler"
	"github.com/av-belyakov/placeholder_misp/internal/logginghandler"
	"github.com/av-belyakov/placeholder_misp/internal/supportingfunctions"
	"github.com/av-belyakov/placeholder_misp/memorytemporarystorage"
	rules "github.com/av-belyakov/placeholder_misp/rulesinteraction"
	"github.com/av-belyakov/simplelogger"
)

func server(ctx context.Context) {
	rootPath, err := supportingfunctions.GetRootPath(constants.Root_Dir)
	if err != nil {
		log.Fatalf("error, it is impossible to form root path (%s)", err.Error())
	}

	// инициализируем модуль чтения конфигурационного файла
	confApp, err := confighandler.New(rootPath, constants.Conf_Dir)
	if err != nil {
		log.Fatalf("error module 'confighandler': %v", err)
	}

	// ****************************************************************************
	// ********************* инициализация модуля логирования *********************
	simpleLogger, err := simplelogger.NewSimpleLogger(ctx, constants.Root_Dir, getLoggerSettings(confApp.GetListLogs()))
	if err != nil {
		log.Fatalf("error module 'simplelogger': %v", err)
	}

	// ****************************************************************************
	// ******* инициализируем модуль чтения правил обработки MISP сообщений *******
	lr, warnings, err := rules.NewListRule(constants.Root_Dir, confApp.RulesProcMSGMISP.Directory, confApp.RulesProcMSGMISP.File)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		_ = simpleLogger.Write("error", fmt.Sprintf("%v %s:%d", err, f, l-2))

		log.Fatalf("error module 'rulesinteraction': %v\n", err)
	}

	//проверяем наличие правил Pass или Passany которые являются обязательными,
	//а также отсутсвие логических ошибок в файле с правилами
	msgWarning, err := checkListRule(lr, warnings)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		_ = simpleLogger.Write("error", fmt.Sprintf("%v %s:%d", err, f, l-2))

		log.Fatal(err)
	}
	if msgWarning != "" {
		_ = simpleLogger.Write("warning", msgWarning)
	}

	// ************************************************************************
	// ************* инициализация модуля взаимодействия с Zabbix *************
	channelZabbix := make(chan commoninterfaces.Messager)
	wzis := wrappers.WrappersZabbixInteractionSettings{
		NetworkPort: confApp.Zabbix.NetworkPort,
		NetworkHost: confApp.Zabbix.NetworkHost,
		ZabbixHost:  confApp.Zabbix.ZabbixHost,
	}
	eventTypes := []wrappers.EventType(nil)
	for _, v := range confApp.Zabbix.EventTypes {
		eventTypes = append(eventTypes, wrappers.EventType{
			IsTransmit: v.IsTransmit,
			EventType:  v.EventType,
			ZabbixKey:  v.ZabbixKey,
			Handshake: wrappers.Handshake{
				TimeInterval: v.Handshake.TimeInterval,
				Message:      v.Handshake.Message,
			},
		})
	}
	wzis.EventTypes = eventTypes
	wrappers.WrappersZabbixInteraction(ctx, wzis, simpleLogger, channelZabbix)

	//***************************************************************************
	//************** инициализация обработчика логирования данных ***************
	logging := logginghandler.New()
	go logginghandler.LoggingHandler(ctx, simpleLogger, channelZabbix, logging.GetChan())

	// ***************************************************************************
	// ************ инициализируем модуль временного хранения информации *********
	storageApp := memorytemporarystorage.NewTemporaryStorage()

	// вывод данных счетчика
	counting := counterHandler(ctx, storageApp, simpleLogger, channelZabbix)

	// ***************************************************************************
	// ************** инициализация модуля для взаимодействия с NATS *************
	// *************** (Данный модуль обязателен для взаимодействия) *************
	natsModule, err := natsapi.NewClientNATS(confApp.AppConfigNATS, confApp.AppConfigTheHive, storageApp, logging, counting)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		_ = simpleLogger.Write("error", fmt.Sprintf(" '%s' %s:%d", err, f, l-2))

		log.Fatal(err)
	}

	// ***************************************************************************
	// *********** инициализация модуля для взаимодействия с СУБД Redis **********
	redisModule := redisapi.HandlerRedis(ctx, *confApp.GetAppRedis(), storageApp, logging)

	// ***************************************************************************
	// *************** инициалиация модуля для взаимодействия с MISP *************
	mispModule, err := mispapi.HandlerMISP(*confApp.GetAppMISP(), confApp.GetListOrganization(), logging)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		_ = simpleLogger.Write("error", fmt.Sprintf(" '%s' %s:%d", err, f, l-2))
	}

	appStatus := fmt.Sprintf("%vproduction%v", constants.Ansi_Bright_Blue, constants.Ansi_Reset)
	envValue, ok := os.LookupEnv("GO_PHMISP_MAIN")
	if ok && envValue == "development" {
		appStatus = fmt.Sprintf("%v%s%v", constants.Ansi_Bright_Red, envValue, constants.Ansi_Reset)
	}

	msg := fmt.Sprintf("Application '%s' v%s was successfully launched", appname.GetAppName(), appversion.GetAppVersion())
	fmt.Printf("\n\n%v%v%s.%v\n", constants.Bold_Font, constants.Ansi_Bright_Green, msg, constants.Ansi_Reset)
	fmt.Printf("%v%vApplication status is '%s'.%v\n", constants.Underlining, constants.Ansi_Bright_Green, appStatus, constants.Ansi_Reset)
	_ = simpleLogger.Write("info", strings.ToLower(msg))

	// выполняется для записи в файл и отправки в Zabbix
	_ = simpleLogger.Write("info", msg)

	core := coremodule.NewCoreHandler(storageApp, logging, counting)
	if err := core.CoreHandler(ctx, natsModule, mispModule, redisModule, lr); err != nil {
		_, f, l, _ := runtime.Caller(0)
		_ = simpleLogger.Write("error", fmt.Sprintf("'%s' %s:%d", err.Error(), f, l-1))
		log.Fatalln(err)
	}
}
