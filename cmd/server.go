package main

import (
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"
	"strings"

	"github.com/av-belyakov/simplelogger"

	"github.com/av-belyakov/placeholder_misp/cmd/coremodule"
	"github.com/av-belyakov/placeholder_misp/cmd/mispapi"
	"github.com/av-belyakov/placeholder_misp/cmd/natsapi"
	"github.com/av-belyakov/placeholder_misp/cmd/redisapi"
	"github.com/av-belyakov/placeholder_misp/cmd/wrappers"
	"github.com/av-belyakov/placeholder_misp/commoninterfaces"
	"github.com/av-belyakov/placeholder_misp/constants"
	"github.com/av-belyakov/placeholder_misp/internal/confighandler"
	"github.com/av-belyakov/placeholder_misp/internal/countermessage"
	"github.com/av-belyakov/placeholder_misp/internal/logginghandler"
	rules "github.com/av-belyakov/placeholder_misp/internal/ruleshandler"
	"github.com/av-belyakov/placeholder_misp/internal/supportingfunctions"
)

func server(ctx context.Context) {
	rootPath, err := supportingfunctions.GetRootPath(constants.Root_Dir)
	if err != nil {
		log.Fatalf("error, it is impossible to form root path (%s)", err.Error())
	}

	// ****************************************************************************
	// *********** инициализируем модуль чтения конфигурационного файла ***********
	confApp, err := confighandler.New(rootPath)
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
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(err).Error())

		log.Fatalf("error module 'rulesinteraction': %v\n", err)
	}

	//проверяем наличие правил Pass или Passany которые являются обязательными,
	//а также отсутсвие логических ошибок в файле с правилами
	msgWarning, err := checkListRule(lr, warnings)
	if err != nil {
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(err).Error())

		log.Fatal(err)
	}
	if msgWarning != "" {
		_ = simpleLogger.Write("warning", msgWarning)
	}

	// ************************************************************************
	// ************* инициализация модуля взаимодействия с Zabbix *************
	chZabbix := make(chan commoninterfaces.Messager)
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
	wrappers.WrappersZabbixInteraction(ctx, wzis, simpleLogger, chZabbix)

	//***************************************************************************
	//************** инициализация обработчика логирования данных ***************
	logging := logginghandler.New()
	go logginghandler.LoggingHandler(ctx, simpleLogger, chZabbix, logging.GetChan())

	// ***************************************************************************
	// *********** инициализируем модуль счётчика для подсчёта сообщений *********
	counting := countermessage.New(chZabbix)
	if err = counting.Handler(ctx); err != nil {
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(err).Error())
	}

	// ***************************************************************************
	// ************** инициализация модуля для взаимодействия с NATS *************
	natsModule, err := natsapi.NewClientNATS(confApp.AppConfigNATS, confApp.AppConfigTheHive, counting, logging)
	if err != nil {
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(err).Error())

		log.Fatal(err)
	}

	// ***************************************************************************
	// *********** инициализация модуля для взаимодействия с СУБД Redis **********
	redisModule := redisapi.HandlerRedis(ctx, *confApp.GetAppRedis(), logging)

	// ***************************************************************************
	// *************** инициалиация модуля для взаимодействия с MISP *************
	mispModule, err := mispapi.HandlerMISP(*confApp.GetAppMISP(), confApp.GetListOrganization(), logging)
	if err != nil {
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(err).Error())
	}

	// вывод информационного сообщения при старте приложения
	msg := getInformationMessage()

	_ = simpleLogger.Write("info", strings.ToLower(msg))

	//для отладки через pprof
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// *****************************************************************
	// *************** инициализируем ядро приложения ******************
	core := coremodule.NewCoreHandler(counting, lr, logging)
	if err := core.Start(ctx, natsModule, mispModule, redisModule); err != nil {
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(err).Error())

		log.Fatalln(err)
	}
}
