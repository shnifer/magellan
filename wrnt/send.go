package wrnt

import (
	"errors"
	"sync"
)

var ErrNotInited = errors.New("not inited")

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

func (s *Send) AddItems(items ...string) {
	s.mu.Lock()
	s.storage.add(items...)
	s.mu.Unlock()
}

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
func (s *Send) init(lastConfirmedN int) {
	if s.inited {
		panic("Send.Init: Already inited!")
	}
	s.inited = true
	s.confirmedN = lastConfirmedN
	s.storage.BaseN = lastConfirmedN + 1
}

func (s *Send) Confirm(n int) {
	s.confirm(n)
	s.storage.cut(s.confirmedN)
}

func (s *Send) DropNotSent() {
	s.mu.Lock()
	s.confirmedN = s.storage.BaseN + len(s.storage.Items) - 1
	s.mu.Unlock()
}
