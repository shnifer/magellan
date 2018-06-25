package commons

import "encoding/json"

type State struct {
	StateID  string `json:"st"`
	ShipID   string `json:"sh"`
	GalaxyID string `json:"gx"`
}

func (s State) Encode() string {
	buf, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	return string(buf)
}

func (State) Decode(data string) State {
	state := State{}
	err := json.Unmarshal([]byte(data), &state)
	if err != nil {
		panic(err)
	}
	return state
}
