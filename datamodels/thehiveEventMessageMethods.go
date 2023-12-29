package datamodels

import (
	"fmt"

	"placeholder_misp/supportingfunctions"
)

// Get возвращает объект типа EventMessageTheHive
func (e *EventMessageTheHive) Get() *EventMessageTheHive {
	return e
}

func (e *EventMessageTheHive) GetBase() bool {
	return e.Base
}

// SetValueBase устанавливает BOOL значение для поля Base
func (e *EventMessageTheHive) SetValueBase(v bool) {
	e.Base = v
}

// SetAnyBase устанавливает ЛЮБОЕ значение для поля Base
func (e *EventMessageTheHive) SetAnyBase(i interface{}) {
	if v, ok := i.(bool); ok {
		e.Base = v
	}
}

func (e *EventMessageTheHive) GetStartDate() uint64 {
	return e.StartDate
}

// SetValueStartDate устанавливает UINT64 значение для поля StartDate
func (e *EventMessageTheHive) SetValueStartDate(v uint64) {
	e.StartDate = v
}

// SetAnyStartDate устанавливает ЛЮБОЕ значение для поля StartDate
func (e *EventMessageTheHive) SetAnyStartDate(i interface{}) {
	if v, ok := i.(uint64); ok {
		e.StartDate = v
	}
}

func (e *EventMessageTheHive) GetRootId() string {
	return e.RootId
}

// SetValueRootId устанавливает STRING значение для поля RootId
func (e *EventMessageTheHive) SetValueRootId(v string) {
	e.RootId = v
}

// SetAnyRootId устанавливает ЛЮБОЕ значение для поля RootId
func (e *EventMessageTheHive) SetAnyRootId(i interface{}) {
	e.RootId = fmt.Sprint(i)
}

func (e *EventMessageTheHive) GetOrganisation() string {
	return e.Organisation
}

// SetValueOrganisation устанавливает STRING значение для поля Organisation
func (e *EventMessageTheHive) SetValueOrganisation(v string) {
	e.Organisation = v
}

// SetAnyOrganisation устанавливает ЛЮБОЕ значение для поля Organisation
func (e *EventMessageTheHive) SetAnyOrganisation(i interface{}) {
	e.Organisation = fmt.Sprint(i)
}

func (e *EventMessageTheHive) GetOrganisationId() string {
	return e.OrganisationId
}

// SetValueOrganisationId устанавливает STRING значение для поля OrganisationId
func (e *EventMessageTheHive) SetValueOrganisationId(v string) {
	e.OrganisationId = v
}

// SetAnyOrganisationId устанавливает ЛЮБОЕ значение для поля OrganisationId
func (e *EventMessageTheHive) SetAnyOrganisationId(i interface{}) {
	e.OrganisationId = fmt.Sprint(i)
}

func (e *EventMessageTheHive) GetObjectId() string {
	return e.ObjectId
}

// SetValueObjectId устанавливает STRING значение для поля ObjectId
func (e *EventMessageTheHive) SetValueObjectId(v string) {
	e.ObjectId = v
}

// SetAnyObjectId устанавливает ЛЮБОЕ значение для поля ObjectId
func (e *EventMessageTheHive) SetAnyObjectId(i interface{}) {
	e.ObjectId = fmt.Sprint(i)
}

func (e *EventMessageTheHive) GetObjectType() string {
	return e.ObjectType
}

// SetValueObjectType устанавливает STRING значение для поля ObjectType
func (e *EventMessageTheHive) SetValueObjectType(v string) {
	e.ObjectType = v
}

// SetAnyObjectType устанавливает ЛЮБОЕ значение для поля ObjectType
func (e *EventMessageTheHive) SetAnyObjectType(i interface{}) {
	e.ObjectType = fmt.Sprint(i)
}

func (e *EventMessageTheHive) GetOperation() string {
	return e.Operation
}

// SetValueOperation устанавливает STRING значение для поля Operation
func (e *EventMessageTheHive) SetValueOperation(v string) {
	e.Operation = v
}

// SetAnyOperation устанавливает ЛЮБОЕ значение для поля Operation
func (e *EventMessageTheHive) SetAnyOperation(i interface{}) {
	e.Operation = fmt.Sprint(i)
}

func (e *EventMessageTheHive) GetRequestId() string {
	return e.RequestId
}

// SetValueRequestId устанавливает STRING значение для поля RequestId
func (e *EventMessageTheHive) SetValueRequestId(v string) {
	e.RequestId = v
}

// SetAnyRequestId устанавливает ЛЮБОЕ значение для поля RequestId
func (e *EventMessageTheHive) SetAnyRequestId(i interface{}) {
	e.RequestId = fmt.Sprint(i)
}

func (e *EventMessageTheHive) GetDetails() EventDetails {
	return e.Details
}

// SetValueDetails устанавливает значение типа EventDetails для поля Details
func (e *EventMessageTheHive) SetValueDetails(v EventDetails) {
	e.Details = v
}

func (e *EventMessageTheHive) GetObject() EventObject {
	return e.Object
}

// SetValueObject устанавливает значение типа EventObject для поля Object
func (e *EventMessageTheHive) SetValueObject(v EventObject) {
	e.Object = v
}

func (em EventMessageTheHive) ToStringBeautiful(num int) string {
	var str string

	ws := supportingfunctions.GetWhitespace(num)

	str += fmt.Sprintf("%soperation: '%s'\n", ws, em.Operation)
	str += fmt.Sprintf("%sobjectId: '%s'\n", ws, em.ObjectId)
	str += fmt.Sprintf("%sobjectType: '%s'\n", ws, em.ObjectType)
	str += fmt.Sprintf("%sbase: '%v'\n", ws, em.Base)
	str += fmt.Sprintf("%sstartDate: '%d'\n", ws, em.StartDate)
	str += fmt.Sprintf("%srootId: '%s'\n", ws, em.RootId)
	str += fmt.Sprintf("%srequestId: '%s'\n", ws, em.RequestId)
	str += fmt.Sprintf("%sdetails:\n", ws)
	str += em.Details.ToStringBeautiful(num + 1)
	str += fmt.Sprintf("%sobject:\n", ws)
	str += em.Object.ToStringBeautiful(num + 1)
	str += fmt.Sprintf("%sorganisationId: '%s'\n", ws, em.OrganisationId)
	str += fmt.Sprintf("%sorganisation: '%s'\n", ws, em.Organisation)

	return str
}

func (e *EventDetails) GetEndDate() uint64 {
	return e.EndDate
}

// SetValueEndDate устанавливает UINT64 значение для поля EndDate
func (e *EventDetails) SetValueEndDate(v uint64) {
	e.EndDate = v
}

// SetAnyEndDate устанавливает ЛЮБОЕ значение для поля EndDate
func (e *EventDetails) SetAnyEndDate(i interface{}) {
	if v, ok := i.(uint64); ok {
		e.EndDate = v
	}
}

func (e *EventDetails) GetResolutionStatus() string {
	return e.ResolutionStatus
}

// SetValueResolutionStatus устанавливает STRING значение для поля ResolutionStatus
func (e *EventDetails) SetValueResolutionStatus(v string) {
	e.ResolutionStatus = v
}

// SetAnyResolutionStatus устанавливает ЛЮБОЕ значение для поля ResolutionStatus
func (e *EventDetails) SetAnyResolutionStatus(i interface{}) {
	e.ResolutionStatus = fmt.Sprint(i)
}

func (e *EventDetails) GetSummary() string {
	return e.Summary
}

// SetValueSummary устанавливает STRING значение для поля Summary
func (e *EventDetails) SetValueSummary(v string) {
	e.Summary = v
}

// SetAnySummary устанавливает ЛЮБОЕ значение для поля Summary
func (e *EventDetails) SetAnySummary(i interface{}) {
	e.Summary = fmt.Sprint(i)
}

func (e *EventDetails) GetStatus() string {
	return e.Status
}

// SetValueStatus устанавливает STRING значение для поля Status
func (e *EventDetails) SetValueStatus(v string) {
	e.Status = v
}

// SetAnyStatus устанавливает ЛЮБОЕ значение для поля Status
func (e *EventDetails) SetAnyStatus(i interface{}) {
	e.Status = fmt.Sprint(i)
}

func (e *EventDetails) GetImpactStatus() string {
	return e.ImpactStatus
}

// SetValueImpactStatus устанавливает STRING значение для поля ImpactStatus
func (e *EventDetails) SetValueImpactStatus(v string) {
	e.ImpactStatus = v
}

// SetAnyImpactStatus устанавливает ЛЮБОЕ значение для поля ImpactStatus
func (e *EventDetails) SetAnyImpactStatus(i interface{}) {
	e.ImpactStatus = fmt.Sprint(i)
}

func (ed EventDetails) ToStringBeautiful(num int) string {
	var str string

	ws := supportingfunctions.GetWhitespace(num)

	str += fmt.Sprintf("%sendDate: '%d'\n", ws, ed.EndDate)
	str += fmt.Sprintf("%sresolutionStatus: '%s'\n", ws, ed.ResolutionStatus)
	str += fmt.Sprintf("%ssummary: '%s'\n", ws, ed.Summary)
	str += fmt.Sprintf("%sstatus: '%s'\n", ws, ed.Status)
	str += fmt.Sprintf("%simpactStatus: '%s'\n", ws, ed.ImpactStatus)
	str += fmt.Sprintf("%scustomFields: \n%s", ws, CustomFieldsToStringBeautiful(ed.CustomFields, num))

	return str
}

func (e *EventObject) GetFlag() bool {
	return e.Flag
}

// SetValueFlag устанавливает BOOL значение для поля Flag
func (e *EventObject) SetValueFlag(v bool) {
	e.Flag = v
}

// SetAnyFlag устанавливает ЛЮБОЕ значение для поля Flag
func (e *EventObject) SetAnyFlag(i interface{}) {
	if v, ok := i.(bool); ok {
		e.Flag = v
	}
}

func (e *EventObject) GetCaseId() int {
	return e.CaseId
}

// SetValueCaseId устанавливает INT значение для поля CaseId
func (e *EventObject) SetValueCaseId(v int) {
	e.CaseId = v
}

// SetAnyCaseId устанавливает ЛЮБОЕ значение для поля CaseId
func (e *EventObject) SetAnyCaseId(i interface{}) {
	if v, ok := i.(int); ok {
		e.CaseId = v
	}
}

func (e *EventObject) GetSeverity() int {
	return e.Severity
}

// SetValueSeverity устанавливает INT значение для поля Severity
func (e *EventObject) SetValueSeverity(v int) {
	e.Severity = v
}

// SetAnySeverity устанавливает ЛЮБОЕ значение для поля Severity
func (e *EventObject) SetAnySeverity(i interface{}) {
	if v, ok := i.(int); ok {
		e.Severity = v
	}
}

func (e *EventObject) GetTlp() int {
	return e.Tlp
}

// SetValueTlp устанавливает INT значение для поля Tlp
func (e *EventObject) SetValueTlp(v int) {
	e.Tlp = v
}

// SetAnyTlp устанавливает ЛЮБОЕ значение для поля Tlp
func (e *EventObject) SetAnyTlp(i interface{}) {
	if v, ok := i.(int); ok {
		e.Tlp = v
	}
}

func (e *EventObject) GetPap() int {
	return e.Pap
}

// SetValuePap устанавливает INT значение для поля Pap
func (e *EventObject) SetValuePap(v int) {
	e.Pap = v
}

// SetAnyPap устанавливает ЛЮБОЕ значение для поля Pap
func (e *EventObject) SetAnyPap(i interface{}) {
	if v, ok := i.(int); ok {
		e.Pap = v
	}
}

func (e *EventObject) GetStartDate() uint64 {
	return e.StartDate
}

// SetValueStartDate устанавливает UINT64 значение для поля StartDate
func (e *EventObject) SetValueStartDate(v uint64) {
	e.StartDate = v
}

// SetAnyStartDate устанавливает ЛЮБОЕ значение для поля StartDate
func (e *EventObject) SetAnyStartDate(i interface{}) {
	if v, ok := i.(uint64); ok {
		e.StartDate = v
	}
}

func (e *EventObject) GetEndDate() uint64 {
	return e.EndDate
}

// SetValueEndDate устанавливает UINT64 значение для поля EndDate
func (e *EventObject) SetValueEndDate(v uint64) {
	e.EndDate = v
}

// SetAnyEndDate устанавливает ЛЮБОЕ значение для поля EndDate
func (e *EventObject) SetAnyEndDate(i interface{}) {
	if v, ok := i.(uint64); ok {
		e.EndDate = v
	}
}

func (e *EventObject) GetCreatedAt() uint64 {
	return e.CreatedAt
}

// SetValueCreatedAt устанавливает UINT64 значение для поля CreatedAt
func (e *EventObject) SetValueCreatedAt(v uint64) {
	e.CreatedAt = v
}

// SetAnyCreatedAt устанавливает ЛЮБОЕ значение для поля CreatedAt
func (e *EventObject) SetAnyCreatedAt(i interface{}) {
	if v, ok := i.(uint64); ok {
		e.CreatedAt = v
	}
}

func (e *EventObject) GetUpdatedAt() uint64 {
	return e.UpdatedAt
}

// SetValueUpdatedAt устанавливает UINT64 значение для поля UpdatedAt
func (e *EventObject) SetValueUpdatedAt(v uint64) {
	e.UpdatedAt = v
}

// SetAnyUpdatedAt устанавливает ЛЮБОЕ значение для поля UpdatedAt
func (e *EventObject) SetAnyUpdatedAt(i interface{}) {
	if v, ok := i.(uint64); ok {
		e.UpdatedAt = v
	}
}

func (e *EventObject) GetUnderliningId() string {
	return e.UnderliningId
}

// SetValueUnderliningId устанавливает STRING значение для поля UnderliningId
func (e *EventObject) SetValueUnderliningId(v string) {
	e.UnderliningId = v
}

// SetAnyUnderliningId устанавливает ЛЮБОЕ значение для поля UnderliningId
func (e *EventObject) SetAnyUnderliningId(i interface{}) {
	e.UnderliningId = fmt.Sprint(i)
}

func (e *EventObject) GetId() string {
	return e.Id
}

// SetValueId устанавливает STRING значение для поля Id
func (e *EventObject) SetValueId(v string) {
	e.Id = v
}

// SetAnyId устанавливает ЛЮБОЕ значение для поля Id
func (e *EventObject) SetAnyId(i interface{}) {
	e.Id = fmt.Sprint(i)
}

func (e *EventObject) GetCreatedBy() string {
	return e.CreatedBy
}

// SetValueCreatedBy устанавливает STRING значение для поля CreatedBy
func (e *EventObject) SetValueCreatedBy(v string) {
	e.CreatedBy = v
}

// SetAnyCreatedBy устанавливает ЛЮБОЕ значение для поля CreatedBy
func (e *EventObject) SetAnyCreatedBy(i interface{}) {
	e.CreatedBy = fmt.Sprint(i)
}

func (e *EventObject) GetUpdatedBy() string {
	return e.UpdatedBy
}

// SetValueUpdatedBy устанавливает STRING значение для поля UpdatedBy
func (e *EventObject) SetValueUpdatedBy(v string) {
	e.UpdatedBy = v
}

// SetAnyUpdatedBy устанавливает ЛЮБОЕ значение для поля UpdatedBy
func (e *EventObject) SetAnyUpdatedBy(i interface{}) {
	e.UpdatedBy = fmt.Sprint(i)
}

func (e *EventObject) GetUnderliningType() string {
	return e.ImpactStatus
}

// SetValueUnderliningType устанавливает STRING значение для поля UnderliningType
func (e *EventObject) SetValueUnderliningType(v string) {
	e.UnderliningType = v
}

// SetAnyUnderliningType устанавливает ЛЮБОЕ значение для поля UnderliningType
func (e *EventObject) SetAnyUnderliningType(i interface{}) {
	e.UnderliningType = fmt.Sprint(i)
}

func (e *EventObject) GetTitle() string {
	return e.Title
}

// SetValueTitle устанавливает STRING значение для поля Title
func (e *EventObject) SetValueTitle(v string) {
	e.Title = v
}

// SetAnyTitle устанавливает ЛЮБОЕ значение для поля Title
func (e *EventObject) SetAnyTitle(i interface{}) {
	e.Title = fmt.Sprint(i)
}

func (e *EventObject) GetDescription() string {
	return e.Description
}

// SetValueDescription устанавливает STRING значение для поля Description
func (e *EventObject) SetValueDescription(v string) {
	e.Description = v
}

// SetAnyDescription устанавливает ЛЮБОЕ значение для поля Description
func (e *EventObject) SetAnyDescription(i interface{}) {
	e.Description = fmt.Sprint(i)
}

func (e *EventObject) GetImpactStatus() string {
	return e.ImpactStatus
}

// SetValueImpactStatus устанавливает STRING значение для поля ImpactStatus
func (e *EventObject) SetValueImpactStatus(v string) {
	e.ImpactStatus = v
}

// SetAnyImpactStatus устанавливает ЛЮБОЕ значение для поля ImpactStatus
func (e *EventObject) SetAnyImpactStatus(i interface{}) {
	e.ImpactStatus = fmt.Sprint(i)
}

func (e *EventObject) GetImpactResolutionStatus() string {
	return e.ResolutionStatus
}

// SetValueResolutionStatus устанавливает STRING значение для поля ResolutionStatus
func (e *EventObject) SetValueResolutionStatus(v string) {
	e.ResolutionStatus = v
}

// SetAnyResolutionStatus устанавливает ЛЮБОЕ значение для поля ResolutionStatus
func (e *EventObject) SetAnyResolutionStatus(i interface{}) {
	e.ResolutionStatus = fmt.Sprint(i)
}

func (e *EventObject) GetStatus() string {
	return e.Status
}

// SetValueStatus устанавливает STRING значение для поля Status
func (e *EventObject) SetValueStatus(v string) {
	e.Status = v
}

// SetAnyStatus устанавливает ЛЮБОЕ значение для поля Status
func (e *EventObject) SetAnyStatus(i interface{}) {
	e.Status = fmt.Sprint(i)
}

func (e *EventObject) GetSummary() string {
	return e.Summary
}

// SetValueSummary устанавливает STRING значение для поля Summary
func (e *EventObject) SetValueSummary(v string) {
	e.Summary = v
}

// SetAnySummary устанавливает ЛЮБОЕ значение для поля Summary
func (e *EventObject) SetAnySummary(i interface{}) {
	e.Summary = fmt.Sprint(i)
}

func (e *EventObject) GetOwner() string {
	return e.Owner
}

// SetValueOwner устанавливает STRING значение для поля Owner
func (e *EventObject) SetValueOwner(v string) {
	e.Owner = v
}

// SetAnyOwner устанавливает ЛЮБОЕ значение для поля Owner
func (e *EventObject) SetAnyOwner(i interface{}) {
	e.Owner = fmt.Sprint(i)
}

func (e *EventObject) GetTags() []string {
	return e.Tags
}

// SetValueTags устанавливает STRING значение для поля Tags
func (e *EventObject) SetValueTags(v string) {
	e.Tags = append(e.Tags, v)
}

// SetAnyTags устанавливает ЛЮБОЕ значение для поля Tags
func (e *EventObject) SetAnyTags(i interface{}) {
	e.Tags = append(e.Tags, fmt.Sprint(i))
}

func (e *EventObject) GetCustomFields() map[string]CustomerFields {
	return e.CustomFields
}

// SetValueCustomFields устанавливает STRING значение для поля CustomFields
func (e *EventObject) SetValueCustomFields(v map[string]CustomerFields) {
	e.CustomFields = v
}

func (eo EventObject) ToStringBeautiful(num int) string {
	var str string

	ws := supportingfunctions.GetWhitespace(num)

	str += fmt.Sprintf("%s_id: '%s'\n", ws, eo.UnderliningId)
	str += fmt.Sprintf("%sid: '%s'\n", ws, eo.Id)
	str += fmt.Sprintf("%screatedBy: '%s'\n", ws, eo.CreatedBy)
	str += fmt.Sprintf("%supdatedBy: '%s'\n", ws, eo.UpdatedBy)
	str += fmt.Sprintf("%screatedAt: '%d'\n", ws, eo.CreatedAt)
	str += fmt.Sprintf("%supdatedAt: '%d'\n", ws, eo.UpdatedAt)
	str += fmt.Sprintf("%s_type: '%s'\n", ws, eo.UnderliningType)
	str += fmt.Sprintf("%scaseId: '%d'\n", ws, eo.CaseId)
	str += fmt.Sprintf("%stitle: '%s'\n", ws, eo.Title)
	str += fmt.Sprintf("%sdescription: '%s'\n", ws, eo.Description)
	str += fmt.Sprintf("%sseverity: '%d'\n", ws, eo.Severity)
	str += fmt.Sprintf("%sstartDate: '%d'\n", ws, eo.StartDate)
	str += fmt.Sprintf("%sendDate: '%d'\n", ws, eo.EndDate)
	str += fmt.Sprintf("%simpactStatus: '%s'\n", ws, eo.ImpactStatus)
	str += fmt.Sprintf("%sresolutionStatus: '%s'\n", ws, eo.ResolutionStatus)
	str += fmt.Sprintf("%stags: \n%s", ws, func(l []string) string {
		var str string
		ws := supportingfunctions.GetWhitespace(num + 1)

		for k, v := range l {
			str += fmt.Sprintf("%s%d. '%s'\n", ws, k+1, v)
		}
		return str
	}(eo.Tags))
	str += fmt.Sprintf("%sflag: '%v'\n", ws, eo.Flag)
	str += fmt.Sprintf("%stlp: '%d'\n", ws, eo.Tlp)
	str += fmt.Sprintf("%spap: '%d'\n", ws, eo.Pap)
	str += fmt.Sprintf("%sstatus: '%s'\n", ws, eo.Status)
	str += fmt.Sprintf("%ssummary: '%s'\n", ws, eo.Summary)
	str += fmt.Sprintf("%sowner: '%s'\n", ws, eo.Owner)
	str += fmt.Sprintf("%scustomFields: \n%s", ws, CustomFieldsToStringBeautiful(eo.CustomFields, num))
	/*str += fmt.Sprintf("%spermissions: \n%s", ws, func(l []string) string {
		var str string
		ws := supportingfunctions.GetWhitespace(num + 1)

		for k, v := range l {
			str += fmt.Sprintf("%s%d. '%s'\n", ws, k+1, v)
		}
		return str
	}(eo.Permissions))*/

	return str
}

func CustomFieldsToStringBeautiful(l map[string]CustomerFields, num int) string {
	var str string
	ws := supportingfunctions.GetWhitespace(num + 2)

	for k, v := range l {
		str += fmt.Sprintf("%s%s:\n", supportingfunctions.GetWhitespace(num+1), k)

		nameOne, dataOne, nameTwo, dataTwo := v.Get()
		str += fmt.Sprintf("%s%s: %d\n", ws, nameOne, dataOne)
		str += fmt.Sprintf("%s%s: %s\n", ws, nameTwo, dataTwo)
	}
	return str
}
