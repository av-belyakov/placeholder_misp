package testmiaspauthdata_test

import (
	"fmt"
	"log"
	"net/url"
	"placeholder_misp/confighandler"
	"placeholder_misp/mispinteractions"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// NewClientMISP возвращает структуру типа ClientMISP с предустановленными значениями
func NewClientMISP(h, a string, v bool) (*mispinteractions.ClientMISP, error) {
	urlBase, err := url.Parse("http://" + h)
	if err != nil {
		return &mispinteractions.ClientMISP{}, err
	}

	return &mispinteractions.ClientMISP{
		BaseURL:  urlBase,
		Host:     h,
		AuthHash: a,
		Verify:   v,
	}, nil
}

// NewStorageAuthorizationDataMISP создает новое хранилище с данными пользователей MISP
func NewStorageAuthorizationDataMISP() *mispinteractions.StorageAuthorizationData {
	return &mispinteractions.StorageAuthorizationData{
		AuthList:         []mispinteractions.UserSettings{},
		OrganisationList: map[string]mispinteractions.OrganisationOptions{},
	}
}

// NewHandlerAuthorizationMISP создает новый обработчик соединений с MISP
func NewHandlerAuthorizationMISP(c mispinteractions.ConnectMISPHandler, s *mispinteractions.StorageAuthorizationData) *mispinteractions.AuthorizationDataMISP {
	return &mispinteractions.AuthorizationDataMISP{
		c,
		s,
	}
}

var _ = Describe("Testmiaspauthdata", Ordered, func() {
	var (
		confOrgs        []confighandler.Organization
		handler         *mispinteractions.AuthorizationDataMISP
		errGetListOrgs  error
		errGetListUsers error
		countUser       int
	)

	BeforeAll(func() {
		confOrgs = []confighandler.Organization{
			{OrgName: "GCM", SourceName: "gcm"},
			{OrgName: "CFO-RCM", SourceName: "rcmmsk"},
			{OrgName: "DFO-RCM", SourceName: "rcmkha"},
			{OrgName: "SFO-RCM", SourceName: "rcmnvs"},
			{OrgName: "SKFO-RCM", SourceName: "rcmsta"},
			{OrgName: "CR-RCM", SourceName: "rcmsr"},
			{OrgName: "PFO-RCM", SourceName: "rcmniz"},
			{OrgName: "UralFO-RCM", SourceName: "rcmsve"},
			{OrgName: "SZFO-RCM", SourceName: "rcmspb"},
			{OrgName: "KGD-RCM", SourceName: "rcmkgd"},
			{OrgName: "UFO-RCM", SourceName: "rcmros"},
			{OrgName: "SMOL-RCM", SourceName: "rcmlnx"},
		}

		client, err := NewClientMISP("misp-center.cloud.gcm", "Z2PwRBdP5lFP7rdDJBzxmSahaLEwIvJoeOuwhRYQ", false)
		if err != nil {
			//_, f, l, _ := runtime.Caller(0)
			//return fmt.Errorf("'%s' %s:%d", err.Error(), f, l-2)
			log.Println(err)
		}

		handler = NewHandlerAuthorizationMISP(client, NewStorageAuthorizationDataMISP())

		errGetListOrgs = handler.GetListAllOrganisation(confOrgs)
		countUser, errGetListUsers = handler.GetListAllUsers()
	})

	Context("Тест 1. Получаем список всех организаций", func() {
		It("При выполнения данного действия не должно быть ошибок, список организаций должен быть из 12 шт.", func() {
			Expect(errGetListOrgs).ShouldNot(HaveOccurred())

			orgList := handler.Storage.GetOptionsAllOrganisations()

			//fmt.Println("__________________________________________")
			//for k, v := range orgList {
			//	fmt.Printf("Key: %s\n  Id: %s\n  Name: %s\n", k, v.Id, v.Name)
			//}

			Expect(len(orgList)).Should(Equal(12))
		})
	})

	Context("Тест 2. Получаем список всех пользователей", func() {
		It("При выполнения данного действия не должно быть ошибок, список пользователей должен быть больше 0", func() {
			Expect(errGetListUsers).ShouldNot(HaveOccurred())
			Expect(countUser).Should(Equal(108))

			//userList := handler.Storage.GetSettingsAllUsers()

			//for k, v := range userList {
			//	fmt.Printf("%d.\n%v\n", k, v)
			//}

			//fmt.Println("____________________ Total users found ______________________")
			//fmt.Println("Users num:", len(userList))
			//fmt.Println("_____________________________________________________________")

			//Expect(len(userList)).Should(BeNumerically(">", 0))
		})
	})

	Context("Тест 3. Получить id и наименование организации по ее псевдониму", func() {
		It("Организация ДОЛЖНА быть найдена", func() {
			res, ok := handler.Storage.GetOrganisationOptions("rcmsta")
			Expect(ok).Should(BeTrue())
			Expect(res.Id).Should(Equal("12"))
			Expect(res.Name).Should(Equal("SKFO-RCM"))

			res, ok = handler.Storage.GetOrganisationOptions("rcmkgd")
			Expect(ok).Should(BeTrue())
			Expect(res.Id).Should(Equal("10"))
			Expect(res.Name).Should(Equal("KGD-RCM"))
		})

		It("Организация НЕ ДОЛЖНА быть найдена", func() {
			_, ok := handler.Storage.GetOrganisationOptions("rcmtest")
			Expect(ok).Should(BeFalse())
		})
	})

	Context("Тест 4.Проверяем наличие данных о пользователе", func() {
		It("В хранилище НЕ ДОЛЖЕН быть найден несуществующий пользователь", func() {
			_, ok := handler.Storage.GetUserSettingsByEmail("feiigg@dffe.h")
			Expect(ok).Should(BeFalse())
		})
		It("В хранилище ДОЛЖЕН быть найден несуществующий пользователь", func() {
			_, ok := handler.Storage.GetUserSettingsByEmail("a.shershnev@cloud.gcm")
			Expect(ok).Should(BeTrue())
		})
		It("Должен быть найден пользователь существующий в MISP", func() {
			us, err := handler.GetUserData("a.shershnev@cloud.gcm")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(us.Email).Should(Equal("a.shershnev@cloud.gcm"))
		})
		It("Неизвестный пользователь не должен быть найден", func() {
			us, err := handler.GetUserData("user.test@example.email")
			Expect(err).Should(HaveOccurred())
			Expect(us.Email).Should(Equal(""))
		})
	})
	Context("Тест 5. Создание нового пользователя в MISP при его отсутствии", func() {
		/*It("Должен быть успешно создан новый пользователь", func() {

			//	При создании нового пользователя нужно сделать запись в лог-файл info

			userEmail := "my.testuser@exampleemail.org"

			userset, err := handler.CreateNewUser(userEmail, "rcmkha")
			Expect(err).ShouldNot(HaveOccurred())
			fmt.Println("userset.AuthKey:", userset.AuthKey)
			Expect(userset.Email).Should(Equal(userEmail))

			nus, ok := handler.Storage.GetUserSettingsByEmail(userEmail)
			Expect(ok).Should(BeTrue())
			Expect(nus.Email).Should(Equal(userEmail))
		})
		It("Должен быть успешно создан новый пользователь со стандартной организацией", func() {
			userEmail := "new.User@exampleemail.org"

			userset, err := handler.CreateNewUser(userEmail, "ff")
			Expect(err).ShouldNot(HaveOccurred())
			fmt.Println("userset.AuthKey:", userset.AuthKey)
			Expect(userset.Email).Should(Equal(userEmail))
		})*/
		It("Должен быть успешно найден, в памяти, определенный пользователь", func() {
			userEmail := "pukucheryaviy@spbfsb.ru"

			us, err := handler.GetUserData(userEmail)

			fmt.Println("DATA for user email", us)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(us.Email).Should(Equal(userEmail))
		})

		It("Информация о пользователе должна быть успешно удалена из хранилища приложения", func() {
			userEmail := "pukucheryaviy@spbfsb.ru"

			//удаляем пользователя
			num, ok := handler.DelUserData(userEmail)
			Expect(ok).Should(BeTrue())
			fmt.Println("USER num = ", num)

			//ищем пользователя в хранилище
			us, ok := handler.Storage.GetUserSettingsByEmail(userEmail)
			fmt.Println("user not found", us)
			Expect(ok).ShouldNot(BeTrue())

			//ищем пользователя в хранилище и MISP
			s, err := handler.GetUserData(userEmail)
			fmt.Println("us", s)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(s.Email).Should(Equal(userEmail))

			//получаем всех пользователей из хранилища
			userList := handler.Storage.GetSettingsAllUsers()
			Expect(len(userList)).Should(Equal(108))

			//ищем УЖЕ находящегося в хранилище пользователя в хранилище и MISP
			s, err = handler.GetUserData(userEmail)
			fmt.Println("----------- us", s)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(s.Email).Should(Equal(userEmail))
		})

		/*
					It("Должен быть успешно найден в MISP, так как в памяти его нет, определенный пользователь", func ()  {
				//pukucheryaviy@spbfsb.ru
			})
		*/
	})

	/*
		Context("", func(){
			It("", func(){

			})
		})
	*/
})
