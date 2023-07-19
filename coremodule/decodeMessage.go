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
// ListConfirmRules - список подтвержденных правил
type ProcessMessageFromHive struct {
	Message          map[string]interface{}
	ListRules        []rules.RuleProcMISPMessageField
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
	newListRule := make([]rules.RuleProcMISPMessageField, 0, len(listRule.Rulles))

	//обрабатываем список правил для исключения правил, которые заведомо не будут использоватся, то есть все правила после 'pass' и 'passany'
	for _, v := range listRule.Rulles {
		if isStop {
			break
		}

		newListRule = append(newListRule, v)

		if v.ActionType == "pass" || v.ActionType == "passany" {
			isStop = true
		}
	}

	pmfh.ListRules = newListRule
	pmfh.ListConfirmRules = make([]bool, 0, len(listRule.Rulles))

	return pmfh, nil
}

func (pm *ProcessMessageFromHive) ProcessMessage() bool {
	//пропускаем сообщение без какой либо обработки
	if len(pm.ListRules) > 0 && pm.ListRules[0].ActionType == "passany" {
		return true
	}

	return false
}

func (pm *ProcessMessageFromHive) GetMessage() ([]byte, error) {
	return json.Marshal(pm.Message)
}

// NewProcessingInputMessageFromHive выполняет обработку сообщений на основе правил
func NewProcessingInputMessageFromHive(b []byte, listRule rules.ListRulesProcMISPMessage) ([]byte, []bool, error) {
	//информирует, какие правила сработали, а какие нет (правила учитываются по их порядковому номеру)
	listConfirmRules := make([]bool, len(listRule.Rulles))
	list := map[string]interface{}{}

	if err := json.Unmarshal(b, &list); err != nil {
		return []byte{}, listConfirmRules, err
	}

	if len(list) == 0 {
		return []byte{}, listConfirmRules, fmt.Errorf("error decoding the json file, it may be empty")
	}

	/*
		Если правило есть,
		   2. Проверить, если есть наличие правила где 'actionType' = 'anypass', то пропускаем обработку всех правил после данного правила,
		   а json объект пропускаем если в правиле 'requiredValue.value' пустое или если оно соответствует указанному значению
		   которое должно находится в поле json объекта с именем хранящемся в 'fieldName'
		   3. Обработать правила 'pass', 'passany', 'replace', 'reject'
	*/

	fmt.Println("---------- func 'NewProcessingInputMessageFromHive', Read list rules -----------")
	for k, v := range listRule.Rulles {
		fmt.Println(k, ". ")
		fmt.Println("  actionType: ", v.ActionType)
		for i, j := range v.ListRequiredValues {
			fmt.Println("    ", i, ". ")
			fmt.Println("      fieldName: ", j.FieldName)
			fmt.Println("      typeValue: ", j.TypeValue)
			fmt.Println("      value: ", j.Value)
			fmt.Println("      replaceValue: ", j.ReplaceValue)
		}
	}
	fmt.Println("--------------------------------------------------------------------------------")

	// если первое правило типа 'passany'
	if len(listRule.Rulles) > 0 && listRule.Rulles[0].ActionType == "passany" {
		listConfirmRules = append(listConfirmRules, true)
		result, err := json.Marshal(list)

		return result, listConfirmRules, err
	}

	/*
		1.
	*/

	newList := processingReflectMap(list, listRule, 0)
	result, err := json.Marshal(newList)

	return result, listConfirmRules, err
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

func processingReflectMap(list map[string]interface{}, listRule rules.ListRulesProcMISPMessage, num int) map[string]interface{} {
	l := map[string]interface{}{}

	for k, v := range list {
		r := reflect.TypeOf(v)

		if r == nil {
			return l
		}

		switch r.Kind() {
		case reflect.Map:
			if v, ok := v.(map[string]interface{}); ok {
				l[k] = processingReflectMap(v, listRule, num+1)
			}

		case reflect.Slice:
			if v, ok := v.([]interface{}); ok {
				l[k] = processingReflectSlice(v, listRule, num+1)
			}

		case reflect.Array:
			//str += fmt.Sprintf("%s: %s (it is array)\n", k, reflect.ValueOf(v).String())

		default:
			l[k] = processingReflectAnySimpleType(k, v, listRule, num)
		}

		if k == "dataType" {
			fmt.Println("func 'processingReflectMapTest', KKKK = ", k, " VVVV = ", l[k])
		}
	}

	return l
}

func processingReflectSlice(list []interface{}, listRule rules.ListRulesProcMISPMessage, num int) []interface{} {
	l := make([]interface{}, 0, len(list))

	for k, v := range list {
		r := reflect.TypeOf(v)

		if r == nil {
			return l
		}

		//l = append(l, processingReflectAnySimpleTypeTest(k, v, listRule, num))

		switch r.Kind() {
		case reflect.Map:
			if v, ok := v.(map[string]interface{}); ok {
				l = append(l, processingReflectMap(v, listRule, num+1))
			}

		case reflect.Slice:
			if v, ok := v.([]interface{}); ok {
				l = append(l, processingReflectSlice(v, listRule, num+1))
			}

		case reflect.Array:
			//str += fmt.Sprintf("%d. %s (it is array)\n", k, reflect.ValueOf(v).String())

		default:
			l = append(l, processingReflectAnySimpleType(k, v, listRule, num))
		}
	}

	return l
}
