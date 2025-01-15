package testcreategalaxy_test

import (
	"bufio"
	"fmt"
	"os"
	"path"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/av-belyakov/placeholder_misp/coremodule"
	"github.com/av-belyakov/placeholder_misp/internal/datamodels"
	"github.com/av-belyakov/placeholder_misp/internal/supportingfunctions"
	rules "github.com/av-belyakov/placeholder_misp/rulesinteraction"
)

func addListGalaxyTags(lgt *coremodule.MispGalaxyTags) func(string, any) {
	var (
		num             int
		fieldBranchList []string
	)

	searchValue := func(list []string, search string) bool {
		for _, v := range list {
			if v == search {
				return true
			}
		}

		return false
	}

	return func(fieldBranch string, value any) {
		//if strings.Contains(fieldBranch, "ttp.extraData.") {
		//	fmt.Println("fieldBranch =", fieldBranch, " value =", value)
		//}

		v := fmt.Sprint(value)
		if searchValue(fieldBranchList, fieldBranch) {
			num++
			fieldBranchList = []string{}
		}

		switch fieldBranch {
		case "ttp.extraData.pattern.patternId":
			//fmt.Println("fieldBranch =", fieldBranch, " value =", value, "NUM:", num)
			//fmt.Println("fieldBranchList:", fieldBranchList)

			lgt.SetPatternId(num, v)
			fieldBranchList = append(fieldBranchList, fieldBranch)
		case "ttp.extraData.pattern.patternType":
			//fmt.Println("fieldBranch =", fieldBranch, " value =", value, "NUM:", num)
			//fmt.Println("fieldBranchList:", fieldBranchList)

			lgt.SetPatternType(num, v)
			fieldBranchList = append(fieldBranchList, fieldBranch)
		case "ttp.extraData.pattern.name":
			//fmt.Println("fieldBranch =", fieldBranch, " value =", value, "NUM:", num)
			//fmt.Println("fieldBranchList:", fieldBranchList)

			lgt.SetName(num, v)
			fieldBranchList = append(fieldBranchList, fieldBranch)
		}
	}
}

func createGalaxyTags(list *coremodule.MispGalaxyTags) []string {
	result := make([]string, 0, len(*list))

	for _, v := range *list {
		result = append(result, fmt.Sprintf("misp-galaxy:mitre-%s=\"%s - %s\"", v.PatternType, v.Name, v.PatternId))
	}

	return result
}

var _ = Describe("Createlistgalaxytags", Ordered, func() {
	var (
		exampleByte    []byte
		errReadFile    error
		errRules       error
		logging        chan datamodels.MessageLogging
		chanInput      chan ChanInputCreateMispFormat
		chanDone       chan struct{}
		stopped        chan bool
		listRules      *rules.ListRule
		listGalaxyTags *coremodule.MispGalaxyTags
	)

	readFileJson := func(fpath, fname string) ([]byte, error) {
		var newResult []byte

		rootPath, err := supportingfunctions.GetRootPath("placeholder_misp")
		if err != nil {
			return newResult, err
		}

		fmt.Println("func 'readFileJson', path = ", path.Join(rootPath, fpath, fname))

		f, err := os.OpenFile(path.Join(rootPath, fpath, fname), os.O_RDONLY, os.ModePerm)
		if err != nil {
			return newResult, err
		}
		defer f.Close()

		sc := bufio.NewScanner(f)
		for sc.Scan() {
			newResult = append(newResult, sc.Bytes()...)
		}

		return newResult, nil
	}

	BeforeAll(func() {
		logging = make(chan datamodels.MessageLogging)
		chanInput = make(chan ChanInputCreateMispFormat)
		chanDone = make(chan struct{})
		stopped = make(chan bool)

		listGalaxyTags = coremodule.NewMispGalaxyTags()

		exampleByte, errReadFile = readFileJson("testing/test_json", "example_caseId_33705.json")

		//инициализация списка правил
		listRules, _, errRules = rules.NewListRule("placeholder_misp", "rules", "mispmsgrule.yaml")

		go func(logging <-chan datamodels.MessageLogging) {
			for {
				select {
				case msg := <-logging:
					fmt.Println("LOG MSG:", msg.MsgData)
				case <-chanDone:
					fmt.Println("---=== STOPED DECODE JSON OBJECT ===---")

					return
				}
			}
		}(logging)

		go DecodeJsonObject(chanInput, exampleByte, listRules, logging, stopped)
	})

	Context("Тест 1. Проверка инициализации модулей", func() {
		It("При инициализации модуля чтения файла примера не должно быть ошибки", func() {
			Expect(errReadFile).ShouldNot(HaveOccurred())
		})

		It("При чтении файла правил ошибки быть не должно", func() {
			Expect(errRules).ShouldNot(HaveOccurred())
		})
	})

	Context("Тест 2. Формирование списка тегов галактики", func() {
		It("Должно быть успешно сформирован список данных для формирования тегов галактик", func() {
			var isStopped bool

			addFunc := addListGalaxyTags(listGalaxyTags)
		DONE:
			for {
				select {
				case data := <-chanInput:
					addFunc(data.FieldBranch, data.Value)

				case <-stopped:
					isStopped = true
					chanDone <- struct{}{}

					break DONE
				}
			}

			lgt := listGalaxyTags.Get()
			galaxyTags := createGalaxyTags(&lgt)

			fmt.Println("********************************")
			fmt.Println(lgt)
			fmt.Println(galaxyTags)
			//"misp-galaxy:mitre-attack-pattern=\"Match Legitimate Name or Location - T1036.005\""
			fmt.Println("********************************")

			Expect(len(lgt)).Should(Equal(3))
			Expect(len(galaxyTags)).Should(Equal(3))
			Expect(isStopped).Should(BeTrue())
		})
	})
})
