package datamodels

import (
	"fmt"
	"placeholder_misp/supportingfunctions"
)

func (o *ObservablesMessageTheHive) Get() []ObservableMessage {
	return o.Observables
}

func (o *ObservablesMessageTheHive) Set(v ObservableMessage) {
	o.Observables = append(o.Observables, v)
}

func (o *ObservableMessage) GetIoc() bool {
	return o.Ioc
}

// SetValueIoc устанавливает BOOL значение для поля Ioc
func (o *ObservableMessage) SetValueIoc(v bool) {
	o.Ioc = v
}

// SetAnyIoc устанавливает ЛЮБОЕ значение для поля Ioc
func (o *ObservableMessage) SetAnyIoc(i interface{}) {
	if v, ok := i.(bool); ok {
		o.Ioc = v
	}
}

func (o *ObservableMessage) GetSighted() bool {
	return o.Sighted
}

// SetValueSighted устанавливает BOOL значение для поля Sighted
func (o *ObservableMessage) SetValueSighted(v bool) {
	o.Sighted = v
}

// SetAnySighted устанавливает ЛЮБОЕ значение для поля Sighted
func (o *ObservableMessage) SetAnySighted(i interface{}) {
	if v, ok := i.(bool); ok {
		o.Sighted = v
	}
}

func (o *ObservableMessage) GetIgnoreSimilarity() bool {
	return o.IgnoreSimilarity
}

// SetValueIgnoreSimilarity устанавливает BOOL значение для поля IgnoreSimilarity
func (o *ObservableMessage) SetValueIgnoreSimilarity(v bool) {
	o.IgnoreSimilarity = v
}

// SetAnyIgnoreSimilarity устанавливает ЛЮБОЕ значение для поля IgnoreSimilarity
func (o *ObservableMessage) SetAnyIgnoreSimilarity(i interface{}) {
	if v, ok := i.(bool); ok {
		o.IgnoreSimilarity = v
	}
}

func (o *ObservableMessage) GetTlp() int {
	return o.Tlp
}

// SetValueTlp устанавливает INT значение для поля Tlp
func (o *ObservableMessage) SetValueTlp(v int) {
	o.Tlp = v
}

// SetAnyTlp устанавливает ЛЮБОЕ значение для поля Tlp
func (o *ObservableMessage) SetAnyTlp(i interface{}) {
	if v, ok := i.(int); ok {
		o.Tlp = v
	}
}

func (o *ObservableMessage) GetCreatedAt() uint64 {
	return o.CreatedAt
}

// SetValueCreatedAt устанавливает UINT64 значение для поля CreatedAt
func (o *ObservableMessage) SetValueCreatedAt(v uint64) {
	o.CreatedAt = v
}

// SetAnyCreatedAt устанавливает ЛЮБОЕ значение для поля CreatedAt
func (o *ObservableMessage) SetAnyCreatedAt(i interface{}) {
	if v, ok := i.(uint64); ok {
		o.CreatedAt = v
	}
}

func (o *ObservableMessage) GetUpdatedAt() uint64 {
	return o.UpdatedAt
}

// SetValueUpdatedAt устанавливает UINT64 значение для поля UpdatedAt
func (o *ObservableMessage) SetValueUpdatedAt(v uint64) {
	o.UpdatedAt = v
}

// SetAnyUpdatedAt устанавливает ЛЮБОЕ значение для поля UpdatedAt
func (o *ObservableMessage) SetAnyUpdatedAt(i interface{}) {
	if v, ok := i.(uint64); ok {
		o.UpdatedAt = v
	}
}

func (o *ObservableMessage) GetStartDate() uint64 {
	return o.StartDate
}

// SetValueStartDate устанавливает UINT64 значение для поля StartDate
func (o *ObservableMessage) SetValueStartDate(v uint64) {
	o.StartDate = v
}

// SetAnyStartDate устанавливает ЛЮБОЕ значение для поля StartDate
func (o *ObservableMessage) SetAnyStartDate(i interface{}) {
	if v, ok := i.(uint64); ok {
		o.StartDate = v
	}
}

func (o *ObservableMessage) GetCreatedBy() string {
	return o.CreatedBy
}

// SetValueCreatedBy устанавливает STRING значение для поля CreatedBy
func (o *ObservableMessage) SetValueCreatedBy(v string) {
	o.CreatedBy = v
}

// SetAnyCreatedBy устанавливает ЛЮБОЕ значение для поля CreatedBy
func (o *ObservableMessage) SetAnyCreatedBy(i interface{}) {
	o.CreatedBy = fmt.Sprint(i)
}

func (o *ObservableMessage) GetUpdatedBy() string {
	return o.UpdatedBy
}

// SetValueUpdatedBy устанавливает STRING значение для поля UpdatedBy
func (o *ObservableMessage) SetValueUpdatedBy(v string) {
	o.UpdatedBy = v
}

// SetAnyUpdatedBy устанавливает ЛЮБОЕ значение для поля UpdatedBy
func (o *ObservableMessage) SetAnyUpdatedBy(i interface{}) {
	o.UpdatedBy = fmt.Sprint(i)
}

func (o *ObservableMessage) GetUnderliningId() string {
	return o.UnderliningId
}

// SetValueUnderliningId устанавливает STRING значение для поля UnderliningId
func (o *ObservableMessage) SetValueUnderliningId(v string) {
	o.UnderliningId = v
}

// SetAnyUnderliningId устанавливает ЛЮБОЕ значение для поля UnderliningId
func (o *ObservableMessage) SetAnyUnderliningId(i interface{}) {
	o.UnderliningId = fmt.Sprint(i)
}

func (o *ObservableMessage) GetUnderliningType() string {
	return o.UnderliningType
}

// SetValueUnderliningType устанавливает STRING значение для поля UnderliningType
func (o *ObservableMessage) SetValueUnderliningType(v string) {
	o.UnderliningType = v
}

// SetAnyUnderliningType устанавливает ЛЮБОЕ значение для поля UnderliningType
func (o *ObservableMessage) SetAnyUnderliningType(i interface{}) {
	o.UnderliningType = fmt.Sprint(i)
}

func (o *ObservableMessage) GetData() string {
	return o.Data
}

// SetValueData устанавливает STRING значение для поля Data
func (o *ObservableMessage) SetValueData(v string) {
	o.Data = v
}

// SetAnyData устанавливает ЛЮБОЕ значение для поля Data
func (o *ObservableMessage) SetAnyData(i interface{}) {
	o.Data = fmt.Sprint(i)
}

func (o *ObservableMessage) GetDataType() string {
	return o.DataType
}

// SetValueDataType устанавливает STRING значение для поля DataType
func (o *ObservableMessage) SetValueDataType(v string) {
	o.DataType = v
}

// SetAnyDataType устанавливает ЛЮБОЕ значение для поля DataType
func (o *ObservableMessage) SetAnyDataType(i interface{}) {
	o.DataType = fmt.Sprint(i)
}

func (o *ObservableMessage) GetMessage() string {
	return o.Message
}

// SetValueMessage устанавливает STRING значение для поля Message
func (o *ObservableMessage) SetValueMessage(v string) {
	o.Message = v
}

// SetAnyMessage устанавливает ЛЮБОЕ значение для поля Message
func (o *ObservableMessage) SetAnyMessage(i interface{}) {
	o.Message = fmt.Sprint(i)
}

func (o *ObservableMessage) GetTags() []string {
	return o.Tags
}

// SetValueTags устанавливает STRING значение для поля Tags
func (o *ObservableMessage) SetValueTags(v string) {
	o.Tags = append(o.Tags, v)
}

// SetAnyTags устанавливает ЛЮБОЕ значение для поля Tags
func (o *ObservableMessage) SetAnyTags(i interface{}) {
	o.Tags = append(o.Tags, fmt.Sprint(i))
}

func (o *ObservableMessage) GetReports() map[string]ReportTaxonomies {
	return o.Reports
}

// SetValueReports устанавливает значение для поля Reports
func (o *ObservableMessage) SetValueReports(v map[string]ReportTaxonomies) {
	o.Reports = v
}

/*
Еще Reports


func (o *ObservableMessage) Get() string {
	return o.
}

// SetValue устанавливает СТРОКОВОЕ значение для поля
func (o *ObservableMessage) SetValue(v string) {
	o. = v
}

// SetAny устанавливает ЛЮБОЕ значение для поля
func (o *ObservableMessage) SetAny(i interface{}) {
	o. = fmt.Sprint(i)
}
*/

func (om ObservablesMessageTheHive) ToStringBeautiful(num int) string {
	var str string

	for _, v := range om.Observables {
		str += v.ToStringBeautiful(num)
	}

	return str
}

func (om ObservableMessage) ToStringBeautiful(num int) string {
	var str string
	ws := supportingfunctions.GetWhitespace(num)

	str += fmt.Sprintf("%s_createdAt: '%d'\n", ws, om.CreatedAt)
	str += fmt.Sprintf("%s_createdBy: '%s'\n", ws, om.CreatedBy)
	str += fmt.Sprintf("%s_id: '%s'\n", ws, om.UnderliningId)
	str += fmt.Sprintf("%s_type: '%s'\n", ws, om.UnderliningType)
	str += fmt.Sprintf("%s_updatedAt: '%d'\n", ws, om.UpdatedAt)
	str += fmt.Sprintf("%s_updatedBy: '%s'\n", ws, om.UpdatedBy)
	str += fmt.Sprintf("%sdata: '%s'\n", ws, om.Data)
	str += fmt.Sprintf("%sdataType: '%s'\n", ws, om.DataType)
	str += fmt.Sprintf("%signoreSimilarity: '%v'\n", ws, om.IgnoreSimilarity)
	//данное поле редко используемое, думаю пока оно не требует реализации
	/*str += fmt.Sprintf("%sextraData: \n%s", ws, func(l map[string]interface{}) string {
		var str string
		ws := supportingfunctions.GetWhitespace(num + 1)

		for k, v := range l {
			str += fmt.Sprintf("%s%s: '%v'\n", ws, k, v)
		}
		return str
	}(om.ExtraData))*/
	str += fmt.Sprintf("%sioc: '%v'\n", ws, om.Ioc)
	str += fmt.Sprintf("%smessage: '%s'\n", ws, om.Message)
	str += fmt.Sprintf("%ssighted: '%v'\n", ws, om.Sighted)
	str += fmt.Sprintf("%sstartDate: '%d'\n", ws, om.StartDate)
	str += fmt.Sprintf("%stags: \n%s", ws, func(l []string) string {
		var str string
		ws := supportingfunctions.GetWhitespace(num + 1)

		for k, v := range l {
			str += fmt.Sprintf("%s%d. '%s'\n", ws, k+1, v)
		}
		return str
	}(om.Tags))
	str += fmt.Sprintf("%stlp: '%d'\n", ws, om.Tlp)
	str += fmt.Sprintf("%sreports: \n%s", ws, func(l map[string]ReportTaxonomies) string {
		var str string
		for key, value := range l {
			str += fmt.Sprintf("%s%s:\n", supportingfunctions.GetWhitespace(num+1), key)
			str += fmt.Sprintf("%staxonomys:\n", supportingfunctions.GetWhitespace(num+2))
			for k, v := range value.Taxonomies {
				str += fmt.Sprintf("%s%d.\n", supportingfunctions.GetWhitespace(num+3), k+1)
				str += fmt.Sprintf("%sLevel: %v\n", supportingfunctions.GetWhitespace(num+4), v.Level)
				str += fmt.Sprintf("%sNamespace: %v\n", supportingfunctions.GetWhitespace(num+4), v.Namespace)
				str += fmt.Sprintf("%sPredicate: %v\n", supportingfunctions.GetWhitespace(num+4), v.Predicate)
				str += fmt.Sprintf("%sValue: %v\n", supportingfunctions.GetWhitespace(num+4), v.Value)
			}
		}
		return str
	}(om.Reports))

	return str
}

// ****************** AttachmentData ********************
func (a *AttachmentData) GetSize() int {
	return a.Size
}

// SetValueSize устанавливает INT значение для поля Size
func (a *AttachmentData) SetValueSize(v int) {
	a.Size = v
}

// SetAnySize устанавливает ЛЮБОЕ значение для поля Size
func (a *AttachmentData) SetAnySize(i interface{}) {
	if v, ok := i.(int); ok {
		a.Size = v
	}
}

func (a *AttachmentData) GetId() string {
	return a.Id
}

// SetValueId устанавливает STRING значение для поля Id
func (a *AttachmentData) SetValueId(v string) {
	a.Id = v
}

// SetAnyId устанавливает ЛЮБОЕ значение для поля Id
func (a *AttachmentData) SetAnyId(i interface{}) {
	a.Id = fmt.Sprint(i)
}

func (a *AttachmentData) GetName() string {
	return a.Name
}

// SetValueName устанавливает STRING значение для поля Name
func (a *AttachmentData) SetValueName(v string) {
	a.Name = v
}

// SetAnyName устанавливает ЛЮБОЕ значение для поля Name
func (a *AttachmentData) SetAnyName(i interface{}) {
	a.Name = fmt.Sprint(i)
}

func (a *AttachmentData) GetContentType() string {
	return a.ContentType
}

// SetValueContentType устанавливает STRING значение для поля ContentType
func (a *AttachmentData) SetValueContentType(v string) {
	a.ContentType = v
}

// SetAnyContentType устанавливает ЛЮБОЕ значение для поля ContentType
func (a *AttachmentData) SetAnyContentType(i interface{}) {
	a.ContentType = fmt.Sprint(i)
}

func (a *AttachmentData) GetHashes() []string {
	return a.Hashes
}

// SetValueHashes устанавливает STRING значение для поля Hashes
func (a *AttachmentData) SetValueHashes(v string) {
	a.Hashes = append(a.Hashes, v)
}

// SetAnyHashes устанавливает ЛЮБОЕ значение для поля Hashes
func (a *AttachmentData) SetAnyHashes(i interface{}) {
	a.Hashes = append(a.Hashes, fmt.Sprint(i))
}

// ********************* ReportTaxonomys *******************
func (t *ReportTaxonomies) GetTaxonomys() []Taxonomy {
	return t.Taxonomies
}

func (t *ReportTaxonomies) GetReportTaxonomys() ReportTaxonomies {
	return *t
}

func (t *ReportTaxonomies) AddTaxonomy(taxonomy Taxonomy) {
	t.Taxonomies = append(t.Taxonomies, taxonomy)
}

// *********************** Taxonomy ************************
func (t *Taxonomy) GetLevel() string {
	return t.Level
}

// SetValueLevel устанавливает STRING значение для поля Level
func (t *Taxonomy) SetValueLevel(v string) {
	t.Level = v
}

// SetAnyLevel устанавливает ЛЮБОЕ значение для поля Level
func (t *Taxonomy) SetAnyLevel(i interface{}) {
	t.Level = fmt.Sprint(i)
}

func (t *Taxonomy) GetNamespace() string {
	return t.Namespace
}

// SetValueNamespace устанавливает STRING значение для поля Namespace
func (t *Taxonomy) SetValueNamespace(v string) {
	t.Namespace = v
}

// SetAnyNamespace устанавливает ЛЮБОЕ значение для поля Namespace
func (t *Taxonomy) SetAnyNamespace(i interface{}) {
	t.Namespace = fmt.Sprint(i)
}

func (t *Taxonomy) GetPredicate() string {
	return t.Predicate
}

// SetValuePredicate устанавливает STRING значение для поля Predicate
func (t *Taxonomy) SetValuePredicate(v string) {
	t.Predicate = v
}

// SetAnyPredicate устанавливает ЛЮБОЕ значение для поля Predicate
func (t *Taxonomy) SetAnyPredicate(i interface{}) {
	t.Predicate = fmt.Sprint(i)
}

func (t *Taxonomy) GetValue() string {
	return t.Value
}

// SetValueValue устанавливает STRING значение для поля Value
func (t *Taxonomy) SetValueValue(v string) {
	t.Value = v
}

// SetAnyValue устанавливает ЛЮБОЕ значение для поля Value
func (t *Taxonomy) SetAnyValue(i interface{}) {
	t.Value = fmt.Sprint(i)
}
