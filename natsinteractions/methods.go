package natsinteractions

// GetDataReceptionChannel возвращает канал приема данных из модуля
func (mnats ModuleNATS) GetDataReceptionChannel() <-chan OutputSettings {
	return mnats.chanOutput
}

// SendingDataInput отправка данных в модуль
func (mnats ModuleNATS) SendingDataInput(data InputSettings) {
	mnats.chanInput <- data
}

// SendingDataOutput отправка данных из модуля
func (mnats ModuleNATS) SendingDataOutput(data OutputSettings) {
	mnats.chanOutput <- data
}
