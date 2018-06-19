package main

import (
	. "github.com/Shnifer/magellan/commons"
	"log"
)

func AddBeacon(msg string) {
	sessionTime := Data.PilotData.SessionTime
	angle := Data.PilotData.Ship.Pos.Dir() / 360
	basePeriod := 5000 * KDev(10)

	N := int(sessionTime / basePeriod)
	log.Println("N", N)
	period := sessionTime / (angle + float64(N))

	b1 := Building{
		Type:     BUILDING_BEACON,
		GalaxyID: Data.State.GalaxyID,
		Period:   period,
		Message:  msg,
	}
	//duplicated into warp on server side
	RequestNewBuilding(Client, b1)
}
