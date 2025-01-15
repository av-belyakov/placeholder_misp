package testcreategalaxy_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"

	"github.com/av-belyakov/placeholder_misp/internal/datamodels"
	rules "github.com/av-belyakov/placeholder_misp/rulesinteraction"
)

// ChanInputCreateMispFormat
// передаваемых данных
// UUID - уникальный идентификатор в формате UUID
// FieldName - наименование поля
// ValueType - тип передаваемого значения (string, int и т.д.)
// Value - любые передаваемые данные
// FieldBranch - 'путь' до значения в как в JSON формате, например 'event.details.customFields.class'
type ChanInputCreateMispFormat struct {
	UUID        string
	FieldName   string
	ValueType   string
	Value       interface{}
	FieldBranch string
}

func DecodeJsonObject(
	cmispf chan<- ChanInputCreateMispFormat,
	b []byte,
	listRule *rules.ListRule,
	logging chan<- datamodels.MessageLogging,
	cmispfDone chan<- bool,
) {
	listMap := map[string]interface{}{}
	if err := json.Unmarshal(b, &listMap); err == nil {
		_ = processingReflectMap(logging, cmispf, listMap, listRule, "")
	}

	cmispfDone <- true
}

func processingReflectAnySimpleType(
	logging chan<- datamodels.MessageLogging,
	chanOutMispFormat chan<- ChanInputCreateMispFormat,
	name interface{},
	anyType interface{},
	lr *rules.ListRule,
	fieldBranch string) interface{} {

	var nameStr string
	r := reflect.TypeOf(anyType)

	if n, ok := name.(int); ok {
		nameStr = fmt.Sprint(n)
	} else if n, ok := name.(string); ok {
		nameStr = n
	}

	if r == nil {
		return anyType
	}

	switch r.Kind() {
	case reflect.String:
		result := reflect.ValueOf(anyType).String()

		ncv, num, err := lr.ReplacementRuleHandler("string", nameStr, result)
		if err != nil {
			_, f, l, _ := runtime.Caller(0)

			logging <- datamodels.MessageLogging{
				MsgData: fmt.Sprintf("'search value \"%s\" from rule number \"%d\" of section \"REPLACE\" is not fulfilled' %s:%d", result, num, f, l-1),
				MsgType: "warning",
			}
		}

		lr.PassRuleHandler(fieldBranch, ncv)

		chanOutMispFormat <- ChanInputCreateMispFormat{
			FieldName:   nameStr,
			ValueType:   "string",
			Value:       ncv,
			FieldBranch: fieldBranch,
		}

		return ncv
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		result := reflect.ValueOf(anyType).Int()

		ncv, num, err := lr.ReplacementRuleHandler("int", nameStr, result)
		if err != nil {
			_, f, l, _ := runtime.Caller(0)

			logging <- datamodels.MessageLogging{
				MsgData: fmt.Sprintf("'search value \"%d\" from rule number \"%d\" of section \"REPLACE\" is not fulfilled' %s:%d", result, num, f, l-1),
				MsgType: "warning",
			}
		}

		lr.PassRuleHandler(fieldBranch, ncv)

		chanOutMispFormat <- ChanInputCreateMispFormat{
			FieldName:   nameStr,
			ValueType:   "int",
			Value:       ncv,
			FieldBranch: fieldBranch,
		}

		return ncv
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		result := reflect.ValueOf(anyType).Uint()

		ncv, num, err := lr.ReplacementRuleHandler("uint", nameStr, result)
		if err != nil {
			_, f, l, _ := runtime.Caller(0)

			logging <- datamodels.MessageLogging{
				MsgData: fmt.Sprintf("'search value \"%d\" from rule number \"%d\" of section \"REPLACE\" is not fulfilled' %s:%d", result, num, f, l-1),
				MsgType: "warning",
			}
		}

		lr.PassRuleHandler(fieldBranch, ncv)

		chanOutMispFormat <- ChanInputCreateMispFormat{
			FieldName:   nameStr,
			ValueType:   "uint",
			Value:       ncv,
			FieldBranch: fieldBranch,
		}

		return ncv
	case reflect.Float32, reflect.Float64:
		result := reflect.ValueOf(anyType).Float()

		ncv, num, err := lr.ReplacementRuleHandler("float", nameStr, result)
		if err != nil {
			_, f, l, _ := runtime.Caller(0)

			logging <- datamodels.MessageLogging{
				MsgData: fmt.Sprintf("'search value \"%v\" from rule number \"%d\" of section \"REPLACE\" is not fulfilled' %s:%d", result, num, f, l-1),
				MsgType: "warning",
			}
		}

		lr.PassRuleHandler(fieldBranch, ncv)

		chanOutMispFormat <- ChanInputCreateMispFormat{
			FieldName:   nameStr,
			ValueType:   "float",
			Value:       ncv,
			FieldBranch: fieldBranch,
		}

		return ncv
	case reflect.Bool:
		result := reflect.ValueOf(anyType).Bool()

		ncv, num, err := lr.ReplacementRuleHandler("bool", nameStr, result)
		if err != nil {
			_, f, l, _ := runtime.Caller(0)

			logging <- datamodels.MessageLogging{
				MsgData: fmt.Sprintf("'search value \"%v\" from rule number \"%d\" of section \"REPLACE\" is not fulfilled' %s:%d", result, num, f, l-1),
				MsgType: "warning",
			}
		}

		lr.PassRuleHandler(fieldBranch, ncv)

		chanOutMispFormat <- ChanInputCreateMispFormat{
			FieldName:   nameStr,
			ValueType:   "bool",
			Value:       ncv,
			FieldBranch: fieldBranch,
		}

		return ncv
	}

	return anyType
}

func processingReflectMap(
	logging chan<- datamodels.MessageLogging,
	chanOutMispFormat chan<- ChanInputCreateMispFormat,
	l map[string]interface{},
	lr *rules.ListRule,
	fieldBranch string) map[string]interface{} {

	var (
		newMap  map[string]interface{}
		newList []interface{}
	)
	nl := map[string]interface{}{}

	for k, v := range l {
		var fbTmp string
		r := reflect.TypeOf(v)

		if r == nil {
			return nl
		}

		fbTmp = fieldBranch
		if fbTmp == "" {
			fbTmp += k
		} else {
			fbTmp += "." + k
		}

		switch r.Kind() {
		case reflect.Map:
			if v, ok := v.(map[string]interface{}); ok {
				newMap = processingReflectMap(logging, chanOutMispFormat, v, lr, fbTmp)
				nl[k] = newMap
			}

		case reflect.Slice:
			if v, ok := v.([]interface{}); ok {
				newList = processingReflectSlice(logging, chanOutMispFormat, v, lr, fbTmp)
				nl[k] = newList
			}

		default:
			nl[k] = processingReflectAnySimpleType(logging, chanOutMispFormat, k, v, lr, fbTmp)
		}
	}

	return nl
}

func processingReflectSlice(
	logging chan<- datamodels.MessageLogging,
	chanOutMispFormat chan<- ChanInputCreateMispFormat,
	l []interface{},
	lr *rules.ListRule,
	fieldBranch string) []interface{} {

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

		switch r.Kind() {
		case reflect.Map:
			if v, ok := v.(map[string]interface{}); ok {
				newMap = processingReflectMap(logging, chanOutMispFormat, v, lr, fieldBranch)

				nl = append(nl, newMap)
			}

		case reflect.Slice:
			if v, ok := v.([]interface{}); ok {
				newList = processingReflectSlice(logging, chanOutMispFormat, v, lr, fieldBranch)

				nl = append(nl, newList...)
			}

		default:
			nl = append(nl, processingReflectAnySimpleType(logging, chanOutMispFormat, k, v, lr, fieldBranch))
		}
	}

	return nl
}
