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
	PathDirectory string
	MaxFileSize   int
	WritingStdout bool
	WritingFile   bool
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
	Prefix   string
	Index    string
	Name     string
	Authtype string
	User     string
	Passwd   string
	Host     string
	Port     int
}

type AppConfigNKCKI struct {
	Host string
	Port int
}

type RulesProcMSGMISP struct {
	Directory, File string
}
