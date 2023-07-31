package elasticsearchinteractions

import (
	"context"
	"fmt"

	"placeholder_misp/confighandler"
	"placeholder_misp/datamodels"
)

var es ModuleElasticSearch

// ModuleElasticSearch инициализированный модуль
// chanInputElasticSearch - канал для отправки данных в модуль
// chanOutputElasticSearch - канал для принятия данных из модуля
// ChanLogging - канал для отправки логов
type ModuleElasticSearch struct {
	chanInputElasticSearch  chan interface{}
	chanOutputElasticSearch chan interface{}
	ChanLogging             chan<- datamodels.MessageLoging
}

func init() {
	es = ModuleElasticSearch{
		chanInputElasticSearch:  make(chan interface{}),
		chanOutputElasticSearch: make(chan interface{}),
	}
}

func NewClientElasticSearch(
	ctx context.Context,
	conf confighandler.AppConfigElasticSearch,
	chanLog chan<- datamodels.MessageLoging) (*ModuleElasticSearch, error) {
	fmt.Println("func 'NewClientElasticSearch', START...")

	es.ChanLogging = chanLog

	/*
		Здесь нужно написать инициализацию подключения
	*/

	return &es, nil
}

func (es ModuleElasticSearch) GetDataReceptionChannel() <-chan interface{} {
	/*
		для формирование правильного сообщения об ошибке
		if err != nil {
			_, f, l, _ := runtime.Caller(0)

			mmisp.ChanLogging <- datamodels.MessageLoging{
				MsgData: fmt.Sprintf("%s %s:%d", fmt.Sprint(err), f, l-2),
				MsgType: "error",
			}
		}
	*/

	return es.chanOutputElasticSearch
}

func (es ModuleElasticSearch) GettingData() interface{} {
	/*
		для формирование правильного сообщения об ошибке
		if err != nil {
			_, f, l, _ := runtime.Caller(0)

			mmisp.ChanLogging <- datamodels.MessageLoging{
				MsgData: fmt.Sprintf("%s %s:%d", fmt.Sprint(err), f, l-2),
				MsgType: "error",
			}
		}
	*/

	return <-es.chanOutputElasticSearch
}

func (es ModuleElasticSearch) SendingData(data interface{}) {
	/*
		для формирование правильного сообщения об ошибке
		if err != nil {
			_, f, l, _ := runtime.Caller(0)

			mmisp.ChanLogging <- datamodels.MessageLoging{
				MsgData: fmt.Sprintf("%s %s:%d", fmt.Sprint(err), f, l-2),
				MsgType: "error",
			}
		}
	*/

	es.chanInputElasticSearch <- data
}
