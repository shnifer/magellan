package wrnt

import (
	"errors"
)

var ErrNotInited = errors.New("not inited")

type Send struct {
	storage    *Storage
	confirmedN int
	inited     bool
}

func NewSend() *Send {
	return &Send{
		storage: newStorage(),
	}
}

func (s *Send) init(lastConfirmedN int) {
	if s.inited {
		panic("Send.Init: Already inited!")
	}
	s.inited = true
	s.confirmedN = lastConfirmedN
	s.storage.BaseN = lastConfirmedN + 1
}

func (s *Send) AddItems(items ...string) {
	s.storage.add(items...)
}

func (s *Send) Pack() (msg Storage, err error) {
	if !s.inited {
		return Storage{}, ErrNotInited
	}

	return Storage{
		BaseN: s.confirmedN + 1,
		Items: s.storage.get(s.confirmedN + 1),
	}, nil
}

func (s *Send) confirm(n int) {
	if !s.inited {
		s.init(n)
	}
	if n > s.confirmedN {
		s.confirmedN = n
	}
}

func (s *Send) Confirm(n int) {
	s.confirm(n)
	s.storage.cut(s.confirmedN)
}

func (s *Send) DropNotSent() {
	s.confirmedN = s.storage.BaseN + len(s.storage.Items) - 1
}
