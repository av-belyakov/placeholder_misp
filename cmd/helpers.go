package main

import (
	"errors"
	"fmt"

	"github.com/av-belyakov/placeholder_misp/internal/confighandler"
	rules "github.com/av-belyakov/placeholder_misp/rulesinteraction"
	"github.com/av-belyakov/simplelogger"
)

func getLoggerSettings(cls []confighandler.LogSet) []simplelogger.Options {
	loggerConf := make([]simplelogger.Options, 0, len(cls))

	for _, v := range cls {
		loggerConf = append(loggerConf, simplelogger.Options{
			WritingToStdout: v.WritingStdout,
			WritingToFile:   v.WritingFile,
			WritingToDB:     v.WritingDB,
			MsgTypeName:     v.MsgTypeName,
			PathDirectory:   v.PathDirectory,
			MaxFileSize:     v.MaxFileSize,
		})
	}

	return loggerConf
}

func checkListRule(listRule *rules.ListRule, warnings []string) (msgWarning string, err error) {
	// поиск логических ошибок в файле с YAML правилами
	if len(warnings) > 0 {
		var warningStr string
		for _, v := range warnings {
			warningStr += fmt.Sprintln(v)
		}

		msgWarning = fmt.Sprintf("the following rules have a number of logical errors: %s\n", warningStr)
	}

	// проверка наличия правил Pass или Passany
	if len(listRule.GetRulePass()) == 0 && !listRule.GetRulePassany() {
		err = errors.New("there are no rules for handling messages received from NATS or all rules have failed validation")
	}

	return
}
