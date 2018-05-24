package storage

import (
	"errors"
	"github.com/peterbourgon/diskv"
	"strconv"
	"strings"
	"sync"
)

const generatorID = "!!ID"
const EventChanSize = 16
const separator = "~"
const glyphDel = "!"

var ErrAlreadyExist = errors.New("already exist")

const (
	Add = iota
	Remove
)

type Event struct {
	Type    int
	Key     string
	FullKey string
	Data    string
}

type Storage struct {
	sync.Mutex
	curID int

	node string

	disk *diskv.Diskv
	//maps[subscribe chan]area
	subs map[chan Event]string
}

func New(nodeName string, diskOpts diskv.Options) *Storage {
	disk := diskv.New(diskOpts)

	curIDs, err := disk.Read(generatorID)
	if err != nil {
		curIDs = []byte("1")
		err := disk.Write(generatorID, curIDs)
		if err != nil {
			panic(err)
		}
	}
	id, err := strconv.Atoi(string(curIDs))
	if err != nil {
		panic(err)
	}

	return &Storage{
		curID: id,
		disk:  disk,
		node:  nodeName,
	}
}

func (s *Storage) Add(area, key string, val string) error {
	s.Lock()
	defer s.Unlock()

	fk := newKey("", area, s.node, key)
	if s.disk.Has(fk) {
		return ErrAlreadyExist
	}

	if err := s.disk.Write(fk, []byte(val)); err != nil {
		return err
	}

	if !s.disk.Has("!" + fk) {
		s.event(Add, key, fk, val)
	}

	return nil
}

func (s *Storage) Remove(fullKey string) error {
	s.Lock()
	defer s.Unlock()

	if s.disk.Has(glyphDel + fullKey) {
		return ErrAlreadyExist
	}
	_, _, _, key := mustSplitKey(fullKey)

	s.disk.Write(glyphDel+fullKey, []byte{})
	if s.disk.Has(fullKey) {
		s.event(Remove, key, fullKey, "")
	}
	return nil
}

func (s *Storage) Subscribe(area string) (data map[string]string, subscribe chan Event) {
	s.Lock()
	defer s.Unlock()

	data = make(map[string]string)

	for key := range s.disk.KeysPrefix(area, nil) {
		//deleted item
		if s.disk.Has(glyphDel + key) {
			continue
		}

		val, err := s.disk.Read(key)
		if err == nil {
			data[key] = string(val)
		}
	}

	subscribe = make(chan Event, EventChanSize)
	s.subs[subscribe] = area

	return data, subscribe
}

func (s *Storage) Unsubscribe(subscribe chan Event) {
	s.Lock()
	defer s.Unlock()

	delete(s.subs, subscribe)
	close(subscribe)
}

func (s *Storage) NextID() int {
	s.Lock()
	defer s.Unlock()

	s.curID++
	str := strconv.Itoa(s.curID)
	s.disk.WriteStream(generatorID, strings.NewReader(str), true)
	return s.curID
}

func (s *Storage) event(t int, key, fullKey string, val string) {
	Event := Event{
		Type:    t,
		Key:     key,
		FullKey: fullKey,
		Data:    val,
	}

	_, area, _, _ := mustSplitKey(fullKey)

	for ch, subArea := range s.subs {
		if subArea == area {
			ch <- Event
		}
	}
}

// strings must not have "/","\" or "~" symbols
func newKey(glyph, area, node, key string) string {
	return glyph + strings.Join([]string{area, node, key}, separator)
}

func mustSplitKey(fullKey string) (glyph, area, node, key string) {
	var err error
	glyph, area, node, key, err = splitKey(fullKey)
	if err != nil {
		panic(err)
	}
	return area, node, key, glyph
}

func isCorrectKey(fullKey string) bool {
	if fullKey == generatorID {
		return true
	}
	_, _, _, _, err := splitKey(fullKey)
	return err == nil
}

func splitKey(fullKey string) (glyph, area, node, key string, err error) {
	parts := strings.Split(fullKey, separator)
	if len(parts) != 3 {
		return "", "", "", "", errors.New("splitKey: len(parts)!=3")
	}
	area, node, key = parts[0], parts[1], parts[2]
	if area == "" || node == "" || key == "" {
		return "", "", "", "", errors.New("splitKey: some fields are empty " + fullKey)
	}
	switch area[:1] {
	case glyphDel:
		glyph = glyphDel
		area = area[1:]
	}
	return glyph, area, node, key, nil
}
