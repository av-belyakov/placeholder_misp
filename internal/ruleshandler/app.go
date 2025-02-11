// Пакет rules выполняет чтение списка специализированных правил
package ruleshandler

import (
	"fmt"
	"path/filepath"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"

	"github.com/av-belyakov/placeholder_misp/internal/supportingfunctions"
)

// NewListRule создает новый список правил
func NewListRule(rootDir, workDir, fileName string) (*ListRule, []string, error) {
	lr := ListRule{}

	rootPath, err := supportingfunctions.GetRootPath(rootDir)
	if err != nil {
		return &lr, []string{}, supportingfunctions.CustomError(err)
	}

	viper.SetConfigFile(filepath.Join(rootPath, workDir, fileName))
	viper.SetConfigType("yml")

	err = viper.ReadInConfig()
	if err != nil {
		return &lr, []string{}, supportingfunctions.CustomError(err)
	}

	if ok := viper.IsSet("RULES"); !ok {
		return &lr, []string{}, supportingfunctions.CustomError(fmt.Errorf("the 'RULES' property is missing in the file '%s'", fileName))
	}

	err = viper.GetViper().Unmarshal(&lr, func(dc *mapstructure.DecoderConfig) {
		dc.Squash = true
	})
	if err != nil {
		return &lr, []string{}, supportingfunctions.CustomError(err)
	}

	warningCheckRules := lr.verification()

	return &lr, warningCheckRules, nil
}
