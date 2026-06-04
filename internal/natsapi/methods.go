package natsapi

// GetChannelFromModule возвращает канал приема данных из модуля
func (api *ApiNatsModule) GetChannelFromModule() <-chan OutputSettings {
	return api.chanOutput
}

// GetChannelToModule возвращает канал приема данных в модуль
func (api *ApiNatsModule) GetChannelToModule() chan InputSettings {
	return api.chanInput
}

// SendingDataInput отправка данных в модуль
func (api *ApiNatsModule) SendingDataInput(data InputSettings) {
	api.chanInput <- data
}

// SendingDataOutput отправка данных из модуля
func (api *ApiNatsModule) SendingDataOutput(data OutputSettings) {
	api.chanOutput <- data
}
