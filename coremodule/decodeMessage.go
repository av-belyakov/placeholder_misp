package coremodule

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"

	"placeholder_misp/datamodels"
	"placeholder_misp/memorytemporarystorage"
	rules "placeholder_misp/rulesinteraction"
	"placeholder_misp/supportingfunctions"
)

type HandlerMessageFromHiveSettings struct {
	StorageApp *memorytemporarystorage.CommonStorageTemporary
	ListRule   *rules.ListRule
	Logging    chan<- datamodels.MessageLogging
	Counting   chan<- datamodels.DataCounterSettings
}

func NewHandlerMessageFromHive(
	storageApp *memorytemporarystorage.CommonStorageTemporary,
	listRule *rules.ListRule,
	logging chan<- datamodels.MessageLogging,
	counting chan<- datamodels.DataCounterSettings) *HandlerMessageFromHiveSettings {
	return &HandlerMessageFromHiveSettings{
		StorageApp: storageApp,
		ListRule:   listRule,
		Logging:    logging,
		Counting:   counting,
	}
}

func (s *HandlerMessageFromHiveSettings) HandlerMessageFromHive(
	cmispf chan<- ChanInputCreateMispFormat,
	b []byte,
	taskId string,
	cmispfDone chan<- bool,
) {
	var (
		f   string
		l   int
		err error
	)

	//для записи событий в лог-файл events
	str, _ := supportingfunctions.NewReadReflectJSONSprint(b)
	s.Logging <- datamodels.MessageLogging{
		MsgData: fmt.Sprintf("\t---------------\n\tEVENTS:\n%s\n", str),
		MsgType: "events",
	}

	//для карт
	_, f, l, _ = runtime.Caller(0)
	listMap := map[string]interface{}{}
	if err = json.Unmarshal(b, &listMap); err == nil {
		if len(listMap) == 0 {
			s.Logging <- datamodels.MessageLogging{
				MsgData: fmt.Sprintf("'error decoding the json message, it may be empty' %s:%d", f, l+2),
				MsgType: "error",
			}

			return
		}

		nlt := processingReflectMap(s.Logging, cmispf, listMap, s.ListRule, 0, "")
		s.StorageApp.SetProcessedDataHiveFormatMessage(taskId, nlt)
	} else {
		// для срезов
		_, f, l, _ = runtime.Caller(0)
		listSlice := []interface{}{}
		if err = json.Unmarshal(b, &listSlice); err != nil {
			s.Logging <- datamodels.MessageLogging{
				MsgData: fmt.Sprintf("'%s' %s:%d", err.Error(), f, l+2),
				MsgType: "error",
			}

			return
		}

		if len(listSlice) == 0 {
			s.Logging <- datamodels.MessageLogging{
				MsgData: fmt.Sprintf("'error decoding the json message, it may be empty' %s:%d", f, l+2),
				MsgType: "error",
			}

			return
		}

		_ = processingReflectSlice(s.Logging, cmispf, listSlice, s.ListRule, 0, "")
	}

	// сетчик обработанных кейсов
	s.Counting <- datamodels.DataCounterSettings{
		DataType: "update processed events",
		Count:    1,
	}

	//проверяем что бы хотя бы одно правило разрешало пропуск кейса
	if s.ListRule.GetRulePassany() || s.ListRule.SomePassRuleIsTrue() {
		s.StorageApp.SetAllowedTransferTrueHiveFormatMessage(taskId)
	}

	//выполняет очистку значения StatementExpression что равно отсутствию совпадений в правилах Pass
	// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	// ВРЕМЕННО ЗАКОМЕНТИРОВАЛ
	// не забыть снять коментарий
	//
	s.ListRule.CleanStatementExpressionRulePass()

	isAllowed, _ := s.StorageApp.GetAllowedTransferHiveFormatMessage(taskId)
	dt := "events do not meet rules"
	if isAllowed {
		dt = "update events meet rules"
	}

	//сетчик кейсов соответствующих или не соответствующих правилам
	s.Counting <- datamodels.DataCounterSettings{
		DataType: dt,
		Count:    1,
	}

	//устанавливаем параметр информирующий о завершении обработки модулем
	s.StorageApp.SetIsProcessedMispHiveFormatMessage(taskId)

	//останавливаем обработчик формирующий MISP формат
	cmispfDone <- isAllowed
}

func processingReflectAnySimpleType(
	logging chan<- datamodels.MessageLogging,
	chanOutMispFormat chan<- ChanInputCreateMispFormat,
	name interface{},
	anyType interface{},
	lr *rules.ListRule,
	num int,
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
			FieldName:           nameStr,
			ValueType:           "string",
			Value:               ncv,
			FieldBranch:         fieldBranch,
			ExclusionRuleWorked: lr.ExcludeRuleHandler(fieldBranch, ncv),
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
			FieldName:           nameStr,
			ValueType:           "int",
			Value:               ncv,
			FieldBranch:         fieldBranch,
			ExclusionRuleWorked: lr.ExcludeRuleHandler(fieldBranch, ncv),
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
			FieldName:           nameStr,
			ValueType:           "uint",
			Value:               ncv,
			FieldBranch:         fieldBranch,
			ExclusionRuleWorked: lr.ExcludeRuleHandler(fieldBranch, ncv),
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
			FieldName:           nameStr,
			ValueType:           "float",
			Value:               ncv,
			FieldBranch:         fieldBranch,
			ExclusionRuleWorked: lr.ExcludeRuleHandler(fieldBranch, ncv),
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
			FieldName:           nameStr,
			ValueType:           "bool",
			Value:               ncv,
			FieldBranch:         fieldBranch,
			ExclusionRuleWorked: lr.ExcludeRuleHandler(fieldBranch, ncv),
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
	num int,
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
				newMap = processingReflectMap(logging, chanOutMispFormat, v, lr, num+1, fbTmp)
				nl[k] = newMap
			}

		case reflect.Slice:
			if v, ok := v.([]interface{}); ok {
				newList = processingReflectSlice(logging, chanOutMispFormat, v, lr, num+1, fbTmp)
				nl[k] = newList
			}

		default:
			nl[k] = processingReflectAnySimpleType(logging, chanOutMispFormat, k, v, lr, num, fbTmp)
		}
	}

	return nl
}

func processingReflectSlice(
	logging chan<- datamodels.MessageLogging,
	chanOutMispFormat chan<- ChanInputCreateMispFormat,
	l []interface{},
	lr *rules.ListRule,
	num int,
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
				newMap = processingReflectMap(logging, chanOutMispFormat, v, lr, num+1, fieldBranch)

				nl = append(nl, newMap)
			}

		case reflect.Slice:
			if v, ok := v.([]interface{}); ok {
				newList = processingReflectSlice(logging, chanOutMispFormat, v, lr, num+1, fieldBranch)

				nl = append(nl, newList...)
			}

		default:
			nl = append(nl, processingReflectAnySimpleType(logging, chanOutMispFormat, k, v, lr, num, fieldBranch))
		}
	}

	return nl
}
