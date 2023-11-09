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
