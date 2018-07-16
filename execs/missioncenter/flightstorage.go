package main

import (
	"github.com/Shnifer/magellan/storage"
	"github.com/peterbourgon/diskv"
)

var flightDisk *storage.Storage

func initFlightStorage() {
	diskOpts := diskv.Options{
		BasePath:     DEFVAL.FlightDiskPath,
		CacheSizeMax: 1024 * 1024,
	}
	flightDisk = storage.New(DEFVAL.NodeName, diskOpts, DEFVAL.DiskRefreshPeriod)
	if DEFVAL.GameExchPort != "" && DEFVAL.GameExchPeriodMs > 0 {
		storage.RunExchanger(flightDisk, DEFVAL.GameExchPort, DEFVAL.GameExchAddrs, DEFVAL.GameExchPeriodMs)
	}
}
