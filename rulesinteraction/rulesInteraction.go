// Пакет rules выполняет чтение списка специализированных правил
package rules

import (
	"fmt"
	"path"
	"runtime"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// NewListRule создает новый список правил
func NewListRule(rootDir, workDir, fileName string) (*ListRule, []string, error) {
	lr := ListRule{}

	_, f, l, _ := runtime.Caller(0)
	rootPath, err := getRootPath(rootDir)
	if err != nil {
		return &lr, []string{}, fmt.Errorf("'%s' %s:%d", err.Error(), f, l+1)
	}

	viper.SetConfigFile(path.Join(rootPath, workDir, fileName))
	viper.SetConfigType("yaml")

	_, f, l, _ = runtime.Caller(0)
	err = viper.ReadInConfig()
	if err != nil {
		return &lr, []string{}, fmt.Errorf("'%s' %s:%d", err.Error(), f, l+1)
	}

	_, f, l, _ = runtime.Caller(0)
	if ok := viper.IsSet("RULES"); !ok {
		return &lr, []string{}, fmt.Errorf("'the \"RULES\" property is missing in the file \"%s\"' %s:%d", fileName, f, l+1)
	}

	_, f, l, _ = runtime.Caller(0)
	err = viper.GetViper().Unmarshal(&lr, func(dc *mapstructure.DecoderConfig) {
		dc.Squash = true
	})
	if err != nil {
		return &lr, []string{}, fmt.Errorf("'%s' %s:%d", err.Error(), f, l+1)
	}

	warningCheckRules := lr.verification()

	return &lr, warningCheckRules, nil
}
