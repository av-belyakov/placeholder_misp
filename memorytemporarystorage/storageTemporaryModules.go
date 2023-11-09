package memorytemporarystorage

import (
	"sync"
	"time"
)

var once sync.Once
var cst CommonStorageTemporary

func NewTemporaryStorage() *CommonStorageTemporary {
	once.Do(func() {
		cst = CommonStorageTemporary{
			temporaryInputCase: TemporaryInputCases{
				Cases: make(map[int]SettingsInputCase),
			},
			HiveFormatMessage: HiveFormatMessages{
				Storages: make(map[string]StorageHiveFormatMessages),
			},
			ListUserSettingsMISP: make([]UserSettingsMISP, 0),
			dataCounter:          DataCounterStorage{},
		}

		go checkTimeDelete(&cst)
	})

	return &cst
}

// checkTimeDeleteTemporaryStorageSearchQueries очистка информации о задаче по истечении определенного времени или неактуальности данных
func checkTimeDelete(cst *CommonStorageTemporary) {
	c := time.Tick(3 * time.Second)

	for range c {
		go func() {
			for k, v := range cst.HiveFormatMessage.Storages {
				if v.isProcessedMisp && v.isProcessedElasticsearsh && v.isProcessedNKCKI {
					deleteHiveFormatMessageElement(k, cst)
				}
			}
		}()

		go func() {
			for k, v := range cst.temporaryInputCase.Cases {
				if time.Now().Unix() > (v.TimeCreate + 54000) {
					deleteTemporaryCase(k, cst)
				}
			}
		}()
	}
}

// GetDataCounter возвращает информацию по сетчикам
func (cst *CommonStorageTemporary) GetDataCounter() DataCounter {
	return DataCounter{
		AcceptedEvents:       cst.dataCounter.acceptedEvents,
		ProcessedEvents:      cst.dataCounter.processedEvents,
		EventsMeetRules:      cst.dataCounter.eventsMeetRules,
		EventsDoNotMeetRules: cst.dataCounter.eventsDoNotMeetRules,
		StartTime:            cst.dataCounter.startTime,
	}
}

// SetStartTimeDataCounter добавляет время начала сетчика
func (cst *CommonStorageTemporary) SetStartTimeDataCounter(t time.Time) {
	cst.dataCounter.Lock()
	defer cst.dataCounter.Unlock()

	cst.dataCounter.startTime = t
}

// SetAcceptedEventsDataCounter увеличивает сетчик принятых событий
func (cst *CommonStorageTemporary) SetAcceptedEventsDataCounter(num int) {
	cst.dataCounter.Lock()
	defer cst.dataCounter.Unlock()

	cst.dataCounter.acceptedEvents += num
}

// SetProcessedEventsDataCounter увеличивает сетчик обработанных событий
func (cst *CommonStorageTemporary) SetProcessedEventsDataCounter(num int) {
	cst.dataCounter.Lock()
	defer cst.dataCounter.Unlock()

	cst.dataCounter.processedEvents += num
}

// SetEventsMeetRulesDataCounter увеличивает сетчик событий соответствующих правилу
func (cst *CommonStorageTemporary) SetEventsMeetRulesDataCounter(num int) {
	cst.dataCounter.Lock()
	defer cst.dataCounter.Unlock()

	cst.dataCounter.eventsMeetRules += num
}

// SetEventsDoNotMeetRulesDataCounter увеличивает сетчик событий не соответствующих правилу
func (cst *CommonStorageTemporary) SetEventsDoNotMeetRulesDataCounter(num int) {
	cst.dataCounter.Lock()
	defer cst.dataCounter.Unlock()

	cst.dataCounter.eventsDoNotMeetRules += num
}

// GetTemporaryCase возвращает информацию из временного списка входящих кейсов
func (cst *CommonStorageTemporary) GetTemporaryCase(id int) (SettingsInputCase, bool) {
	s, ok := cst.temporaryInputCase.Cases[id]

	return s, ok
}

// SetTemporaryCase добавляет информацию о кейсах во временное хранилище
func (cst *CommonStorageTemporary) SetTemporaryCase(id int, s SettingsInputCase) {
	cst.temporaryInputCase.Lock()
	defer cst.temporaryInputCase.Unlock()

	s.TimeCreate = time.Now().Unix()
	cst.temporaryInputCase.Cases[id] = s
}

// GetTemporaryCases возвращает список временных кейсов
func (cst *CommonStorageTemporary) GetListTemporaryCases() map[int]SettingsInputCase {
	return cst.temporaryInputCase.Cases
}

// GetCountHiveFormatMessage возвращает количество сообщений полученных от TheHive и еще не обработанных
func (cst *CommonStorageTemporary) GetCountHiveFormatMessage() int {
	return len(cst.HiveFormatMessage.Storages)
}

// SetOriginalHaveFormatMessage добавляет сырые данные полученные от TheHive
func (cst *CommonStorageTemporary) SetRawDataHiveFormatMessage(uuid string, data []byte) {
	cst.HiveFormatMessage.Lock()
	defer cst.HiveFormatMessage.Unlock()

	if _, ok := cst.HiveFormatMessage.Storages[uuid]; !ok {
		cst.HiveFormatMessage.Storages[uuid] = StorageHiveFormatMessages{}
	}

	storage := cst.HiveFormatMessage.Storages[uuid]
	cst.HiveFormatMessage.Storages[uuid] = StorageHiveFormatMessages{
		rawMessage:               data,
		processedMessage:         storage.processedMessage,
		allowedTransfer:          storage.allowedTransfer,
		isProcessedMisp:          storage.isProcessedMisp,
		isProcessedElasticsearsh: storage.isProcessedElasticsearsh,
		isProcessedNKCKI:         storage.isProcessedNKCKI,
	}
}

// GetRawDataHiveFormatMessage возвращает из хранилища сырые данные полученные от TheHive
func (cst *CommonStorageTemporary) GetRawDataHiveFormatMessage(uuid string) ([]byte, bool) {
	if s, ok := cst.HiveFormatMessage.Storages[uuid]; ok {
		return s.rawMessage, true
	}

	return []byte{}, false
}

// SetProcessedDataHiveFormatMessage добавляет данные частично разобранные Unmarshal JSON
func (cst *CommonStorageTemporary) SetProcessedDataHiveFormatMessage(uuid string, data map[string]interface{}) {
	cst.HiveFormatMessage.Lock()
	defer cst.HiveFormatMessage.Unlock()

	if _, ok := cst.HiveFormatMessage.Storages[uuid]; !ok {
		cst.HiveFormatMessage.Storages[uuid] = StorageHiveFormatMessages{}
	}

	storage := cst.HiveFormatMessage.Storages[uuid]
	cst.HiveFormatMessage.Storages[uuid] = StorageHiveFormatMessages{
		rawMessage:               storage.rawMessage,
		processedMessage:         data,
		allowedTransfer:          storage.allowedTransfer,
		isProcessedMisp:          storage.isProcessedMisp,
		isProcessedElasticsearsh: storage.isProcessedElasticsearsh,
		isProcessedNKCKI:         storage.isProcessedNKCKI,
	}
}

// GetProcessedDataHiveFormatMessage возвращает из хранилища сырые данные полученные от TheHive
func (cst *CommonStorageTemporary) GetProcessedDataHiveFormatMessage(uuid string) (map[string]interface{}, bool) {
	if s, ok := cst.HiveFormatMessage.Storages[uuid]; ok {
		return s.processedMessage, true
	}

	return map[string]interface{}{}, false
}

// SetAllowedTransferTrueHiveFormatMessage устанавливает поле информирующее о разрешении пропуска события в TRUE
func (cst *CommonStorageTemporary) SetAllowedTransferTrueHiveFormatMessage(uuid string) {
	cst.HiveFormatMessage.Lock()
	defer cst.HiveFormatMessage.Unlock()

	if _, ok := cst.HiveFormatMessage.Storages[uuid]; !ok {
		cst.HiveFormatMessage.Storages[uuid] = StorageHiveFormatMessages{}
	}

	storage := cst.HiveFormatMessage.Storages[uuid]
	cst.HiveFormatMessage.Storages[uuid] = StorageHiveFormatMessages{
		rawMessage:               storage.rawMessage,
		processedMessage:         storage.processedMessage,
		allowedTransfer:          true,
		isProcessedMisp:          storage.isProcessedMisp,
		isProcessedElasticsearsh: storage.isProcessedElasticsearsh,
		isProcessedNKCKI:         storage.isProcessedNKCKI,
	}
}

// SetAllowedTransferFalseHiveFormatMessage устанавливает поле информирующее о разрешении пропуска события в FALSE
func (cst *CommonStorageTemporary) SetAllowedTransferFalseHiveFormatMessage(uuid string) {
	if _, ok := cst.HiveFormatMessage.Storages[uuid]; !ok {
		cst.HiveFormatMessage.Storages[uuid] = StorageHiveFormatMessages{}
	}

	storage := cst.HiveFormatMessage.Storages[uuid]
	cst.HiveFormatMessage.Storages[uuid] = StorageHiveFormatMessages{
		rawMessage:               storage.rawMessage,
		processedMessage:         storage.processedMessage,
		allowedTransfer:          false,
		isProcessedMisp:          storage.isProcessedMisp,
		isProcessedElasticsearsh: storage.isProcessedElasticsearsh,
		isProcessedNKCKI:         storage.isProcessedNKCKI,
	}
}

// GetAllowedTransferHiveFormatMessage возвращает значение поля информирующее о разрешении пропуска события
func (cst *CommonStorageTemporary) GetAllowedTransferHiveFormatMessage(uuid string) (bool, bool) {
	if s, ok := cst.HiveFormatMessage.Storages[uuid]; ok {
		return s.allowedTransfer, true
	}

	return false, false
}

func (cst *CommonStorageTemporary) SetIsProcessedMispHiveFormatMessage(uuid string) bool {
	if _, ok := cst.HiveFormatMessage.Storages[uuid]; !ok {
		return false
	}

	storage := cst.HiveFormatMessage.Storages[uuid]
	cst.HiveFormatMessage.Storages[uuid] = StorageHiveFormatMessages{
		rawMessage:               storage.rawMessage,
		processedMessage:         storage.processedMessage,
		allowedTransfer:          storage.allowedTransfer,
		isProcessedMisp:          true,
		isProcessedElasticsearsh: storage.isProcessedElasticsearsh,
		isProcessedNKCKI:         storage.isProcessedNKCKI,
	}

	return true
}

func (cst *CommonStorageTemporary) SetIsProcessedElasticsearshHiveFormatMessage(uuid string) bool {
	if _, ok := cst.HiveFormatMessage.Storages[uuid]; !ok {
		return false
	}

	storage := cst.HiveFormatMessage.Storages[uuid]
	cst.HiveFormatMessage.Storages[uuid] = StorageHiveFormatMessages{
		rawMessage:               storage.rawMessage,
		processedMessage:         storage.processedMessage,
		allowedTransfer:          storage.allowedTransfer,
		isProcessedMisp:          storage.isProcessedMisp,
		isProcessedElasticsearsh: true,
		isProcessedNKCKI:         storage.isProcessedNKCKI,
	}

	return true
}

func (cst *CommonStorageTemporary) SetIsProcessedNKCKIHiveFormatMessage(uuid string) bool {
	if _, ok := cst.HiveFormatMessage.Storages[uuid]; !ok {
		return false
	}

	storage := cst.HiveFormatMessage.Storages[uuid]
	cst.HiveFormatMessage.Storages[uuid] = StorageHiveFormatMessages{
		rawMessage:               storage.rawMessage,
		processedMessage:         storage.processedMessage,
		allowedTransfer:          storage.allowedTransfer,
		isProcessedMisp:          storage.isProcessedMisp,
		isProcessedElasticsearsh: storage.isProcessedElasticsearsh,
		isProcessedNKCKI:         true,
	}

	return true
}

// AddUserSettingsMISP добавляет настройки для пользователя
func (cst *CommonStorageTemporary) AddUserSettingsMISP(usmisp UserSettingsMISP) {
	cst.ListUserSettingsMISP = append(cst.ListUserSettingsMISP, usmisp)
}

// возвращает настройки пользователя по его email
func (cst *CommonStorageTemporary) GetUserSettingsMISP(email string) (UserSettingsMISP, bool) {
	for _, v := range cst.ListUserSettingsMISP {
		if v.Email == email {
			return v, true
		}
	}

	return UserSettingsMISP{}, false
}

// возвращает весь список настроек пользователей
func (cst *CommonStorageTemporary) GetListUserSettingsMISP() *CommonStorageTemporary {
	return cst
}

func deleteHiveFormatMessageElement(uuid string, cst *CommonStorageTemporary) {
	cst.HiveFormatMessage.Lock()
	defer cst.HiveFormatMessage.Unlock()

	delete(cst.HiveFormatMessage.Storages, uuid)
}

func deleteTemporaryCase(id int, cst *CommonStorageTemporary) {
	cst.temporaryInputCase.Lock()
	defer cst.temporaryInputCase.Unlock()

	delete(cst.temporaryInputCase.Cases, id)
}
