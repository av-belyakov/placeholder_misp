package rules

import (
	"fmt"
	"path"
	"runtime"

	"github.com/spf13/viper"

	"placeholder_misp/supportingfunctions"
)

// GetRuleProcessedMISPMsg получает правила обработки сообщений из файла конфигурации и выполняет их верификацию
// принимает workDir - рабочую директорию и имя файла
// возвращает список верифицированых правил, список предупреждений о возникших при верификации и ошибку
func GetRuleProcessedMISPMsg(workDir, fn string) (ListRulesProcMISPMessage, []string, error) {
	lrp := ListRulesProcMISPMessage{}

	rootPath, err := supportingfunctions.GetRootPath("placeholder_misp")
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return lrp, []string{}, fmt.Errorf("%s %s:%d", fmt.Sprint(err), f, l+1)
	}

	viper.SetConfigFile(path.Join(rootPath, workDir, fn))
	viper.SetConfigType("yaml")
	err = viper.ReadInConfig()
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return lrp, []string{}, fmt.Errorf("%s %s:%d", fmt.Sprint(err), f, l+1)
	}

	if ok := viper.IsSet("RULLES"); !ok {
		_, f, l, _ := runtime.Caller(0)
		return lrp, []string{}, fmt.Errorf("the 'RULLES' property is missing in the file '%s' %s:%d", fn, f, l+1)
	}

	err = viper.GetViper().Unmarshal(&lrp)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return lrp, []string{}, fmt.Errorf("%s %s:%d", fmt.Sprint(err), f, l+1)
	}

	//выполняем анализ конфигурационного файла (проверяем правильность сформированной структуры)
	lrp, warningCheckRules := verificationRules(lrp)

	return lrp, warningCheckRules, nil
}

func verificationRules(lrp ListRulesProcMISPMessage) (ListRulesProcMISPMessage, []string) {
	rp, vr := []RuleProcMISPMessageField{}, []string{}
	lat := []string{"pass", "passany", "reject", "replace"}
	ltv := []string{"bool", "int", "string"}

	for k, v := range lrp.Rulles {
		_, ok := searchStr(lat, v.ActionType)
		if !ok {
			vr = append(vr, fmt.Sprintf("warning: number rule '%d', the 'actionType' property contains an invalid value '%s'", k, v.ActionType))

			continue
		}

		if v.ActionType != "passany" && len(v.ListRequiredValues) == 0 {
			vr = append(vr, fmt.Sprintf("warning: number rule '%d', the 'listRequiredValues' property should not be empty", k))

			continue
		}

		lrv := []ListRequiredValue{}
		for key, value := range v.ListRequiredValues {
			if value.FieldSearchName == "" && v.ActionType != "replace" {
				vr = append(vr, fmt.Sprintf("warning: number rule '%d.%d', the 'fieldSearchName' property should not be empty", k, key))

				continue
			}

			_, ok := searchStr(ltv, value.TypeValue)
			if !ok {
				vr = append(vr, fmt.Sprintf("warning: number rule '%d.%d', the 'typeValue' property contains an invalid value '%s'", k, key, value.TypeValue))

				continue
			}

			if value.FindValue == "" {
				vr = append(vr, fmt.Sprintf("warning: number rule '%d.%d', missing 'findValue' property, to indicate an empty value for this property, use the value 'null'", k, key))

				continue
			}

			if value.ReplaceValue == "" && v.ActionType == "replace" {
				vr = append(vr, fmt.Sprintf("warning: number rule '%d.%d', missing 'replaceValue' property, to indicate an empty value for this property, use the value 'null'", k, key))

				continue
			}

			lrv = append(lrv, ListRequiredValue{
				FieldSearchName: value.FieldSearchName,
				TypeValue:       value.TypeValue,
				FindValue:       value.FindValue,
				ReplaceValue:    value.ReplaceValue,
			})
		}

		if v.ActionType != "passany" && len(lrv) == 0 {
			continue
		}

		rp = append(rp, RuleProcMISPMessageField{
			ActionType:         v.ActionType,
			ListRequiredValues: lrv,
		})
	}

	return ListRulesProcMISPMessage{Rulles: rp}, vr
}

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
