package ranma

import (
	"sync"
)

type system struct {
	sync.Mutex

	lastId int
	dat    uint16

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

func (s *system) xor(x uint16) {
	s.Lock()
	defer s.Unlock()

	s.dat = x ^ s.dat
}

func (s *system) set(x uint16, id int) {
	s.Lock()
	defer s.Unlock()

	if s.lastId == id {
		return
	}
	s.lastId = id

	s.cursor = (s.cursor + 1) % s.depth
	s.history[s.cursor] = x

	for i := 0; i < s.depth; i++ {
		if s.history[i] != x {
			return
		}
	}

	s.dat = x
}
