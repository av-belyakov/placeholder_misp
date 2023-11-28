package datamodels

import (
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"
)

func NewListAttributesMispFormat() *ListAttributesMispFormat {
	return &ListAttributesMispFormat{
		attributes: make(map[int]AttributesMispFormat)}
}

func createNewAttributesMisp() AttributesMispFormat {
	return AttributesMispFormat{
		Category:       "Other",
		Type:           "other",
		Timestamp:      "0",
		Distribution:   getDistributionAttributes(),
		Uuid:           uuid.New().String(),
		FirstSeen:      fmt.Sprint(time.Now().Format(time.RFC3339)),
		LastSeen:       fmt.Sprint(time.Now().Format(time.RFC3339)),
		ToIds:          true,
		SharingGroupId: "1",
	}
}

func (lambda *ListAttributesMispFormat) GetCountListAttributesMisp() int {
	return len(lambda.attributes)
}

func (lambda *ListAttributesMispFormat) CleanListAttributesMisp() {
	lambda.Lock()
	defer lambda.Unlock()

	lambda.attributes = map[int]AttributesMispFormat{}
}

func (lambda *ListAttributesMispFormat) DelElementListAttributesMisp(num int) (AttributesMispFormat, bool) {
	lambda.Lock()
	defer lambda.Unlock()

	var (
		ok   bool
		attr AttributesMispFormat
	)

	if attr, ok = lambda.attributes[num]; ok {
		delete(lambda.attributes, num)
	}

	return attr, ok
}

func (lamisp *ListAttributesMispFormat) GetListAttributesMisp() map[int]AttributesMispFormat {
	return lamisp.attributes
}

func (lamisp *ListAttributesMispFormat) SetValueObjectIdAttributesMisp(v interface{}, num int) {
	var tmp AttributesMispFormat
	lamisp.Lock()
	defer lamisp.Unlock()

	if attr, ok := lamisp.attributes[num]; ok {
		tmp = attr
	} else {
		tmp = createNewAttributesMisp()
	}

	tmp.ObjectId = fmt.Sprint(v)
	lamisp.attributes[num] = tmp
}

func (lamisp *ListAttributesMispFormat) SetValueObjectRelationAttributesMisp(v interface{}, num int) {
	var tmp AttributesMispFormat
	lamisp.Lock()
	defer lamisp.Unlock()

	if attr, ok := lamisp.attributes[num]; ok {
		tmp = attr
	} else {
		tmp = createNewAttributesMisp()
	}

	tmp.ObjectRelation = fmt.Sprint(v)
	lamisp.attributes[num] = tmp
}

func (lamisp *ListAttributesMispFormat) SetValueCategoryAttributesMisp(v interface{}, num int) {
	var tmp AttributesMispFormat
	lamisp.Lock()
	defer lamisp.Unlock()

	if attr, ok := lamisp.attributes[num]; ok {
		tmp = attr
	} else {
		tmp = createNewAttributesMisp()
	}

	tmp.Category = fmt.Sprint(v)
	lamisp.attributes[num] = tmp
}

func (lamisp *ListAttributesMispFormat) SetValueTypeAttributesMisp(v interface{}, num int) {
	var tmp AttributesMispFormat
	lamisp.Lock()
	defer lamisp.Unlock()

	if attr, ok := lamisp.attributes[num]; ok {
		tmp = attr
	} else {
		tmp = createNewAttributesMisp()
	}

	tmp.Type = fmt.Sprint(v)
	lamisp.attributes[num] = tmp
}

func (lamisp *ListAttributesMispFormat) SetValueValueAttributesMisp(v interface{}, num int) {
	var tmp AttributesMispFormat
	lamisp.Lock()
	defer lamisp.Unlock()

	if attr, ok := lamisp.attributes[num]; ok {
		tmp = attr
	} else {
		tmp = createNewAttributesMisp()
	}

	tmp.Value = fmt.Sprint(v)
	lamisp.attributes[num] = tmp
}

func (lamisp *ListAttributesMispFormat) SetValueUuidAttributesMisp(v interface{}, num int) {
	var tmp AttributesMispFormat
	lamisp.Lock()
	defer lamisp.Unlock()

	if attr, ok := lamisp.attributes[num]; ok {
		tmp = attr
	} else {
		tmp = createNewAttributesMisp()
	}

	tmp.Uuid = fmt.Sprint(v)
	lamisp.attributes[num] = tmp
}

func (lamisp *ListAttributesMispFormat) SetValueTimestampAttributesMisp(v interface{}, num int) {
	var tmp AttributesMispFormat
	lamisp.Lock()
	defer lamisp.Unlock()

	if attr, ok := lamisp.attributes[num]; ok {
		tmp = attr
	} else {
		tmp = createNewAttributesMisp()
	}

	if dt, ok := v.(float64); ok {
		tmp.Timestamp = fmt.Sprintf("%10.f", dt)[:10]
	}

	lamisp.attributes[num] = tmp
}

func (lamisp *ListAttributesMispFormat) SetValueDistributionAttributesMisp(v interface{}, num int) {
	var tmp AttributesMispFormat
	lamisp.Lock()
	defer lamisp.Unlock()

	if attr, ok := lamisp.attributes[num]; ok {
		tmp = attr
	} else {
		tmp = createNewAttributesMisp()
	}

	tmp.Distribution = fmt.Sprint(v)
	lamisp.attributes[num] = tmp
}

func (lamisp *ListAttributesMispFormat) SetValueSharingGroupIdAttributesMisp(v interface{}, num int) {
	var tmp AttributesMispFormat
	lamisp.Lock()
	defer lamisp.Unlock()

	if attr, ok := lamisp.attributes[num]; ok {
		tmp = attr
	} else {
		tmp = createNewAttributesMisp()
	}

	tmp.SharingGroupId = fmt.Sprint(v)
	lamisp.attributes[num] = tmp
}

func (lamisp *ListAttributesMispFormat) SetValueCommentAttributesMisp(v interface{}, num int) {
	var tmp AttributesMispFormat
	lamisp.Lock()
	defer lamisp.Unlock()

	if attr, ok := lamisp.attributes[num]; ok {
		tmp = attr
	} else {
		tmp = createNewAttributesMisp()
	}

	tmp.Comment = fmt.Sprint(v)
	lamisp.attributes[num] = tmp
}

func (lamisp *ListAttributesMispFormat) SetValueFirstSeenAttributesMisp(v interface{}, num int) {
	var tmp AttributesMispFormat
	lamisp.Lock()
	defer lamisp.Unlock()

	if attr, ok := lamisp.attributes[num]; ok {
		tmp = attr
	} else {
		tmp = createNewAttributesMisp()
	}

	if dt, ok := v.(float64); ok {
		tmp.FirstSeen = time.UnixMilli(int64(dt)).Format(time.RFC3339)
	}

	lamisp.attributes[num] = tmp
}

func (lamisp *ListAttributesMispFormat) SetValueLastSeenAttributesMisp(v interface{}, num int) {
	var tmp AttributesMispFormat
	lamisp.Lock()
	defer lamisp.Unlock()

	if attr, ok := lamisp.attributes[num]; ok {
		tmp = attr
	} else {
		tmp = createNewAttributesMisp()
	}

	if dt, ok := v.(float64); ok {
		tmp.LastSeen = time.UnixMilli(int64(dt)).Format(time.RFC3339)
	}

	lamisp.attributes[num] = tmp
}

func (lamisp *ListAttributesMispFormat) SetValueToIdsAttributesMisp(v interface{}, num int) {
	var tmp AttributesMispFormat
	lamisp.Lock()
	defer lamisp.Unlock()

	if attr, ok := lamisp.attributes[num]; ok {
		tmp = attr
	} else {
		tmp = createNewAttributesMisp()
	}

	if data, ok := v.(bool); ok {
		tmp.ToIds = data
	}

	lamisp.attributes[num] = tmp
}

func (lamisp *ListAttributesMispFormat) SetValueDeletedAttributesMisp(v interface{}, num int) {
	var tmp AttributesMispFormat
	lamisp.Lock()
	defer lamisp.Unlock()

	if attr, ok := lamisp.attributes[num]; ok {
		tmp = attr
	} else {
		tmp = createNewAttributesMisp()
	}

	if data, ok := v.(bool); ok {
		tmp.Deleted = data
	}

	lamisp.attributes[num] = tmp
}

func (lamisp *ListAttributesMispFormat) SetValueDisableCorrelationAttributesMisp(v interface{}, num int) {
	var tmp AttributesMispFormat
	lamisp.Lock()
	defer lamisp.Unlock()

	if attr, ok := lamisp.attributes[num]; ok {
		tmp = attr
	} else {
		tmp = createNewAttributesMisp()
	}

	if data, ok := v.(bool); ok {
		tmp.DisableCorrelation = data
	}

	lamisp.attributes[num] = tmp
}

func (lamisp *ListAttributesMispFormat) HandlingValueEventIdAttributesMisp(v interface{}, num int) {
	var tmp AttributesMispFormat
	lamisp.Lock()
	defer lamisp.Unlock()

	if attr, ok := lamisp.attributes[num]; ok {
		tmp = attr
	} else {
		tmp = createNewAttributesMisp()
	}

	tmp.EventId = fmt.Sprint(v)
	lamisp.attributes[num] = tmp
}

// HandlingValueDataTypeAttributesMisp изменяет некоторые поля объекта типа Attributes
// при этом на эти поля возможно вляиние других функций корректировщиков, например,
// функции getNewListAttributes, которая применяется для совмещения списков Attributes
// и Tags
func (lamisp *ListAttributesMispFormat) HandlingValueDataTypeAttributesMisp(v interface{}, num int) {
	snortSidCategory := func(lamisp *ListAttributesMispFormat, num int) {
		lamisp.SetValueCategoryAttributesMisp("Network activity", num)
	}
	snortSidType := func(lamisp *ListAttributesMispFormat, num int) {
		lamisp.SetValueTypeAttributesMisp("snort", num)
	}

	collection := map[string][]func(lamisp *ListAttributesMispFormat, num int){
		"snort_sid": {snortSidCategory, snortSidType},
	}

	if l, ok := collection[fmt.Sprint(v)]; ok {
		for _, f := range l {
			f(lamisp, num)
		}
	}
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

func getDistributionAttributes() string {
	return "2"
}
