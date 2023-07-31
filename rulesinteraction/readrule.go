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

	fmt.Println("func 'GetRuleProcessingMsgForMISP', work dir = ", workDir, " file name = ", fn)

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

	return lr, warningCheckRules, nil
}

func verificationRules(lr ListRulesProcessingMsgMISP) (ListRulesProcessingMsgMISP, []string) {
	replace, passListAnd := []RuleReplace{}, []PassListAnd{}
	warning := []string{}

	for _, value := range lr.Rules.Pass {
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

		passListAnd = append(passListAnd, PassListAnd{ListAnd: listand})
	}

	for k, v := range lr.Rules.Replace {
		if v.SearchField == "" && v.SearchValue == "" {
			warning = append(warning, fmt.Sprintf("warning: rule type 'REPLACE', number rule '%d', one of the properties 'searchField' or 'searchValue' must be filled in", k))

			continue
		}

		replace = append(replace, RuleReplace{
			SearchField:  v.SearchField,
			SearchValue:  v.SearchValue,
			ReplaceValue: v.ReplaceValue,
		})
	}

	return ListRulesProcessingMsgMISP{
		Rules: RuleSetProcessingMsgMISP{
			Passany: lr.Rules.Passany,
			Pass:    passListAnd,
			Replace: replace,
		}}, warning
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
