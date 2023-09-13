package datamodels

import (
	"fmt"
	"regexp"
	"sync"
	"time"
)

func NewListAttributesMispFormat() *ListAttributesMispFormat {
	return &ListAttributesMispFormat{
		attributes:    make([]AttributesMispFormat, 0),
		attributeTags: make(map[int][][2]string),
		mutex:         sync.Mutex{},
	}
}

func createNewAttributesMisp() AttributesMispFormat {
	return AttributesMispFormat{
		Category:       "Other",
		Type:           "other",
		Timestamp:      "0",
		Distribution:   "3",
		FirstSeen:      fmt.Sprint(time.Now().Format(time.RFC3339)),
		LastSeen:       fmt.Sprint(time.Now().Format(time.RFC3339)),
		ToIds:          true,
		SharingGroupId: "1",
	}
}

func (lambda *ListAttributesMispFormat) GetCountListAttributesMisp() int {
	return len(lambda.attributes)
}

func (lambda *ListAttributesMispFormat) GetCountListAttributeTags() int {
	return len(lambda.attributeTags)
}

func (lambda *ListAttributesMispFormat) CleanListAttributesMisp() {
	lambda.mutex.Lock()
	defer lambda.mutex.Unlock()

	lambda.attributes = []AttributesMispFormat{}
}

func (lamisp *ListAttributesMispFormat) GetListAttributesMisp() []AttributesMispFormat {
	return lamisp.attributes
}

func (lamisp *ListAttributesMispFormat) SetValueObjectIdAttributesMisp(v interface{}, isNew bool) {
	lamisp.mutex.Lock()
	defer lamisp.mutex.Unlock()

	if isNew {
		lamisp.attributes = append(lamisp.attributes, createNewAttributesMisp())
	}

	//lamisp.attributes[getCountList(lamisp.GetCountListAttributesMisp())].ObjectId = fmt.Sprint(v)
	lamisp.attributes[lamisp.GetCountListAttributesMisp()-1].ObjectId = fmt.Sprint(v)
}

func (lamisp *ListAttributesMispFormat) SetValueObjectRelationAttributesMisp(v interface{}, isNew bool) {
	lamisp.mutex.Lock()
	defer lamisp.mutex.Unlock()

	if isNew {
		lamisp.attributes = append(lamisp.attributes, createNewAttributesMisp())
	}

	//lamisp.attributes[getCountList(lamisp.GetCountListAttributesMisp())].ObjectRelation = fmt.Sprint(v)
	lamisp.attributes[lamisp.GetCountListAttributesMisp()-1].ObjectRelation = fmt.Sprint(v)
}

func (lamisp *ListAttributesMispFormat) SetValueCategoryAttributesMisp(v interface{}, isNew bool) {
	lamisp.mutex.Lock()
	defer lamisp.mutex.Unlock()

	if isNew {
		lamisp.attributes = append(lamisp.attributes, createNewAttributesMisp())
	}

	//lamisp.attributes[getCountList(lamisp.GetCountListAttributesMisp())].Category = fmt.Sprint(v)
	lamisp.attributes[lamisp.GetCountListAttributesMisp()-1].Category = fmt.Sprint(v)
}

func (lamisp *ListAttributesMispFormat) SetValueTypeAttributesMisp(v interface{}, isNew bool) {
	lamisp.mutex.Lock()
	defer lamisp.mutex.Unlock()

	if isNew {
		lamisp.attributes = append(lamisp.attributes, createNewAttributesMisp())
	}

	//lamisp.attributes[getCountList(lamisp.GetCountListAttributesMisp())].Type = fmt.Sprint(v)
	lamisp.attributes[lamisp.GetCountListAttributesMisp()-1].Type = fmt.Sprint(v)
}

func (lamisp *ListAttributesMispFormat) SetValueValueAttributesMisp(v interface{}, isNew bool) {
	lamisp.mutex.Lock()
	defer lamisp.mutex.Unlock()

	if isNew {
		lamisp.attributes = append(lamisp.attributes, createNewAttributesMisp())
	}

	//lamisp.attributes[getCountList(lamisp.GetCountListAttributesMisp())].Value = fmt.Sprint(v)
	lamisp.attributes[lamisp.GetCountListAttributesMisp()-1].Value = fmt.Sprint(v)
}

func (lamisp *ListAttributesMispFormat) SetValueUuidAttributesMisp(v interface{}, isNew bool) {
	lamisp.mutex.Lock()
	defer lamisp.mutex.Unlock()

	if isNew {
		lamisp.attributes = append(lamisp.attributes, createNewAttributesMisp())
	}

	//lamisp.attributes[getCountList(lamisp.GetCountListAttributesMisp())].Uuid = fmt.Sprint(v)
	lamisp.attributes[lamisp.GetCountListAttributesMisp()-1].Uuid = fmt.Sprint(v)
}

func (lamisp *ListAttributesMispFormat) SetValueTimestampAttributesMisp(v interface{}, isNew bool) {
	lamisp.mutex.Lock()
	defer lamisp.mutex.Unlock()

	if isNew {
		lamisp.attributes = append(lamisp.attributes, createNewAttributesMisp())
	}

	if dt, ok := v.(float64); ok {
		//lamisp.attributes[getCountList(lamisp.GetCountListAttributesMisp())].Timestamp = fmt.Sprintf("%10.f", dt)[:10]
		lamisp.attributes[lamisp.GetCountListAttributesMisp()-1].Timestamp = fmt.Sprintf("%10.f", dt)[:10]
	}
}

func (lamisp *ListAttributesMispFormat) SetValueDistributionAttributesMisp(v interface{}, isNew bool) {
	lamisp.mutex.Lock()
	defer lamisp.mutex.Unlock()

	if isNew {
		lamisp.attributes = append(lamisp.attributes, createNewAttributesMisp())
	}

	//lamisp.attributes[getCountList(lamisp.GetCountListAttributesMisp())].Distribution = fmt.Sprint(v)
	lamisp.attributes[lamisp.GetCountListAttributesMisp()-1].Distribution = fmt.Sprint(v)
}

func (lamisp *ListAttributesMispFormat) SetValueSharingGroupIdAttributesMisp(v interface{}, isNew bool) {
	lamisp.mutex.Lock()
	defer lamisp.mutex.Unlock()

	if isNew {
		lamisp.attributes = append(lamisp.attributes, createNewAttributesMisp())
	}

	//lamisp.attributes[getCountList(lamisp.GetCountListAttributesMisp())].SharingGroupId = fmt.Sprint(v)
	lamisp.attributes[lamisp.GetCountListAttributesMisp()-1].SharingGroupId = fmt.Sprint(v)
}

func (lamisp *ListAttributesMispFormat) SetValueCommentAttributesMisp(v interface{}, isNew bool) {
	lamisp.mutex.Lock()
	defer lamisp.mutex.Unlock()

	if isNew {
		lamisp.attributes = append(lamisp.attributes, createNewAttributesMisp())
	}

	//lamisp.attributes[getCountList(lamisp.GetCountListAttributesMisp())].Comment = fmt.Sprint(v)
	lamisp.attributes[lamisp.GetCountListAttributesMisp()-1].Comment = fmt.Sprint(v)
}

func (lamisp *ListAttributesMispFormat) SetValueFirstSeenAttributesMisp(v interface{}, isNew bool) {
	lamisp.mutex.Lock()
	defer lamisp.mutex.Unlock()

	if isNew {
		lamisp.attributes = append(lamisp.attributes, createNewAttributesMisp())
	}

	if dt, ok := v.(float64); ok {
		//lamisp.attributes[getCountList(lamisp.GetCountListAttributesMisp())].FirstSeen = time.UnixMilli(int64(dt)).Format(time.RFC3339)
		lamisp.attributes[lamisp.GetCountListAttributesMisp()-1].FirstSeen = time.UnixMilli(int64(dt)).Format(time.RFC3339)
	}
}

func (lamisp *ListAttributesMispFormat) SetValueLastSeenAttributesMisp(v interface{}, isNew bool) {
	lamisp.mutex.Lock()
	defer lamisp.mutex.Unlock()

	if isNew {
		lamisp.attributes = append(lamisp.attributes, createNewAttributesMisp())
	}

	if dt, ok := v.(float64); ok {
		//lamisp.attributes[getCountList(lamisp.GetCountListAttributesMisp())].LastSeen = time.UnixMilli(int64(dt)).Format(time.RFC3339)
		lamisp.attributes[lamisp.GetCountListAttributesMisp()-1].LastSeen = time.UnixMilli(int64(dt)).Format(time.RFC3339)
	}
}

func (lamisp *ListAttributesMispFormat) SetValueToIdsAttributesMisp(v interface{}, isNew bool) {
	lamisp.mutex.Lock()
	defer lamisp.mutex.Unlock()

	if isNew {
		lamisp.attributes = append(lamisp.attributes, createNewAttributesMisp())
	}

	if data, ok := v.(bool); ok {
		//lamisp.attributes[getCountList(lamisp.GetCountListAttributesMisp())].ToIds = data
		lamisp.attributes[lamisp.GetCountListAttributesMisp()-1].ToIds = data
	}
}

func (lamisp *ListAttributesMispFormat) SetValueDeletedAttributesMisp(v interface{}, isNew bool) {
	lamisp.mutex.Lock()
	defer lamisp.mutex.Unlock()

	if isNew {
		lamisp.attributes = append(lamisp.attributes, createNewAttributesMisp())
	}

	if data, ok := v.(bool); ok {
		//lamisp.attributes[getCountList(lamisp.GetCountListAttributesMisp())].Deleted = data
		lamisp.attributes[lamisp.GetCountListAttributesMisp()-1].Deleted = data
	}
}

func (lamisp *ListAttributesMispFormat) SetValueDisableCorrelationAttributesMisp(v interface{}, isNew bool) {
	lamisp.mutex.Lock()
	defer lamisp.mutex.Unlock()

	if isNew {
		lamisp.attributes = append(lamisp.attributes, createNewAttributesMisp())
	}

	if data, ok := v.(bool); ok {
		//lamisp.attributes[getCountList(lamisp.GetCountListAttributesMisp())].DisableCorrelation = data
		lamisp.attributes[lamisp.GetCountListAttributesMisp()-1].DisableCorrelation = data
	}
}

func (lamisp *ListAttributesMispFormat) HandlingValueEventIdAttributesMisp(v interface{}, num int) {
	lamisp.mutex.Lock()
	defer lamisp.mutex.Unlock()

	if lamisp.GetCountListAttributesMisp() < num {
		return
	}

	if num > 0 {
		num = num - 1
	}

	lamisp.attributes[num].EventId = fmt.Sprint(v)
}

func (lamisp *ListAttributesMispFormat) HandlingValueTagsAttributesMisp(v interface{}, isNew bool) {
	lamisp.mutex.Lock()
	defer lamisp.mutex.Unlock()

	if l, ok := v.([]string); ok {
		//lamisp.attributeTags[getCountList(lamisp.GetCountListAttributesMisp())] = HandlingListTags(l)
		lamisp.attributeTags[lamisp.GetCountListAttributesMisp()-1] = HandlingListTags(l)
	}
}

func (lamisp *ListAttributesMispFormat) GetListAttributeTags() map[int][][2]string {
	return lamisp.attributeTags
}

func HandlingListTags(l []string) [][2]string {
	nl := make([][2]string, 0, len(l))
	patter := regexp.MustCompile(`^misp:([\w\-].*)=\"([\w\-].*)\"$`)

	for _, v := range l {
		if !patter.MatchString(v) {
			continue
		}

		result := patter.FindAllStringSubmatch(v, -1)
		if len(result) > 0 && len(result[0]) == 3 {
			nl = append(nl, [2]string{result[0][1], result[0][2]})
		}
	}

	return nl
}

func getCountList(num int) int {
	if num == 0 {
		return 0
	}

	return num - 1
}
