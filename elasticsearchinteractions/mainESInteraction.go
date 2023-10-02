package elasticsearchinteractions

import (
	"placeholder_misp/confighandler"
	"placeholder_misp/datamodels"
	"placeholder_misp/memorytemporarystorage"
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

func init() {
	es = ModuleElasticSearch{
		chanInputElasticSearch:  make(chan SettingsInputChan),
		chanOutputElasticSearch: make(chan interface{}),
	}
}

func HandlerElasticSearch(
	conf confighandler.AppConfigElasticSearch,
	storageApp *memorytemporarystorage.CommonStorageTemporary,
	loging chan<- datamodels.MessageLoging) (*ModuleElasticSearch, error) {

	go func() {
		for data := range es.chanInputElasticSearch {
			//Так как пока никакой обработки нет просто устанавливаем
			//статус для хранилища как 'обработанный модулем'

			storageApp.SetIsProcessedElasticsearshHiveFormatMessage(data.UUID)
		}
	}()

	return &es, nil
}

func (es ModuleElasticSearch) SendingData(data SettingsInputChan) {
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
