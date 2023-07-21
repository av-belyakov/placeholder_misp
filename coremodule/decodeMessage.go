package coremodule

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	rules "placeholder_misp/rulesinteraction"
)

// ProcessMessageFromHive содержит информацию необходимую для обработки сообщений
// Message - сообщение полученное от Hive
// ListRules - список правил для обработки
// ListConfirmRules - список подтвержденных правил (ЗДЕСЬ НАДО ПОДУМАТЬ В КАКОМ ВИДЕ ЛУЧШЕ УЧИТЫВАТЬ ПРАВИЛА)
type ProcessMessageFromHive struct {
	Message          map[string]interface{}
	ListRules        rules.ListRulesProcMISPMessage
	ListConfirmRules []bool
}

func NewProcessMessageFromHive(b []byte, listRule rules.ListRulesProcMISPMessage) (ProcessMessageFromHive, error) {
	pmfh := ProcessMessageFromHive{}

	list := map[string]interface{}{}
	if err := json.Unmarshal(b, &list); err != nil {
		return pmfh, err
	}

	if len(list) == 0 {
		return pmfh, fmt.Errorf("error decoding the json file, it may be empty")
	}

	pmfh.Message = list

	var isStop bool
	newListRule := make([]rules.RuleProcMISPMessageField, 0, len(listRule.Rules))

	//обрабатываем список правил для исключения правил, которые заведомо не будут использоватся, то есть все правила после 'pass' и 'passany'
	for _, v := range listRule.Rules {
		if isStop {
			break
		}

		newListRule = append(newListRule, v)

		if v.ActionType == "pass" || v.ActionType == "passany" {
			isStop = true
		}
	}

	pmfh.ListRules.Rules = newListRule
	pmfh.ListConfirmRules = make([]bool, 0, len(listRule.Rules))

	return pmfh, nil
}

func (pm *ProcessMessageFromHive) ProcessMessage() bool {
	//пропускаем сообщение без какой либо обработки
	if len(pm.ListRules.Rules) > 0 && pm.ListRules.Rules[0].ActionType == "passany" {
		return true
	}

	/*

		единственная точка обработки всех правил это функция processingReflectAnySimpleType
		там стоит разделить типы действий для каждого правила

		нужно как то помечать какое правило для какого значения сработало

		reject (если совпадает, помечаем сообщение как отбрасываемое)
		pass (если совпадает, то помечаем как разрешонное к пропуску)
		replase (если совпадает, выполняем замену)

	*/

	newList, skipMsg := processingReflectMap(pm.Message, pm.ListRules, 0)

	pm.Message = newList

	return skipMsg
}

func (pm *ProcessMessageFromHive) GetMessage() ([]byte, error) {
	return json.Marshal(pm.Message)
}

func processingReflectAnySimpleType(
	name interface{},
	anyType interface{},
	listRule rules.ListRulesProcMISPMessage,
	num int) interface{} {

	var str, nameStr string
	r := reflect.TypeOf(anyType)

	if n, ok := name.(int); ok {
		nameStr = fmt.Sprintln(n)
	} else if n, ok := name.(string); ok {
		nameStr = n
	}

	if r == nil {
		return ""
	}

	/*
	   1. Проверяем совпадение reflect.ValueOf(anyType).String() или другого типа со значениями из listRule.SearchValuesName
	   2. Проверяем совпадение reflect.ValueOf(anyType).String() или другого типа со значениями из listRul.SearchFieldsName
	   3. Сравниваем индексы типа [0,0]
	   4. Если совпадает то получаем по этим индексам ActionType
	   5. Если не совпадает то получаем только по индексу из listRule.SearchValuesName и проверяем ActionType на равенство replace

	   для ActionType = pass, reject, RuleProcMISPMessageField.ListRequiredValues работает как 'И'
	   а для ActionType = replace ????
	*/

	switch r.Kind() {
	case reflect.String:
		//
		// ЭТО тестовая обработка поля "dataType", дальше обработку нужно сделать на основе правил
		//
		if strings.Contains(nameStr, "dataType") && strings.Contains(reflect.ValueOf(anyType).String(), "snort") {

			fmt.Printf("--- func 'processingReflectAnySimpleType' BEFORE ------%s = '%s'\n", nameStr, reflect.ValueOf(anyType).String())

			r := reflect.ValueOf(&anyType).Elem()
			fmt.Printf("reflect.ValueOf(&anyType).Elem() = %s\n", r)
			r.Set(reflect.ValueOf("TEST_SNORT_ID"))

			fmt.Printf("---- func 'processingReflectAnySimpleType' AFTER -----%s = '%s'\n", nameStr, reflect.ValueOf(anyType).String())

			return reflect.ValueOf(anyType).String()
		}

		return reflect.ValueOf(anyType).String()

	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		str += fmt.Sprintf("%s %d\n", nameStr, reflect.ValueOf(anyType).Int())

		return reflect.ValueOf(anyType).Int()

	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		str += fmt.Sprintf("%s %d\n", nameStr, reflect.ValueOf(anyType).Uint())

		return reflect.ValueOf(anyType).Uint()

	case reflect.Float32, reflect.Float64:
		str += fmt.Sprintf("%s %v\n", nameStr, int(reflect.ValueOf(anyType).Float()))

		return reflect.ValueOf(anyType).Float()

	case reflect.Bool:
		str += fmt.Sprintf("%s %v\n", nameStr, reflect.ValueOf(anyType).Bool())

		return reflect.ValueOf(anyType).Bool()
	}

	return ""
}

func processingReflectMap(
	l map[string]interface{},
	lr rules.ListRulesProcMISPMessage,
	num int) (map[string]interface{}, bool) {

	var (
		newMap  map[string]interface{}
		newList []interface{}
		skipMsg bool
	)
	nl := map[string]interface{}{}

	for k, v := range l {
		r := reflect.TypeOf(v)

		if r == nil {
			return nl, skipMsg
		}

		switch r.Kind() {
		case reflect.Map:
			if v, ok := v.(map[string]interface{}); ok {
				newMap, skipMsg = processingReflectMap(v, lr, num+1)
				nl[k] = newMap
			}

		case reflect.Slice:
			if v, ok := v.([]interface{}); ok {
				newList, skipMsg = processingReflectSlice(v, lr, num+1)
				nl[k] = newList
			}

		case reflect.Array:
			//str += fmt.Sprintf("%s: %s (it is array)\n", k, reflect.ValueOf(v).String())

		default:
			nl[k] = processingReflectAnySimpleType(k, v, lr, num)
		}

		if k == "dataType" {
			fmt.Println("func 'processingReflectMapTest', KKKK = ", k, " VVVV = ", l[k])
		}
	}

	return nl, skipMsg
}

func processingReflectSlice(
	l []interface{},
	lr rules.ListRulesProcMISPMessage,
	num int) ([]interface{}, bool) {

	var (
		newMap  map[string]interface{}
		newList []interface{}
		skipMsg bool
	)
	nl := make([]interface{}, 0, len(l))

	for k, v := range l {
		r := reflect.TypeOf(v)

		if r == nil {
			return nl, skipMsg
		}

		//l = append(l, processingReflectAnySimpleTypeTest(k, v, listRule, num))

		switch r.Kind() {
		case reflect.Map:
			if v, ok := v.(map[string]interface{}); ok {
				newMap, skipMsg = processingReflectMap(v, lr, num+1)

				nl = append(nl, newMap)
			}

		case reflect.Slice:
			if v, ok := v.([]interface{}); ok {
				newList, skipMsg = processingReflectSlice(v, lr, num+1)

				nl = append(nl, newList...)
			}

		case reflect.Array:
			//str += fmt.Sprintf("%d. %s (it is array)\n", k, reflect.ValueOf(v).String())

		default:
			nl = append(nl, processingReflectAnySimpleType(k, v, lr, num))
		}
	}

	return nl, skipMsg
}
