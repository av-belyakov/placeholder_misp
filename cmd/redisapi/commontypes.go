package redisapi

// ModuleRedis инициализированный модуль
// chanInputRedis - канал для отправки данных в модуль
// chanOutputRedis - канал для отправки данных из модуля
type ModuleRedis struct {
	chanInputRedis  chan SettingsChanInputRedis
	chanOutputRedis chan SettingChanOutputRedis
}

type SettingsChanInputRedis struct {
	Command string
	Data    string
	RawData []byte
}

type SettingChanOutputRedis struct {
	CommandResult string
	Result        interface{}
}

func (mmisp ModuleRedis) GetDataReceptionChannel() <-chan SettingChanOutputRedis {
	return mmisp.chanOutputRedis
}

func (mmisp ModuleRedis) SendingDataInput(data SettingsChanInputRedis) {
	mmisp.chanInputRedis <- data
}

func (mmisp ModuleRedis) SendingDataOutput(data SettingChanOutputRedis) {
	mmisp.chanOutputRedis <- data
}
