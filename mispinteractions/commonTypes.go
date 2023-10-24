package mispinteractions

// ModuleMISP инициализированный модуль
// chanInputMISP - канал для отправки данных в модуль
// chanOutputMISP - канал для отправки данных из модуля
/*type ModuleMISP struct {
	chanInputMISP  chan SettingsChanInputMISP
	chanOutputMISP chan SettingChanOutputMISP
}*/
type ModuleMISP struct {
	ChanInputMISP  chan SettingsChanInputMISP
	ChanOutputMISP chan SettingChanOutputMISP
}

type SettingsChanInputMISP struct {
	Command   string
	CaseId    float64
	EventId   string
	UserEmail string
	MajorData map[string]interface{}
}

type SettingChanOutputMISP struct {
	Command, CaseId, EventId string
}

func (mmisp ModuleMISP) GetDataReceptionChannel() <-chan SettingChanOutputMISP {
	return mmisp.ChanOutputMISP
}

func (mmisp ModuleMISP) SendingDataOutput(data SettingChanOutputMISP) {
	mmisp.ChanOutputMISP <- data
}

func (mmisp ModuleMISP) GetInputChannel() <-chan SettingsChanInputMISP {
	return mmisp.ChanInputMISP
}

func (mmisp ModuleMISP) SendingDataInput(data SettingsChanInputMISP) {
	mmisp.ChanInputMISP <- data
}

/*
func (mmisp ModuleMISP) GetDataReceptionChannel() <-chan SettingChanOutputMISP {
	return mmisp.chanOutputMISP
}

func (mmisp ModuleMISP) SendingDataInput(data SettingsChanInputMISP) {
	mmisp.chanInputMISP <- data
}

func (mmisp ModuleMISP) SendingDataOutput(data SettingChanOutputMISP) {
	mmisp.chanOutputMISP <- data
}
*/
