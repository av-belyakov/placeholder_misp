package cachestorage

type CacheStorageFuncHandler interface {
	SetFunc(func(int) bool)
	GetFunc() func(int) bool
	Comparison(FormatImplementer) bool
}

type FormatImplementer interface {
	//здесь надо написать все геттеры для объектов
	//MISP перечисленных в listFormatsMISP
	EventGetter
	EventSetter
	EventReportsGetter
	EventReportsSetter
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

type EventSetter interface {
	SetOrgId(string)
	SetOrgcId(string)
	SetDistribution(string)
	SetInfo(string)
	SetUUID(string)
	SetDate(string)
	SetAnalysis(string)
	SetAttributeCount(string)
	SetTimestamp(string)
	SetSharingGroupId(string)
	SetThreatLevelId(string)
	SetPublishTimestamp(string)
	SetSightingTimestamp(string)
	SetExtendsUUID(string)
	SetEventCreatorEmail(string)
	SetPublished(bool)
	SetProposalEmailLock(bool)
	SetLocked(bool)
	SetDisableCorrelation(bool)
}

type EventReportsGetter interface {
	GetEventReportsName() string
	GetEventReportsContent() string
	GetEventReportsDistribution() string
}

type EventReportsSetter interface {
	SetEventReportsName(v string)
	SetEventReportsContent(v string)
	SetEventReportsDistribution(v string)
}
