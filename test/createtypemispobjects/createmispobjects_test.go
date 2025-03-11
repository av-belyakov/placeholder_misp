package createtypemispobjects_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	"github.com/av-belyakov/placeholder_misp/cmd/coremodule"
	"github.com/av-belyakov/placeholder_misp/cmd/mispapi"
	"github.com/av-belyakov/placeholder_misp/cmd/sqlite3api"
	"github.com/av-belyakov/placeholder_misp/commoninterfaces"
	"github.com/av-belyakov/placeholder_misp/internal/countermessage"
	"github.com/av-belyakov/placeholder_misp/internal/logginghandler"
	rules "github.com/av-belyakov/placeholder_misp/internal/ruleshandler"
	"github.com/av-belyakov/placeholder_misp/test/createtypemispobjects"
	"github.com/av-belyakov/simplelogger"
)

const (
	Root_Dir      string = "placeholder_misp"
	Rules_Dir     string = "rules"
	Rules_File    string = "mispmsgrule.yml"
	Example_File  string = "../test_json/event_39100.json"
	Task_Id       string = "7s7qeytyyy2e27tr73213143a"
	sqlite3FileDb string = "../../backupdb/sqlite3_backup.db"
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
	if err := godotenv.Load("../../.env"); err != nil {
		log.Panicln(err)
	}

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

	// инициализация модуля взаимодействия с Sqlite3
	sqlite3Module, err := sqlite3api.New(ctx, sqlite3FileDb, logging)
	assert.NoError(t, err)

	//чтение файла с примером
	b, err := os.ReadFile(Example_File)
	assert.NoError(t, err)

	handler := coremodule.NewHandlerJSON(counting, logging)
	chDecode := handler.Start(b, Task_Id)

	moduleMisp := NewModuleMISPForTest()
	go coremodule.CreateObjectsFormatMISP(chDecode, Task_Id, moduleMisp, sqlite3Module, listRules, counting, logging)

	/*
		!!!!!!!!!!!!!!!!!!!!!!!!!!

			надо потестировать добавление события в MISP с учётом взаимодействия
			с Sqlite3 database

		!!!!!!!!!!!!!!!!!!!!!!!!!!
	*/

	t.Run("Формирование документов в формате MISP", func(t *testing.T) {
		//var eventId string
		msg := <-moduleMisp.GetInputChannel()

		/*rmisp, err := mispapi.NewMispRequest(
			mispapi.WithHost("misp-center.cloud.gcm"),
			mispapi.WithUserAuthKey("GO_PHMISP_MAUTH"),
			mispapi.WithMasterAuthKey("GO_PHMISP_UAUTH"))
		assert.NoError(t, err)

		t.Run("Добавление event", func(t *testing.T) {
			_, resBodyByte, err := rmisp.SendEvent_ForTest(ctx, msg.Data.GetEvent())
			assert.NoError(t, err)

			resMisp := mispapi.MispResponse{}
			err = json.Unmarshal(resBodyByte, &resMisp)
			assert.NoError(t, err)

			log.Println("func 'specialObject.SetFunc', MISP response:", resMisp)

			//получаем уникальный id MISP
			for key, value := range resMisp.Event {
				if key == "id" {
					if str, ok := value.(string); ok {
						eventId = str

						break
					}
				}
			}

			log.Println("func 'specialObject.SetFunc', MISP eventId:", eventId)
			assert.NotEmpty(t, eventId)
		})

		t.Run("Добавление event_reports", func(t *testing.T) {
			err := rmisp.SendEventReports_ForTest(ctx, eventId, msg.Data.GetReports())
			assert.NoError(t, err)
		})

		t.Run("Добавление attribytes", func(t *testing.T) {
			_, _, warning, err := rmisp.SendAttribytes_ForTest(ctx, eventId, msg.Data.GetAttributes())
			assert.NoError(t, err)
			t.Log("warning:", warning)
		})

		t.Run("Добавление objects", func(t *testing.T) {
			_, _, err := rmisp.SendObjects_ForTest(ctx, eventId, msg.Data.GetObjects())
			assert.NoError(t, err)
		})

		t.Run("Добавление event_tags", func(t *testing.T) {
			err := rmisp.SendEventTags_ForTest(ctx, eventId, msg.Data.GetObjectTags())
			assert.NoError(t, err)
		})

		t.Run("Публикация события", func(t *testing.T) {
			resMsg, err := rmisp.SendRequestPublishEvent_ForTest(ctx, eventId)
			assert.NoError(t, err)

			t.Log("response:", resMsg)
			assert.NotEmpty(t, resMsg)
		})*/

		createtypemispobjects.AddNewObject(
			context.Background(),
			msg,
			sqlite3Module,
			createtypemispobjects.OptionsAddNewObject{
				Host:        "misp-center.cloud.gcm",
				AuthKey:     os.Getenv("GO_PHMISP_MAUTH"),
				UserAuthKey: os.Getenv("GO_PHMISP_UAUTH"),
			})

		time.Sleep(3 * time.Second)
		//b, err = json.MarshalIndent(msg, "", " ")
		//assert.NoError(t, err)

		chRes := make(chan sqlite3api.Response)
		sqlite3Module.SendDataToModule(sqlite3api.Request{
			Command:    "search caseId",
			ChResponse: chRes,
			Payload:    fmt.Append(nil, 39100),
		})

		res := <-chRes
		eventId := string(res.Payload)
		t.Log("EventId:", eventId)
		assert.NotEmpty(t, eventId)

		//fmt.Println("MSG:", string(b))
	})
}
