package datamodels

import (
	"fmt"
	"sync"
	"time"
)

func NewListAttributesMispFormat() *ListAttributesMispFormat {
	return &ListAttributesMispFormat{
		attributes: make([]AttributesMispFormat, 0),
		mutex:      sync.Mutex{},
	}
}

func createNewAttributesMisp() AttributesMispFormat {
	return AttributesMispFormat{
		Category:     "Other",
		Type:         "other",
		Timestamp:    "0",
		Distribution: "3",
		//FirstSeen:      fmt.Sprint(time.Now().UnixMicro()),
		//LastSeen:       fmt.Sprint(time.Now().UnixMicro()),
		FirstSeen:      fmt.Sprint(time.Now().Format(time.RFC3339)),
		LastSeen:       fmt.Sprint(time.Now().Format(time.RFC3339)),
		ToIds:          true,
		SharingGroupId: "1",
	}
}

func (lambda *ListAttributesMispFormat) getCountListAttributesMisp() int {
	return len(lambda.attributes)
}

func (lamisp *ListAttributesMispFormat) GetListAttributesMisp() []AttributesMispFormat {
	return lamisp.attributes
}

func (lamisp *ListAttributesMispFormat) SetValueEventIdAttributesMisp(v interface{}, isNew bool) {
	lamisp.mutex.Lock()
	for k := range lamisp.attributes {
		(*lamisp).attributes[k].EventId = fmt.Sprint(v)
	}
	lamisp.mutex.Unlock()
}

func (lamisp *ListAttributesMispFormat) SetValueObjectIdAttributesMisp(v interface{}, isNew bool) {
	lamisp.mutex.Lock()
	if isNew {
		lamisp.attributes = append(lamisp.attributes, createNewAttributesMisp())
	}

	lamisp.attributes[lamisp.getCountListAttributesMisp()-1].ObjectId = fmt.Sprint(v)
	lamisp.mutex.Unlock()
}

func (lamisp *ListAttributesMispFormat) SetValueObjectRelationAttributesMisp(v interface{}, isNew bool) {
	lamisp.mutex.Lock()
	if isNew {
		lamisp.attributes = append(lamisp.attributes, createNewAttributesMisp())
	}

	lamisp.attributes[lamisp.getCountListAttributesMisp()-1].ObjectRelation = fmt.Sprint(v)
	lamisp.mutex.Unlock()
}

func (lamisp *ListAttributesMispFormat) SetValueCategoryAttributesMisp(v interface{}, isNew bool) {
	lamisp.mutex.Lock()
	if isNew {
		lamisp.attributes = append(lamisp.attributes, createNewAttributesMisp())
	}

	lamisp.attributes[lamisp.getCountListAttributesMisp()-1].Category = fmt.Sprint(v)
	lamisp.mutex.Unlock()
}

func (lamisp *ListAttributesMispFormat) SetValueTypeAttributesMisp(v interface{}, isNew bool) {
	lamisp.mutex.Lock()
	if isNew {
		lamisp.attributes = append(lamisp.attributes, createNewAttributesMisp())
	}

	lamisp.attributes[lamisp.getCountListAttributesMisp()-1].Type = fmt.Sprint(v)
	lamisp.mutex.Unlock()
}

func (lamisp *ListAttributesMispFormat) SetValueValueAttributesMisp(v interface{}, isNew bool) {
	lamisp.mutex.Lock()
	if isNew {
		lamisp.attributes = append(lamisp.attributes, createNewAttributesMisp())
	}

	lamisp.attributes[lamisp.getCountListAttributesMisp()-1].Value = fmt.Sprint(v)
	lamisp.mutex.Unlock()
}

func (lamisp *ListAttributesMispFormat) SetValueUuidAttributesMisp(v interface{}, isNew bool) {
	lamisp.mutex.Lock()
	if isNew {
		lamisp.attributes = append(lamisp.attributes, createNewAttributesMisp())
	}

	lamisp.attributes[lamisp.getCountListAttributesMisp()-1].Uuid = fmt.Sprint(v)
	lamisp.mutex.Unlock()
}

func (lamisp *ListAttributesMispFormat) SetValueTimestampAttributesMisp(v interface{}, isNew bool) {
	lamisp.mutex.Lock()
	if isNew {
		lamisp.attributes = append(lamisp.attributes, createNewAttributesMisp())
	}

	if dt, ok := v.(float64); ok {
		lamisp.attributes[lamisp.getCountListAttributesMisp()-1].Timestamp = fmt.Sprintf("%10.f", dt)[:10]
	}

	lamisp.mutex.Unlock()
}

func (lamisp *ListAttributesMispFormat) SetValueDistributionAttributesMisp(v interface{}, isNew bool) {
	lamisp.mutex.Lock()
	if isNew {
		lamisp.attributes = append(lamisp.attributes, createNewAttributesMisp())
	}

	lamisp.attributes[lamisp.getCountListAttributesMisp()-1].Distribution = fmt.Sprint(v)
	lamisp.mutex.Unlock()
}

func (lamisp *ListAttributesMispFormat) SetValueSharingGroupIdAttributesMisp(v interface{}, isNew bool) {
	lamisp.mutex.Lock()
	if isNew {
		lamisp.attributes = append(lamisp.attributes, createNewAttributesMisp())
	}

	lamisp.attributes[lamisp.getCountListAttributesMisp()-1].SharingGroupId = fmt.Sprint(v)
	lamisp.mutex.Unlock()
}

func (lamisp *ListAttributesMispFormat) SetValueCommentAttributesMisp(v interface{}, isNew bool) {
	lamisp.mutex.Lock()
	if isNew {
		lamisp.attributes = append(lamisp.attributes, createNewAttributesMisp())
	}

	lamisp.attributes[lamisp.getCountListAttributesMisp()-1].Comment = fmt.Sprint(v)
	lamisp.mutex.Unlock()
}

func (lamisp *ListAttributesMispFormat) SetValueFirstSeenAttributesMisp(v interface{}, isNew bool) {
	lamisp.mutex.Lock()
	if isNew {
		lamisp.attributes = append(lamisp.attributes, createNewAttributesMisp())
	}

	if dt, ok := v.(float64); ok {
		lamisp.attributes[lamisp.getCountListAttributesMisp()-1].FirstSeen = fmt.Sprintf("%13.f000", dt)
	}

	lamisp.mutex.Unlock()
}

func (lamisp *ListAttributesMispFormat) SetValueLastSeenAttributesMisp(v interface{}, isNew bool) {
	lamisp.mutex.Lock()
	if isNew {
		lamisp.attributes = append(lamisp.attributes, createNewAttributesMisp())
	}

	if dt, ok := v.(float64); ok {
		lamisp.attributes[lamisp.getCountListAttributesMisp()-1].LastSeen = fmt.Sprintf("%13.f000", dt)
	}

	lamisp.mutex.Unlock()
}

func (lamisp *ListAttributesMispFormat) SetValueToIdsAttributesMisp(v interface{}, isNew bool) {
	lamisp.mutex.Lock()
	if isNew {
		lamisp.attributes = append(lamisp.attributes, createNewAttributesMisp())
	}

	if data, ok := v.(bool); ok {
		lamisp.attributes[lamisp.getCountListAttributesMisp()-1].ToIds = data
	}
	lamisp.mutex.Unlock()
}

func (lamisp *ListAttributesMispFormat) SetValueDeletedAttributesMisp(v interface{}, isNew bool) {
	lamisp.mutex.Lock()
	if isNew {
		lamisp.attributes = append(lamisp.attributes, createNewAttributesMisp())
	}

	if data, ok := v.(bool); ok {
		lamisp.attributes[lamisp.getCountListAttributesMisp()-1].Deleted = data
	}
	lamisp.mutex.Unlock()
}

func (lamisp *ListAttributesMispFormat) SetValueDisableCorrelationAttributesMisp(v interface{}, isNew bool) {
	lamisp.mutex.Lock()
	if isNew {
		lamisp.attributes = append(lamisp.attributes, createNewAttributesMisp())
	}

	if data, ok := v.(bool); ok {
		lamisp.attributes[lamisp.getCountListAttributesMisp()-1].DisableCorrelation = data
	}
	lamisp.mutex.Unlock()
}
