package coremodule

import (
	"fmt"
	"runtime"

	"placeholder_misp/datamodels"
	"placeholder_misp/elasticsearchinteractions"
	"placeholder_misp/mispinteractions"
	"placeholder_misp/natsinteractions"
	"placeholder_misp/nkckiinteractions"
	rules "placeholder_misp/rulesinteraction"
	//"placeholder_misp/supportingfunctions"
)

func NewCore(
	natsmodule natsinteractions.ModuleNATS,
	mispmodule mispinteractions.ModuleMISP,
	elmodule elasticsearchinteractions.ModuleElasticSearch,
	nkckimodule nkckiinteractions.ModuleNKCKI,
	listRules rules.ListRulesProcessingMsgMISP,
	msgOutChan chan<- datamodels.MessageLoging) {
	fmt.Println("func 'NewCore', START...")

	natsChanReception := natsmodule.GetDataReceptionChannel()
	mispChanReception := mispmodule.GetDataReceptionChannel()

	for {
		select {
		case data := <-natsChanReception:
			//вот здесь будет основной обработчик
			//newByte, listOk, err := coremodule.NewProcessingInputMessageFromHive(eb, listRules)

			procMsg, err := NewHandleMessageFromHive(data, listRules)
			if err != nil {
				_, f, l, _ := runtime.Caller(0)

				msgOutChan <- datamodels.MessageLoging{
					MsgData: fmt.Sprintf("%s %s:%d", fmt.Sprint(err), f, l-2),
					MsgType: "error",
				}

				continue
			}

			chanInput, chanDone := NewMispFormat()

			//обработка сообщения полученного от Hive на основе правил
			skipMsg, warningMsg := procMsg.HandleMessage(chanInput)

			if len(warningMsg) > 0 {
				for _, v := range warningMsg {
					_, f, l, _ := runtime.Caller(0)

					msgOutChan <- datamodels.MessageLoging{
						MsgData: fmt.Sprintf("%s %s:%d", v, f, l-2),
						MsgType: "warning",
					}
				}
			}

			chanDone <- struct{}{}

			//если сообщение соответствовало правилам
			if skipMsg {
				//отправка данных через канал в модуль взаимодействия с MISP
				mispmodule.SendingData(procMsg.Message)
			}

			/*
				Внимание!!!

				все полученные от NATS сообщения нужно дублировать по всем модулям MISP, NKCKI и ElasticSearch
				Вопрос только их дублировать после обработки по правилам (правила именно в контексте замены Replace)
				или обработка по правилам должна выполнятся только для модуля MISP

			*/
			// отправка сообщения в Elasticshearch (пока заглушка)
			elmodule.SendingData(procMsg.Message)
			// отправка сообщения в НКЦКИ (пока заглушка)
			nkckimodule.SendingData(procMsg.Message)

			//	ДЛЯ ТЕСТОВ
			//это только вывод данных в файл по сообщению из Have
			//
			/*strMsg, err := supportingfunctions.NewReadReflectJSONSprint(data)
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
			}*/

		case data := <-mispChanReception:
			fmt.Println("func 'NewCore', MISP reseived message from chanOutMISP: ", data)

		}
	}
}
