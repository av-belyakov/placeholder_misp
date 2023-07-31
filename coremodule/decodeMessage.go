package coremodule

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

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

func (pm *ProcessMessageFromHive) HandleMessage(chanOutMispFormat chan<- ChanInputCreateMispFormat) (bool, []string) {
	var skipMsg bool
	warningMsg := []string{}
	chanWarningMsg := make(chan string)

	go func(cwmsg <-chan string, wmsg []string) {
		wmsg = append(wmsg, <-cwmsg)
	}(chanWarningMsg, warningMsg)

	pm.Message = processingReflectMap(chanWarningMsg, chanOutMispFormat, pm.Message, pm.ListRules, 0)
	close(chanWarningMsg)

	/* ДЛЯ ТЕСТА */
	fmt.Println("func 'HandleMessage', 22222 pm.ListRules.Rules.Passtest = ", pm.ListRules.Rules.Pass)
	fmt.Println("func 'HandleMessage', pm.Rules.Pass: ")
	for key, value := range pm.ListRules.Rules.Pass {
		fmt.Printf("%d.\n", key)
		for k, v := range value.ListAnd {
			fmt.Printf("  %d.\n", k)
			fmt.Printf("  SearchField: %s\n", v.SearchField)
			fmt.Printf("  SearchValue: %s\n", v.SearchValue)
			fmt.Printf("  StatementExpression: %t\n", v.StatementExpression)
		}
	}

	//сообщение пропускается в независимости от результата обработки правил PASS
	if pm.ListRules.Rules.Passany {
		skipMsg = true
	}

	skipMsg = true
	//проверяем соответствие сообщения правилам из раздела Pass
	for _, v := range pm.ListRules.Rules.Pass {
		for _, value := range v.ListAnd {
			if !value.StatementExpression {
				skipMsg = false

				break
			}
		}

		if skipMsg {
			break
		}
	}

	/*
		   Протестировал pm.ListRules.Rules.Pass, если
		   					SearchField: dataType
		     				SearchValue: i2p_111home

		   				то получаем skipMsg так как нет совпадения, если
		         			SearchField: "dataType"
		         			SearchValue: "ip_home"

					Вроде все работает
	*/

	fmt.Println("")
	fmt.Println("___________________func 'HandleMessage' ____________________")
	fmt.Println("-------- TEST skipMsg:", skipMsg, " ---------------")
	fmt.Println("____________________________________________________________")

	return skipMsg, warningMsg
}

func (pm *ProcessMessageFromHive) GetMessage() ([]byte, error) {
	return json.Marshal(pm.Message)
}

func PassRuleHandler(rulePass []rules.PassListAnd, fn string, cv interface{}) {
	cvstr := fmt.Sprint(cv)

	for key, value := range rulePass {
		for k, v := range value.ListAnd {
			if fn != v.SearchField {
				continue
			}

			if cvstr != v.SearchValue {
				continue
			}

			/*
				с strings.Contains не работает почему то
				зато с if cvstr != v.SearchValue работае

				if !strings.Contains(v.SearchValue, cvstr) {
					continue
				}*/

			rulePass[key].ListAnd[k].StatementExpression = true
		}
	}
}

func ReplacementRuleHandler(lr rules.ListRulesProcessingMsgMISP, svt, fn string, cv interface{}) (interface{}, int, error) {
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
		if v.SearchValue != fmt.Sprint(cv) {
			continue
		}

		if v.SearchField != "" {
			if v.SearchField == fn {
				nv, err := getReplaceValue(svt, v.ReplaceValue)

				return nv, k, err
			}

			continue
		}

		nv, err := getReplaceValue(svt, v.ReplaceValue)

		return nv, k, err
	}

	return cv, 0, nil
}

func processingReflectAnySimpleType(
	wmsg chan<- string,
	chanOutMispFormat chan<- ChanInputCreateMispFormat,
	name interface{},
	anyType interface{},
	listRule rules.ListRulesProcessingMsgMISP,
	num int) interface{} {

	var nameStr string
	r := reflect.TypeOf(anyType)

	if n, ok := name.(int); ok {
		nameStr = fmt.Sprintln(n)
	} else if n, ok := name.(string); ok {
		nameStr = n
	}

	if r == nil {
		return anyType
	}

	switch r.Kind() {
	case reflect.String:
		result := reflect.ValueOf(anyType).String()

		ncv, num, err := ReplacementRuleHandler(listRule, "string", nameStr, result)
		if err != nil {
			wmsg <- fmt.Sprintf("search value '%s' from rule number '%d' of section 'REPLACE' is not fulfilled", result, num)
		}

		PassRuleHandler(listRule.Rules.Pass, nameStr, ncv)

		chanOutMispFormat <- ChanInputCreateMispFormat{
			FieldName: nameStr,
			ValueType: "string",
			Value:     ncv,
		}

		return ncv

	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		result := reflect.ValueOf(anyType).Int()

		ncv, num, err := ReplacementRuleHandler(listRule, "int", nameStr, result)
		if err != nil {
			wmsg <- fmt.Sprintf("search value '%v' from rule number '%d' of section 'REPLACE' is not fulfilled", result, num)
		}

		PassRuleHandler(listRule.Rules.Pass, nameStr, ncv)

		chanOutMispFormat <- ChanInputCreateMispFormat{
			FieldName: nameStr,
			ValueType: "int",
			Value:     ncv,
		}

		return ncv

	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		result := reflect.ValueOf(anyType).Uint()

		ncv, num, err := ReplacementRuleHandler(listRule, "uint", nameStr, result)
		if err != nil {
			wmsg <- fmt.Sprintf("search value '%v' from rule number '%d' of section 'REPLACE' is not fulfilled", result, num)
		}

		PassRuleHandler(listRule.Rules.Pass, nameStr, ncv)

		chanOutMispFormat <- ChanInputCreateMispFormat{
			FieldName: nameStr,
			ValueType: "uint",
			Value:     ncv,
		}

		return ncv

	case reflect.Float32, reflect.Float64:
		result := reflect.ValueOf(anyType).Float()

		ncv, num, err := ReplacementRuleHandler(listRule, "float", nameStr, result)
		if err != nil {
			wmsg <- fmt.Sprintf("search value '%v' from rule number '%d' of section 'REPLACE' is not fulfilled", result, num)
		}

		PassRuleHandler(listRule.Rules.Pass, nameStr, ncv)

		chanOutMispFormat <- ChanInputCreateMispFormat{
			FieldName: nameStr,
			ValueType: "float",
			Value:     ncv,
		}

		return ncv

	case reflect.Bool:
		result := reflect.ValueOf(anyType).Bool()

		ncv, num, err := ReplacementRuleHandler(listRule, "bool", nameStr, result)
		if err != nil {
			wmsg <- fmt.Sprintf("search value '%v' from rule number '%d' of section 'REPLACE' is not fulfilled", result, num)
		}

		PassRuleHandler(listRule.Rules.Pass, nameStr, ncv)

		chanOutMispFormat <- ChanInputCreateMispFormat{
			FieldName: nameStr,
			ValueType: "bool",
			Value:     ncv,
		}

		return ncv
	}

	return anyType
}

func processingReflectMap(
	wmsg chan<- string,
	chanOutMispFormat chan<- ChanInputCreateMispFormat,
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
				newMap = processingReflectMap(wmsg, chanOutMispFormat, v, lr, num+1)
				nl[k] = newMap
			}

		case reflect.Slice:
			if v, ok := v.([]interface{}); ok {
				newList = processingReflectSlice(wmsg, chanOutMispFormat, v, lr, num+1)
				nl[k] = newList
			}

		case reflect.Array:
			//str += fmt.Sprintf("%s: %s (it is array)\n", k, reflect.ValueOf(v).String())

		default:
			nl[k] = processingReflectAnySimpleType(wmsg, chanOutMispFormat, k, v, lr, num)
		}
	}

	return nl
}

func processingReflectSlice(
	wmsg chan<- string,
	chanOutMispFormat chan<- ChanInputCreateMispFormat,
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
				newMap = processingReflectMap(wmsg, chanOutMispFormat, v, lr, num+1)

				nl = append(nl, newMap)
			}

		case reflect.Slice:
			if v, ok := v.([]interface{}); ok {
				newList = processingReflectSlice(wmsg, chanOutMispFormat, v, lr, num+1)

				nl = append(nl, newList...)
			}

		case reflect.Array:
			//str += fmt.Sprintf("%d. %s (it is array)\n", k, reflect.ValueOf(v).String())

		default:
			nl = append(nl, processingReflectAnySimpleType(wmsg, chanOutMispFormat, k, v, lr, num))
		}
	}

	return nl
}
