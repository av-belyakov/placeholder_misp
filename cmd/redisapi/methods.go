package redisapi

// GetDataReceptionChannel получить канал для отправки данных ИЗ модуля
func (mmisp ModuleRedis) GetReceptionChannel() <-chan SettingsOutput {
	return mmisp.chOutput
}

// SendingDataInput отправка данных В модуль
func (mmisp ModuleRedis) SendDataInput(data SettingsInput) {
	mmisp.chInput <- data
}

// SendingDataOutput отправка данных ИЗ модуля
func (mmisp ModuleRedis) SendDataOutput(data SettingsOutput) {
	mmisp.chOutput <- data
}
