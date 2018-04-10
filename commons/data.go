package commons

import "sync"

type TData struct {
	Mu sync.RWMutex
	State
	StateData
	CommonData
}

func (d *TData) SetState(state State) {
	d.Mu.Lock()
	d.State = state
	d.Mu.Unlock()
}

func (d *TData) SetStateData(stateData StateData) {
	d.Mu.Lock()
	d.StateData = stateData
	d.Mu.Unlock()
}

func (d *TData) CommonPartEncoded(roleName string) []byte {
	d.Mu.RLock()
	defer d.Mu.RUnlock()
	return d.Part(roleName).Encode()
}

func (d *TData) LoadCommonData(src CommonData, roleName string, readAll bool) {
	d.Mu.Lock()
	if !readAll {
		src.ClearRole(roleName)
	}
	src.FillNotNil(&d.CommonData)
	d.Mu.Unlock()
}

func (d TData) Encode() {
	panic("Don't do it! use methods of embeded structs")
}

func (d TData) Decode() {
	panic("Don't do it! use methods of embeded structs")
}
