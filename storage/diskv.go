package storage

import (
	"errors"
	"github.com/peterbourgon/diskv"
	"strconv"
	"strings"
	"sync"
	"time"
)

const generatorID = "!!ID"

var ErrAlreadyExist = errors.New("already exist")

//Diskv wrapper for immutable disk storage of string Key-value pairs
//with subscribe on addition of new keys
//and unique Id generator
type disk struct {
	*diskv.Diskv

	sync.RWMutex
	curID   int
	keys    map[string]struct{}
	keySubs map[chan string]struct{}
}

func newDisk(diskOpts diskv.Options, refreshFilesPeriod int) *disk {
	diskV := diskv.New(diskOpts)
	curIDs, err := diskV.Read(generatorID)
	if err != nil {
		curIDs = []byte("1")
		err := diskV.Write(generatorID, curIDs)
		if err != nil {
			panic(err)
		}
	}
	id, err := strconv.Atoi(string(curIDs))
	if err != nil {
		curIDs = []byte("1")
		id = 1
		err := diskV.Write(generatorID, curIDs)
		if err != nil {
			panic(err)
		}
	}

	keys := make(map[string]struct{}, 0)

	//preload in cache
	for key := range diskV.Keys(nil) {
		_ = diskV.ReadString(key)
		if key == generatorID {
			continue
		}
		keys[key] = struct{}{}
	}

	res:= &disk{
		curID:   id,
		Diskv:   diskV,
		keySubs: make(map[chan string]struct{}, 0),
		keys:    keys,
	}

	if refreshFilesPeriod>0 {
		go daemonFilesChecker(res, refreshFilesPeriod)
	}

	return res
}

//use this to add new pairs.
//announce new key for subscribers
func (d *disk) append(key, val string) error {
	d.Lock()
	defer d.Unlock()

	if d.has(key) {
		return ErrAlreadyExist
	}

	err := d.Write(key, []byte(val))
	if err != nil {
		return err
	}

	d.registerNewKey(key)

	return nil
}

//get list of all keys in storage and subscribe for new
func (d *disk) subscribe() (fullKeys map[string]struct{}, subscribe chan string) {
	d.Lock()
	defer d.Unlock()

	fullKeys = make(map[string]struct{}, 0)

	for key := range d.keys {
		fullKeys[key] = struct{}{}
	}

	subscribe = make(chan string, EventChanSize)
	d.keySubs[subscribe] = struct{}{}

	return fullKeys, subscribe
}

func (d *disk) unsubscribe(subscribe chan string) {
	d.Lock()
	defer d.Unlock()

	delete(d.keySubs, subscribe)
	close(subscribe)
}

func (d *disk) nextID() int {
	d.Lock()
	defer d.Unlock()

	d.curID++
	str := strconv.Itoa(d.curID)
	d.WriteStream(generatorID, strings.NewReader(str), true)
	return d.curID
}

func (d *disk) getKeys() []string {
	d.Lock()
	defer d.Unlock()
	res := make([]string, 0, len(d.keys))
	for key := range d.keys {
		res = append(res, key)
	}
	return res
}

func (d *disk) has(key string) bool {
	d.RLock()
	defer d.RUnlock()
	_, exist := d.keys[key]
	return exist
}

//run under mutex
func (d *disk) registerNewKey(key string) {
	for ch := range d.keySubs {
		ch <- key
	}
	d.keys[key] = struct{}{}
}

func daemonFilesChecker(d *disk, period int) {
	var keys []string
	for {
		time.Sleep(time.Duration(period)*time.Second)
		for key:=range d.Keys(nil){
			keys = append(keys, key)
		}

		d.Lock()
		for _,key:=range keys{
			if key==generatorID{
				continue
			}
			if _,exist:=d.keys[key]; !exist{
				d.registerNewKey(key)
			}
		}
		d.Unlock()
	}
}