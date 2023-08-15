package mispinteractions

// ModuleMISP инициализированный модуль
// chanInputMISP - канал для отправки данных в модуль
// chanOutputMISP - канал для отправки данных из модуля
type ModuleMISP struct {
	chanInputMISP  chan map[string]interface{}
	chanOutputMISP chan interface{}
}

func (mmisp ModuleMISP) GetDataReceptionChannel() <-chan interface{} {
	return mmisp.chanOutputMISP
}

func (mmisp ModuleMISP) GettingData() interface{} {
	return <-mmisp.chanOutputMISP
}

func (mmisp ModuleMISP) SendingDataInputMisp(data map[string]interface{}) {
	mmisp.chanInputMISP <- data
}

func (mmisp ModuleMISP) SendingDataOutputMisp(data interface{}) {
	mmisp.chanOutputMISP <- data
}
