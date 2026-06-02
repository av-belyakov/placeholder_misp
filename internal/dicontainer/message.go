package dicontainer

import (
	"context"
	"log"
	"os"

	"github.com/av-belyakov/simplelogger"

	"github.com/av-belyakov/placeholder_misp/constants"
	"github.com/av-belyakov/placeholder_misp/internal/confighandler"
	"github.com/av-belyakov/placeholder_misp/internal/countermessage"
	"github.com/av-belyakov/placeholder_misp/internal/elasticsearchapi"
	"github.com/av-belyakov/placeholder_misp/internal/logginghandler"
	"github.com/av-belyakov/placeholder_misp/internal/mispapi"
	"github.com/av-belyakov/placeholder_misp/internal/natsapi"
	"github.com/av-belyakov/placeholder_misp/internal/ruleshandler"
	"github.com/av-belyakov/placeholder_misp/internal/sqlite3api"
	"github.com/av-belyakov/placeholder_misp/internal/supportingfunctions"
)

// Configer чтения конфигурационного файла
func (d *DiContainer) Configer() Configer {
	if d.configer == nil {
		rootPath, err := supportingfunctions.GetRootPath(d.rootDir)
		if err != nil {
			log.Fatalf("error, it is impossible to form root path (%s)", err.Error())
		}

		cfg, err := confighandler.New(rootPath)
		if err != nil {
			log.Fatal("error module 'confighandler':", err)
		}

		d.configer = cfg
	}

	return d.configer
}

// SimpleLogger простое логирование с помощью стороннего пакета
func (d *DiContainer) SimpleLogger(ctx context.Context) SimpleLogger {
	if d.simpleLogger == nil {
		listLog := make([]simplelogger.OptionsManager, 0, len(d.Configer().GetListLogs()))
		for _, v := range d.Configer().GetListLogs() {
			listLog = append(listLog, v)
		}

		opts := simplelogger.CreateOptions(listLog...)
		simpleLogger, err := simplelogger.NewSimpleLogger(ctx, d.rootDir, opts)
		if err != nil {
			log.Fatal("error module 'simplelogger':", err)
		}

		d.simpleLogger = simpleLogger

		//подключение логирования в БД
		simpleLogger.SetDataBaseInteraction(d.DbLogger())
	}

	return d.simpleLogger
}

// Logger основное логирование
func (d *DiContainer) Logger(ctx context.Context) Logger {
	if d.logger == nil {
		logger := logginghandler.New(d.SimpleLogger(ctx), d.ch)
		logger.Start(ctx)

		d.logger = logger
	}

	return d.logger
}

// Counter счетчик сообщений
func (d *DiContainer) Counter(ctx context.Context) Counter {
	if d.counter == nil {
		counter := countermessage.New(d.ch)
		counter.Start(ctx)

		d.counter = counter
	}

	return d.counter
}

// DbLogger запись логов в БД
func (d *DiContainer) DbLogger() DbLogger {
	if d.dbLogger == nil {
		var nameRegionalObject = "gcm"
		if os.Getenv("GO_PHMISP_MAIN") == "development" {
			nameRegionalObject = "gcm-test"
		}

		conn, err := elasticsearchapi.NewElasticsearchConnect(elasticsearchapi.Settings{
			Port:               d.Configer().GetLogDB().Port,
			Host:               d.Configer().GetLogDB().Host,
			User:               d.Configer().GetLogDB().User,
			Passwd:             d.Configer().GetLogDB().Passwd,
			IndexDB:            d.Configer().GetLogDB().StorageNameDB,
			NameRegionalObject: nameRegionalObject,
		})
		if err != nil {
			log.Fatal("error module 'elasticsearchapi':", err)
		}

		d.dbLogger = conn
	}

	return d.dbLogger
}

// DB подключение к БД
func (d *DiContainer) DB(ctx context.Context) DB {
	if d.db == nil {
		rootPath, err := supportingfunctions.GetRootPath(d.rootDir)
		if err != nil {
			log.Fatalf("error, it is impossible to form root path (%s)", err.Error())
		}

		// инициализируем файл базы данных sqlite3
		newPathSqlite3Db, err := sqlite3api.Sqlite3DbFileIsExist(rootPath, d.Configer().GetSqlite3().PathFileDb)
		if err != nil {
			log.Fatal("error file sqlite3 database:", err)
		}

		sqlite3Module, err := sqlite3api.New(ctx, newPathSqlite3Db, d.Logger(ctx))
		if err != nil {
			log.Fatal("error module 'database':", err)
		}

		d.db = sqlite3Module
	}

	return d.db
}

// NatsConnecter подключение к NATS
func (d *DiContainer) NatsConnecter(ctx context.Context) NatsConnecter {
	if d.nats == nil {
		apiNats, err := natsapi.New(
			d.Logger(ctx),
			d.Counter(ctx),
			natsapi.WithHost(d.Configer().GetNATS().Host),
			natsapi.WithPort(d.Configer().GetNATS().Port),
			natsapi.WithCacheTTL(d.Configer().GetNATS().CacheTTL),
			natsapi.WithSubcriptionListenerCase(d.Configer().GetNATS().Subscriptions.ListenerCase),
			natsapi.WithSubcriptionSenderCommand(d.Configer().GetNATS().Subscriptions.SenderCommand),
			natsapi.WithSubcriptionGetSensorInfo(d.Configer().GetNATS().Subscriptions.GetSensorInfo),
		)
		if err != nil {
			log.Fatal("error initialization module 'natsapi':", err)
		}

		if err = apiNats.Start(ctx); err != nil {
			log.Fatal("error start module 'natsapi':", err)
		}

		d.nats = apiNats
	}

	return d.nats
}

// MispConnecter подключение к MISP
func (d *DiContainer) MispConnecter(ctx context.Context) MispConnecter {
	if d.misp == nil {
		apiMisp, err := mispapi.NewModuleMISP(d.Configer().GetMISP().Host, d.Configer().GetMISP().Auth, d.Configer().GetListOrganization(), d.Logger(ctx))
		if err != nil {
			log.Fatalln("error initialization module 'mispapi':", err)
		}

		if err = apiMisp.Start(ctx); err != nil {
			log.Fatal("error start module 'natsapi':", err)
		}

		d.misp = apiMisp
	}

	return d.misp
}

// RulesHandler обработчик правил
func (d *DiContainer) Rules(ctx context.Context) RulesHandler {
	if d.rules == nil {
		listRules, warnings, err := ruleshandler.NewListRule(constants.Root_Dir, d.Configer().GetRules().Directory, d.Configer().GetRules().File)
		if err != nil {
			log.Fatal("error module 'ruleshandler':", err)
		}

		// проверяем наличие правил Pass или Passany которые являются обязательными, а также отсутсвие
		// логических ошибок в файле с правилами
		msgWarning, err := ruleshandler.CheckListRule(listRules, warnings)
		if err != nil {
			log.Fatal("error module 'ruleshandler':", err)
		}

		if msgWarning != "" {
			d.SimpleLogger(ctx).Write("warning", msgWarning)
		}

		d.rules = listRules
	}

	return d.rules
}
