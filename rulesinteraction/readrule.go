package rules

import (
	"fmt"
	"path"
	"runtime"

	"github.com/spf13/viper"

	"placeholder_misp/supportingfunctions"
)

func GetRuleProcessingMsgForMISP(workDir, fn string) (ListRulesProcessingMsgMISP, []string, error) {
	lr := ListRulesProcessingMsgMISP{}

	fmt.Println("func 'GetRuleProcessingMsgForMISP', START. workDir = ", workDir, " fn = ", fn)

	rootPath, err := supportingfunctions.GetRootPath("placeholder_misp")
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return lr, []string{}, fmt.Errorf("%s %s:%d", fmt.Sprint(err), f, l+1)
	}

	viper.SetConfigFile(path.Join(rootPath, workDir, fn))
	viper.SetConfigType("yaml")
	err = viper.ReadInConfig()
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return lr, []string{}, fmt.Errorf("%s %s:%d", fmt.Sprint(err), f, l+1)
	}

	if ok := viper.IsSet("RULES"); !ok {
		_, f, l, _ := runtime.Caller(0)
		return lr, []string{}, fmt.Errorf("the 'RULES' property is missing in the file '%s' %s:%d", fn, f, l+1)
	}

	err = viper.GetViper().Unmarshal(&lr)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return lr, []string{}, fmt.Errorf("%s %s:%d", fmt.Sprint(err), f, l+1)
	}

	lr, warningCheckRules := verificationRules(lr)
	lr.RulesIndex = createIndexRules(lr.Rules)

	return lr, warningCheckRules, nil
}

func verificationRules(lr ListRulesProcessingMsgMISP) (ListRulesProcessingMsgMISP, []string) {
	replace, pass, passtestListAnd := []RuleReplace{}, []RulePass{}, []PasstestListAnd{}
	warning := []string{}

	for k, v := range lr.Rules.Pass {
		if v.SearchField == "" {
			warning = append(warning, fmt.Sprintf("warning: rule type 'PASS', number rule '%d', the 'searchField' property should not be empty", k))

			continue
		}

		if v.SearchValue == "" {
			warning = append(warning, fmt.Sprintf("warning: rule type 'PASS', number rule '%d', the 'searchValue' property should not be empty", k))

			continue
		}

		pass = append(pass, RulePass{SearchField: v.SearchField, SearchValue: v.SearchValue})
	}

	for k, v := range lr.Rules.Replace {
		if v.SearchValue == "" {
			warning = append(warning, fmt.Sprintf("warning: rule type 'REPLACE', number rule '%d', the 'searchValue' property should not be empty", k))

			continue
		}

		if v.ReplaceValue == "" {
			warning = append(warning, fmt.Sprintf("warning: rule type 'REPLACE', number rule '%d', the 'replaceValue' property should not be empty", k))

			continue
		}

		replace = append(replace, RuleReplace{
			SearchField:  v.SearchField,
			SearchValue:  v.SearchValue,
			ReplaceValue: v.ReplaceValue,
		})
	}

	/*
		ТЕСТОВЫЙ РАЗДЕЛ начало
	*/
	for _, value := range lr.Rules.Passtest {
		if len(value.ListAnd) == 0 {
			continue
		}

		listand := []RulePass{}

		for k, v := range value.ListAnd {
			if v.SearchField == "" {
				warning = append(warning, fmt.Sprintf("warning: rule type 'PASS', number rule '%d', the 'searchField' property should not be empty", k))

				continue
			}

			if v.SearchValue == "" {
				warning = append(warning, fmt.Sprintf("warning: rule type 'PASS', number rule '%d', the 'searchValue' property should not be empty", k))

				continue
			}

			listand = append(listand, RulePass{SearchField: v.SearchField, SearchValue: v.SearchValue})
		}

		passtestListAnd = append(passtestListAnd, PasstestListAnd{ListAnd: listand})
	}
	/*
		ТЕСТОВЫЙ РАЗДЕЛ конец
	*/

	return ListRulesProcessingMsgMISP{
		Rules: RuleSetProcessingMsgMISP{
			Passany:  lr.Rules.Passany,
			Pass:     pass,
			Replace:  replace,
			Passtest: passtestListAnd,
		}}, warning
}

func createIndexRules(rules RuleSetProcessingMsgMISP) map[string][]RuleIndex {
	ri := map[string][]RuleIndex{}

	for _, v := range rules.Pass {
		if _, ok := ri[v.SearchValue]; !ok {
			ri[v.SearchValue] = []RuleIndex{}
		}

		ri[v.SearchValue] = append(ri[v.SearchValue], RuleIndex{
			RuleType:    "PASS",
			SearchField: v.SearchField,
		})
	}

	for _, v := range rules.Replace {
		if _, ok := ri[v.SearchValue]; !ok {
			ri[v.SearchValue] = []RuleIndex{}
		}

		ri[v.SearchValue] = append(ri[v.SearchValue], RuleIndex{
			RuleType:     "REPLACE",
			SearchField:  v.SearchField,
			ReplaceValue: v.ReplaceValue,
		})
	}

	return ri
}

/*
func searchStr(l []string, d string) (int, bool) {
	var (
		i  int
		ok bool
	)

	for k, v := range l {
		if d == v {
			i, ok = k, true

			break
		}
	}

	return i, ok
}
*/
