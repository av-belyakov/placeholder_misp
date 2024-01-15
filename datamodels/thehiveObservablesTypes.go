package datamodels

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
	Ioc              bool                        `json:"ioc"`
	Sighted          bool                        `json:"sighted"`
	IgnoreSimilarity bool                        `json:"ignoreSimilarity"`
	Tlp              int                         `json:"tlp"`
	CreatedAt        uint64                      `json:"_createdAt"`
	UpdatedAt        uint64                      `json:"_updatedAt"`
	StartDate        uint64                      `json:"startDate"`
	CreatedBy        string                      `json:"_createdBy"`
	UpdatedBy        string                      `json:"_updatedBy"`
	UnderliningId    string                      `json:"_id"`
	UnderliningType  string                      `json:"_type"`
	Data             string                      `json:"data"`
	DataType         string                      `json:"dataType"`
	Message          string                      `json:"message"`
	Tags             []string                    `json:"tags"`
	Attachment       AttachmentData              `json:"attachment"`
	Reports          map[string]ReportTaxonomies `json:"reports"`
	//данное поле редко используемое, думаю пока оно не требует реализации
	//ExtraData        map[string]interface{}                         `json:"extraData"`
}

// AttachmentData прикрепленные данные
type AttachmentData struct {
	Size        int      `json:"size"`
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	ContentType string   `json:"contentType"`
	Hashes      []string `json:"hashes"`
}

// ReportTaxonomies
type ReportTaxonomies struct {
	Taxonomies []Taxonomy `json:"taxonomies"`
}

// Taxonomy
type Taxonomy struct {
	Level     string `json:"level"`
	Namespace string `json:"namespace"`
	Predicate string `json:"predicate"`
	Value     string `json:"value"`
}
