package mispinteractions_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"placeholder_misp/mispinteractions"
)

var _ = Describe("Requestapimisp", Ordered, func() {
	var (
		errJson, errRes, errWriteFile, errReadAll, errClientMisp error
		clientMisp                                               mispinteractions.ClientMISP
		res                                                      *http.Response
		resByte                                                  []byte
		filePath                                                 string
		listTmp                                                  []interface{}
	)

	getRootPath := func(rootDir string) (string, error) {
		currentDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			return "", err
		}

		tmp := strings.Split(currentDir, "/")

		if tmp[len(tmp)-1] == rootDir {
			return currentDir, nil
		}

		var path string = ""
		for _, v := range tmp {
			path += v + "/"

			if v == rootDir {
				return path, nil
			}
		}

		return path, nil
	}

	rootPath, _ := getRootPath("placeholder_misp")
	filePath = path.Join(rootPath, "logs", "list-Attributes.txt")

	BeforeAll(func() {
		clientMisp, errClientMisp = mispinteractions.NewClientMISP("192.168.9.37", "Z2PwRBdP5lFP7rdDJBzxmSahaLEwIvJoeOuwhRYQ", false)
		res, errRes = clientMisp.Get("/attributes", []byte{})
		resByte, errReadAll = io.ReadAll(res.Body)

		fmt.Println("Root path: ", rootPath)

		errWriteFile = os.WriteFile(filePath, resByte, 0666)

		errJson = json.Unmarshal(resByte, &listTmp)
	})

	Context("Тест 1. Выполняем запрос к API MISP с целью полчения информации по объекту Attributes", func() {
		It("При выполнении запроса типа attributes ошибки быть не должно", func() {
			Expect(errClientMisp).ShouldNot(HaveOccurred())
			Expect(errRes).ShouldNot(HaveOccurred())
			Expect(errReadAll).ShouldNot(HaveOccurred())

			fmt.Printf("RESPONSE: \nHead:%v\n", res.Header)
		})

		It("При выполнении записи в файл ошибки быть не должно", func() {

			fmt.Printf("\nCount attributes: %d\n", len(listTmp))

			Expect(errWriteFile).ShouldNot(HaveOccurred())
			Expect(errJson).ShouldNot(HaveOccurred())
		})
	})
})
