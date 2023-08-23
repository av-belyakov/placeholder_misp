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
		Category:       "Other",
		Type:           "other",
		Timestamp:      "0",
		Distribution:   "3",
		FirstSeen:      fmt.Sprint(time.Now().UnixMicro()),
		LastSeen:       fmt.Sprint(time.Now().UnixMicro()),
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

/*
func (amisp *AttributesMispFormat) SetValueEventIdAttributesMisp(v interface{}) bool {
	var isSuccess bool

	if str, ok := v.(string); ok {
		amisp.EventId = str

		isSuccess = true
	}

	return isSuccess
}
*/

func (lamisp *ListAttributesMispFormat) SetValueEventIdAttributesMisp(v interface{}, isNew bool) {
	lamisp.mutex.Lock()
	for k := range lamisp.attributes {
		(*lamisp).attributes[k].EventId = fmt.Sprint(v)
	}
	lamisp.mutex.Unlock()
}

/*func (amisp *AttributesMispFormat) SetValueObjectIdAttributesMisp(v interface{}) bool {
	var isSuccess bool

	if str, ok := v.(string); ok {
		amisp.ObjectId = str

		isSuccess = true
	}

	return isSuccess
}*/

func (lamisp *ListAttributesMispFormat) SetValueObjectIdAttributesMisp(v interface{}, isNew bool) {
	fmt.Println("func 'SetValueObjectIdAttributesMisp', START")
	fmt.Println("111 isNew = ", isNew, " value = ", v)

	lamisp.mutex.Lock()
	if isNew {
		lamisp.attributes = append(lamisp.attributes, createNewAttributesMisp())
	}

	lamisp.attributes[lamisp.getCountListAttributesMisp()-1].ObjectId = fmt.Sprint(v)
	lamisp.mutex.Unlock()
}

/*func (amisp *AttributesMispFormat) SetValueObjectRelationAttributesMisp(v interface{}) bool {
	var isSuccess bool

	if str, ok := v.(string); ok {
		amisp.ObjectRelation = str

		isSuccess = true
	}

	return isSuccess
}*/

func (lamisp *ListAttributesMispFormat) SetValueObjectRelationAttributesMisp(v interface{}, isNew bool) {
	fmt.Println("func 'SetValueObjectRelationAttributesMisp', START")
	fmt.Println("222 isNew = ", isNew, " value = ", v)

	lamisp.mutex.Lock()
	if isNew {
		lamisp.attributes = append(lamisp.attributes, createNewAttributesMisp())
	}

	lamisp.attributes[lamisp.getCountListAttributesMisp()-1].ObjectRelation = fmt.Sprint(v)
	lamisp.mutex.Unlock()
}

/*func (amisp *AttributesMispFormat) SetValueCategoryAttributesMisp(v interface{}) bool {
	var isSuccess bool

	if str, ok := v.(string); ok {
		amisp.Category = str

		isSuccess = true
	}

	return isSuccess
}*/

func (lamisp *ListAttributesMispFormat) SetValueCategoryAttributesMisp(v interface{}, isNew bool) {
	fmt.Println("func 'SetValueCategoryAttributesMisp', START")
	fmt.Println("333 isNew = ", isNew, " value = ", v)

	lamisp.mutex.Lock()
	if isNew {
		lamisp.attributes = append(lamisp.attributes, createNewAttributesMisp())
	}

	lamisp.attributes[lamisp.getCountListAttributesMisp()-1].Category = fmt.Sprint(v)
	lamisp.mutex.Unlock()
}

/*func (amisp *AttributesMispFormat) SetValueTypeAttributesMisp(v interface{}) bool {
	var isSuccess bool

	if str, ok := v.(string); ok {
		amisp.Type = str

		isSuccess = true
	}

	return isSuccess
}*/

func (lamisp *ListAttributesMispFormat) SetValueTypeAttributesMisp(v interface{}, isNew bool) {
	fmt.Println("func 'SetValueTypeAttributesMisp', START")
	fmt.Println("444 isNew = ", isNew, " value = ", v)

	lamisp.mutex.Lock()
	if isNew {
		lamisp.attributes = append(lamisp.attributes, createNewAttributesMisp())
	}

	lamisp.attributes[lamisp.getCountListAttributesMisp()-1].Type = fmt.Sprint(v)
	lamisp.mutex.Unlock()
}

/*func (amisp *AttributesMispFormat) SetValueValueAttributesMisp(v interface{}) bool {
	var isSuccess bool

	if str, ok := v.(string); ok {
		amisp.Value = str

		isSuccess = true
	}

	return isSuccess
}*/

func (lamisp *ListAttributesMispFormat) SetValueValueAttributesMisp(v interface{}, isNew bool) {
	fmt.Println("func 'SetValueValueAttributesMisp', START")
	fmt.Println("555 isNew = ", isNew, " value = ", v)

	lamisp.mutex.Lock()
	if isNew {
		lamisp.attributes = append(lamisp.attributes, createNewAttributesMisp())
	}

	lamisp.attributes[lamisp.getCountListAttributesMisp()-1].Value = fmt.Sprint(v)
	lamisp.mutex.Unlock()
}

/*func (amisp *AttributesMispFormat) SetValueUuidAttributesMisp(v interface{}) bool {
	var isSuccess bool

	if str, ok := v.(string); ok {
		amisp.Uuid = str

		isSuccess = true
	}

	return isSuccess
}*/

func (lamisp *ListAttributesMispFormat) SetValueUuidAttributesMisp(v interface{}, isNew bool) {
	fmt.Println("func 'SetValueUuidAttributesMisp', START")
	fmt.Println("666 isNew = ", isNew, " value = ", v)

	lamisp.mutex.Lock()
	if isNew {
		lamisp.attributes = append(lamisp.attributes, createNewAttributesMisp())
	}

	lamisp.attributes[lamisp.getCountListAttributesMisp()-1].Uuid = fmt.Sprint(v)
	lamisp.mutex.Unlock()
}

/*func (amisp *AttributesMispFormat) SetValueTimestampAttributesMisp(v interface{}) bool {
	var isSuccess bool

	if str, ok := v.(string); ok {
		amisp.Timestamp = str

		isSuccess = true
	}

	return isSuccess
}*/

func (lamisp *ListAttributesMispFormat) SetValueTimestampAttributesMisp(v interface{}, isNew bool) {
	fmt.Println("func 'SetValueTimestampAttributesMisp', START")
	fmt.Println("777 isNew = ", isNew, " value = ", v)

	lamisp.mutex.Lock()
	if isNew {
		lamisp.attributes = append(lamisp.attributes, createNewAttributesMisp())
	}

	if dt, ok := v.(float64); ok {
		lamisp.attributes[lamisp.getCountListAttributesMisp()-1].Timestamp = fmt.Sprintf("%10.f", dt)[:10]
	}

	lamisp.mutex.Unlock()
}

/*func (amisp *AttributesMispFormat) SetValueDistributionAttributesMisp(v interface{}) bool {
	var isSuccess bool

	if str, ok := v.(string); ok {
		amisp.Distribution = str

		isSuccess = true
	}

	return isSuccess
}*/

func (lamisp *ListAttributesMispFormat) SetValueDistributionAttributesMisp(v interface{}, isNew bool) {
	fmt.Println("func 'SetValueDistributionAttributesMisp', START")
	fmt.Println("888 isNew = ", isNew, " value = ", v)

	lamisp.mutex.Lock()
	if isNew {
		lamisp.attributes = append(lamisp.attributes, createNewAttributesMisp())
	}

	lamisp.attributes[lamisp.getCountListAttributesMisp()-1].Distribution = fmt.Sprint(v)
	lamisp.mutex.Unlock()
}

/*func (amisp *AttributesMispFormat) SetValueSharingGroupIdAttributesMisp(v interface{}) bool {
	var isSuccess bool

	if str, ok := v.(string); ok {
		amisp.SharingGroupId = str

		isSuccess = true
	}

	return isSuccess
}*/

func (lamisp *ListAttributesMispFormat) SetValueSharingGroupIdAttributesMisp(v interface{}, isNew bool) {
	fmt.Println("func 'SetValueSharingGroupIdAttributesMisp', START")
	fmt.Println("999 isNew = ", isNew, " value = ", v)

	lamisp.mutex.Lock()
	if isNew {
		lamisp.attributes = append(lamisp.attributes, createNewAttributesMisp())
	}

	lamisp.attributes[lamisp.getCountListAttributesMisp()-1].SharingGroupId = fmt.Sprint(v)
	lamisp.mutex.Unlock()
}

/*func (amisp *AttributesMispFormat) SetValueCommentAttributesMisp(v interface{}) bool {
	var isSuccess bool

	if str, ok := v.(string); ok {
		amisp.Comment = str

		isSuccess = true
	}

	return isSuccess
}*/

func (lamisp *ListAttributesMispFormat) SetValueCommentAttributesMisp(v interface{}, isNew bool) {
	fmt.Println("func 'SetValueCommentAttributesMisp', START")
	fmt.Println("1010 isNew = ", isNew, " value = ", v)

	lamisp.mutex.Lock()
	if isNew {
		lamisp.attributes = append(lamisp.attributes, createNewAttributesMisp())
	}

	lamisp.attributes[lamisp.getCountListAttributesMisp()-1].Comment = fmt.Sprint(v)
	lamisp.mutex.Unlock()
}

/*func (amisp *AttributesMispFormat) SetValueFirstSeenAttributesMisp(v interface{}) bool {
	var isSuccess bool

	if str, ok := v.(string); ok {
		amisp.FirstSeen = str

		isSuccess = true
	}

	return isSuccess
}*/

func (lamisp *ListAttributesMispFormat) SetValueFirstSeenAttributesMisp(v interface{}, isNew bool) {
	fmt.Println("func 'SetValueFirstSeenAttributesMisp', START")
	fmt.Println("1111 isNew = ", isNew, " value = ", v)

	lamisp.mutex.Lock()
	if isNew {
		lamisp.attributes = append(lamisp.attributes, createNewAttributesMisp())
	}

	if dt, ok := v.(float64); ok {
		lamisp.attributes[lamisp.getCountListAttributesMisp()-1].FirstSeen = fmt.Sprintf("%13.f000", dt)
	}

	lamisp.mutex.Unlock()
}

/*func (amisp *AttributesMispFormat) SetValueLastSeenAttributesMisp(v interface{}) bool {
	var isSuccess bool

	if str, ok := v.(string); ok {
		amisp.LastSeen = str

		isSuccess = true
	}

	return isSuccess
}*/

func (lamisp *ListAttributesMispFormat) SetValueLastSeenAttributesMisp(v interface{}, isNew bool) {
	fmt.Println("func 'SetValueLastSeenAttributesMisp', START")
	fmt.Println("1212 isNew = ", isNew, " value = ", v)

	lamisp.mutex.Lock()
	if isNew {
		lamisp.attributes = append(lamisp.attributes, createNewAttributesMisp())
	}

	if dt, ok := v.(float64); ok {
		lamisp.attributes[lamisp.getCountListAttributesMisp()-1].LastSeen = fmt.Sprintf("%13.f000", dt)
	}

	lamisp.mutex.Unlock()
}

/*func (amisp *AttributesMispFormat) SetValueToIdsAttributesMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(bool); ok {
		amisp.ToIds = data

		isSuccess = true
	}

	return isSuccess
}*/

func (lamisp *ListAttributesMispFormat) SetValueToIdsAttributesMisp(v interface{}, isNew bool) {
	fmt.Println("func 'SetValueToIdsAttributesMisp', START")
	fmt.Println("1313 isNew = ", isNew)

	lamisp.mutex.Lock()
	if isNew {
		lamisp.attributes = append(lamisp.attributes, createNewAttributesMisp())
	}

	if data, ok := v.(bool); ok {
		lamisp.attributes[lamisp.getCountListAttributesMisp()-1].ToIds = data
	}
	lamisp.mutex.Unlock()
}

/*func (amisp *AttributesMispFormat) SetValueDeletedAttributesMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(bool); ok {
		amisp.Deleted = data

		isSuccess = true
	}

	return isSuccess
}*/

func (lamisp *ListAttributesMispFormat) SetValueDeletedAttributesMisp(v interface{}, isNew bool) {
	fmt.Println("func 'SetValueDeletedAttributesMisp', START")
	fmt.Println("1414 isNew = ", isNew)

	lamisp.mutex.Lock()
	if isNew {
		lamisp.attributes = append(lamisp.attributes, createNewAttributesMisp())
	}

	if data, ok := v.(bool); ok {
		lamisp.attributes[lamisp.getCountListAttributesMisp()-1].Deleted = data
	}
	lamisp.mutex.Unlock()
}

/*func (amisp *AttributesMispFormat) SetValueDisableCorrelationAttributesMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(bool); ok {
		amisp.DisableCorrelation = data

		isSuccess = true
	}

	return isSuccess
}*/

func (lamisp *ListAttributesMispFormat) SetValueDisableCorrelationAttributesMisp(v interface{}, isNew bool) {
	fmt.Println("func 'SetValueDisableCorrelationAttributesMisp', START")

	lamisp.mutex.Lock()
	if isNew {
		lamisp.attributes = append(lamisp.attributes, createNewAttributesMisp())
	}

	if data, ok := v.(bool); ok {
		lamisp.attributes[lamisp.getCountListAttributesMisp()-1].DisableCorrelation = data
	}
	lamisp.mutex.Unlock()
}
