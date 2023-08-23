package datamodels

import "sync"

type EventsMispFormat struct {
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
	Published          bool   `json:"published"`
	ProposalEmailLock  bool   `json:"proposal_email_lock"`
	Locked             bool   `json:"locked"`
	DisableCorrelation bool   `json:"disable_correlation"`
}

type ListAttributesMispFormat struct {
	attributes []AttributesMispFormat
	mutex      sync.Mutex
}

type AttributesMispFormat struct {
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
	ToIds              bool   `json:"to_ids"`
	Deleted            bool   `json:"deleted"`
	DisableCorrelation bool   `json:"disable_correlation"`
}

type GalaxyClustersMispFormat struct {
	Id             string                    `json:"id"`
	Uuid           string                    `json:"uuid"`
	CollectionUuid string                    `json:"collection_uuid"`
	Type           string                    `json:"type"`
	Value          string                    `json:"value"`
	TagName        string                    `json:"tag_name"`
	Description    string                    `json:"description"`
	GalaxyId       string                    `json:"galaxy_id"`
	Source         string                    `json:"source"`
	Authors        []string                  `json:"authors"`
	Version        string                    `json:"version"`
	Distribution   string                    `json:"distribution"` //цифры в виде строки из списка
	SharingGroupId string                    `json:"sharing_group_id"`
	OrgId          string                    `json:"org_id"`
	OrgcId         string                    `json:"orgc_id"`
	ExtendsUuid    string                    `json:"extends_uuid"`
	ExtendsVersion string                    `json:"extends_version"`
	Default        bool                      `json:"default"`
	Locked         bool                      `json:"locked"`
	Published      bool                      `json:"published"`
	Deleted        bool                      `json:"deleted"`
	GalaxyElement  []GalaxyElementMispFormat `json:"GalaxyElement"`
}

type GalaxyElementMispFormat struct {
	Id              string `json:"id"`
	GalaxyClusterId string `json:"galaxy_cluster_id"`
	Key             string `json:"key"`
	Value           string `json:"value"`
}

type UsersMispFormat struct {
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
	Autoalert     bool   `json:"autoalert"`
	Termsaccepted bool   `json:"termsaccepted"`
	Contactalert  bool   `json:"contactalert"`
	Disabled      bool   `json:"disabled"`
	ForceLogout   bool   `json:"force_logout"`
}

type OrganisationsMispFormat struct {
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
	Local              bool     `json:"local"`
}

type ServersMispFormat struct {
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
}

type FeedsMispFormat struct {
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
	Enabled         bool   `json:"enabled"`
	FixedEvent      bool   `json:"fixed_event"`
	DeltaMerge      bool   `json:"delta_merge"`
	Publish         bool   `json:"publish"`
	OverrideIds     bool   `json:"override_ids"`
	DeleteLocalFile bool   `json:"delete_local_file"`
	LookupVisible   bool   `json:"lookup_visible"`
	CachingEnabled  bool   `json:"caching_enabled"`
	ForceToIds      bool   `json:"force_to_ids"`
}

type TagsMispFormat struct {
	Name           string `json:"name"`
	Colour         string `json:"colour"`
	OrgId          string `json:"org_id"`
	UserId         string `json:"user_id"`
	NumericalValue string `json:"numerical_value"`
	Inherited      int    `json:"inherited"`
	HideTag        bool   `json:"hide_tag"`
	IsGalaxy       bool   `json:"is_galaxy"`
	Exportable     bool   `json:"exportable"`
	IsCustomGalaxy bool   `json:"is_custom_galaxy"`
}

/*
{
  "name": "ORGNAME",
  "date_created": "2021-06-14 14:29:19",
  "date_modified": "2021-06-14 14:29:19",
  "description": "string",
  "type": "ADMIN",
  "nationality": "string",
  "sector": "string",
  "created_by": "12345",
  "uuid": "string",
  "contacts": "string",
  "local": true,
  "restricted_to_domain": [
    "example.com"
  ],
  "landingpage": "string",
  "user_count": "3",
  "created_by_email": "string"
}
*/

/*
То что пришло из MISP при запросе /events/


  {
    {
        "id": "1",
        "org_id": "1",
        "date": "2014-10-02",
        "info": "OSINT ShellShock scanning IPs from OpenDNS",
        "uuid": "542e4c9c-cadc-4f8f-bb11-6d13950d210b",
        "published": true,
        "analysis": "2",
        "attribute_count": "1067",
        "orgc_id": "2",
        "timestamp": "1517817037",
        "distribution": "3",
        "sharing_group_id": "0",
        "proposal_email_lock": false,
        "locked": false,
        "threat_level_id": "3",
        "publish_timestamp": "1615380763",
        "sighting_timestamp": "0",
        "disable_correlation": false,
        "extends_uuid": "",
        "Org": {
          "id": "1",
          "name": "ORGNAME",
          "uuid": "9b912cb0-3079-4c08-83dd-9a58554f0385"
        },
        "Orgc": {
          "id": "2",
          "name": "CthulhuSPRL.be",
          "uuid": "55f6ea5f-fd34-43b8-ac1d-40cb950d210f"
        },
        "EventTag": [
          {
            "id": "1",
            "event_id": "1",
            "tag_id": "1",
            "local": false,
            "Tag": {
              "id": "1",
              "name": "type:OSINT",
              "colour": "#004646",
              "is_galaxy": false
            }
          },
          {
            "id": "2",
            "event_id": "1",
            "tag_id": "2",
            "local": false,
            "Tag": {
              "id": "2",
              "name": "tlp:green",
              "colour": "#33FF00",
              "is_galaxy": false
            }
          },
          {
            "id": "3",
            "event_id": "1",
            "tag_id": "3",
            "local": false,
            "Tag": {
              "id": "3",
              "name": "tlp:white",
              "colour": "#ffffff",
              "is_galaxy": false
            }
          }
        ]
      },
  }
*/
