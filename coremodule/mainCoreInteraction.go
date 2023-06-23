package coremodule

import (
	"fmt"
	"runtime"

	"placeholder_misp/datamodels"
	"placeholder_misp/mispinteractions"
	"placeholder_misp/natsinteractions"
	"placeholder_misp/supportingfunctions"
)

func NewCore(
	natsmodule natsinteractions.ModuleNATS,
	mispmodule mispinteractions.ModuleMISP,
	msgOutChan chan<- datamodels.MessageLoging) {
	fmt.Println("func 'NewCore', START...")

	natsChanReception := natsmodule.GetDataReceptionChannel()
	mispChanReception := mispmodule.GetDataReceptionChannel()

	for {
		select {
		case data := <-natsChanReception:
			strMsg, err := supportingfunctions.NewReadReflectJSONSprint(data)
			if err != nil {
				_, f, l, _ := runtime.Caller(0)

				msgOutChan <- datamodels.MessageLoging{
					MsgData: fmt.Sprintf("%s %s:%d", fmt.Sprint(err), f, l-2),
					MsgType: "error",
				}

				continue
			}

			msgOutChan <- datamodels.MessageLoging{
				MsgData: strMsg,
				MsgType: "info",
			}

		case data := <-mispChanReception:
			fmt.Println("func 'NewCore', MISP reseived message from chanOutMISP: ", data)

		}
	}
}
