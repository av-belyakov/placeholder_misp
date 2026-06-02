package confighandler

import "errors"

// GetCommonApp общие настройки приложения
func (cfg *ConfigApp) GetCommonApp() CommonCfg {
	return cfg.Common
}

// GetListLogs список логов
func (cfg *ConfigApp) GetListLogs() []*LogSet {
	return cfg.Common.LogList
}

// GetListOrganization список организаций
func (cfg *ConfigApp) GetListOrganization() []Organization {
	return cfg.Common.Organizations
}

// GetSqlite3 настройки доступа к БД
func (cfg *ConfigApp) GetSqlite3() CfgSqlite3 {
	return cfg.Sqlite3
}

// GetNATS настройки доступа к NATS
func (cfg *ConfigApp) GetNATS() CfgNATS {
	return cfg.NATS
}

// GetMISP настройки доступа к MISP
func (cfg *ConfigApp) GetMISP() CfgMISP {
	return cfg.MISP
}

// GetTheHive настройки доступа к TheHive
func (cfg *ConfigApp) GetTheHive() CfgTheHive {
	return cfg.TheHive
}

// GetRules настройки правил
func (cfg *ConfigApp) GetRules() CfgRules {
	return cfg.Rules
}

// GetLogDB настройки доступа к БД для логирования данных
func (cfg *ConfigApp) GetLogDB() CfgWriteLogDB {
	return cfg.WriteLogDB
}

// GetDebugServer настройки доступа к БД для логирования данных
func (cfg *ConfigApp) GetDebugServer() CfgDebugServer {
	return cfg.DebugServer
}

func (cfg *ConfigApp) Clean() {
	cfg = &ConfigApp{}
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
