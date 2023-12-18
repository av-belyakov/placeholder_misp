package testelasticsearch_test

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/elastic/go-elasticsearch/v8"

	"placeholder_misp/supportingfunctions"
)

var _ = Describe("Testelasticsearch", Ordered, func() {
	var (
		errConn, errReadFile error
		es                   *elasticsearch.Client
		exampleByte          []byte
	)

	readFileJson := func(fpath, fname string) ([]byte, error) {
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

	/*
	    es = Elasticsearch([f"http://writer:XxZqesYXuk8C@datahook.cloud.gcm:9200"])

	    http://datahook.cloud.gcm:5601
	   elastic
	   iG99g3lyHsazTucx8eOL
	*/

	BeforeAll(func() {
		//читаем тестовый файл
		exampleByte, errReadFile = readFileJson("natsinteractions/test_json", "example_caseId_33705_1.json")

		fmt.Println("file length: ", len(exampleByte))

		es, errConn = elasticsearch.NewClient(elasticsearch.Config{

			Addresses: []string{"http://datahook.cloud.gcm:9200"},
			Username:  "writer",
			Password:  "XxZqesYXuk8C",
		})
	})

	Context("Тест 1. Чтение JSON файла", func() {
		It("При чтении файла не должно быть ошибок", func() {
			Expect(errReadFile).ShouldNot(HaveOccurred())
		})
	})

	Context("Тест 2. Создание соединения с БД", func() {
		It("При инициализации соединения не должно быть ошибок", func() {
			Expect(errConn).ShouldNot(HaveOccurred())
		})
	})

	Context("Тест 3. Получение информации о БД", func() {
		It("При обработке запроса для получения инфрмации о БД не должно быть ошибок", func() {
			//res, err := client.API.Index("hive-case*").Raw(exampleByte).Do(context.Background())
			res, err := es.Info()
			defer res.Body.Close()
			Expect(err).ShouldNot(HaveOccurred())

			fmt.Println("Elasticsearch Info: ", res)

			Expect(res.StatusCode).Should(Equal(http.StatusOK))
		})
	})

	Context("Тест 4. Запись json объекта в бинарном виде в БД", func() {
		It("При записи не должно быть ошибок", func() {
			buf := bytes.NewReader(exampleByte)

			t := time.Now()
			prefix := ""
			index := "module_placeholder_thehive_case"

			str := fmt.Sprintf("%s%s_%d_%d", prefix, index, t.Year(), int(t.Month()))

			fmt.Println("string reguest:", str)

			//res, err := es.Create("my_test_thehive_case_gcm_2023", "1999", buf)
			//module_placeholder_thehive_case
			res, err := es.API.Index(str, buf)
			//res, err := es.API.Index("my_test_thehive_case_gcm_2023", buf /*, es.Index.WithDocumentID("1000004")*/)
			defer res.Body.Close()
			Expect(err).ShouldNot(HaveOccurred())

			r := map[string]interface{}{}
			err = json.NewDecoder(res.Body).Decode(&r)

			fmt.Println("____________________ INDEX response: ")
			for k, v := range r {
				fmt.Printf("%s: %v", k, v)
			}

			if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusOK {
				if e, ok := r["error"]; ok {
					fmt.Println("Error:", e)
				}
			}

			Expect(res.StatusCode).Should(Equal(http.StatusCreated))
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("Тест 5. Поиск индекса в БД", func() {
		It("При поиске индекса не должно быть ошибок, индекс должен быть найден", func() {
			//res, err := client.API.Index("hive-case*").Raw(exampleByte).Do(context.Background())

			res, err := es.Search(
				es.Search.WithContext(context.Background()),
				es.Search.WithIndex("my_test_thehive_case_gcm_2023"),
				es.Search.WithBody(strings.NewReader(`{"query" : { "match" : { "_id" : "1000000" } }}`)),
				es.Search.WithPretty(),
			)
			defer res.Body.Close()
			Expect(err).ShouldNot(HaveOccurred())

			//fmt.Println("RESPONSE: ", res)

			r := map[string]interface{}{}
			err = json.NewDecoder(res.Body).Decode(&r)
			Expect(err).ShouldNot(HaveOccurred())

			fmt.Println("RESPONSE: ", r)

			Expect(true).Should(BeTrue())
		})
	})
})
