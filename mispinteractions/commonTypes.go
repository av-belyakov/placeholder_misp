package mispinteractions

// ModuleMISP инициализированный модуль
// chanInputMISP - канал для отправки данных в модуль
// chanOutputMISP - канал для отправки данных из модуля
type ModuleMISP struct {
	chanInputMISP  chan SettingsChanInputMISP
	chanOutputMISP chan SettingChanOutputMISP
}

type SettingsChanInputMISP struct {
	ObjectId  string
	UserEmail string
	MajorData map[string]interface{}
}

type SettingChanOutputMISP struct {
	Command, EventId string
}

func (mmisp ModuleMISP) GetDataReceptionChannel() <-chan SettingChanOutputMISP {
	return mmisp.chanOutputMISP
}

func (mmisp ModuleMISP) SendingDataInput(data SettingsChanInputMISP) {
	mmisp.chanInputMISP <- data
}

func (mmisp ModuleMISP) SendingDataOutput(data SettingChanOutputMISP) {
	mmisp.chanOutputMISP <- data
}
