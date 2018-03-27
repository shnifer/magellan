package commons

import "encoding/json"

type State struct{
	Special string
	ShipID string
	GalaxyID string
}

func (s State) Encode()string{
	buf,err:=json.Marshal(s)
	if err!=nil{
		panic(err)
	}
	return string(buf)
}

func (State) Decode(data string) State{
	state:=State{}
	err:=json.Unmarshal([]byte(data), &state)
	if err!=nil{
		panic(err)
	}
	return state
}
