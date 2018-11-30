package main

import (
	"github.com/peterbourgon/diskv"
	"github.com/shnifer/magellan/commons"
	. "github.com/shnifer/magellan/log"
	"github.com/shnifer/magellan/storage"
)

var flightDisk *storage.Storage
var buildingSubscribe chan storage.Event
var buildingData map[storage.ObjectKey]string

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

func getBuildings(galaxyID string) {
	if buildingSubscribe != nil {
		flightDisk.Unsubscribe(buildingSubscribe)
	}
	//for warp we acumulate ALL buidings and sort+sum mines
	if galaxyID == commons.WARP_Galaxy_ID {
		galaxyID = ""
	}
	buildingData, buildingSubscribe = flightDisk.SubscribeAndData(galaxyID)
}

func updateBuildings() {
	if buildingSubscribe == nil || Scene == nil {
		return
	}
loop:
	for {
		select {
		case event := <-buildingSubscribe:
			fk := event.Key.FullKey()
			building, err := commons.Building{}.Decode([]byte(event.Data))
			if err != nil {
				Log(LVL_ERROR, "updateBuildings can't decode data: ", event.Data, " ", err)
				continue
			}
			switch event.Type {
			case storage.Add:
				Scene.EventAddBuilding(building, fk)
			case storage.Remove:
				Scene.EventDelBuilding(building, fk)
			}
		default:
			break loop
		}
	}
}
