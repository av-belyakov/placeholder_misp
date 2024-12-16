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
	performComparison := func(current, added AttributesMispFormat) bool {
		if current.ToIds != added.ToIds {
			return false
		}

		if current.Deleted != added.Deleted {
			return false
		}

		if current.DisableCorrelation != added.DisableCorrelation {
			return false
		}

		if current.EventId != added.EventId {
			return false
		}

		if current.ObjectRelation != added.ObjectRelation {
			return false
		}

		if current.Category != added.Category {
			return false
		}

		if current.Type != added.Type {
			return false
		}

		if current.Value != added.Value {
			return false
		}

		if current.Distribution != added.Distribution {
			return false
		}
		if current.SharingGroupId != added.SharingGroupId {
			return false
		}
		if current.Comment != added.Comment {
			return false
		}

		//Uuid `json:"uuid"` не стал так как для каждого объекта через конструктор
		//автоматически формируется свой идентификатор
		//
		//Timestamp `json:"timestamp"`
		//FirstSeen `json:"first_seen"`
		//LastSeen `json:"last_seen"`
		//а с временем вообще не ясно что и когда может поменять TheHive

		return true
	}

	if len(o.Attributes) != len(v) {
		return false
	}

	for _, currentAttribute := range o.GetAttributes() {
		var isExist bool
		for _, addedAttribute := range v {
			if currentAttribute.ObjectId == addedAttribute.ObjectId {
				isExist = true

				if !performComparison(*currentAttribute, *addedAttribute) {
					return false
				}
			}
		}

		if !isExist {
			return false
		}
	}

	return true
}

// ComparisonObjects выполняет сравнение свойств объекта Objects
func (o *ListFormatsMISP) ComparisonObjects(v map[int]*ObjectsMispFormat) bool {
	performComparison := func(current, added ObjectsMispFormat) bool {
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

		// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		// зесь надо дописать сравнение для каждого свойства
		// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

		return true
	}

	if len(o.Objects) != len(v) {
		return false
	}

	for _, currentObject := range o.GetObjects() {
		var isExist bool
		for _, addedObject := range v {
			if currentObject.ID == addedObject.ID {
				isExist = true

				if !performComparison(*currentObject, *addedObject) {
					return false
				}
			}
		}

		if !isExist {
			return false
		}
	}

	return true
}

// ComparisonObjectTags выполняет сравнение свойств объекта ObjectTags
func (o *ListFormatsMISP) ComparisonObjectTags(v *ListEventObjectTags) bool {
	if len(*o.ObjectTags) != len(v.GetListTags()) {
		return false
	}

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
