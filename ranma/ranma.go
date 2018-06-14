package ranma

import (
	"net/http"
	"sync"
	"time"
)

const sCount = 8

type Ranma struct {
	addr      string
	client    *http.Client
	corrected [sCount]*system

	mu         sync.Mutex
	programmed [sCount]uint16
}

func NewRanma(addr string, timeoutMs int, depth int) *Ranma {
	res := &Ranma{
		addr: addr,
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

	r.programmed[sn] = x
	go r.send(sn, x)
}

func (r *Ranma) XorIn(sn int, x uint16) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.programmed[sn] = x ^ r.programmed[sn]
	go r.send(sn, r.programmed[sn])
}

func (r *Ranma) XorInByte(sn, bn int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.programmed[sn] = (1 << uint(bn)) ^ r.programmed[sn]
	go r.send(sn, r.programmed[sn])
}

func (r *Ranma) GetOut(sn int) uint16 {
	return r.corrected[sn].get()
}

func (r *Ranma) GetOutByte(sn, bn int) bool {
	return r.corrected[sn].getByte(bn)
}

func (r *Ranma) reset() {
	for i := 0; i < sCount; i++ {
		r.send(i, 0)
	}
}

func reqDaemon(r *Ranma) {
	for {
		time.Sleep(time.Second)
		for i := 0; i < sCount; i++ {
			r.recv(i)
		}
	}
}
