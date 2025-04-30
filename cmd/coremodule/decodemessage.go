package coremodule

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/av-belyakov/placeholder_misp/commoninterfaces"
	"github.com/av-belyakov/placeholder_misp/internal/supportingfunctions"
)

type HandlerJsonMessageSettings struct {
	logger  commoninterfaces.Logger
	counter commoninterfaces.Counter
}

// NewHandlerJSON конструктор нового обработчика JSON сообщения
func NewHandlerJSON(counter commoninterfaces.Counter, logger commoninterfaces.Logger) *HandlerJsonMessageSettings {
	return &HandlerJsonMessageSettings{
		logger:  logger,
		counter: counter,
	}
}

// Start инициализация обработки
func (s *HandlerJsonMessageSettings) Start(b []byte, taskId string) chan ChanInputCreateMispFormat {
	chanInput := make(chan ChanInputCreateMispFormat)

	go func() {
		//для карт
		listMap := map[string]any{}
		if err := json.Unmarshal(b, &listMap); err == nil {
			if len(listMap) == 0 {
				s.logger.Send("error", supportingfunctions.CustomError(errors.New("error decoding the json message, it may be empty")).Error())

				return
			}

			_ = processingReflectMap(chanInput, listMap, "")
		} else {
			// для срезов
			listSlice := []any{}
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
		s.counter.SendMessage("update processed events", 1)

		close(chanInput)
	}()

	return chanInput
}

func processingReflectAnySimpleType(
	chanOutMispFormat chan<- ChanInputCreateMispFormat,
	name any,
	anyType any,
	fieldBranch string) any {

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
	list map[string]any,
	fieldBranch string) map[string]any {
	var (
		newMap  map[string]any
		newList []any

		nl map[string]any = map[string]any{}
	)

	for k, v := range list {
		var fbTmp string
		r := reflect.TypeOf(v)

		if r == nil {
			continue
		}

		fbTmp = fieldBranch
		if fbTmp == "" {
			fbTmp += k
		} else {
			fbTmp += "." + k
		}

		switch r.Kind() {
		case reflect.Map:
			if v, ok := v.(map[string]any); ok {
				newMap = processingReflectMap(chanOutMispFormat, v, fbTmp)
				nl[k] = newMap
			}

		case reflect.Slice:
			if v, ok := v.([]any); ok {
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
	list []any,
	fieldBranch string) []any {
	var (
		newMap  map[string]any
		newList []any

		nl []any = make([]any, 0, len(list))
	)

	for k, v := range list {
		r := reflect.TypeOf(v)

		if r == nil {
			continue
		}

		switch r.Kind() {
		case reflect.Map:
			if v, ok := v.(map[string]any); ok {
				newMap = processingReflectMap(chanOutMispFormat, v, fieldBranch)

				nl = append(nl, newMap)
			}

		case reflect.Slice:
			if v, ok := v.([]any); ok {
				newList = processingReflectSlice(chanOutMispFormat, v, fieldBranch)

				nl = append(nl, newList...)
			}

		default:
			nl = append(nl, processingReflectAnySimpleType(chanOutMispFormat, k, v, fieldBranch))
		}
	}

	return nl
}
