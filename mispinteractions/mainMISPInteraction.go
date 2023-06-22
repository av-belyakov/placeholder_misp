package mispinteractions

import (
	"context"
	"fmt"

	"placeholder_misp/confighandler"
	"placeholder_misp/datamodels"
)

var mmisp ModuleMISP

// ModuleMISP инициализированный модуль
// chanInputMISP - канал для отправки данных в модуль
// chanOutputMISP - канал для принятия данных из модуля
// ChanLogging - канал для отправки логов
type ModuleMISP struct {
	chanInputMISP  chan interface{}
	chanOutputMISP chan interface{}
	ChanLogging    chan<- datamodels.MessageLoging
}

func init() {
	mmisp = ModuleMISP{
		chanInputMISP:  make(chan interface{}),
		chanOutputMISP: make(chan interface{}),
	}
}

func NewClientMISP(
	ctx context.Context,
	conf confighandler.AppConfigMISP,
	chanLog chan<- datamodels.MessageLoging) (*ModuleMISP, error) {
	fmt.Println("func 'NewClientMISP', START...")

	mmisp.ChanLogging = chanLog

	/*
		Здесь нужно написать инициализацию подключения к MISP
	*/

	return &mmisp, nil
}

func (mmisp ModuleMISP) GetDataReceptionChannel() <-chan interface{} {
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

	return mmisp.chanOutputMISP
}

func (mmisp ModuleMISP) GettingData() interface{} {
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

	return <-mmisp.chanOutputMISP
}

func (mmisp ModuleMISP) SendingData(data interface{}) {
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

	mmisp.chanInputMISP <- data
}
