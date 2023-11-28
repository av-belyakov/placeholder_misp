package coremodule

import (
	"fmt"
	"regexp"
)

// CheckObservablesTag проверяет соответствие "observables.tags" определенному значению
func CheckObservablesTag(tag string) ([2]string, error) {
	nl := [2]string{}
	patter := regexp.MustCompile(`^misp:([\w\-].*)=\"([\w\-].*)\"$`)

	if !patter.MatchString(tag) {
		return nl, fmt.Errorf("the accepted value does not match the regular expression")
	}

	result := patter.FindAllStringSubmatch(tag, -1)

	if len(result) > 0 && len(result[0]) == 3 {
		nl = [2]string{result[0][1], result[0][2]}
	}

	return nl, nil
}
