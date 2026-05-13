package natsapi

type ModuleNatsHandler interface {
	GetDataReceptionChannel() <-chan OutputSettings
	SendingDataInput(InputSettings)
	SendingDataOutput(OutputSettings)
}
