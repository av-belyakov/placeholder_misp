package commoninterfaces

import (
	"github.com/av-belyakov/placeholder_misp/cmd/mispapi"
	"github.com/av-belyakov/placeholder_misp/cmd/natsapi"
)

type ModuleMispHandler interface {
	GetDataReceptionChannel() <-chan mispapi.OutputSetting
	SendingDataOutput(mispapi.OutputSetting)
	SendingDataInput(mispapi.InputSettings)
}

type ModuleNatsHandler interface {
	GetDataReceptionChannel() <-chan natsapi.OutputSettings
	SendingDataInput(natsapi.InputSettings)
	SendingDataOutput(natsapi.OutputSettings)
}

//************** логирование ***************

type Logger interface {
	GetChan() <-chan Messager
	Send(msgType, msgData string)
}

type Messager interface {
	GetType() string
	SetType(v string)
	GetMessage() string
	SetMessage(v string)
}

type WriterLoggingData interface {
	Write(typeLogFile, str string) bool
}
