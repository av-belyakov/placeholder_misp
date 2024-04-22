// Пакет zabbixinteractions предназначен для взаимодействия с API Zabbix
package zabbixinteractions

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"runtime"
	"time"
)

// NewZabbixConnection создает обработчик соединения с Zabbix
// ctx - должен быть context.WithCancel()
// settings - настройки
func NewZabbixConnection(ctx context.Context, settings SettingsZabbixConnection) (*HandlerZabbixConnection, error) {
	var hzc HandlerZabbixConnection

	if settings.Host == "" {
		return &hzc, fmt.Errorf("the value 'Host' should not be empty")
	}

	if settings.Port == 0 {
		return &hzc, fmt.Errorf("the value 'Port' should not be equal '0'")
	}

	if settings.ZabbixHost == "" {
		return &hzc, fmt.Errorf("the value 'ZabbixHost' should not be empty")
	}

	if settings.NetProto != "tcp" && settings.NetProto != "udp" {
		settings.NetProto = "tcp"
	}

	if settings.ConnectionTimeout == nil {
		t := time.Duration(5 * time.Second)
		settings.ConnectionTimeout = &t
	}

	hzc = HandlerZabbixConnection{
		ctx:         ctx,
		host:        settings.Host,
		port:        settings.Port,
		netProto:    settings.NetProto,
		zabbixHost:  settings.ZabbixHost,
		connTimeout: settings.ConnectionTimeout,
		chanErr:     make(chan error),
	}

	return &hzc, nil
}

// GetChanErr возвращает канал в который отправляются ошибки возникающие при соединении с Zabbix
func (hzc *HandlerZabbixConnection) GetChanErr() chan error {
	return hzc.chanErr
}

func (hzc *HandlerZabbixConnection) Handler(events []EventType, msgChan <-chan MessageSettings) error {
	countEvents := len(events)
	if countEvents == 0 {
		_, f, l, _ := runtime.Caller(0)
		return fmt.Errorf("'invalid configuration file for Zabbix, the number of event types (ZABBIX.zabbixHosts.eventTypes) is 0' %s:%d", f, l-1)
	}

	listChans := make(map[string]chan<- string, countEvents)

	go func() {
		<-hzc.ctx.Done()

		for _, channel := range listChans {
			close(channel)
		}

		listChans = nil
		close(hzc.chanErr)
	}()

	for _, v := range events {
		if !v.IsTransmit {
			continue
		}

		newChan := make(chan string)
		listChans[v.EventType] = newChan

		go func(cm <-chan string, zkey string, hs Handshake) {
			var t *time.Ticker
			if hs.TimeInterval > 0 && hs.Message != "" {
				t = time.NewTicker(time.Duration(hs.TimeInterval) * time.Minute)
				defer t.Stop()
			}

			if t == nil {
				for msg := range cm {
					if _, err := hzc.SendData(zkey, []string{msg}); err != nil {
						hzc.chanErr <- err
					}
				}
			} else {
				for {
					select {
					case <-t.C:
						if _, err := hzc.SendData(zkey, []string{hs.Message}); err != nil {
							hzc.chanErr <- err
						}

					case msg, open := <-cm:
						if !open {
							cm = nil

							return
						}

						if _, err := hzc.SendData(zkey, []string{msg}); err != nil {
							hzc.chanErr <- err
						}
					}
				}
			}
		}(newChan, v.ZabbixKey, v.Handshake)
	}

	go func() {
		for data := range msgChan {
			if c, ok := listChans[data.EventType]; ok {
				c <- data.Message
			}
		}
	}()

	return nil
}

func (hzc *HandlerZabbixConnection) SendData(zkey string, data []string) (int, error) {
	if len(data) == 0 {
		return 0, fmt.Errorf("the list of transmitted data should not be empty")
	}

	ldz := make([]DataZabbix, 0, len(data))
	for _, v := range data {
		ldz = append(ldz, DataZabbix{
			Host:  hzc.zabbixHost,
			Key:   zkey,
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

	var d net.Dialer = net.Dialer{}
	ctx, cancel := context.WithTimeout(context.Background(), *hzc.connTimeout)
	defer cancel()

	conn, err := d.DialContext(ctx, hzc.netProto, fmt.Sprintf("%s:%d", hzc.host, hzc.port))
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
