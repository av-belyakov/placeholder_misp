package createtypemispobjects_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/av-belyakov/placeholder_misp/cmd/coremodule"
	"github.com/av-belyakov/placeholder_misp/cmd/mispapi"
	"github.com/av-belyakov/placeholder_misp/commoninterfaces"
	"github.com/av-belyakov/placeholder_misp/internal/countermessage"
	"github.com/av-belyakov/placeholder_misp/internal/logginghandler"
	rules "github.com/av-belyakov/placeholder_misp/internal/ruleshandler"
)

//GetDataReceptionChannel() <-chan mispiapi.ChanOutputSetting
//SendingDataOutput(mispiapi.ChanOutputSetting)
//SendingDataInput(mispiapi.ChanInputSettings)

var _ = Describe("Createtypemispobjects", Ordered, func() {
	var (
		rootDir     string = "placeholder_misp"
		rulesDir    string = "rules"
		rulesFile   string = "mispmsgrule.yaml"
		exampleFile string = "../test_json/example_3.json"
		taskId      string = "new_task_1"

		moduleMisp *mispapi.ModuleMISP

		ctx       context.Context
		ctxCancel context.CancelFunc
	)

	BeforeAll(func() {
		chZabbix := make(chan commoninterfaces.Messager)

		ctx, ctxCancel = context.WithCancel(context.Background())
		counting := countermessage.New(chZabbix)

		moduleMisp = &mispapi.ModuleMISP{
			ChanInput:  make(chan mispapi.InputSettings),
			ChanOutput: make(chan mispapi.OutputSetting),
		}

		logging := logginghandler.New()

		b, err := os.ReadFile(exampleFile)
		Expect(err).ShouldNot(HaveOccurred())

		//инициализируем модуль чтения правил обработки MISP сообщений
		lr, _, err := rules.NewListRule(rootDir, rulesDir, rulesFile)
		Expect(err).ShouldNot(HaveOccurred())

		//сообщения для логирования
		go func(ctx context.Context, logging commoninterfaces.Logger) {
			for {
				select {
				case <-ctx.Done():
					return

				case msg := <-logging.GetChan():
					fmt.Printf("log type:'%s', message:'%s'\n", msg.GetType(), msg.GetMessage())
				}
			}
		}(ctx, logging)

		//заглушка для счетчика
		go func(ctx context.Context, chz <-chan commoninterfaces.Messager) {
			for {
				select {
				case <-ctx.Done():
					return

				case c := <-chz:
					fmt.Printf("\tmessage type:'%s', message:'%s'\n", c.GetType(), c.GetMessage())
				}
			}
		}(ctx, chZabbix)

		handler := coremodule.NewHandlerJsonMessage(counting, logging)
		chDecode := handler.HandlerJsonMessage(b, "twe83475js")

		go coremodule.NewMispFormat(chDecode, taskId, moduleMisp, lr, counting, logging)
	})

	Context("Тест 1. Формирование документов в формате MISP", func() {
		It("", func() {
			msg := <-moduleMisp.GetInputChannel()

			b, err := json.MarshalIndent(msg, "", " ")
			Expect(err).ShouldNot(HaveOccurred())

			fmt.Println("MSG:", string(b))

		})
	})

	AfterAll(func() {
		ctxCancel()
	})
})
