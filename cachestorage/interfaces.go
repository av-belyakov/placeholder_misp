package cachestorage

type CacheStorageFuncHandler interface {
	SetFunc(func(int) bool)
	GetFunc() func(int) bool
	Comparison(FormatImplementer) bool
}

type FormatImplementer interface {
	//здесь надо написать все геттеры для объектов
	//MISP перечисленных в listFormatsMISP
	EventSetter
	EventGetter
}

type EventSetter interface {
	SetOrgId(v interface{}, num int)
	SetOrgcId(v interface{}, num int)
	SetDistribution(v interface{}, num int)
	SetInfo(v interface{}, num int)
	SetUUID(v interface{}, num int)
	SetDate(v interface{}, num int)
	SetAnalysis(v interface{}, num int)
	SetAttributeCount(v interface{}, num int)
	SetTimestamp(v interface{}, num int)
	SetSharingGroupId(v interface{}, num int)
	SetThreatLevelId(v interface{}, num int)
	SetPublishTimestamp(v interface{}, num int)
	SetSightingTimestamp(v interface{}, num int)
	SetExtendsUUID(v interface{}, num int)
	SetEventCreatorEmail(v interface{}, num int)
	SetPublished(v interface{}, num int)
	SetProposalEmailLock(v interface{}, num int)
	SetLocked(v interface{}, num int)
	SetDisableCorrelation(v interface{}, num int)
}

type EventGetter interface {
	GetOrgId() string
	GetOrgcId() string
	GetDistribution() string
	GetInfo() string
	GetUUID() string
	GetDate() string
	GetAnalysis() string
	GetAttributeCount() string
	GetTimestamp() string
	GetSharingGroupId() string
	GetThreatLevelId() string
	GetPublishTimestamp() string
	GetSightingTimestamp() string
	GetExtendsUUID() string
	GetEventCreatorEmail() string
	GetPublished() bool
	GetProposalEmailLock() bool
	GetLocked() bool
	GetDisableCorrelation() bool
}
