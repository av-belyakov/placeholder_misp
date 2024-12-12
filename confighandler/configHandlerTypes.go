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
	WritingStdout bool   `validate:"required" yaml:"writingStdout"`
	WritingFile   bool   `validate:"required" yaml:"writingFile"`
	MaxFileSize   int    `validate:"min=1000" yaml:"maxFileSize"`
	MsgTypeName   string `validate:"oneof=error info warning" yaml:"msgTypeName"`
	PathDirectory string `validate:"required" yaml:"pathDirectory"`
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
	NetworkPort int         `validate:"gt=0,lte=65535" yaml:"networkPort"`
	NetworkHost string      `validate:"required" yaml:"networkHost"`
	ZabbixHost  string      `validate:"required" yaml:"zabbixHost"`
	EventTypes  []EventType `yaml:"eventType"`
}

type EventType struct {
	IsTransmit bool      `yaml:"isTransmit"`
	EventType  string    `validate:"required" yaml:"eventType"`
	ZabbixKey  string    `validate:"required" yaml:"zabbixKey"`
	Handshake  Handshake `yaml:"handshake"`
}

type Handshake struct {
	TimeInterval int    `yaml:"timeInterval"`
	Message      string `validate:"required" yaml:"message"`
}

type AppConfigNATS struct {
	Port          int               `validate:"gt=0,lte=65535" yaml:"port"`
	Host          string            `validate:"required" yaml:"host"`
	Subscriptions SubscriptionsNATS `yaml:"subscriptions"`
}

type SubscriptionsNATS struct {
	SenderCase      string `validate:"required" yaml:"sender_case"`
	ListenerCommand string `validate:"required" yaml:"listener_command"`
}

type AppConfigMISP struct {
	Host string `yaml:"host"`
	Auth string `yaml:"auth"`
}

type AppConfigRedis struct {
	Port int    `validate:"gt=0,lte=65535" yaml:"port"`
	Host string `yaml:"host"`
}

type AppConfigTheHive struct {
	Send bool `yaml:"send"`
}

type RulesProcMSGMISP struct {
	Directory, File string
}
