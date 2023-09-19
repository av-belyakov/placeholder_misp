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
	LogList []LogSet
}

type Logs struct {
	Logging []LogSet
}

type LogSet struct {
	MsgTypeName   string
	WritingFile   bool
	PathDirectory string
	WritingStdout bool
	MaxFileSize   int
}

type AppConfigNATS struct {
	Host string
	Port int
}

type AppConfigMISP struct {
	Host string
	Auth string
}

type AppConfigRedis struct {
	Host string
	Port int
}

type AppConfigElasticSearch struct {
	Host     string
	Port     int
	Prefix   string
	Index    string
	Name     string
	Authtype string
	User     string
	Passwd   string
}

type AppConfigNKCKI struct {
	Host string
	Port int
}

type RulesProcMSGMISP struct {
	Directory, File string
}
