package mispapi

import (
	"context"
	"net/http"

	"github.com/av-belyakov/objectsmispformat"
)

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
	Get(ctx context.Context, path string, data []byte) (*http.Response, []byte, error)
	Post(ctx context.Context, path string, data []byte) (*http.Response, []byte, error)
	Delete(ctx context.Context, path string) (*http.Response, []byte, error)
}

type SetterAuthData interface {
	SetAuthData(ah string)
	GetAuthData() string
}

type SpecialObjectComparator interface {
	ComparisonID(string) bool
	ComparisonEvent(*objectsmispformat.EventsMispFormat) bool
	ComparisonReports(*objectsmispformat.EventReports) bool
	ComparisonAttributes([]*objectsmispformat.AttributesMispFormat) bool
	ComparisonObjects(map[int]*objectsmispformat.ObjectsMispFormat) bool
	ComparisonObjectTags(*objectsmispformat.ListEventObjectTags) bool
	SpecialObjectGetter
}

type SpecialObjectGetter interface {
	GetID() string
	GetEvent() *objectsmispformat.EventsMispFormat
	GetReports() *objectsmispformat.EventReports
	GetAttributes() []*objectsmispformat.AttributesMispFormat
	GetObjects() map[int]*objectsmispformat.ObjectsMispFormat
	GetObjectTags() *objectsmispformat.ListEventObjectTags
}
