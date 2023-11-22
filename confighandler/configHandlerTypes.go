package confighandler

type ConfigApp struct {
	CommonAppConfig
	AppConfigNATS
	AppConfigMISP
	AppConfigRedis
	AppConfigElasticSearch
	AppConfigNKCKI
	RulesProcMSGMISP
}

type CommonAppConfig struct {
	LogList       []LogSet
	Organizations []Organization
	Zabbix        ZabbixOptions
}

type Logs struct {
	Logging []LogSet
}

type LogSet struct {
	WritingStdout bool   `yaml:"writingStdout"`
	WritingFile   bool   `yaml:"writingFile"`
	MaxFileSize   int    `yaml:"maxFileSize"`
	MsgTypeName   string `yaml:"msgTypeName"`
	PathDirectory string `yaml:"pathDirectory"`
}

type Orgs struct {
	Organizations []Organization
}

type Organization struct {
	OrgName    string `yaml:"orgName"`
	SourceName string `yaml:"sourceName"`
}

type ZabbixSet struct {
	Zabbix ZabbixOptions
}

type ZabbixOptions struct {
	IsTransmit   bool   `yaml:"isTransmit"`
	TimeInterval int    `yaml:"timeInterval"`
	NetworkPort  int    `yaml:"networkPort"`
	NetworkHost  string `yaml:"networkHost"`
	ZabbixHost   string `yaml:"zabbixHost"`
	ZabbixKey    string `yaml:"zabbixKey"`
	Handshake    string `yaml:"handshake"`
}

type AppConfigNATS struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type AppConfigMISP struct {
	Host string `yaml:"host"`
	Auth string `yaml:"auth"`
}

type AppConfigRedis struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type AppConfigElasticSearch struct {
	Prefix   string `yaml:"prefix"`
	Index    string `yaml:"index"`
	Name     string `yaml:"name"`
	Authtype string `yaml:"authtype"`
	User     string `yaml:"user"`
	Passwd   string `yaml:"passwd"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
}

type AppConfigNKCKI struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type RulesProcMSGMISP struct {
	Directory, File string
}
