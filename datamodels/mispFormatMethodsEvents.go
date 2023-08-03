package datamodels

func (emisp EventsMispFormat) GetEventsMisp() EventsMispFormat {
	return emisp
}

func (emisp *EventsMispFormat) SetValueOrgIdEventsMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		emisp.OrgId = data

		isSuccess = true
	}

	return isSuccess
}

func (emisp *EventsMispFormat) SetValueOrgcIdEventsMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		emisp.OrgcId = data

		isSuccess = true
	}

	return isSuccess
}

func (emisp *EventsMispFormat) SetValueDistributionEventsMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		emisp.Distribution = data

		isSuccess = true
	}

	return isSuccess
}

func (emisp *EventsMispFormat) SetValueInfoEventsMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		emisp.Info = data

		isSuccess = true
	}

	return isSuccess
}

func (emisp *EventsMispFormat) SetValueUuidEventsMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		emisp.Uuid = data

		isSuccess = true
	}

	return isSuccess
}

func (emisp *EventsMispFormat) SetValueDateEventsMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		emisp.Date = data

		isSuccess = true
	}

	return isSuccess
}

func (emisp *EventsMispFormat) SetValueAnalysisEventsMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		emisp.Analysis = data

		isSuccess = true
	}

	return isSuccess
}

func (emisp *EventsMispFormat) SetValueAttributeCountEventsMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		emisp.AttributeCount = data

		isSuccess = true
	}

	return isSuccess
}

func (emisp *EventsMispFormat) SetValueTimestampEventsMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		emisp.Timestamp = data

		isSuccess = true
	}

	return isSuccess
}

func (emisp *EventsMispFormat) SetValueSharingGroupIdEventsMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		emisp.SharingGroupId = data

		isSuccess = true
	}

	return isSuccess
}

func (emisp *EventsMispFormat) SetValueThreatLevelIdEventsMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		emisp.ThreatLevelId = data

		isSuccess = true
	}

	return isSuccess
}

func (emisp *EventsMispFormat) SetValuePublishTimestampEventsMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		emisp.PublishTimestamp = data

		isSuccess = true
	}

	return isSuccess
}

func (emisp *EventsMispFormat) SetValueSightingTimestampEventsMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		emisp.SightingTimestamp = data

		isSuccess = true
	}

	return isSuccess
}

func (emisp *EventsMispFormat) SetValueExtendsUuidEventsMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		emisp.ExtendsUuid = data

		isSuccess = true
	}

	return isSuccess
}

func (emisp *EventsMispFormat) SetValueEventCreatorEmailEventsMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		emisp.EventCreatorEmail = data

		isSuccess = true
	}

	return isSuccess
}

func (emisp *EventsMispFormat) SetValuePublishedEventsMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(bool); ok {
		emisp.Published = data

		isSuccess = true
	}

	return isSuccess
}

func (emisp *EventsMispFormat) SetValueProposalEmailLockEventsMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(bool); ok {
		emisp.ProposalEmailLock = data

		isSuccess = true
	}

	return isSuccess
}

func (emisp *EventsMispFormat) SetValueLockedEventsMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(bool); ok {
		emisp.Locked = data

		isSuccess = true
	}

	return isSuccess
}

func (emisp *EventsMispFormat) SetValueDisableCorrelationEventsMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(bool); ok {
		emisp.DisableCorrelation = data

		isSuccess = true
	}

	return isSuccess
}
