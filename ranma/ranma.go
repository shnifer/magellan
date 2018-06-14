package ranma

import (
	"net/http"
	"sync"
	"time"
)

const sCount = 8

type Ranma struct {
	addr           string
	dropInOnRepair bool
	client         *http.Client
	corrected      [sCount]*system

	mu         sync.Mutex
	programmed [sCount]uint16
	broken     [sCount]bool
}

func NewRanma(addr string, dropInOnRepair bool, timeoutMs int, depth int) *Ranma {
	res := &Ranma{
		addr:           addr,
		dropInOnRepair: dropInOnRepair,
		client: &http.Client{
			Timeout: time.Duration(timeoutMs) * time.Millisecond,
		},
	}
	for i := 0; i < sCount; i++ {
		res.corrected[i] = newSystem(depth)
	}
	res.reset()
	go reqDaemon(res)
	return res
}

func (r *Ranma) GetIn(sn int) uint16 {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.programmed[sn]
}

func (r *Ranma) GetInByte(sn, bn int) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	return (r.programmed[sn]>>uint(bn))&1 > 0
}

func (r *Ranma) SetIn(sn int, x uint16) {
	r.mu.Lock()
	defer r.mu.Unlock()

	xorV := r.programmed[sn] ^ x
	r.programmed[sn] = x
	r.corrected[sn].hardXor(xorV)
	if xorV != 0 {
		r.broken[sn] = true
	}
	go r.send(sn, x)
}

func (r *Ranma) XorIn(sn int, x uint16) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.programmed[sn] = x ^ r.programmed[sn]
	r.corrected[sn].hardXor(x)
	if x != 0 {
		r.broken[sn] = true
	}
	go r.send(sn, r.programmed[sn])
}

func (r *Ranma) XorInByte(sn, bn int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	xorV := uint16(1 << uint(bn))
	r.programmed[sn] = xorV ^ r.programmed[sn]
	r.corrected[sn].hardXor(xorV)
	if xorV != 0 {
		r.broken[sn] = true
	}
	go r.send(sn, r.programmed[sn])
}

func (r *Ranma) GetOut(sn int) uint16 {
	if !r.broken[sn] {
		return 0
	}
	return r.corrected[sn].get()
}

func (r *Ranma) GetOutByte(sn, bn int) bool {
	if !r.broken[sn] {
		return false
	}
	return r.corrected[sn].getByte(bn)
}

func (r *Ranma) reset() {
	for i := 0; i < sCount; i++ {
		r.send(i, 0)
	}
}

func reqDaemon(r *Ranma) {
	tick := time.Tick(time.Second)
	for {
		<-tick
		for i := 0; i < sCount; i++ {
			r.recv(i)
			v, isHardware := r.corrected[i].getWithFlag()
			if v == 0 && isHardware {
				//clear broken
				r.mu.Lock()
				r.broken[i] = false
				if r.dropInOnRepair {
					r.programmed[i] = 0
				}
				r.mu.Unlock()
			}
		}
	}
}
