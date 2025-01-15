package datamodels

func (smisp ServersMispFormat) GetServersMisp() ServersMispFormat {
	return smisp
}

func (smisp *ServersMispFormat) SetValueNameServersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		smisp.Name = data

		isSuccess = true
	}

	return isSuccess
}

func (smisp *ServersMispFormat) SetValueUrlServersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		smisp.Url = data

		isSuccess = true
	}

	return isSuccess
}

func (smisp *ServersMispFormat) SetValueAuthkeyServersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		smisp.Authkey = data

		isSuccess = true
	}

	return isSuccess
}

func (smisp *ServersMispFormat) SetValueOrgIdServersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		smisp.OrgId = data

		isSuccess = true
	}

	return isSuccess
}

func (smisp *ServersMispFormat) SetValueLastpulledidServersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		smisp.Lastpulledid = data

		isSuccess = true
	}

	return isSuccess
}

func (smisp *ServersMispFormat) SetValueLastpushedidServersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		smisp.Lastpushedid = data

		isSuccess = true
	}

	return isSuccess
}

func (smisp *ServersMispFormat) SetValueOrganizationServersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		smisp.Organization = data

		isSuccess = true
	}

	return isSuccess
}

func (smisp *ServersMispFormat) SetValueRemoteOrgIdServersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		smisp.RemoteOrgId = data

		isSuccess = true
	}

	return isSuccess
}

func (smisp *ServersMispFormat) SetValuePullRulesServersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		smisp.PullRules = data

		isSuccess = true
	}

	return isSuccess
}

func (smisp *ServersMispFormat) SetValuePushRulesServersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		smisp.PushRules = data

		isSuccess = true
	}

	return isSuccess
}

func (smisp *ServersMispFormat) SetValueCertFileServersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		smisp.CertFile = data

		isSuccess = true
	}

	return isSuccess
}

func (smisp *ServersMispFormat) SetValueClientCertFileServersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		smisp.ClientCertFile = data

		isSuccess = true
	}

	return isSuccess
}

func (smisp *ServersMispFormat) SetValuePriorityServersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(string); ok {
		smisp.Priority = data

		isSuccess = true
	}

	return isSuccess
}

func (smisp *ServersMispFormat) SetValuePushServersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(bool); ok {
		smisp.Push = data

		isSuccess = true
	}

	return isSuccess
}

func (smisp *ServersMispFormat) SetValuePullServersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(bool); ok {
		smisp.Pull = data

		isSuccess = true
	}

	return isSuccess
}

func (smisp *ServersMispFormat) SetValuePushSightingsServersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(bool); ok {
		smisp.PushSightings = data

		isSuccess = true
	}

	return isSuccess
}

func (smisp *ServersMispFormat) SetValuePushGalaxyClustersServersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(bool); ok {
		smisp.PushGalaxyClusters = data

		isSuccess = true
	}

	return isSuccess
}

func (smisp *ServersMispFormat) SetValuePullGalaxyClustersServersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(bool); ok {
		smisp.PullGalaxyClusters = data

		isSuccess = true
	}

	return isSuccess
}

func (smisp *ServersMispFormat) SetValuePublishWithoutEmailServersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(bool); ok {
		smisp.PublishWithoutEmail = data

		isSuccess = true
	}

	return isSuccess
}

func (smisp *ServersMispFormat) SetValueUnpublishEventServersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(bool); ok {
		smisp.UnpublishEvent = data

		isSuccess = true
	}

	return isSuccess
}

func (smisp *ServersMispFormat) SetValueSelfSignedServersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(bool); ok {
		smisp.SelfSigned = data

		isSuccess = true
	}

	return isSuccess
}

func (smisp *ServersMispFormat) SetValueInternalServersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(bool); ok {
		smisp.Internal = data

		isSuccess = true
	}

	return isSuccess
}

func (smisp *ServersMispFormat) SetValueSkipProxyServersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(bool); ok {
		smisp.SkipProxy = data

		isSuccess = true
	}

	return isSuccess
}

func (smisp *ServersMispFormat) SetValueCachingEnabledServersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(bool); ok {
		smisp.CachingEnabled = data

		isSuccess = true
	}

	return isSuccess
}

func (smisp *ServersMispFormat) SetValueCacheTimestampServersMisp(v interface{}) bool {
	var isSuccess bool

	if data, ok := v.(bool); ok {
		smisp.CacheTimestamp = data

		isSuccess = true
	}

	return isSuccess
}
