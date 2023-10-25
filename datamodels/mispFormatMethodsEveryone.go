package datamodels

import (
	"strings"

	"placeholder_misp/supportingfunctions"
)

func NewListAttribute() ListAttribute {
	return ListAttribute{}
}

func (la *ListAttribute) AddAttribute(branch, value string) {
	t := "other"
	objr := "other"

	nameIsExist := strings.Contains(branch, "attachment.name")
	hashesIsExist := strings.Contains(branch, "attachment.hashes")

	if !nameIsExist && !hashesIsExist {
		return
	}

	if nameIsExist {
		t = "filename"
		objr = "filename"
	}

	if hashesIsExist {
		t = supportingfunctions.CheckHashSum(value)
		objr = "hashsum"
	}

	*la = append(*la, AttributeMispFormat{
		Category:       "Payload delivery",
		Distribution:   "0",
		Value:          value,
		Type:           t,
		ObjectRelation: objr,
	})
}

func (la *ListAttribute) GetListAttribute() ListAttribute {
	return *la
}

func (la *ListAttribute) DelAttribute() {
	la = &ListAttribute{}
}
