package natsapi

import (
	"fmt"

	"github.com/nats-io/nats.go"
)

// SendRequestCommandExecute выполняет отправку запроса с командой в NATS
func SendRequestCommandExecute(nc *nats.Conn, listenerCommand string, data InputSettings) error {
	if err := nc.Publish(listenerCommand, data.Data); err != nil {
		return fmt.Errorf("error processing command '%s' for caseId '%s' (rootId '%s') '%s' %s:%d", data.Command, data.CaseId, data.RootId, err)
	}

	return nil
}
