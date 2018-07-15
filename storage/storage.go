package storage

import (
	. "github.com/Shnifer/magellan/log"
	"github.com/peterbourgon/diskv"
	"sync"
)

const EventChanSize = 16
const glyphDel = "!"

const (
	Add = iota
	Remove
)

type Event struct {
	Type int
	Key  ObjectKey
	Data string
}

type Storage struct {
	sync.Mutex

	node string
	disk *disk

	newKeys chan string
	//maps[subscribe chan]Area
	subs map[chan Event]string
}

func New(nodeName string, diskOpts diskv.Options, refreshFilesPeriod int) *Storage {
	LogFunc("storage.New")

	disk := newDisk(diskOpts, refreshFilesPeriod)
	_, keySub := disk.subscribe()

	s := &Storage{
		disk:    disk,
		node:    nodeName,
		subs:    make(map[chan Event]string),
		newKeys: keySub,
	}

	go subscribeLoop(s)

	return s
}

func (s *Storage) Add(area, key string, val string) error {
	LogFunc("storage.Add area: " + area + " key: " + key + " val: " + val)

	s.Lock()
	defer s.Unlock()

	objectKey := newKey(area, s.node, key)
	fk := objectKey.FullKey()
	err := s.disk.append(fk, val)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) Remove(objectKey ObjectKey) error {
	LogFunc("storage.Remove fullKey:" + objectKey.FullKey())

	s.Lock()
	defer s.Unlock()

	fk := objectKey.FullKey()
	err := s.disk.append(glyphDel+fk, "")
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) NextID() int {
	LogFunc("storage.NextID")
	return s.disk.nextID()
}

func (s *Storage) SubscribeAndData(area string) (data map[ObjectKey]string, subscribe chan Event) {
	LogFunc("storage.SubscribeAndData area: " + area)

	s.Lock()
	defer s.Unlock()

	data = make(map[ObjectKey]string)

	for key := range s.disk.KeysPrefix(area, nil) {
		//deleted item
		if s.disk.has(glyphDel + key) {
			continue
		}

		val, err := s.disk.Read(key)
		if err == nil {
			objKey, err := ReadKey(key)
			if err == nil {
				data[objKey] = string(val)
			}
		}
	}

	subscribe = make(chan Event, EventChanSize)
	s.subs[subscribe] = area

	return data, subscribe
}

func (s *Storage) Unsubscribe(subscribe chan Event) {
	LogFunc("storage.Unsubscribe")

	s.Lock()
	defer s.Unlock()

	delete(s.subs, subscribe)
	close(subscribe)
}

func subscribeLoop(s *Storage) {
	for {
		newKey := <-s.newKeys

		s.procNewKey(newKey)
	}
}

func (s *Storage) procNewKey(newKey string) {
	LogFunc("storage.procNewKey: " + newKey)

	objKey, err := ReadKey(newKey)
	if err != nil {
		Log(LVL_ERROR, "procNewKey: can't read objKey ", newKey)
		return
	}

	s.Lock()
	defer s.Unlock()

	switch objKey.glyph {
	//object added
	case "":
		if !s.disk.has(glyphDel + newKey) {
			val := s.disk.ReadString(newKey)
			s.sendEvent(Add, objKey, val)
		}
	//object removed
	case glyphDel:
		//send event about the object, not deleteKey
		objKey.glyph = ""
		fk := objKey.FullKey()
		if s.disk.has(fk) {
			val := s.disk.ReadString(fk)
			s.sendEvent(Remove, objKey, val)
		}
	}

}

func (s *Storage) sendEvent(t int, objKey ObjectKey, val string) {
	Event := Event{
		Type: t,
		Key:  objKey,
		Data: val,
	}

	area := objKey.Area

	for ch, subArea := range s.subs {
		if subArea == area {
			ch <- Event
		}
	}
}
