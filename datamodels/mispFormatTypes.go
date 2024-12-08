package datamodels

import "sync"

// описание формата MISP типа Events для загрузки в API MISP
type EventsMispFormat struct {
	Published          bool   `json:"published"`
	ProposalEmailLock  bool   `json:"proposal_email_lock"`
	Locked             bool   `json:"locked"`
	DisableCorrelation bool   `json:"disable_correlation"`
	OrgId              string `json:"org_id"`
	OrgcId             string `json:"orgc_id"`
	Distribution       string `json:"distribution"` //цифры в виде строки из списка
	Info               string `json:"info"`
	Uuid               string `json:"uuid"`
	Date               string `json:"date"`
	Analysis           string `json:"analysis"` //цифры в виде строки из списка
	AttributeCount     string `json:"attribute_count"`
	Timestamp          string `json:"timestamp"`
	SharingGroupId     string `json:"sharing_group_id"`
	ThreatLevelId      string `json:"threat_level_id"`    //цифры в виде строки из списка
	PublishTimestamp   string `json:"publish_timestamp"`  //по умолчанию "0"
	SightingTimestamp  string `json:"sighting_timestamp"` //по умолчанию "0"
	ExtendsUuid        string `json:"extends_uuid"`
	EventCreatorEmail  string `json:"event_creator_email"`
}

type EventReports struct {
	Name         string `json:"name"`
	Distribution string `json:"distribution"`
	Content      string `json:"content"`
}

type ListAttributesMispFormat struct {
	attributes map[int]AttributesMispFormat
	sync.Mutex
}

// описание формата MISP типа Attributes для загрузки в API MISP
type AttributesMispFormat struct {
	ToIds              bool   `json:"to_ids"`
	Deleted            bool   `json:"deleted"`
	DisableCorrelation bool   `json:"disable_correlation"`
	EventId            string `json:"event_id"`
	ObjectId           string `json:"object_id"`
	ObjectRelation     string `json:"object_relation"`
	Category           string `json:"category"` //содержит одно из значений предустановленного списка
	Type               string `json:"type"`     //содержит одно из значений предустановленного списка
	Value              string `json:"value"`
	Uuid               string `json:"uuid"`
	Timestamp          string `json:"timestamp"`    //по умолчанию "0"
	Distribution       string `json:"distribution"` //цифры в виде строки из списка
	SharingGroupId     string `json:"sharing_group_id"`
	Comment            string `json:"comment"`
	FirstSeen          string `json:"first_seen"`
	LastSeen           string `json:"last_seen"`
}

// описание формата MISP типа GalaxyClusters для загрузки в API MISP
type GalaxyClustersMispFormat struct {
	Default        bool                      `json:"default"`
	Locked         bool                      `json:"locked"`
	Published      bool                      `json:"published"`
	Deleted        bool                      `json:"deleted"`
	Id             string                    `json:"id"`
	Uuid           string                    `json:"uuid"`
	CollectionUuid string                    `json:"collection_uuid"`
	Type           string                    `json:"type"`
	Value          string                    `json:"value"`
	TagName        string                    `json:"tag_name"`
	Description    string                    `json:"description"`
	GalaxyId       string                    `json:"galaxy_id"`
	Source         string                    `json:"source"`
	Version        string                    `json:"version"`
	Distribution   string                    `json:"distribution"` //цифры в виде строки из списка
	SharingGroupId string                    `json:"sharing_group_id"`
	OrgId          string                    `json:"org_id"`
	OrgcId         string                    `json:"orgc_id"`
	ExtendsUuid    string                    `json:"extends_uuid"`
	ExtendsVersion string                    `json:"extends_version"`
	Authors        []string                  `json:"authors"`
	GalaxyElement  []GalaxyElementMispFormat `json:"GalaxyElement"`
}

// описание формата MISP типа Galaxy для загрузки в API MISP
type GalaxyElementMispFormat struct {
	Id              string `json:"id"`
	GalaxyClusterId string `json:"galaxy_cluster_id"`
	Key             string `json:"key"`
	Value           string `json:"value"`
}

// описание формата MISP типа Users для загрузки в API MISP
type UsersMispFormat struct {
	Autoalert     bool   `json:"autoalert"`
	Termsaccepted bool   `json:"termsaccepted"`
	Contactalert  bool   `json:"contactalert"`
	Disabled      bool   `json:"disabled"`
	ForceLogout   bool   `json:"force_logout"`
	OrgId         string `json:"org_id"`
	ServerId      string `json:"server_id"`
	Email         string `json:"email"`
	Authkey       string `json:"authkey"`
	InvitedBy     string `json:"invited_by"`
	Gpgkey        string `json:"gpgkey"`
	CertifPublic  string `json:"certif_public"`
	NidsSid       string `json:"nids_sid"`
	Newsread      string `json:"newsread"`
	RoleId        string `json:"role_id"`
	ChangePw      string `json:"change_pw"`
	Expiration    string `json:"expiration"`
	CurrentLogin  string `json:"current_login"`
	LastLogin     string `json:"last_login"`
	DateCreated   string `json:"date_created"`
	DateModified  string `json:"date_modified"`
}

// описание формата MISP типа Organisations для загрузки в API MISP
type OrganisationsMispFormat struct {
	Local              bool     `json:"local"`
	Name               string   `json:"name"`
	DateCreated        string   `json:"date_created"`
	DateModified       string   `json:"date_modified"`
	Description        string   `json:"description"`
	Type               string   `json:"type"`
	Nationality        string   `json:"nationality"`
	Sector             string   `json:"sector"`
	CreatedBy          string   `json:"created_by"`
	Uuid               string   `json:"uuid"`
	Contacts           string   `json:"contacts"`
	Landingpage        string   `json:"landingpage"`
	UserCount          string   `json:"user_count"`
	CreatedByEmail     string   `json:"created_by_email"`
	RestrictedToDomain []string `json:"restricted_to_domain"`
}

// описание формата MISP типа Servers для загрузки в API MISP
type ServersMispFormat struct {
	Push                bool   `json:"push"`
	Pull                bool   `json:"pull"`
	PushSightings       bool   `json:"push_sightings"`
	PushGalaxyClusters  bool   `json:"push_galaxy_clusters"`
	PullGalaxyClusters  bool   `json:"pull_galaxy_clusters"`
	PublishWithoutEmail bool   `json:"publish_without_email"`
	UnpublishEvent      bool   `json:"unpublish_event"`
	SelfSigned          bool   `json:"self_signed"`
	Internal            bool   `json:"internal"`
	SkipProxy           bool   `json:"skip_proxy"`
	CachingEnabled      bool   `json:"caching_enabled"`
	CacheTimestamp      bool   `json:"cache_timestamp"`
	Name                string `json:"name"`
	Url                 string `json:"url"`
	Authkey             string `json:"authkey"`
	OrgId               string `json:"org_id"`
	Lastpulledid        string `json:"lastpulledid"`
	Lastpushedid        string `json:"lastpushedid"`
	Organization        string `json:"organization"`
	RemoteOrgId         string `json:"remote_org_id"`
	PullRules           string `json:"pull_rules"`
	PushRules           string `json:"push_rules"`
	CertFile            string `json:"cert_file"`
	ClientCertFile      string `json:"client_cert_file"`
	Priority            string `json:"priority"`
}

// описание формата MISP типа Feeds для загрузки в API MISP
type FeedsMispFormat struct {
	Enabled         bool   `json:"enabled"`
	FixedEvent      bool   `json:"fixed_event"`
	DeltaMerge      bool   `json:"delta_merge"`
	Publish         bool   `json:"publish"`
	OverrideIds     bool   `json:"override_ids"`
	DeleteLocalFile bool   `json:"delete_local_file"`
	LookupVisible   bool   `json:"lookup_visible"`
	CachingEnabled  bool   `json:"caching_enabled"`
	ForceToIds      bool   `json:"force_to_ids"`
	Name            string `json:"name"`
	Provider        string `json:"provider"`
	Url             string `json:"url"`
	Rules           string `json:"rules"`
	Distribution    string `json:"distribution"` //цифры в виде строки из списка
	SharingGroupId  string `json:"sharing_group_id"`
	TagId           string `json:"tag_id"`
	SourceFormat    string `json:"source_format"`
	EventId         string `json:"event_id"`
	InputSource     string `json:"input_source"`
	Headers         string `json:"headers"`
	OrgcId          string `json:"orgc_id"`
}

// описание формата MISP типа Tags для загрузки в API MISP
type TagsMispFormat struct {
	HideTag        bool   `json:"hide_tag"`
	IsGalaxy       bool   `json:"is_galaxy"`
	Exportable     bool   `json:"exportable"`
	IsCustomGalaxy bool   `json:"is_custom_galaxy"`
	Inherited      int    `json:"inherited"`
	Name           string `json:"name"`
	Colour         string `json:"colour"`
	OrgId          string `json:"org_id"`
	UserId         string `json:"user_id"`
	NumericalValue string `json:"numerical_value"`
}

// EventObjectTagsMispFormat описание формата MISP для загрузки в event.object.tags
type EventObjectTagsMispFormat struct {
	Event string `json:"event"`
	Tag   string `json:"tag"`
}

// описание формата MISP типа Users для данных приходящих из API MISP
// на GET запрос типа /admin/users
type UsersSettingsMispFormat struct {
	User         UserSettingsMispFormat         `json:"User"`
	Organisation OrganisationSettingsMispFormat `json:"Organisation"`
	Role         RoleSettingsMispFormat         `json:"Role"`
}

// описание формата сообщения типа 'User' приходящего от MISP на запрос /admin/users
// так как весь перечень информации о пользователе в настоящее время не нужен
// 'лишние' свойства отключены
type UserSettingsMispFormat struct {
	Id       string `json:"id"`
	OrgId    string `json:"org_id"`
	ServerId string `json:"server_id"`
	Email    string `json:"email"`
	Authkey  string `json:"authkey"`
	//InvitedBy     string `json:"invited_by"`
	//Gpgkey        string `json:"gpgkey"`
	//CertifPublic  string `json:"certif_public"`
	//NidsSid       string `json:"nids_sid"`
	//Newsread      string `json:"newsread"`
	RoleId string `json:"role_id"`
	//Expiration    string `json:"expiration"`
	CurrentLogin string `json:"current_login"`
	//LastLogin     string `json:"last_login"`
	//LastApiAccess string `json:"last_api_access"`
	//DateCreated   string `json:"date_created"`
	//DateModified  string `json:"date_modified"`
	//ChangePw      string `json:"change_pw"`
	//Autoalert     bool   `json:"autoalert"`
	//Termsaccepted bool   `json:"termsaccepted"`
	//Contactalert  bool   `json:"contactalert"`
	//Disabled      bool   `json:"disabled"`
	//ForceLogout   bool   `json:"force_logout"`
}

// описание формата сообщения типа 'Role' приходящего от MISP на запрос /admin/users
type RoleSettingsMispFormat struct {
	PermAuth     bool   `json:"perm_auth"`
	PermSiteAmin bool   `json:"perm_site_admin"`
	Id           string `json:"id"`
	Name         string `json:"name"`
}

// описание формата сообщения типа 'Organisation' приходящего от MISP на запрос /admin/users
type OrganisationSettingsMispFormat struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// описание формата сообщения типа 'Objects' (нет в спецификации API)
// формируется на основе содержимого поля observables получаемого от TheHive
// которое дополнительно содержит поле attachment
type ListObjectsMispFormat struct {
	objects map[int]ObjectsMispFormat
	sync.Mutex
}

type ObjectsMispFormat struct {
	TemplateUUID    string        `json:"template_uuid"`
	TemplateVersion string        `json:"template_version"`
	FirstSeen       string        `json:"first_seen"`
	Timestamp       string        `json:"timestamp"`
	Name            string        `json:"name"`
	Description     string        `json:"description"`
	EventId         string        `json:"event_id"`
	MetaCategory    string        `json:"meta-category"`
	Distribution    string        `json:"distribution"`
	Attribute       ListAttribute `json:"Attribute"`
}

type ListAttribute []AttributeMispFormat

type AttributeMispFormat struct {
	Category       string `json:"category"`
	Type           string `json:"type"`
	Value          string `json:"value"`
	Distribution   string `json:"distribution"`
	ObjectRelation string `json:"object_relation"`
}
