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
// FieldName - наименование искомого поля
// TypeValue - тип значения искомого поля, должно содержать одно из значений string, int, bool
// Value - значение искомого поля
// ReplaceValue - значение на которое нужно заменить (может быть в том числе и пустым)
type ListRequiredValue struct {
	FieldName    string
	TypeValue    string
	Value        string
	ReplaceValue string
}

// RuleProcMISPMessageFields описание обрабатываемых полей misp сообщения
// ActionType - тип действия (обрабатывать, обрабатывать все, заменить, очистить значение поля)
// ListRequiredValues требуемые значения для обработки правила
type RuleProcMISPMessageField struct {
	ActionType         string
	ListRequiredValues []ListRequiredValue
}

// ListRulesProcessedMISPMessage список обрабатываемых полей misp сообщения
// где свойство map есть имя поля
type ListRulesProcMISPMessage struct {
	Rulles []RuleProcMISPMessageField
}

/*
func (rp RuleProcessedMISPMessageFields) GetActionType() string {
	return rp.ActionType
}

func (rp RuleProcessedMISPMessageFields) GetRequiredType() string {
	return rp.RequiredValue.ValueType
}

func (rp RuleProcessedMISPMessageFields) GetRequiredValueString() string {
	if str, ok := rp.RequiredValue.Value.(string); ok {
		return str
	}

	return ""
}

func (rp RuleProcessedMISPMessageFields) GetRequiredValueInt() int {
	if num, ok := rp.RequiredValue.Value.(int); ok {
		return num
	}

	return 0
}

func (rp RuleProcessedMISPMessageFields) GetRequiredValueBool() bool {
	if b, ok := rp.RequiredValue.Value.(bool); ok {
		return b
	}

	return false
}

func (rp RuleProcessedMISPMessageFields) GetDesiredType() string {
	return rp.DesiredValue.ValueType
}

func (rp RuleProcessedMISPMessageFields) GetDesiredValueString() string {
	if str, ok := rp.DesiredValue.Value.(string); ok {
		return str
	}

	return ""
}

func (rp RuleProcessedMISPMessageFields) GetDesiredValueInt() int {
	if num, ok := rp.DesiredValue.Value.(int); ok {
		return num
	}

	return 0
}

func (rp RuleProcessedMISPMessageFields) GetDesiredValueBool() bool {
	if b, ok := rp.DesiredValue.Value.(bool); ok {
		return b
	}

	return false
}

func (rp RuleProcessedMISPMessageFields) GetReplacedType() string {
	return rp.DesiredValue.ValueType
}

func (rp RuleProcessedMISPMessageFields) GetReplacedValueString() string {
	if str, ok := rp.ReplacedValue.Value.(string); ok {
		return str
	}

	return ""
}

func (rp RuleProcessedMISPMessageFields) GetReplacedValueInt() int {
	if num, ok := rp.ReplacedValue.Value.(int); ok {
		return num
	}

	return 0
}

func (rp RuleProcessedMISPMessageFields) GetReplacedValueBool() bool {
	if b, ok := rp.ReplacedValue.Value.(bool); ok {
		return b
	}

	return false
}

func (rf RuleProcessedMISPMessageFields) GetListRelatedFields() map[string]interface{} {
	return rf.ListRelatedFields
}
*/
