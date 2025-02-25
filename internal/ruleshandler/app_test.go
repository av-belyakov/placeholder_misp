package ruleshandler_test

import (
	"fmt"
	"testing"

	rules "github.com/av-belyakov/placeholder_misp/internal/ruleshandler"
	"github.com/stretchr/testify/assert"
)

const fileName = "mispmsgrule.yaml"

var (
	lr  *rules.ListRule
	lw  []string
	err error
)

func printRuleResult(r *rules.ListRule) string {
	resultPrint := fmt.Sprintln("RULES:")

	resultPrint += fmt.Sprintln("  REPLACE:")
	for k, v := range r.Rules.Replace {
		resultPrint += fmt.Sprintln("  ", k+1, ".")
		resultPrint += fmt.Sprintf("    searchField: '%s'\n", v.SearchField)
		resultPrint += fmt.Sprintf("    searchValue: '%s'\n", v.SearchValue)
		resultPrint += fmt.Sprintf("    replaceValue: '%s'\n", v.ReplaceValue)
	}

	resultPrint += fmt.Sprintln("  PASS:")
	for key, value := range r.Rules.Pass {
		resultPrint += fmt.Sprintln("  ", key+1, ".")
		for k, v := range value.ListAnd {
			resultPrint += fmt.Sprintln("    ", k+1, ".")
			resultPrint += fmt.Sprintf("      searchField: '%s'\n", v.SearchField)
			resultPrint += fmt.Sprintf("      searchValue: '%s'\n", v.SearchValue)
			resultPrint += fmt.Sprintf("      statementExpression: '%v'\n", v.StatementExpression)
		}
	}

	resultPrint += fmt.Sprintln("  EXCLUDE:")
	for key, value := range r.Rules.Exclude {
		resultPrint += fmt.Sprintln("  ", key+1, ".")
		for k, v := range value.ListAnd {
			resultPrint += fmt.Sprintln("    ", k+1, ".")
			resultPrint += fmt.Sprintf("      searchField: '%s'\n", v.SearchField)
			resultPrint += fmt.Sprintf("      searchValue: '%s'\n", v.SearchValue)
			resultPrint += fmt.Sprintf("      accurateComparison: '%v'\n", v.AccurateComparison)
		}
	}

	resultPrint += fmt.Sprintf("  PASSANY: '%v'\n", r.Rules.Passany)

	return resultPrint
}

func TestApp(t *testing.T) {
	//инициализация списка правил
	lr, lw, err = rules.NewListRule("placeholder_misp", "rules", "mispmsgrule.yml")

	//инициализация списка правил
	t.Log("NEW RULES FILE", fileName, ":")
	for k, v := range lw {
		t.Logf("%d. %s\n", k, v)
	}
	t.Log("new rule result:")
	t.Log(printRuleResult(lr))

	assert.NoError(t, err)
}
