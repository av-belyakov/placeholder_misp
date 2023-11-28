package confighandler

func (conf *ConfigApp) GetCommonApp() *CommonAppConfig {
	return &conf.CommonAppConfig
}

func (conf *ConfigApp) GetListLogs() []LogSet {
	return conf.LogList
}

func (conf *ConfigApp) GetListOrganization() []Organization {
	return conf.Organizations
}

func (conf *ConfigApp) GetAppRedis() *AppConfigRedis {
	return &conf.AppConfigRedis
}

func (conf *ConfigApp) GetAppNATS() *AppConfigNATS {
	return &conf.AppConfigNATS
}

func (conf *ConfigApp) GetAppMISP() *AppConfigMISP {
	return &conf.AppConfigMISP
}

func (conf *ConfigApp) GetAppTheHive() *AppConfigTheHive {
	return &conf.AppConfigTheHive
}

func (conf *ConfigApp) GetAppES() *AppConfigElasticSearch {
	return &conf.AppConfigElasticSearch
}

func (conf *ConfigApp) Clean() {
	conf = &ConfigApp{}
}
