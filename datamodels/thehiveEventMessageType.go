package datamodels

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
	Base           bool         `json:"base"`
	StartDate      uint64       `json:"startDate"`
	RootId         string       `json:"rootId"`
	Organisation   string       `json:"organisation"`
	OrganisationId string       `json:"organisationId"`
	ObjectId       string       `json:"objectId"`
	ObjectType     string       `json:"objectType"`
	Operation      string       `json:"operation"`
	RequestId      string       `json:"requestId"`
	Details        EventDetails `json:"details"`
	Object         EventObject  `json:"object"`
}

// EventDetails детальная информация о событии
// EndDate - конечное дата и время
// ResolutionStatus - статус постановления
// Summary - резюме
// Status - статус
// ImpactStatus - краткое описание воздействия
// CustomFields - настраиваемые поля
type EventDetails struct {
	EndDate          uint64                    `json:"endDate"`
	ResolutionStatus string                    `json:"resolutionStatus"`
	Summary          string                    `json:"summary"`
	Status           string                    `json:"status"`
	ImpactStatus     string                    `json:"impactStatus"`
	CustomFields     map[string]CustomerFields `json:"customFields"`
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
	Flag             bool                      `json:"flag"`
	CaseId           int                       `json:"caseId"`
	Severity         int                       `json:"severity"`
	Tlp              int                       `json:"tlp"`
	Pap              int                       `json:"pap"`
	StartDate        uint64                    `json:"startDate"`
	EndDate          uint64                    `json:"endDate"`
	CreatedAt        uint64                    `json:"createdAt"`
	UpdatedAt        uint64                    `json:"updatedAt"`
	UnderliningId    string                    `json:"_id"`
	Id               string                    `json:"id"`
	CreatedBy        string                    `json:"createdBy"`
	UpdatedBy        string                    `json:"updatedBy"`
	UnderliningType  string                    `json:"_type"`
	Title            string                    `json:"title"`
	Description      string                    `json:"description"`
	ImpactStatus     string                    `json:"impactStatus"`
	ResolutionStatus string                    `json:"resolutionStatus"`
	Status           string                    `json:"status"`
	Summary          string                    `json:"summary"`
	Owner            string                    `json:"owner"`
	Tags             []string                  `json:"tags"`
	CustomFields     map[string]CustomerFields `json:"customFields"`
	//данное поле редко используемое, думаю пока оно не требует реализации
	//Stats            map[string]interface{} `json:"stats"`
	//данное поле редко используемое, думаю пока оно не требует реализации
	//Permissions  []string              `json:"permissions"`
}
