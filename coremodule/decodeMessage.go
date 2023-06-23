package coremodule

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"placeholder_misp/rules"
)

// NewProcessingInputMessageFromHive выполняет обработку сообщений на основе правил
func NewProcessingInputMessageFromHive(b []byte, listRule rules.ListRulesProcessedMISPMessage) ([]byte, bool, error) {
	var isSuccess bool
	list := map[string]interface{}{}

	if err := json.Unmarshal(b, &list); err != nil {
		return []byte{}, isSuccess, err
	}

	if len(list) == 0 {
		return []byte{}, isSuccess, fmt.Errorf("error decoding the json file, it may be empty")
	}

	/*
	   1. Проверить наличие правил в listRule, если правил нет, то правила не обрабатываем, а json объект пропускаем сразу
	   2. Проверить, если есть наличие правила где 'actionType' = 'anypass', то пропускаем обработку всех правил после данного правила,
	   а json объект пропускаем если в правиле 'requiredValue.value' пустое или если оно соответствует указанному значению
	   которое должно находится в поле json объекта с именем хранящемся в 'fieldName'
	   3. Обработать правила с где 'actionType' = 'process', 'replace' или 'clean'.
	*/

	newList := processingReflectMap(list, listRule, 0)

	fmt.Println("---------- Read list rules -----------")
	fmt.Println("func 'NewProcessingInputMessageFromHive', listRule:", listRule)
	fmt.Println("--------------------------------------")

	//
	// Сделал тестовое измененние значение некоторых полей
	// теперь нужно сделать логику работы с правилами обработки полей
	//

	//fmt.Println("func NewProcessingInputMessageFromHive, new list = ", newList)
	result, err := json.Marshal(newList)

	return result, isSuccess, err
}

func processingReflectAnySimpleType(
	name interface{},
	anyType interface{},
	listRule rules.ListRulesProcessedMISPMessage,
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

func processingReflectMap(list map[string]interface{}, listRule rules.ListRulesProcessedMISPMessage, num int) map[string]interface{} {
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

func processingReflectSlice(list []interface{}, listRule rules.ListRulesProcessedMISPMessage, num int) []interface{} {
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

/*-------- ТЕСТОВЫЕ ФУНКЦИИ --------*/
/*
	ЭТА ФУНКЦИЯ НЕ РАБОТАЕТ ДОЛЖНЫМ ОБРАЗОМ
func NewProcessingInputMessageFromHiveTest(b []byte, listRule rules.ListRulesProcessedMISPMessage) ([]byte, error) {
	list := map[string]interface{}{}

	//
	//Не могу сохранить или протащить по ссылке данные которые изменил в processingReflectAnySimpleTypeTest
	// в первую очередь сторку
	//
	if err := json.Unmarshal(b, &list); err != nil {
		return []byte{}, err
	}

	if len(list) == 0 {
		return []byte{}, fmt.Errorf("error decoding the json file, it may be empty")
	}

	processingReflectMapTest(list, listRule, 0)

	return json.Marshal(list)
}

func processingReflectAnySimpleTypeTest(
	name interface{},
	anyType *interface{},
	listRule rules.ListRulesProcessedMISPMessage,
	num int) {

	var str, nameStr string
	r := reflect.TypeOf(*anyType)

	if n, ok := name.(int); ok {
		nameStr = fmt.Sprintln(n)
	} else if n, ok := name.(string); ok {
		nameStr = n
	}

	if r == nil {
		return
	}

	//if strings.Contains(nameStr, "dataType") {
	//	fmt.Println("func 'processingReflectAnySimpleType' nameStr = ", nameStr, " reflect.ValueOf(anyType).String() = ", reflect.ValueOf(*anyType).String())
	//	fmt.Println("--- strings.Contains(reflect.ValueOf(*anyType).String(), 'snort') == ", strings.Contains(reflect.ValueOf(*anyType).String(), "snort"))
	//}

	switch r.Kind() {
	case reflect.String:
		if strings.Contains(nameStr, "dataType") && strings.Contains(reflect.ValueOf(*anyType).String(), "snort") {

			fmt.Printf("--- func 'processingReflectAnySimpleType' BEFORE ------%s = '%s'\n", nameStr, reflect.ValueOf(*anyType).String())

			// r := reflect.ValueOf(&anyType).Elem()
			r := reflect.ValueOf(anyType).Elem()
			fmt.Printf("reflect.ValueOf(&anyType).Elem() = %s\n", r)
			r.Set(reflect.ValueOf("TEST_SNORT_ID"))

			fmt.Printf("---- func 'processingReflectAnySimpleType' AFTER -----%s = '%s'\n", nameStr, reflect.ValueOf(*anyType).String())
		}

	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		str += fmt.Sprintf("%s %d\n", nameStr, reflect.ValueOf(*anyType).Int())

	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		str += fmt.Sprintf("%s %d\n", nameStr, reflect.ValueOf(*anyType).Uint())

	case reflect.Float32, reflect.Float64:
		str += fmt.Sprintf("%s %v\n", nameStr, int(reflect.ValueOf(*anyType).Float()))

	case reflect.Bool:
		str += fmt.Sprintf("%s %v\n", nameStr, reflect.ValueOf(*anyType).Bool())
	}
}

func processingReflectMapTest(list map[string]interface{}, listRule rules.ListRulesProcessedMISPMessage, num int) {

	//fmt.Printf("func 'processingReflectMapTest', LIST = '%T'\n", list)

	newList := make(map[string]interface{}, len(list))

	for k, v := range list {
		r := reflect.TypeOf(v)

		//fmt.Println("____________-- ", (*list)[k])

		if r == nil {
			continue
		}

		processingReflectAnySimpleTypeTest(k, &v, listRule, num)

		newList[k] = v

		if k == "dataType" {
			fmt.Println("func 'processingReflectMapTest', KKKK = ", k, " VVVV = ", v)
		}

		switch r.Kind() {
		case reflect.Map:
			if v, ok := v.(map[string]interface{}); ok {
				//if v, ok := (*list)[k].(*map[string]interface{}); ok {
				//processingReflectMapTest(v, listRule, num+1)
				processingReflectMapTest(v, listRule, num+1)

				newList[k] = v

				for i, j := range v {
					if i == "dataType" {
						fmt.Println("FFFFFFF :", j)
					}
				}
			}

		case reflect.Slice:
			if v, ok := v.([]interface{}); ok {
				//if v, ok := (*list)[k].(*[]interface{}); ok {
				//processingReflectSliceTest(v, listRule, num+1)
				processingReflectSliceTest(&v, listRule, num+1)

				newList[k] = v
			}

		case reflect.Array:
			//str += fmt.Sprintf("%s: %s (it is array)\n", k, reflect.ValueOf(v).String())
		}
	}

	for k, v := range newList {
		if k == "dataType" {
			fmt.Println("JJJJJJJJ :", v)
		}
	}

	list = newList
}

func processingReflectSliceTest(list *[]interface{}, listRule rules.ListRulesProcessedMISPMessage, num int) {
	for k, v := range *list {
		//for i := 0; i < len(*list); i++ {
		r := reflect.TypeOf(v)
		//r := reflect.TypeOf((*list)[i])

		if r == nil {
			return
		}

		processingReflectAnySimpleTypeTest(k, &v, listRule, num)
		//processingReflectAnySimpleTypeTest(i, &(*list)[i], listRule, num)

		switch r.Kind() {
		case reflect.Map:
			if v, ok := v.(map[string]interface{}); ok {
				//if v, ok := (*list)[i].(*map[string]interface{}); ok {
				//processingReflectMapTest(v, listRule, num+1)
				processingReflectMapTest(v, listRule, num+1)
			}

		case reflect.Slice:
			if v, ok := v.([]interface{}); ok {
				//if v, ok := (*list)[i].(*[]interface{}); ok {
				//processingReflectSliceTest(v, listRule, num+1)
				processingReflectSliceTest(&v, listRule, num+1)
			}

		case reflect.Array:
			//str += fmt.Sprintf("%d. %s (it is array)\n", k, reflect.ValueOf(v).String())
		}
	}
}*/
