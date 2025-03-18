package main

import (
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"

	"github.com/av-belyakov/placeholder_misp/cmd/elasticsearchapi"
	"github.com/av-belyakov/simplelogger"

	"github.com/av-belyakov/placeholder_misp/cmd/coremodule"
	"github.com/av-belyakov/placeholder_misp/cmd/mispapi"
	"github.com/av-belyakov/placeholder_misp/cmd/natsapi"
	"github.com/av-belyakov/placeholder_misp/cmd/sqlite3api"
	"github.com/av-belyakov/placeholder_misp/cmd/wrappers"
	"github.com/av-belyakov/placeholder_misp/commoninterfaces"
	"github.com/av-belyakov/placeholder_misp/constants"
	"github.com/av-belyakov/placeholder_misp/internal/appversion"
	"github.com/av-belyakov/placeholder_misp/internal/confighandler"
	"github.com/av-belyakov/placeholder_misp/internal/countermessage"
	"github.com/av-belyakov/placeholder_misp/internal/logginghandler"
	rules "github.com/av-belyakov/placeholder_misp/internal/ruleshandler"
	"github.com/av-belyakov/placeholder_misp/internal/supportingfunctions"
)

func server(ctx context.Context) {
	version, err := appversion.GetAppVersion()
	if err != nil {
		log.Println(err)
	}

	rootPath, err := supportingfunctions.GetRootPath(constants.Root_Dir)
	if err != nil {
		log.Fatalf("error, it is impossible to form root path (%s)", err.Error())
	}

	// ****************************************************************************
	// *********** инициализируем модуль чтения конфигурационного файла ***********
	conf, err := confighandler.New(rootPath)
	if err != nil {
		log.Fatalf("error module 'confighandler': %v", err)
	}

	// ****************************************************************************
	// ****************** инициализируем файл базы данных sqlite3 *****************
	newPathSqlite3Db, err := checkSqlite3DbFileExist(rootPath, conf.AppConfigSqlite3.PathFileDb)
	if err != nil {
		log.Fatalf("error file sqlite3 database: %v", err)
	}

	// ****************************************************************************
	// ********************* инициализация модуля логирования *********************
	var listLog []simplelogger.OptionsManager
	for _, v := range conf.GetListLogs() {
		listLog = append(listLog, v)
	}
	opts := simplelogger.CreateOptions(listLog...)
	simpleLogger, err := simplelogger.NewSimpleLogger(ctx, constants.Root_Dir, opts)
	if err != nil {
		log.Fatalf("error module 'simplelogger': %v", err)
	}

	//*********************************************************************************
	//********** инициализация модуля взаимодействия с БД для передачи логов **********
	confDB := conf.GetApplicationWriteLogDB()
	if esc, err := elasticsearchapi.NewElasticsearchConnect(elasticsearchapi.Settings{
		Port:               confDB.Port,
		Host:               confDB.Host,
		User:               confDB.User,
		Passwd:             confDB.Passwd,
		IndexDB:            confDB.StorageNameDB,
		NameRegionalObject: "gcm",
	}); err != nil {
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(err).Error())
	} else {
		//подключение логирования в БД
		simpleLogger.SetDataBaseInteraction(esc)
	}

	// ****************************************************************************
	// ******* инициализируем модуль чтения правил обработки MISP сообщений *******
	listRules, warnings, err := rules.NewListRule(constants.Root_Dir, conf.RulesProcMSGMISP.Directory, conf.RulesProcMSGMISP.File)
	if err != nil {
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(err).Error())

		log.Fatalf("error module 'rulesinteraction': %v\n", err)
	}

	//проверяем наличие правил Pass или Passany которые являются обязательными,
	//а также отсутсвие логических ошибок в файле с правилами
	msgWarning, err := checkListRule(listRules, warnings)
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
		NetworkPort: conf.Zabbix.NetworkPort,
		NetworkHost: conf.Zabbix.NetworkHost,
		ZabbixHost:  conf.Zabbix.ZabbixHost,
	}
	eventTypes := []wrappers.EventType(nil)
	for _, v := range conf.Zabbix.EventTypes {
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
	logging := logginghandler.New(simpleLogger, chZabbix)
	logging.Start(ctx)

	// ***************************************************************************
	// *********** инициализируем модуль счётчика для подсчёта сообщений *********
	counting := countermessage.New(chZabbix)
	counting.Start(ctx)

	// ***************************************************************************
	// ************** инициализация модуля для взаимодействия с NATS *************
	confNats := conf.AppConfigNATS
	apiNats, err := natsapi.New(
		logging,
		counting,
		natsapi.WithHost(confNats.Host),
		natsapi.WithPort(confNats.Port),
		natsapi.WithCacheTTL(confNats.CacheTTL),
		natsapi.WithListenerCase(confNats.Subscriptions.ListenerCase),
		natsapi.WithSenderCommand(confNats.Subscriptions.SenderCommand))
	if err != nil {
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(err).Error())

		log.Fatal(err)
	}
	if err = apiNats.Start(ctx); err != nil {
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(err).Error())

		log.Fatal(err)
	}

	// ***************************************************************************
	// ************ инициализация модуля для взаимодействия с Sqlite3 ************
	sqlite3Module, err := sqlite3api.New(ctx, newPathSqlite3Db, logging)
	if err != nil {
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(err).Error())

		log.Fatal(err)
	}

	// ***************************************************************************
	// *************** инициалиация модуля для взаимодействия с MISP *************
	mispModule, err := mispapi.NewModuleMISP(conf.GetAppMISP().Host, conf.GetAppMISP().Auth, conf.GetListOrganization(), logging)
	if err != nil {
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(err).Error())
	}
	if err = mispModule.Start(ctx); err != nil {
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(err).Error())

		log.Fatalln(err)
	}

	// вывод информационного сообщения при старте приложения
	msg := getInformationMessage(version)
	_ = simpleLogger.Write("info", strings.ToLower(msg))

	//для отладки через pprof
	//http://localhost:6161/debug/pprof/
	//go tool pprof http://localhost:6161/debug/pprof/heap
	//go tool pprof http://localhost:6161/debug/pprof/goroutine
	//go tool pprof http://localhost:6161/debug/pprof/allocs
	if os.Getenv("GO_HIVEHOOK_MAIN") == "development" {
		go func() {
			log.Println(http.ListenAndServe("localhost:6161", nil))
		}()
	}
	//------------------------

	// *****************************************************************
	// *************** инициализируем ядро приложения ******************
	core := coremodule.NewCoreHandler(counting, listRules, logging)
	core.Start(ctx, apiNats, mispModule, sqlite3Module)
}
