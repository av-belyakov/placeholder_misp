package ruleshandler

import (
	"errors"
	"fmt"
)

// CheckListRule проверяет наличие правил Pass или Passany которые являются обязательными,
// а также отсутсвие логических ошибок в файле с правилами
func CheckListRule(listRule *ListRule, warnings []string) (msgWarning string, err error) {
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
