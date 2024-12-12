package datamodels

// ListFormatsMISP содержит описание типов добавляемых в MISP и их порядок добавления.
// По результатам добавления Event, MISP возвращает id котрый необходим как для добавления
// следующих типов объектов MISP, так и для добавления этого значения в поле
// 'customFields' TheHive. Не все из этих объектов могут сразу добавлятся в MISP 'как есть',
// некоторые из них подлежат дополнительной обработке, см. обработчик для каждого из объектов.
// После добавления всех объектов, событие MISP необходимо опобликовать, как это сделать
// см. обработчик публикации.
type ListFormatsMISP struct {
	ID         string
	Event      *EventsMispFormat
	Reports    *EventReports
	Attributes []*AttributesMispFormat
	Objects    map[int]*ObjectsMispFormat
	ObjectTags *ListEventObjectTags
}
