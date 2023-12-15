package mispinteractions

import (
	"sync"
)

// ModuleMISP инициализированный модуль
// chanInputMISP - канал для отправки данных в модуль
// chanOutputMISP - канал для отправки данных из модуля
type ModuleMISP struct {
	ChanInputMISP  chan SettingsChanInputMISP
	ChanOutputMISP chan SettingChanOutputMISP
}

type SettingsChanInputMISP struct {
	Command    string
	TaskId     string
	CaseId     float64
	CaseSource string
	EventId    string
	UserEmail  string
	MajorData  map[string]interface{}
}

type SettingChanOutputMISP struct {
	Command, CaseId, EventId, TaskId string
}

func (nmisp ModuleMISP) TestSend() {}

func (mmisp ModuleMISP) GetDataReceptionChannel() <-chan SettingChanOutputMISP {
	return mmisp.ChanOutputMISP
}

func (mmisp ModuleMISP) SendingDataOutput(data SettingChanOutputMISP) {
	mmisp.ChanOutputMISP <- data
}

func (mmisp ModuleMISP) GetInputChannel() <-chan SettingsChanInputMISP {
	return mmisp.ChanInputMISP
}

func (mmisp ModuleMISP) SendingDataInput(data SettingsChanInputMISP) {
	mmisp.ChanInputMISP <- data
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
