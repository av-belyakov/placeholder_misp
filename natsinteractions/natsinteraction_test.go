package natsinteractions_test

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"placeholder_misp/confighandler"
	"placeholder_misp/coremodule"
	"placeholder_misp/datamodels"
	"placeholder_misp/natsinteractions"
	rules "placeholder_misp/rulesinteraction"
	"placeholder_misp/supportingfunctions"
	"placeholder_misp/tmpdata"
)

var _ = Describe("Natsinteraction", Ordered, func() {
	var (
		ctx      context.Context
		errConn  error
		closeCtx context.CancelFunc
		mnats    *natsinteractions.ModuleNATS
		chanLog  chan<- datamodels.MessageLoging
	)

	/*
		Для отправки логов в zabbix см. https://habr.com/ru/companies/nixys/news/503104/
	*/

	printVerificationWarning := func(lvw []string) string {
		var resultPrint string

		for _, v := range lvw {
			resultPrint += fmt.Sprintln(v)
		}

		return resultPrint
	}

	BeforeAll(func() {
		chanLog = make(chan<- datamodels.MessageLoging)

		ctx, closeCtx = context.WithTimeout(context.Background(), 2*time.Second)

		mnats, errConn = natsinteractions.NewClientNATS(ctx, confighandler.AppConfigNATS{
			//Host: "nats.cloud.gcm",
			Host: "127.0.0.1",
			Port: 4222,
		}, chanLog)
	})

	Context("Тест 1.1. Проверка декодирования тестовых данных из файла 'binaryDataOne'", func() {
		var exampleByte []byte

		for _, v := range strings.Split(tmpdata.GetExampleDataOne(), " ") {
			i, err := strconv.Atoi(v)
			if err != nil {
				continue
			}

			exampleByte = append(exampleByte, uint8(i))
		}

		/*It("Должно нормально отрабатывать функция  GetWhitespace", func() {
			fmt.Printf("%s none Whitespace\n", datamodels.GetWhitespace(0))
			fmt.Printf("%s one Whitespace\n", datamodels.GetWhitespace(1))
			fmt.Printf("%s two Whitespace\n", datamodels.GetWhitespace(2))
			fmt.Printf("%s three Whitespace\n", datamodels.GetWhitespace(3))

			Expect(true).Should(BeTrue())
		})*/

		It("При анмаршалинге данных в ИЗВЕСТНЫЙ ТИП ошибки быть не должно", func() {
			mm := datamodels.MainMessage{}
			err := json.Unmarshal(exampleByte, &mm)

			//fmt.Println("---- ExampleDataOne ----")
			//fmt.Println(mm.ToStringBeautiful(0))
			//fmt.Println("------------------------")

			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("Тест 1.2. Проверка декодирования тестовых данных из файла 'binaryDataTwo'", func() {
		var exampleByte []byte

		for _, v := range strings.Split(tmpdata.GetExampleDataTwo(), " ") {
			i, err := strconv.Atoi(v)
			if err != nil {
				continue
			}

			exampleByte = append(exampleByte, uint8(i))
		}

		It("При анмаршалинге данных в ИЗВЕСТНЫЙ ТИП ошибки быть не должно", func() {
			mm := datamodels.MainMessage{}
			err := json.Unmarshal(exampleByte, &mm)

			//fmt.Println("---- ExampleDataTwo ----")
			//fmt.Println(mm.ToStringBeautiful(0))
			//fmt.Println("------------------------")

			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("Тест 1.3. Проверка декодирования тестовых данных из файла 'binaryDataThree'", func() {
		var exampleByte []byte

		for _, v := range strings.Split(tmpdata.GetExampleDataThree(), " ") {
			i, err := strconv.Atoi(v)
			if err != nil {
				continue
			}

			exampleByte = append(exampleByte, uint8(i))
		}

		It("При анмаршалинге данных в ИЗВЕСТНЫЙ ТИП ошибки быть не должно", func() {
			mm := datamodels.MainMessage{}
			err := json.Unmarshal(exampleByte, &mm)

			Expect(err).ShouldNot(HaveOccurred())
		})

		It("При анмаршалинге в НЕИЗВЕСТНЫЙ тип ошибки быть не должно", func() {

			/*
				Теперь для вывода данных из JSON сообщения используется функция
			*/

			_, err := supportingfunctions.NewReadReflectJSONSprint(exampleByte)

			/*
				!!!!!!
				В coremodule.decodeMessageReflect тестово реализованна возможность замены
				строковых данных в определенных полях, подробнее смотреть в функции readReflectAnyTypeSprint
				надо написать для других типов, таких как int, bool. И сделать замену на основе правил из
				пакета rules
				!!!!!!

				result := map[string]interface{}{}
				err := json.Unmarshal(exampleByte, &result)
				Expect(err).ShouldNot(HaveOccurred())

				strData := coremodule.ReadReflectMapSprint(result, rules.ListRulesProcessedMISPMessage{}, 0)
				Expect(err).ShouldNot(HaveOccurred())
			*/

			//fmt.Println("---- REFLECTION MAPPING ExampleDataThree ----")
			//fmt.Println(strData)
			//fmt.Println("---------------------------------------------")

			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("Тест 2. Проверка замены некоторых значений в json файле", func() {
		var eb []byte

		for _, v := range strings.Split(tmpdata.GetExampleDataThree(), " ") {
			i, err := strconv.Atoi(v)
			if err != nil {
				continue
			}

			eb = append(eb, uint8(i))
		}

		It("Должны быть заменины некоторые значения на основе правил из файла 'procmispmsg.yaml'", func() {
			//lrp, lvw, err := rules.GetRuleProcessedMISPMsg("rules", "procmispmsg.yaml")
			var (
				strData     string
				procMsgHive coremodule.ProcessMessageFromHive
			)

			lr, lw, err := rules.GetRuleProcessingMsgForMISP("rules", "procmispmsg_test.yaml")

			fmt.Println("list verification warning:")
			fmt.Println(printVerificationWarning(lw))

			Expect(err).ShouldNot(HaveOccurred())

			reg := regexp.MustCompile(`_createdBy: \'[a-zA-Z_.@0-9]+\'`)
			regCaseId := regexp.MustCompile(`caseId: [0-9]+`)
			regRevoked := regexp.MustCompile(`revoked: [a-zA-Z]+`)
			regStartDate := regexp.MustCompile(`startDate: [a-zA-Z_.@0-9]+`)
			regPatternId := regexp.MustCompile(`capecId: \'[a-zA-Z0-9\-]+\'`)

			fmt.Println("____BEFORE reflect modify____")
			strData, err = supportingfunctions.NewReadReflectJSONSprint(eb)
			bl := reg.FindAllString(strData, 10)
			for k, v := range bl {
				fmt.Printf("%d. %s\n", k+1, v)
			}
			cibl := regCaseId.FindAllString(strData, 10)
			for k, v := range cibl {
				fmt.Printf("%d. %s\n", k+1, v)
			}
			rbl := regRevoked.FindAllString(strData, 10)
			for k, v := range rbl {
				fmt.Printf("%d. %s\n", k+1, v)
			}
			pidbl := regPatternId.FindAllString(strData, 10)
			for k, v := range pidbl {
				fmt.Printf("%d. %s\n", k+1, v)
			}

			Expect(err).ShouldNot(HaveOccurred())

			procMsgHive, err = coremodule.NewHandleMessageFromHive(eb, lr)

			Expect(err).ShouldNot(HaveOccurred())

			ok, warningMsg := procMsgHive.HandleMessage()
			neb, err := procMsgHive.GetMessage()

			Expect(err).ShouldNot(HaveOccurred())

			fmt.Println("____AFTER reflect modify____:")
			strData, err = supportingfunctions.NewReadReflectJSONSprint(neb)
			bl = reg.FindAllString(strData, 10)
			for k, v := range bl {
				fmt.Printf("%d. %s\n", k+1, v)
			}
			cibl = regCaseId.FindAllString(strData, 10)
			for k, v := range cibl {
				fmt.Printf("%d. %s\n", k+1, v)
			}
			rbl = regRevoked.FindAllString(strData, 10)
			for k, v := range rbl {
				fmt.Printf("%d. %s\n", k+1, v)
			}
			pidbl = regPatternId.FindAllString(strData, 10)
			for k, v := range pidbl {
				fmt.Printf("%d. %s\n", k+1, v)
			}
			sdbl := regStartDate.FindAllString(strData, 10)
			for k, v := range sdbl {
				fmt.Printf("%d. %s\n", k+1, v)
			}

			fmt.Println("procMsgHive.ProcessMessage() is true: ", ok)
			fmt.Println("warningMsg: ", warningMsg)
			fmt.Println("")

			Expect(err).ShouldNot(HaveOccurred())
		})

		It("Должен быть получен полный текстовый вывод объекта из хайва", func() {
			str, err := supportingfunctions.NewReadReflectJSONSprint(eb)

			//fmt.Println(str)

			Expect(str).ShouldNot(Equal(""))
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("Тест 3. Проверка инициализации соединения с NATS", func() {
		It("При инициализации соединения с NATS не должно быть ошибки", func() {
			Expect(errConn).ShouldNot(HaveOccurred())

			fmt.Println("Resevid message = ", <-mnats.GetDataReceptionChannel())
		})
	})

	AfterAll(func() {
		closeCtx()
	})
})

/*
{
"source": "gcm-test",
"event": {
    "operation": "update",
    "details": {
        "endDate": 1690468334490,
        "customFields": {
            "class-attack": {
                "order": 1,
                "string": "Exploit"
            },
            "first-time": {
                "date": 1690468300490,
                "order": 0
            },
            "last-time": {
                "date": 1690468320490,
                "order": 0
            },
            "misp-event-id": {
                "order": 2,
                "string": "7481"
            },
            "ncircc-bulletine-id": {
                "order": 0,
                "string": "21-09-45"
            },
            "ncircc-class-attack": {
                "order": 1,
                "string": "Попытки эксплуатации уязвимости;attack"
            }
        },
        "resolutionStatus": "TruePositive",
        "summary": "Попытка эксплуатации уязвимости CVE-2019-11043 ",
        "status": "Resolved",
        "impactStatus": "NoImpact"
    },
    "objectType": "case",
    "objectId": "~1938194504",
    "base": true,
    "startDate": 1636622316515,
    "rootId": "~1938194504",
    "requestId": "821288ed2fe50017: -5dad3d10: 17d0da95f31: -8000: 12615",
    "object": {
        "_id": "~1938194504",
        "id": "~1938194504",
        "createdBy": "i.monahov@cloud.gcm",
        "updatedBy": "m.miroshnichenko@cloud.gcm",
        "createdAt": 1630597296241,
        "updatedAt": 1636622316508,
        "_type": "case",
        "caseId": 99990,
        "title": "Попытка эксплуатации уязвимости CVE-2019-11043 (Распределенная атака)",
        "description": "Попытка эксплуатации уязвимости CVE-2019-11043 (Распределенная атака)\n  \n#### Merged with alert #TSK-CENTER-2-ZPM-210902-136913 Зафиксирована подозрительная активность по признакам НКЦКИ с 68.183.207.206\n\n**Задача переданная из смежной системы: Заслон-Пост-Модерн**\n\nВ формате ГЦМ: **`TSK-CENTER-2-ZPM-210902-136913`** ID: `136913`\n\n[http: //siem.cloud.gcm:3000/tasks/card/TSK-CENTER-2-ZPM-210902-136913](http://siem.cloud.gcm:3000/tasks/card/TSK-CENTER-2-ZPM-210902-136913)\n\nАвтор задачи: **`admin`**\n\nТип: **`snort_alert`**\n\n**Причина по которой создана задача**\n\nНазвание: `Зафиксирована подозрительная активность по признакам НКЦКИ с 68.183.207.206`\n\nОписание: `## Данная задача создана автоматически\n Время начала: 2021-09-02 06:08:30\n Время окончания: 2021-09-02 06:08:40\n Продолжительность воздействий: 0:00:10`\n\nОтработало на СОА: \n- **`1052`**   ФАС, Установлен: г. Москва, IP адрес: 10.20.0.52\n\n\n\n**Полное описаие события IDS:**\n\n- Время начала: **`02.09.2021 06:08:30`**\n- Время окончания: **`02.09.2021 06:08:40`**\n- **IP из домашней подсети**\n\n1. **`194.226.26.252`**\n- **IP из внешней подсети**\n\n1. **`68.183.207.206`**\n\n\n**Сигнатуры на которых отработал анализатор сетевого трафика:**\n\n1. РП: **`3005391`**, Сообщение: AM Exploit Apache HTTP Server mod-status Race Condition Heap Buffer Overflow, Добавлена: 10.05.2021 10:32:33\n2. РП: **`3005422`**, Сообщение: AM Exploit HTTP Apache mod_mylo Module Possible Buffer Overflow, Добавлена: 10.05.2021 10:32:33\n3. РП: **`3066221`**, Сообщение: AM Exploit Possible RCE in PHP-fpm, Добавлена: 10.05.2021 10:32:33\n4. РП: **`99010424`**, Сообщение: NCIRCC PHP-RCE check (D-Pisos: 8=D), Добавлена: 10.05.2021 10:27:52\n\n\n\n\n\n**Фильтрация и выгрузка от** Thu Sep 02 2021 18:31:58 GMT+0300 \n\nРазмер: **`171.4 KB`**, [Скачать файл](ftp://ftp.cloud.gcm//traffic/1052/1630596715_2021_09_02____18_31_55_492399.pcap)\n\nКонтент: Фильтрация успешно завершена. Включена автоматическая выгрузка\n\nФайл на СОА: `/opt/zaslon/zmanager/data/pfilter_storage/1630596715_2021_09_02____18_31_55_492399.pcap`\n  \n#### Merged with alert #TSK-CENTER-2-ZPM-210902-136777 Зафиксирована подозрительная активность по признакам НКЦКИ с 68.183.207.206\n\n**Задача переданная из смежной системы: Заслон-Пост-Модерн**\n\nВ формате ГЦМ: **`TSK-CENTER-2-ZPM-210902-136777`** ID: `136777`\n\n[http://siem.cloud.gcm:3000/tasks/card/TSK-CENTER-2-ZPM-210902-136777](http://siem.cloud.gcm:3000/tasks/card/TSK-CENTER-2-ZPM-210902-136777)\n\nАвтор задачи: **`admin`**\n\nТип: **`snort_alert`**\n\n**Причина по которой создана задача**\n\nНазвание: `Зафиксирована подозрительная активность по признакам НКЦКИ с 68.183.207.206`\n\nОписание: `## Данная задача создана автоматически\n Время начала: 2021-09-02 05:37:38\n Время окончания: 2021-09-02 05:38:44\n Продолжительность воздействий: 0:01:06`\n\nОтработало на СОА: \n- **`1052`**   ФАС, Установлен: г. Москва, IP адрес: 10.20.0.52\n\n\n\n**Полное описаие события IDS:**\n\n- Время начала: **`02.09.2021 05:37:38`**\n- Время окончания: **`02.09.2021 05:38:44`**\n- **IP из домашней подсети**\n\n1. **`194.226.26.127`**\n- **IP из внешней подсети**\n\n1. **`68.183.207.206`**\n\n\n**Сигнатуры на которых отработал анализатор сетевого трафика:**\n\n1. РП: **`3005422`**, Сообщение: AM Exploit HTTP Apache mod_mylo Module Possible Buffer Overflow, Добавлена: 10.05.2021 10:32:33\n2. РП: **`3066221`**, Сообщение: AM Exploit Possible RCE in PHP-fpm, Добавлена: 10.05.2021 10:32:33\n3. РП: **`99010424`**, Сообщение: NCIRCC PHP-RCE check (D-Pisos: 8=D), Добавлена: 10.05.2021 10:27:52\n4. РП: **`99010425`**, Сообщение: NCIRCC PHP-RCE Target is probably vulnerable (Status code 502, adding as a candidate), Добавлена: 10.05.2021 10:27:52\n\n\n\n\n\n**Фильтрация и выгрузка от** Thu Sep 02 2021 16:59:35 GMT+0300 \n\nРазмер: **`231.6 KB`**, [Скачать файл](ftp://ftp.cloud.gcm//traffic/1052/1630591171_2021_09_02____16_59_31_905897.pcap)\n\nКонтент: Фильтрация успешно завершена. Включена автоматическая выгрузка\n\nФайл на СОА: `/opt/zaslon/zmanager/data/pfilter_storage/1630591171_2021_09_02____16_59_31_905897.pcap`\n  \n#### Merged with alert #TSK-CENTER-2-ZPM-210902-136758 Зафиксирована подозрительная активность по признакам НКЦКИ с 68.183.207.206\n\n**Задача переданная из смежной системы: Заслон-Пост-Модерн**\n\nВ формате ГЦМ: **`TSK-CENTER-2-ZPM-210902-136758`** ID: `136758`\n\n[http://siem.cloud.gcm:3000/tasks/card/TSK-CENTER-2-ZPM-210902-136758](http://siem.cloud.gcm:3000/tasks/card/TSK-CENTER-2-ZPM-210902-136758)\n\nАвтор задачи: **`admin`**\n\nТип: **`snort_alert`**\n\n**Причина по которой создана задача**\n\nНазвание: `Зафиксирована подозрительная активность по признакам НКЦКИ с 68.183.207.206`\n\nОписание: `## Данная задача создана автоматически\n Время начала: 2021-09-02 06:30:39\n Время окончания: 2021-09-02 06:30:50\n Продолжительность воздействий: 0:00:11`\n\nОтработало на СОА: \n- **`1052`**   ФАС, Установлен: г. Москва, IP адрес: 10.20.0.52\n\n\n\n**Полное описаие события IDS:**\n\n- Время начала: **`02.09.2021 06:30:39`**\n- Время окончания: **`02.09.2021 06:30:50`**\n- **IP из домашней подсети**\n\n1. **`194.226.26.250`**\n- **IP из внешней подсети**\n\n1. **`68.183.207.206`**\n\n\n**Сигнатуры на которых отработал анализатор сетевого трафика:**\n\n1. РП: **`3005391`**, Сообщение: AM Exploit Apache HTTP Server mod-status Race Condition Heap Buffer Overflow, Добавлена: 10.05.2021 10:32:33\n2. РП: **`3005422`**, Сообщение: AM Exploit HTTP Apache mod_mylo Module Possible Buffer Overflow, Добавлена: 10.05.2021 10:32:33\n3. РП: **`3066221`**, Сообщение: AM Exploit Possible RCE in PHP-fpm, Добавлена: 10.05.2021 10:32:33\n4. РП: **`99010424`**, Сообщение: NCIRCC PHP-RCE check (D-Pisos: 8=D), Добавлена: 10.05.2021 10:27:52\n\n\n\n\n\n**Фильтрация и выгрузка от** Thu Sep 02 2021 15:51:37 GMT+0300 \n\nРазмер: **`192.7 KB`**, [Скачать файл](ftp://ftp.cloud.gcm//traffic/1052/1630587093_2021_09_02____15_51_33_911281.pcap)\n\nКонтент: Фильтрация успешно завершена. Включена автоматическая выгрузка\n\nФайл на СОА: `/opt/zaslon/zmanager/data/pfilter_storage/1630587093_2021_09_02____15_51_33_911281.pcap`\n  \n#### Merged with alert #TSK-CENTER-2-ZPM-210902-136735 Зафиксирована подозрительная активность по признакам НКЦКИ с 68.183.207.206\n\n**Задача переданная из смежной системы: Заслон-Пост-Модерн**\n\nВ формате ГЦМ: **`TSK-CENTER-2-ZPM-210902-136735`** ID: `136735`\n\n[http://siem.cloud.gcm:3000/tasks/card/TSK-CENTER-2-ZPM-210902-136735](http://siem.cloud.gcm:3000/tasks/card/TSK-CENTER-2-ZPM-210902-136735)\n\nАвтор задачи: **`admin`**\n\nТип: **`snort_alert`**\n\n**Причина по которой создана задача**\n\nНазвание: `Зафиксирована подозрительная активность по признакам НКЦКИ с 68.183.207.206`\n\nОписание: `## Данная задача создана автоматически\n Время начала: 2021-09-02 03:46:48\n Время окончания: 2021-09-02 05:17:57\n Продолжительность воздействий: 1:31:09`\n\nОтработало на СОА: \n- **`1052`**   ФАС, Установлен: г. Москва, IP адрес: 10.20.0.52\n\n\n\n**Полное описаие события IDS:**\n\n- Время начала: **`02.09.2021 03:46:48`**\n- Время окончания: **`02.09.2021 05:17:57`**\n- **IP из домашней подсети**\n\n1. **`194.226.26.115`**\n2. **`194.226.26.68`**\n3. **`194.226.26.185`**\n4. **`194.226.26.36`**\n5. **`194.226.26.120`**\n- **IP из внешней подсети**\n\n1. **`68.183.207.206`**\n\n\n**Сигнатуры на которых отработал анализатор сетевого трафика:**\n\n1. РП: **`3005391`**, Сообщение: AM Exploit Apache HTTP Server mod-status Race Condition Heap Buffer Overflow, Добавлена: 10.05.2021 10:32:33\n2. РП: **`3005422`**, Сообщение: AM Exploit HTTP Apache mod_mylo Module Possible Buffer Overflow, Добавлена: 10.05.2021 10:32:33\n3. РП: **`3066221`**, Сообщение: AM Exploit Possible RCE in PHP-fpm, Добавлена: 10.05.2021 10:32:33\n4. РП: **`99010424`**, Сообщение: NCIRCC PHP-RCE check (D-Pisos: 8=D), Добавлена: 10.05.2021 10:27:52\n\n\n\n\n\n**Фильтрация и выгрузка от** Thu Sep 02 2021 15:00:31 GMT+0300 \n\nРазмер: **`1.3 MB`**, [Скачать файл](ftp://ftp.cloud.gcm//traffic/1052/1630584012_2021_09_02____15_00_12_719724.pcap)\n\nКонтент: Фильтрация успешно завершена. Включена автоматическая выгрузка\n\nФайл на СОА: `/opt/zaslon/zmanager/data/pfilter_storage/1630584012_2021_09_02____15_00_12_719724.pcap`\n  \n#### Merged with alert #TSK-CENTER-2-ZPM-210902-136682 Редко встречающиеся признаки ВПО, Зафиксирована подозрительная активность по признакам НКЦКИ с 68.183.207.206\n\n**Задача переданная из смежной системы: Заслон-Пост-Модерн**\n\nВ формате ГЦМ: **`TSK-CENTER-2-ZPM-210902-136682`** ID: `136682`\n\n[http://siem.cloud.gcm:3000/tasks/card/TSK-CENTER-2-ZPM-210902-136682](http://siem.cloud.gcm:3000/tasks/card/TSK-CENTER-2-ZPM-210902-136682)\n\nАвтор задачи: **`admin`**\n\nТип: **`snort_alert`**\n\n**Причина по которой создана задача**\n\nНазвание: `Редко встречающиеся признаки ВПО, Зафиксирована подозрительная активность по признакам НКЦКИ с 68.183.207.206`\n\nОписание: `## Данная задача создана автоматически\n Время начала: 2021-09-02 03:46:48\n Время окончания: 2021-09-02 03:46:59\n Продолжительность воздействий: 0:00:11`\n\nОтработало на СОА: \n- **`1052`**   ФАС, Установлен: г. Москва, IP адрес: 10.20.0.52\n\n\n\n**Полное описаие события IDS:**\n\n- Время начала: **`02.09.2021 03:46:48`**\n- Время окончания: **`02.09.2021 03:46:59`**\n- **IP из домашней подсети**\n\n1. **`194.226.26.185`**\n- **IP из внешней подсети**\n\n1. **`68.183.207.206`**\n\n\n**Сигнатуры на которых отработал анализатор сетевого трафика:**\n\n1. РП: **`99010424`**, Сообщение: NCIRCC PHP-RCE check (D-Pisos: 8=D), Добавлена: 10.05.2021 10:27:52\n\n\n\n\n\n**Фильтрация и выгрузка от** Thu Sep 02 2021 11:31:42 GMT+0300 \n\nРазмер: **`404.6 KB`**, [Скачать файл](ftp://ftp.cloud.gcm//traffic/1052/1630571500_2021_09_02____11_31_40_570968.pcap)\n\nКонтент: Фильтрация успешно завершена. Включена автоматическая выгрузка\n\nФайл на СОА: `/opt/zaslon/zmanager/data/pfilter_storage/1630571500_2021_09_02____11_31_40_570968.pcap`",
        "severity": 2,
        "startDate": 1630597260000,
        "endDate": 1636622316490,
        "impactStatus": "NoImpact",
        "resolutionStatus": "TruePositive",
        "tags": [
            "ATs:geoip=Канада",
            "ATs:reason=Редко встречающиеся признаки ВПО",
            "Webhook:MISP",
            "Sensor_id=1052",
            "WebHook:send=ATD",
            "ZK:send",
            "WebHook:send=NCIRCC",
            "Webhook:ES",
            "WebHook:send=PCAP",
            "ATs:reason=Зафиксирована подозрительная активность по признакам НКЦКИ"
        ],
        "flag": false,
        "tlp": 2,
        "pap": 2,
        "status": "Resolved",
        "summary": "Попытка эксплуатации уязвимости CVE-2019-11043 ",
        "owner": "i.monahov@cloud.gcm",
        "customFields": {
            "class-attack": {
                "order": 1,
                "string": "Exploit"
            },
            "first-time": {
                "date": 1630543560000,
                "order": 0
            },
            "last-time": {
                "date": 1630543560000,
                "order": 0
            },
            "misp-event-id": {
                "order": 2,
                "string": "7481"
            },
            "ncircc-bulletine-id": {
                "order": 0,
                "string": "21-09-45"
            },
            "ncircc-class-attack": {
                "order": 1,
                "string": "Попытки эксплуатации уязвимости;attack"
            }
        },
        "stats": {},
        "permissions": []
    },
    "organisationId": "~4192",
    "organisation": "GCM-TEST"
},
"observables": [],
"ttp": []
}

ttps:
  ttp:
    1.
      _createdAt: '1686652917436'
      _createdBy: 'p.delyukin@cloud.rcm'
      _id: '~248369200'
        _createdAt: '1679910385644'
        _createdBy: 'admin@thehive.local'
        _id: '~282632'
        _type: 'Pattern'
        dataSources:
          1. 'User Account: User Account Authentication'
          2. 'Command: Command Execution'
          3. 'Application Log: Application Log Content'
        defenseBypassed:
        description: 'Adversaries may use brute force techniques to gain access to accounts when passwords are unknown or when password hashes are obtained. Without knowledge of the password for an account or set of accounts, an adversary may systematically guess the password using a repetitive or iterative mechanism. Brute forcing passwords can take place via interaction with a service that will check the validity of those credentials or offline against previously acquired credential data, such as password hashes.

Brute forcing credentials may take place at various points during a breach. For example, adversaries may attempt to brute force access to [Valid Accounts](https://attack.mitre.org/techniques/T1078) within a victim environment leveraging knowledge gathered from other post-compromise behaviors such as [OS Credential Dumping](https://attack.mitre.org/techniques/T1003), [Account Discovery](https://attack.mitre.org/techniques/T1087), or [Password Policy Discovery](https://attack.mitre.org/techniques/T1201). Adversaries may also combine brute forcing activity with behaviors such as [External Remote Services](https://attack.mitre.org/techniques/T1133) as part of Initial Access.'
        extraData:
        name: 'Brute Force'
        patternId: 'T1110'
        patternType: 'attack-pattern'
        permissionsRequired:
        platforms:
          1. 'Windows'
          2. 'Azure AD'
          3. 'Office 365'
          4. 'SaaS'
          5. 'IaaS'
          6. 'Linux'
          7. 'macOS'
          8. 'Google Workspace'
          9. 'Containers'
          10. 'Network'
        remoteSupport: 'false'
        revoked: 'false'
        systemRequirements:
        tactics:
          1. 'credential-access'
        URL: 'https://attack.mitre.org/techniques/T1110'
        version: '2.4'
        _createdAt: '0'
        _createdBy: ''
        _id: ''
        _type: ''
        dataSources:
        defenseBypassed:
        description: ''
        extraData:
        name: ''
        patternId: ''
        patternType: ''
        permissionsRequired:
        platforms:
        remoteSupport: 'false'
        revoked: 'false'
        systemRequirements:
        tactics:
        URL: ''
        version: ''
      occurDate: '1686652860000'
      patternId: 'T1110'
      tactic: 'credential-access'

	REFLECTION MAPPING
	ttp:
    occurDate: 1686652860000
    patternId: 'T1110'
    tactic: 'credential-access'
    _createdAt: 1686652917436
    _createdBy: 'p.delyukin@cloud.rcm'
    _id: '~248369200'
    extraData:
      pattern:
        platforms:
          1. 'Windows'
          2. 'Azure AD'
          3. 'Office 365'
          4. 'SaaS'
          5. 'IaaS'
          6. 'Linux'
          7. 'macOS'
          8. 'Google Workspace'
          9. 'Containers'
          10. 'Network'
        remoteSupport: false
        url: 'https://attack.mitre.org/techniques/T1110'
        version: '2.4'
        capecId: 'CAPEC-49'
        _createdBy: 'admin@thehive.local'
        _id: '~282632'
        dataSources:
          1. 'User Account: User Account Authentication'
          2. 'Command: Command Execution'
          3. 'Application Log: Application Log Content'
        description: 'Adversaries may use brute force techniques to gain access to accounts when passwords are unknown or when password hashes are obtained. Without knowledge of the password for an account or set of accounts, an adversary may systematically guess the password using a repetitive or iterative mechanism. Brute forcing passwords can take place via interaction with a service that will check the validity of those credentials or offline against previously acquired credential data, such as password hashes.

Brute forcing credentials may take place at various points during a breach. For example, adversaries may attempt to brute force access to [Valid Accounts](https://attack.mitre.org/techniques/T1078) within a victim environment leveraging knowledge gathered from other post-compromise behaviors such as [OS Credential Dumping](https://attack.mitre.org/techniques/T1003), [Account Discovery](https://attack.mitre.org/techniques/T1087), or [Password Policy Discovery](https://attack.mitre.org/techniques/T1201). Adversaries may also combine brute forcing activity with behaviors such as [External Remote Services](https://attack.mitre.org/techniques/T1133) as part of Initial Access.'
        extraData:
        name: 'Brute Force'
        patternType: 'attack-pattern'
        _createdAt: 1679910385644
        revoked: false
        capecUrl: 'https://capec.mitre.org/data/definitions/49.html'
        defenseBypassed:
        detection: 'Monitor authentication logs for system and application login failures of [Valid Accounts](https://attack.mitre.org/techniques/T1078). If authentication failures are high, then there may be a brute force attempt to gain access to a system using legitimate credentials. Also monitor for many failed authentication attempts across various accounts that may result from password spraying attempts. It is difficult to detect when hashes are cracked, since this is generally done outside the scope of the target network.'
        patternId: 'T1110'
        systemRequirements:
        _type: 'Pattern'
        tactics:
          1. 'credential-access'
        permissionsRequired:
*/

/*
	Извесный тип
ttps:
  ttp:
    1.
      _createdAt: '1686652917436'
      _createdBy: 'p.delyukin@cloud.rcm'
      _id: '~248369200'
        _createdAt: '1679910385644'
        _createdBy: 'admin@thehive.local'
        _id: '~282632'
        _type: 'Pattern'
        dataSources:
          1. 'User Account: User Account Authentication'
          2. 'Command: Command Execution'
          3. 'Application Log: Application Log Content'
        defenseBypassed:
        description: 'Adversaries may use brute force techniques to gain access to accounts when passwords are unknown or when password hashes are obtained. Without knowledge of the password for an account or set of accounts, an adversary may systematically guess the password using a repetitive or iterative mechanism. Brute forcing passwords can take place via interaction with a service that will check the validity of those credentials or offline against previously acquired credential data, such as password hashes.

Brute forcing credentials may take place at various points during a breach. For example, adversaries may attempt to brute force access to [Valid Accounts](https://attack.mitre.org/techniques/T1078) within a victim environment leveraging knowledge gathered from other post-compromise behaviors such as [OS Credential Dumping](https://attack.mitre.org/techniques/T1003), [Account Discovery](https://attack.mitre.org/techniques/T1087), or [Password Policy Discovery](https://attack.mitre.org/techniques/T1201). Adversaries may also combine brute forcing activity with behaviors such as [External Remote Services](https://attack.mitre.org/techniques/T1133) as part of Initial Access.'
        extraData:
        name: 'Brute Force'
        patternId: 'T1110'
        patternType: 'attack-pattern'
        permissionsRequired:
        platforms:
          1. 'Windows'
          2. 'Azure AD'
          3. 'Office 365'
          4. 'SaaS'
          5. 'IaaS'
          6. 'Linux'
          7. 'macOS'
          8. 'Google Workspace'
          9. 'Containers'
          10. 'Network'
        remoteSupport: 'false'
        revoked: 'false'
        systemRequirements:
        tactics:
          1. 'credential-access'
        URL: 'https://attack.mitre.org/techniques/T1110'
        version: '2.4'
        _createdAt: '0'
        _createdBy: ''
        _id: ''
        _type: ''
        dataSources:
        defenseBypassed:
        description: ''
        extraData:
        name: ''
        patternId: ''
        patternType: ''
        permissionsRequired:
        platforms:
        remoteSupport: 'false'
        revoked: 'false'
        systemRequirements:
        tactics:
        URL: ''
        version: ''
      occurDate: '1686652860000'
      patternId: 'T1110'
      tactic: 'credential-access'

	REFLECTION MAPPING
	ttp:
    occurDate: 1686652860000
    patternId: 'T1110'
    tactic: 'credential-access'
    _createdAt: 1686652917436
    _createdBy: 'p.delyukin@cloud.rcm'
    _id: '~248369200'
    extraData:
      pattern:
        platforms:
          1. 'Windows'
          2. 'Azure AD'
          3. 'Office 365'
          4. 'SaaS'
          5. 'IaaS'
          6. 'Linux'
          7. 'macOS'
          8. 'Google Workspace'
          9. 'Containers'
          10. 'Network'
        remoteSupport: false
        url: 'https://attack.mitre.org/techniques/T1110'
        version: '2.4'
        capecId: 'CAPEC-49'
        _createdBy: 'admin@thehive.local'
        _id: '~282632'
        dataSources:
          1. 'User Account: User Account Authentication'
          2. 'Command: Command Execution'
          3. 'Application Log: Application Log Content'
        description: 'Adversaries may use brute force techniques to gain access to accounts when passwords are unknown or when password hashes are obtained. Without knowledge of the password for an account or set of accounts, an adversary may systematically guess the password using a repetitive or iterative mechanism. Brute forcing passwords can take place via interaction with a service that will check the validity of those credentials or offline against previously acquired credential data, such as password hashes.

Brute forcing credentials may take place at various points during a breach. For example, adversaries may attempt to brute force access to [Valid Accounts](https://attack.mitre.org/techniques/T1078) within a victim environment leveraging knowledge gathered from other post-compromise behaviors such as [OS Credential Dumping](https://attack.mitre.org/techniques/T1003), [Account Discovery](https://attack.mitre.org/techniques/T1087), or [Password Policy Discovery](https://attack.mitre.org/techniques/T1201). Adversaries may also combine brute forcing activity with behaviors such as [External Remote Services](https://attack.mitre.org/techniques/T1133) as part of Initial Access.'
        extraData:
        name: 'Brute Force'
        patternType: 'attack-pattern'
        _createdAt: 1679910385644
        revoked: false
        capecUrl: 'https://capec.mitre.org/data/definitions/49.html'
        defenseBypassed:
        detection: 'Monitor authentication logs for system and application login failures of [Valid Accounts](https://attack.mitre.org/techniques/T1078). If authentication failures are high, then there may be a brute force attempt to gain access to a system using legitimate credentials. Also monitor for many failed authentication attempts across various accounts that may result from password spraying attempts. It is difficult to detect when hashes are cracked, since this is generally done outside the scope of the target network.'
        patternId: 'T1110'
        systemRequirements:
        _type: 'Pattern'
        tactics:
          1. 'credential-access'
        permissionsRequired:
*/
