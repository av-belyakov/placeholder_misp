package rules

import (
	"fmt"
	"path"
	"runtime"

	"github.com/spf13/viper"

	"placeholder_misp/supportingfunctions"
)

func GetRuleProcessedMISPMsg(workDir, fn string) (ListRulesProcessedMISPMessage, error) {
	fmt.Println("func 'GetRuleProcessedMISPMsg', START")

	lrp := ListRulesProcessedMISPMessage{}

	rootPath, err := supportingfunctions.GetRootPath("placeholder_misp")
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return lrp, fmt.Errorf("%s %s:%d", fmt.Sprint(err), f, l+1)
	}

	viper.SetConfigFile(path.Join(rootPath, workDir, "processedmispmsg.yaml"))
	viper.SetConfigType("yaml")
	err = viper.ReadInConfig()
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return lrp, fmt.Errorf("%s %s:%d", fmt.Sprint(err), f, l+1)
	}

	if ok := viper.IsSet("RULLES"); !ok {
		_, f, l, _ := runtime.Caller(0)
		return lrp, fmt.Errorf("the 'RULLES' property is missing in the file '%s' %s:%d", fn, f, l+1)
	}

	err = viper.GetViper().Unmarshal(&lrp)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return lrp, fmt.Errorf("%s %s:%d", fmt.Sprint(err), f, l+1)
	}

	return lrp, nil
}
