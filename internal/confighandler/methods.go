package confighandler

import "errors"

func (conf *ConfigApp) GetCommonApp() *CommonCfg {
	return &conf.CommonCfg
}

func (conf *ConfigApp) GetListLogs() []*LogSet {
	return conf.LogList
}

func (conf *ConfigApp) GetListOrganization() []Organization {
	return conf.Organizations
}

func (conf *ConfigApp) GetSqlite3() *CfgSqlite3 {
	return &conf.CfgSqlite3
}

func (conf *ConfigApp) GetNATS() *CfgNATS {
	return &conf.CfgNATS
}

func (conf *ConfigApp) GetMISP() *CfgMISP {
	return &conf.CfgMISP
}

func (conf *ConfigApp) GetTheHive() *CfgTheHive {
	return &conf.CfgTheHive
}

// GetApplicationWriteLogDB настройки доступа к БД для логирования данных
func (conf *ConfigApp) GetApplicationWriteLogDB() *CfgWriteLogDB {
	return &conf.CfgWriteLogDB
}

func (conf *ConfigApp) Clean() {
	conf = &ConfigApp{}
}

// SetNameMessageType наименование типа логирования
func (l *LogSet) SetNameMessageType(v string) error {
	if v == "" {
		return errors.New("the value 'MsgTypeName' must not be empty")
	}

	return nil
}

// SetMaxLogFileSize максимальный размер файла для логирования
func (l *LogSet) SetMaxLogFileSize(v int) error {
	if v < 1000 {
		return errors.New("the value 'MaxFileSize' must not be less than 1000")
	}

	return nil
}

// SetPathDirectory путь к директории логирования
func (l *LogSet) SetPathDirectory(v string) error {
	if v == "" {
		return errors.New("the value 'PathDirectory' must not be empty")
	}

	return nil
}

// SetWritingStdout запись логов на вывод stdout
func (l *LogSet) SetWritingStdout(v bool) {
	l.WritingStdout = v
}

// SetWritingFile запись логов в файл
func (l *LogSet) SetWritingFile(v bool) {
	l.WritingFile = v
}

// SetWritingDB запись логов  в БД
func (l *LogSet) SetWritingDB(v bool) {
	l.WritingDB = v
}

// GetNameMessageType наименование тпа логирования
func (l *LogSet) GetNameMessageType() string {
	return l.MsgTypeName
}

// GetMaxLogFileSize максимальный размер файла для логирования
func (l *LogSet) GetMaxLogFileSize() int {
	return l.MaxFileSize
}

// GetPathDirectory путь к директории логирования
func (l *LogSet) GetPathDirectory() string {
	return l.PathDirectory
}

// GetWritingStdout запись логов на вывод stdout
func (l *LogSet) GetWritingStdout() bool {
	return l.WritingStdout
}

// GetWritingFile запись логов в файл
func (l *LogSet) GetWritingFile() bool {
	return l.WritingFile
}

// GetWritingDB запись логов  в БД
func (l *LogSet) GetWritingDB() bool {
	return l.WritingDB
}
