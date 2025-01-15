package wrappers

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/av-belyakov/placeholder_misp/commoninterfaces"
	zabbixapicommunicator "github.com/av-belyakov/zabbixapicommunicator/cmd"
)

// WrappersZabbixInteraction обертка для взаимодействия с модулем zabbixapi
func WrappersZabbixInteraction(
	ctx context.Context,
	settings WrappersZabbixInteractionSettings,
	writerLoggingData commoninterfaces.WriterLoggingData,
	channelZabbix <-chan commoninterfaces.Messager) {

	connTimeout := time.Duration(7 * time.Second)
	zc, err := zabbixapicommunicator.New(zabbixapicommunicator.SettingsZabbixConnection{
		Port:              settings.NetworkPort,
		Host:              settings.NetworkHost,
		NetProto:          "tcp",
		ZabbixHost:        settings.ZabbixHost,
		ConnectionTimeout: &connTimeout,
	})
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		writerLoggingData.Write("error", fmt.Sprintf("zabbix module: '%s' %s:%d", err.Error(), f, l-1))

		return
	}

	et := make([]zabbixapicommunicator.EventType, len(settings.EventTypes))
	for _, v := range settings.EventTypes {
		et = append(et, zabbixapicommunicator.EventType{
			IsTransmit: v.IsTransmit,
			EventType:  v.EventType,
			ZabbixKey:  v.ZabbixKey,
			Handshake:  zabbixapicommunicator.Handshake(v.Handshake),
		})
	}

	recipient := make(chan zabbixapicommunicator.Messager)
	if err = zc.Start(ctx, et, recipient); err != nil {
		_, f, l, _ := runtime.Caller(0)
		writerLoggingData.Write("error", fmt.Sprintf("zabbix module: '%s' %s:%d", err.Error(), f, l-1))

		return
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case msg := <-channelZabbix:
				newMessageSettings := &zabbixapicommunicator.MessageSettings{}
				newMessageSettings.SetType(msg.GetType())
				newMessageSettings.SetMessage(msg.GetMessage())

				recipient <- newMessageSettings

			case errMsg := <-zc.GetChanErr():
				writerLoggingData.Write("error", fmt.Sprintf("zabbix module: '%s'", errMsg.Error()))

			}
		}
	}()
}
