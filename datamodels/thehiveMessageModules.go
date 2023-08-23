package datamodels

import (
	"fmt"

	"placeholder_misp/supportingfunctions"
)

// NewResponseMessage формирует новый тип ResponseMessageFromMispToTheHave с предустановленными значениями
func NewResponseMessage() *ResponseMessageFromMispToTheHave {
	return &ResponseMessageFromMispToTheHave{
		Success: true,
		Service: "MISP",
		Commands: []ResponseCommandForTheHive{
			{Command: "addtag", String: "Webhook: send=\"MISP\""},
		},
	}
}

func (rm *ResponseMessageFromMispToTheHave) ResponseMessageAddNewCommand(rcm ResponseCommandForTheHive) {
	rm.Commands = append(rm.Commands, rcm)
}

func (rm *ResponseMessageFromMispToTheHave) GetResponseMessageFromMispToTheHave() ResponseMessageFromMispToTheHave {
	return *rm
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

func (sm SourceMessageTheHive) ToStringBeautiful(num int) string {
	return fmt.Sprintf("source: '%s'\n", sm.Source)
}

func (em EventMessageTheHive) ToStringBeautiful(num int) string {
	var str string

	ws := supportingfunctions.GetWhitespace(num)

	str += fmt.Sprintf("%soperation: '%s'\n", ws, em.Operation)
	str += fmt.Sprintf("%sobjectId: '%s'\n", ws, em.ObjectId)
	str += fmt.Sprintf("%sobjectType: '%s'\n", ws, em.ObjectType)
	str += fmt.Sprintf("%sbase: '%v'\n", ws, em.Base)
	str += fmt.Sprintf("%sstartDate: '%d'\n", ws, em.StartDate)
	str += fmt.Sprintf("%srootId: '%s'\n", ws, em.RootId)
	str += fmt.Sprintf("%srequestId: '%s'\n", ws, em.RequestId)
	str += fmt.Sprintf("%sdetails:\n", ws)
	str += em.Details.ToStringBeautiful(num + 1)
	str += fmt.Sprintf("%sobject:\n", ws)
	str += em.Object.ToStringBeautiful(num + 1)
	str += fmt.Sprintf("%sorganisationId: '%s'\n", ws, em.OrganisationId)
	str += fmt.Sprintf("%sorganisation: '%s'\n", ws, em.Organisation)

	return str
}

func (ed EventDetails) ToStringBeautiful(num int) string {
	var str string

	ws := supportingfunctions.GetWhitespace(num)

	str += fmt.Sprintf("%sendDate: '%d'\n", ws, ed.EndDate)
	str += fmt.Sprintf("%sresolutionStatus: '%s'\n", ws, ed.ResolutionStatus)
	str += fmt.Sprintf("%ssummary: '%s'\n", ws, ed.Summary)
	str += fmt.Sprintf("%sstatus: '%s'\n", ws, ed.Status)
	str += fmt.Sprintf("%simpactStatus: '%s'\n", ws, ed.ImpactStatus)
	str += ed.CustomFields.ToStringBeautiful(num)

	return str
}

func (eo EventObject) ToStringBeautiful(num int) string {
	var str string

	ws := supportingfunctions.GetWhitespace(num)

	str += fmt.Sprintf("%s_id: '%s'\n", ws, eo.UnderliningId)
	str += fmt.Sprintf("%sid: '%s'\n", ws, eo.Id)
	str += fmt.Sprintf("%screatedBy: '%s'\n", ws, eo.CreatedBy)
	str += fmt.Sprintf("%supdatedBy: '%s'\n", ws, eo.UpdatedBy)
	str += fmt.Sprintf("%screatedAt: '%d'\n", ws, eo.CreatedAt)
	str += fmt.Sprintf("%supdatedAt: '%d'\n", ws, eo.UpdatedAt)
	str += fmt.Sprintf("%s_type: '%s'\n", ws, eo.UnderliningType)
	str += fmt.Sprintf("%scaseId: '%d'\n", ws, eo.CaseId)
	str += fmt.Sprintf("%stitle: '%s'\n", ws, eo.Title)
	str += fmt.Sprintf("%sdescription: '%s'\n", ws, eo.Description)
	str += fmt.Sprintf("%sseverity: '%d'\n", ws, eo.Severity)
	str += fmt.Sprintf("%sstartDate: '%d'\n", ws, eo.StartDate)
	str += fmt.Sprintf("%sendDate: '%d'\n", ws, eo.EndDate)
	str += fmt.Sprintf("%simpactStatus: '%s'\n", ws, eo.ImpactStatus)
	str += fmt.Sprintf("%sresolutionStatus: '%s'\n", ws, eo.ResolutionStatus)
	str += fmt.Sprintf("%stags: \n%s", ws, func(l []string) string {
		var str string
		ws := supportingfunctions.GetWhitespace(num + 1)

		for k, v := range l {
			str += fmt.Sprintf("%s%d. '%s'\n", ws, k+1, v)
		}
		return str
	}(eo.Tags))
	str += fmt.Sprintf("%sflag: '%v'\n", ws, eo.Flag)
	str += fmt.Sprintf("%stlp: '%d'\n", ws, eo.Tlp)
	str += fmt.Sprintf("%spap: '%d'\n", ws, eo.Pap)
	str += fmt.Sprintf("%sstatus: '%s'\n", ws, eo.Status)
	str += fmt.Sprintf("%ssummary: '%s'\n", ws, eo.Summary)
	str += fmt.Sprintf("%sowner: '%s'\n", ws, eo.Owner)
	str += eo.CustomFields.ToStringBeautiful(num)
	str += fmt.Sprintf("%sstats: \n%s", ws, func(l map[string]interface{}) string {
		var str string
		ws := supportingfunctions.GetWhitespace(num + 1)

		for k, v := range l {
			str += fmt.Sprintf("%s%s: '%v'\n", ws, k, v)
		}
		return str
	}(eo.Stats))
	str += fmt.Sprintf("%spermissions: \n%s", ws, func(l []string) string {
		var str string
		ws := supportingfunctions.GetWhitespace(num + 1)

		for k, v := range l {
			str += fmt.Sprintf("%s%d. '%s'\n", ws, k+1, v)
		}
		return str
	}(eo.Permissions))

	return str
}

func (cf CustomFields) ToStringBeautiful(num int) string {
	var str string
	ws := supportingfunctions.GetWhitespace(num)

	str += fmt.Sprintf("%scustomFields:\n", ws)
	for key, value := range cf {
		str += fmt.Sprintf("%s%s:\n", supportingfunctions.GetWhitespace(num+1), key)
		for k, v := range value {
			str += fmt.Sprintf("%s%s: '%v'\n", supportingfunctions.GetWhitespace(num+2), k, v)
		}
	}

	return str
}

func (om ObservablesMessageTheHive) ToStringBeautiful(num int) string {
	var str string

	for _, v := range om.Observables {
		str += v.ToStringBeautiful(num)
	}

	return str
}

func (om ObservableMessage) ToStringBeautiful(num int) string {
	var str string
	ws := supportingfunctions.GetWhitespace(num)

	str += fmt.Sprintf("%s_createdAt: '%d'\n", ws, om.CreatedAt)
	str += fmt.Sprintf("%s_createdBy: '%s'\n", ws, om.CreatedBy)
	str += fmt.Sprintf("%s_id: '%s'\n", ws, om.UnderliningId)
	str += fmt.Sprintf("%s_type: '%s'\n", ws, om.UnderliningType)
	str += fmt.Sprintf("%s_updatedAt: '%d'\n", ws, om.UpdatedAt)
	str += fmt.Sprintf("%s_updatedBy: '%s'\n", ws, om.UpdatedBy)
	str += fmt.Sprintf("%sdata: '%s'\n", ws, om.Data)
	str += fmt.Sprintf("%sdataType: '%s'\n", ws, om.DataType)
	str += fmt.Sprintf("%signoreSimilarity: '%v'\n", ws, om.IgnoreSimilarity)
	str += fmt.Sprintf("%sextraData: \n%s", ws, func(l map[string]interface{}) string {
		var str string
		ws := supportingfunctions.GetWhitespace(num + 1)

		for k, v := range l {
			str += fmt.Sprintf("%s%s: '%v'\n", ws, k, v)
		}
		return str
	}(om.ExtraData))
	str += fmt.Sprintf("%sioc: '%v'\n", ws, om.Ioc)
	str += fmt.Sprintf("%smessage: '%s'\n", ws, om.Message)
	str += fmt.Sprintf("%ssighted: '%v'\n", ws, om.Sighted)
	str += fmt.Sprintf("%sstartDate: '%d'\n", ws, om.StartDate)
	str += fmt.Sprintf("%stags: \n%s", ws, func(l []string) string {
		var str string
		ws := supportingfunctions.GetWhitespace(num + 1)

		for k, v := range l {
			str += fmt.Sprintf("%s%d. '%s'\n", ws, k+1, v)
		}
		return str
	}(om.Tags))
	str += fmt.Sprintf("%stlp: '%d'\n", ws, om.Tlp)
	str += fmt.Sprintf("%sreports: \n%s", ws, func(l map[string]map[string][]map[string]interface{}) string {
		var str string
		for key, value := range l {
			str += fmt.Sprintf("%s%s:\n", supportingfunctions.GetWhitespace(num+1), key)
			for k, v := range value {
				str += fmt.Sprintf("%s%s:\n", supportingfunctions.GetWhitespace(num+2), k)
				for i, j := range v {
					str += fmt.Sprintf("%s%d.\n", supportingfunctions.GetWhitespace(num+3), i+1)
					for n, m := range j {
						str += fmt.Sprintf("%s%s: %v\n", supportingfunctions.GetWhitespace(num+4), n, m)
					}
				}
			}
		}
		return str
	}(om.Reports))

	return str
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
	str += fmt.Sprintf("%sdefenseBypassed: \n%v", ws, func(l []string) string {
		var str string
		for k, v := range l {
			str += fmt.Sprintf("%s%d. '%s'\n", supportingfunctions.GetWhitespace(num+1), k+1, v)
		}
		return str
	}(ped.DefenseBypassed))
	str += fmt.Sprintf("%sdescription: '%s'\n", ws, ped.Description)
	str += fmt.Sprintf("%sextraData: \n%s", ws, func(l map[string]interface{}) string {
		var str string
		for k, v := range l {
			str += fmt.Sprintf("%s%s: '%v'\n", supportingfunctions.GetWhitespace(num+1), k, v)
		}
		return str
	}(ped.ExtraData))
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
	str += fmt.Sprintf("%ssystemRequirements: \n%s", ws, func(l []string) string {
		var str string
		for k, v := range l {
			str += fmt.Sprintf("%s%d. '%s'\n", supportingfunctions.GetWhitespace(num+1), k+1, v)
		}
		return str
	}(ped.SystemRequirements))
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
