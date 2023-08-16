package elasticsearchinteractions_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	//"placeholder_misp/elasticsearchinteractions"
)

var _ = Describe("Esconnection", Ordered, func() {

	var (
		errEs  error
		esConf elasticsearch.Config
		es     *elasticsearch.Client
	)

	BeforeAll(func() {
		esConf = elasticsearch.Config{
			Addresses: []string{"datahook.cloud.gcm"},
			Username:  "elasticsearch",
			Password:  "1zex3TvB",
		}

		es, errEs = elasticsearch.NewClient(esConf)
		//els, err := elasticsearch.NewTypedClient(esConf)
		es.Get("thehive33", "test")
	})

	Context("Тест 1. Проверка подклюбчения к Elasticsearch", func() {
		It("Должно быть успешно установлено подключение к БД", func() {

			fmt.Println("Elasticsearch client ERROR: ", errEs)
			fmt.Println("Elasticsearch client: ")

			Expect(errEs).ShouldNot(HaveOccurred())
		})
	})
})
