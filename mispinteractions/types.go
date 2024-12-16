package mispinteractions

import (
	"sync"
)

// ModuleMISP инициализированный модуль
type ModuleMISP struct {
	ChanInput  chan InputSettings //канал для отправки данных В модуль
	ChanOutput chan OutputSetting //канал для отправки данных ИЗ модуля
}

// ChanInputSettings параметры канала для приема данных в модуль
type InputSettings struct {
	CaseId     float64
	Command    string
	TaskId     string
	RootId     string
	CaseSource string
	EventId    string
	UserEmail  string
	MajorData  map[string]interface{}
}

// SettingChanOutputMISP параметры канала для передачи данных из модуля
type OutputSetting struct {
	Command, CaseId, EventId, TaskId, RootId string
}

type StorageAuthorizationData struct {
	AuthList         []UserSettings
	OrganisationList map[string]OrganisationOptions
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
