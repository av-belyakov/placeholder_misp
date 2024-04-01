package confighandler

type ConfigApp struct {
	CommonAppConfig
	AppConfigNATS
	AppConfigMISP
	AppConfigRedis
	AppConfigTheHive
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
	NetworkPort int         `yaml:"networkPort"`
	NetworkHost string      `yaml:"networkHost"`
	ZabbixHost  string      `yaml:"zabbixHost"`
	EventTypes  []EventType `yaml:"eventType"`
}

type ZabbixHostOptions struct {
	ZabbixHost string      `yaml:"zabbixHost"`
	EventTypes []EventType `yaml:"eventType"`
}

type EventType struct {
	IsTransmit bool      `yaml:"isTransmit"`
	EventType  string    `yaml:"eventType"`
	ZabbixKey  string    `yaml:"zabbixKey"`
	Handshake  Handshake `yaml:"handshake"`
}

type Handshake struct {
	TimeInterval int    `yaml:"timeInterval"`
	Message      string `yaml:"message"`
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

type AppConfigTheHive struct {
	Send bool `yaml:"send"`
}

type RulesProcMSGMISP struct {
	Directory, File string
}
