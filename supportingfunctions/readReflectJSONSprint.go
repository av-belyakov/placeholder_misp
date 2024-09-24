package supportingfunctions

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// NewReadReflectJSONSprint функция выполняет вывод JSON сообщения в виде текста
// Для данной функции не требуется описание текст, так как обработка JSON сообщения
// осуществляется с помощью пакета reflect
func NewReadReflectJSONSprint(b []byte) (string, error) {
	var str string
	errSrc := "error decoding the json file, it may be empty"

	listMap := map[string]interface{}{}
	if err := json.Unmarshal(b, &listMap); err == nil {
		if len(listMap) == 0 {
			return str, errors.New(errSrc)
		}

		return readReflectMapSprint(listMap, 0), err
	}

	listSlice := []interface{}{}
	if err := json.Unmarshal(b, &listSlice); err == nil {
		if len(listSlice) == 0 {
			return str, errors.New(errSrc)
		}

		return readReflectSliceSprint(listSlice, 0), err
	}

	return str, fmt.Errorf("the contents of the file are not Map or Slice")
}

func readReflectAnyTypeSprint(name interface{}, anyType interface{}, num int) string {
	var (
		nameStr string
		str     strings.Builder = strings.Builder{}

		isCleanLine bool
	)

	r := reflect.TypeOf(anyType)
	ws := GetWhitespace(num)

	if n, ok := name.(int); ok {
		nameStr = fmt.Sprintf("%s%v.", ws, n+1)
	} else if n, ok := name.(string); ok {
		nameStr = fmt.Sprintf("%s\"%s\":", ws, n)

		if n == "description" {
			isCleanLine = true
		}
	}

	if r == nil {
		return str.String()
	}

	switch r.Kind() {
	case reflect.String:
		dataStr := reflect.ValueOf(anyType).String()

		if isCleanLine {
			dataStr = strings.ReplaceAll(dataStr, "\t", "")
			dataStr = strings.ReplaceAll(dataStr, "\n", "")
		}

		str.WriteString(fmt.Sprintf("%s \"%s\"\n", nameStr, dataStr))

	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		str.WriteString(fmt.Sprintf("%s %d\n", nameStr, reflect.ValueOf(anyType).Int()))

	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		str.WriteString(fmt.Sprintf("%s %d\n", nameStr, reflect.ValueOf(anyType).Uint()))

	case reflect.Float32, reflect.Float64:
		str.WriteString(fmt.Sprintf("%s %v\n", nameStr, int(reflect.ValueOf(anyType).Float())))

	case reflect.Bool:
		str.WriteString(fmt.Sprintf("%s %v\n", nameStr, reflect.ValueOf(anyType).Bool()))
	}

	return str.String()
}

func readReflectMapSprint(list map[string]interface{}, num int) string {
	var str string
	ws := GetWhitespace(num)

	for k, v := range list {
		r := reflect.TypeOf(v)

		if r == nil {
			return str
		}

		switch r.Kind() {
		case reflect.Map:
			if v, ok := v.(map[string]interface{}); ok {
				str += fmt.Sprintf("%s%s:\n", ws, k)
				str += readReflectMapSprint(v, num+1)
			}

		case reflect.Slice:
			if v, ok := v.([]interface{}); ok {
				str += fmt.Sprintf("%s%s:\n", ws, k)
				str += readReflectSliceSprint(v, num+1)
			}

		case reflect.Array:
			str += fmt.Sprintf("%s: %s (it is array)\n", k, reflect.ValueOf(v).String())

		default:
			str += readReflectAnyTypeSprint(k, v, num)
		}
	}

	return str
}

func readReflectSliceSprint(list []interface{}, num int) string {
	var str string
	ws := GetWhitespace(num)

	for k, v := range list {
		r := reflect.TypeOf(v)

		if r == nil {
			return str
		}

		switch r.Kind() {
		case reflect.Map:
			if v, ok := v.(map[string]interface{}); ok {
				str += fmt.Sprintf("%s%d.\n", ws, k+1)
				str += readReflectMapSprint(v, num+1)
			}

		case reflect.Slice:
			if v, ok := v.([]interface{}); ok {
				str += fmt.Sprintf("%s%d.\n", ws, k+1)
				str += readReflectSliceSprint(v, num+1)
			}

		case reflect.Array:
			str += fmt.Sprintf("%d. %s (it is array)\n", k, reflect.ValueOf(v).String())

		default:
			str += readReflectAnyTypeSprint(k, v, num)
		}
	}

	return str
}
