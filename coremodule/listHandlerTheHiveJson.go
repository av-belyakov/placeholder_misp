package coremodule

import "placeholder_misp/datamodels"

var (
	source datamodels.SourceMessageTheHive = datamodels.SourceMessageTheHive{}

	event                    datamodels.EventMessageTheHive       = datamodels.EventMessageTheHive{}
	eventObject              datamodels.EventObject               = datamodels.EventObject{}
	eventDetails             datamodels.EventDetails              = datamodels.EventDetails{}
	eventObjectCustomFields  map[string]datamodels.CustomerFields = make(map[string]datamodels.CustomerFields)
	eventDetailsCustomFields map[string]datamodels.CustomerFields = make(map[string]datamodels.CustomerFields)

	observables datamodels.ObservablesMessageTheHive = datamodels.ObservablesMessageTheHive{}
)

// ------ EVENT ------
var listHandlerEvent map[string][]func(interface{}) = map[string][]func(interface{}){
	"event.rootId":         {event.SetAnyRootId},
	"event.objectId":       {event.SetAnyObjectId},
	"event.objectType":     {event.SetAnyObjectType},
	"event.base":           {event.SetAnyBase},
	"event.startDate":      {event.SetAnyStartDate},
	"event.requestId":      {event.SetAnyRequestId},
	"event.organisation":   {event.SetAnyOperation},
	"event.organisationId": {event.SetAnyOrganisationId},
	"event.operation":      {event.SetAnyOperation},
}

// ------ EVENT DETAILS ------
var listHandlerEventDetails map[string][]func(interface{}) = map[string][]func(interface{}){
	"event.details.endDate":          {eventDetails.SetAnyEndDate},
	"event.details.resolutionStatus": {eventDetails.SetAnyResolutionStatus},
	"event.details.summary":          {eventDetails.SetAnySummary},
	"event.details.status":           {eventDetails.SetAnyStatus},
	"event.details.impactStatus":     {eventDetails.SetAnyImpactStatus},
}

// ------ EVENT DETAILS CUSTOMFIELDS ------
var listHandlerEventDetailsCustomFields map[string][]func(interface{}) = map[string][]func(interface{}){
	//--- attack-type ---
	"event.details.customFields.attack-type.order": {func(i interface{}) {
		//создаем элемент "attack-type" если его нет
		newCustomFieldsElement("attack-type", "string", &eventDetailsCustomFields)
		_, _, _, str := eventDetailsCustomFields["attack-type"].Get()
		eventDetailsCustomFields["attack-type"].Set(i, str)
	}},
	"event.details.customFields.attack-type.string": {func(i interface{}) {
		newCustomFieldsElement("attack-type", "string", &eventDetailsCustomFields)
		_, order, _, _ := eventDetailsCustomFields["attack-type"].Get()
		eventDetailsCustomFields["attack-type"].Set(order, i)
	}},
	//--- class-attack ---
	"event.details.customFields.class-attack.order": {func(i interface{}) {
		newCustomFieldsElement("class-attack", "string", &eventDetailsCustomFields)
		_, _, _, str := eventDetailsCustomFields["class-attack"].Get()
		eventDetailsCustomFields["class-attack"].Set(i, str)
	}},
	"event.details.customFields.class-attack.string": {func(i interface{}) {
		newCustomFieldsElement("class-attack", "string", &eventDetailsCustomFields)
		_, order, _, _ := eventDetailsCustomFields["class-attack"].Get()
		eventDetailsCustomFields["class-attack"].Set(order, i)
	}},
	//--- event-source ---
	"event.details.customFields.event-source.order": {func(i interface{}) {
		newCustomFieldsElement("event-source", "string", &eventDetailsCustomFields)
		_, _, _, str := eventDetailsCustomFields["event-source"].Get()
		eventDetailsCustomFields["event-source"].Set(i, str)
	}},
	"event.details.customFields.event-source.string": {func(i interface{}) {
		newCustomFieldsElement("event-source", "string", &eventDetailsCustomFields)
		_, order, _, _ := eventDetailsCustomFields["event-source"].Get()
		eventDetailsCustomFields["event-source"].Set(order, i)
	}},
	//--- ncircc-class-attack ---
	"event.details.customFields.ncircc-class-attack.order": {func(i interface{}) {
		newCustomFieldsElement("ncircc-class-attack", "string", &eventDetailsCustomFields)
		_, _, _, str := eventDetailsCustomFields["ncircc-class-attack"].Get()
		eventDetailsCustomFields["ncircc-class-attack"].Set(i, str)
	}},
	"event.details.customFields.ncircc-class-attack.string": {func(i interface{}) {
		newCustomFieldsElement("ncircc-class-attack", "string", &eventDetailsCustomFields)
		_, order, _, _ := eventDetailsCustomFields["ncircc-class-attack"].Get()
		eventDetailsCustomFields["ncircc-class-attack"].Set(order, i)
	}},
	//--- ncircc-bulletin-id ---
	"event.details.customFields.ncircc-bulletin-id.order": {func(i interface{}) {
		newCustomFieldsElement("ncircc-bulletin-id", "string", &eventDetailsCustomFields)
		_, _, _, str := eventDetailsCustomFields["ncircc-bulletin-id"].Get()
		eventDetailsCustomFields["ncircc-bulletin-id"].Set(i, str)
	}},
	"event.details.customFields.ncircc-bulletin-id.string": {func(i interface{}) {
		newCustomFieldsElement("ncircc-bulletin-id", "string", &eventDetailsCustomFields)
		_, order, _, _ := eventDetailsCustomFields["ncircc-bulletin-id"].Get()
		eventDetailsCustomFields["ncircc-bulletin-id"].Set(order, i)
	}},
	//--- ir-name ---
	"event.details.customFields.ir-name.order": {func(i interface{}) {
		newCustomFieldsElement("ir-name", "string", &eventDetailsCustomFields)
		_, _, _, str := eventDetailsCustomFields[" "].Get()
		eventDetailsCustomFields[" "].Set(i, str)
	}},
	"event.details.customFields.ir-name.string": {func(i interface{}) {
		newCustomFieldsElement("ir-name", "string", &eventDetailsCustomFields)
		_, order, _, _ := eventDetailsCustomFields["ir-name"].Get()
		eventDetailsCustomFields["ir-name"].Set(order, i)
	}},
	//--- sphere ---
	"event.details.customFields.sphere.order": {func(i interface{}) {
		newCustomFieldsElement("sphere", "string", &eventDetailsCustomFields)
		_, _, _, str := eventDetailsCustomFields["sphere"].Get()
		eventDetailsCustomFields["sphere"].Set(i, str)
	}},
	"event.details.customFields.sphere.string": {func(i interface{}) {
		newCustomFieldsElement("sphere", "string", &eventDetailsCustomFields)
		_, order, _, _ := eventDetailsCustomFields["sphere"].Get()
		eventDetailsCustomFields["sphere"].Set(order, i)
	}},
	//--- id-soa ---
	"event.details.customFields.id-soa.order": {func(i interface{}) {
		newCustomFieldsElement("id-soa", "string", &eventDetailsCustomFields)
		_, _, _, str := eventDetailsCustomFields["id-soa"].Get()
		eventDetailsCustomFields["id-soa"].Set(i, str)
	}},
	"event.details.customFields.id-soa.string": {func(i interface{}) {
		newCustomFieldsElement("id-soa", "string", &eventDetailsCustomFields)
		_, order, _, _ := eventDetailsCustomFields["id-soa"].Get()
		eventDetailsCustomFields["id-soa"].Set(order, i)
	}},
	//--- state ---
	"event.details.customFields.state.order": {func(i interface{}) {
		newCustomFieldsElement("state", "string", &eventDetailsCustomFields)
		_, _, _, str := eventDetailsCustomFields["state"].Get()
		eventDetailsCustomFields["state"].Set(i, str)
	}},
	"event.details.customFields.state.string": {func(i interface{}) {
		newCustomFieldsElement("state", "string", &eventDetailsCustomFields)
		_, order, _, _ := eventDetailsCustomFields["state"].Get()
		eventDetailsCustomFields["state"].Set(order, i)
	}},
	//--- external-letter ---
	"event.details.customFields.external-letter.order": {func(i interface{}) {
		newCustomFieldsElement("external-letter", "string", &eventDetailsCustomFields)
		_, _, _, str := eventDetailsCustomFields["external-letter"].Get()
		eventDetailsCustomFields["external-letter"].Set(i, str)
	}},
	//--- inbox1 ---
	"event.details.customFields.inbox1.order": {func(i interface{}) {
		newCustomFieldsElement("inbox1", "string", &eventDetailsCustomFields)
		_, _, _, str := eventDetailsCustomFields["inbox1"].Get()
		eventDetailsCustomFields["inbox1"].Set(i, str)
	}},
	//--- inner-letter ---
	"event.details.customFields.inner-letter.order": {func(i interface{}) {
		newCustomFieldsElement("inner-letter", "string", &eventDetailsCustomFields)
		_, _, _, str := eventDetailsCustomFields["inner-letter"].Get()
		eventDetailsCustomFields["inner-letter"].Set(i, str)
	}},
	//--- notification ---
	"event.details.customFields.notification.order": {func(i interface{}) {
		newCustomFieldsElement("notification", "string", &eventDetailsCustomFields)
		_, _, _, str := eventDetailsCustomFields["notification"].Get()
		eventDetailsCustomFields["notification"].Set(i, str)
	}},
	//--- report ---
	"event.details.customFields.report.order": {func(i interface{}) {
		newCustomFieldsElement("report", "string", &eventDetailsCustomFields)
		_, _, _, str := eventDetailsCustomFields["report"].Get()
		eventDetailsCustomFields["report"].Set(i, str)
	}},
	//--- first-time ---
	"event.details.customFields.first-time.order": {func(i interface{}) {
		newCustomFieldsElement("first-time", "string", &eventDetailsCustomFields)
		_, _, _, str := eventDetailsCustomFields["first-time"].Get()
		eventDetailsCustomFields["first-time"].Set(i, str)
	}},
	"event.details.customFields.first-time.date": {func(i interface{}) {
		newCustomFieldsElement("first-time", "date", &eventDetailsCustomFields)
		_, order, _, _ := eventDetailsCustomFields["first-time"].Get()
		eventDetailsCustomFields["first-time"].Set(order, i)
	}},
	//--- last-time ---
	"event.details.customFields.last-time.order": {func(i interface{}) {
		newCustomFieldsElement("last-time", "string", &eventDetailsCustomFields)
		_, _, _, str := eventDetailsCustomFields["last-time"].Get()
		eventDetailsCustomFields["last-time"].Set(i, str)
	}},
	"event.details.customFields.last-time.date": {func(i interface{}) {
		newCustomFieldsElement("last-time", "date", &eventDetailsCustomFields)
		_, order, _, _ := eventDetailsCustomFields["last-time"].Get()
		eventDetailsCustomFields["last-time"].Set(order, i)
	}},
	//--- b2mid ---
	"event.details.customFields.b2mid.order": {func(i interface{}) {
		newCustomFieldsElement("b2mid", "integer", &eventDetailsCustomFields)
		_, _, _, str := eventDetailsCustomFields["b2mid"].Get()
		eventDetailsCustomFields["b2mid"].Set(i, str)
	}},
	"event.details.customFields.b2mid.integer": {func(i interface{}) {
		newCustomFieldsElement("b2mid", "integer", &eventDetailsCustomFields)
		_, order, _, _ := eventDetailsCustomFields["b2mid"].Get()
		eventDetailsCustomFields["b2mid"].Set(order, i)
	}},
}

// ------ EVENT OBJECT ------
var listHandlerEventObject map[string][]func(interface{}) = map[string][]func(interface{}){
	"event.object.flag":             {eventObject.SetAnyFlag},
	"event.object.caseId":           {eventObject.SetAnyCaseId},
	"event.object.severity":         {eventObject.SetAnySeverity},
	"event.object.tlp":              {eventObject.SetAnyTlp},
	"event.object.pap":              {eventObject.SetAnyPap},
	"event.object.startDate":        {eventObject.SetAnyStartDate},
	"event.object.endDate":          {eventObject.SetAnyEndDate},
	"event.object.createdAt":        {eventObject.SetAnyCreatedAt},
	"event.object.updatedAt":        {eventObject.SetAnyUpdatedAt},
	"event.object._id":              {eventObject.SetAnyUnderliningId},
	"event.object.id":               {eventObject.SetAnyId},
	"event.object.createdBy":        {eventObject.SetAnyCreatedBy},
	"event.object.updatedBy":        {eventObject.SetAnyUpdatedBy},
	"event.object._type":            {eventObject.SetAnyUnderliningType},
	"event.object.title":            {eventObject.SetAnyTitle},
	"event.object.description":      {eventObject.SetAnyDescription},
	"event.object.impactStatus":     {eventObject.SetAnyImpactStatus},
	"event.object.resolutionStatus": {eventObject.SetAnyResolutionStatus},
	"event.object.status":           {eventObject.SetAnyStatus},
	"event.object.summary":          {eventObject.SetAnySummary},
	"event.object.owner":            {eventObject.SetAnyOwner},
	"event.object.tags":             {eventObject.SetAnyTags},

	//ниже следующие поля редко используются, думаю пока они не требуют реализации
	//"event.object.stats.impactStatus":    {},
	//"event.object.permissions.id":        {},
	//"event.object.permissions.createdAt": {},
	//"event.object.permissions.pap":       {},
}

// ------ EVENT OBJECT CUSTOMFIELDS ------
var listHandlerEventObjectCustomFields map[string][]func(interface{}) = map[string][]func(interface{}){
	//--- attack-type ---
	"event.object.customFields.attack-type.order": {func(i interface{}) {
		//создаем элемент "attack-type" если его нет
		newCustomFieldsElement("attack-type", "string", &eventObjectCustomFields)
		_, _, _, str := eventObjectCustomFields["attack-type"].Get()
		eventObjectCustomFields["attack-type"].Set(i, str)
	}},
	"event.object.customFields.attack-type.string": {func(i interface{}) {
		newCustomFieldsElement("attack-type", "string", &eventObjectCustomFields)
		_, order, _, _ := eventObjectCustomFields["attack-type"].Get()
		eventObjectCustomFields["attack-type"].Set(order, i)
	}},
	//--- class-attack ---
	"event.object.customFields.class-attack.order": {func(i interface{}) {
		newCustomFieldsElement("class-attack", "string", &eventObjectCustomFields)
		_, _, _, str := eventObjectCustomFields["class-attack"].Get()
		eventObjectCustomFields["class-attack"].Set(i, str)
	}},
	"event.object.customFields.class-attack.string": {func(i interface{}) {
		newCustomFieldsElement("class-attack", "string", &eventObjectCustomFields)
		_, order, _, _ := eventObjectCustomFields["class-attack"].Get()
		eventObjectCustomFields["class-attack"].Set(order, i)
	}},
	//--- ncircc-class-attack ---
	"event.object.customFields.ncircc-class-attack.order": {func(i interface{}) {
		newCustomFieldsElement("ncircc-class-attack", "string", &eventObjectCustomFields)
		_, _, _, str := eventObjectCustomFields["ncircc-class-attack"].Get()
		eventObjectCustomFields["ncircc-class-attack"].Set(i, str)
	}},
	"event.object.customFields.ncircc-class-attack.string": {func(i interface{}) {
		newCustomFieldsElement("ncircc-class-attack", "string", &eventObjectCustomFields)
		_, order, _, _ := eventObjectCustomFields["ncircc-class-attack"].Get()
		eventObjectCustomFields["ncircc-class-attack"].Set(order, i)
	}},
	//--- inbox1 ---
	"event.object.customFields.inbox1.order": {func(i interface{}) {
		newCustomFieldsElement("inbox1", "string", &eventObjectCustomFields)
		_, _, _, str := eventObjectCustomFields["inbox1"].Get()
		eventObjectCustomFields["inbox1"].Set(i, str)
	}},
	//--- inner-letter ---
	"event.object.customFields.inner-letter.order": {func(i interface{}) {
		newCustomFieldsElement("inner-letter", "string", &eventObjectCustomFields)
		_, _, _, str := eventObjectCustomFields["inner-letter"].Get()
		eventObjectCustomFields["inner-letter"].Set(i, str)
	}},
	//--- notification ---
	"event.object.customFields.notification.order": {func(i interface{}) {
		newCustomFieldsElement("notification", "string", &eventObjectCustomFields)
		_, _, _, str := eventObjectCustomFields["notification"].Get()
		eventObjectCustomFields["notification"].Set(i, str)
	}},
	//--- report ---
	"event.object.customFields.report.order": {func(i interface{}) {
		newCustomFieldsElement("report", "string", &eventObjectCustomFields)
		_, _, _, str := eventObjectCustomFields["report"].Get()
		eventObjectCustomFields["report"].Set(i, str)
	}},
}

// ------ OBSERVABLES ------
var listHandlerObservables map[string][]func(interface{}) = map[string][]func(interface{}){

	/*
		Тут надо подумать как обрабатывать пути начинающиеся на observables.* так как
		datamodels.ObservablesMessageTheHive является списком
	*/

	"observables.rootId": {observables.UnderliningId},
}

func newCustomFieldsElement(elem, objType string, customFields *map[string]datamodels.CustomerFields) {
	if _, ok := (*customFields)[elem]; !ok {
		switch objType {
		case "string":
			(*customFields)[elem] = &datamodels.CustomFieldStringType{}
		case "date":
			(*customFields)[elem] = &datamodels.CustomFieldDateType{}
		case "integer":
			(*customFields)[elem] = &datamodels.CustomFieldIntegerType{}
		}
	}
}

func joinListHandler(listSeveral []map[string][]func(interface{})) map[string][]func(interface{}) {
	result := make(map[string][]func(interface{}))

	for _, list := range listSeveral {
		for k, v := range list {
			result[k] = v
		}
	}

	return result
}
