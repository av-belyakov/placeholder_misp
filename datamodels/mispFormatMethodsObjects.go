package datamodels

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func NewListObjectsMispFormat() *ListObjectsMispFormat {
	return &ListObjectsMispFormat{objects: map[int]ObjectsMispFormat{}}
}

func createNewObjectsMisp() ObjectsMispFormat {
	return ObjectsMispFormat{
		TemplateUUID:    uuid.NewString(),
		TemplateVersion: "1",
		FirstSeen:       fmt.Sprint(time.Now().UnixMicro()),
		Timestamp:       fmt.Sprint(time.Now().Unix()),
		MetaCategory:    "file",
		Distribution:    "5",
		Attribute:       []AttributeMispFormat{},
	}
}

func (lomf *ListObjectsMispFormat) GetCountListObjectsMisp() int {
	return len(lomf.objects)
}

func (lomf *ListObjectsMispFormat) CleanListObjectsMisp() {
	lomf.Lock()
	defer lomf.Unlock()

	lomf.objects = map[int]ObjectsMispFormat{}
}

func (lomf *ListObjectsMispFormat) GetListObjectsMisp() map[int]ObjectsMispFormat {
	return lomf.objects
}

func (lomf *ListObjectsMispFormat) SetValueEventIdObjectsMisp(v interface{}, num int) {
	var tmp ObjectsMispFormat
	lomf.Lock()
	defer lomf.Unlock()

	if obj, ok := lomf.objects[num]; ok {
		tmp = obj
	} else {
		tmp = createNewObjectsMisp()
	}

	tmp.EventId = fmt.Sprint(v)
	lomf.objects[num] = tmp
}

func (lomf *ListObjectsMispFormat) SetValueNameObjectsMisp(v interface{}, num int) {
	var tmp ObjectsMispFormat
	lomf.Lock()
	defer lomf.Unlock()

	if obj, ok := lomf.objects[num]; ok {
		tmp = obj
	} else {
		tmp = createNewObjectsMisp()
	}

	tmp.Name = fmt.Sprint(v)
	lomf.objects[num] = tmp
}

func (lomf *ListObjectsMispFormat) SetValueDescriptionObjectsMisp(v interface{}, num int) {
	var tmp ObjectsMispFormat
	lomf.Lock()
	defer lomf.Unlock()

	if obj, ok := lomf.objects[num]; ok {
		tmp = obj
	} else {
		tmp = createNewObjectsMisp()
	}

	tmp.Description = fmt.Sprint(v)
	lomf.objects[num] = tmp
}

func (lomf *ListObjectsMispFormat) SetValueFirstSeenObjectsMisp(v interface{}, num int) {
	var tmp ObjectsMispFormat
	lomf.Lock()
	defer lomf.Unlock()

	if obj, ok := lomf.objects[num]; ok {
		tmp = obj
	} else {
		tmp = createNewObjectsMisp()
	}

	if dt, ok := v.(float64); ok {
		tmp.FirstSeen = time.UnixMilli(int64(dt)).Format(time.RFC3339)
	}

	lomf.objects[num] = tmp
}

func (lomf *ListObjectsMispFormat) SetValueTimestampObjectsMisp(v interface{}, num int) {
	var tmp ObjectsMispFormat
	lomf.Lock()
	defer lomf.Unlock()

	if obj, ok := lomf.objects[num]; ok {
		tmp = obj
	} else {
		tmp = createNewObjectsMisp()
	}

	if dt, ok := v.(float64); ok {
		tmp.Timestamp = fmt.Sprintf("%10.f", dt)[:10]
	}

	lomf.objects[num] = tmp
}

func (lomf *ListObjectsMispFormat) SetValueSizeObjectsMisp(v interface{}, num int) {
	var tmp ObjectsMispFormat
	lomf.Lock()
	defer lomf.Unlock()

	if obj, ok := lomf.objects[num]; ok {
		tmp = obj
	} else {
		tmp = createNewObjectsMisp()
	}

	tmp.Description = fmt.Sprintf("размер %v байт", v)
	lomf.objects[num] = tmp
}

func (lomf *ListObjectsMispFormat) SetValueAttributeObjectsMisp(v interface{}, num int) {
	var tmp ObjectsMispFormat
	lomf.Lock()
	defer lomf.Unlock()

	if obj, ok := lomf.objects[num]; ok {
		tmp = obj
	} else {
		tmp = createNewObjectsMisp()
	}

	if newSlice, ok := v.([]AttributeMispFormat); ok {
		tmp.Attribute = newSlice
		lomf.objects[num] = tmp
	}
}
