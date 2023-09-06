package datamodels

import "fmt"

func NewEventMisp() EventsMispFormat {
	return EventsMispFormat{
		Analysis:          getAnalysis(),
		Distribution:      getDistributionEvent(),
		Timestamp:         "0",
		ThreatLevelId:     getThreatLevelId(),
		PublishTimestamp:  "0",
		SightingTimestamp: "0",
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
	emisp.Uuid = ""
	emisp.Date = ""
	emisp.AttributeCount = ""
	emisp.ExtendsUuid = ""
	emisp.EventCreatorEmail = ""
	emisp.Published = false
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

func (emisp *EventsMispFormat) SetValueOrgIdEventsMisp(v interface{}, isNew bool) {
	if data, ok := v.(string); ok {
		emisp.OrgId = data
	}
}

func (emisp *EventsMispFormat) SetValueOrgcIdEventsMisp(v interface{}, isNew bool) {
	if data, ok := v.(string); ok {
		emisp.OrgcId = data
	}
}

func (emisp *EventsMispFormat) SetValueDistributionEventsMisp(v interface{}, isNew bool) {
	if data, ok := v.(string); ok {
		emisp.Distribution = data
	}
}

func (emisp *EventsMispFormat) SetValueInfoEventsMisp(v interface{}, isNew bool) {
	if data, ok := v.(string); ok {
		emisp.Info = data
	}
}

func (emisp *EventsMispFormat) SetValueUuidEventsMisp(v interface{}, isNew bool) {
	if data, ok := v.(string); ok {
		emisp.Uuid = data
	}
}

func (emisp *EventsMispFormat) SetValueDateEventsMisp(v interface{}, isNew bool) {
	if data, ok := v.(string); ok {
		emisp.Date = data
	}
}

func (emisp *EventsMispFormat) SetValueAnalysisEventsMisp(v interface{}, isNew bool) {
	if data, ok := v.(string); ok {
		emisp.Analysis = data
	}
}

func (emisp *EventsMispFormat) SetValueAttributeCountEventsMisp(v interface{}, isNew bool) {
	if data, ok := v.(string); ok {
		emisp.AttributeCount = data
	}
}

func (emisp *EventsMispFormat) SetValueTimestampEventsMisp(v interface{}, isNew bool) {
	if data, ok := v.(float64); ok {
		//emisp.Timestamp = fmt.Sprintf("%13.f", data)
		emisp.Timestamp = fmt.Sprintf("%13.f", data)[:10]
	}
}

func (emisp *EventsMispFormat) SetValueSharingGroupIdEventsMisp(v interface{}, isNew bool) {
	if data, ok := v.(string); ok {
		emisp.SharingGroupId = data
	}
}

func (emisp *EventsMispFormat) SetValueThreatLevelIdEventsMisp(v interface{}, isNew bool) {
	if data, ok := v.(string); ok {
		emisp.ThreatLevelId = data
	}

	if data, ok := v.(float64); ok {
		emisp.ThreatLevelId = fmt.Sprint(data)
	}
}

func (emisp *EventsMispFormat) SetValuePublishTimestampEventsMisp(v interface{}, isNew bool) {
	if data, ok := v.(float64); ok {
		//emisp.PublishTimestamp = fmt.Sprintf("%13.f", data)
		emisp.PublishTimestamp = fmt.Sprintf("%13.f", data)[:10]
	}
}

func (emisp *EventsMispFormat) SetValueSightingTimestampEventsMisp(v interface{}, isNew bool) {
	if data, ok := v.(float64); ok {
		//emisp.SightingTimestamp = fmt.Sprintf("%13.f", data)
		emisp.SightingTimestamp = fmt.Sprintf("%13.f", data)[:10]
	}
}

func (emisp *EventsMispFormat) SetValueExtendsUuidEventsMisp(v interface{}, isNew bool) {
	if data, ok := v.(string); ok {
		emisp.ExtendsUuid = data
	}
}

func (emisp *EventsMispFormat) SetValueEventCreatorEmailEventsMisp(v interface{}, isNew bool) {
	if data, ok := v.(string); ok {
		emisp.EventCreatorEmail = data
	}
}

func (emisp *EventsMispFormat) SetValuePublishedEventsMisp(v interface{}, isNew bool) {
	if data, ok := v.(bool); ok {
		emisp.Published = data
	}
}

func (emisp *EventsMispFormat) SetValueProposalEmailLockEventsMisp(v interface{}, isNew bool) {
	if data, ok := v.(bool); ok {
		emisp.ProposalEmailLock = data
	}
}

func (emisp *EventsMispFormat) SetValueLockedEventsMisp(v interface{}, isNew bool) {
	if data, ok := v.(bool); ok {
		emisp.Locked = data
	}
}

func (emisp *EventsMispFormat) SetValueDisableCorrelationEventsMisp(v interface{}, isNew bool) {
	if data, ok := v.(bool); ok {
		emisp.DisableCorrelation = data
	}
}

func getAnalysis() string {
	return "2"
}

func getDistributionEvent() string {
	return "3"
}

func getThreatLevelId() string {
	return "4"
}

func getSharingGroupId() string {
	return "1"
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
