package coremodule

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	rules "placeholder_misp/rulesinteraction"
)

// ProcessMessageFromHive содержит информацию необходимую для обработки сообщений
// Message - сообщение полученное от Hive
// ListRules - список правил для обработки
// ListConfirmRules - список подтвержденных правил (ЗДЕСЬ НАДО ПОДУМАТЬ В КАКОМ ВИДЕ ЛУЧШЕ УЧИТЫВАТЬ ПРАВИЛА)
type ProcessMessageFromHive struct {
	Message          map[string]interface{}
	ListRules        rules.ListRulesProcessingMsgMISP
	ListConfirmRules []bool
}

func NewHandleMessageFromHive(b []byte, listRule rules.ListRulesProcessingMsgMISP) (ProcessMessageFromHive, error) {
	pmfh := ProcessMessageFromHive{}

	fmt.Println("func 'NewHandleMessageFromHive', START")

	list := map[string]interface{}{}
	if err := json.Unmarshal(b, &list); err != nil {
		return pmfh, err
	}

	if len(list) == 0 {
		return pmfh, fmt.Errorf("error decoding the json file, it may be empty")
	}

	pmfh.Message = list
	pmfh.ListRules = listRule

	return pmfh, nil
}

func (pm *ProcessMessageFromHive) HandleMessage() (bool, []string) {
	var skipMsg bool

	warningMsg := []string{}
	listMissed := [][]bool{}
	chanWarningMsg := make(chan string)

	fmt.Println("func 'HandleMessage', START")

	//если правило PASSANY false, а совпадений в PASS нет то сообщение отбрасывается

	for _, v := range pm.ListRules.Rules.Passtest {
		listMissed = append(listMissed, make([]bool, len(v.ListAnd)))
	}

	go func(cwmsg <-chan string, wmsg *[]string) {
		*wmsg = append(*wmsg, <-cwmsg)
	}(chanWarningMsg, &warningMsg)

	fmt.Println("func 'HandleMessage', 11111")

	pm.Message = processingReflectMap(chanWarningMsg, pm.Message, pm.ListRules, 0)
	close(chanWarningMsg)

	fmt.Println("func 'HandleMessage', 22222")

	//сообщение пропускается в независимости от результата обработки правил PASS
	if pm.ListRules.Rules.Passany {
		skipMsg = true
	}

	for _, l := range listMissed {
		var isTrue bool = true

		for _, v := range l {
			if !v {
				isTrue = false
			}
		}

		if isTrue {
			skipMsg = true

			break
		}
	}

	fmt.Println("___________________func 'HandleMessage' ____________________")
	fmt.Println("test listMissed: ", listMissed, " skipMsg:", skipMsg)

	return skipMsg, warningMsg
}

func (pm *ProcessMessageFromHive) GetMessage() ([]byte, error) {
	return json.Marshal(pm.Message)
}

func ReplacementRuleHandler(svt, rv string, sv, cv interface{}) (interface{}, error) {
	searchValue, ok := sv.(string)
	if !ok {
		return cv, fmt.Errorf("unable to convert value 'searchValue'")
	}

	switch svt {
	case "string":
		cv, ok := cv.(string)
		if !ok {
			break
		}

		if sv == cv {
			return rv, nil
		}

	case "int":
		cv, ok := cv.(int64)
		if !ok {
			break
		}

		sv, err := strconv.ParseInt(rv, 10, 64)
		if err != nil {
			return cv, err
		}

		if sv == cv {
			return strconv.ParseInt(rv, 10, 64)
		}

	case "uint":
		cv, ok := cv.(uint64)
		if !ok {
			break
		}

		sv, err := strconv.ParseUint(searchValue, 10, 64)
		if err != nil {
			return cv, err
		}

		if sv == cv {
			return strconv.ParseUint(rv, 10, 64)
		}

	case "float":
		cv, ok := cv.(float64)
		if !ok {
			break
		}

		sv, err := strconv.ParseFloat(searchValue, 64)
		if err != nil {
			return cv, err
		}

		if sv == cv {
			return strconv.ParseFloat(rv, 64)
		}

	case "bool":
		cv, ok := cv.(bool)
		if !ok {
			break
		}

		sv, err := strconv.ParseBool(searchValue)
		if err != nil {
			return cv, err
		}

		if sv == cv {
			return strconv.ParseBool(rv)
		}
	}

	return cv, nil
}

func processingReflectAnySimpleType(
	wmsg chan<- string,
	name interface{},
	anyType interface{},
	listRule rules.ListRulesProcessingMsgMISP,
	num int) interface{} {

	var (
		err     error
		nameStr string
	)
	r := reflect.TypeOf(anyType)

	if n, ok := name.(int); ok {
		nameStr = fmt.Sprintln(n)
	} else if n, ok := name.(string); ok {
		nameStr = n
	}

	if r == nil {
		return anyType
	}

	fmt.Println("func 'processingReflectAnySimpleType' nameStr: ", nameStr, " AnyType: ", anyType, " TypeOf: ", reflect.TypeOf(anyType))

	switch r.Kind() {
	case reflect.String:
		result := reflect.ValueOf(anyType).String()

		/*

			Похоже все числовые значения преобразовываются в FLOAT32, FLOAT64
			с индексами как доп. список правил нельзя использовать null и поиск
			пустых значений

			Надо все же попробовать обрабатывать обычный список правил REPLACE
			и сделать обработку правил PASS по образу PASSTEST
		*/

		//может так надо попробовать обрабатыват правила из раздела REPLACE???
		for k, v := range listRule.Rules.Replace {
			if result != v.SearchValue {
				continue
			}

			var t interface{}
			if v.SearchField != "" {
				if strings.Contains(v.SearchField, nameStr) {
					t, err = ReplacementRuleHandler("string", v.ReplaceValue, v.SearchValue, result)
				}
			} else {
				t, err = ReplacementRuleHandler("string", v.ReplaceValue, v.SearchValue, result)
			}

			if err != nil {
				wmsg <- fmt.Sprintf("search value '%s' from rule number '%d' of section 'REPLACE' is not fulfilled", result, k)
			}

			if s, ok := t.(string); ok {
				result = s

				//при совпадении искомого значения стоит ли продалжать поиск и повторную замену
				//или выйти из цикла???
				break
			}
		}

		/*indexs, ok := listRule.RulesIndex[result]
		if !ok {
			return result
		}

		fmt.Println("func 'processingReflectAnySimpleType' STRING searchValue: ", result, " added settings: ", indexs)

		for k, index := range indexs {
			//if index.RuleType == "PASS" {}

			if index.RuleType == "REPLACE" {
				var t interface{}

				if index.SearchField != "" {
					if strings.Contains(nameStr, index.SearchField) {
						t, err = ReplacementRuleHandler("string", index.ReplaceValue, result, result)
					}
				} else {
					t, err = ReplacementRuleHandler("string", index.ReplaceValue, result, result)
				}

				if err != nil {
					//ВОТ ТУТ НАДО ПОДУМАТЬ КАК ОБРАБАТЫВАТЬ ОШИБКУ
					wmsg <- fmt.Sprintf("search value '%s' from rule number '%d' of section '%s' is not fulfilled", result, k, index.RuleType)
				}

				if s, ok := t.(string); ok {
					result = s
				}
			}
		}*/

		return result

	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:

		fmt.Println("func 'processingReflectAnySimpleType' INT, INT16, INT32, INT64 searchValue: |||||||||||||||||||||")

		result := reflect.ValueOf(anyType).Int()
		indexs, ok := listRule.RulesIndex[fmt.Sprintln(result)]
		if !ok {
			return result
		}

		fmt.Println("func 'processingReflectAnySimpleType' INT, INT16, INT32, INT64 searchValue: ", result, " added settings: ", indexs)

		for k, index := range indexs {
			//if index.RuleType == "PASS" {}

			if index.RuleType == "REPLACE" {
				var t interface{}

				if index.SearchField != "" {
					if strings.Contains(nameStr, index.SearchField) {
						t, err = ReplacementRuleHandler("int", index.ReplaceValue, fmt.Sprintln(result), result)
					}
				} else {
					t, err = ReplacementRuleHandler("int", index.ReplaceValue, fmt.Sprintln(result), result)
				}

				if err != nil {
					//ВОТ ТУТ НАДО ПОДУМАТЬ КАК ОБРАБАТЫВАТЬ ОШИБКУ
					wmsg <- fmt.Sprintf("search value '%v' from rule number '%d' of section '%s' is not fulfilled", result, k, index.RuleType)
				}

				if s, ok := t.(int64); ok {
					result = s
				}
			}
		}

		return result

	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		result := reflect.ValueOf(anyType).Uint()

		fmt.Println("func 'processingReflectAnySimpleType' UINT, UINT16, UINT32, UINT64 searchValue: |||||||||||||||||||||")

		indexs, ok := listRule.RulesIndex[fmt.Sprintln(result)]
		if !ok {
			return result
		}

		fmt.Println("func 'processingReflectAnySimpleType' UINT, UINT16, UINT32, UINT64 searchValue: ", result, " added settings: ", indexs)

		for k, index := range indexs {
			//if index.RuleType == "PASS" {}

			if index.RuleType == "REPLACE" {
				var t interface{}

				if index.SearchField != "" {
					if strings.Contains(nameStr, index.SearchField) {
						t, err = ReplacementRuleHandler("uint", index.ReplaceValue, fmt.Sprintln(result), result)
					}
				} else {
					t, err = ReplacementRuleHandler("uint", index.ReplaceValue, fmt.Sprintln(result), result)
				}

				if err != nil {
					//ВОТ ТУТ НАДО ПОДУМАТЬ КАК ОБРАБАТЫВАТЬ ОШИБКУ
					wmsg <- fmt.Sprintf("search value '%v' from rule number '%d' of section '%s' is not fulfilled", result, k, index.RuleType)
				}

				if s, ok := t.(uint64); ok {
					result = s
				}
			}
		}

		return result

	case reflect.Float32, reflect.Float64:
		result := reflect.ValueOf(anyType).Float()

		fmt.Println("func 'processingReflectAnySimpleType' FLOAT32, FLOAT64 searchValue: ||||||||||||||||||||| listRule.RulesIndex = ", listRule.RulesIndex)
		fmt.Println("float value: ", fmt.Sprintln(result))

		/*

			Похоже все числовые значения преобразовываются в FLOAT32, FLOAT64
			с индексами как доп. список правил нельзя использовать null и поиск
			пустых значений

			Надо все же попробовать обрабатывать обычный список правил REPLACE
			и сделать обработку правил PASS по образу PASSTEST
		*/

		indexs, ok := listRule.RulesIndex[fmt.Sprintln(result)]
		if !ok {

			fmt.Println("OOOOOOOOOOOOOOOOoo")

			return result
		}

		fmt.Println("func 'processingReflectAnySimpleType' FLOAT32, FLOAT64 searchValue: ", result, " added settings: ", indexs)

		for k, index := range indexs {
			//if index.RuleType == "PASS" {}

			if index.RuleType == "REPLACE" {
				var t interface{}

				if index.SearchField != "" {
					if strings.Contains(nameStr, index.SearchField) {
						t, err = ReplacementRuleHandler("float", index.ReplaceValue, fmt.Sprintln(result), result)
					}
				} else {
					t, err = ReplacementRuleHandler("float", index.ReplaceValue, fmt.Sprintln(result), result)
				}

				if err != nil {
					//ВОТ ТУТ НАДО ПОДУМАТЬ КАК ОБРАБАТЫВАТЬ ОШИБКУ
					wmsg <- fmt.Sprintf("search value '%v' from rule number '%d' of section '%s' is not fulfilled", result, k, index.RuleType)
				}

				if s, ok := t.(float64); ok {
					result = s
				}
			}
		}

		return result

		/*case reflect.Bool:
		result := reflect.ValueOf(anyType).Bool()

		sv, err := strconv.ParseBool(searchValue)
		if err != nil {
			wmsg <- fmt.Sprintf("rule error, search value '%s' of section 'Replace', %v", searchValue, err)

			return result
		}

		if result != sv {
			return result
		}

		for k, index := range indexs {
			//if index.RuleType == "PASS" {}

			if index.RuleType == "REPLACE" {
				var t interface{}

				if index.SearchField != "" {
					if strings.Contains(nameStr, index.SearchField) {
						t, err = ReplacementRuleHandler("bool", index.ReplaceValue, searchValue, result)
					}
				} else {
					t, err = ReplacementRuleHandler("bool", index.ReplaceValue, searchValue, result)
				}

				if err != nil {
					//ВОТ ТУТ НАДО ПОДУМАТЬ КАК ОБРАБАТЫВАТЬ ОШИБКУ
					wmsg <- fmt.Sprintf("search value '%s' from rule number '%d' of section '%s' is not fulfilled", searchValue, k, index.RuleType)
				}

				if s, ok := t.(bool); ok {
					result = s
				}
			}
		}

		return result*/
	}

	/*switch r.Kind() {
	case reflect.String:

		//
		// ЭТО тестовая обработка поля "dataType", дальше обработку нужно сделать на основе правил
		//
		//if strings.Contains(nameStr, "dataType") && strings.Contains(reflect.ValueOf(anyType).String(), "snort") {
		//
		//	fmt.Printf("--- func 'processingReflectAnySimpleType' BEFORE ------%s = '%s'\n", nameStr, reflect.ValueOf(anyType).String())
		//
		//	r := reflect.ValueOf(&anyType).Elem()
		//	fmt.Printf("reflect.ValueOf(&anyType).Elem() = %s\n", r)
		//	r.Set(reflect.ValueOf("TEST_SNORT_ID"))
		//
		//	fmt.Printf("---- func 'processingReflectAnySimpleType' AFTER -----%s = '%s'\n", nameStr, reflect.ValueOf(anyType).String())
		//
		//	return reflect.ValueOf(anyType).String()
		//}

		result := reflect.ValueOf(anyType).String()

		for k, v := range listRule.Rules.Replace {
			var t interface{}

			if v.SearchField != "" {
				if strings.Contains(nameStr, v.SearchField) {
					t, err = ReplacementRuleHandler("string", v.ReplaceValue, v.SearchValue, reflect.ValueOf(anyType).String())
				}
			} else {
				t, err = ReplacementRuleHandler("string", v.ReplaceValue, v.SearchValue, reflect.ValueOf(anyType).String())
			}

			if err != nil {
				//ВОТ ТУТ НАДО ПОДУМАТЬ КАК ОБРАБАТЫВАТЬ ОШИБКУ
				wmsg <- fmt.Sprintf("rule number '%d' of section 'Replace' is not fulfilled", k)
			}

			if s, ok := t.(string); ok {
				result = s

				break
			}
		}

		// 			!!!!
		// Надо написать обработку раздела PASS правил и все потестировать
		// 			!!!!

		return result

	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		result := reflect.ValueOf(anyType).Int()

		for k, v := range listRule.Rules.Replace {
			var t interface{}

			if v.SearchField != "" {
				if strings.Contains(nameStr, v.SearchField) {
					t, err = ReplacementRuleHandler("int", v.ReplaceValue, v.SearchValue, reflect.ValueOf(anyType).Int())
				}
			} else {
				t, err = ReplacementRuleHandler("int", v.ReplaceValue, v.SearchValue, reflect.ValueOf(anyType).Int())
			}

			if err != nil {
				wmsg <- fmt.Sprintf("rule number '%d' of section 'Replace' is not fulfilled", k)
			}

			if s, ok := t.(int64); ok {
				result = s

				break
			}
		}

		return result

	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		result := reflect.ValueOf(anyType).Uint()

		for k, v := range listRule.Rules.Replace {
			var t interface{}

			if v.SearchField != "" {
				if strings.Contains(nameStr, v.SearchField) {
					t, err = ReplacementRuleHandler("uint", v.ReplaceValue, v.SearchValue, reflect.ValueOf(anyType).Uint())
				}
			} else {
				t, err = ReplacementRuleHandler("uint", v.ReplaceValue, v.SearchValue, reflect.ValueOf(anyType).Uint())
			}

			if err != nil {
				wmsg <- fmt.Sprintf("rule number '%d' of section 'Replace' is not fulfilled", k)
			}

			if s, ok := t.(uint64); ok {
				result = s

				break
			}
		}

		return result

	case reflect.Float32, reflect.Float64:
		result := reflect.ValueOf(anyType).Float()

		for k, v := range listRule.Rules.Replace {
			var t interface{}

			if v.SearchField != "" {
				if strings.Contains(nameStr, v.SearchField) {
					t, err = ReplacementRuleHandler("float", v.ReplaceValue, v.SearchValue, reflect.ValueOf(anyType).Float())
				}
			} else {
				t, err = ReplacementRuleHandler("float", v.ReplaceValue, v.SearchValue, reflect.ValueOf(anyType).Float())
			}

			if err != nil {
				wmsg <- fmt.Sprintf("rule number '%d' of section 'Replace' is not fulfilled", k)
			}

			if s, ok := t.(float64); ok {
				result = s

				break
			}
		}

		return result

	case reflect.Bool:
		result := reflect.ValueOf(anyType).Bool()

		for k, v := range listRule.Rules.Replace {
			var t interface{}

			if v.SearchField != "" {
				if strings.Contains(nameStr, v.SearchField) {
					t, err = ReplacementRuleHandler("bool", v.ReplaceValue, v.SearchValue, reflect.ValueOf(anyType).Bool())
				}
			} else {
				t, err = ReplacementRuleHandler("bool", v.ReplaceValue, v.SearchValue, reflect.ValueOf(anyType).Bool())
			}

			if err != nil {
				wmsg <- fmt.Sprintf("rule number '%d' of section 'Replace' is not fulfilled", k)
			}

			if s, ok := t.(bool); ok {
				result = s

				break
			}
		}

		return result
	}*/

	return anyType
}

func processingReflectMap(
	wmsg chan<- string,
	l map[string]interface{},
	lr rules.ListRulesProcessingMsgMISP,
	num int) map[string]interface{} {

	var (
		newMap  map[string]interface{}
		newList []interface{}
	)
	nl := map[string]interface{}{}

	for k, v := range l {
		r := reflect.TypeOf(v)

		if r == nil {
			return nl
		}

		switch r.Kind() {
		case reflect.Map:
			if v, ok := v.(map[string]interface{}); ok {
				newMap = processingReflectMap(wmsg, v, lr, num+1)
				nl[k] = newMap
			}

		case reflect.Slice:
			if v, ok := v.([]interface{}); ok {
				newList = processingReflectSlice(wmsg, v, lr, num+1)
				nl[k] = newList
			}

		case reflect.Array:
			//str += fmt.Sprintf("%s: %s (it is array)\n", k, reflect.ValueOf(v).String())

		default:
			nl[k] = processingReflectAnySimpleType(wmsg, k, v, lr, num)
		}

		if k == "dataType" {
			fmt.Println("func 'processingReflectMapTest', KKKK = ", k, " VVVV = ", l[k])
		}
	}

	return nl
}

func processingReflectSlice(
	wmsg chan<- string,
	l []interface{},
	lr rules.ListRulesProcessingMsgMISP,
	num int) []interface{} {

	var (
		newMap  map[string]interface{}
		newList []interface{}
	)
	nl := make([]interface{}, 0, len(l))

	for k, v := range l {
		r := reflect.TypeOf(v)

		if r == nil {
			return nl
		}

		//l = append(l, processingReflectAnySimpleTypeTest(k, v, listRule, num))

		switch r.Kind() {
		case reflect.Map:
			if v, ok := v.(map[string]interface{}); ok {
				newMap = processingReflectMap(wmsg, v, lr, num+1)

				nl = append(nl, newMap)
			}

		case reflect.Slice:
			if v, ok := v.([]interface{}); ok {
				newList = processingReflectSlice(wmsg, v, lr, num+1)

				nl = append(nl, newList...)
			}

		case reflect.Array:
			//str += fmt.Sprintf("%d. %s (it is array)\n", k, reflect.ValueOf(v).String())

		default:
			nl = append(nl, processingReflectAnySimpleType(wmsg, k, v, lr, num))
		}
	}

	return nl
}
