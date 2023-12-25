package datamodels

import (
	"fmt"
	"time"

	"placeholder_misp/supportingfunctions"
)

// NewResponseMessage формирует новый тип ResponseMessageFromMispToTheHave с предустановленными значениями
func NewResponseMessage() *ResponseMessageFromMispToTheHave {
	return &ResponseMessageFromMispToTheHave{
		Success: true,
		Service: "MISP",
		Commands: []ResponseCommandForTheHive{
			{
				Command: "addtag",
				String:  "Webhook: send=\"MISP\"",
				//String:  "Webhook: send=\"MISP-WORLD\"",
				//String: "Webhook: send=\"MISP-CENTER\"",
			},
		},
	}
}

func (rm *ResponseMessageFromMispToTheHave) ResponseMessageAddNewCommand(rcm ResponseCommandForTheHive) {
	rm.Commands = append(rm.Commands, rcm)
}

func (rm *ResponseMessageFromMispToTheHave) GetResponseMessageFromMispToTheHave() ResponseMessageFromMispToTheHave {
	return *rm
}

// Get возвращает MainMessageTheHive
func (mm *MainMessageTheHive) Get() *MainMessageTheHive {
	return mm
}

func (mm MainMessageTheHive) ToStringBeautiful(num int) string {
	var str string

	str += mm.SourceMessageTheHive.ToStringBeautiful(num + 1)
	str += fmt.Sprintln("event:")
	str += mm.EventMessageTheHive.ToStringBeautiful(num + 1)
	str += fmt.Sprintln("observables:")
	str += mm.ObservablesMessageTheHive.ToStringBeautiful(num + 1)
	str += fmt.Sprintln("ttps:")
	str += mm.TtpsMessageTheHive.ToStringBeautiful(num + 1)

	return str
}

// GetSource возвращает содержимое поля Source
func (s *SourceMessageTheHive) GetSource() string {
	return s.Source
}

// SetValueSource устанавливает СТРОКОВОЕ значение для поля Source
func (s *SourceMessageTheHive) SetValueSource(v string) {
	s.Source = v
}

// SetAnySource устанавливает ЛЮБОЕ значение для поля Source
func (s *SourceMessageTheHive) SetAnySource(i interface{}) {
	s.Source = fmt.Sprint(i)
}

func (s SourceMessageTheHive) ToStringBeautiful(num int) string {
	return fmt.Sprintf("source: '%s'\n", s.Source)
}

func (tm TtpsMessageTheHive) ToStringBeautiful(num int) string {
	return fmt.Sprintf("%sttp: \n%s", supportingfunctions.GetWhitespace(num), func(l []TtpMessage) string {
		var str string
		for k, v := range l {
			str += fmt.Sprintf("%s%d.\n", supportingfunctions.GetWhitespace(num+1), k+1)
			str += v.ToStringBeautiful(num + 2)
		}
		return str
	}(tm.Ttp))
}

func (tm TtpMessage) ToStringBeautiful(num int) string {
	var str string

	ws := supportingfunctions.GetWhitespace(num)

	str += fmt.Sprintf("%s_createdAt: '%d'\n", ws, tm.CreatedAt)
	str += fmt.Sprintf("%s_createdBy: '%s'\n", ws, tm.CreatedBy)
	str += fmt.Sprintf("%s_id: '%s'\n", ws, tm.UnderliningId)
	str += tm.ExtraData.ToStringBeautiful(num + 1)
	str += fmt.Sprintf("%soccurDate: '%d'\n", ws, tm.OccurDate)
	str += fmt.Sprintf("%spatternId: '%s'\n", ws, tm.PatternId)
	str += fmt.Sprintf("%stactic: '%s'\n", ws, tm.Tactic)

	return str
}

func (edtm ExtraDataTtpMessage) ToStringBeautiful(num int) string {
	var str string

	str += edtm.Pattern.ToStringBeautiful(num)
	str += edtm.PatternParent.ToStringBeautiful(num)

	return str
}

func (ped PatternExtraData) ToStringBeautiful(num int) string {
	var str string

	ws := supportingfunctions.GetWhitespace(num)

	str += fmt.Sprintf("%s_createdAt: '%d'\n", ws, ped.CreatedAt)
	str += fmt.Sprintf("%s_createdBy: '%s'\n", ws, ped.CreatedBy)
	str += fmt.Sprintf("%s_id: '%s'\n", ws, ped.UnderliningId)
	str += fmt.Sprintf("%s_type: '%s'\n", ws, ped.UnderliningType)
	str += fmt.Sprintf("%sdataSources: \n%v", ws, func(l []string) string {
		var str string
		for k, v := range l {
			str += fmt.Sprintf("%s%d. '%s'\n", supportingfunctions.GetWhitespace(num+1), k+1, v)
		}
		return str
	}(ped.DataSources))
	/*str += fmt.Sprintf("%sdefenseBypassed: \n%v", ws, func(l []string) string {
		var str string
		for k, v := range l {
			str += fmt.Sprintf("%s%d. '%s'\n", supportingfunctions.GetWhitespace(num+1), k+1, v)
		}
		return str
	}(ped.DefenseBypassed))*/
	str += fmt.Sprintf("%sdescription: '%s'\n", ws, ped.Description)
	/*str += fmt.Sprintf("%sextraData: \n%s", ws, func(l map[string]interface{}) string {
		var str string
		for k, v := range l {
			str += fmt.Sprintf("%s%s: '%v'\n", supportingfunctions.GetWhitespace(num+1), k, v)
		}
		return str
	}(ped.ExtraData))*/
	str += fmt.Sprintf("%sname: '%s'\n", ws, ped.Name)
	str += fmt.Sprintf("%spatternId: '%s'\n", ws, ped.PatternId)
	str += fmt.Sprintf("%spatternType: '%s'\n", ws, ped.PatternType)
	str += fmt.Sprintf("%spermissionsRequired: \n%s", ws, func(l []string) string {
		var str string
		for k, v := range l {
			str += fmt.Sprintf("%s%d. '%s'\n", supportingfunctions.GetWhitespace(num+1), k+1, v)
		}
		return str
	}(ped.PermissionsRequired))
	str += fmt.Sprintf("%splatforms: \n%s", ws, func(l []string) string {
		var str string
		for k, v := range l {
			str += fmt.Sprintf("%s%d. '%s'\n", supportingfunctions.GetWhitespace(num+1), k+1, v)
		}
		return str
	}(ped.Platforms))
	str += fmt.Sprintf("%sremoteSupport: '%v'\n", ws, ped.RemoteSupport)
	str += fmt.Sprintf("%srevoked: '%v'\n", ws, ped.Revoked)
	/*str += fmt.Sprintf("%ssystemRequirements: \n%s", ws, func(l []string) string {
		var str string
		for k, v := range l {
			str += fmt.Sprintf("%s%d. '%s'\n", supportingfunctions.GetWhitespace(num+1), k+1, v)
		}
		return str
	}(ped.SystemRequirements))*/
	str += fmt.Sprintf("%stactics: \n%s", ws, func(l []string) string {
		var str string
		for k, v := range l {
			str += fmt.Sprintf("%s%d. '%s'\n", supportingfunctions.GetWhitespace(num+1), k+1, v)
		}
		return str
	}(ped.Tactics))
	str += fmt.Sprintf("%sURL: '%s'\n", ws, ped.URL)
	str += fmt.Sprintf("%sversion: '%s'\n", ws, ped.Version)

	return str
}

// Get возвращает значения CustomFieldStringType, где 1 и 3 значение это
// наименование поля
func (cf *CustomFieldStringType) Get() (string, int, string, string) {
	return "order", cf.Order, "string", cf.String
}

// Set устанавливает значения CustomFieldStringType
func (cf *CustomFieldStringType) Set(order int, value interface{}) {
	cf.Order = order
	cf.String = fmt.Sprint(value)
}

// Get возвращает значения CustomFieldDateType, где 1 и 3 значение это
// наименование поля
func (cf *CustomFieldDateType) Get() (string, int, string, string) {
	t := time.UnixMilli(int64(cf.Date))

	return "order", cf.Order, "date", t.String()
}

// Set устанавливает значения CustomFieldDateType, при этом
// значение value должно быть типа uint64
func (cf *CustomFieldDateType) Set(order int, value interface{}) {
	cf.Order = order
	if v, ok := value.(uint64); ok {
		cf.Date = v
	}
}

/*
	Остались следующие типы и проверить

// TtpsMessage список TTP сообщений
type TtpsMessage struct {
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
*/

/*
"reports": {
                "Drill_1_2": {
                    "taxonomies": [
                        {
                            "level": "info",
                            "namespace": "Drill",
                            "predicate": "Score",
                            "value": 0
                        }
                    ]
                },
                "Moloch_1_7": {
                    "taxonomies": [
                        {
                            "level": "info",
                            "namespace": "Moloch",
                            "predicate": "Uploading",
                            "value": "1 pcap(s)"
                        }
                    ]
                }
            },

type ObservableMessage struct {
+	CreatedAt        uint64                              `json:"_createdAt"`
+	CreatedBy        string                              `json:"_createdBy"`
+	UnderliningId    string                              `json:"_id"`
+	UnderliningType  string                              `json:"_type"`
+	UpdatedAt        uint64                              `json:"_updatedAt"`
+	UpdatedBy        string                              `json:"_updatedBy"`
+	Data             string                              `json:"data"`
+	DataType         string                              `json:"dataType"`
+	IgnoreSimilarity bool                                `json:"ignoreSimilarity"`
+	ExtraData        map[string]interface{}              `json:"extraData"`
+	Ioc              bool                                `json:"ioc"`
+	Message          string                              `json:"message"`
+	Sighted          bool                                `json:"sighted"`
+	StartDate        uint64                              `json:"startDate"`
+	Tags             []string                            `json:"tags"`
+	Tlp              int                                 `json:"tlp"`
	Reports          map[string]map[string][]interface{} `json:"reports"`
}*/

// CustomFields настраиваемые поля
//type CustomFields map[string]map[string]interface{}
/*
// ToStringBeautiful выполняет красивое представление информации содержащейся в типе
func (istix IntrusionSetDomainObjectsSTIX) ToStringBeautiful() string {
	str := istix.CommonPropertiesObjectSTIX.ToStringBeautiful()
	str += istix.CommonPropertiesDomainObjectSTIX.ToStringBeautiful()
	str += fmt.Sprintf("name: '%s'\n", istix.Name)
	str += fmt.Sprintf("description: '%s'\n", istix.Description)
	str += fmt.Sprintf("aliases: \n%v", func(l []string) string {
		var str string
		for k, v := range l {
			str += fmt.Sprintf("\taliase '%d': '%s'\n", k, v)
		}
		return str
	}(istix.Aliases))
	str += fmt.Sprintf("first_seen: '%v'\n", istix.FirstSeen)
	str += fmt.Sprintf("last_seen: '%v'\n", istix.LastSeen)
	str += fmt.Sprintf("goals: \n%v", func(l []string) string {
		var str string
		for k, v := range l {
			str += fmt.Sprintf("\tgoal '%d': '%s'\n", k, v)
		}
		return str
	}(istix.Goals))
	str += fmt.Sprintf("resource_level: '%s'\n", istix.FirstSeen)
	str += fmt.Sprintf("primary_motivation: '%s'\n", istix.LastSeen)
	str += fmt.Sprintf("secondary_motivations: \n%v", func(l []OpenVocabTypeSTIX) string {
		var str string
		for k, v := range l {
			str += fmt.Sprintf("\tsecondary_motivation '%d': '%v'\n", k, v)
		}
		return str
	}(istix.SecondaryMotivations))

	return str
}

*/
