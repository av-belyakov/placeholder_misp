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
	errSrc := "error decoding the json file, it may be empty"
	str := strings.Builder{}
	defer str.Reset()

	listMap := map[string]any{}
	if err := json.Unmarshal(b, &listMap); err == nil {
		if len(listMap) == 0 {
			return "", errors.New(errSrc)
		}

		readReflectMapSprint(listMap, &str, 0)

		return str.String(), err
	}

	listSlice := []any{}
	if err := json.Unmarshal(b, &listSlice); err == nil {
		if len(listSlice) == 0 {
			return "", errors.New(errSrc)
		}

		readReflectSliceSprint(listSlice, &str, 0)

		return str.String(), err
	}

	return str.String(), fmt.Errorf("the contents of the file are not Map or Slice")
}

func readReflectAnyTypeSprint(name interface{}, anyType interface{}, num int) string {
	var (
		str         string
		nameStr     string
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
		return str
	}

	switch r.Kind() {
	case reflect.String:
		dataStr := reflect.ValueOf(anyType).String()

		if isCleanLine {
			dataStr = strings.ReplaceAll(dataStr, "\t", "")
			dataStr = strings.ReplaceAll(dataStr, "\n", "")
		}

		str = fmt.Sprintf("%s \"%s\"\n", nameStr, dataStr)

	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		str = fmt.Sprintf("%s %d\n", nameStr, reflect.ValueOf(anyType).Int())

	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		str = fmt.Sprintf("%s %d\n", nameStr, reflect.ValueOf(anyType).Uint())

	case reflect.Float32, reflect.Float64:
		str = fmt.Sprintf("%s %v\n", nameStr, int(reflect.ValueOf(anyType).Float()))

	case reflect.Bool:
		str = fmt.Sprintf("%s %v\n", nameStr, reflect.ValueOf(anyType).Bool())
	}

	return str
}

func readReflectMapSprint(list map[string]interface{}, str *strings.Builder, num int) {
	ws := GetWhitespace(num)

	for k, v := range list {
		r := reflect.TypeOf(v)

		if r == nil {
			return
		}

		switch r.Kind() {
		case reflect.Map:
			if v, ok := v.(map[string]interface{}); ok {
				str.WriteString(fmt.Sprintf("%s%s:\n", ws, k))
				readReflectMapSprint(v, str, num+1)
			}

		case reflect.Slice:
			if v, ok := v.([]interface{}); ok {
				str.WriteString(fmt.Sprintf("%s%s:\n", ws, k))
				readReflectSliceSprint(v, str, num+1)
			}

		case reflect.Array:
			str.WriteString(fmt.Sprintf("%s: %s (it is array)\n", k, reflect.ValueOf(v).String()))

		default:
			str.WriteString(readReflectAnyTypeSprint(k, v, num))
		}
	}
}

func readReflectSliceSprint(list []interface{}, str *strings.Builder, num int) {
	ws := GetWhitespace(num)

	for k, v := range list {
		r := reflect.TypeOf(v)

		if r == nil {
			return
		}

		switch r.Kind() {
		case reflect.Map:
			if v, ok := v.(map[string]interface{}); ok {
				str.WriteString(fmt.Sprintf("%s%d.\n", ws, k+1))
				readReflectMapSprint(v, str, num+1)
			}

		case reflect.Slice:
			if v, ok := v.([]interface{}); ok {
				str.WriteString(fmt.Sprintf("%s%d.\n", ws, k+1))
				readReflectSliceSprint(v, str, num+1)
			}

		case reflect.Array:
			str.WriteString(fmt.Sprintf("%d. %s (it is array)\n", k, reflect.ValueOf(v).String()))

		default:
			str.WriteString(readReflectAnyTypeSprint(k, v, num))
		}
	}
}
