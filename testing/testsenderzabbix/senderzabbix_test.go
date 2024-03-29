package testsenderzabbix_test

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"placeholder_misp/datamodels"
	"placeholder_misp/zabbixinteractions"
)

var _ = Describe("Senderzabbix", Ordered, func() {
	var (
		testuchetdb *zabbixinteractions.HandlerZabbixConnection
		zcErr       error
		ctx         context.Context
		ctxCancel   context.CancelFunc

		eventTypes []zabbixinteractions.EventType = []zabbixinteractions.EventType{
			{
				IsTransmit: true,
				EventType:  "error",
				ZabbixKey:  "error_bav",
			},
			{
				IsTransmit: true,
				EventType:  "info",
				ZabbixKey:  "info_bav",
			},
			{
				IsTransmit: true,
				EventType:  "handshake",
				ZabbixKey:  "handshake_bav",
				Handshake: zabbixinteractions.Handshake{
					TimeInterval: 1,
					Message:      "I'm still alive",
				},
			},
		}
		events []datamodels.MessageLogging = []datamodels.MessageLogging{
			{
				MsgType: "error",
				MsgData: "ERROR: test error message",
			},
			{
				MsgType: "warning",
				MsgData: "WARNING: test warning message",
			},
			{
				MsgType: "info",
				MsgData: "test message with information about app",
			},
		}

		chanErr chan error
	)

	BeforeAll(func() {
		chanErr = make(chan error)
		connTimeout := time.Duration(7 * time.Second)

		go func() {
			for err := range chanErr {
				fmt.Println("------------- ERROR -------------")
				fmt.Println(err)
			}
		}()

		ctx, ctxCancel = context.WithCancel(context.Background())
		testuchetdb, zcErr = zabbixinteractions.NewZabbixConnection(
			ctx,
			zabbixinteractions.SettingsZabbixConnection{
				Port: 10051,
				//Host:              "192.168.9.45", //правильный
				Host:              "192.168.9.145", //не правильный
				NetProto:          "tcp",
				ZabbixHost:        "test-uchet-db.cloud.gcm",
				ConnectionTimeout: &connTimeout,
			},
			chanErr)
	})

	Context("Тест 0. Проверяем на наличие ошибок при выполнении NewZabbixConnection", func() {
		It("Не должно быть ошибок при инициализации NewZabbixConnection", func() {
			Expect(zcErr).ShouldNot(HaveOccurred())
		})
	})

	Context("Тест 1. Пробуем выполнить соединение с Zabbix", Ordered, func() {
		var (
			err      error
			chanDone chan struct{}
			msgChan  chan zabbixinteractions.MessageSettings
		)

		BeforeAll(func() {
			msgChan = make(chan zabbixinteractions.MessageSettings)
			chanDone = make(chan struct{})
			err = testuchetdb.Handler(eventTypes, msgChan)

			go func() {
				for k, v := range events {
					fmt.Printf("%d. send message type %s\n", k, v.MsgType)

					msgChan <- zabbixinteractions.MessageSettings{
						EventType: v.MsgType,
						Message:   v.MsgData,
					}

					time.Sleep(time.Duration(1 * time.Second))
				}

				time.Sleep(time.Duration(2 * time.Second))

				ctxCancel()
				close(chanErr)

				chanDone <- struct{}{}
			}()
		})

		It("Соединение с Zabbix должно быть успешно установлено", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("Должно придти сообщение о завершении обработки", func() {
			/*
					!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
				Написал обработчик для взаимодействия с Zabbix
				тесты вроде проходят нормально, НО лучьше на
				свежую голову внимательно проверить все тесты
				ЕЩЕ раз
					!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
				Похоже ошибки подключения не обрабатываются, из
				канала chanErr НИЧЕГО не приходит, даже если
				ip адрес сервера Zabbix не правильно указан
			*/

			Expect(<-chanDone).Should(Equal(struct{}{}))
		})
	})
	/*
		Context("Тест 1. Пробуем выполнить соединение с Zabbix", func() {
			It("Соединение с Zabbix должно быть успешно установлено", func() {
				var d net.Dialer
				ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
				defer cancel()

				conn, err := d.DialContext(ctx, "tcp", "zabbix.cloud.gcm:10051")
				Expect(err).ShouldNot(HaveOccurred())
				defer conn.Close()

				Expect(true).Should(BeTrue())
			})
		})

		Context("Тест 2. Проверяем возможность подключения и отправки данных в Zabbix", func() {
			It("При отправки данных в Zabbix не должно быть ошибок", func() {

				//для подтверждения что модуль
				num, err := zc.SendData([]string{"I'm still alive"})

				fmt.Println("Count sended byte:", num)

				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	*/
})
