package mispinteractions

import (
	"context"
	"fmt"

	"placeholder_misp/confighandler"
)

var enumChannels EnumChannelsMISP

// enumChannelsMISP перечень каналов для взаимодействия с модулем
// chanInput - канал для отправки данных в модуль
// chanOutput - канал для принятия данных из модуля
type EnumChannelsMISP struct {
	chanInput  chan interface{}
	chanOutput chan interface{}
}

func init() {
	enumChannels = EnumChannelsMISP{
		chanInput:  make(chan interface{}),
		chanOutput: make(chan interface{}),
	}
}

func NewClientMISP(ctx context.Context, conf confighandler.AppConfigMISP) (*EnumChannelsMISP, error) {
	fmt.Println("func 'NewClientMISP', START...")

	return &enumChannels, nil
}

func (ecmisp EnumChannelsMISP) GetDataReceptionChannel() <-chan interface{} {
	return ecmisp.chanOutput
}

func (ecmisp EnumChannelsMISP) GettingData() interface{} {
	return <-ecmisp.chanOutput
}

func (ecmisp EnumChannelsMISP) SendingData(data interface{}) {
	ecmisp.chanInput <- data
}
