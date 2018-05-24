package storage

import (
	"github.com/peterbourgon/diskv"
	"strings"
	"sync"
)

const EventChanSize = 128
const separator = "~"

const (
	Add = iota
	Remove
)

type Event struct{
	Type int
	Key string
	Data []byte
}

// strings must not have "/","\" or "~" symbols
func newKey(area, node, key string) string{
	return strings.Join([]string{area,node,key}, separator )
}

func splitKey(fullKey string) (area,node,key string) {
	parts:=strings.Split(fullKey,separator)
	if len(parts)!=3 {
		panic("incorrect fullKey "+ fullKey)
	}
	return parts[0],parts[1],parts[2]
}

type Storage struct{
	sync.RWMutex

	node string

	disk *diskv.Diskv
	subs map[chan Event]string
}

func New(nodeName string, diskOpts diskv.Options) *Storage{
	disk:=diskv.New(diskOpts)

	return &Storage{
		disk: disk,
		node: nodeName,
	}
}

func (s *Storage) Add(area, key string, val []byte) {
	s.RLock()
	defer s.RUnlock()

	fk:=newKey(area, s.node, key)

	s.disk.Write(fk, val)

	s.event(Add, fk, val)
}

func (s *Storage) Remove(key string) {
	s.RLock()
	defer s.RUnlock()

	s.disk.Erase(key)

	s.event(Remove, key, nil)
}

func (s *Storage) event (t int, key string, val []byte) {
	Event:=Event{
		Type:t,
		Key:key,
		Data: val,
	}

	for ch, prefix:=range s.subs{
		if strings.HasPrefix(key, prefix){
			ch<-Event
		}
	}
}

func (s *Storage) Subscribe(prefix string) (data map[string][]byte, subscribe chan Event) {
	s.Lock()
	defer s.Unlock()

	data = make(map[string][]byte)

	for key:=range s.disk.KeysPrefix(prefix,nil){
		val,err:=s.disk.Read(key)
		if err==nil{
			data[key] = val
		}
	}

	subscribe=make(chan Event,EventChanSize)
	s.subs[subscribe] = prefix

	return data,subscribe
}

func (s *Storage) Unsubscribe(subscribe chan Event) {
	s.Lock()
	defer s.Unlock()

	delete(s.subs, subscribe)
	close(subscribe)
}