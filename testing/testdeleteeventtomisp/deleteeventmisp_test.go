package testdeleteeventtomisp_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"placeholder_misp/confighandler"
	"placeholder_misp/datamodels"
	"placeholder_misp/memorytemporarystorage"
	"placeholder_misp/mispinteractions"
)

var _ = Describe("Deleteeventmisp", Ordered, func() {
	var (
		confApp    confighandler.ConfigApp
		storageApp *memorytemporarystorage.CommonStorageTemporary
		//redismodule               *redisinteractions.ModuleRedis
		mispmodule                *mispinteractions.ModuleMISP
		errConfApp, errMispModule error
	)

	logging := make(chan datamodels.MessageLogging)

	BeforeAll(func() {
		confApp, errConfApp = confighandler.NewConfig()
		confApp.AppConfigRedis = confighandler.AppConfigRedis{
			Host: "192.168.9.208",
			Port: 16379,
		}
		confApp.AppConfigMISP = confighandler.AppConfigMISP{
			Host: "misp-world.cloud.gcm",
			Auth: "TvHkjH8jVQEIdvAxjxnL4H6wDoKyV7jobDjndvAo",
		}

		//инициализируем модуль временного хранения информации
		storageApp = memorytemporarystorage.NewTemporaryStorage()
		//redismodule = redisinteractions.HandlerRedis(context.Background(), confApp.AppConfigRedis, storageApp, logging)
		mispmodule, errMispModule = mispinteractions.HandlerMISP(confApp.AppConfigMISP, storageApp, logging)
	})

	Context("Тест 1. Проверка успешной инициализации модулей", func() {
		It("При инициализации модулей ошибок быть не должно", func() {
			Expect(errConfApp).ShouldNot(HaveOccurred())
			Expect(errMispModule).ShouldNot(HaveOccurred())
		})
	})

	Context("Тест 2. Проверка работы канала передачи в MISP события для удаления", func() {
		It("Должно быть успешно переданно событие", func() {
			chanDone := make(chan struct{})

			go func() {
				fmt.Println("___ Logging START")
				defer fmt.Println("___ Logging STOP")

				for log := range logging {
					if log.MsgData == "TEST_INFO STOP" {
						chanDone <- struct{}{}

						return
					}

					fmt.Println("----", log, "----")
				}
			}()

			mispmodule.SendingDataInput(mispinteractions.SettingsChanInputMISP{
				Command: "del event by id",
				EventId: "7418",
			})

			<-chanDone

			time.Sleep(1 * time.Second)
			fmt.Println("STOP CHAN TESTING")

			Expect(true).Should(BeTrue())
		})
	})

	/*
		Context("Тест 3. Проверка удаления события из MISP по его id", func() {
			It("Должно быть выполненно успешное удаление из MISP события по его id", func() {
				chanDone := make(chan bool)

				go func() {
					fmt.Println("___ Logging START")
					defer fmt.Println("___ Logging STOP")

					for log := range logging {
						fmt.Println("----", log, "----")
					}
				}()

				go func() {
					mispChanReception := mispmodule.GetDataReceptionChannel()
					redisChanReception := redismodule.GetDataReceptionChannel()

					for {
						select {
						case data := <-mispChanReception:
							switch data.Command {
							case "set new event id":
								// ***********************************
								// Это логирование только для теста!!!
								// ***********************************
								logging <- datamodels.MessageLogging{
									MsgData: fmt.Sprintf("TEST_INFO func 'NewCore', надо отправить инфу CaseID '%s' и EventId '%s' to REDIS DB\n", data.CaseId, data.EventId),
									MsgType: "testing",
								}
								//
								//

								//обработка запроса на добавления новой связки caseId:eventId в Redis
								redismodule.SendingDataInput(redisinteractions.SettingsChanInputRedis{
									Command: "set case id",
									Data:    fmt.Sprintf("%s:%s", data.CaseId, data.EventId),
								})
							case "TEST STOP":
								fmt.Println("TEST STOP --====---==-=")

								chanDone <- true
							}

						case data := <-redisChanReception:
							switch data.CommandResult {
							case "found event id":
								// ***********************************
								// Это логирование только для теста!!!
								// ***********************************
								logging <- datamodels.MessageLogging{
									MsgData: fmt.Sprintf("TEST_INFO func 'NewCore', Здесь, получаем event id: '%v' из Redis для удаления события в MISP", data.Result),
									MsgType: "testing",
								}
								//
								//

								// Здесь, получаем eventId из Redis для удаления события в MISP
								eventId, ok := data.Result.(string)
								if !ok {
									_, f, l, _ := runtime.Caller(0)

									logging <- datamodels.MessageLogging{
										MsgData: fmt.Sprintf("'it is not possible to convert a value to a string' %s:%d", f, l-1),
										MsgType: "warning",
									}

									break
								}

								// ***********************************
								// Это логирование только для теста!!!
								// ***********************************
								logging <- datamodels.MessageLogging{
									MsgData: fmt.Sprintf("TEST_INFO func 'NewCore', отправляем event id: '%s' в MISP для удаления события", eventId),
									MsgType: "testing",
								}
								//
								//

								mispmodule.SendingDataInput(mispinteractions.SettingsChanInputMISP{
									Command: "del event by id",
									EventId: eventId,
								})
							}
						}
					}
				}()

				mispmodule.SendingDataOutput(mispinteractions.SettingChanOutputMISP{
					Command: "set new event id",
					CaseId:  fmt.Sprint(33669),
					EventId: fmt.Sprint(7342),
				})

				isDone := <-chanDone
				close(logging)

				time.Sleep(1 * time.Second)
				fmt.Println("STOP TESTING")

				Expect(isDone).Should(BeTrue())
			})
		})
	*/
})
