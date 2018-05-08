package commons

import "sync"

type TData struct {
	//Main game loop will handle this
	//Mu sync.RWMutex
	State
	StateData
	CommonData

	//StateBSP * Common.Engi.BSPDegrade
	SP *BSP

	actionQ chan func()

	mu         sync.Mutex
	partToSend CommonData
}

func NewData() TData {
	return TData{
		CommonData: CommonData{}.Empty(),
		actionQ:    make(chan func(), 8),
		SP:         &BSP{},
	}
}

//Main cycle
func (d *TData) Update(roleName string) {
	defer LogFunc("Data.Update")()

loop:
	for {
		select {
		case f := <-d.actionQ:
			f()
		default:
			break loop
		}
	}

	d.mu.Lock()
	d.partToSend = d.CommonData.Part(roleName).Copy()
	d.mu.Unlock()
}

//Network cycle
func (d *TData) SetState(state State) {
	d.actionQ <- func() {
		d.State = state
	}
}

//Network cycle
func (d *TData) SetStateData(stateData StateData) {
	d.actionQ <- func() {
		d.StateData = stateData
	}
}

//Network cycle
func (d *TData) LoadCommonData(src CommonData) {
	d.actionQ <- func() {
		src.FillNotNil(&d.CommonData)
		d.SP = BSP{}.CalcDegrade(d.BSP, d.EngiData.BSPDegrade)
	}
}

//Network cycle
func (d *TData) WaitDone() {
	defer LogFunc("Data.WaitDone")()
	done := make(chan struct{})
	d.actionQ <- func() {
		close(done)
	}
	<-done
}

func (d *TData) GetState() State {
	stateCh := make(chan State)
	defer close(stateCh)
	d.actionQ <- func() {
		stateCh <- d.State
	}
	return <-stateCh
}

func (d *TData) GetStateData() StateData {
	stateDataCh := make(chan StateData)
	defer close(stateDataCh)
	d.actionQ <- func() {
		stateDataCh <- d.StateData.Copy()
	}
	return <-stateDataCh
}

//Network cycle
func (d *TData) MyPartToSend() []byte {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.partToSend.Encode()
}

func (d *TData) Encode() {
	panic("Don't do it! use methods of embeded structs")
}

func (d *TData) Decode() {
	panic("Don't do it! use methods of embeded structs")
}
