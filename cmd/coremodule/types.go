package coremodule

import (
	"github.com/av-belyakov/placeholder_misp/cmd/mispapi"
	"github.com/av-belyakov/placeholder_misp/cmd/sqlite3api"
	"github.com/av-belyakov/placeholder_misp/commoninterfaces"
	rules "github.com/av-belyakov/placeholder_misp/internal/ruleshandler"
)

// ChanInputCreateMispFormat входные данные для канала приёма информации используемого
// при создании объектов  вформате MISP
type ChanInputCreateMispFormat struct {
	Value       interface{} //любые передаваемые данные
	UUID        string      //уникальный идентификатор в формате UUID
	FieldName   string      //наименование поля
	ValueType   string      //тип передаваемого значения (string, int и т.д.)
	FieldBranch string      //'путь' до значения в как в JSON формате, например 'event.details.customFields.class'
}

type FieldsNameMapping struct {
	InputFieldName, MispFieldName string
}

// storageValueName временное хранилище свойств элементов observables
type storageValueName []string

// ExclusionRules содержит информацию об объектах которые нужно исключить из
// передачи в MISP
type ExclusionRules []ExclusionRule

// ExclusionRule
type ExclusionRule struct {
	NameList       string //наименование объекта
	SequenceNumber int    //порядковый номер в списке объектов
}

type MispGalaxyTags map[int]MispGalaxyOptions

type MispGalaxyOptions struct {
	Name, PatternId, PatternType string
}

// GenerateObjectsFormatMISP генерирует объекты в формате MISP
type GenerateObjectsFormatMISP struct {
	mispModule    mispapi.ModuleMispHandler
	sqlite3Module *sqlite3api.ApiSqlite3Module
	listRule      *rules.ListRule
	counter       commoninterfaces.Counter
	logger        commoninterfaces.Logger
}

type SettingsGenerateObjectsFormatMISP struct {
	MispModule    mispapi.ModuleMispHandler
	Sqlite3Module *sqlite3api.ApiSqlite3Module
	ListRule      *rules.ListRule
	Counter       commoninterfaces.Counter
	Logger        commoninterfaces.Logger
}
