package coremodule

import (
	"fmt"

	"github.com/av-belyakov/simplelogger"

	"placeholder_misp/mispinteractions"
	"placeholder_misp/natsinteractions"
)

func NewCore(
	natsmodule natsinteractions.EnumChannelsNATS,
	mispmodule mispinteractions.EnumChannelsMISP,
	sl simplelogger.SimpleLoggerSettings) {
	fmt.Println("func 'NewCore', START...")

	natsChanReception := natsmodule.GetDataReceptionChannel()
	mispChanReception := mispmodule.GetDataReceptionChannel()

	for {
		select {
		case data := <-natsChanReception:
			//fmt.Println("func 'NewCore', NATS reseived message from chanOutNATS: ", data)
			_ = sl.WriteLoggingData(fmt.Sprintln(data), "info")

		case data := <-mispChanReception:
			fmt.Println("func 'NewCore', MISP reseived message from chanOutMISP: ", data)

		}
	}
}
