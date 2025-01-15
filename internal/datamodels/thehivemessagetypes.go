package datamodels

// ResponseMessageFromMispToTheHave содержит ответ для TheHive получаемый от MISP
type ResponseMessageFromMispToTheHave struct {
	Commands []ResponseCommandForTheHive `json:"commands"`
	Error    error                       `json:"error"`
	Service  string                      `json:"service"`
	Success  bool                        `json:"success"`
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
	RootId         string       `json:"rootId"`
	RequestId      string       `json:"requestId"`
	OrganisationId string       `json:"organisationId"`
	Organisation   string       `json:"organisation"`
	Details        EventDetails `json:"details"`
	Object         EventObject  `json:"object"`
	StartDate      uint64       `json:"startDate"`
	Base           bool         `json:"base"`
}

// EventDetails детальная информация о событии
// EndDate - конечное дата и время
// ResolutionStatus - статус постановления
// Summary - резюме
// Status - статус
// ImpactStatus - краткое описание воздействия
// CustomFields - настраиваемые поля
type EventDetails struct {
	CustomFields     `json:"customFields"`
	ResolutionStatus string `json:"resolutionStatus"`
	Summary          string `json:"summary"`
	Status           string `json:"status"`
	ImpactStatus     string `json:"impactStatus"`
	EndDate          uint64 `json:"endDate"`
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
	CustomFields     `json:"customFields"`
	Stats            map[string]interface{} `json:"stats"`
	Tags             []string               `json:"tags"`
	Permissions      []string               `json:"permissions"`
	UnderliningId    string                 `json:"_id"`
	Id               string                 `json:"id"`
	CreatedBy        string                 `json:"createdBy"`
	UpdatedBy        string                 `json:"updatedBy"`
	UnderliningType  string                 `json:"_type"`
	Title            string                 `json:"title"`
	Description      string                 `json:"description"`
	ImpactStatus     string                 `json:"impactStatus"`
	ResolutionStatus string                 `json:"resolutionStatus"`
	Status           string                 `json:"status"`
	Summary          string                 `json:"summary"`
	Owner            string                 `json:"owner"`
	CreatedAt        uint64                 `json:"createdAt"`
	UpdatedAt        uint64                 `json:"updatedAt"`
	StartDate        uint64                 `json:"startDate"`
	EndDate          uint64                 `json:"endDate"`
	CaseId           int                    `json:"caseId"`
	Severity         int                    `json:"severity"`
	Tlp              int                    `json:"tlp"`
	Pap              int                    `json:"pap"`
	Flag             bool                   `json:"flag"`
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
	Reports          map[string]map[string][]map[string]interface{} `json:"reports"`
	ExtraData        map[string]interface{}                         `json:"extraData"`
	Tags             []string                                       `json:"tags"`
	CreatedBy        string                                         `json:"_createdBy"`
	UnderliningId    string                                         `json:"_id"`
	UnderliningType  string                                         `json:"_type"`
	UpdatedBy        string                                         `json:"_updatedBy"`
	Data             string                                         `json:"data"`
	DataType         string                                         `json:"dataType"`
	Message          string                                         `json:"message"`
	CreatedAt        uint64                                         `json:"_createdAt"`
	UpdatedAt        uint64                                         `json:"_updatedAt"`
	StartDate        uint64                                         `json:"startDate"`
	Tlp              int                                            `json:"tlp"`
	IgnoreSimilarity bool                                           `json:"ignoreSimilarity"`
	Ioc              bool                                           `json:"ioc"`
	Sighted          bool                                           `json:"sighted"`
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
	ExtraData     ExtraDataTtpMessage `json:"extraData"`
	CreatedBy     string              `json:"_createdBy"`
	UnderliningId string              `json:"_id"`
	PatternId     string              `json:"patternId"`
	Tactic        string              `json:"tactic"`
	CreatedAt     uint64              `json:"_createdAt"`
	OccurDate     uint64              `json:"occurDate"`
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
	ExtraData           map[string]interface{} `json:"extraData"`
	DataSources         []string               `json:"dataSources"`
	DefenseBypassed     []string               `json:"defenseBypassed"` //надо проверить тип
	PermissionsRequired []string               `json:"permissionsRequired"`
	Platforms           []string               `json:"platforms"`
	SystemRequirements  []string               `json:"systemRequirements"` //надо проверить тип
	Tactics             []string               `json:"tactics"`
	URL                 string                 `json:"url"`
	Version             string                 `json:"version"`
	CreatedBy           string                 `json:"_createdBy"`
	UnderliningId       string                 `json:"_id"`
	UnderliningType     string                 `json:"_type"`
	Description         string                 `json:"description"`
	Name                string                 `json:"name"`
	PatternId           string                 `json:"patternId"`
	PatternType         string                 `json:"patternType"`
	CreatedAt           uint64                 `json:"_createdAt"`
	RemoteSupport       bool                   `json:"remoteSupport"`
	Revoked             bool                   `json:"revoked"`
}
