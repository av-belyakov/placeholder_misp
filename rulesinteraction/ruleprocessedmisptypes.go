package rules

// ListRule содержит список правил
type ListRule struct {
	Rules RuleOptions `yaml:"RULES"`
}

// RuleOptions содержит опции правил
type RuleOptions struct {
	Passany bool          `yaml:"PASSANY"`
	Pass    []PassListAnd `yaml:"PASS"`
	Replace []RuleReplace `yaml:"REPLACE"`
	Exclude []RuleExclude `yaml:"EXCLUDE"`
}

// PassListAnd список правил с логикой 'И'
type PassListAnd struct {
	ListAnd []RulePass `yaml:"listAnd"`
}

// CommonRuleFields общие поля которые могут использоватся для описания большинства типов правил
// SearchField искомое поле
// SearchValue искомое значение
type CommonRuleFields struct {
	SearchField string `yaml:"searchField"`
	SearchValue string `yaml:"searchValue"`
}

// RulePassany содержит тип правила для пропуска всех сообщений
// IsPass разрешен ли пропуск всех сообщений
type RulePassany struct {
	IsPass bool
}

// RulePass содержит тип правила для пропуска сообщений подходящих под определенные критерии
// StatementExpression утверждение выражения
type RulePass struct {
	CommonRuleFields    `mapstructure:",squash"`
	StatementExpression bool
}

// RuleReplace содержит тип правила для замены определенных значений
// ReplaceValue заменяемое значение
type RuleReplace struct {
	CommonRuleFields `mapstructure:",squash"`
	ReplaceValue     string `yaml:"replaceValue"`
}

// RuleExclude содержит тип правила для исключения объекта из списка объектов
// предназначенных для передачи
// AccurateComparison содержит тригер, информирующий о типе поиска, 'true' является
// 'строгим' поиском, 'false', нет
type RuleExclude struct {
	CommonRuleFields   `mapstructure:",squash"`
	AccurateComparison bool `yaml:"accurateComparison"`
}
