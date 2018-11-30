package wrnt

import (
	. "github.com/shnifer/magellan/log"
	"sync"
)

//message to be send from send to recv
//holds all non-confirmed items
type Message struct {
	//starting index
	BaseN int
	//items, with index from
	Items []string
}

//storage is a send-side storage of items
//concurrent-safe
type storage struct {
	mu sync.RWMutex
	Message
}

func newStorage() *storage {
	return &storage{
		Message: Message{
			BaseN: 0,
			Items: make([]string, 0),
		},
	}
}

//adds new items to storage
func (s *storage) add(items ...string) {
	s.mu.Lock()
	s.Items = append(s.Items, items...)
	s.mu.Unlock()
}

//get returns a copy of storage part with index>=fromN
//if fromN is not in range of [s.BaseN;s.BaseN+len(s.Items)]
//get returns empty []string without error
func (s *storage) get(fromN int) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if fromN < s.BaseN {
		Log(LVL_ERROR, "storage.get fromN<BaseN: ", fromN, "<", s.BaseN)
		return []string{}
	}
	count := len(s.Items) - (fromN - s.BaseN)
	if count < 0 {
		Log(LVL_ERROR, "storage.get count<0")
		return []string{}
	}
	res := make([]string, count)
	copy(res, s.Items[fromN-s.BaseN:])
	return res
}

//cut remove items with index<=toN
//if toN is not valid do nothing
func (s *storage) cut(toN int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	startInd := toN - s.BaseN + 1
	if startInd < 0 {
		Log(LVL_ERROR, "storage.cut startInd ", startInd, "<0")
		return
	} else if startInd > len(s.Items) {
		Log(LVL_ERROR, "storage.cut startInd<len(storage.Items) ", startInd, "<", len(s.Items))
		return
	}
	s.Items = s.Items[startInd:]
	s.BaseN += startInd
}

//set baseN manually
//may be called only once, if BaseN iz zero-value
func (s *storage) initBaseN(baseN int) {
	s.mu.Lock()
	if s.BaseN == 0 {
		s.BaseN = baseN
	} else {
		Log(LVL_ERROR, "storage.initBaseN called for storage with non-zero baseN=", s.BaseN)
	}
	s.mu.Unlock()
}
