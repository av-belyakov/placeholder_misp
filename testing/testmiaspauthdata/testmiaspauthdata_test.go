package testmiaspauthdata_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"placeholder_misp/datamodels"
	"runtime"
	"sync"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type ConnectMISPHandler interface {
	NetworkSender
	SetterAuthData
}

type NetworkSender interface {
	Get(path string, data []byte) (*http.Response, []byte, error)
	Post(path string, data []byte) (*http.Response, []byte, error)
	Delete(path string) (*http.Response, []byte, error)
}

type SetterAuthData interface {
	SetAuthData(ah string)
}

type StorageAuthorizationData struct {
	authList []UserSettings
	sync.Mutex
}

type UserSettings struct {
	UserId  string
	OrgId   string
	Email   string
	AuthKey string
	OrgName string
	Role    string
}

type ClientMISP struct {
	BaseURL  *url.URL
	Host     string
	AuthHash string
	Verify   bool
}

// SetUserSettings добавляет настройки пользователя в хранилище, если
// пользователь id или email уже есть, ничего не делает
func (s *StorageAuthorizationData) SetUserSettings(us UserSettings) bool {
	s.Lock()
	defer s.Unlock()

	for _, v := range s.authList {
		idIsExist := us.UserId == v.UserId
		emailIsExist := us.Email == v.Email

		if idIsExist || emailIsExist {
			return false
		}
	}

	s.authList = append(s.authList, us)

	return true
}

// GetUserSettingsByEmail получает настройки пользователя по его email
func (s *StorageAuthorizationData) GetUserSettingsByEmail(email string) (*UserSettings, bool) {
	s.Lock()
	defer s.Unlock()

	for _, v := range s.authList {
		if v.Email == email {
			return &v, true
		}
	}

	return nil, false
}

// CleanUsers удаляет из памяти данные о всех пользователях
func (s *StorageAuthorizationData) CleanUsers() {
	s.Lock()
	defer s.Unlock()

	s.authList = []UserSettings{}
}

// NewStorageAuthorizationDataMISP создает новое хранилище с данными пользователей MISP
func NewStorageAuthorizationDataMISP() *StorageAuthorizationData {
	return &StorageAuthorizationData{}
}

type AuthorizationDataMISP struct {
	ConnectMISPHandler
	storage *StorageAuthorizationData
}

func NewHandlerAuthorizationMISP(c ConnectMISPHandler, s *StorageAuthorizationData) *AuthorizationDataMISP {
	return &AuthorizationDataMISP{
		c,
		s,
	}
}

func (ad *AuthorizationDataMISP) getListAllUsers() ([]datamodels.UsersSettingsMispFormat, error) {
	usmispf := []datamodels.UsersSettingsMispFormat{}
	_, resByte, err := ad.Get("/admin/users", nil)
	if err != nil {
		return usmispf, err
	}

	err = json.Unmarshal(resByte, &usmispf)
	if err != nil {
		return usmispf, err
	}

	return usmispf, nil
}

// GetListAllUsers получает список всех пользователей которые есть в MISP
// и добавляет их в хранилище. Внимание, перед добавлением списка пользователей,
// данные о пользователях которые оставались в хранилище будут удалены. Но только
// если при обращении к MISP не будет ошибок
func (ad *AuthorizationDataMISP) GetListAllUsers() error {
	lus, err := ad.getListAllUsers()
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return fmt.Errorf("'%s' %s:%d", err.Error(), f, l-2)
	}

	//очищаем хранилище с данными пользователей
	ad.storage.CleanUsers()

	for _, v := range lus {
		ad.storage.SetUserSettings(UserSettings{
			UserId:  v.User.Id,
			OrgId:   v.Organisation.Id,
			OrgName: v.Organisation.Name,
			Email:   v.User.Email,
			AuthKey: v.User.Authkey,
			Role:    v.Role.Name,
		})
	}

	return nil
}

// GetUserData проверяет наличие пользователя в памяти, если его там нет то
// делает запрос к MISP и обновляет данные в памяти если что то пришло
func (ad *AuthorizationDataMISP) GetUserData(user string) (UserSettings, error) {
	if us, ok := ad.storage.GetUserSettingsByEmail(user); ok {
		return *us, nil
	}

	lus, err := ad.getListAllUsers()
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return UserSettings{}, fmt.Errorf("'%s' %s:%d", err.Error(), f, l-2)
	}

	for _, v := range lus {
		if _, ok := ad.storage.GetUserSettingsByEmail(v.User.Email); ok {
			continue
		}

		us := UserSettings{
			UserId:  v.User.Id,
			OrgId:   v.Organisation.Id,
			OrgName: v.Organisation.Name,
			Email:   v.User.Email,
			AuthKey: v.User.Authkey,
			Role:    v.Role.Name,
		}

		ad.storage.SetUserSettings(us)

		if user == v.User.Email {
			return us, nil
		}
	}

	return UserSettings{}, fmt.Errorf("information about the user '%s' was not found", user)
}

/*
"gcm": "RU-MOW",
            "rcmlnx": "RU-SMO",
            "rcmros": "RU-ROS",
            "rcmkgd": "RU-KGD",
            "rcmspb": "RU-SPE",
            "rcmsve": "RU-SVE",
            "rcmniz": "RU-NIZ",
            "rcmsr": "RU-CR",
            "rcmsta": "RU-STA",
            "rcmnvs": "RU-NVS",
            "rcmkha": "RU-KHA",
            "rcmmsk": "RU-MOW"

"gcm": "Москва",
            "rcmlnx": "Смоленск",
            "rcmros": "Ростов-на-дону",
            "rcmkgd": "Калининград",
            "rcmspb": "Санкт-Петербург",
            "rcmsve": "Екатеринбург",
            "rcmniz": "Нижний Новгород",
            "rcmsr": "Симферополь",
            "rcmsta": "Ставрополь",
            "rcmnvs": "Новосибирск",
            "rcmkha": "Хабаровск",
            "rcmmsk": "Москва"

"gcm": "ГЦМ (г.Москва)",
            "rcmlnx": "РЦМ (г. Смоленск)",
            "rcmros": "РЦМ (г. Ростов-на-дону)",
            "rcmkgd": "РЦМ (г. Калининград)",
            "rcmspb": "РЦМ (г. Санкт-Петербург)",
            "rcmsve": "РЦМ (г. Екатеринбург)",
            "rcmniz": "РЦМ (г. Нижний Новгород)",
            "rcmsr": "РЦМ (г. Симферополь)",
            "rcmsta": "РЦМ (г. Ставрополь)",
            "rcmnvs": "РЦМ (г. Новосибирск)",
            "rcmkha": "РЦМ (г. Хабаровск)",
            "rcmmsk": "РЦМ (г. Москва и МО)"

*/
//CreateNewUser создает нового пользователя в MISP и добавляет его данные в хранилище
func (ad *AuthorizationDataMISP) CreateNewUser(email, source string) (UserSettings, error) {
	b, err := json.Marshal(struct {
		Email  string `json:"email"`
		OrgId  string `json:"org_id"`
		RoleId string `json:"role_id"`
	}{
		Email:  email,
		OrgId:  "",  //тут по sourceId надо найти наименование организации
		RoleId: "3", //3 это роль User в MISP
	})
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return UserSettings{}, fmt.Errorf("'%s' %s:%d", err.Error(), f, l-2)
	}

	_, resByte, err := ad.Post("/admin/users/add", b)
}

/*
1.+ получить список всех пользователей при старте приложения
2.+ проверять, по запросу, наличие пользователя с таким email
2.1.+ если пользователь есть в памяти то возвращать его авторизационные данные
2.2.+- запрашивать пользователя у MISP
2.2.1. если пользователя нет, то выполнять запрос на создание нового
 и добавлять этого пользователя в память (с учетом поля source,
 принадлежность к определенному региону)
2.2.2.+ если пользователь есть в MISP то возвращать его данные и записывать
 в память хранилища
*/

/*



	client, err := mispinteractions.NewClientMISP(host, auth, false)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return &sadmisp, fmt.Errorf("'%s' %s:%d", err.Error(), f, l-2)
	}

	_, resByte, err := client.Get("/admin/users", nil)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return &sadmisp, fmt.Errorf("'%s' %s:%d", err.Error(), f, l-2)
	}

	usmispf := []datamodels.UsersSettingsMispFormat{}
	err = json.Unmarshal(resByte, &usmispf)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return &sadmisp, fmt.Errorf("'%s' %s:%d", err.Error(), f, l-2)
	}

	for _, v := range usmispf {
		sadmisp.authList = append(sadmisp.authList, UserSettingsMISP{
			UserId:  v.User.Id,
			OrgId:   v.Organisation.Id,
			OrgName: v.Organisation.Name,
			Email:   v.User.Email,
			AuthKey: v.User.Authkey,
			Role:    v.Role.Name,
		})
	}*/

var _ = Describe("Testmiaspauthdata", Ordered, func() {
	BeforeAll(func() {

	})

	Context("", func() {
		It("", func() {

		})
	})

	/*
		Context("", func(){
			It("", func(){

			})
		})
	*/
})
