package guarantee

import "log"

type Storage struct {
	BaseN int
	Items []string
}

func newStorage() *Storage {
	return &Storage{
		BaseN: 0,
		Items: make([]string, 0),
	}
}

func (s *Storage) add(items ...string) {
	s.Items = append(s.Items, items...)
}
func (s *Storage) get(fromN int) []string {
	if fromN < s.BaseN {
		log.Panicln("Storage.get fromN<BaseN", fromN, "<", s.BaseN)
	}
	count := len(s.Items) - (fromN - s.BaseN)
	if count < 0 {
		log.Println("Storage.get count<0")
	}
	res := make([]string, count)
	copy(res, s.Items[fromN-s.BaseN:])
	return res
}

func (s *Storage) cut(toN int) {
	startInd := toN - s.BaseN + 1
	s.Items = s.Items[startInd:]
	s.BaseN += startInd
}
