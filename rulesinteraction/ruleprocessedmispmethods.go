package rules

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// NewListRule создает новый список правил
func NewListRule(rootDir, workDir, fileName string) (*ListRule, []string, error) {
	lr := ListRule{}

	_, f, l, _ := runtime.Caller(0)
	rootPath, err := getRootPath(rootDir)
	if err != nil {
		return &lr, []string{}, fmt.Errorf("'%s' %s:%d", err.Error(), f, l+1)
	}

	viper.SetConfigFile(path.Join(rootPath, workDir, fileName))
	viper.SetConfigType("yaml")

	_, f, l, _ = runtime.Caller(0)
	err = viper.ReadInConfig()
	if err != nil {
		return &lr, []string{}, fmt.Errorf("'%s' %s:%d", err.Error(), f, l+1)
	}

	_, f, l, _ = runtime.Caller(0)
	if ok := viper.IsSet("RULES"); !ok {
		return &lr, []string{}, fmt.Errorf("'the \"RULES\" property is missing in the file \"%s\"' %s:%d", fileName, f, l+1)
	}

	_, f, l, _ = runtime.Caller(0)
	err = viper.GetViper().Unmarshal(&lr, func(dc *mapstructure.DecoderConfig) {
		dc.Squash = true
	})
	if err != nil {
		return &lr, []string{}, fmt.Errorf("'%s' %s:%d", err.Error(), f, l+1)
	}

	warningCheckRules := lr.verification()

	return &lr, warningCheckRules, nil
}

// ReplacementRuleHandler выполняет замену значения, на значение, из поля replaceValue
// файла правил, если свойство fieldName равно содержимому поля searchField правила, а
// свойство currentValue совпадает со значением поля searchValue правила
func (lr *ListRule) ReplacementRuleHandler(searchValueType, fieldName string, currentValue interface{}) (interface{}, int, error) {
	getReplaceValue := func(svt, rv string) (interface{}, error) {
		switch svt {
		case "string":
			return rv, nil

		case "int":
			return strconv.ParseInt(rv, 10, 64)

		case "uint":
			return strconv.ParseUint(rv, 10, 64)

		case "float":
			return strconv.ParseFloat(rv, 64)

		case "bool":
			return strconv.ParseBool(rv)
		}

		return rv, nil
	}

	for k, v := range lr.Rules.Replace {
		if v.SearchValue != fmt.Sprint(currentValue) {
			continue
		}

		if v.SearchField != "" {
			if v.SearchField == fieldName {
				nv, err := getReplaceValue(searchValueType, v.ReplaceValue)

				return nv, k, err
			}

			continue
		}

		nv, err := getReplaceValue(searchValueType, v.ReplaceValue)

		return nv, k, err
	}

	return currentValue, 0, nil
}

// PassRuleHandler выполняет сравнение имени поля из свойства fieldName и поля searchField,
// а также значения свойства currentValue и поля searchValue правила Pass. При совпадении
// этих значений изменяется состояние поля StatementExpression соответствующего правил
// на true
func (lr *ListRule) PassRuleHandler(fieldName string, currentValue interface{}) {
	cvstr := fmt.Sprint(currentValue)

	for key, value := range lr.Rules.Pass {
		for k, v := range value.ListAnd {
			if fieldName != v.SearchField {
				continue
			}

			if strings.Contains(v.SearchValue, "not:") && v.SearchValue[:4] == "not:" {
				if cvstr == v.SearchValue[4:] {
					continue
				}

				lr.Rules.Pass[key].ListAnd[k].StatementExpression = true
			} else {
				if cvstr != v.SearchValue {
					continue
				}

				lr.Rules.Pass[key].ListAnd[k].StatementExpression = true
			}
		}
	}
}

// CleanStatementExpressionRulePass приводит поле StatementExpression к значению false
// это поле проверяется на соответствие заданным правилам обрабатываемым значениям
// обязательно нужно выполнять данный метод после проверки значения StatementExpression
func (lr *ListRule) CleanStatementExpressionRulePass() {
	for k, v := range lr.Rules.Pass {
		for key := range v.ListAnd {
			lr.Rules.Pass[k].ListAnd[key].StatementExpression = false
		}
	}
}

// SomePassRuleIsTrue выполняет проверку на совпадение хотя бы одного блока
// 'listAnd' из правил Pass
func (lr *ListRule) SomePassRuleIsTrue() bool {
	list := make([]bool, len(lr.Rules.Pass))

	for k, v := range lr.Rules.Pass {
		skipMsg := true
		for _, value := range v.ListAnd {
			if !value.StatementExpression {
				skipMsg = false

				break
			}
		}
		//для каждого блока ListAnd свое значение, так как внутри блока
		// действует правило 'И'
		list[k] = skipMsg
	}

	for _, v := range list {
		if v {
			//если есть хотя бы одно положительное значение, так как между блоками
			// ListAnd действует правило 'ИЛИ'
			return true
		}
	}

	return false
}

// ExcludeRuleHandler выполняет сравнение имени поля из свойства fieldName и поля searchField,
// а также значения свойства currentValue и поля searchValue правила Exclude. При этом учитывается
// состояние поля AccurateComparison, если его значение 'true', то тогда применяется 'строгое'
// сравнение, то есть содержимое currentValue должно в точности соответствовать содержимому
// поля searchValue. А если состояние поля AccurateComparison 'false', то тогда currentValue
// должно содержать в себе значение из поля searchValue вместе с любыми другими значениями.
// При совпадения значений, функция возвращает 'true'.
func (lr *ListRule) ExcludeRuleHandler(fieldName string, currentValue interface{}) bool {
	cvstr := fmt.Sprint(currentValue)

	for _, v := range lr.Rules.Exclude {
		if fieldName != v.SearchField {
			continue
		}

		//fmt.Println("func 'ExcludeRuleHandler'")
		//fmt.Println("fieldName: ", fieldName, " currentValue:", cvstr)

		if v.AccurateComparison && cvstr == v.SearchValue {
			return true
		}

		if !v.AccurateComparison && strings.Contains(cvstr, v.SearchValue) {
			return true
		}
	}

	return false
}

// GetRulePass возвращает список правил типа Pass
func (lr *ListRule) GetRulePass() []PassListAnd {
	return lr.Rules.Pass
}

// GetRulePassany возвращает значения правила Passany
func (lr *ListRule) GetRulePassany() bool {
	return lr.Rules.Passany
}

// GetRuleExclude  возвращает список правил типа Exclude
func (lr *ListRule) GetRuleExclude() []RuleExclude {
	return lr.Rules.Exclude
}

func getRootPath(rootDir string) (string, error) {
	currentDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}

	tmp := strings.Split(currentDir, "/")

	if tmp[len(tmp)-1] == rootDir {
		return currentDir, nil
	}

	var path string = ""
	for _, v := range tmp {
		path += v + "/"

		if v == rootDir {
			return path, nil
		}
	}

	return path, nil
}

func (lr *ListRule) verification() []string {
	rr, rp, re := []RuleReplace{}, []PassListAnd{}, []RuleExclude{}
	warning := []string{}

	//проверяем правила типа 'Pass'
	for _, value := range lr.Rules.Pass {
		if len(value.ListAnd) == 0 {
			continue
		}

		listAnd := []RulePass{}
		for k, v := range value.ListAnd {
			k += 1
			if v.SearchField == "" {
				warning = append(warning, fmt.Sprintf("warning: rule type 'PASS', number rule '%d', the 'searchField' property should not be empty", k))

				continue
			}

			if v.SearchValue == "" {
				warning = append(warning, fmt.Sprintf("warning: rule type 'PASS', number rule '%d', the 'searchValue' property should not be empty", k))

				continue
			}

			listAnd = append(listAnd, RulePass{
				CommonRuleFields: CommonRuleFields{
					SearchField: v.SearchField,
					SearchValue: v.SearchValue,
				},
			})
		}

		rp = append(rp, PassListAnd{ListAnd: listAnd})
	}

	//проверяем правила типа 'Replace'
	for k, v := range lr.Rules.Replace {
		k += 1
		if v.SearchField == "" && v.SearchValue == "" {
			warning = append(warning, fmt.Sprintf("warning: rule type 'REPLACE', number rule '%d', one of the properties 'searchField' or 'searchValue' must be filled in", k))

			continue
		}

		rr = append(rr, RuleReplace{
			CommonRuleFields: CommonRuleFields{
				SearchField: v.SearchField,
				SearchValue: v.SearchValue,
			},
			ReplaceValue: v.ReplaceValue,
		})
	}

	//проверяем правила типа 'Exclude'
	for k, v := range lr.Rules.Exclude {
		k += 1
		if v.SearchField == "" {
			warning = append(warning, fmt.Sprintf("warning: rule type 'EXCLUDE', number rule '%d', the 'searchField' property should not be empty", k))

			continue
		}

		re = append(re, RuleExclude{
			CommonRuleFields: CommonRuleFields{
				SearchField: v.SearchField,
				SearchValue: v.SearchValue,
			},
			AccurateComparison: v.AccurateComparison,
		})
	}

	if len(rp) == 0 && !lr.Rules.Passany {
		warning = append(warning, fmt.Sprintf("warning: rule type 'PASSANY' is '%v', however rule type 'PASS' is empty too", lr.Rules.Passany))
	}

	lr.Rules.Pass = rp
	lr.Rules.Replace = rr
	lr.Rules.Exclude = re

	return warning
}
