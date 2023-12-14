Placeholder_MISP v0.9.8

Конфигурационные параметры для сервиса могут быть заданы как через конфигурационный файл так и методом установки переменных окружения.

Типы конфигурационных файлов:

- config.yaml общий конфигурационный файл
- config_dev.yaml конфигурационный файл используемый для тестов при разработке
- config_prod.yaml конфигурационный файл применяемый в продуктовом режиме

Основная переменная окружения для данного приложения - GO_PHMISP_MAIN. На основании
значения этой переменной принимается решение какой из конфигурационных файлов config_dev.yaml или config_prod.yaml использовать. При GO_PHMISP_MAIN=development
будет использоваться config_dev.yaml, во всех остальных случаях, в том числе и при отсутствии переменной окружения GO_PHMISP_MAIN будет использоваться конфигурационный файл config_prod.yaml. Перечень переменных окружения которые можно использовать для настройки приложения:

//Подключение к MISP
GO_PHMISP_MHOST
GO_PHMISP_MAUTH

//Подключение к NATS
GO_PHMISP_NHOST
GO_PHMISP_NPORT

//Подключение к СУБД Redis
GO_PHMISP_REDISHOST
GO_PHMISP_REDISPORT

//Подключение к СУБД Elasticsearch
GO_PHMISP_ESSEND
GO_PHMISP_ESHOST
GO_PHMISP_ESPORT
GO_PHMISP_ESPREFIX
GO_PHMISP_ESINDEX
GO_PHMISP_ESUSER
GO_PHMISP_ESPASSWD

//Подключение к
GO_PHMISP_NKCKIHOST
GO_PHMISP_NKCKIPORT
//Место расположения и наименования файла правил
GO_PHMISP_RULES_DIR
GO_PHMISP_RULES_FILE

Приоритет значений заданных через переменные окружения выше чем значений полученных из конфигурационных файлов. Таким образом можно осуществлять гибкую временную настройку приложения.

Сервис выполняет сделующие действия:

1.  Получает, через API MISP, список всех пользователей MISP с их авторизационным ключем. Это нужно для того что бы загружать json сообщения в форматах MISP от имени любого пользователя. Имя пользователя кейса TheHive берется из 'event.object.owner'. Если имени пользователя, полученного из TheHive, нет в MISP то такой пользователь автоматически создается.
2.  Осуществляет соединение с NATS.
3.  Получает кейсы от TheHive в формате json.
4.  Выполняет их разбор и анализ на основе правил. Есть два типа правил для анализа принятых кейсов, это правила разрешающие дальнейшую передачу кейса в MISP и правила, при совпадении параметров которых выполняется модификация данных в кейсе.
5.  Из кейсов, пропущенных для отправки MISP, формируются json сообщения на основе MISP форматов типа Events и Attributes. Данные сообщения загружаются через API MISP от имени пользователя создавшего кейс в TheHive (путь в json от TheHive 'event.object.owner')
    Тип Attributes формата MISP формируется по следующим условиям:
    - если, свойство observables.tags содержит значение вида misp:Attribution="whois-registrar", то осуществляется разбор данной строки, где
      значение Attribution добавляется в AttributesMispFormat.Category, а значение
      whois-registrar в AttributesMispFormat.Type,
    - если, свойство observables.dataType содержит одно из свойств определенного перечня значений, то свойства AttributesMispFormat.Category и AttributesMispFormat.Type будут заполненны на основании найденного значения в
      observables.dataType,
    - если, observables.tags содержит значение отличное от знаяения вида misp:Attribution="whois-registrar", но которое совпадает со значением в коде
      приложения, например, одно из подобных значений это type:<значение> которое может совпадать или быть схожем со содержимым в observables.dataType, то это
      значение добавляется в AttributesMispFormat.ObjectRelation
6.  После успешной отправки в MISP сформированных сообщений в TheHive отправляется json сообщение формата "{'success': True, 'service': 'MISP', 'commands': [{'command': 'addtag', 'string': 'Webhook: send="MISP"'}, {'command': 'setcustomfield', 'name': 'misp-event-id.string', 'string': '115199'}]}" содержащее идентификационный номер события полученного от MISP.
7.  Выполняеся передача кейсов TheHive в СУБД Elasticsearch.
