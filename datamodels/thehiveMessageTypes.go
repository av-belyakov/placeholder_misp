package datamodels

// ResponseMessageFromMispToTheHave содержит ответ для TheHive получаемый от MISP
type ResponseMessageFromMispToTheHave struct {
	Success  bool                        `json:"success"`
	Service  string                      `json:"string"`
	Commands []ResponseCommandForTheHive `json:"commands"`
}

// ResponseCommandForTheHive ответы с командами для TheHive
type ResponseCommandForTheHive struct {
	Command string `json:"command"`
	Name    string `json:"name"`
	String  string `json:"string"`
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

// EventMessageTheHive сообщение с информацией о событии
// Operation - операция
// ObjectId - уникальный идентификатор объекта
// ObjectType - тип объекта
// Base - основа
// StartDate - начальная дата
// RootId - главный уникальный идентификатор
// RequestId - уникальный идентификатор запроса
// EventDetails - детальная информация о событии
// Object - объект события
// OrganisationId - уникальный идентификатор организации
// Organisation - наименование организации
type EventMessageTheHive struct {
	Operation      string       `json:"operation"`
	ObjectId       string       `json:"objectId"`
	ObjectType     string       `json:"objectType"`
	Base           bool         `json:"base"`
	StartDate      uint64       `json:"startDate"`
	RootId         string       `json:"rootId"`
	RequestId      string       `json:"requestId"`
	Details        EventDetails `json:"details"`
	Object         EventObject  `json:"object"`
	OrganisationId string       `json:"organisationId"`
	Organisation   string       `json:"organisation"`
}

// EventDetails детальная информация о событии
// EndDate - конечное дата и время
// ResolutionStatus - статус постановления
// Summary - резюме
// Status - статус
// ImpactStatus - краткое описание воздействия
// CustomFields - настраиваемые поля
type EventDetails struct {
	EndDate          uint64 `json:"endDate"`
	ResolutionStatus string `json:"resolutionStatus"`
	Summary          string `json:"summary"`
	Status           string `json:"status"`
	ImpactStatus     string `json:"impactStatus"`
	CustomFields     `json:"customFields"`
}

// EventObject объект события
// UnderliningId - уникальный идентификатор
// Id - уникальный идентификатор
// CreatedBy - кем создан
// UpdatedBy - кем обновлен
// CreatedAt - дата создания
// UpdatedAt - дата обновления
// UnderliningType - тип
// CaseId - уникальный идентификатор дела
// Title - заголовок
// Description - описание
// Severity - строгость
// StartDate - начальная дата
// EndDate - конечная дата
// ImpactStatus - краткое описание воздействия
// ResolutionStatus - статус разрешения
// Tags - список тегов
// Flag - флаг
// Tlp - tlp
// Pap - pap
// Status - статус
// Summary - резюме
// Owner - владелец
// CustomFields - настраиваемые поля
// Stats - статистика
// Permissions - разрешения
type EventObject struct {
	UnderliningId    string   `json:"_id"`
	Id               string   `json:"id"`
	CreatedBy        string   `json:"createdBy"`
	UpdatedBy        string   `json:"updatedBy"`
	CreatedAt        uint64   `json:"createdAt"`
	UpdatedAt        uint64   `json:"updatedAt"`
	UnderliningType  string   `json:"_type"`
	CaseId           int      `json:"caseId"`
	Title            string   `json:"title"`
	Description      string   `json:"description"`
	Severity         int      `json:"severity"`
	StartDate        uint64   `json:"startDate"`
	EndDate          uint64   `json:"endDate"`
	ImpactStatus     string   `json:"impactStatus"`
	ResolutionStatus string   `json:"resolutionStatus"`
	Tags             []string `json:"tags"`
	Flag             bool     `json:"flag"`
	Tlp              int      `json:"tlp"`
	Pap              int      `json:"pap"`
	Status           string   `json:"status"`
	Summary          string   `json:"summary"`
	Owner            string   `json:"owner"`
	CustomFields     `json:"customFields"`
	Stats            map[string]interface{} `json:"stats"`
	Permissions      []string               `json:"permissions"`
}

// CustomFields настраиваемые поля
type CustomFields map[string]map[string]interface{}

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
	CreatedAt        uint64                                         `json:"_createdAt"`
	CreatedBy        string                                         `json:"_createdBy"`
	UnderliningId    string                                         `json:"_id"`
	UnderliningType  string                                         `json:"_type"`
	UpdatedAt        uint64                                         `json:"_updatedAt"`
	UpdatedBy        string                                         `json:"_updatedBy"`
	Data             string                                         `json:"data"`
	DataType         string                                         `json:"dataType"`
	IgnoreSimilarity bool                                           `json:"ignoreSimilarity"`
	ExtraData        map[string]interface{}                         `json:"extraData"`
	Ioc              bool                                           `json:"ioc"`
	Message          string                                         `json:"message"`
	Sighted          bool                                           `json:"sighted"`
	StartDate        uint64                                         `json:"startDate"`
	Tags             []string                                       `json:"tags"`
	Tlp              int                                            `json:"tlp"`
	Reports          map[string]map[string][]map[string]interface{} `json:"reports"`
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
	CreatedAt     uint64              `json:"_createdAt"`
	CreatedBy     string              `json:"_createdBy"`
	UnderliningId string              `json:"_id"`
	ExtraData     ExtraDataTtpMessage `json:"extraData"`
	OccurDate     uint64              `json:"occurDate"`
	PatternId     string              `json:"patternId"`
	Tactic        string              `json:"tactic"`
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
	CreatedAt           uint64                 `json:"_createdAt"`
	CreatedBy           string                 `json:"_createdBy"`
	UnderliningId       string                 `json:"_id"`
	UnderliningType     string                 `json:"_type"`
	DataSources         []string               `json:"dataSources"`
	DefenseBypassed     []string               `json:"defenseBypassed"` //надо проверить тип
	Description         string                 `json:"description"`
	ExtraData           map[string]interface{} `json:"extraData"`
	Name                string                 `json:"name"`
	PatternId           string                 `json:"patternId"`
	PatternType         string                 `json:"patternType"`
	PermissionsRequired []string               `json:"permissionsRequired"`
	Platforms           []string               `json:"platforms"`
	RemoteSupport       bool                   `json:"remoteSupport"`
	Revoked             bool                   `json:"revoked"`
	SystemRequirements  []string               `json:"systemRequirements"` //надо проверить тип
	Tactics             []string               `json:"tactics"`
	URL                 string                 `json:"url"`
	Version             string                 `json:"version"`
}
