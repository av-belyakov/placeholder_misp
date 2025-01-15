package main

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/av-belyakov/placeholder_misp/cmd/commoninterfaces"
	"github.com/av-belyakov/placeholder_misp/internal/confighandler"
	"github.com/av-belyakov/placeholder_misp/internal/datamodels"
	"github.com/av-belyakov/placeholder_misp/internal/supportingfunctions"
	"github.com/av-belyakov/placeholder_misp/memorytemporarystorage"
	"github.com/av-belyakov/placeholder_misp/zabbixinteractions"
	"github.com/av-belyakov/simplelogger"
)

func getLoggerSettings(cls []confighandler.LogSet) []simplelogger.Options {
	loggerConf := make([]simplelogger.Options, 0, len(cls))

	for _, v := range cls {
		loggerConf = append(loggerConf, simplelogger.Options{
			WritingToStdout: v.WritingStdout,
			WritingToFile:   v.WritingFile,
			WritingToDB:     v.WritingDB,
			MsgTypeName:     v.MsgTypeName,
			PathDirectory:   v.PathDirectory,
			MaxFileSize:     v.MaxFileSize,
		})
	}

	return loggerConf
}

// counterHandler обработчик счетчиков
func counterHandler(
	//channelZabbix chan<- zabbixinteractions.MessageSettings,
	channelZabbix chan<- commoninterfaces.Messager,
	storageApp *memorytemporarystorage.CommonStorageTemporary,
	simpleLogger *simplelogger.SimpleLoggerSettings,
	counting <-chan datamodels.DataCounterSettings) {
	for data := range counting {
		d, h, m, s := supportingfunctions.GetDifference(storageApp.GetStartTimeDataCounter(), time.Now())
		patternTime := fmt.Sprintf("со старта приложения: дней %d, часов %d, минут %d, секунд %d", d, h, m, s)
		var msg string

		switch data.DataType {
		case "update accepted events":
			storageApp.SetAcceptedEventsDataCounter(data.Count)
			msg = fmt.Sprintf("принято: %d, %s", storageApp.GetAcceptedEventsDataCounter(), patternTime)
		case "update processed events":
			storageApp.SetProcessedEventsDataCounter(data.Count)
			msg = fmt.Sprintf("обработано: %d, %s", storageApp.GetProcessedEventsDataCounter(), patternTime)
		case "update events meet rules":
			storageApp.SetEventsMeetRulesDataCounter(data.Count)
			msg = fmt.Sprintf("соответствует правилам: %d, %s", storageApp.GetEventsMeetRulesDataCounter(), patternTime)
		}

		_ = simpleLogger.Write("debug", msg)

		msg := NewMessageLogging()
		msg.SetType("error")
		msg.SetMessage(fmt.Sprintf("%s: %s", msg.GetType(), msg.GetMessage()))

		channelZabbix <- msg

		channelZabbix <- zabbixinteractions.MessageSettings{
			EventType: "info",
			Message:   msg,
		}
	}
}

// interactionZabbix осуществляет взаимодействие с Zabbix
func interactionZabbix(
	ctx context.Context,
	confApp *confighandler.ConfigApp,
	simpleLogger *simplelogger.SimpleLoggerSettings,
	channelZabbix <-chan zabbixinteractions.MessageSettings) error {

	connTimeout := time.Duration(7 * time.Second)
	hz, err := zabbixinteractions.NewZabbixConnection(
		ctx,
		zabbixinteractions.SettingsZabbixConnection{
			Port:              confApp.Zabbix.NetworkPort,
			Host:              confApp.Zabbix.NetworkHost,
			NetProto:          "tcp",
			ZabbixHost:        confApp.Zabbix.ZabbixHost,
			ConnectionTimeout: &connTimeout,
		})
	if err != nil {
		return err
	}

	et := make([]zabbixinteractions.EventType, len(confApp.Zabbix.EventTypes))
	for _, v := range confApp.Zabbix.EventTypes {
		et = append(et, zabbixinteractions.EventType{
			IsTransmit: v.IsTransmit,
			EventType:  v.EventType,
			ZabbixKey:  v.ZabbixKey,
			Handshake: zabbixinteractions.Handshake{
				TimeInterval: v.Handshake.TimeInterval,
				Message:      v.Handshake.Message,
			},
		})
	}

	if err = hz.Handler(et, channelZabbix); err != nil {
		return err
	}

	go func() {
		for err := range hz.GetChanErr() {
			_, f, l, _ := runtime.Caller(0)
			_ = simpleLogger.Write("error", fmt.Sprintf("zabbix module: '%s' %s:%d", err.Error(), f, l-1))
		}
	}()

	return nil
}
