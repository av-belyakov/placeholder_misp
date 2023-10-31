package datamodels

// MessageLogging содержит информацию используемую при логировании
// MsgData - сообщение
// MsgType - тип сообщения
type MessageLogging struct {
	MsgData, MsgType string
}

// DataCounterSettings содержит информацию для подсчета
type DataCounterSettings struct {
	DataType string
	Count    int
}

// ListEventObjectTags временное хранилище для тегов полученных из event.object.tags
type ListEventObjectTags []string

func NewListEventObjectTags() *ListEventObjectTags {
	return &ListEventObjectTags{}
}

func (leot *ListEventObjectTags) SetTag(v string) {
	*leot = append(*leot, v)
}

func (leot *ListEventObjectTags) GetListTags() ListEventObjectTags {
	return *leot
}

func (leot *ListEventObjectTags) CleanListTags() {
	leot = &ListEventObjectTags{}
}
