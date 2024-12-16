package createtypemispobjects_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"placeholder_misp/coremodule"
	"placeholder_misp/datamodels"
	"placeholder_misp/memorytemporarystorage"
	"placeholder_misp/mispinteractions"
	rules "placeholder_misp/rulesinteraction"
)

//GetDataReceptionChannel() <-chan mispinteractions.ChanOutputSetting
//SendingDataOutput(mispinteractions.ChanOutputSetting)
//SendingDataInput(mispinteractions.ChanInputSettings)

var _ = Describe("Createtypemispobjects", Ordered, func() {
	var (
		rootDir     string = "placeholder_misp"
		rulesDir    string = "rules"
		rulesFile   string = "mispmsgrule.yaml"
		exampleFile string = "../test_json/example_3.json"
		taskId      string = "new_task_1"

		moduleMisp *mispinteractions.ModuleMISP

		ctx       context.Context
		ctxCancel context.CancelFunc
	)

	BeforeAll(func() {
		ctx, ctxCancel = context.WithCancel(context.Background())
		logging := make(chan datamodels.MessageLogging)
		counting := make(chan datamodels.DataCounterSettings)

		moduleMisp = &mispinteractions.ModuleMISP{
			ChanInput:  make(chan mispinteractions.InputSettings),
			ChanOutput: make(chan mispinteractions.OutputSetting),
		}

		b, err := os.ReadFile(exampleFile)
		Expect(err).ShouldNot(HaveOccurred())

		//f, err = os.Open(exampleFile)

		//b := []byte(nil)
		//_, err = f.Read(b)
		//Expect(err).ShouldNot(HaveOccurred())

		//wr := bytes.Buffer{}
		//sc := bufio.NewScanner(f)
		//for sc.Scan() {
		//wr.WriteString(sc.Text())
		//if _, err := wr.Write(sc.Bytes()); err != nil {
		//log.Fatalln(err)
		//}
		//}
		//инициализируем модуль чтения правил обработки MISP сообщений
		lr, _, err := rules.NewListRule(rootDir, rulesDir, rulesFile)
		Expect(err).ShouldNot(HaveOccurred())

		//сообщения для логирования
		go func(ctx context.Context, chl <-chan datamodels.MessageLogging) {
			for {
				select {
				case <-ctx.Done():
					return

				case msg := <-chl:
					fmt.Printf("log type:'%s', message:'%s'\n", msg.MsgType, msg.MsgData)
				}
			}
		}(ctx, logging)

		//заглушка для счетчика
		go func(ctx context.Context, chc <-chan datamodels.DataCounterSettings) {
			for {
				select {
				case <-ctx.Done():
					return

				case c := <-chc:
					fmt.Printf("\tcounting:'%d'\n", c.Count)
				}
			}
		}(ctx, counting)

		//инициализируем модуль временного хранения информации
		storageApp := memorytemporarystorage.NewTemporaryStorage()

		//добавляем время инициализации счетчика хранения
		storageApp.SetStartTimeDataCounter(time.Now())

		handler := coremodule.NewHandlerJsonMessage(storageApp, logging, counting)
		chDecode := handler.HandlerJsonMessage(b, "twe83475js")

		go coremodule.NewMispFormat(chDecode, taskId, moduleMisp, lr, logging, counting)
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
