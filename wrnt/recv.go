package wrnt

import (
	"log"
	"sync"
)

type Recv struct {
	mu    sync.Mutex
	lastN int
}

func NewRecv() *Recv {
	return &Recv{
		lastN: 0,
	}
}

func (r *Recv) Unpack(msg Message) []string {
	r.mu.Lock()
	defer r.mu.Unlock()

	firstInd := r.lastN - msg.BaseN + 1
	if firstInd < 0 {
		firstInd = 0
	} else if firstInd > len(msg.Items) {
		firstInd = len(msg.Items)
	}

	msgLast := msg.BaseN + len(msg.Items) - 1
	if msgLast > r.lastN {
		r.lastN = msgLast
	}

	return msg.Items[firstInd:]
}

func (r *Recv) LastRecv() int {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.lastN
}
