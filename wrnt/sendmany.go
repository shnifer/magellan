package wrnt

import "log"

type SendMany struct {
	storage *Storage
	names   []string
	sends   map[string]*Send
}

func NewSendMany(names []string) *SendMany {
	storage := newStorage()
	sends := make(map[string]*Send, len(names))
	for _, name := range names {
		sends[name] = NewSend()
		sends[name].storage = storage
	}
	return &SendMany{
		storage: storage,
		names:   names,
		sends:   sends,
	}
}

func (sm *SendMany) AddItems(items ...string) {
	sm.storage.add(items...)
}

func (sm *SendMany) Pack(name string) (msg Storage, err error) {
	send, ok := sm.sends[name]
	if !ok {
		panic("SendMany.Pack: Unknown name " + name)
	}
	return send.Pack()
}

func (sm *SendMany) Confirm(name string, n int) {
	send, ok := sm.sends[name]
	if !ok {
		panic("SendMany.Confirm: Unknown name " + name)
	}
	send.confirm(n)

	minN := sm.storage.BaseN + len(sm.storage.Items) - 1
	for _, send := range sm.sends {
		conf := send.confirmedN
		if conf < sm.storage.BaseN {
			return
		}
		if conf < minN {
			minN = conf
		}
	}
	sm.storage.cut(minN)
	log.Println("sm.storage.len,.cap=", len(sm.storage.Items), cap(sm.storage.Items))
}

func (sm *SendMany) Reset() {
	sm.storage.doEmpty()
	for _, send := range sm.sends {
		send.inited = false
	}
}
