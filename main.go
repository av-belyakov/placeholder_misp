package main

import (
	"fmt"
	"log"
	"os"

	"placeholder_misp/confighandler"
)

var (
	err     error
	confApp confighandler.ConfigApp
)

func init() {
	fmt.Println("func 'init', START...")

	// + 1. Прочитать переменные окружения, пока одну
	// + 2. Инициировать модуль для чтения конфигурационных файлов. При этом сначало читается общий конфиг, а затем
	// тот конфиг, выбор которого зависит от переменной окружения 'GO_PH_MISP_MAIN'
	//3. Инициализировать обработчик ошибок (запись логов) или отправка их на stdout
	//4. Инициализировать модуль соединения с NATS
	//5. Инициализировать модуль соединения с MISP
	//6. Инициализировать модуль обработчик

	confApp, err = confighandler.NewConfig()
	if err != nil {
		log.Fatalf("error module 'confighandler': %v\n", err)
	}
}

func main() {
	fmt.Println("func 'main', START...")
	fmt.Println("config application:", confApp)

	i, err := os.Stdout.Write([]byte("test writing to stdout"))
	fmt.Println("os.Stdout.Write i = ", i, " error = ", err)
}
