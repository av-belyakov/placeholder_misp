package appversion

import "os"

// GetVersion версия приложения
func GetAppVersion() (string, error) {
	b, err := os.ReadFile("version")
	if err != nil {
		return "", err
	}

	return string(b), nil
}
