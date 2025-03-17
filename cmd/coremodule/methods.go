package coremodule

import (
	"strings"
)

//***************** хранилище данных ******************

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

// ********************* генератор объектов MISP **********************
// NewGenerateObjectsFormatMISP новый генератор объектов MISP
func NewGenerateObjectsFormatMISP(settings SettingsGenerateObjectsFormatMISP) *GenerateObjectsFormatMISP {
	return &GenerateObjectsFormatMISP{
		mispModule:    settings.MispModule,
		sqlite3Module: settings.Sqlite3Module,
		listRule:      settings.ListRule,
		counter:       settings.Counter,
		logger:        settings.Logger,
	}
}

//********************* правила обработки сообщений *********************

func NewExclusionRules() *ExclusionRules {
	return &ExclusionRules{}
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

func getObjName(objName string) string {
	l := strings.Split(objName, ".")

	if len(l) == 0 {
		return ""
	}

	return l[0]
}

//****************** временный MISP формат *********************

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
