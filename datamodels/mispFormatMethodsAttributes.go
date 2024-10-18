package datamodels

import (
	"fmt"
	"placeholder_misp/supportingfunctions"
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

func (lamisp *ListAttributesMispFormat) AutoSetValueCategoryAttributesMisp(v string, num int) {
	networkActivityCategory := func(lamisp *ListAttributesMispFormat, num int) {
		lamisp.SetValueCategoryAttributesMisp("Network activity", num)
	}
	payloadDeliveryCategory := func(lamisp *ListAttributesMispFormat, num int) {
		lamisp.SetValueCategoryAttributesMisp("Payload delivery", num)
	}

	collection := map[string]func(lamisp *ListAttributesMispFormat, num int){
		"snort_sid": networkActivityCategory,
		"url":       payloadDeliveryCategory,
		"domain":    networkActivityCategory,
		"md5":       payloadDeliveryCategory,
		"sha1":      payloadDeliveryCategory,
		"sha224":    payloadDeliveryCategory,
		"sha256":    payloadDeliveryCategory,
		"sha512":    payloadDeliveryCategory,
		"filename":  payloadDeliveryCategory,
		"ja3":       payloadDeliveryCategory,
		"ip_home":   networkActivityCategory,
	}

	if f, ok := collection[v]; ok {
		f(lamisp, num)
	}
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

func (lamisp *ListAttributesMispFormat) AutoSetValueTypeAttributesMisp(v string, num int) {
	snortSidType := func(lamisp *ListAttributesMispFormat, num int) {
		lamisp.SetValueTypeAttributesMisp("snort", num)
	}
	urlType := func(lamisp *ListAttributesMispFormat, num int) {
		lamisp.SetValueTypeAttributesMisp("url", num)
	}
	domainType := func(lamisp *ListAttributesMispFormat, num int) {
		lamisp.SetValueTypeAttributesMisp("domain", num)
	}
	md5Type := func(lamisp *ListAttributesMispFormat, num int) {
		lamisp.SetValueTypeAttributesMisp("md5", num)
	}
	sha1Type := func(lamisp *ListAttributesMispFormat, num int) {
		lamisp.SetValueTypeAttributesMisp("sha1", num)
	}
	sha224Type := func(lamisp *ListAttributesMispFormat, num int) {
		lamisp.SetValueTypeAttributesMisp("sha224", num)
	}
	sha256Type := func(lamisp *ListAttributesMispFormat, num int) {
		lamisp.SetValueTypeAttributesMisp("sha256", num)
	}
	sha512Type := func(lamisp *ListAttributesMispFormat, num int) {
		lamisp.SetValueTypeAttributesMisp("sha512", num)
	}
	filenameType := func(lamisp *ListAttributesMispFormat, num int) {
		lamisp.SetValueTypeAttributesMisp("filename", num)
	}
	ja3Type := func(lamisp *ListAttributesMispFormat, num int) {
		lamisp.SetValueTypeAttributesMisp("ja3-fingerprint-md5", num)
	}
	ipHome := func(lamisp *ListAttributesMispFormat, num int) {
		lamisp.SetValueTypeAttributesMisp("other", num)
	}

	collection := map[string]func(lamisp *ListAttributesMispFormat, num int){
		"snort_sid": snortSidType,
		"url":       urlType,
		"domain":    domainType,
		"md5":       md5Type,
		"sha1":      sha1Type,
		"sha224":    sha224Type,
		"sha256":    sha256Type,
		"sha512":    sha512Type,
		"filename":  filenameType,
		"ja3":       ja3Type,
		"ip_home":   ipHome,
	}

	if f, ok := collection[v]; ok {
		f(lamisp, num)
	}

	//это для определения типа хеша
	if v == "hash" {
		if hashName, _, err := supportingfunctions.CheckStringHash(v); err == nil {
			lamisp.SetValueTypeAttributesMisp(hashName, num)
		}
	}
}

func (lamisp *ListAttributesMispFormat) SetValueValueAttributesMisp(v interface{}, num int) {
	var tmp AttributesMispFormat
	lamisp.Lock()

	if attr, ok := lamisp.attributes[num]; ok {
		tmp = attr
	} else {
		tmp = createNewAttributesMisp()
	}

	value := fmt.Sprint(v)

	tmp.Value = value
	lamisp.attributes[num] = tmp

	//надо разблокировать Mutex до того как использовать lamisp.AutoSetValueCategoryAttributesMisp и
	//AutoSetValueTypeAttributesMisp так как эти два метода используют методы
	//AutoSetValueCategoryAttributesMisp и AutoSetValueTypeAttributesMisp вызывающие
	//повторную блокировку Mutex
	lamisp.Unlock()

	//дополнительно, если значение подподает под рег. выражение типа "8030073:193.29.19.55"
	//то устанавливаем дополнительное значение типа в поле "object_relation"
	patter := regexp.MustCompile(`^[\d]+:((25[0-5]|2[0-4]\d|[01]?\d\d?)[.]){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)$`)
	if patter.MatchString(value) {
		//выполняем автоматическое изменение значения свойства Category
		lamisp.AutoSetValueCategoryAttributesMisp("ip_home", num)

		//выполняем автоматическое изменение значения свойства Type
		lamisp.AutoSetValueTypeAttributesMisp("ip_home", num)

		np := regexp.MustCompile(`^([\d]+):([\d]+\.[\d]+\.[\d]+\.[\d]+)$`)
		result := np.FindAllStringSubmatch(value, -1)
		if len(result) > 0 && len(result[0]) == 3 {
			lamisp.SetValueValueAttributesMisp(result[0][2], num)
		}
	}
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
func (lamisp *ListAttributesMispFormat) HandlingValueDataTypeAttributesMisp(i interface{}, num int) {
	v := fmt.Sprint(i)

	//выполняем автоматическое изменение значения свойства Category
	lamisp.AutoSetValueCategoryAttributesMisp(v, num)

	//выполняем автоматическое изменение значения свойства Type
	lamisp.AutoSetValueTypeAttributesMisp(v, num)

	if v == "ip_home" {
		lamisp.SetValueObjectRelationAttributesMisp(v, num)
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
