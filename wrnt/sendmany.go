package wrnt

import "github.com/pkg/errors"

type SendMany struct {
	sends map[string]*Send
}

func NewSendMany(names []string) *SendMany {
	sends := make(map[string]*Send, len(names))
	for _, name := range names {
		sends[name] = NewSend()
	}
	return &SendMany{
		sends: sends,
	}
}

func (sm *SendMany) AddItems(items ...string) {
	for _, send := range sm.sends {
		send.AddItems(items...)
	}
}

func (sm *SendMany) Pack(name string) (msg Message, err error) {
	send, ok := sm.sends[name]
	if !ok {
		return Message{}, errors.New("SendMany.Pack: unknown name " + name)
	}
	return send.Pack()
}

func (sm *SendMany) Confirm(name string, n int) error {

	send, ok := sm.sends[name]

	if !ok {
		return errors.New("SendMany.Confirm: Unknown name " + name)
	}
	send.Confirm(n)
	return nil
}

func (sm *SendMany) DropNotSent() {
	for _, send := range sm.sends {
		send.DropNotSent()
	}
}
