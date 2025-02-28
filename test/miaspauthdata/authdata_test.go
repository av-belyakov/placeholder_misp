package testmiaspauthdata_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/av-belyakov/placeholder_misp/cmd/mispapi"
	"github.com/av-belyakov/placeholder_misp/internal/confighandler"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var (
	client      *mispapi.ClientMISP
	handlerAuth *mispapi.AuthorizationDataMISP

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

	err error
)

const Test_User_Email string = "my.testuser@exampleemail.org"

func TestMain(m *testing.M) {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalln(err)
	}

	client, err = mispapi.NewClientMISP("misp-center.cloud.gcm", os.Getenv("GO_PHMISP_MAUTH"), false)
	if err != nil {
		log.Fatalln(err)
	}

	handlerAuth = mispapi.NewHandlerAuthorizationMISP(client, mispapi.NewStorageAuthorizationDataMISP())

	os.Exit(m.Run())
}

func TestAuthData(t *testing.T) {
	t.Run("Тест 1. Получаем список всех организаций.", func(t *testing.T) {
		ctx, CancelFunc := context.WithTimeout(context.Background(), time.Second*10)
		defer CancelFunc()

		err := handlerAuth.GetListAllOrganisation(ctx, confOrgs)
		assert.NoError(t, err)

		newListOrg := handlerAuth.Storage.GetOptionsAllOrganisations()
		fmt.Printf("Test 1.\nOrganization list:\n")
		for k, v := range newListOrg {
			fmt.Printf("Key:'%s' Id:'%s' Name:'%s'\n", k, v.Id, v.Name)
		}
		assert.Equal(t, len(newListOrg), 12)
	})

	t.Run("Тест 2. Получаем список всех пользователей.", func(t *testing.T) {
		ctx, CancelFunc := context.WithTimeout(context.Background(), time.Second*10)
		defer CancelFunc()

		countUsers, err := handlerAuth.GetListAllUsers(ctx)
		assert.NoError(t, err)

		fmt.Printf("\nTest 2.\nIt was added %d users\n", countUsers)

		usersSettings := handlerAuth.Storage.GetSettingsAllUsers()
		fmt.Println("\nSettings all users:")
		var num int
		for ; num < len(usersSettings); num++ {
			if num == 3 {
				break
			}

			fmt.Printf("user Id:%s, role:%s, org. name:%s\n", usersSettings[num].UserId, usersSettings[num].Role, usersSettings[num].OrgName)
		}
		fmt.Printf("and more %d users...\n\n", len(usersSettings)-num)
	})

	t.Run("Тест 3. Поиск ID и наименование организаций по её псевдониму.", func(t *testing.T) {
		res, ok := handlerAuth.Storage.GetOrganisationOptions("rcmsta")
		assert.True(t, ok)
		assert.Equal(t, res.Id, "12")
		assert.Equal(t, res.Name, "SKFO-RCM")

		res, ok = handlerAuth.Storage.GetOrganisationOptions("rcmkgd")
		assert.True(t, ok)
		assert.Equal(t, res.Id, "10")
		assert.Equal(t, res.Name, "KGD-RCM")

		_, ok = handlerAuth.Storage.GetOrganisationOptions("rcmtest")
		assert.False(t, ok)
	})

	t.Run("Тест 4. Проверяем наличие данных о пользователе.", func(t *testing.T) {
		ctx, CancelFunc := context.WithTimeout(context.Background(), time.Second*10)
		defer CancelFunc()

		us, err := handlerAuth.GetUserData(ctx, "a.shershnev@cloud.gcm")
		assert.NoError(t, err)
		assert.Equal(t, us.Email, "a.shershnev@cloud.gcm")

		us, err = handlerAuth.GetUserData(ctx, "user.test@example.email")
		assert.Error(t, err)
	})

	t.Run("Тест 5. Создание нового пользователя", func(t *testing.T) {
		ctx, CancelFunc := context.WithTimeout(context.Background(), time.Second*10)
		defer CancelFunc()

		userset, err := handlerAuth.CreateNewUser(ctx, Test_User_Email, "rcmkha")
		assert.NoError(t, err)
		assert.Equal(t, userset.Email, Test_User_Email)

		fmt.Printf("\nTest 5.\nCreated new user:%+v\n", userset)

		//поиск в хранилище
		nus, ok := handlerAuth.Storage.GetUserSettingsByEmail(Test_User_Email)
		assert.True(t, ok)
		assert.Equal(t, nus.Email, Test_User_Email)
	})

	t.Run("Тест 6. Удаление только созданного нового пользователя", func(t *testing.T) {
		ctx, CancelFunc := context.WithTimeout(context.Background(), time.Second*10)
		defer CancelFunc()

		//удаляем пользователя с хранилища
		num, ok := handlerAuth.DelUserData(Test_User_Email)
		assert.True(t, ok)
		fmt.Printf("Test 6.\nUSER num = %d\n", num)

		// ищем пользователя в хранилище
		user, ok := handlerAuth.Storage.GetUserSettingsByEmail(Test_User_Email)
		fmt.Println("user not found", user)
		assert.False(t, ok)

		// получаем данные пользователя включая его userId
		userSettings, err := handlerAuth.GetUserData(ctx, Test_User_Email)
		assert.NoError(t, err)

		//удаляем пользователя из MISP
		err = handlerAuth.DeleteUser(ctx, userSettings.UserId)
		assert.NoError(t, err)

		// ищем пользователя в хранилище и MISP
		userSettings, err = handlerAuth.GetUserData(ctx, Test_User_Email)
		fmt.Printf("user settings:%+v\n", userSettings)
		assert.NoError(t, err)
		assert.Equal(t, userSettings.Email, Test_User_Email)
	})
}
