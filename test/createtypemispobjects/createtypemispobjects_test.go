package createtypemispobjects_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/av-belyakov/placeholder_misp/cmd/coremodule"
	"github.com/av-belyakov/placeholder_misp/cmd/mispapi"
	"github.com/av-belyakov/placeholder_misp/commoninterfaces"
	"github.com/av-belyakov/placeholder_misp/internal/datamodels"
	"github.com/av-belyakov/placeholder_misp/internal/logginghandler"
	"github.com/av-belyakov/placeholder_misp/memorytemporarystorage"
	rules "github.com/av-belyakov/placeholder_misp/rulesinteraction"
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
		ctx, ctxCancel = context.WithCancel(context.Background())
		counting := make(chan datamodels.DataCounterSettings)

		moduleMisp = &mispapi.ModuleMISP{
			ChanInput:  make(chan mispapi.InputSettings),
			ChanOutput: make(chan mispapi.OutputSetting),
		}

		logging := logginghandler.New()

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
