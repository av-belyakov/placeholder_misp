package mispapi

import "net/http"

type ModuleMispHandler interface {
	GetDataReceptionChannel() <-chan OutputSetting
	SendingDataOutput(OutputSetting)
	SendingDataInput(InputSettings)
}

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
	GetAuthData() string
}
