package coremodule

import (
	"encoding/json"
	"fmt"
	"runtime"

	"placeholder_misp/datamodels"
	"placeholder_misp/mispinteractions"
	"placeholder_misp/natsinteractions"
	"placeholder_misp/rules"
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
			//fmt.Println("func 'NewCore', NATS reseived message from chanOutNATS: ", data)
			result := map[string]interface{}{}

			err := json.Unmarshal(data, &result)
			if err != nil {
				_, f, l, _ := runtime.Caller(0)

				msgOutChan <- datamodels.MessageLoging{
					MsgData: fmt.Sprintf("%s %s:%d", fmt.Sprint(err), f, l-2),
					MsgType: "error",
				}

				continue
			}

			strMsg := ReadReflectMapSprint(result, rules.ListRulesProcessedMISPMessage{}, 0)
			//_ = sl.WriteLoggingData(strMsg, "info")

			msgOutChan <- datamodels.MessageLoging{
				MsgData: strMsg,
				MsgType: "info",
			}

		case data := <-mispChanReception:
			fmt.Println("func 'NewCore', MISP reseived message from chanOutMISP: ", data)

		}
	}
}
