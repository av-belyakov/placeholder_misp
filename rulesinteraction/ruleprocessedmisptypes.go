package rules

// ListRulesProcessingMsgMISP содержит список правил обработки сообщений предназначенных для MISP
// Rules основной список правил полученный из конфигурационного файла
type ListRulesProcessingMsgMISP struct {
	Rules RuleSetProcessingMsgMISP
}

// RuleSetProcessingMsgMISP содержит правила обработки сообщений
// Passany тип правила для пропуска всех сообщений
// Pass тип правила для пропуска сообщений подходящих под определенные критерии
// Replace тип правила для замены определенных значений подходящих под определенные критерии
type RuleSetProcessingMsgMISP struct {
	Passany bool
	Pass    []PassListAnd
	Replace []RuleReplace
}

// PassListAnd список правил
type PassListAnd struct {
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
// StatementExpression утверждение выражения
type RulePass struct {
	SearchField         string
	SearchValue         string
	StatementExpression bool
}

// RuleReplace содержит тип правила для замены определенных значений
// SearchField искомое поле
// SearchValue искомое значение
// ReplaceValue заменяемое значение
type RuleReplace struct {
	SearchField, SearchValue, ReplaceValue string
}

// CleanStatementExpressionRulePass приводит поле StatementExpression к значению false
// это поле проверяется на соответствие заданным правилам обрабатываемым значениям
// обязательно нужно выполнять данный метод после проверки значения StatementExpression
func (lr *ListRulesProcessingMsgMISP) CleanStatementExpressionRulePass() {
	for k, v := range lr.Rules.Pass {
		for key := range v.ListAnd {
			lr.Rules.Pass[k].ListAnd[key].StatementExpression = false
		}
	}
}
