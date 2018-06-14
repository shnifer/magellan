package ranma

import (
	"sync"
)

type system struct {
	sync.Mutex

	lastId       int
	dat          uint16
	fromHardware bool

	depth   int
	cursor  int
	history []uint16
}

func newSystem(depth int) *system {
	return &system{
		depth:   depth,
		history: make([]uint16, depth),
	}
}

func (s *system) getWithFlag() (dat uint16, fromHardware bool) {
	s.Lock()
	defer s.Unlock()

	return s.dat, s.fromHardware
}

func (s *system) get() uint16 {
	s.Lock()
	defer s.Unlock()

	return s.dat
}

func (s *system) getByte(b int) bool {
	s.Lock()
	defer s.Unlock()

	n := s.dat >> uint(b)
	return n&1 > 0
}

func (s *system) hardXor(x uint16) {
	s.Lock()
	defer s.Unlock()

	//hard set
	s.dat = x ^ s.dat
	s.fromHardware = false
	s.addHistory(s.dat)
}

//call internally, under lock
func (s *system) addHistory(x uint16) {
	s.cursor = (s.cursor + 1) % s.depth
	s.history[s.cursor] = x
}

func (s *system) setMsg(x uint16, id int) {
	s.Lock()
	defer s.Unlock()

	if s.lastId == id {
		return
	}
	s.lastId = id

	s.addHistory(x)

	//is history coherent -- approve as new value
	for i := 0; i < s.depth; i++ {
		if s.history[i] != x {
			return
		}
	}
	s.dat = x
	s.fromHardware = true
}
