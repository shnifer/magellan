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

const (
	//google-disney
	OWNER_1 = "gd"
	//pony-express
	OWNER_2 = "pre"
	//mars-stroy-trest
	OWNER_3 = "mst"
	//mitsibishi-autovaz
	OWNER_4 = "mat"
	//red-cross
	OWNER_5 = "kkg"
)

var CorpNames = [...]string{OWNER_1, OWNER_2, OWNER_3, OWNER_4, OWNER_5}

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
	//for mines and fishhouses
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
	if b.OwnerID != "" && b.OwnerID != OWNER_1 && b.OwnerID != OWNER_2 &&
		b.OwnerID != OWNER_3 && b.OwnerID != OWNER_4 && b.OwnerID != OWNER_5 {
		Log(LVL_ERROR, "RequestNewBuilding: strange new building OwnerID ", b.OwnerID)
	}
	buf := string(b.Encode())
	client.SendRequest(CMD_ADDBUILDREQ + buf)
}

func RequestRemoveBuilding(client *network.Client, fullKey string) {
	client.SendRequest(CMD_DELBUILDREQ + fullKey)
}

func ColorByOwner(owner string) color.Color {
	switch owner {
	case OWNER_1:
		return colornames.Lightcyan
	case OWNER_2:
		return colornames.Lightgoldenrodyellow
	case OWNER_3:
		return colornames.Darkolivegreen
	case OWNER_4:
		return colornames.Steelblue
	case OWNER_5:
		return colornames.Firebrick
	default:
		Log(LVL_ERROR, "ColorByOwner unknown owner:", owner)
		return colornames.White
	}
}
func CompanyNameByOwner(owner string) string {
	switch owner {
	case OWNER_1:
		return "Google Disney"
	case OWNER_2:
		return "Pony Roscosmos Express"
	case OWNER_3:
		return "MarsStroyTrest"
	case OWNER_4:
		return "Mitsubishi AutoVAZ Technology"
	case OWNER_5:
		return "Red Cross Genetics"
	default:
		Log(LVL_ERROR, "CompanyNameByOwner unknown owner:", owner)
		return "#COMPANYNAME"
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
