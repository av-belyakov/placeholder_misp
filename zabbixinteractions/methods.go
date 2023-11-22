package zabbixinteractions

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"
)

func (hzc *HandlerZabbixConnection) SendData(data []string) (int, error) {
	if len(data) == 0 {
		return 0, fmt.Errorf("the list of transmitted data should not be empty")
	}

	ldz := make([]DataZabbix, 0, len(data))
	for _, v := range data {
		ldz = append(ldz, DataZabbix{
			Host:  hzc.ZabbixHost,
			Key:   hzc.ZabbixKey,
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

	//заголовок пакета
	pkg := []byte("ZBXD\x01")

	//длинна пакета с данными
	dataLen := make([]byte, 8)
	binary.LittleEndian.PutUint32(dataLen, uint32(len(jsonReg)))

	pkg = append(pkg, dataLen...)
	pkg = append(pkg, jsonReg...)

	var d net.Dialer
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	conn, err := d.DialContext(ctx, "tcp", hzc.NetHost)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	num, err := conn.Write(pkg)
	if err != nil {
		return 0, err
	}

	_, err = io.ReadAll(conn)
	if err != nil {
		return num, err
	}

	return num, nil
}
