package redisapi

import (
	"github.com/av-belyakov/placeholder_misp/commoninterfaces"
	"github.com/redis/go-redis/v9"
)

// ModuleRedis инициализированный модуль
type ModuleRedis struct {
	client   *redis.Client           //дескриптор соединения с Redis БД
	logger   commoninterfaces.Logger //интерфейс логгирования
	chInput  chan SettingsInput      //канал для отправки данных В модуль
	chOutput chan SettingsOutput     //канал для отправки данных ИЗ модуля
	host     string
	port     int
}

type SettingsInput struct {
	RawData []byte
	Command string
	Data    string
}

type SettingsOutput struct {
	Result        interface{}
	CommandResult string
}
