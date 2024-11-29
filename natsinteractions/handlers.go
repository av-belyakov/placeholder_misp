package natsinteractions

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"time"

	"github.com/nats-io/nats.go"
)

// SendRequestCommandExecute выполняет отправку запроса с командой в NATS
func SendRequestCommandExecute(nc *nats.Conn, listenerCommand string, data SettingsInputChan) (string, error) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer ctxCancel()

	res, err := nc.RequestWithContext(ctx, listenerCommand, data.Data)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return "", fmt.Errorf("error processing command '%s' for caseId '%s' (rootId '%s') '%s' %s:%d", data.Command, data.CaseId, data.RootId, err.Error(), f, l-2)
	}

	resToComm := ResponseToCommand{}
	if err = json.Unmarshal(res.Data, &resToComm); err != nil {
		_, f, l, _ := runtime.Caller(0)
		return "", fmt.Errorf("%s %s:%d", err.Error(), f, l-2)
	}

	if resToComm.Error != "" {
		return "", fmt.Errorf("command '%s' for caseId '%s' (rootId '%s') return status code '%d' with error message '%s'", data.Command, data.CaseId, data.RootId, resToComm.StatusCode, resToComm.Error)
	}

	return fmt.Sprintf("command '%s' for caseId '%s' (rootId '%s') return status code '%d'", data.Command, data.CaseId, data.RootId, resToComm.StatusCode), nil
}
