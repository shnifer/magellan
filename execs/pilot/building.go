package main

import . "github.com/Shnifer/magellan/commons"

func AddBeacon(msg string) {
	b1 := Building{
		Type:     BUILDING_BEACON,
		GalaxyID: Data.State.GalaxyID,
		Period:   5000 * KDev(10),
		Message:  msg,
	}
	requestNewBuilding(b1)
	//duplicated into warp on server side
}

func requestNewBuilding(b Building) {
	buf := string(b.Encode())
	Client.SendRequest(CMD_ADDBUILDREQ + buf)
}

func requestRemoveBuilding(fullKey string) {
	Client.SendRequest(CMD_DELBUILDREQ + fullKey)
}
