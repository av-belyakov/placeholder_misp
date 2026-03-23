package mispapi

import (
	"sync"

	"github.com/av-belyakov/cachingstoragewithqueue"
	"github.com/av-belyakov/objectsmispformat"
	"github.com/av-belyakov/placeholder_misp/commoninterfaces"
	"github.com/av-belyakov/placeholder_misp/internal/confighandler"
)

// ModuleMISP инициализированный модуль
type ModuleMISP struct {
	cache        *cachingstoragewithqueue.CacheStorageWithQueue[*objectsmispformat.ListFormatsMISP]
	organistions []confighandler.Organization
	host         string
	authKey      string
	logger       commoninterfaces.Logger
	chInput      chan InputSettings //канал для отправки данных В модуль
	chOutput     chan OutputSetting //канал для отправки данных ИЗ модуля
}

// ChanInputSettings параметры канала для приема данных в модуль
type InputSettings struct {
	//MajorData  map[string]interface{}
	Data       objectsmispformat.ListFormatsMISP
	Command    string
	TaskId     string
	RootId     string
	CaseSource string
	EventId    string
	UserEmail  string
	CaseId     float64
}

// SettingChanOutputMISP параметры канала для передачи данных из модуля
type OutputSetting struct {
	Command, CaseId, EventId, TaskId, RootId, CaseSource string
}

// StorageAuthorizationData хранилище с авторизационными настройками пользователя
type StorageAuthorizationData struct {
	AuthList         []UserSettings
	OrganisationList map[string]OrganisationOptions
	sync.Mutex
}

// UserSettings авторизационныt настройки пользователя
type UserSettings struct {
	UserId  string
	OrgId   string
	Email   string
	AuthKey string
	OrgName string
	Role    string
}

// OrganisationOptions yfcnhjqrb jhufybpfwbb
type OrganisationOptions struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type RecivedOrganisations []RecivedOrganisation

type RecivedOrganisation struct {
	Organisation OrganisationOptions
}

type AuthorizationDataMISP struct {
	ConnectMISPHandler
	Storage *StorageAuthorizationData
}

type MispResponse struct {
	Event map[string]interface{} `json:"event"`
}

type requestMISP struct {
	host          string
	userAuthKey   string
	masterAuthKey string
}

type RequestMISPOptions func(*requestMISP) error

// CacheSpecialObject специальный объект соответствующий интерфейсу cachingstoragewithqueue.CacheStorageHandler
type CacheSpecialObject[T SpecialObject] struct {
	object      T
	handlerFunc func(int) bool
	id          string
}

type ResponseTagsMISPFormat []ResponseTagMISPFormat

type ResponseTagMISPFormat struct {
	Tag TagSettings `json:"tag"`
}

type TagSettings struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Colour         string `json:"colour"`
	OrgId          string `json:"org_id"`
	UserId         string `json:"user_id"`
	NumericalValue any    `json:"numerical_value"`
	HideTag        bool   `json:"hide_tag"`
	IsGalaxy       bool   `json:"is_galaxy"`
	Exportable     bool   `json:"exportable"`
	LocalOnly      bool   `json:"local_only"`
	IsCustomGalaxy bool   `json:"is_custom_galaxy"`
}
