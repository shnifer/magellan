package main

import (
	. "github.com/Shnifer/magellan/log"
	"github.com/Shnifer/magellan/storage"
	"github.com/peterbourgon/diskv"
	"github.com/pkg/errors"
	"strings"
)

const namesPath = "names"

var nameDisk *storage.Storage
var namesSubscribe chan storage.Event
var namesData map[storage.ObjectKey]string

func initNamesStorage() {
	diskOpts := diskv.Options{
		BasePath:     namesPath,
		CacheSizeMax: 1024 * 1024,
	}
	nameDisk = storage.New(DEFVAL.NodeName, diskOpts, DEFVAL.DiskRefreshPeriod)
}

func getNames(galaxyID string) {
	if namesSubscribe != nil {
		nameDisk.Unsubscribe(namesSubscribe)
	}
	namesData, namesSubscribe = nameDisk.SubscribeAndData(galaxyID)
}

type nameRec struct {
	planetID string
	name     string
}

func (r nameRec) encode() string {
	return r.planetID + "!" + r.name
}

func (nameRec) decode(str string) (nameRec, error) {
	a := strings.SplitN(str, "!", 2)
	if len(a) < 2 {
		return nameRec{}, errors.New("nameRec.decode can't split " + str)
	}
	return nameRec{
		planetID: a[0],
		name:     a[1],
	}, nil
}

func updateNamesAndNotes() {
	if namesSubscribe == nil || Scene == nil {
		return
	}
loop:
	for {
		select {
		case event := <-namesSubscribe:
			fk := event.Key.FullKey()
			switch event.Type {
			case storage.Add:
				nameRec, err := nameRec{}.decode(event.Data)
				if err != nil {
					Log(LVL_ERROR, "updateNamesAndNotes can't decode data: ", event.Data, " ", err)
					continue
				}
				Scene.EventAddName(nameRec, fk)
			case storage.Remove:
				Scene.EventDelName(fk)
			}
		default:
			break loop
		}
	}
}
