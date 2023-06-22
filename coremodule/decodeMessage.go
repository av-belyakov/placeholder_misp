package coremodule

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"placeholder_misp/rules"
	"placeholder_misp/supportingfunctions"
)

// NewProcessingInputMessageFromHive выполняет обработку сообщений на основе правил
func NewProcessingInputMessageFromHive(b *[]byte, listRule rules.ListRulesProcessedMISPMessage) error {
	list := map[string]interface{}{}

	//
	// тут проблемма из-за которой не изменяется данные b, потому что в Unmarshal передется копия данных
	//
	if err := json.Unmarshal(*b, &list); err != nil {
		return err
	}

	if len(list) == 0 {
		return fmt.Errorf("error decoding the json file, it may be empty")
	}

	processingReflectMap(&list, listRule, 0)

	return nil
}

func processingReflectAnySimpleType(
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

	/*if strings.Contains(nameStr, "dataType") {
		fmt.Println("func 'processingReflectAnySimpleType' nameStr = ", nameStr, " reflect.ValueOf(anyType).String() = ", reflect.ValueOf(*anyType).String())
		fmt.Println("--- strings.Contains(reflect.ValueOf(*anyType).String(), 'snort') == ", strings.Contains(reflect.ValueOf(*anyType).String(), "snort"))
	}*/

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

func processingReflectMap(list *map[string]interface{}, listRule rules.ListRulesProcessedMISPMessage, num int) {
	for k, v := range *list {
		r := reflect.TypeOf(v)

		if r == nil {
			return
		}

		processingReflectAnySimpleType(k, &v, listRule, num)

		switch r.Kind() {
		case reflect.Map:
			if v, ok := v.(map[string]interface{}); ok {
				processingReflectMap(&v, listRule, num+1)
			}

		case reflect.Slice:
			if v, ok := v.([]interface{}); ok {
				processingReflectSlice(&v, listRule, num+1)
			}

		case reflect.Array:
			//str += fmt.Sprintf("%s: %s (it is array)\n", k, reflect.ValueOf(v).String())
		}
	}
}

func processingReflectSlice(list *[]interface{}, listRule rules.ListRulesProcessedMISPMessage, num int) {
	for k, v := range *list {
		r := reflect.TypeOf(v)

		if r == nil {
			return
		}

		processingReflectAnySimpleType(k, &v, listRule, num)

		switch r.Kind() {
		case reflect.Map:
			if v, ok := v.(map[string]interface{}); ok {
				processingReflectMap(&v, listRule, num+1)
			}

		case reflect.Slice:
			if v, ok := v.([]interface{}); ok {
				processingReflectSlice(&v, listRule, num+1)
			}

		case reflect.Array:
			//str += fmt.Sprintf("%d. %s (it is array)\n", k, reflect.ValueOf(v).String())
		}
	}
}

/*-------- ТЕСТОВЫЕ ФУНКЦИИ --------*/
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

	processingReflectMapTest(&list, listRule, 0)

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

func processingReflectMapTest(list *map[string]interface{}, listRule rules.ListRulesProcessedMISPMessage, num int) {

	//	fmt.Println("func 'processingReflectMapTest', LIST = ", *list)

	for k, v := range *list {
		r := reflect.TypeOf(v)

		//fmt.Println("____________-- ", (*list)[k])

		if r == nil {
			return
		}

		processingReflectAnySimpleTypeTest(k, &v, listRule, num)

		if k == "dataType" {
			fmt.Println("func 'processingReflectMapTest', KKKK = ", k, " VVVV = ", v)
		}

		switch r.Kind() {
		case reflect.Map:
			//if v, ok := v.(map[string]interface{}); ok {
			if v, ok := (*list)[k].(map[string]interface{}); ok {
				processingReflectMapTest(&v, listRule, num+1)
			}

		case reflect.Slice:
			//			if v, ok := v.([]interface{}); ok {
			if v, ok := (*list)[k].([]interface{}); ok {
				processingReflectSliceTest(&v, listRule, num+1)
			}

		case reflect.Array:
			//str += fmt.Sprintf("%s: %s (it is array)\n", k, reflect.ValueOf(v).String())
		}
	}
}

func processingReflectSliceTest(list *[]interface{}, listRule rules.ListRulesProcessedMISPMessage, num int) {
	for k, v := range *list {
		r := reflect.TypeOf(v)

		if r == nil {
			return
		}

		processingReflectAnySimpleTypeTest(k, &v, listRule, num)

		switch r.Kind() {
		case reflect.Map:
			if v, ok := v.(map[string]interface{}); ok {
				processingReflectMapTest(&v, listRule, num+1)
			}

		case reflect.Slice:
			if v, ok := v.([]interface{}); ok {
				processingReflectSliceTest(&v, listRule, num+1)
			}

		case reflect.Array:
			//str += fmt.Sprintf("%d. %s (it is array)\n", k, reflect.ValueOf(v).String())
		}
	}
}

/*func NewProcessingInputMessageFromHiveTest(b []byte, listRule rules.ListRulesProcessedMISPMessage) ([]byte, error) {
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

	newList := processingReflectMapTest(list, listRule, 0)

	fmt.Println("func NewProcessingInputMessageFromHiveTest, new list = ", newList)

	return json.Marshal(newList)
}

func processingReflectAnySimpleTypeTest(
	name interface{},
	anyType interface{},
	listRule rules.ListRulesProcessedMISPMessage,
	num int) reflect.Value {

	var str, nameStr string
	r := reflect.TypeOf(anyType)

	if n, ok := name.(int); ok {
		nameStr = fmt.Sprintln(n)
	} else if n, ok := name.(string); ok {
		nameStr = n
	}

	if r == nil {
		return reflect.ValueOf(anyType)
	}

	//if strings.Contains(nameStr, "dataType") {
	//	fmt.Println("func 'processingReflectAnySimpleType' nameStr = ", nameStr, " reflect.ValueOf(anyType).String() = ", reflect.ValueOf(*anyType).String())
	//	fmt.Println("--- strings.Contains(reflect.ValueOf(*anyType).String(), 'snort') == ", strings.Contains(reflect.ValueOf(*anyType).String(), "snort"))
	//}

	switch r.Kind() {
	case reflect.String:
		if strings.Contains(nameStr, "dataType") && strings.Contains(reflect.ValueOf(anyType).String(), "snort") {

			fmt.Printf("--- func 'processingReflectAnySimpleType' BEFORE ------%s = '%s'\n", nameStr, reflect.ValueOf(anyType).String())

			// r := reflect.ValueOf(&anyType).Elem()
			r := reflect.ValueOf(&anyType).Elem()
			fmt.Printf("reflect.ValueOf(&anyType).Elem() = %s\n", r)
			r.Set(reflect.ValueOf("TEST_SNORT_ID"))

			fmt.Printf("---- func 'processingReflectAnySimpleType' AFTER -----%s = '%s'\n", nameStr, reflect.ValueOf(anyType).String())

			return r
		}

	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		str += fmt.Sprintf("%s %d\n", nameStr, reflect.ValueOf(anyType).Int())

	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		str += fmt.Sprintf("%s %d\n", nameStr, reflect.ValueOf(anyType).Uint())

	case reflect.Float32, reflect.Float64:
		str += fmt.Sprintf("%s %v\n", nameStr, int(reflect.ValueOf(anyType).Float()))

	case reflect.Bool:
		str += fmt.Sprintf("%s %v\n", nameStr, reflect.ValueOf(anyType).Bool())
	}

	return reflect.ValueOf(&anyType).Elem()
}

func processingReflectMapTest(list map[string]interface{}, listRule rules.ListRulesProcessedMISPMessage, num int) map[string]interface{} {
	l := map[string]interface{}{}

	//	fmt.Println("func 'processingReflectMapTest', LIST = ", *list)

	for k, v := range list {
		r := reflect.TypeOf(v)

		if r == nil {
			return l
		}

		l[k] = processingReflectAnySimpleTypeTest(k, v, listRule, num)

		if k == "dataType" {
			fmt.Println("func 'processingReflectMapTest', KKKK = ", k, " VVVV = ", l[k])
		}

		switch r.Kind() {
		case reflect.Map:
			if v, ok := v.(map[string]interface{}); ok {
				l[k] = processingReflectMapTest(v, listRule, num+1)
			}

		case reflect.Slice:
			if v, ok := v.([]interface{}); ok {
				l[k] = processingReflectSliceTest(v, listRule, num+1)
			}

		case reflect.Array:
			//str += fmt.Sprintf("%s: %s (it is array)\n", k, reflect.ValueOf(v).String())
		}
	}

	return l
}

func processingReflectSliceTest(list []interface{}, listRule rules.ListRulesProcessedMISPMessage, num int) []interface{} {
	l := make([]interface{}, 0, len(list))

	for k, v := range list {
		r := reflect.TypeOf(v)

		if r == nil {
			return l
		}

		l = append(l, processingReflectAnySimpleTypeTest(k, v, listRule, num))

		switch r.Kind() {
		case reflect.Map:
			if v, ok := v.(map[string]interface{}); ok {
				l = append(l, processingReflectMapTest(v, listRule, num+1))
			}

		case reflect.Slice:
			if v, ok := v.([]interface{}); ok {
				l = append(l, processingReflectSliceTest(v, listRule, num+1))
			}

		case reflect.Array:
			//str += fmt.Sprintf("%d. %s (it is array)\n", k, reflect.ValueOf(v).String())
		}
	}

	return l
}*/

/* дальше неиспользуемые функции */

func readReflectAnyTypeSprint(name, anyType interface{}, lpmfmisp rules.ListRulesProcessedMISPMessage, num int) string {
	var str, nameStr string
	r := reflect.TypeOf(anyType)
	ws := supportingfunctions.GetWhitespace(num)

	if n, ok := name.(int); ok {
		nameStr = fmt.Sprintf("%s%v.", ws, n+1)
	} else if n, ok := name.(string); ok {
		nameStr = fmt.Sprintf("%s%s:", ws, n)
	}

	if r == nil {
		return str
	}

	switch r.Kind() {
	case reflect.String:
		if strings.Contains(nameStr, "dataType") && strings.Contains(reflect.ValueOf(anyType).String(), "url") {

			fmt.Println("func 'ReadReflectAnyTypeSprint', type String, nameStr:", nameStr, " reflect.ValueOf(anyType).String():", reflect.ValueOf(anyType).String())

			r := reflect.ValueOf(&anyType).Elem()
			fmt.Printf("reflect.ValueOf(&anyType).Elem() = %s", r)
			//r.Set(reflect.ValueOf("TEST_SNORT_ID"))
			r.Set(reflect.ValueOf("TEST_URL_URL_TEST"))

			/*

				Урааа!!!
				Получилось заменить содержимое поля "dataType"="snort" на "TEST_SNORT_ID"

				Тепер эту функцию нужно полностью переписать для модификации данных,
				а для просмотра использовать функцию NewReadReflectJSONSprint из пакета supportingfunctions.readReflectJSONSprint

			*/

		}

		str += fmt.Sprintf("%s '%s'\n", nameStr, reflect.ValueOf(anyType).String())

	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		str += fmt.Sprintf("%s %d\n", nameStr, reflect.ValueOf(anyType).Int())

	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		str += fmt.Sprintf("%s %d\n", nameStr, reflect.ValueOf(anyType).Uint())

	case reflect.Float32, reflect.Float64:
		str += fmt.Sprintf("%s %v\n", nameStr, int(reflect.ValueOf(anyType).Float()))

	case reflect.Bool:
		str += fmt.Sprintf("%s %v\n", nameStr, reflect.ValueOf(anyType).Bool())
	}

	return str
}

func ReadReflectMapSprint(list map[string]interface{}, lpmfmisp rules.ListRulesProcessedMISPMessage, num int) string {
	var str string
	ws := supportingfunctions.GetWhitespace(num)

	for k, v := range list {
		r := reflect.TypeOf(v)

		if r == nil {
			return str
		}

		str += readReflectAnyTypeSprint(k, v, lpmfmisp, num)

		switch r.Kind() {
		case reflect.Map:
			if v, ok := v.(map[string]interface{}); ok {
				str += fmt.Sprintf("%s%s:\n", ws, k)
				//str += fmt.Sprintln(k, ":")
				str += ReadReflectMapSprint(v, lpmfmisp, num+1)
			}

		case reflect.Slice:
			if v, ok := v.([]interface{}); ok {
				str += fmt.Sprintf("%s%s:\n", ws, k)
				str += readReflectSliceSprint(v, lpmfmisp, num+1)
			}

		case reflect.Array:
			str += fmt.Sprintf("%s: %s (it is array)\n", k, reflect.ValueOf(v).String())
		}
	}

	return str
}

func readReflectSliceSprint(list []interface{}, lpmfmisp rules.ListRulesProcessedMISPMessage, num int) string {
	var str string

	for k, v := range list {
		r := reflect.TypeOf(v)

		if r == nil {
			return str
		}

		str += readReflectAnyTypeSprint(k, v, lpmfmisp, num)

		switch r.Kind() {
		case reflect.Map:
			if v, ok := v.(map[string]interface{}); ok {
				str += ReadReflectMapSprint(v, lpmfmisp, num+1)
			}

		case reflect.Slice:
			if v, ok := v.([]interface{}); ok {
				str += readReflectSliceSprint(v, lpmfmisp, num+1)
			}

		case reflect.Array:
			str += fmt.Sprintf("%d. %s (it is array)\n", k, reflect.ValueOf(v).String())
		}
	}

	return str
}
