package nkckiinteractions

import (
	"placeholder_misp/confighandler"
	"placeholder_misp/datamodels"
)

var nkcki ModuleNKCKI

// ModuleNKCKI инициализированный модуль
// chanInputNKCKI - канал для отправки данных в модуль
// chanOutputNKCKI - канал для принятия данных из модуля
// ChanLogging - канал для отправки логов
type ModuleNKCKI struct {
	chanInputNKCKI  chan interface{}
	chanOutputNKCKI chan interface{}
	ChanLogging     chan<- datamodels.MessageLoging
}

func init() {
	nkcki = ModuleNKCKI{
		chanInputNKCKI:  make(chan interface{}),
		chanOutputNKCKI: make(chan interface{}),
	}
}

func NewClientNKCKI(
	conf confighandler.AppConfigNKCKI,
	chanLog chan<- datamodels.MessageLoging) (*ModuleNKCKI, error) {

	nkcki.ChanLogging = chanLog

	/*
		Здесь нужно написать инициализацию подключения
	*/

	return &nkcki, nil
}

func (nkcki ModuleNKCKI) GetDataReceptionChannel() <-chan interface{} {
	/*
		для формирование правильного сообщения об ошибке
		if err != nil {
			_, f, l, _ := runtime.Caller(0)

			mmisp.ChanLogging <- datamodels.MessageLoging{
				MsgData: fmt.Sprintf("%s %s:%d", fmt.Sprint(err), f, l-2),
				MsgType: "error",
			}
		}
	*/

	return nkcki.chanOutputNKCKI
}

func (nkcki ModuleNKCKI) GettingData() interface{} {
	/*
		для формирование правильного сообщения об ошибке
		if err != nil {
			_, f, l, _ := runtime.Caller(0)

			mmisp.ChanLogging <- datamodels.MessageLoging{
				MsgData: fmt.Sprintf("%s %s:%d", fmt.Sprint(err), f, l-2),
				MsgType: "error",
			}
		}
	*/

	return <-nkcki.chanOutputNKCKI
}

func (nkcki ModuleNKCKI) SendingData(data interface{}) {
	/*
		для формирование правильного сообщения об ошибке
		if err != nil {
			_, f, l, _ := runtime.Caller(0)

			mmisp.ChanLogging <- datamodels.MessageLoging{
				MsgData: fmt.Sprintf("%s %s:%d", fmt.Sprint(err), f, l-2),
				MsgType: "error",
			}
		}
	*/

	nkcki.chanInputNKCKI <- data
}
