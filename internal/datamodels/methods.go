package datamodels

func NewMessageLogging() *MessageLogging {
	return &MessageLogging{}
}

// GetMessage возвращает сообщение
func (ml *MessageLogging) GetMessage() string {
	return ml.Message
}

// SetMessage устанавливает сообщение
func (ml *MessageLogging) SetMessage(v string) {
	ml.Message = v
}

// GetType возвращает тип сообщения
func (ml *MessageLogging) GetType() string {
	return ml.Type
}

// SetType устанавливает тип сообщения
func (ml *MessageLogging) SetType(v string) {
	ml.Type = v
}
