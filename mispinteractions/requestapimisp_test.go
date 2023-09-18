package mispinteractions_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"placeholder_misp/datamodels"
	"placeholder_misp/memorytemporarystorage"
	"placeholder_misp/mispinteractions"
)

var _ = Describe("Requestapimisp", Ordered, func() {
	var (
		errUnmar, errHttpRes error
		errClientMisp        error
		httpRes              *http.Response
		resByte              []byte
		clientMisp           mispinteractions.ClientMISP
		filePath             string
		listTmp              []interface{}
		usmispf              []datamodels.UsersSettingsMispFormat
		mts                  *memorytemporarystorage.CommonStorageTemporary
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
		mts = memorytemporarystorage.NewTemporaryStorage()
	})

	Context("Тест 1. Выполняем запрос к API MISP с целью полчения информации по объекту Attributes", func() {
		BeforeEach(func() {
			httpRes, resByte, errHttpRes = clientMisp.Get("/attributes", []byte{})
		})

		It("При выполнении запроса типа attributes ошибки быть не должно", func() {
			Expect(errClientMisp).ShouldNot(HaveOccurred())
			Expect(httpRes.StatusCode).Should(Equal(http.StatusOK))
			Expect(errHttpRes).ShouldNot(HaveOccurred())

			//fmt.Printf("RESPONSE: \nHead:%v\n", httpRes.Header)
		})

		It("При выполнении записи в файл ошибки быть не должно", func() {
			errWriteFile := os.WriteFile(filePath, resByte, 0666)
			errJson := json.Unmarshal(resByte, &listTmp)

			fmt.Printf("\nCount attributes: %d\n", len(listTmp))

			Expect(errWriteFile).ShouldNot(HaveOccurred())
			Expect(errJson).ShouldNot(HaveOccurred())
		})
	})

	Context("Тест 2. Выполняем запрос к API MISP с целью получения списка авторизационных ключей", func() {
		BeforeEach(func() {
			httpRes, resByte, errHttpRes = clientMisp.Get("/admin/users", []byte{})
			usmispf = []datamodels.UsersSettingsMispFormat{}
			errUnmar = json.Unmarshal(resByte, &usmispf)
		})

		It("При выполнении запроса к API MISP ошибок быть не должно и должно быть получено определенной количество авторизационных ключей", func() {

			//fmt.Println("http body: ", string(resByte))

			Expect(httpRes.StatusCode).Should(Equal(http.StatusOK))
			Expect(errHttpRes).ShouldNot(HaveOccurred())
			Expect(errUnmar).ShouldNot(HaveOccurred())

			fmt.Println("COUNT USERS = ", len(usmispf))
			for k, v := range usmispf {
				fmt.Printf("%d. userEmail:%s\n", k, v.User.Email)
			}

			Expect(len(usmispf)).Should(Equal(83))
		})

		It("После заполнения временного хранилища должен быть найден пользователь по его email", func() {
			for _, v := range usmispf {
				mts.AddUserSettingsMISP(memorytemporarystorage.UserSettingsMISP{
					UserId:  v.User.Id,
					OrgId:   v.Organisation.Id,
					OrgName: v.Organisation.Name,
					Email:   v.User.Email,
					AuthKey: v.User.Authkey,
					Role:    v.Role.Name,
				})
			}

			ts, ok := mts.GetUserSettingsMISP("a.makarov@cloud.gcm")

			fmt.Println("_________________________")
			fmt.Println("Find user: ", ts)

			Expect(ok).Should(BeTrue())
			Expect(ts.AuthKey).Should(Equal("yoEJPVwXDGXELzOVuel9HueXxkrMrtZ9UoJGkh4y"))

			Expect(len(mts.ListUserSettingsMISP)).Should(Equal(78))
		})
	})
})
