package wrnt

import (
	"errors"
	"sync"
)

var ErrNotInited = errors.New("not inited")

//Send is a send-side component of wrnt package
//new Send doesn't count as inited before first Confirm call
//so Confirm with last received item must be guaranteed
type Send struct {
	mu         sync.Mutex
	storage    *storage
	confirmedN int
	inited     bool
}

func NewSend() *Send {
	return &Send{
		storage: newStorage(),
	}
}

//AddItems add new items to send
//may be called before initiation
//all added items are not confirmed
//and must be sent
func (s *Send) AddItems(items ...string) {
	s.mu.Lock()
	s.storage.add(items...)
	s.mu.Unlock()
}

//Pack returns a Message according to ConfirmedN
//user must transport it to receiver-side
//if send is not inited Pack returns ErrNotInited error
func (s *Send) Pack() (msg Message, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.inited {
		return Message{}, ErrNotInited
	}

	return Message{
		BaseN: s.confirmedN + 1,
		Items: s.storage.get(s.confirmedN + 1),
	}, nil
}

//Confirm the fact of receiving messages up to index N
//after Confirm call Send is inited
func (s *Send) Confirm(n int) {
	s.confirm(n)
	s.storage.cut(s.confirmedN)
}

//DropNotSent moves ConfirmedN to the end of storage
//marking all items as confirmed.
//by fact, they may be already received or not received
//but they will never be sent again
func (s *Send) DropNotSent() {
	s.mu.Lock()
	if s.inited {
		s.confirmedN = s.storage.BaseN + len(s.storage.Items) - 1
	} else {
		s.storage = newStorage()
	}
	s.mu.Unlock()
}

//init Send
//if n is less than already confirmed - do nothing
func (s *Send) confirm(n int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.inited {
		s.init(n)
	}
	if n > s.confirmedN {
		s.confirmedN = n
	}
}

//runned from confirm under mutex
//executed only once.
//saves confirmedN and moves storage.BaseN accordingly
func (s *Send) init(lastConfirmedN int) {
	if s.inited {
		panic("Send.Init: Already inited!")
	}
	s.inited = true
	s.confirmedN = lastConfirmedN
	s.storage.initBaseN(lastConfirmedN + 1)
}
