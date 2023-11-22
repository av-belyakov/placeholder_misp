package zabbixinteractions

// NewHandlerZabbixConnection создает обработчик соединения с Zabbix
// h - ip адрес или доменное имя сервера Zabbix и сет. порт, строка
// в формате "zabix.example.nets:10051"
// zh - приемник сообщений
// zk - ключ для приемника сообщений
func NewHandlerZabbixConnection(h, zh, zk string) *HandlerZabbixConnection {
	return &HandlerZabbixConnection{
		NetHost:    h,
		ZabbixHost: zh,
		ZabbixKey:  zk,
	}
}
