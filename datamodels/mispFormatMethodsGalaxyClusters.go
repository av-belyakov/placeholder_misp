package datamodels

func (gcmisp GalaxyClustersMispFormat) GetGalaxyClustersMisp() GalaxyClustersMispFormat {
	return gcmisp
}

func (gcmisp *GalaxyClustersMispFormat) SetValueIdGalaxyClustersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		gcmisp.Id = data

		isSuccess = true
	}

	return isSuccess
}

func (gcmisp *GalaxyClustersMispFormat) SetValueUuidGalaxyClustersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		gcmisp.Uuid = data

		isSuccess = true
	}

	return isSuccess
}

func (gcmisp *GalaxyClustersMispFormat) SetValueCollectionUuidGalaxyClustersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		gcmisp.CollectionUuid = data

		isSuccess = true
	}

	return isSuccess
}

func (gcmisp *GalaxyClustersMispFormat) SetValueTypeGalaxyClustersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		gcmisp.Type = data

		isSuccess = true
	}

	return isSuccess
}

func (gcmisp *GalaxyClustersMispFormat) SetValueValueGalaxyClustersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		gcmisp.Value = data

		isSuccess = true
	}

	return isSuccess
}

func (gcmisp *GalaxyClustersMispFormat) SetValueTagNameGalaxyClustersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		gcmisp.TagName = data

		isSuccess = true
	}

	return isSuccess
}

func (gcmisp *GalaxyClustersMispFormat) SetValueDescriptionGalaxyClustersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		gcmisp.Description = data

		isSuccess = true
	}

	return isSuccess
}

func (gcmisp *GalaxyClustersMispFormat) SetValueGalaxyIdGalaxyClustersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		gcmisp.GalaxyId = data

		isSuccess = true
	}

	return isSuccess
}

func (gcmisp *GalaxyClustersMispFormat) SetValueSourceGalaxyClustersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		gcmisp.Source = data

		isSuccess = true
	}

	return isSuccess
}

func (gcmisp *GalaxyClustersMispFormat) SetValueAuthorsGalaxyClustersMisp(v interface{}) bool {
	var isSuccess bool

	switch v := v.(type) {
	case string:
		gcmisp.Authors = append(gcmisp.Authors, v)

		isSuccess = true
	case []string:
		gcmisp.Authors = append(gcmisp.Authors, v...)

		isSuccess = true
	}

	return isSuccess
}

func (gcmisp *GalaxyClustersMispFormat) SetValueVersionGalaxyClustersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		gcmisp.Version = data

		isSuccess = true
	}

	return isSuccess
}

func (gcmisp *GalaxyClustersMispFormat) SetValueDistributionGalaxyClustersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		gcmisp.Description = data

		isSuccess = true
	}

	return isSuccess
}

func (gcmisp *GalaxyClustersMispFormat) SetValueSharingGroupIdGalaxyClustersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		gcmisp.SharingGroupId = data

		isSuccess = true
	}

	return isSuccess
}

func (gcmisp *GalaxyClustersMispFormat) SetValueOrgIdGalaxyClustersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		gcmisp.OrgId = data

		isSuccess = true
	}

	return isSuccess
}

func (gcmisp *GalaxyClustersMispFormat) SetValueOrgcIdGalaxyClustersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		gcmisp.OrgcId = data

		isSuccess = true
	}

	return isSuccess
}

func (gcmisp *GalaxyClustersMispFormat) SetValueExtendsUuidGalaxyClustersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		gcmisp.ExtendsUuid = data

		isSuccess = true
	}

	return isSuccess
}

func (gcmisp *GalaxyClustersMispFormat) SetValueExtendsVersionGalaxyClustersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		gcmisp.ExtendsVersion = data

		isSuccess = true
	}

	return isSuccess
}

func (gcmisp *GalaxyClustersMispFormat) SetValueDefaultGalaxyClustersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(bool); ok {
		gcmisp.Default = data

		isSuccess = true
	}

	return isSuccess
}

func (gcmisp *GalaxyClustersMispFormat) SetValueLockedGalaxyClustersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(bool); ok {
		gcmisp.Locked = data

		isSuccess = true
	}

	return isSuccess
}

func (gcmisp *GalaxyClustersMispFormat) SetValuePublishedGalaxyClustersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(bool); ok {
		gcmisp.Published = data

		isSuccess = true
	}

	return isSuccess
}

func (gcmisp *GalaxyClustersMispFormat) SetValueDeletedGalaxyClustersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(bool); ok {
		gcmisp.Deleted = data

		isSuccess = true
	}

	return isSuccess
}

func (gcmisp *GalaxyClustersMispFormat) SetValueGalaxyElementGalaxyClustersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.([]GalaxyElementMispFormat); ok {
		gcmisp.GalaxyElement = append(gcmisp.GalaxyElement, data...)

		isSuccess = true
	}

	return isSuccess
}
