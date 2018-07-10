package main

import (
	"github.com/Shnifer/magellan/storage"
	"github.com/peterbourgon/diskv"
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
	nameDisk = storage.New(DEFVAL.NodeName, diskOpts)
}

func getNames(galaxyID string) {
	if namesSubscribe != nil {
		nameDisk.Unsubscribe(namesSubscribe)
	}
	namesData, namesSubscribe = nameDisk.SubscribeAndData(galaxyID)
}
