package redisapi_test

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"testing"
	"time"

	"github.com/av-belyakov/placeholder_misp/cmd/redisapi"
	"github.com/av-belyakov/placeholder_misp/commoninterfaces"
	"github.com/av-belyakov/placeholder_misp/constants"
	"github.com/av-belyakov/placeholder_misp/internal/logginghandler"
	"github.com/av-belyakov/placeholder_misp/internal/supportingfunctions"
	"github.com/av-belyakov/simplelogger"
	"github.com/stretchr/testify/assert"
)

const (
	Host string = "127.0.0.1"
	Port int    = 6379
)

var (
	module *redisapi.ModuleRedis

	chZabbix chan commoninterfaces.Messager = make(chan commoninterfaces.Messager)
	logging  *logginghandler.LoggingChan
)

func readFileJson(fpath, fname string) ([]byte, error) {
	var newResult []byte

	rootPath, err := supportingfunctions.GetRootPath("placeholder_misp")
	if err != nil {
		return newResult, err
	}

	//fmt.Println("func 'readFileJson', path = ", path.Join(rootPath, fpath, fname))

	f, err := os.OpenFile(path.Join(rootPath, fpath, fname), os.O_RDONLY, os.ModePerm)
	if err != nil {
		return newResult, err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		newResult = append(newResult, sc.Bytes()...)
	}

	return newResult, nil
}

func TestMain(m *testing.M) {
	simpleLogger, err := simplelogger.NewSimpleLogger(context.Background(), constants.Root_Dir, []simplelogger.Options{})
	if err != nil {
		log.Fatalf("error module 'simplelogger': %v", err)
	}

	chZabbix = make(chan commoninterfaces.Messager)
	logging = logginghandler.New(simpleLogger, chZabbix)
	//logging.Start(ctx)

	os.Exit(m.Run())
}

func TestModuleRedis(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(chZabbix)
				logging.Close()

			case log := <-logging.GetChan():
				fmt.Println("LOGGING: ", log)

			case msg := <-chZabbix:
				t.Logf("counting: type = %s, message = %s\n", msg.GetType(), msg.GetMessage())

			}
		}
	}()

	module = redisapi.NewModuleRedis(Host, Port, logging)
	assert.NoError(t, module.Start(ctx))

	t.Run("Тест 1. Добавление значений", func(t *testing.T) {
		ch := module.GetReceptionChannel()

		//добавляем значение
		module.SendDataInput(redisapi.SettingsInput{
			Command: "set caseId",
			Data:    "12003:789", //caseId:eventId
		})

		//ищем значение
		module.SendDataInput(redisapi.SettingsInput{
			Command: "search caseId",
			Data:    "12003",
		})

		info := <-ch

		assert.Equal(t, info.CommandResult, "found caseId")

		strRes, ok := info.Result.(string)
		assert.True(t, ok)
		assert.Equal(t, strRes, "112340")
	})
}
