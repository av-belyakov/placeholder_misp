package datamodels

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

func NewListObjectsMispFormat() *ListObjectsMispFormat {
	return &ListObjectsMispFormat{
		objects: map[int]ObjectsMispFormat{},
		mutex:   sync.Mutex{},
	}
}

/*
		"template_uuid": "c8cc27a6-4bd31-1f72-afa5-7b9bb4ac3b3b",
	      "template_version": "1",
	      "first_seen": "1581984000000000",
	      "timestamp": "1617875568",
	      "name": "file",
	      "description": "size 817 byte",
	      "event_id": "12660",
	      "meta-category": "file",
	      "distribution": "5",
		return AttributesMispFormat{
			Category:       "Other",
			Type:           "other",
			Timestamp:      "0",
			Distribution:   "3",
			FirstSeen:      fmt.Sprint(time.Now().Format(time.RFC3339)),
			LastSeen:       fmt.Sprint(time.Now().Format(time.RFC3339)),
			ToIds:          true,
			SharingGroupId: "1",
		}*/

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
	lomf.mutex.Lock()
	defer lomf.mutex.Unlock()

	lomf.objects = map[int]ObjectsMispFormat{}
}

func (lomf *ListObjectsMispFormat) GetListObjectsMisp() map[int]ObjectsMispFormat {
	return lomf.objects
}

func (lomf *ListObjectsMispFormat) SetValueEventIdObjectsMisp(v interface{}, num int) {
	var tmp ObjectsMispFormat
	lomf.mutex.Lock()
	defer lomf.mutex.Unlock()

	if obj, ok := lomf.objects[num]; ok {
		tmp = obj
	} else {
		tmp = createNewObjectsMisp()
	}

	tmp.EventID = fmt.Sprint(v)
	lomf.objects[num] = tmp
}

func (lomf *ListObjectsMispFormat) SetValueNameObjectsMisp(v interface{}, num int) {
	var tmp ObjectsMispFormat
	lomf.mutex.Lock()
	defer lomf.mutex.Unlock()

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
	lomf.mutex.Lock()
	defer lomf.mutex.Unlock()

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
	lomf.mutex.Lock()
	defer lomf.mutex.Unlock()

	if obj, ok := lomf.objects[num]; ok {
		tmp = obj
	} else {
		tmp = createNewObjectsMisp()
	}

	if dt, ok := v.(float64); ok {
		fst := fmt.Sprint(dt)
		fslen := len(fst)
		if fslen < 16 {
			fst = fst + strings.Repeat("0", 16-fslen)
		} else if fslen > 16 {
			fst = fst[:16]
		}

		tmp.FirstSeen = fst
		//tmp.FirstSeen = time.UnixMilli(int64(dt)).Format(time.RFC3339)
	}

	lomf.objects[num] = tmp
}

func (lomf *ListObjectsMispFormat) SetValueTimestampObjectsMisp(v interface{}, num int) {
	var tmp ObjectsMispFormat
	lomf.mutex.Lock()
	defer lomf.mutex.Unlock()

	if obj, ok := lomf.objects[num]; ok {
		tmp = obj
	} else {
		tmp = createNewObjectsMisp()
	}

	if dt, ok := v.(float64); ok {
		ts := fmt.Sprint(dt)
		fslen := len(ts)
		if fslen < 10 {
			ts = ts + strings.Repeat("0", 10-fslen)
		} else if fslen > 10 {
			ts = ts[:10]
		}

		tmp.Timestamp = ts
		//tmp.FirstSeen = time.UnixMilli(int64(dt)).Format(time.RFC3339)
	}

	lomf.objects[num] = tmp
}
