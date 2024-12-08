package datamodels

import (
	"fmt"

	"github.com/google/uuid"
)

func NewEventMisp() EventsMispFormat {
	return EventsMispFormat{
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

// SetOrgId устанавливает значение для OrgId
func (emisp *EventsMispFormat) SetOrgId(v interface{}, num int) {
	if data, ok := v.(string); ok {
		emisp.OrgId = data
	}
}

// GetOrgId возвращает значение OrgId
func (emisp *EventsMispFormat) GetOrgId() string {
	return emisp.OrgId
}

// SetOrgcId устанавливает значение для OrgcId
func (emisp *EventsMispFormat) SetOrgcId(v interface{}, num int) {
	if data, ok := v.(string); ok {
		emisp.OrgcId = data
	}
}

// GetOrgcId возвращает значение OrgcId
func (emisp *EventsMispFormat) GetOrgcId() string {
	return emisp.OrgcId
}

// SetDistribution устанавливает значение для Distribution
func (emisp *EventsMispFormat) SetDistribution(v interface{}, num int) {
	if data, ok := v.(string); ok {
		emisp.Distribution = data
	}
}

// GetDistribution возвращает значение Distribution
func (emisp *EventsMispFormat) GetDistribution() string {
	return emisp.Distribution
}

// SetInfo устанавливает значение для Info
func (emisp *EventsMispFormat) SetInfo(v interface{}, num int) {
	if data, ok := v.(string); ok {
		emisp.Info = data
	}
}

// GetInfo возвращает значение Info
func (emisp *EventsMispFormat) GetInfo() string {
	return emisp.Info
}

// SetUUID устанавливает значение для UUID
func (emisp *EventsMispFormat) SetUUID(v interface{}, num int) {
	if data, ok := v.(string); ok {
		emisp.Uuid = data
	}
}

// GetUUID возвращает значение UUID
func (emisp *EventsMispFormat) GetUUID() string {
	return emisp.Uuid
}

// SetDate устанавливает значение для Date
func (emisp *EventsMispFormat) SetDate(v interface{}, num int) {
	if data, ok := v.(string); ok {
		emisp.Date = data
	}
}

// GetDate возвращает значение Date
func (emisp *EventsMispFormat) GetDate() string {
	return emisp.Date
}

// SetAnalysis устанавливает значение для Analysis
func (emisp *EventsMispFormat) SetAnalysis(v interface{}, num int) {
	if data, ok := v.(string); ok {
		emisp.Analysis = data
	}
}

// GetAnalysis возвращает значение Analysis
func (emisp *EventsMispFormat) GetAnalysis() string {
	return emisp.Analysis
}

// SetAttributeCount устанавливает значение для AttributeCount
func (emisp *EventsMispFormat) SetAttributeCount(v interface{}, num int) {
	if data, ok := v.(string); ok {
		emisp.AttributeCount = data
	}
}

// GetAttributeCount возвращает значение AttributeCount
func (emisp *EventsMispFormat) GetAttributeCount() string {
	return emisp.AttributeCount
}

// SetTimestamp устанавливает значение для Timestamp
func (emisp *EventsMispFormat) SetTimestamp(v interface{}, num int) {
	if data, ok := v.(float64); ok {
		//emisp.Timestamp = fmt.Sprintf("%13.f", data)
		emisp.Timestamp = fmt.Sprintf("%13.f", data)[:10]
	}
}

// GetTimestamp возвращает значение Timestamp
func (emisp *EventsMispFormat) GetTimestamp() string {
	return emisp.Timestamp
}

// SetSharingGroupId устанавливает значение для SharingGroupId
func (emisp *EventsMispFormat) SetSharingGroupId(v interface{}, num int) {
	if data, ok := v.(string); ok {
		emisp.SharingGroupId = data
	}
}

// GetSharingGroupId возвращает значение SharingGroupId
func (emisp *EventsMispFormat) GetSharingGroupId() string {
	return emisp.SharingGroupId
}

// SetThreatLevelId устанавливает значение для ThreatLevelId
func (emisp *EventsMispFormat) SetThreatLevelId(v interface{}, num int) {
	if data, ok := v.(string); ok {
		emisp.ThreatLevelId = data
	}

	if data, ok := v.(float64); ok {
		emisp.ThreatLevelId = fmt.Sprint(data)
	}
}

// GetThreatLevelId возвращает значение ThreatLevelId
func (emisp *EventsMispFormat) GetThreatLevelId() string {
	return emisp.ThreatLevelId
}

// SetPublishTimestamp устанавливает значение для PublishTimestamp
func (emisp *EventsMispFormat) SetPublishTimestamp(v interface{}, num int) {
	if data, ok := v.(float64); ok {
		//emisp.PublishTimestamp = fmt.Sprintf("%13.f", data)
		emisp.PublishTimestamp = fmt.Sprintf("%13.f", data)[:10]
	}
}

// GetPublishTimestamp возвращает значение PublishTimestamp
func (emisp *EventsMispFormat) GetPublishTimestamp() string {
	return emisp.PublishTimestamp
}

// SetSightingTimestamp устанавливает значение для SightingTimestamp
func (emisp *EventsMispFormat) SetSightingTimestamp(v interface{}, num int) {
	if data, ok := v.(float64); ok {
		//emisp.SightingTimestamp = fmt.Sprintf("%13.f", data)
		emisp.SightingTimestamp = fmt.Sprintf("%13.f", data)[:10]
	}
}

// GetSightingTimestamp возвращает значение SightingTimestamp
func (emisp *EventsMispFormat) GetSightingTimestamp() string {
	return emisp.SightingTimestamp
}

// SetExtendsUUID устанавливает значение для ExtendsUUID
func (emisp *EventsMispFormat) SetExtendsUUID(v interface{}, num int) {
	if data, ok := v.(string); ok {
		emisp.ExtendsUuid = data
	}
}

// GetExtendsUuid возвращает значение ExtendsUuid
func (emisp *EventsMispFormat) GetExtendsUuid() string {
	return emisp.ExtendsUuid
}

// SetEventCreatorEmail устанавливает значение для EventCreatorEmail
func (emisp *EventsMispFormat) SetEventCreatorEmail(v interface{}, num int) {
	if data, ok := v.(string); ok {
		emisp.EventCreatorEmail = data
	}
}

// GetEventCreatorEmail возвращает значение EventCreatorEmail
func (emisp *EventsMispFormat) GetEventCreatorEmail() string {
	return emisp.EventCreatorEmail
}

// SetPublished устанавливает значение для Published
func (emisp *EventsMispFormat) SetPublished(v interface{}, num int) {
	if data, ok := v.(bool); ok {
		emisp.Published = data
	}
}

// GetPublished возвращает значение Published
func (emisp *EventsMispFormat) GetPublished() bool {
	return emisp.Published
}

// SetProposalEmailLock устанавливает значение для ProposalEmailLock
func (emisp *EventsMispFormat) SetProposalEmailLock(v interface{}, num int) {
	if data, ok := v.(bool); ok {
		emisp.ProposalEmailLock = data
	}
}

// GetProposalEmailLock возвращает значение ProposalEmailLock
func (emisp *EventsMispFormat) GetProposalEmailLock() bool {
	return emisp.ProposalEmailLock
}

// SetLocked устанавливает значение для Locked
func (emisp *EventsMispFormat) SetLocked(v interface{}, num int) {
	if data, ok := v.(bool); ok {
		emisp.Locked = data
	}
}

// GetLocked возвращает значение Locked
func (emisp *EventsMispFormat) GetLocked() bool {
	return emisp.Locked
}

// SetDisableCorrelation устанавливает значение для DisableCorrelation
func (emisp *EventsMispFormat) SetDisableCorrelation(v interface{}, num int) {
	if data, ok := v.(bool); ok {
		emisp.DisableCorrelation = data
	}
}

// GetDisableCorrelation возвращает значение DisableCorrelation
func (emisp *EventsMispFormat) GetDisableCorrelation() bool {
	return emisp.DisableCorrelation
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
