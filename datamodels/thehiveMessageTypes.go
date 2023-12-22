package datamodels

type CustomerFields interface {
	Set(int, interface{})
	Get() (int, string)
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

// CustomFields настраиваемые поля
type CustomFields map[string]CustomerFields

type CustomFieldStringType struct {
	Order  int    `json:"order"`
	String string `json:"string"`
}

type CustomFieldDateType struct {
	Order int    `json:"order"`
	Date  uint64 `json:"date"`
}

// ObservablesMessageTheHive список наблюдаемых сообщений
// Observables - наблюдаемые сообщения
type ObservablesMessageTheHive struct {
	Observables []ObservableMessage `json:"observables"`
}

// ObservableMessage наблюдаемое сообщение
// CreatedAt - время создания
// CreatedBy - кем создан
// UnderliningId - уникальный идентификатор
// UnderliningType - тип
// UpdatedAt - время обновления
// UpdatedBy - кем обновлен
// Data - данные
// DataType - тип данных
// IgnoreSimilarity - игнорировать сходство
// ExtraData - дополнительные данные
// Ioc - индикатор компрометации
// Message - сообщение
// Sighted - видящий
// StartDate - дата начала
// Tags - список тегов
// Tlp - tlp
// Reports - список отчетов
type ObservableMessage struct {
	Ioc              bool                             `json:"ioc"`
	Sighted          bool                             `json:"sighted"`
	IgnoreSimilarity bool                             `json:"ignoreSimilarity"`
	Tlp              int                              `json:"tlp"`
	CreatedAt        uint64                           `json:"_createdAt"`
	UpdatedAt        uint64                           `json:"_updatedAt"`
	StartDate        uint64                           `json:"startDate"`
	CreatedBy        string                           `json:"_createdBy"`
	UpdatedBy        string                           `json:"_updatedBy"`
	UnderliningId    string                           `json:"_id"`
	UnderliningType  string                           `json:"_type"`
	Data             string                           `json:"data"`
	DataType         string                           `json:"dataType"`
	Message          string                           `json:"message"`
	Tags             []string                         `json:"tags"`
	Attachment       AttachmentData                   `json:"attachment"`
	Reports          map[string]map[string][]Taxonomy `json:"reports"`
	//данное поле редко используемое, думаю пока оно не требует реализации
	//ExtraData        map[string]interface{}                         `json:"extraData"`
}

// AttachmentData прикрепленные данные
type AttachmentData struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Size        int      `json:"size"`
	ContentType string   `json:"contentType"`
	Hashes      []string `json:"hashes"`
}

// Taxonomy
type Taxonomy struct {
	Level     string `json:"level"`
	Namespace string `json:"namespace"`
	Predicate string `json:"predicate"`
	Value     string `json:"value"`
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
