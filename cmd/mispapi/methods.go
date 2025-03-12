package mispapi

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/av-belyakov/placeholder_misp/internal/confighandler"
	"github.com/av-belyakov/placeholder_misp/internal/datamodels"
	"github.com/av-belyakov/placeholder_misp/internal/supportingfunctions"
	"golang.org/x/net/context"
)

//****** каналы *******
//

func (m ModuleMISP) GetReceptionChannel() <-chan OutputSetting {
	return m.chOutput
}

func (m ModuleMISP) SendDataOutput(data OutputSetting) {
	m.chOutput <- data
}

func (m ModuleMISP) GetInputChannel() <-chan InputSettings {
	return m.chInput
}

func (m ModuleMISP) SendDataInput(data InputSettings) {
	m.chInput <- data
}

//****** хранилище пользовательских данных и настроек организаций ******
//

// setUserSettings добавляет настройки пользователя в хранилище, если
// пользователь id или email уже есть, ничего не делает
func (s *StorageAuthorizationData) setUserSettings(us UserSettings) bool {
	s.Lock()
	defer s.Unlock()

	for _, v := range s.AuthList {
		idIsExist := us.UserId == v.UserId
		emailIsExist := us.Email == v.Email

		if idIsExist || emailIsExist {
			return false
		}
	}

	s.AuthList = append(s.AuthList, us)

	return true
}

// GetUserSettingsByEmail получает настройки пользователя по его email
func (s *StorageAuthorizationData) GetUserSettingsByEmail(email string) (*UserSettings, bool) {
	s.Lock()
	defer s.Unlock()

	for _, v := range s.AuthList {
		if v.Email == email {
			return &v, true
		}
	}

	return nil, false
}

// GetSettingsAllUsers возвращает настройки всех пользователей
func (s *StorageAuthorizationData) GetSettingsAllUsers() []UserSettings {
	return s.AuthList
}

// CleanUsers удаляет из памяти данные о всех пользователях
func (s *StorageAuthorizationData) cleanUsers() {
	s.Lock()
	defer s.Unlock()

	s.AuthList = []UserSettings{}
}

// setOrganisationOptions добавляет информацию об организации в хранилище
// key - наименование источника (свойство 'source' в json TheHive)
// values - первое значение id организации, второе наименование организации (
// соответствующие поля в MISP)
func (s *StorageAuthorizationData) setOrganisationOptions(key string, values [2]string) {
	s.Lock()
	defer s.Unlock()

	s.OrganisationList[key] = OrganisationOptions{Id: values[0], Name: values[1]}
}

// GetOrganisationOptions возвращает информацию об организации по ее наименованию
func (s *StorageAuthorizationData) GetOrganisationOptions(key string) (*OrganisationOptions, bool) {
	s.Lock()
	defer s.Unlock()

	if v, ok := s.OrganisationList[key]; ok {
		return &v, ok
	}

	return &OrganisationOptions{}, false
}

// GetOptionsAllOrganisations возвращает опции по всем организациям
func (s *StorageAuthorizationData) GetOptionsAllOrganisations() map[string]OrganisationOptions {
	return s.OrganisationList
}

//****** управление авторизационными данными *******
//

func (ad *AuthorizationDataMISP) getListAllUsers(ctx context.Context) ([]datamodels.UsersSettingsMispFormat, error) {
	usmispf := []datamodels.UsersSettingsMispFormat{}
	_, resByte, err := ad.Get(ctx, "/admin/users", nil)
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
func (ad *AuthorizationDataMISP) GetListAllUsers(ctx context.Context) (int, error) {
	var countUser int
	lus, err := ad.getListAllUsers(ctx)
	if err != nil {
		return countUser, supportingfunctions.CustomError(err)
	}

	//очищаем хранилище с данными пользователей
	ad.Storage.cleanUsers()

	for _, v := range lus {
		countUser++
		ad.Storage.setUserSettings(UserSettings{
			UserId:  v.User.Id,
			OrgId:   v.Organisation.Id,
			OrgName: v.Organisation.Name,
			Email:   v.User.Email,
			AuthKey: v.User.Authkey,
			Role:    v.Role.Name,
		})
	}

	return countUser, nil
}

// GetUserData проверяет наличие пользователя в памяти, если его там нет то
// делает запрос к MISP и обновляет данные в памяти если что то пришло
func (ad *AuthorizationDataMISP) GetUserData(ctx context.Context, user string) (UserSettings, error) {
	if us, ok := ad.Storage.GetUserSettingsByEmail(user); ok {
		return *us, nil
	}

	lus, err := ad.getListAllUsers(ctx)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return UserSettings{}, fmt.Errorf("'%s' %s:%d", err.Error(), f, l-2)
	}

	for _, v := range lus {
		if _, ok := ad.Storage.GetUserSettingsByEmail(v.User.Email); ok {
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

		ad.Storage.setUserSettings(us)

		if user == v.User.Email {
			return us, nil
		}
	}

	return UserSettings{}, fmt.Errorf("information about the user '%s' was not found", user)
}

// только для теста
func (ad *AuthorizationDataMISP) DelUserData(user string) (int, bool) {
	ad.Storage.Lock()
	defer ad.Storage.Unlock()

	var (
		num     int
		isExist bool
	)

	newList := []UserSettings{}
	for k, v := range ad.Storage.AuthList {
		if v.Email == user {
			num = k
			isExist = true

			newList = append(ad.Storage.AuthList[:k], ad.Storage.AuthList[k+1:]...)
			break
		}
	}

	ad.Storage.AuthList = newList

	return num, isExist
}

// DeleteUser удаление пользователя из MISP
func (ad *AuthorizationDataMISP) DeleteUser(ctx context.Context, userId string) error {
	_, _, err := ad.Delete(ctx, fmt.Sprintf("/admin/users/delete/%s", userId))
	if err != nil {
		return err
	}

	return nil
}

// GetListAllOrganisation получает список всех организаций добавленых в MISP
// коррелирует с данными из конфигурационного файла и сохраняет в хранилище
func (ad *AuthorizationDataMISP) GetListAllOrganisation(ctx context.Context, confOrg []confighandler.Organization) error {
	_, resByte, err := ad.Get(ctx, "/organisations", nil)
	if err != nil {
		return err
	}

	orgs := RecivedOrganisations{}
	if err := json.Unmarshal(resByte, &orgs); err != nil {
		_, f, l, _ := runtime.Caller(0)
		return fmt.Errorf("'%s' %s:%d", err.Error(), f, l-2)
	}

	for _, v := range orgs {
		for _, value := range confOrg {
			if v.Organisation.Name == value.OrgName {
				ad.Storage.setOrganisationOptions(value.SourceName, [2]string{v.Organisation.Id, v.Organisation.Name})
			}
		}
	}

	return nil
}

// CreateNewUser создает нового пользователя в MISP и добавляет его данные в хранилище
// source name:
//
//	"gcm": "ГЦМ (г.Москва)", "GCM"
//	"rcmsr": "РЦМ (г. Симферополь)", "CR-RCM"
//	"rcmlnx": "РЦМ (г. Смоленск)", "SMOL-RCM"
//	"rcmros": "РЦМ (г. Ростов-на-дону)", "UFO-RCM"
//	"rcmkgd": "РЦМ (г. Калининград)", "KGD-RCM"
//	"rcmspb": "РЦМ (г. Санкт-Петербург)", "SZFO-RCM"
//	"rcmsve": "РЦМ (г. Екатеринбург)", "UralFO-RCM"
//	"rcmniz": "РЦМ (г. Нижний Новгород)", "PFO-RCM"
//	"rcmsta": "РЦМ (г. Ставрополь)", "SKFO-RCM"
//	"rcmnvs": "РЦМ (г. Новосибирск)", "SFO-RCM"
//	"rcmkha": "РЦМ (г. Хабаровск)", "DFO-RCM"
//	"rcmmsk": "РЦМ (г. Москва и МО)", "CFO-RCM"
func (ad *AuthorizationDataMISP) CreateNewUser(ctx context.Context, email, source string) (UserSettings, error) {
	orgId := "1"
	if org, ok := ad.Storage.GetOrganisationOptions(source); ok {
		orgId = org.Id
	}

	b, err := json.Marshal(struct {
		Email  string `json:"email"`
		OrgId  string `json:"org_id"`
		RoleId string `json:"role_id"`
	}{
		Email:  email,
		OrgId:  orgId,
		RoleId: "3", //3 это роль группы 'User' в MISP
	})
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return UserSettings{}, fmt.Errorf("'%s' %s:%d", err.Error(), f, l-2)
	}

	_, resByte, err := ad.Post(ctx, "/admin/users/add", b)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return UserSettings{}, fmt.Errorf("'%s' %s:%d", err.Error(), f, l-2)
	}

	usmispf := datamodels.UsersSettingsMispFormat{}
	if err := json.Unmarshal(resByte, &usmispf); err != nil {
		_, f, l, _ := runtime.Caller(0)
		return UserSettings{}, fmt.Errorf("'%s' %s:%d", err.Error(), f, l-2)
	}

	newUser := UserSettings{
		UserId:  usmispf.User.Id,
		OrgId:   usmispf.Organisation.Id,
		OrgName: usmispf.Organisation.Name,
		Email:   usmispf.User.Email,
		AuthKey: usmispf.User.Authkey,
		Role:    usmispf.Role.Name,
	}

	_ = ad.Storage.setUserSettings(newUser)

	return newUser, nil
}

//******* методы вспомогательного типа, используемого для кэша ********

// NewCacheSpecialObject конструктор вспомогательного типа реализующий интерфейс CacheStorageFuncHandler[T any]
func NewCacheSpecialObject[T SpecialObject]() *CacheSpecialObject[T] {
	return &CacheSpecialObject[T]{}
}

func (o *CacheSpecialObject[T]) SetID(v string) {
	o.id = v
}

func (o *CacheSpecialObject[T]) GetID() string {
	return o.id
}

func (o *CacheSpecialObject[T]) SetObject(v T) {
	o.object = v
}

func (o *CacheSpecialObject[T]) GetObject() T {
	return o.object
}

func (o *CacheSpecialObject[T]) SetFunc(f func(int) bool) {
	o.handlerFunc = f
}

func (o *CacheSpecialObject[T]) GetFunc() func(int) bool {
	return o.handlerFunc
}

func (o *CacheSpecialObject[T]) Comparison(objFromCache T) bool {
	if !o.object.ComparisonID(objFromCache.GetID()) {
		return false
	}

	if !o.object.ComparisonEvent(objFromCache.GetEvent()) {
		return false
	}

	if !o.object.ComparisonReports(objFromCache.GetReports()) {
		return false
	}

	if !o.object.ComparisonAttributes(objFromCache.GetAttributes()) {
		return false

	}

	if !o.object.ComparisonObjects(objFromCache.GetObjects()) {
		return false
	}

	if !o.object.ComparisonObjectTags(objFromCache.GetObjectTags()) {
		return false
	}

	return true
}

func (o *CacheSpecialObject[T]) MatchingAndReplacement(objFromCache T) T {

	/*
		MatchingAndReplacementEvents(v objectsmispformat.EventsMispFormat) objectsmispformat.EventsMispFormat
		MatchingAndReplacementReport(v objectsmispformat.EventReports) objectsmispformat.EventReports
		MatchingAndReplacementAttributes(v []*objectsmispformat.AttributesMispFormat) []*objectsmispformat.AttributesMispFormat
		MatchingAndReplacementObjects(v map[int]*objectsmispformat.ObjectsMispFormat) map[int]*objectsmispformat.ObjectsMispFormat
		MatchingAndReplacementListEventObjectTags(v objectsmispformat.ListEventObjectTags) objectsmispformat.ListEventObjectTags
	*/
	objFromCache.SetEvent(o.object.MatchingAndReplacementEvents(*objFromCache.GetEvent()))
	objFromCache.SetReports(o.object.MatchingAndReplacementReport(*objFromCache.GetReports()))
	objFromCache.SetAttributes(o.object.MatchingAndReplacementAttributes(objFromCache.GetAttributes()))
	objFromCache.SetObjects(o.object.MatchingAndReplacementObjects(objFromCache.GetObjects()))
	objFromCache.SetObjectTags(o.object.MatchingAndReplacementListEventObjectTags(*objFromCache.GetObjectTags()))

	return objFromCache
}
