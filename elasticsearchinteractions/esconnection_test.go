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
		//es     *elasticsearch.TypedClient
		es *elasticsearch.Client
	)

	BeforeAll(func() {
		esConf = elasticsearch.Config{
			Addresses: []string{"http://es-siem-db.cloud.gcm:9200" /*"http://datahook.cloud.gcm:9200"*/},
			Username:  "elasticsearch",
			//Username: "elastic",
			Password: "1zex3TvB",
			//Password: "iG99g3lyHsazTucx8eOL",
		}

		fmt.Println(esConf)
		//es, errEs = elasticsearch.NewTypedClient(esConf)
		es, errEs = elasticsearch.NewClient(esConf)
		//els, err := elasticsearch.NewTypedClient(esConf)

		es.Get("thehive33", "test")
	})

	Context("Тест 1. Проверка подключения к Elasticsearch", func() {
		It("Должно быть успешно установлено подключение к БД", func() {

			fmt.Println("Elasticsearch client ERROR: ", errEs)
			fmt.Println("Elasticsearch client: ", es)

			Expect(errEs).ShouldNot(HaveOccurred())
		})
	})
})
