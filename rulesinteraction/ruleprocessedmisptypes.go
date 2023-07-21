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

type ListRulesProcessingMsgMISP struct {
	Rules RuleSetProcessingMsgMISP
}

type RuleSetProcessingMsgMISP struct {
	Passany bool
	Pass    []RulePass
	Replace []RuleReplace
}

type RulePassany struct {
	IsPass bool
}

type RulePass struct {
	FieldSearchName string
	SearchValue     string
}

type RuleReplace struct {
	FieldSearchName string
	SearchValue     string
	ReplaceValue    string
}
