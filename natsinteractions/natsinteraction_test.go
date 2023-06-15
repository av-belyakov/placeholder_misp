package natsinteractions_test

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"placeholder_misp/confighandler"
	"placeholder_misp/natsinteractions"
)

var _ = Describe("Natsinteraction", Ordered, func() {
	var (
		ctx          context.Context
		errConn      error
		closeCtx     context.CancelFunc
		enumChannels *natsinteractions.EnumChannelsNATS
	)

	BeforeAll(func() {
		ctx, closeCtx = context.WithTimeout(context.Background(), 2*time.Second)

		enumChannels, errConn = natsinteractions.NewClientNATS(ctx, confighandler.AppConfigNATS{
			//Host: "nats.cloud.gcm",
			Host: "127.0.0.1",
			Port: 4222,
		})
	})

	Context("Тест 1. Проверка инициализации соединения с NATS", func() {
		It("При инициализации соединения с NATS не должно быть ошибки", func() {
			Expect(errConn).ShouldNot(HaveOccurred())

			/*

				Тест не проходит успешно. Поставил и запустил nats-server в консоле, пробую с ним, тот же результат

			*/

			fmt.Println("Resevid message = ", <-enumChannels.GetDataReceptionChannel())
		})
	})

	AfterAll(func() {
		closeCtx()
	})
})
