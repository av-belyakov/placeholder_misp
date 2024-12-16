package interfaces

import (
	"placeholder_misp/mispinteractions"
	"placeholder_misp/natsinteractions"
)

type ModuleMispHandler interface {
	GetDataReceptionChannel() <-chan mispinteractions.OutputSetting
	SendingDataOutput(mispinteractions.OutputSetting)
	SendingDataInput(mispinteractions.InputSettings)
}

type ModuleNatsHandler interface {
	GetDataReceptionChannel() <-chan natsinteractions.OutputSettings
	SendingDataInput(natsinteractions.InputSettings)
	SendingDataOutput(natsinteractions.OutputSettings)
}
