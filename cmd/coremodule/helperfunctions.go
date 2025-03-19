package coremodule

import (
	"fmt"

	"slices"

	"github.com/av-belyakov/objectsmispformat"
	"github.com/av-belyakov/placeholder_misp/commoninterfaces"
	rules "github.com/av-belyakov/placeholder_misp/internal/ruleshandler"
)

// удаляет элемент из списка атрибутов
func delElementAttributes(er *ExclusionRules, la *objectsmispformat.ListAttributesMispFormat, logger commoninterfaces.Logger) {
	for _, v := range er.SearchObjectName("observables") {
		if attr, ok := la.DelElementList(v.SequenceNumber); ok {
			logger.Send("warning", fmt.Sprintf("'an attribute with a value of '%s' has been removed'", attr.Value))
		}
	}
}

// обрабатывает значение свойства observables.tags
func handlerObservablesTags(v interface{},
	listTags map[int][2]string,
	listAttributesMisp *objectsmispformat.ListAttributesMispFormat,
	seqNumObservable int) map[int][2]string {
	tag, ok := v.(string)
	if !ok {
		return listTags
	}

	//проверка значения на соответствию определенному шаблону
	//начинающемуся на misp: при этом значения целиком берутся из
	//этого шаблона
	result, err := CheckMISPObservablesTag(tag)
	if err == nil {
		listTags[seqNumObservable] = result
	} else {
		result := GetTypeNameObservablesTag(tag)
		if result == "" {
			return listTags
		}
		//добавляем значение из tags в поле object_relation
		listAttributesMisp.SetAnyObjectRelation(result, seqNumObservable)

		//Добавляем в свойство Category соответствующее значение
		//если наименование похоже на наименование типа хеширования
		if checkHashName(result) {
			listAttributesMisp.AutoSetValueCategory(result, seqNumObservable)
		}
	}

	return listTags
}

// SupportiveExcludeRuleTmp тестовый
type SupportiveExcludeRuleTmp struct {
	rules   map[int][]SupportiveExcludeRuleOptionsListAnd
	pattern []rules.ExcludeListAnd
}

type SupportiveExcludeRuleOptionsListAnd []SupportiveExcludeRuleOptionsTmp

type SupportiveExcludeRuleOptionsTmp struct {
	fieldName string
	isEqual   bool
}

func NewSupportiveListExcludeRuleTmp(listRules []rules.ExcludeListAnd) *SupportiveExcludeRuleTmp {
	return &SupportiveExcludeRuleTmp{
		rules:   make(map[int][]SupportiveExcludeRuleOptionsListAnd),
		pattern: listRules,
	}
}

func (ser *SupportiveExcludeRuleTmp) Get() map[int][]SupportiveExcludeRuleOptionsListAnd {
	return ser.rules
}

func (ser *SupportiveExcludeRuleTmp) Add(num int, fieldName string, searchValue any, addrRule [2]int, isEqual bool) {
	//если нет совпадений с правилами
	if !isEqual {
		return
	}

	if _, ok := ser.rules[num]; !ok {
		ser.rules[num] = ser.createPattern()
	}

	for k, v := range ser.rules[num] {
		for key, value := range v {
			if value.fieldName != fieldName {
				continue
			}

			if k != addrRule[0] || key != addrRule[1] {
				continue
			}

			ser.rules[num][k][key].isEqual = true
		}
	}
}

func (ser *SupportiveExcludeRuleTmp) CheckRuleTrue(num int) bool {
	r, ok := ser.rules[num]
	if !ok {
		return false
	}

	searchEqual := func(list SupportiveExcludeRuleOptionsListAnd) bool {
		for _, v := range list {
			if !v.isEqual {
				return false
			}
		}

		return true
	}

	listTmp := make([]bool, 0, len(r))
	for _, v := range r {
		listTmp = append(listTmp, searchEqual(v))
	}

	for _, v := range listTmp {
		if v {
			return true
		}
	}

	return false
}

func (ser *SupportiveExcludeRuleTmp) createPattern() []SupportiveExcludeRuleOptionsListAnd {
	result := []SupportiveExcludeRuleOptionsListAnd(nil)

	for _, v := range ser.pattern {
		listTmp := []SupportiveExcludeRuleOptionsTmp(nil)

		for _, value := range v.ListAnd {
			listTmp = append(listTmp, SupportiveExcludeRuleOptionsTmp{fieldName: value.SearchField})
		}

		result = append(result, listTmp)
	}

	return result
}

// ---------------
// SupportiveExcludeRule вспомогательный тип для обработки правил типа Exclude (исключение)
type SupportiveExcludeRule struct {
	rules map[int][]SupportiveExcludeRuleOptions
}

type SupportiveExcludeRuleOptions struct {
	fieldName string
	value     string
	isEqual   bool
}

func NewSupportiveListExcludeRule() *SupportiveExcludeRule {
	return &SupportiveExcludeRule{
		rules: make(map[int][]SupportiveExcludeRuleOptions),
	}
}

func (ser *SupportiveExcludeRule) Get() map[int][]SupportiveExcludeRuleOptions {
	return ser.rules
}

func (ser *SupportiveExcludeRule) Add(num int, fieldName, value string, isEqual bool) {
	//если нет совпадений с правилами
	if !isEqual {
		return
	}

	if _, ok := ser.rules[num]; !ok {
		ser.rules[num] = []SupportiveExcludeRuleOptions(nil)
	}

	if fieldName == "observables.data" {
		fmt.Println("func 'Add' 333 value:", value)
	}

	ser.rules[num] = append(ser.rules[num], SupportiveExcludeRuleOptions{
		fieldName: fieldName,
		value:     value,
		isEqual:   isEqual,
	})
}

func (ser *SupportiveExcludeRule) CheckRuleTrue(num int) bool {
	if _, ok := ser.rules[num]; !ok {
		fmt.Println("func 'CheckRuleTrue', num =", num)

		return false
	}

	var isTrue = true
	for _, value := range ser.rules[num] {
		if !value.isEqual {
			isTrue = false

			break
		}
	}

	return isTrue

	/*	for k, v := range ser.rules {
			var isTrue = true
			for _, value := range v {
				if !value.isEqual {
					isTrue = false

					break
				}
			}

			if isTrue {
				return k, true
			}
		}

		return false*/
}

// searchEventSource выполняет поиск источника события
func searchEventSource(tmf ChanInputCreateMispFormat) (string, bool) {
	var (
		source string
		ok     bool
	)

	if tmf.FieldBranch == "source" {
		source, ok = tmf.Value.(string)
	}

	return source, ok
}

// searchCaseId выполняет поиск id кейса
func searchCaseId(tmf ChanInputCreateMispFormat) (float64, bool) {
	var (
		cid float64
		ok  bool
	)

	if tmf.FieldBranch == "event.object.caseId" {
		cid, ok = tmf.Value.(float64)
	}

	return cid, ok
}

// searchOwnerEmail выполняет поиск email владельца события
func searchOwnerEmail(tmf ChanInputCreateMispFormat) (string, bool) {
	var (
		email string
		ok    bool
	)

	if tmf.FieldBranch == "event.object.owner" {
		email, ok = tmf.Value.(string)
	}

	return email, ok
}

// getNewListAttributes устанавливает значения в свойствах Category и Type и в том
// числе изменяет состояние свойства DisableCorrelation (которое по умолчанию ВЫКЛЮЧЕНО)
// на значение информирующее MISP нужно ли выполнять корреляцию или нет
func getNewListAttributes(al map[int]objectsmispformat.AttributesMispFormat, lat map[int][2]string) []*objectsmispformat.AttributesMispFormat {
	nal := make([]*objectsmispformat.AttributesMispFormat, 0, len(al))

	for k, v := range al {
		if elem, ok := lat[k]; ok {
			v.Category = elem[0]
			v.Type = elem[1]

			//очищаем данное свойство, являющееся вспомогательным, так как оно может
			//быть заполненно ранее, а приоритетными являются значения в Category и Type
			v.ObjectRelation = ""
		}

		//выключаем автоматическую коореляцию с другими событиями для MISP
		if (v.Type == "other" || v.Type == "snort") && (v.ObjectRelation == "" || v.ObjectRelation == "snort") {
			v.DisableCorrelation = true
		}

		nal = append(nal, &v)
	}

	return nal
}

// addListGalaxyTags заполняет mispGalaxyTags значениями при условии совпадения
// пути соответствующие следующим значениям:
// ttp.extraData.pattern.patternId - содержит что то похожее на T1036.005
// ttp.extraData.pattern.patternType - содержит что то похожее на attack-pattern
// ttp.extraData.pattern.name - наименование
func addListGalaxyTags(lgt *MispGalaxyTags) func(string, any) {
	var (
		num             int
		fieldBranchList []string
	)

	searchValue := func(list []string, search string) bool {
		for _, v := range list {
			if v == search {
				return true
			}
		}

		return false
	}

	return func(fieldBranch string, value any) {
		v := fmt.Sprint(value)
		if searchValue(fieldBranchList, fieldBranch) {
			num++
			fieldBranchList = []string{}
		}

		switch fieldBranch {
		case "ttp.extraData.pattern.patternId":
			lgt.SetPatternId(num, v)
			fieldBranchList = append(fieldBranchList, fieldBranch)

		case "ttp.extraData.pattern.patternType":
			lgt.SetPatternType(num, v)
			fieldBranchList = append(fieldBranchList, fieldBranch)

		case "ttp.extraData.pattern.name":
			lgt.SetName(num, v)
			fieldBranchList = append(fieldBranchList, fieldBranch)
		}
	}
}

// создает список тегов которые MISP использует для формирования галактик,
// теги меют подобную структуру:
// "misp-galaxy:mitre-attack-pattern=\"Match Legitimate Name or Location - T1036.005\""
func createGalaxyTags(list *MispGalaxyTags) []string {
	result := make([]string, 0, len(*list))

	for _, v := range *list {
		result = append(result, fmt.Sprintf("misp-galaxy:mitre-%s=\"%s - %s\"", v.PatternType, v.Name, v.PatternId))
	}

	return result
}

func getNewListObjects(
	listObjects map[int]objectsmispformat.ObjectsMispFormat,
	attachment map[int][]objectsmispformat.AttributeMispFormat) map[int]*objectsmispformat.ObjectsMispFormat {
	nlo := make(map[int]*objectsmispformat.ObjectsMispFormat, len(attachment))

	for k, v := range attachment {
		if obj, ok := listObjects[k]; ok {
			obj.Attribute = v
			nlo[k] = &obj
		}
	}

	return nlo
}

func joinEventTags(listTags *objectsmispformat.ListEventObjectTags, galaxyTags []string) {
	for _, v := range galaxyTags {
		listTags.SetTag(v)
	}
}

func checkHashName(name string) bool {
	list := []string{
		"md5",
		"sha1",
		"sha128",
		"sha224",
		"sha256",
		"sha384",
		"sha512",
		"ja3",
	}

	return slices.Contains(list, name)
}
