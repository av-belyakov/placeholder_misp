package zabbixinteractions

type HandlerZabbixConnection struct {
	NetHost    string
	ZabbixHost string
	ZabbixKey  string
}

type PatternZabbix struct {
	Request string       `json:"request"`
	Data    []DataZabbix `json:"data"`
}

type DataZabbix struct {
	Host  string `json:"host"`
	Key   string `json:"key"`
	Value string `json:"value"`
}
