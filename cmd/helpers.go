package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/av-belyakov/placeholder_misp/constants"
	"github.com/av-belyakov/placeholder_misp/internal/appname"
	"github.com/av-belyakov/placeholder_misp/internal/appversion"
	rules "github.com/av-belyakov/placeholder_misp/internal/ruleshandler"
)

// checkSqlite3DbFileExist проверяет наличие файла базф данных Sqlite3
// и при необходимости создает его из резервного файла
func checkSqlite3DbFileExist(pathFileDb string) error {
	backupFile := "../internal/backupdb/sqlite3_backup.db"

	// наличие файла backup
	if _, err := os.Stat(backupFile); err != nil {
		return err
	}

	//файл с основной БД
	_, err := os.Stat(pathFileDb)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		fr, err := os.OpenFile(backupFile, os.O_RDONLY, 0666)
		if err != nil {
			return err
		}
		defer fr.Close()

		fw, err := os.OpenFile(pathFileDb, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return err
		}
		defer fw.Close()

		if _, err := io.Copy(fw, fr); err != nil {
			return err
		}
	}

	return nil
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

func getInformationMessage() string {
	appStatus := fmt.Sprintf("%vproduction%v", constants.Ansi_Bright_Blue, constants.Ansi_Reset)
	envValue, ok := os.LookupEnv("GO_PHMISP_MAIN")
	if ok && envValue == "development" {
		appStatus = fmt.Sprintf("%v%s%v", constants.Ansi_Bright_Red, envValue, constants.Ansi_Reset)
	}

	msg := fmt.Sprintf("Application '%s' v%s was successfully launched", appname.GetAppName(), appversion.GetAppVersion())

	fmt.Printf("\n%v%v%s.%v\n", constants.Bold_Font, constants.Ansi_Bright_Green, msg, constants.Ansi_Reset)
	fmt.Printf("%v%vApplication status is '%s'.%v\n", constants.Underlining, constants.Ansi_Bright_Green, appStatus, constants.Ansi_Reset)

	return msg
}
