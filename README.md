Placeholder_MISP v0.4.1

Сервис выполняет сделующие действия:
 1. Получает, через API MISP, список всех пользователей MISP с их авторизационным ключем. Это нужно для того что бы загружать json сообщения в форматах MISP от имени любого пользователя.
 2. Осуществляет соединение с NATS.
 3. Получает кейсы от TheHive в формате json.
 4. Выполняет их разбор и анализ на основе правил. Есть два типа правил для анализа принятых кейсов, это правила разрешающие дальнейшую передачу кейса в MISP и правила, при совпадении параметров которых выполняется модификация данных в кейсе.
 5. Из кейсов, пропущенных для отправки MISP, формируются json сообщения на основе MISP форматов типа Events и Attributes. Данные сообщения загружаются через API MISP от имени пользователя создавшего кейс в TheHive (путь в json от TheHive 'event.object.owner')
 6. После успешной отправки в MISP сформированных сообщений в TheHive отправляется json сообщение формата "{'success': True, 'service': 'MISP', 'commands': [{'command': 'addtag', 'string': 'Webhook: send="MISP"'}, {'command': 'setcustomfield', 'name': 'misp-event-id.string', 'string': '115199'}]}" содержащее идентификационный номер события полученного от MISP.
 7. Выполняеся передача кейсов TheHive в СУБД Elasticsearch.


