package natsinteractions

// ModuleNATS инициализированный модуль
// ChanOutputMISP - канал для отправки полученных данных из модуля
// chanInputNATS - канал для принятия данных в модуль
type ModuleNATS struct {
	chanOutputNATS chan SettingsOutputChan
	chanInputNATS  chan SettingsInputChan
}

type SettingsOutputChan struct {
	MsgId string
	Data  []byte
}

type SettingsInputChan struct {
	Command, EventId, TaskId, RootId, CaseId string
}

func (mnats ModuleNATS) GetDataReceptionChannel() <-chan SettingsOutputChan {
	return mnats.chanOutputNATS
}

func (mnats ModuleNATS) SendingDataInput(data SettingsInputChan) {
	mnats.chanInputNATS <- data
}

func (mnats ModuleNATS) SendingDataOutput(data SettingsOutputChan) {
	mnats.chanOutputNATS <- data
}

type ResponseToCommand struct {
	StatusCode int         `json:"status_code"`
	ID         string      `json:"id"`
	Error      string      `json:"error"`
	Command    string      `json:"command"`
	Data       interface{} `json:"data"`
}
