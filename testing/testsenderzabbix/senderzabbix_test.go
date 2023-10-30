package testsenderzabbix_test

import (
	"encoding/json"
	"fmt"

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

func (hzc *HandlerZabbixConnection) SendData(data []string) error {
	if len(data) == 0 {
		return fmt.Errorf("the list of transmitted data should not be empty")
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
		return err
	}

package := []byte("ZBXD\1")

	return nil
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
*/

var _ = Describe("Senderzabbix", func() {

	Context("", func() {
		It("", func() {

		})
	})

	/*
		Context("", func(){
			It("", func(){

			})
		})
	*/
})
