package confighandler

type ConfigApp struct {
	CommonAppConfig
	AppConfigNATS
	AppConfigMISP
	AppConfigSqlite3
	AppConfigTheHive
	RulesProcMSGMISP
}

type CommonAppConfig struct {
	LogList       []*LogSet
	Organizations []Organization
	Zabbix        ZabbixOptions
}

type Logs struct {
	Logging []*LogSet
}

type LogSet struct {
	MsgTypeName   string `validate:"oneof=error info warning" yaml:"msgTypeName"`
	PathDirectory string `validate:"required" yaml:"pathDirectory"`
	MaxFileSize   int    `validate:"min=1000" yaml:"maxFileSize"`
	WritingStdout bool   `validate:"required" yaml:"writingStdout"`
	WritingFile   bool   `validate:"required" yaml:"writingFile"`
	WritingDB     bool   `validate:"required" yaml:"writingDB"`
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
	EventTypes  []EventType `yaml:"eventType"`
	NetworkHost string      `validate:"required" yaml:"networkHost"`
	ZabbixHost  string      `validate:"required" yaml:"zabbixHost"`
	NetworkPort int         `validate:"gt=0,lte=65535" yaml:"networkPort"`
}

type EventType struct {
	EventType  string    `validate:"required" yaml:"eventType"`
	ZabbixKey  string    `validate:"required" yaml:"zabbixKey"`
	Handshake  Handshake `yaml:"handshake"`
	IsTransmit bool      `yaml:"isTransmit"`
}

type Handshake struct {
	Message      string `validate:"required" yaml:"message"`
	TimeInterval int    `yaml:"timeInterval"`
}

type AppConfigNATS struct {
	Subscriptions SubscriptionsNATS `yaml:"subscriptions"`
	Host          string            `validate:"required" yaml:"host"`
	Port          int               `validate:"gt=0,lte=65535" yaml:"port"`
	CacheTTL      int               `validate:"gt=10,lte=86400" yaml:"cache_ttl"`
}

type SubscriptionsNATS struct {
	ListenerCase  string `validate:"required" yaml:"listener_case"`
	SenderCommand string `validate:"required" yaml:"sender_command"`
}

type AppConfigMISP struct {
	Host string `validate:"required" yaml:"host"`
	Auth string `validate:"required" yaml:"auth"`
}

type AppConfigSqlite3 struct {
	PathFileDb string `validate:"required" yaml:"path_file_db"`
}

type AppConfigTheHive struct {
	Send bool `yaml:"send"`
}

type RulesProcMSGMISP struct {
	Directory, File string
}
