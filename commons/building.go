package commons

import (
	"encoding/json"
	. "github.com/Shnifer/magellan/log"
	"github.com/Shnifer/magellan/network"
	"github.com/Shnifer/magellan/storage"
	"golang.org/x/image/colornames"
	"image/color"
)

const (
	BUILDING_BLACKBOX  = "BUILDING_BLACKBOX"
	BUILDING_MINE      = "BUILDING_MINE"
	BUILDING_BEACON    = "BUILDING_BEACON"
	BUILDING_FISHHOUSE = "BUILDING_FISHHOUSE"
)

type Building struct {
	FullKey string

	Type string
	//where is it
	GalaxyID string
	//for mines
	PlanetID string
	//beckon and boxes are auto placed on far reach of system
	//very slow and some random if there are many
	Period float64

	Message string
	//for mine
	OwnerID string
}

func (b Building) Encode() []byte {
	buf, err := json.Marshal(b)
	if err != nil {
		Log(LVL_ERROR, "can't marshal Building", err)
		return nil
	}
	return buf
}

func (Building) Decode(buf []byte) (b Building, err error) {
	err = json.Unmarshal(buf, &b)
	if err != nil {
		return Building{}, err
	}
	return b, nil
}

func EventToCommand(e storage.Event) string {
	buf, err := json.Marshal(e)
	if err != nil {
		Log(LVL_ERROR, "can't marshal event", err)
		return ""
	}
	return CMD_BUILDINGEVENT + string(buf)
}

func DecodeEvent(buf []byte) (t int, b Building, err error) {
	var e storage.Event
	err = json.Unmarshal(buf, &e)
	if err != nil {
		return 0, Building{}, err
	}
	if e.Data != "" {
		b, err = Building{}.Decode([]byte(e.Data))
		if err != nil {
			return 0, Building{}, err
		}
		b.FullKey = e.Key.FullKey()
	} else {
		Log(LVL_ERROR, "decode event empty data field")
	}

	return e.Type, b, nil
}

func RequestNewBuilding(client *network.Client, b Building) {
	buf := string(b.Encode())
	client.SendRequest(CMD_ADDBUILDREQ + buf)
}

func RequestRemoveBuilding(client *network.Client, fullKey string) {
	client.SendRequest(CMD_DELBUILDREQ + fullKey)
}

func ColorByOwner(owner string) color.Color {
	switch owner {
	case "corp1":
		return colornames.Red
	case "corp2":
		return colornames.Green
	case "corp3":
		return colornames.Blue
	default:
		Log(LVL_ERROR, "ColorByOwner unknown owner:", owner)
		return colornames.White
	}
}

func AddBeacon(Data TData, Client *network.Client, msg string) {
	sessionTime := Data.PilotData.SessionTime
	angle := Data.PilotData.Ship.Pos.Dir() / 360
	basePeriod := 5000 * KDev(10)

	N := int(sessionTime / basePeriod)
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

func AddMine(Data TData, Client *network.Client, planetID string, owner string) {

	b := Building{
		Type:     BUILDING_MINE,
		GalaxyID: Data.State.GalaxyID,
		PlanetID: planetID,
		OwnerID:  owner,
	}
	RequestNewBuilding(Client, b)
}

func AddFishHouse(Data TData, Client *network.Client, planetID string, owner string) {

	b := Building{
		Type:     BUILDING_FISHHOUSE,
		GalaxyID: Data.State.GalaxyID,
		PlanetID: planetID,
		OwnerID:  owner,
	}
	RequestNewBuilding(Client, b)
}