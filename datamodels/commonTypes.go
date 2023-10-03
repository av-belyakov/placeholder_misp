package datamodels

// MessageLoging содержит информацию используемую при логировании
// MsgData - сообщение
// MsgType - тип сообщения
type MessageLoging struct {
	MsgData, MsgType string
}

// DataCounterSettings содержит информацию для подсчета
type DataCounterSettings struct {
	DataType string
	Count    int
}
