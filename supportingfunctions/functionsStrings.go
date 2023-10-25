package supportingfunctions

import (
	"bufio"
	"os"
	"regexp"
)

func GetWhitespace(num int) string {
	var str string

	if num == 0 {
		return str
	}

	for i := 0; i < num; i++ {
		str += "  "
	}

	return str
}

func GetAppName(pf string, nl int) (string, error) {
	var line string

	f, err := os.OpenFile(pf, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return line, err
	}
	defer f.Close()

	num := 1
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		if num == nl {
			return sc.Text(), nil
		}

		num++
	}

	return line, nil
}

func GetAppVersion(str string) string {
	version := "версия не определена"
	patter := regexp.MustCompile(`v(\d)+\.(\d)+.(\d)+`)
	ls := patter.FindStringSubmatch(str)

	if len(ls) > 0 {
		version = ls[0]
	}

	return version
}

func CheckHashSum(hsum string) string {
	switch len(hsum) {
	case 32:
		return "md5"
	case 40:
		return "sha1"
	case 64:
		return "sha256"
	}

	return "other"
}
