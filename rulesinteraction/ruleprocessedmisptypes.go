package rules

type ProcesserMISPRule interface {
	GetterMISPRuleProcessed
}

type GetterMISPRuleProcessed interface {
	GetActionType() string
	GetListRelatedFields() map[string]interface{}
	GetterRequiredValue
	GetterDesiredValue
	GetterReplacedValue
}

type GetterRequiredValue interface {
	GetRequiredType() string
	GetRequiredValueString() string
	GetRequiredValueInt() int
	GetRequiredValueBool() bool
}

type GetterDesiredValue interface {
	GetDesiredType() string
	GetDesiredValueString() string
	GetDesiredValueInt() int
	GetDesiredValueBool() bool
}

type GetterReplacedValue interface {
	GetReplacedType() string
	GetReplacedValueString() string
	GetReplacedValueInt() int
	GetReplacedValueBool() bool
}

/*
type ValueParameters struct {
	ValueType string
	Value     interface{}
}

// RuleProcessedMISPMessageFields описание обрабатываемых полей misp сообщения
// FieldName - наименование поля
// ActionType - тип действия (обрабатывать, обрабатывать все, заменить, очистить значение поля)
// anypass, process, replace, clean.
// RequiredValue - обязательное значение, если поле пустое то не учитывается, если заполненно то все события без этого
// значения дальше не обрабатываются и отбрасываются
// DesiredValue - значение ПО которому выполняется поиск (должен выполнятся поиск по содержимому в строке, например snort_id, а ищем только snort)
// ReplacedValue - значение НА которое следует заменить
// ListRelatedFields - список сопутствующих значений (нужен для обработки дополнительных полей). Например, если поле dataType содержит
// snort, то заменить содержимое поля data: '39609404, 39609404' на data: ”
type RuleProcessedMISPMessageFields struct {
	FieldName         string
	ActionType        string
	RequiredValue     ValueParameters
	DesiredValue      ValueParameters
	ReplacedValue     ValueParameters
	ListRelatedFields map[string]interface{}
}
*/

// ListRequiredValue описание полей содержащихся в свойстве listRequiredValues
// FieldSearchName - наименование искомого поля
// TypeValue - тип значения искомого поля, должно содержать одно из значений string, int, bool
// SearchValue - значение искомого поля
// ReplaceValue - значение на которое нужно заменить (может быть в том числе и пустым)
type ListRequiredValue struct {
	FieldSearchName string
	TypeValue       string //может это и не надо
	SearchValue     string
	ReplaceValue    string
}

// RuleProcMISPMessageFields описание обрабатываемых полей misp сообщения
// ActionType - тип действия (обрабатывать, обрабатывать все, заменить, очистить значение поля)
// ListRequiredValues требуемые значения для обработки правила
type RuleProcMISPMessageField struct {
	ActionType         string
	ListRequiredValues []ListRequiredValue
}

// ListRulesProcessedMISPMessage список обрабатываемых полей misp сообщения
// Rules - описание основных правил
// SearchFieldsName - список полей по которым осуществляется поиск
// SearchValuesName - список значений по которым осуществляется поиск
type ListRulesProcMISPMessage struct {
	Rules            []RuleProcMISPMessageField
	SearchFieldsName map[string][][2]int
	SearchValuesName map[string][][2]int
}

/*

Ниже основное описание типов полей правил

*/

// ListRulesProcessingMsgMISP содержит список правил обработки сообщений предназначенных для MISP
// Rules основной список правил полученный из конфигурационного файла
// RulesIndex список индексированных правил для более удобной обработки
type ListRulesProcessingMsgMISP struct {
	Rules      RuleSetProcessingMsgMISP
	RulesIndex map[string][]RuleIndex
}

// RuleSetProcessingMsgMISP содержит правила обработки сообщений
// Passany тип правила для пропуска всех сообщений
// Pass тип правила для пропуска сообщений подходящих под определенные критерии
// Replace тип правила для замены определенных значений подходящих под определенные критерии
type RuleSetProcessingMsgMISP struct {
	Passany  bool
	Pass     []RulePass
	Replace  []RuleReplace
	Passtest []PasstestListAnd
}

type PasstestListAnd struct {
	ListAnd []RulePass
}

// RulePassany содержит тип правила для пропуска всех сообщений
// IsPass разрешен ли пропуск всех сообщений
type RulePassany struct {
	IsPass bool
}

// RulePass содержит тип правила для пропуска сообщений подходящих под определенные критерии
// SearchField искомое поле
// SearchValue искомое значение
type RulePass struct {
	SearchField, SearchValue string
}

// RuleReplace содержит тип правила для замены определенных значений
// SearchField искомое поле
// SearchValue искомое значение
// ReplaceValue заменяемое значение
type RuleReplace struct {
	SearchField, SearchValue, ReplaceValue string
}

// RuleIndex содержит список индексированных правил для более удобной обработки
// RuleType тип правила
// SearchField искомое поле
// ReplaceValue заменяемое значение
type RuleIndex struct {
	RuleType, SearchField, ReplaceValue string
}
