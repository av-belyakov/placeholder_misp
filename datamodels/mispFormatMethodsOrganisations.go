package datamodels

func (omisp OrganisationsMispFormat) GetOrganisationsMisp() OrganisationsMispFormat {
	return omisp
}

func (omisp *OrganisationsMispFormat) SetValueNameOrganisationsMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		omisp.Name = data

		isSuccess = true
	}

	return isSuccess
}

func (omisp *OrganisationsMispFormat) SetValueDateCreatedOrganisationsMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		omisp.DateCreated = data

		isSuccess = true
	}

	return isSuccess
}

func (omisp *OrganisationsMispFormat) SetValueDateModifiedOrganisationsMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		omisp.DateModified = data

		isSuccess = true
	}

	return isSuccess
}

func (omisp *OrganisationsMispFormat) SetValueDescriptionOrganisationsMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		omisp.Description = data

		isSuccess = true
	}

	return isSuccess
}

func (omisp *OrganisationsMispFormat) SetValueTypeOrganisationsMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		omisp.Type = data

		isSuccess = true
	}

	return isSuccess
}

func (omisp *OrganisationsMispFormat) SetValueNationalityOrganisationsMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		omisp.Nationality = data

		isSuccess = true
	}

	return isSuccess
}

func (omisp *OrganisationsMispFormat) SetValueSectorOrganisationsMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		omisp.Sector = data

		isSuccess = true
	}

	return isSuccess
}

func (omisp *OrganisationsMispFormat) SetValueCreatedByOrganisationsMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		omisp.CreatedBy = data

		isSuccess = true
	}

	return isSuccess
}

func (omisp *OrganisationsMispFormat) SetValueUuidOrganisationsMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		omisp.Uuid = data

		isSuccess = true
	}

	return isSuccess
}

func (omisp *OrganisationsMispFormat) SetValueContactsOrganisationsMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		omisp.Contacts = data

		isSuccess = true
	}

	return isSuccess
}

func (omisp *OrganisationsMispFormat) SetValueLandingpageOrganisationsMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		omisp.Landingpage = data

		isSuccess = true
	}

	return isSuccess
}

func (omisp *OrganisationsMispFormat) SetValueUserCountOrganisationsMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		omisp.UserCount = data

		isSuccess = true
	}

	return isSuccess
}

func (omisp *OrganisationsMispFormat) SetValueCreatedByEmailOrganisationsMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		omisp.CreatedByEmail = data

		isSuccess = true
	}

	return isSuccess
}

func (omisp *OrganisationsMispFormat) SetValueRestrictedToDomainOrganisationsMisp(v interface{}) bool {
	var isSuccess bool

	switch v := v.(type) {
	case string:
		omisp.RestrictedToDomain = append(omisp.RestrictedToDomain, v)

		isSuccess = true
	case []string:
		omisp.RestrictedToDomain = append(omisp.RestrictedToDomain, v...)

		isSuccess = true
	}

	return isSuccess
}

func (omisp *OrganisationsMispFormat) SetValueLocalOrganisationsMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(bool); ok {
		omisp.Local = data

		isSuccess = true
	}

	return isSuccess
}
