package wrappers

import (
	"context"
	"fmt"
	"time"

	"github.com/av-belyakov/placeholder_misp/commoninterfaces"
	"github.com/av-belyakov/placeholder_misp/internal/supportingfunctions"
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
		writerLoggingData.Write("error", supportingfunctions.CustomError(fmt.Errorf("zabbix module: %w", err)).Error())

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
		writerLoggingData.Write("error", supportingfunctions.CustomError(fmt.Errorf("zabbix module: %w", err)).Error())

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
				writerLoggingData.Write("error", supportingfunctions.CustomError(fmt.Errorf("zabbix module: %W", errMsg)).Error())

			}
		}
	}()
}
