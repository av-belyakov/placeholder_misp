package coremodule

import (
	"fmt"
	"regexp"
	"strings"
)

// CheckMISPObservablesTag проверяет соответствие "observables.tags" значению
// соответствующему шаблону misp тегов
func CheckMISPObservablesTag(tag string) ([2]string, error) {
	nl := [2]string{}
	pattern := regexp.MustCompile(`^misp:([\w\-].*)=\"([\w\-].*)\"$`)

	if !pattern.MatchString(tag) {
		return nl, fmt.Errorf("the accepted value does not match the regular expression")
	}

	result := pattern.FindAllStringSubmatch(tag, -1)

	if len(result) > 0 && len(result[0]) == 3 {
		nl = [2]string{result[0][1], result[0][2]}
	}

	return nl, nil
}

func GetTypeNameObservablesTag(tag string) string {
	if !strings.Contains(tag, ":") {
		return tag
	}

	patterns := []*regexp.Regexp{
		regexp.MustCompile(`^type:([\w\-].*)$`),

		//это просто для примера с целью дальнейшего расширения
		regexp.MustCompile(`^example:([\w\-].*)$`),
	}

	for _, v := range patterns {
		if !v.MatchString(tag) {
			continue
		}

		result := v.FindAllStringSubmatch(tag, -1)

		if len(result) > 0 && len(result[0]) >= 2 {
			return result[0][1]
		}
	}

	return ""
}
