package wrnt

import "sync"

type SendMany struct {
	mu      sync.RWMutex
	storage *storage
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

func (sm *SendMany) AddName(name string) {
	sm.mu.Lock()

	if _, ok := sm.sends[name]; ok {
		panic("name " + name + " is already in list")
	}
	send := NewSend()
	send.storage = sm.storage
	sm.sends[name] = send

	sm.mu.Unlock()
}

func (sm *SendMany) AddItems(items ...string) {
	sm.mu.RLock()
	sm.storage.add(items...)
	sm.mu.RUnlock()
}

func (sm *SendMany) Pack(name string) (msg Message, err error) {
	sm.mu.RLock()
	send, ok := sm.sends[name]
	sm.mu.RUnlock()

	if !ok {
		panic("SendMany.Pack: Unknown name " + name)
	}
	return send.Pack()
}

func (sm *SendMany) Confirm(name string, n int) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	send, ok := sm.sends[name]

	if !ok {
		panic("SendMany.Confirm: Unknown name " + name)
	}
	send.confirm(n)

	sm.storage.mu.RLock()
	baseN := sm.storage.BaseN
	minN := sm.storage.BaseN + len(sm.storage.Items) - 1
	sm.storage.mu.RUnlock()

	for _, send := range sm.sends {
		send.mu.Lock()
		conf := send.confirmedN
		send.mu.Unlock()

		if conf < baseN {
			return
		}
		if conf < minN {
			minN = conf
		}
	}
	sm.storage.cut(minN)
}

func (sm *SendMany) DropNotSent() {
	lastN := sm.storage.BaseN + len(sm.storage.Items) - 1
	for _, send := range sm.sends {
		send.confirmedN = lastN
	}
}
