package datamodels

import (
	"fmt"

	"github.com/google/uuid"
)

func NewEventMisp() *EventsMispFormat {
	return &EventsMispFormat{
		Timestamp:         "0",
		Published:         getPublished(),
		PublishTimestamp:  "0",
		SightingTimestamp: "0",
		Uuid:              uuid.New().String(),
		Analysis:          getAnalysis(),
		Distribution:      getDistributionEvent(),
		ThreatLevelId:     getThreatLevelId(),
		SharingGroupId:    getSharingGroupId(),
	}
}

func (emisp EventsMispFormat) GetEventsMisp() EventsMispFormat {
	return emisp
}

func (emisp *EventsMispFormat) CleanEventsMispFormat() {
	emisp.OrgId = ""
	emisp.OrgcId = ""
	emisp.Info = ""
	emisp.Uuid = uuid.New().String()
	emisp.Date = ""
	emisp.AttributeCount = ""
	emisp.ExtendsUuid = ""
	emisp.EventCreatorEmail = ""
	emisp.Published = getPublished()
	emisp.ProposalEmailLock = false
	emisp.Locked = false
	emisp.DisableCorrelation = false
	emisp.Analysis = getAnalysis()
	emisp.Distribution = getDistributionEvent()
	emisp.Timestamp = "0"
	emisp.ThreatLevelId = getThreatLevelId()
	emisp.PublishTimestamp = "0"
	emisp.SightingTimestamp = "0"
	emisp.SharingGroupId = getSharingGroupId()
}

// Comparison выполняет сравнение двух объектов типа EventsMispFormat
func (emisp *EventsMispFormat) Comparison(newEvents *EventsMispFormat) bool {
	if emisp.Analysis != newEvents.Analysis {
		return false
	}

	if emisp.Analysis != newEvents.Analysis {
		return false
	}

	if emisp.AttributeCount != newEvents.AttributeCount {
		return false
	}

	if emisp.OrgId != newEvents.OrgId {
		return false
	}

	if emisp.OrgcId != newEvents.OrgcId {
		return false
	}

	if emisp.Distribution != newEvents.Distribution {
		return false
	}

	if emisp.Info != newEvents.Info {
		return false
	}

	if emisp.Uuid != newEvents.Uuid {
		return false
	}

	if emisp.Date != newEvents.Date {
		return false
	}

	if emisp.SharingGroupId != newEvents.SharingGroupId {
		return false
	}

	if emisp.ThreatLevelId != newEvents.ThreatLevelId {
		return false
	}

	if emisp.ExtendsUuid != newEvents.ExtendsUuid {
		return false
	}

	if emisp.EventCreatorEmail != newEvents.EventCreatorEmail {
		return false
	}

	if emisp.Published != newEvents.Published {
		return false
	}

	if emisp.ProposalEmailLock != newEvents.ProposalEmailLock {
		return false
	}

	if emisp.Locked != newEvents.Locked {
		return false
	}

	if emisp.DisableCorrelation != newEvents.DisableCorrelation {
		return false
	}

	// думаю время сравнивать не стоит, потому что большая вероятность получить идентичный
	//во всех параметрах объект у которого будет отличатся только время, что в данном случае не очень важно
	//Timestamp
	//PublishTimestamp
	//SightingTimestamp

	return true
}

// GetOrgId возвращает значение OrgId
func (emisp *EventsMispFormat) GetOrgId() string {
	return emisp.OrgId
}

// SetAnyOrgId устанавливает значение для OrgId
func (emisp *EventsMispFormat) SetAnyOrgId(v interface{}, num int) {
	if data, ok := v.(string); ok {
		emisp.OrgId = data
	}
}

// GetOrgcId возвращает значение OrgcId
func (emisp *EventsMispFormat) GetOrgcId() string {
	return emisp.OrgcId
}

// SetAnyOrgcId устанавливает значение для OrgcId
func (emisp *EventsMispFormat) SetAnyOrgcId(v interface{}, num int) {
	if data, ok := v.(string); ok {
		emisp.OrgcId = data
	}
}

// SetOrgcId устанавливает значение для OrgcId
func (emisp *EventsMispFormat) SetOrgcId(v string) {
	emisp.OrgcId = v
}

// GetDistribution возвращает значение Distribution
func (emisp *EventsMispFormat) GetDistribution() string {
	return emisp.Distribution
}

// SetAnyDistribution устанавливает значение для Distribution
func (emisp *EventsMispFormat) SetAnyDistribution(v interface{}, num int) {
	if data, ok := v.(string); ok {
		emisp.Distribution = data
	}
}

// SetDistribution устанавливает значение для Distribution
func (emisp *EventsMispFormat) SetDistribution(v string) {
	emisp.Distribution = v
}

// GetInfo возвращает значение Info
func (emisp *EventsMispFormat) GetInfo() string {
	return emisp.Info
}

// SetAnyInfo устанавливает значение для Info
func (emisp *EventsMispFormat) SetAnyInfo(v interface{}, num int) {
	if data, ok := v.(string); ok {
		emisp.Info = data
	}
}

// SetInfo устанавливает значение для Info
func (emisp *EventsMispFormat) SetInfo(v string) {
	emisp.Info = v
}

// GetUUID возвращает значение UUID
func (emisp *EventsMispFormat) GetUUID() string {
	return emisp.Uuid
}

// SetAnyUUID устанавливает значение для UUID
func (emisp *EventsMispFormat) SetAnyUUID(v interface{}, num int) {
	if data, ok := v.(string); ok {
		emisp.Uuid = data
	}
}

// SetUUID устанавливает значение для UUID
func (emisp *EventsMispFormat) SetUUID(v string) {
	emisp.Uuid = v
}

// GetDate возвращает значение Date
func (emisp *EventsMispFormat) GetDate() string {
	return emisp.Date
}

// SetAnyDate устанавливает значение для Date
func (emisp *EventsMispFormat) SetAnyDate(v interface{}, num int) {
	if data, ok := v.(string); ok {
		emisp.Date = data
	}
}

// SetDate устанавливает значение для Date
func (emisp *EventsMispFormat) SetDate(v string) {
	emisp.Date = v
}

// GetAnalysis возвращает значение Analysis
func (emisp *EventsMispFormat) GetAnalysis() string {
	return emisp.Analysis
}

// SetAnyAnalysis устанавливает значение для Analysis
func (emisp *EventsMispFormat) SetAnyAnalysis(v interface{}, num int) {
	if data, ok := v.(string); ok {
		emisp.Analysis = data
	}
}

// SetAnalysis устанавливает значение для Analysis
func (emisp *EventsMispFormat) SetAnalysis(v string) {
	emisp.Analysis = v
}

// GetAttributeCount возвращает значение AttributeCount
func (emisp *EventsMispFormat) GetAttributeCount() string {
	return emisp.AttributeCount
}

// SetAnyAttributeCount устанавливает значение для AttributeCount
func (emisp *EventsMispFormat) SetAnyAttributeCount(v interface{}, num int) {
	if data, ok := v.(string); ok {
		emisp.AttributeCount = data
	}
}

// SetAttributeCount устанавливает значение для AttributeCount
func (emisp *EventsMispFormat) SetAttributeCount(v string) {
	emisp.AttributeCount = v
}

// GetTimestamp возвращает значение Timestamp
func (emisp *EventsMispFormat) GetTimestamp() string {
	return emisp.Timestamp
}

// SetAnyTimestamp устанавливает значение для Timestamp
func (emisp *EventsMispFormat) SetAnyTimestamp(v interface{}, num int) {
	if data, ok := v.(float64); ok {
		//emisp.Timestamp = fmt.Sprintf("%13.f", data)
		emisp.Timestamp = fmt.Sprintf("%13.f", data)[:10]
	}
}

// SetTimestamp устанавливает значение для Timestamp
func (emisp *EventsMispFormat) SetTimestamp(v string) {
	emisp.Timestamp = v
}

// GetSharingGroupId возвращает значение SharingGroupId
func (emisp *EventsMispFormat) GetSharingGroupId() string {
	return emisp.SharingGroupId
}

// SetAnySharingGroupId устанавливает значение для SharingGroupId
func (emisp *EventsMispFormat) SetAnySharingGroupId(v interface{}, num int) {
	if data, ok := v.(string); ok {
		emisp.SharingGroupId = data
	}
}

// SetSharingGroupId устанавливает значение для SharingGroupId
func (emisp *EventsMispFormat) SetSharingGroupId(v string) {
	emisp.SharingGroupId = v
}

// GetThreatLevelId возвращает значение ThreatLevelId
func (emisp *EventsMispFormat) GetThreatLevelId() string {
	return emisp.ThreatLevelId
}

// SetAnyThreatLevelId устанавливает значение для ThreatLevelId
func (emisp *EventsMispFormat) SetAnyThreatLevelId(v interface{}, num int) {
	if data, ok := v.(string); ok {
		emisp.ThreatLevelId = data
	}

	if data, ok := v.(float64); ok {
		emisp.ThreatLevelId = fmt.Sprint(data)
	}
}

// SetThreatLevelId устанавливает значение для ThreatLevelId
func (emisp *EventsMispFormat) SetThreatLevelId(v string) {
	emisp.ThreatLevelId = v
}

// GetPublishTimestamp возвращает значение PublishTimestamp
func (emisp *EventsMispFormat) GetPublishTimestamp() string {
	return emisp.PublishTimestamp
}

// SetAnyPublishTimestamp устанавливает значение для PublishTimestamp
func (emisp *EventsMispFormat) SetAnyPublishTimestamp(v interface{}, num int) {
	if data, ok := v.(float64); ok {
		//emisp.PublishTimestamp = fmt.Sprintf("%13.f", data)
		emisp.PublishTimestamp = fmt.Sprintf("%13.f", data)[:10]
	}
}

// SetPublishTimestamp устанавливает значение для PublishTimestamp
func (emisp *EventsMispFormat) SetPublishTimestamp(v string) {
	emisp.PublishTimestamp = v
}

// GetSightingTimestamp возвращает значение SightingTimestamp
func (emisp *EventsMispFormat) GetSightingTimestamp() string {
	return emisp.SightingTimestamp
}

// SetAnySightingTimestamp устанавливает значение для SightingTimestamp
func (emisp *EventsMispFormat) SetAnySightingTimestamp(v interface{}, num int) {
	if data, ok := v.(float64); ok {
		//emisp.SightingTimestamp = fmt.Sprintf("%13.f", data)
		emisp.SightingTimestamp = fmt.Sprintf("%13.f", data)[:10]
	}
}

// SetSightingTimestamp устанавливает значение для SightingTimestamp
func (emisp *EventsMispFormat) SetSightingTimestamp(v string) {
	emisp.SightingTimestamp = v
}

// GetExtendsUuid возвращает значение ExtendsUuid
func (emisp *EventsMispFormat) GetExtendsUuid() string {
	return emisp.ExtendsUuid
}

// SetAnyExtendsUUID устанавливает значение для ExtendsUUID
func (emisp *EventsMispFormat) SetAnyExtendsUUID(v interface{}, num int) {
	if data, ok := v.(string); ok {
		emisp.ExtendsUuid = data
	}
}

// SetExtendsUUID устанавливает значение для ExtendsUUID
func (emisp *EventsMispFormat) SetExtendsUUID(v string) {
	emisp.ExtendsUuid = v
}

// GetEventCreatorEmail возвращает значение EventCreatorEmail
func (emisp *EventsMispFormat) GetEventCreatorEmail() string {
	return emisp.EventCreatorEmail
}

// SetAnyEventCreatorEmail устанавливает значение для EventCreatorEmail
func (emisp *EventsMispFormat) SetAnyEventCreatorEmail(v interface{}, num int) {
	if data, ok := v.(string); ok {
		emisp.EventCreatorEmail = data
	}
}

// SetEventCreatorEmail устанавливает значение для EventCreatorEmail
func (emisp *EventsMispFormat) SetEventCreatorEmail(v string) {
	emisp.EventCreatorEmail = v
}

// GetPublished возвращает значение Published
func (emisp *EventsMispFormat) GetPublished() bool {
	return emisp.Published
}

// SetAnyPublished устанавливает значение для Published
func (emisp *EventsMispFormat) SetAnyPublished(v interface{}, num int) {
	if data, ok := v.(bool); ok {
		emisp.Published = data
	}
}

// SetPublished устанавливает значение для Published
func (emisp *EventsMispFormat) SetPublished(v bool) {
	emisp.Published = v
}

// GetProposalEmailLock возвращает значение ProposalEmailLock
func (emisp *EventsMispFormat) GetProposalEmailLock() bool {
	return emisp.ProposalEmailLock
}

// SetAnyProposalEmailLock устанавливает значение для ProposalEmailLock
func (emisp *EventsMispFormat) SetAnyProposalEmailLock(v interface{}, num int) {
	if data, ok := v.(bool); ok {
		emisp.ProposalEmailLock = data
	}
}

// SetProposalEmailLock устанавливает значение для ProposalEmailLock
func (emisp *EventsMispFormat) SetProposalEmailLock(v bool) {
	emisp.ProposalEmailLock = v
}

// GetLocked возвращает значение Locked
func (emisp *EventsMispFormat) GetLocked() bool {
	return emisp.Locked
}

// SetAnyLocked устанавливает значение для Locked
func (emisp *EventsMispFormat) SetAnyLocked(v interface{}, num int) {
	if data, ok := v.(bool); ok {
		emisp.Locked = data
	}
}

// SetLocked устанавливает значение для Locked
func (emisp *EventsMispFormat) SetLocked(v bool) {
	emisp.Locked = v
}

// GetDisableCorrelation возвращает значение DisableCorrelation
func (emisp *EventsMispFormat) GetDisableCorrelation() bool {
	return emisp.DisableCorrelation
}

// SetAnyDisableCorrelation устанавливает значение для DisableCorrelation
func (emisp *EventsMispFormat) SetAnyDisableCorrelation(v interface{}, num int) {
	if data, ok := v.(bool); ok {
		emisp.DisableCorrelation = data
	}
}

// SetDisableCorrelation устанавливает значение для DisableCorrelation
func (emisp *EventsMispFormat) SetDisableCorrelation(v bool) {
	emisp.DisableCorrelation = v
}

func getAnalysis() string {
	return "2"
}

func getDistributionEvent() string {
	return "1"
}

func getThreatLevelId() string {
	return "4"
}

func getSharingGroupId() string {
	return "1"
}

func getPublished() bool {
	return false
}

/*func getTagTLP(tlp int) datamodels.TagsMispFormat {
	tag := datamodels.TagsMispFormat{Name: "tlp:red", Colour: "#cc0033"}

	switch tlp {
	case 0:
		tag.Name = "tlp:white"
		tag.Colour = "#ffffff"
	case 1:
		tag.Name = "tlp:green"
		tag.Colour = "#339900"
	case 2:
		tag.Name = "tlp:amber"
		tag.Colour = "#ffc000"
	}

	return tag
}*/
