package zabbixinteractions

import (
	"context"
	"time"
)

// SettingsZabbixConnection настройки Zabbix соединения
// Port - сетевой порт
// Host - ip адрес или доменное имя
// NetProto - сетевой протокол (по умолчанию используется tcp)
// ZabbixHost - имя Zabbix хоста
// ConnectionTimeout - время ожидания подключения (по умолчанию используется 5 сек)
type SettingsZabbixConnection struct {
	Port              int
	Host              string
	NetProto          string
	ZabbixHost        string
	ConnectionTimeout *time.Duration
}

type HandlerZabbixConnection struct {
	ctx         context.Context
	port        int
	host        string
	netProto    string
	zabbixHost  string
	connTimeout *time.Duration
	chanErr     chan error
}

type ZabbixOptions struct {
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

type MessageSettings struct {
	Message, EventType string
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
