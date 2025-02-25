package datamodels

func (umisp UsersMispFormat) GetUsersMisp() UsersMispFormat {
	return umisp
}

func (umisp *UsersMispFormat) SetValueOrgIdUsersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		umisp.OrgId = data

		isSuccess = true
	}

	return isSuccess
}

func (umisp *UsersMispFormat) SetValueServerIdUsersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		umisp.ServerId = data

		isSuccess = true
	}

	return isSuccess
}

func (umisp *UsersMispFormat) SetValueEmailUsersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		umisp.Email = data

		isSuccess = true
	}

	return isSuccess
}

func (umisp *UsersMispFormat) SetValueAuthkeyUsersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		umisp.Authkey = data

		isSuccess = true
	}

	return isSuccess
}

func (umisp *UsersMispFormat) SetValueInvitedByUsersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		umisp.InvitedBy = data

		isSuccess = true
	}

	return isSuccess
}

func (umisp *UsersMispFormat) SetValueGpgkeyUsersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		umisp.Gpgkey = data

		isSuccess = true
	}

	return isSuccess
}

func (umisp *UsersMispFormat) SetValueCertifPublicUsersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		umisp.CertifPublic = data

		isSuccess = true
	}

	return isSuccess
}

func (umisp *UsersMispFormat) SetValueNidsSidUsersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		umisp.NidsSid = data

		isSuccess = true
	}

	return isSuccess
}

func (umisp *UsersMispFormat) SetValueNewsreadUsersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		umisp.Newsread = data

		isSuccess = true
	}

	return isSuccess
}

func (umisp *UsersMispFormat) SetValueRoleIdUsersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		umisp.RoleId = data

		isSuccess = true
	}

	return isSuccess
}

func (umisp *UsersMispFormat) SetValueChangePwUsersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		umisp.ChangePw = data

		isSuccess = true
	}

	return isSuccess
}

func (umisp *UsersMispFormat) SetValueExpirationUsersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		umisp.Expiration = data

		isSuccess = true
	}

	return isSuccess
}

func (umisp *UsersMispFormat) SetValueCurrentLoginUsersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		umisp.CurrentLogin = data

		isSuccess = true
	}

	return isSuccess
}

func (umisp *UsersMispFormat) SetValueLastLoginUsersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		umisp.LastLogin = data

		isSuccess = true
	}

	return isSuccess
}

func (umisp *UsersMispFormat) SetValueDateCreatedUsersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		umisp.DateCreated = data

		isSuccess = true
	}

	return isSuccess
}

func (umisp *UsersMispFormat) SetValueDateModifiedUsersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		umisp.DateModified = data

		isSuccess = true
	}

	return isSuccess
}

func (umisp *UsersMispFormat) SetValueAutoalertUsersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(bool); ok {
		umisp.Autoalert = data

		isSuccess = true
	}

	return isSuccess
}

func (umisp *UsersMispFormat) SetValueTermsacceptedUsersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(bool); ok {
		umisp.Termsaccepted = data

		isSuccess = true
	}

	return isSuccess
}

func (umisp *UsersMispFormat) SetValueContactalertUsersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(bool); ok {
		umisp.Contactalert = data

		isSuccess = true
	}

	return isSuccess
}

func (umisp *UsersMispFormat) SetValueDisabledUsersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(bool); ok {
		umisp.Disabled = data

		isSuccess = true
	}

	return isSuccess
}

func (umisp *UsersMispFormat) SetValueForceLogoutUsersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(bool); ok {
		umisp.ForceLogout = data

		isSuccess = true
	}

	return isSuccess
}
