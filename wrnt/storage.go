package wrnt

import (
	. "github.com/Shnifer/magellan/log"
	"sync"
)

type Message struct {
	BaseN int
	Items []string
}

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

func (s *storage) add(items ...string) {
	s.mu.Lock()
	s.Items = append(s.Items, items...)
	s.mu.Unlock()
}
func (s *storage) get(fromN int) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if fromN < s.BaseN {
		//todo:check this again
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

func (s *storage) cut(toN int) {
	s.mu.Lock()

	startInd := toN - s.BaseN + 1
	s.Items = s.Items[startInd:]
	s.BaseN += startInd

	s.mu.Unlock()
}
