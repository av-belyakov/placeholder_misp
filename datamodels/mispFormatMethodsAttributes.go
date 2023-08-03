package datamodels

func (amisp AttributesMispFormat) GetAttributesMisp() AttributesMispFormat {
	return amisp
}

func (amisp *AttributesMispFormat) SetValueEventIdAttributesMisp(v interface{}) bool {
	var isSuccess bool

	if str, ok := v.(string); ok {
		amisp.EventId = str

		isSuccess = true
	}

	return isSuccess
}

func (amisp *AttributesMispFormat) SetValueObjectIdAttributesMisp(v interface{}) bool {
	var isSuccess bool

	if str, ok := v.(string); ok {
		amisp.ObjectId = str

		isSuccess = true
	}

	return isSuccess
}

func (amisp *AttributesMispFormat) SetValueObjectRelationAttributesMisp(v interface{}) bool {
	var isSuccess bool

	if str, ok := v.(string); ok {
		amisp.ObjectRelation = str

		isSuccess = true
	}

	return isSuccess
}

func (amisp *AttributesMispFormat) SetValueCategoryAttributesMisp(v interface{}) bool {
	var isSuccess bool

	if str, ok := v.(string); ok {
		amisp.Category = str

		isSuccess = true
	}

	return isSuccess
}

func (amisp *AttributesMispFormat) SetValueTypeAttributesMisp(v interface{}) bool {
	var isSuccess bool

	if str, ok := v.(string); ok {
		amisp.Type = str

		isSuccess = true
	}

	return isSuccess
}

func (amisp *AttributesMispFormat) SetValueValueAttributesMisp(v interface{}) bool {
	var isSuccess bool

	if str, ok := v.(string); ok {
		amisp.Value = str

		isSuccess = true
	}

	return isSuccess
}

func (amisp *AttributesMispFormat) SetValueUuidAttributesMisp(v interface{}) bool {
	var isSuccess bool

	if str, ok := v.(string); ok {
		amisp.Uuid = str

		isSuccess = true
	}

	return isSuccess
}

func (amisp *AttributesMispFormat) SetValueTimestampAttributesMisp(v interface{}) bool {
	var isSuccess bool

	if str, ok := v.(string); ok {
		amisp.Timestamp = str

		isSuccess = true
	}

	return isSuccess
}

func (amisp *AttributesMispFormat) SetValueDistributionAttributesMisp(v interface{}) bool {
	var isSuccess bool

	if str, ok := v.(string); ok {
		amisp.Distribution = str

		isSuccess = true
	}

	return isSuccess
}

func (amisp *AttributesMispFormat) SetValueSharingGroupIdAttributesMisp(v interface{}) bool {
	var isSuccess bool

	if str, ok := v.(string); ok {
		amisp.SharingGroupId = str

		isSuccess = true
	}

	return isSuccess
}

func (amisp *AttributesMispFormat) SetValueCommentAttributesMisp(v interface{}) bool {
	var isSuccess bool

	if str, ok := v.(string); ok {
		amisp.Comment = str

		isSuccess = true
	}

	return isSuccess
}

func (amisp *AttributesMispFormat) SetValueFirstSeenAttributesMisp(v interface{}) bool {
	var isSuccess bool

	if str, ok := v.(string); ok {
		amisp.FirstSeen = str

		isSuccess = true
	}

	return isSuccess
}

func (amisp *AttributesMispFormat) SetValueLastSeenAttributesMisp(v interface{}) bool {
	var isSuccess bool

	if str, ok := v.(string); ok {
		amisp.LastSeen = str

		isSuccess = true
	}

	return isSuccess
}

func (amisp *AttributesMispFormat) SetValueToIdsAttributesMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(bool); ok {
		amisp.ToIds = data

		isSuccess = true
	}

	return isSuccess
}

func (amisp *AttributesMispFormat) SetValueDeletedAttributesMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(bool); ok {
		amisp.Deleted = data

		isSuccess = true
	}

	return isSuccess
}

func (amisp *AttributesMispFormat) SetValueDisableCorrelationAttributesMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(bool); ok {
		amisp.DisableCorrelation = data

		isSuccess = true
	}

	return isSuccess
}
