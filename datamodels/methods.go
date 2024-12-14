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

// GetID возвращает уникальный идентификатор
func (o *ListFormatsMISP) GetID() string {
	return o.ID
}

// GetEvent возвращает объект Event
func (o *ListFormatsMISP) GetEvent() *EventsMispFormat {
	return o.Event
}

// GetReports возвращает объект Reports
func (o *ListFormatsMISP) GetReports() *EventReports {
	return o.Reports
}

// GetAttributes возвращает объект Attributes
func (o *ListFormatsMISP) GetAttributes() []*AttributesMispFormat {
	return o.Attributes
}

// GetObjects возвращает объект Objects
func (o *ListFormatsMISP) GetObjects() map[int]*ObjectsMispFormat {
	return o.Objects
}

// GetObjectTags возвращает объект ObjectTags
func (o *ListFormatsMISP) GetObjectTags() *ListEventObjectTags {
	return o.ObjectTags
}

// ComparisonID выполняет сравнение уникальных идентификаторов
func (o *ListFormatsMISP) ComparisonID(v string) bool {
	return o.ID == v
}

// ComparisonEvent выполняет сравнение свойств объекта Event
func (o *ListFormatsMISP) ComparisonEvent(v *EventsMispFormat) bool {
	if o.Event.Analysis != v.Analysis {
		return false
	}

	if o.Event.Analysis != v.Analysis {
		return false
	}

	if o.Event.AttributeCount != v.AttributeCount {
		return false
	}

	if o.Event.OrgId != v.OrgId {
		return false
	}

	if o.Event.OrgcId != v.OrgcId {
		return false
	}

	if o.Event.Distribution != v.Distribution {
		return false
	}

	if o.Event.Info != v.Info {
		return false
	}

	if o.Event.Uuid != v.Uuid {
		return false
	}

	if o.Event.Date != v.Date {
		return false
	}

	if o.Event.SharingGroupId != v.SharingGroupId {
		return false
	}

	if o.Event.ThreatLevelId != v.ThreatLevelId {
		return false
	}

	if o.Event.ExtendsUuid != v.ExtendsUuid {
		return false
	}

	if o.Event.EventCreatorEmail != v.EventCreatorEmail {
		return false
	}

	if o.Event.Published != v.Published {
		return false
	}

	if o.Event.ProposalEmailLock != v.ProposalEmailLock {
		return false
	}

	if o.Event.Locked != v.Locked {
		return false
	}

	if o.Event.DisableCorrelation != v.DisableCorrelation {
		return false
	}

	// думаю время сравнивать не стоит, потому что большая вероятность получить идентичный
	//во всех параметрах объект у которого будет отличатся только время, что в данном случае не очень важно
	//Timestamp
	//PublishTimestamp
	//SightingTimestamp

	return true
}

// ComparisonReports выполняет сравнение свойств объекта Reports
func (o *ListFormatsMISP) ComparisonReports(v *EventReports) bool {
	if o.Reports.Name != v.Name {
		return false
	}

	if o.Reports.Content != v.Content {
		return false
	}

	if o.Reports.Distribution != v.Distribution {
		return false
	}

	return true
}

// ComparisonAttributes выполняет сравнение свойств объекта Attributes
func (o *ListFormatsMISP) ComparisonAttributes(v []*AttributesMispFormat) bool {

	/*
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
		Distribution       string `json:"distribution"` //цифры в виде строки из списка
		SharingGroupId     string `json:"sharing_group_id"`
		Comment            string `json:"comment"`
		FirstSeen          string `json:"first_seen"` //время
		LastSeen           string `json:"last_seen"` //время
		Timestamp          string `json:"timestamp"`    //по умолчанию "0"
	*/

	return true
}

// ComparisonObjects выполняет сравнение свойств объекта Objects
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

// ComparisonObjectTags выполняет сравнение свойств объекта ObjectTags
func (o *ListFormatsMISP) ComparisonObjectTags(v *ListEventObjectTags) bool {
	for _, currentObject := range *o.ObjectTags {
		var isEqual bool
		for _, objectAdded := range v.GetListTags() {
			if currentObject == objectAdded {
				isEqual = true

				break
			}
		}

		if !isEqual {
			return false
		}
	}

	return true
}
