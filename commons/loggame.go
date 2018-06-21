package commons

import (
	"encoding/json"
	"fmt"
	. "github.com/Shnifer/magellan/log"
	"github.com/Shnifer/magellan/network"
)

type LogGameEvent struct {
	Key         string
	Args        string
	StateFields string
}

//for clients
func ClientLogGame(Client *network.Client, key string, args ...interface{}) {
	LogGame(key, args)
	requestLogGame(Client, key, fmt.Sprint(args))
}

func requestLogGame(Client *network.Client, key string, args string) {
	lge := LogGameEvent{
		Key:         key,
		Args:        args,
		StateFields: GetLogStateFieldsStr(),
	}
	buf, err := json.Marshal(lge)
	if err != nil {
		Log(LVL_ERROR, "LogGameEvent can't marshal", err)
	}
	Client.SendRequest(CMD_LOGGAMEEVENT + string(buf))
}

func (LogGameEvent) Decode(data []byte) (res LogGameEvent, err error) {
	err = json.Unmarshal(data, &res)
	if err != nil {
		return LogGameEvent{}, err
	} else {
		return res, nil
	}
}
