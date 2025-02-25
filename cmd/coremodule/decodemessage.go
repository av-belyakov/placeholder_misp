package coremodule

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/av-belyakov/placeholder_misp/commoninterfaces"
	"github.com/av-belyakov/placeholder_misp/internal/countermessage"
	"github.com/av-belyakov/placeholder_misp/internal/supportingfunctions"
)

type HandlerJsonMessageSettings struct {
	logger   commoninterfaces.Logger
	counting *countermessage.CounterMessage
}

// NewHandlerJSON конструктор нового обработчика JSON сообщения
func NewHandlerJSON(counting *countermessage.CounterMessage, logger commoninterfaces.Logger) *HandlerJsonMessageSettings {
	return &HandlerJsonMessageSettings{
		logger:   logger,
		counting: counting,
	}
}

// Start инициализация обработки
func (s *HandlerJsonMessageSettings) Start(b []byte, taskId string) chan ChanInputCreateMispFormat {
	chanInput := make(chan ChanInputCreateMispFormat)

	go func() {
		//для карт
		listMap := map[string]interface{}{}
		if err := json.Unmarshal(b, &listMap); err == nil {
			if len(listMap) == 0 {
				s.logger.Send("error", supportingfunctions.CustomError(errors.New("error decoding the json message, it may be empty")).Error())

				return
			}

			_ = processingReflectMap(chanInput, listMap, "")
		} else {
			// для срезов
			listSlice := []interface{}{}
			if err = json.Unmarshal(b, &listSlice); err != nil {
				s.logger.Send("error", supportingfunctions.CustomError(err).Error())

				return
			}

			if len(listSlice) == 0 {
				s.logger.Send("error", supportingfunctions.CustomError(errors.New("error decoding the json message, it may be empty")).Error())

				return
			}

			_ = processingReflectSlice(chanInput, listSlice, "")
		}

		// сетчик обработанных кейсов
		s.counting.SendMessage("update processed events", 1)

		// ***********************************
		// Это логирование только для теста!!!
		// ***********************************
		s.logger.Send("testing", "TEST_INFO func 'HandlerJsonMessage', handling json message")
		//
		//

		fmt.Println("HandlerJsonMessageSettings.Start, close channel")
		close(chanInput)
	}()

	return chanInput
}

func processingReflectAnySimpleType(
	chanOutMispFormat chan<- ChanInputCreateMispFormat,
	name interface{},
	anyType interface{},
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
		chanOutMispFormat <- ChanInputCreateMispFormat{
			FieldName:   nameStr,
			ValueType:   "string",
			Value:       result,
			FieldBranch: fieldBranch,
		}

		return result

	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		result := reflect.ValueOf(anyType).Int()
		chanOutMispFormat <- ChanInputCreateMispFormat{
			FieldName:   nameStr,
			ValueType:   "int",
			Value:       result,
			FieldBranch: fieldBranch,
		}

		return result

	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		result := reflect.ValueOf(anyType).Uint()
		chanOutMispFormat <- ChanInputCreateMispFormat{
			FieldName:   nameStr,
			ValueType:   "uint",
			Value:       result,
			FieldBranch: fieldBranch,
		}

		return result

	case reflect.Float32, reflect.Float64:
		result := reflect.ValueOf(anyType).Float()
		chanOutMispFormat <- ChanInputCreateMispFormat{
			FieldName:   nameStr,
			ValueType:   "float",
			Value:       result,
			FieldBranch: fieldBranch,
		}

		return result

	case reflect.Bool:
		result := reflect.ValueOf(anyType).Bool()
		chanOutMispFormat <- ChanInputCreateMispFormat{
			FieldName:   nameStr,
			ValueType:   "bool",
			Value:       result,
			FieldBranch: fieldBranch,
		}

		return result
	}

	return anyType
}

func processingReflectMap(
	chanOutMispFormat chan<- ChanInputCreateMispFormat,
	l map[string]interface{},
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
				newMap = processingReflectMap(chanOutMispFormat, v, fbTmp)
				nl[k] = newMap
			}

		case reflect.Slice:
			if v, ok := v.([]interface{}); ok {
				newList = processingReflectSlice(chanOutMispFormat, v, fbTmp)
				nl[k] = newList
			}

		default:
			nl[k] = processingReflectAnySimpleType(chanOutMispFormat, k, v, fbTmp)
		}
	}

	return nl
}

func processingReflectSlice(
	chanOutMispFormat chan<- ChanInputCreateMispFormat,
	l []interface{},
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
				newMap = processingReflectMap(chanOutMispFormat, v, fieldBranch)

				nl = append(nl, newMap)
			}

		case reflect.Slice:
			if v, ok := v.([]interface{}); ok {
				newList = processingReflectSlice(chanOutMispFormat, v, fieldBranch)

				nl = append(nl, newList...)
			}

		default:
			nl = append(nl, processingReflectAnySimpleType(chanOutMispFormat, k, v, fieldBranch))
		}
	}

	return nl
}
