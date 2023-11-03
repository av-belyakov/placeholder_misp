package testsenderzabbix_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type HandlerZabbixConnection struct {
	Host, Key string
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

func NewHandlerZabbixConnection(host, key string) *HandlerZabbixConnection {
	return &HandlerZabbixConnection{
		Host: host,
		Key:  key,
	}
}

func (hzc *HandlerZabbixConnection) SendData(data []string) (int, error) {
	if len(data) == 0 {
		return 0, fmt.Errorf("the list of transmitted data should not be empty")
	}

	ldz := make([]DataZabbix, 0, len(data))
	for _, v := range data {
		ldz = append(ldz, DataZabbix{
			Host:  hzc.Host,
			Key:   hzc.Key,
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

	fmt.Println(jsonReg)

	pkg := append([]byte("ZBXD"), []uint8{0x01}...)

	fmt.Println("Lenght header = ", len(pkg))

	pkgDataLen := make([]byte, 0, 4)
	pkgDataLen = append(pkgDataLen, uint8(len(jsonReg)))

	fmt.Println("Lenght data len = ", len(pkgDataLen))

	pkg = append(pkg, pkgDataLen...)

	pkg = append(pkg, []uint8{0x00, 0x00, 0x00, 0x00}...)

	fmt.Println("Size pkg = ", len(pkg), " - ", pkg)

	pkg = append(pkg, jsonReg...)

	//	fmt.Println("1 to uint8 =", uint8(1))
	//	fmt.Println("1 to uint8 =", uint8(len(jsonReg)))

	var d net.Dialer
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	conn, err := d.DialContext(ctx, "tcp", "zabbix.cloud.gcm:10051")
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	num, err := conn.Write(pkg)
	if err != nil {
		return 0, err
	}

	return num, nil
}

/*
def zabbix_sender(key,value):
    data_json = {
        "request": "sender data",
        "data": [
            {
                "host": "sib-server",
                "key": "test_bav",
                "value": value
            }
        ]
    }
    data = json.dumps(data_json, separators=(',', ':')).encode()
    packet = b"ZBXD\1" + struct.pack('<Q', len(data)) + data
    try:
        s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        s.connect((zabbix_ip_or_domain, zabbix_port))
        s.sendall(packet)
        data = s.recv(1024)
        s.close()
        logger.debug(f"Received', repr({repr(data)})")
    except Exception as e:
        logger.error(f"Не отправлено в заббикс({e})")

		//для боевого использования
		zabbix_ip_or_domain : "zabbix.cloud.gcm"
  controlled_host : "sib-server"
  zabbix_port : 10051
*/

var _ = Describe("Senderzabbix", Ordered, func() {
	var hzc *HandlerZabbixConnection

	BeforeAll(func() {
		hzc = NewHandlerZabbixConnection("zabbix.cloud.gcm:10051", "test_bav")
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
			num, err := hzc.SendData([]string{"simply not many test"})

			fmt.Println("Count sended byte:", num)

			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
