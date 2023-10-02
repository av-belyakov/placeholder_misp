package coremodule

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"strconv"

	"placeholder_misp/datamodels"
	"placeholder_misp/memorytemporarystorage"
	rules "placeholder_misp/rulesinteraction"
	"placeholder_misp/supportingfunctions"
)

func HandlerMessageFromHive(
	b []byte,
	uuidTask string,
	storageApp *memorytemporarystorage.CommonStorageTemporary,
	listRule rules.ListRulesProcessingMsgMISP,
	cmispf chan<- ChanInputCreateMispFormat,
	cmispfDone chan<- bool,
	loging chan<- datamodels.MessageLoging) {
	listTmp := map[string]interface{}{}

	//для записи событий в файл events
	str, _ := supportingfunctions.NewReadReflectJSONSprint(b)
	loging <- datamodels.MessageLoging{
		MsgData: fmt.Sprintf("\t---------------\n\tEVENTS:\n%s\n", str),
		MsgType: "events",
	}

	if err := json.Unmarshal(b, &listTmp); err != nil {
		_, f, l, _ := runtime.Caller(0)

		loging <- datamodels.MessageLoging{
			MsgData: fmt.Sprintf("'%s' %s:%d", err.Error(), f, l-1),
			MsgType: "error",
		}

		return
	}

	if len(listTmp) == 0 {
		_, f, l, _ := runtime.Caller(0)

		loging <- datamodels.MessageLoging{
			MsgData: fmt.Sprintf("'%s' %s:%d", fmt.Errorf("error decoding the json message, it may be empty"), f, l-1),
			MsgType: "error",
		}

		return
	}

	nlt := processingReflectMap(loging, cmispf, listTmp, &listRule, 0, "")
	storageApp.SetProcessedDataHiveFormatMessage(uuidTask, nlt)

	if listRule.Rules.Passany {
		storageApp.SetAllowedTransferTrueHiveFormatMessage(uuidTask)
	} else {
		//проверяем соответствие сообщения правилам из раздела Pass
		for _, v := range listRule.Rules.Pass {
			skipMsg := true
			for _, value := range v.ListAnd {
				if !value.StatementExpression {
					skipMsg = false

					break
				}
			}

			if skipMsg {
				storageApp.SetAllowedTransferTrueHiveFormatMessage(uuidTask)

				break
			}
		}
	}

	isAllowed, _ := storageApp.GetAllowedTransferHiveFormatMessage(uuidTask)

	//устанавливаем параметр информирующий о завершении обработки модулем
	storageApp.SetIsProcessedMispHiveFormatMessage(uuidTask)

	//останавливаем обработчик формирующий MISP формат
	cmispfDone <- isAllowed
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

			rulePass[key].ListAnd[k].StatementExpression = true
		}
	}
}

func ReplacementRuleHandler(lr *rules.ListRulesProcessingMsgMISP, svt, fn string, cv interface{}) (interface{}, int, error) {
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
	loging chan<- datamodels.MessageLoging,
	chanOutMispFormat chan<- ChanInputCreateMispFormat,
	name interface{},
	anyType interface{},
	listRule *rules.ListRulesProcessingMsgMISP,
	num int,
	fieldBranch string) interface{} {

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
			_, f, l, _ := runtime.Caller(0)

			loging <- datamodels.MessageLoging{
				MsgData: fmt.Sprintf("'search value \"%s\" from rule number \"%d\" of section \"REPLACE\" is not fulfilled' %s:%d", result, num, f, l-1),
				MsgType: "warning",
			}
		}

		PassRuleHandler(listRule.Rules.Pass, nameStr, ncv)

		chanOutMispFormat <- ChanInputCreateMispFormat{
			FieldName:   nameStr,
			ValueType:   "string",
			Value:       ncv,
			FieldBranch: fieldBranch,
		}

		return ncv

	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		result := reflect.ValueOf(anyType).Int()

		ncv, num, err := ReplacementRuleHandler(listRule, "int", nameStr, result)
		if err != nil {
			_, f, l, _ := runtime.Caller(0)

			loging <- datamodels.MessageLoging{
				MsgData: fmt.Sprintf("'search value \"%d\" from rule number \"%d\" of section \"REPLACE\" is not fulfilled' %s:%d", result, num, f, l-1),
				MsgType: "warning",
			}
		}

		PassRuleHandler(listRule.Rules.Pass, nameStr, ncv)

		chanOutMispFormat <- ChanInputCreateMispFormat{
			FieldName:   nameStr,
			ValueType:   "int",
			Value:       ncv,
			FieldBranch: fieldBranch,
		}

		return ncv

	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		result := reflect.ValueOf(anyType).Uint()

		ncv, num, err := ReplacementRuleHandler(listRule, "uint", nameStr, result)
		if err != nil {
			_, f, l, _ := runtime.Caller(0)

			loging <- datamodels.MessageLoging{
				MsgData: fmt.Sprintf("'search value \"%d\" from rule number \"%d\" of section \"REPLACE\" is not fulfilled' %s:%d", result, num, f, l-1),
				MsgType: "warning",
			}
		}

		PassRuleHandler(listRule.Rules.Pass, nameStr, ncv)

		chanOutMispFormat <- ChanInputCreateMispFormat{
			FieldName:   nameStr,
			ValueType:   "uint",
			Value:       ncv,
			FieldBranch: fieldBranch,
		}

		return ncv

	case reflect.Float32, reflect.Float64:
		result := reflect.ValueOf(anyType).Float()

		ncv, num, err := ReplacementRuleHandler(listRule, "float", nameStr, result)
		if err != nil {
			_, f, l, _ := runtime.Caller(0)

			loging <- datamodels.MessageLoging{
				MsgData: fmt.Sprintf("'search value \"%v\" from rule number \"%d\" of section \"REPLACE\" is not fulfilled' %s:%d", result, num, f, l-1),
				MsgType: "warning",
			}
		}

		PassRuleHandler(listRule.Rules.Pass, nameStr, ncv)

		chanOutMispFormat <- ChanInputCreateMispFormat{
			FieldName:   nameStr,
			ValueType:   "float",
			Value:       ncv,
			FieldBranch: fieldBranch,
		}

		return ncv

	case reflect.Bool:
		result := reflect.ValueOf(anyType).Bool()

		ncv, num, err := ReplacementRuleHandler(listRule, "bool", nameStr, result)
		if err != nil {
			_, f, l, _ := runtime.Caller(0)

			loging <- datamodels.MessageLoging{
				MsgData: fmt.Sprintf("'search value \"%v\" from rule number \"%d\" of section \"REPLACE\" is not fulfilled' %s:%d", result, num, f, l-1),
				MsgType: "warning",
			}
		}

		PassRuleHandler(listRule.Rules.Pass, nameStr, ncv)

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
	loging chan<- datamodels.MessageLoging,
	chanOutMispFormat chan<- ChanInputCreateMispFormat,
	l map[string]interface{},
	lr *rules.ListRulesProcessingMsgMISP,
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
				newMap = processingReflectMap(loging, chanOutMispFormat, v, lr, num+1, fbTmp)
				nl[k] = newMap
			}

		case reflect.Slice:
			if v, ok := v.([]interface{}); ok {
				newList = processingReflectSlice(loging, chanOutMispFormat, v, lr, num+1, fbTmp)
				nl[k] = newList
			}

		default:
			nl[k] = processingReflectAnySimpleType(loging, chanOutMispFormat, k, v, lr, num, fbTmp)
		}
	}

	return nl
}

func processingReflectSlice(
	loging chan<- datamodels.MessageLoging,
	chanOutMispFormat chan<- ChanInputCreateMispFormat,
	l []interface{},
	lr *rules.ListRulesProcessingMsgMISP,
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
				newMap = processingReflectMap(loging, chanOutMispFormat, v, lr, num+1, fieldBranch)

				nl = append(nl, newMap)
			}

		case reflect.Slice:
			if v, ok := v.([]interface{}); ok {
				newList = processingReflectSlice(loging, chanOutMispFormat, v, lr, num+1, fieldBranch)

				nl = append(nl, newList...)
			}

		default:
			nl = append(nl, processingReflectAnySimpleType(loging, chanOutMispFormat, k, v, lr, num, fieldBranch))
		}
	}

	return nl
}
