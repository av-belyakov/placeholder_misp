package datamodels

import (
	"strings"
	"sync"

	"placeholder_misp/supportingfunctions"
)

type ListAttributeTmp struct {
	attributes map[int][]AttributeMispFormat
	sync.Mutex
}

func NewListAttributeTmp() *ListAttributeTmp {
	return &ListAttributeTmp{
		attributes: map[int][]AttributeMispFormat{},
	}
}

func (la *ListAttributeTmp) AddAttribute(branch string, value interface{}, num int) {
	t := "other"
	objr := "other"

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
			t = supportingfunctions.CheckHashSum(str)
			objr = "hashsum"
		}

		tmp = append(tmp, AttributeMispFormat{
			Category:       "Payload delivery",
			Distribution:   "0",
			Value:          str,
			Type:           t,
			ObjectRelation: objr,
		})
	}

	la.attributes[num] = tmp
}

func (la *ListAttributeTmp) GetListAttribute() map[int][]AttributeMispFormat {
	return la.attributes
}

func (la *ListAttributeTmp) CleanAttribute() {
	la.Lock()
	defer la.Unlock()

	la.attributes = map[int][]AttributeMispFormat{}
}

func createNewAttributeMisp() []AttributeMispFormat {
	return []AttributeMispFormat{}
}
