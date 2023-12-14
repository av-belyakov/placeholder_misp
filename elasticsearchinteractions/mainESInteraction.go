package elasticsearchinteractions

import (
	"bytes"
	"fmt"
	"net/http"
	"placeholder_misp/confighandler"
	"placeholder_misp/datamodels"
	"placeholder_misp/memorytemporarystorage"
	"runtime"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
)

var es ModuleElasticSearch

type SettingsInputChan struct {
	UUID string
}

// ModuleElasticSearch инициализированный модуль
// chanInputElasticSearch - канал для отправки данных в модуль
// chanOutputElasticSearch - канал для принятия данных из модуля
type ModuleElasticSearch struct {
	chanInputElasticSearch  chan SettingsInputChan
	chanOutputElasticSearch chan interface{}
}

type handlerSendData struct {
	conf       confighandler.AppConfigElasticSearch
	storageApp *memorytemporarystorage.CommonStorageTemporary
	logging    chan<- datamodels.MessageLogging
}

func (hsd handlerSendData) sendingData(uuid string) {
	defer hsd.storageApp.SetIsProcessedElasticsearshHiveFormatMessage(uuid)

	if !hsd.conf.Send {
		return
	}

	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{fmt.Sprintf("http://%s:%d", hsd.conf.Host, hsd.conf.Port)},
		Username:  hsd.conf.User,
		Password:  hsd.conf.Passwd,
	})
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		hsd.logging <- datamodels.MessageLogging{
			MsgData: fmt.Sprintf("'%s' %s:%d", err.Error(), f, l-6),
			MsgType: "error",
		}

		return
	}

	b, ok := hsd.storageApp.GetRawDataHiveFormatMessage(uuid)
	if !ok {
		_, f, l, _ := runtime.Caller(0)
		hsd.logging <- datamodels.MessageLogging{
			MsgData: fmt.Sprintf("'the raw data was not found in the repository' %s:%d", f, l-2),
			MsgType: "warning",
		}

		return
	}

	t := time.Now()
	buf := bytes.NewReader(b)
	res, err := es.API.Index(fmt.Sprintf("%s%s_%d_%d", hsd.conf.Prefix, hsd.conf.Index, t.Year(), int(t.Month())), buf)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		hsd.logging <- datamodels.MessageLogging{
			MsgData: fmt.Sprintf("'%s' %s:%d", err.Error(), f, l-2),
			MsgType: "error",
		}

		return
	}

	if res.StatusCode == http.StatusCreated || res.StatusCode == http.StatusOK {
		return
	}

	hsd.logging <- datamodels.MessageLogging{
		MsgData: fmt.Sprintf("received from module Elsaticsearch: %d %s", res.StatusCode, res.Status()),
		MsgType: "warning",
	}
}

func init() {
	es = ModuleElasticSearch{
		chanInputElasticSearch:  make(chan SettingsInputChan),
		chanOutputElasticSearch: make(chan interface{}),
	}
}

func HandlerElasticSearch(
	conf confighandler.AppConfigElasticSearch,
	storageApp *memorytemporarystorage.CommonStorageTemporary,
	logging chan<- datamodels.MessageLogging) *ModuleElasticSearch {

	hsd := handlerSendData{conf, storageApp, logging}

	go func() {
		for data := range es.chanInputElasticSearch {
			go hsd.sendingData(data.UUID)
		}
	}()

	return &es
}

func (es ModuleElasticSearch) HandlerData(data SettingsInputChan) {
	es.chanInputElasticSearch <- data
}

/*
func (es ModuleElasticSearch) GetDataReceptionChannel() <-chan interface{} {
	return es.chanOutputElasticSearch
}

func (es ModuleElasticSearch) GettingData() interface{} {
	return <-es.chanOutputElasticSearch
}
*/
