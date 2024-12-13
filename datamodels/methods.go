package datamodels

func NewListFormatsMISP() *ListFormatsMISP {
	return &ListFormatsMISP{
		Event:      NewEventMisp(),
		Reports:    NewEventReports(),
		Attributes: []*AttributesMispFormat(nil),
		Objects:    map[int]*ObjectsMispFormat(nil),
		ObjectTags: &ListEventObjectTags{},
	}
}

func (o *ListFormatsMISP) GetID() string {
	return o.ID
}

func (o *ListFormatsMISP) GetEvent() *EventsMispFormat {
	return o.Event
}

func (o *ListFormatsMISP) GetReports() *EventReports {
	return o.Reports
}

func (o *ListFormatsMISP) GetAttributes() []*AttributesMispFormat {
	return o.Attributes
}

func (o *ListFormatsMISP) GetObjects() map[int]*ObjectsMispFormat {
	return o.Objects
}

func (o *ListFormatsMISP) GetObjectTags() *ListEventObjectTags {
	return o.ObjectTags
}

func (o *ListFormatsMISP) ComparisonID(v string) bool {
	return o.ID == v
}

func (o *ListFormatsMISP) ComparisonEvent(v *EventsMispFormat) bool {
	if o.Event.Analysis != v.Analysis {
		return false
	}

	if o.Event.AttributeCount != v.AttributeCount {
		return false
	}

	// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	//и т.д. для всех свойств
	/*
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
	*/

	return true
}

func (o *ListFormatsMISP) ComparisonReports(v *EventReports) bool {
	/*
		type EventReports struct {
		Name         string `json:"name"`
		Content      string `json:"content"`
		Distribution string `json:"distribution"`
	}
	*/

	return true
}

func (o *ListFormatsMISP) ComparisonAttributes(v []*AttributesMispFormat) bool {
	/*
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
	*/

	return true
}

func (o *ListFormatsMISP) ComparisonObjects(v map[int]*ObjectsMispFormat) bool {
	/*
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
	*/

	return true
}

func (o *ListFormatsMISP) ComparisonObjectTags(v *ListEventObjectTags) bool {
	//type ListEventObjectTags []string

	return true
}
