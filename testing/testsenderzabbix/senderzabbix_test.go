package testsenderzabbix_test

import (
	"context"
	"fmt"
	"net"
	"time"

	"placeholder_misp/zabbixinteractions"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*type HandlerZabbixConnection struct {
	Host string

}

type PatternZabbix struct {
	Request string       `json:"request"`
	Data    []DataZabbix `json:"data"`
}

type DataZabbix struct {
	Host  string `json:"host"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

func NewHandlerZabbixConnection(host string) *HandlerZabbixConnection {
	return &HandlerZabbixConnection{
		Host: host,
	}
}

func (hzc *HandlerZabbixConnection) SendData(data []string) (int, error) {
	if len(data) == 0 {
		return 0, fmt.Errorf("the list of transmitted data should not be empty")
	}

	ldz := make([]DataZabbix, 0, len(data))
	for _, v := range data {
		ldz = append(ldz, DataZabbix{
			Host:  "sib-server",
			Key:   "test_bav",
			Value: v,
		})
	}

	jsonReg, err := json.Marshal(PatternZabbix{
		Request: "sender data",
		Data:    ldz,
	})
	if err != nil {
		return 0, err
	}

	pkg := []byte("ZBXD\x01")

	dataLen := make([]byte, 8)
	binary.LittleEndian.PutUint32(dataLen, uint32(len(jsonReg)))

	pkg = append(pkg, dataLen...)
	pkg = append(pkg, jsonReg...)

	var d net.Dialer
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	conn, err := d.DialContext(ctx, "tcp", hzc.Host)
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	num, err := conn.Write(pkg)
	if err != nil {
		return 0, err
	}

	buf, read_err := io.ReadAll(conn)
	if read_err != nil {
		fmt.Println("failed:", read_err)
	}
	fmt.Println(string(buf))

	return num, nil
}*/

var _ = Describe("Senderzabbix", Ordered, func() {
	//var zc *HandlerZabbixConnection
	var zc *zabbixinteractions.HandlerZabbixConnection

	BeforeAll(func() {
		//zc = NewHandlerZabbixConnection("zabbix.cloud.gcm:10051")
		zc = zabbixinteractions.NewHandlerZabbixConnection("zabbix.cloud.gcm:10051", "sib-server", "test_bav")
	})

	Context("Тест 1. Пробуем выполнить соединение с Zabbix", func() {
		It("Соединение с Zabbix должно быть успешно установлено", func() {
			var d net.Dialer
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			conn, err := d.DialContext(ctx, "tcp", "zabbix.cloud.gcm:10051")
			Expect(err).ShouldNot(HaveOccurred())
			defer conn.Close()

			Expect(true).Should(BeTrue())
		})
	})

	Context("Тест 2. Проверяем возможность подключения и отправки данных в Zabbix", func() {
		It("При отправки данных в Zabbix не должно быть ошибок", func() {

			//для подтверждения что модуль
			num, err := zc.SendData([]string{"I'm still alive"})

			fmt.Println("Count sended byte:", num)

			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
