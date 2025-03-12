package mispapi

import (
	"context"
	"net/http"

	"github.com/av-belyakov/objectsmispformat"
)

type ModuleMispHandler interface {
	GetReceptionChannel() <-chan OutputSetting
	SendDataOutput(OutputSetting)
	SendDataInput(InputSettings)
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

type SpecialObject interface {
	SpecialObjectComparator
	SpecialObjectMatchingAndReplacement
}

type SpecialObjectComparator interface {
	SpecialObjectGetter
	ComparisonID(string) bool
	ComparisonEvent(*objectsmispformat.EventsMispFormat) bool
	ComparisonReports(*objectsmispformat.EventReports) bool
	ComparisonAttributes([]*objectsmispformat.AttributesMispFormat) bool
	ComparisonObjects(map[int]*objectsmispformat.ObjectsMispFormat) bool
	ComparisonObjectTags(*objectsmispformat.ListEventObjectTags) bool
}

type SpecialObjectMatchingAndReplacement interface {
	SpecialObjectGetter
	MatchingAndReplacementEvents(v objectsmispformat.EventsMispFormat) objectsmispformat.EventsMispFormat
	MatchingAndReplacementReport(v objectsmispformat.EventReports) objectsmispformat.EventReports
	MatchingAndReplacementAttributes(v []*objectsmispformat.AttributesMispFormat) []*objectsmispformat.AttributesMispFormat
	MatchingAndReplacementObjects(v map[int]*objectsmispformat.ObjectsMispFormat) map[int]*objectsmispformat.ObjectsMispFormat
	MatchingAndReplacementListEventObjectTags(v objectsmispformat.ListEventObjectTags) objectsmispformat.ListEventObjectTags
}

type SpecialObjectGetter interface {
	GetID() string
	GetEvent() *objectsmispformat.EventsMispFormat
	GetReports() *objectsmispformat.EventReports
	GetAttributes() []*objectsmispformat.AttributesMispFormat
	GetObjects() map[int]*objectsmispformat.ObjectsMispFormat
	GetObjectTags() *objectsmispformat.ListEventObjectTags
}
