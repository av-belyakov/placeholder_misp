package createtypemispobjects_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/av-belyakov/placeholder_misp/cmd/coremodule"
	"github.com/av-belyakov/placeholder_misp/cmd/mispapi"
	"github.com/av-belyakov/placeholder_misp/commoninterfaces"
	"github.com/av-belyakov/placeholder_misp/internal/countermessage"
	"github.com/av-belyakov/placeholder_misp/internal/logginghandler"
	rules "github.com/av-belyakov/placeholder_misp/internal/ruleshandler"
	"github.com/av-belyakov/placeholder_misp/test/createtypemispobjects"
	"github.com/av-belyakov/simplelogger"
)

const (
	Root_Dir     string = "placeholder_misp"
	Rules_Dir    string = "rules"
	Rules_File   string = "mispmsgrule.yml"
	Example_File string = "../test_json/event_39100.json"
	Task_Id      string = "7s7qeytyyy2e27tr73213143a"
)

var (
	counting  *countermessage.CounterMessage
	logging   *logginghandler.LoggingChan
	chZabbix  chan commoninterfaces.Messager
	listRules *rules.ListRule
)

// ModuleMISPForTest имитация подключения к MISP API (только для тестов)
type ModuleMISPForTest struct {
	chInput  chan mispapi.InputSettings
	chOutput chan mispapi.OutputSetting
}

func NewModuleMISPForTest() *ModuleMISPForTest {
	return &ModuleMISPForTest{
		chInput:  make(chan mispapi.InputSettings),
		chOutput: make(chan mispapi.OutputSetting),
	}
}

func (m *ModuleMISPForTest) GetReceptionChannel() <-chan mispapi.OutputSetting {
	return m.chOutput
}

func (m *ModuleMISPForTest) SendDataOutput(s mispapi.OutputSetting) {
	m.chOutput <- s
}

func (m *ModuleMISPForTest) GetOutputChannel() <-chan mispapi.OutputSetting {
	return m.chOutput
}

func (m *ModuleMISPForTest) SendDataInput(s mispapi.InputSettings) {
	m.chInput <- s
}

func (m *ModuleMISPForTest) GetInputChannel() <-chan mispapi.InputSettings {
	return m.chInput
}

func TestMain(m *testing.M) {
	chZabbix = make(chan commoninterfaces.Messager)
	counting = countermessage.New(chZabbix)

	simpleLogger, err := simplelogger.NewSimpleLogger(context.Background(), "palceholder_misp", simplelogger.CreateOptions())
	if err != nil {
		log.Fatalf("error module 'simplelogger': %v", err)
	}

	logging = logginghandler.New(simpleLogger, chZabbix)

	//инициализируем модуль чтения правил обработки MISP сообщений
	listRules, _, err = rules.NewListRule(Root_Dir, Rules_Dir, Rules_File)
	if err != nil {
		log.Fatalln(err)
	}

	os.Exit(m.Run())
}

func TestCreateMispObjects(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//это что бы работал подсчёт, без него каналы будут не активны
	counting.Start(ctx)

	// сообщения для логирования
	go func(ctx context.Context, logging commoninterfaces.Logger) {
		for {
			select {
			case <-ctx.Done():
				t.Log("logging STOP")

				return

			case msg := <-logging.GetChan():
				fmt.Printf("log type:'%s', message:'%s'\n", msg.GetType(), msg.GetMessage())
			}
		}
	}(ctx, logging)
	logging.Send("test_logging_message", "check logging module")

	// заглушка для счётчика
	go func(ctx context.Context, chz <-chan commoninterfaces.Messager) {
		for {
			select {
			case <-ctx.Done():
				t.Log("counting STOP")

				return

			case c := <-chz:
				fmt.Printf("\tmessage type:'%s', message:'%s'\n", c.GetType(), c.GetMessage())
			}
		}
	}(ctx, chZabbix)
	counting.SendMessage("test_countiong_message", 100)

	//чтение файла с примером
	b, err := os.ReadFile(Example_File)
	assert.NoError(t, err)

	handler := coremodule.NewHandlerJSON(counting, logging)
	chDecode := handler.Start(b, Task_Id)

	moduleMisp := NewModuleMISPForTest()
	go coremodule.CreateObjectsFormatMISP(chDecode, Task_Id, moduleMisp, listRules, counting, logging)

	t.Run("Формирование документов в формате MISP", func(t *testing.T) {
		msg := <-moduleMisp.GetInputChannel()

		/*
			!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

			выполнить этот тест через debug

			!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		*/

		createtypemispobjects.AddNewObject(
			context.Background(),
			msg,
			createtypemispobjects.OptionsAddNewObject{
				Host:        "misp-center.cloud.gcm",
				AuthKey:     os.Getenv("GO_PHMISP_MAUTH"),
				UserAuthKey: os.Getenv("GO_PHMISP_MAUTH"),
			})

		b, err = json.MarshalIndent(msg, "", " ")
		assert.NoError(t, err)

		fmt.Println("MSG:", string(b))
	})
}
