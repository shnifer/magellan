package main

import (
	"github.com/peterbourgon/diskv"
	"github.com/Shnifer/magellan/storage"
)

const storagePath = "xstore"

var disk *storage.Storage

func initStorage(){
	diskOpts := diskv.Options{
		BasePath:     storagePath,
		CacheSizeMax: 1024 * 1024,
	}
	disk := storage.New(DEFVAL.NodeName, diskOpts)
	if DEFVAL.GameExchPort != "" && DEFVAL.GameExchPeriodMs > 0 {
		storage.RunExchanger(disk, DEFVAL.GameExchPort, DEFVAL.GameExchAddrs, DEFVAL.GameExchPeriodMs)
	}
}