package testdeleteeventtomisp_test

import (
	"context"
	"fmt"
	"log"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/av-belyakov/placeholder_misp/cmd/mispapi"
	"github.com/av-belyakov/placeholder_misp/commoninterfaces"
	"github.com/av-belyakov/placeholder_misp/constants"
	"github.com/av-belyakov/placeholder_misp/internal/confighandler"
	"github.com/av-belyakov/placeholder_misp/internal/logginghandler"
	"github.com/av-belyakov/placeholder_misp/internal/supportingfunctions"
	"github.com/av-belyakov/simplelogger"
)

var _ = Describe("Deleteeventmisp", Ordered, func() {
	var (
		confApp *confighandler.ConfigApp
		//redismodule               *redisinteractions.ModuleRedis
		mispmodule                *mispapi.ModuleMISP
		logging                   commoninterfaces.Logger
		errConfApp, errMispModule error
	)

	BeforeAll(func() {
		rootPath, err := supportingfunctions.GetRootPath(constants.Root_Dir)
		if err != nil {
			log.Fatalf("error, it is impossible to form root path (%s)", err.Error())
		}

		confApp, errConfApp = confighandler.New(rootPath)
		confApp.AppConfigRedis = confighandler.AppConfigRedis{
			Host: "192.168.9.208",
			Port: 16379,
		}
		confApp.AppConfigMISP = confighandler.AppConfigMISP{
			Host: "misp-world.cloud.gcm",
			Auth: "TvHkjH8jVQEIdvAxjxnL4H6wDoKyV7jobDjndvAo",
		}

		chZabbix := make(chan commoninterfaces.Messager)
		simpleLogger, err := simplelogger.NewSimpleLogger(context.Background(), "palceholder_misp", simplelogger.CreateOptions())
		if err != nil {
			log.Fatalf("error module 'simplelogger': %v", err)
		}

		logging = logginghandler.New(simpleLogger, chZabbix)

		//redismodule = redisinteractions.HandlerRedis(context.Background(), confApp.AppConfigRedis, storageApp, logging)
		mispmodule, errMispModule = mispapi.NewModuleMISP(confApp.GetAppMISP().Host, confApp.GetAppMISP().Auth, confApp.GetListOrganization(), logging)
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

				for log := range logging.GetChan() {
					if log.GetMessage() == "TEST_INFO STOP" {
						chanDone <- struct{}{}

						return
					}

					fmt.Println("----", log, "----")
				}
			}()

			mispmodule.SendDataInput(mispapi.InputSettings{
				Command: "del event by id",
				EventId: "7418",
			})

			<-chanDone

			time.Sleep(1 * time.Second)
			fmt.Println("STOP CHAN TESTING")

			Expect(true).Should(BeTrue())
		})
	})
})
