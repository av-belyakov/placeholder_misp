package coremodule

import (
	"strings"
)

// ChanInputCreateMispFormat
// UUID - уникальный идентификатор в формате UUID
// FieldName - наименование поля
// ValueType - тип передаваемого значения (string, int и т.д.)
// Value - любые передаваемые данные
// FieldBranch - 'путь' до значения в как в JSON формате, например 'event.details.customFields.class'
type ChanInputCreateMispFormat struct {
	UUID        string
	FieldName   string
	ValueType   string
	Value       interface{}
	FieldBranch string
}

type FieldsNameMapping struct {
	InputFieldName, MispFieldName string
}

// storageValueName временное хранилище свойств элементов observables
type storageValueName []string

func NewStorageValueName() *storageValueName {
	return &storageValueName{}
}

func (svn *storageValueName) SetValueName(value string) {
	*svn = append(*svn, value)
}

func (svn *storageValueName) GetValueName(value string) bool {
	for _, v := range *svn {
		if v == value {
			return true
		}
	}

	return false
}

func (svn *storageValueName) CleanValueName() {
	*svn = storageValueName{}
}

// ExclusionRules содержит информацию об объектах которые нужно исключить из
// передачи в MISP
type ExclusionRules []ExclusionRule

// ExclusionRule
// SequenceNumber - порядковый номер в списке объектов
// NameList - наименование объекта
type ExclusionRule struct {
	SequenceNumber int
	NameList       string
}

// Add добавляет информацию об объекте подлежащего исключению из списка передаваемых в MISP
func (er *ExclusionRules) Add(sn int, n string) {
	var isExist bool
	n = getObjName(n)

	for _, v := range *er {
		if v.SequenceNumber == sn && v.NameList == n {
			isExist = true

			break
		}
	}

	if !isExist {
		*er = append(*er, ExclusionRule{SequenceNumber: sn, NameList: n})
	}
}

// SearchObjectName ищет, и возвращает список с информацией об объектах которые нужно
// исключить из передачи в MISP
func (er *ExclusionRules) SearchObjectName(objName string) []ExclusionRule {
	list := []ExclusionRule{}

	for _, v := range *er {
		if v.NameList != getObjName(objName) {
			continue
		}

		list = append(list, v)
	}

	return list
}

// SearchSeqNum ищет, и возвращает список с информацией об объектах которые нужно
// исключить из передачи в MISP
func (er *ExclusionRules) SearchSeqNum(sn int) []ExclusionRule {
	list := []ExclusionRule{}

	for _, v := range *er {
		if v.SequenceNumber != sn {
			continue
		}

		list = append(list, v)
	}

	return list
}

// Clean выполняет очистку списка
func (er *ExclusionRules) Clean() {
	er = &ExclusionRules{}
}

func NewExclusionRules() *ExclusionRules {
	return &ExclusionRules{}
}

func getObjName(objName string) string {
	l := strings.Split(objName, ".")

	if len(l) == 0 {
		return ""
	}

	return l[0]
}

type MispGalaxyTags map[int]MispGalaxyOptions

type MispGalaxyOptions struct {
	Name, PatternId, PatternType string
}

func NewMispGalaxyTags() *MispGalaxyTags {
	return &MispGalaxyTags{}
}

func (gt *MispGalaxyTags) Get() MispGalaxyTags {
	return *gt
}

func (gt *MispGalaxyTags) SetPatternId(key int, value string) {
	options := gt.GetGalaxyOptions(key)
	options.PatternId = value
	(*gt)[key] = options
}

func (gt *MispGalaxyTags) SetPatternType(key int, value string) {
	options := gt.GetGalaxyOptions(key)
	options.PatternType = value
	(*gt)[key] = options
}

func (gt *MispGalaxyTags) SetName(key int, value string) {
	options := gt.GetGalaxyOptions(key)
	options.Name = value
	(*gt)[key] = options
}

func (gt *MispGalaxyTags) GetGalaxyOptions(num int) MispGalaxyOptions {
	if elem, ok := (*gt)[num]; ok {
		return elem
	}

	return MispGalaxyOptions{}
}
