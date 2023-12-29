package datamodels

type CustomerFields interface {
	Set(interface{}, interface{})
	Get() (string, int, string, string)
}

// ResponseMessageFromMispToTheHave содержит ответ для TheHive получаемый от MISP
type ResponseMessageFromMispToTheHave struct {
	Success  bool                        `json:"success"`
	Service  string                      `json:"service"`
	Error    error                       `json:"error"`
	Commands []ResponseCommandForTheHive `json:"commands"`
}

// ResponseCommandForTheHive ответы с командами для TheHive
type ResponseCommandForTheHive struct {
	Command string `json:"command"`
	String  string `json:"string"`
	Name    string `json:"name"`
}

// MainMessageTheHive основное сообщение получаемое через NATS
type MainMessageTheHive struct {
	SourceMessageTheHive
	EventMessageTheHive
	ObservablesMessageTheHive
	TtpsMessageTheHive
}

// SourceMessageTheHive сообщение с информацией об источнике
// Source - источник
type SourceMessageTheHive struct {
	Source string `json:"source"`
}

// TtpsMessageTheHive список TTP сообщений
type TtpsMessageTheHive struct {
	Ttp []TtpMessage `json:"ttp"`
}

// TtpMessage TTP сообщения
// CreatedAt - время создания
// CreatedBy - кем создан
// UnderliningId - уникальный идентификатор
// ExtraData - дополнительные данные
// OccurDate - дата возникновения
// PatternId - уникальный идентификатор шаблона
// Tactic - тактика
type TtpMessage struct {
	OccurDate     uint64              `json:"occurDate"`
	CreatedAt     uint64              `json:"_createdAt"`
	UnderliningId string              `json:"_id"`
	CreatedBy     string              `json:"_createdBy"`
	PatternId     string              `json:"patternId"`
	Tactic        string              `json:"tactic"`
	ExtraData     ExtraDataTtpMessage `json:"extraData"`
}

// ExtraDataTtpMessage дополнительные данные TTP сообщения
// Pattern - шаблон
// PatternParent - родительский шаблон
type ExtraDataTtpMessage struct {
	Pattern       PatternExtraData `json:"pattern"`
	PatternParent PatternExtraData `json:"patternParent"`
}

// PatternExtraData шаблон дополнительных данных
// CreatedAt - время создания
// CreatedBy - кем создан
// UnderliningId - уникальный идентификатор
// UnderliningType - тип
// DataSources - источники данных
// DefenseBypassed - чем выполнен обход защиты
// Description - описание
// ExtraData - дополнительные данные
// Name - наименование
// PatternId - уникальный идентификатор шаблона
// PatternType - тип шаблона
// PermissionsRequired - требуемые разрешения
// Platforms - список платформ
// RemoteSupport - удаленная поддержка
// Revoked - аннулированный
// SystemRequirements - системные требования
// Tactics - список тактик
// URL - URL
// Version - версия
type PatternExtraData struct {
	RemoteSupport       bool     `json:"remoteSupport"`
	Revoked             bool     `json:"revoked"`
	CreatedAt           uint64   `json:"_createdAt"`
	CreatedBy           string   `json:"_createdBy"`
	UnderliningId       string   `json:"_id"`
	UnderliningType     string   `json:"_type"`
	Description         string   `json:"description"`
	Name                string   `json:"name"`
	PatternId           string   `json:"patternId"`
	PatternType         string   `json:"patternType"`
	URL                 string   `json:"url"`
	Version             string   `json:"version"`
	Platforms           []string `json:"platforms"`
	PermissionsRequired []string `json:"permissionsRequired"`
	DataSources         []string `json:"dataSources"`
	Tactics             []string `json:"tactics"`
	//данное поле редко используемое, думаю пока оно не требует реализации
	//DefenseBypassed     []string               `json:"defenseBypassed"` //надо проверить тип
	//данное поле редко используемое, думаю пока оно не требует реализации
	//SystemRequirements  []string               `json:"systemRequirements"` //надо проверить тип
	//данное поле редко используемое, думаю пока оно не требует реализации
	//ExtraData           map[string]interface{} `json:"extraData"`
}

type CustomFieldStringType struct {
	Order  int    `json:"order"`
	String string `json:"string"`
}

type CustomFieldDateType struct {
	Order int    `json:"order"`
	Date  uint64 `json:"date"`
}

type CustomFieldIntegerType struct {
	Order   int `json:"order"`
	Integer int `json:"integer"`
}
