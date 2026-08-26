package dicontainer

import (
	"context"

	"github.com/av-belyakov/simplelogger"

	"github.com/av-belyakov/placeholder_misp/v2/commoninterfaces"
	"github.com/av-belyakov/placeholder_misp/v2/internal/confighandler"
	"github.com/av-belyakov/placeholder_misp/v2/internal/mispapi"
	"github.com/av-belyakov/placeholder_misp/v2/internal/natsapi"
	"github.com/av-belyakov/placeholder_misp/v2/internal/ruleshandler"
	"github.com/av-belyakov/placeholder_misp/v2/internal/sqlite3api"
)

type Logger interface {
	GetChan() <-chan commoninterfaces.Messager
	Send(msgType, message string)
	Close()
}

type Counter interface {
	SendMessage(msgType string, count int)
}

type SimpleLogger interface {
	SetDataBaseInteraction(dbi simplelogger.DataBaseInteractor)
	GetCountFileDescription() int
	GetListTypeFiles() []string
	Write(typeLog, msg string) bool
}

type Configer interface {
	GetCommonApp() confighandler.CommonCfg
	GetNATS() confighandler.CfgNATS
	GetMISP() confighandler.CfgMISP
	GetTheHive() confighandler.CfgTheHive
	GetSqlite3() confighandler.CfgSqlite3
	GetRules() confighandler.CfgRules
	GetListLogs() []*confighandler.LogSet
	GetListOrganization() []confighandler.Organization
	GetLogDB() confighandler.CfgWriteLogDB
	GetDebugServer() confighandler.CfgDebugServer
}

type NatsConnecter interface {
	GetChannelFromModule() <-chan natsapi.OutputSettings
	GetChannelToModule() chan natsapi.InputSettings
	SendingDataInput(data natsapi.InputSettings)
	SendingDataOutput(data natsapi.OutputSettings)
}

type MispConnecter interface {
	GetReceptionChannel() <-chan mispapi.OutputSetting
	GetInputChannel() <-chan mispapi.InputSettings
	SendDataOutput(data mispapi.OutputSetting)
	SendDataInput(data mispapi.InputSettings)
}

type DB interface {
	GetChRequest() <-chan sqlite3api.Request
	SendData(req sqlite3api.Request)
	SearchCaseId(ctx context.Context, source string, caseId int) (int, error)
	UpdateCaseId(ctx context.Context, source string, caseId, eventId int) error
	DeleteCaseId(ctx context.Context, caseId int) error
	Ping(ctx context.Context) error
	ConnectionClose()
}

type DbLogger interface {
	Write(msgType, msg string) error
}

type RulesHandler interface {
	GetRulePass() []ruleshandler.PassListAnd
	GetRuleExclude() []ruleshandler.ExcludeListAnd
	GetRulePassany() bool
	PassRuleHandler(string, any)
	ReplacementRuleHandler(string, string, any) (any, int, error)
	ExcludeRuleHandler(string, any) ([2]int, bool)
	SomePassRuleIsTrue() bool
	CleanStatementExpressionRulePass()
}
