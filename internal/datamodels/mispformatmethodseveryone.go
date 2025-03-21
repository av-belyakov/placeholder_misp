package datamodels

import (
	"strings"
	"sync"

	"github.com/av-belyakov/placeholder_misp/internal/supportingfunctions"
)

// ListAttributeTmp временный список атрибутов
type ListAttributeTmp struct {
	attributes map[int][]AttributeMispFormat
	sync.Mutex
}

// NewListAttributeTmp генерирует временный список атрибутов
func NewListAttributeTmp() *ListAttributeTmp {
	return &ListAttributeTmp{
		attributes: map[int][]AttributeMispFormat{},
	}
}

// AddAttribute добавляет атрибуты в список атрибутов
func (la *ListAttributeTmp) AddAttribute(branch string, value interface{}, num int) {
	var (
		t                  string = "other"
		objr               string = "other"
		disableCorrelation bool

		err error
	)

	nameIsExist := strings.Contains(branch, "attachment.name")
	hashesIsExist := strings.Contains(branch, "attachment.hashes")

	if !nameIsExist && !hashesIsExist {
		return
	}

	var tmp []AttributeMispFormat
	la.Lock()
	defer la.Unlock()

	if attr, ok := la.attributes[num]; ok {
		tmp = attr
	} else {
		tmp = createNewAttributeMisp()
	}

	if nameIsExist {
		t = "filename"
		objr = "filename"
	}

	if str, ok := value.(string); ok {
		if hashesIsExist {
			if t, _, err = supportingfunctions.CheckStringHash(str); err == nil {
				objr = t
			}
		}

		if objr == "filename" || objr == "other" {
			disableCorrelation = true
		}

		tmp = append(tmp, AttributeMispFormat{
			Category:           "Payload delivery",
			Distribution:       "0",
			DisableCorrelation: disableCorrelation,
			Value:              str,
			Type:               t,
			ObjectRelation:     objr,
		})
	}

	la.attributes[num] = tmp
}

// GetListAttribute возвращает список атрибутов
func (la *ListAttributeTmp) GetListAttribute() map[int][]AttributeMispFormat {
	return la.attributes
}

// CleanAttribute очищает временный список атрибутов
func (la *ListAttributeTmp) CleanAttribute() {
	la.Lock()
	defer la.Unlock()

	la.attributes = map[int][]AttributeMispFormat{}
}

func createNewAttributeMisp() []AttributeMispFormat {
	return []AttributeMispFormat{}
}
