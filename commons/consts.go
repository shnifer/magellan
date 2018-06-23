package commons

import "time"

const (
	ROLE_Pilot  = "Pilot"
	ROLE_Engi   = "Engi"
	ROLE_Navi   = "Navi"
	ROLE_Server = "Server"
)

const (
	STATE_login = "login"
	STATE_cosmo = "cosmo"
	STATE_warp  = "warp"
)

const (
	START_Galaxy_ID = "solar"
	WARP_Galaxy_ID  = "warp"
)

const (
	CMD_STATECHANGEFAIL = "FailedStateChange"
	CMD_BUILDINGEVENT   = "BuildEvent"        //from server to subscribed clients
	CMD_ADDBUILDREQ     = "AddBuildingReq"    //from client to server
	CMD_DELBUILDREQ     = "DeleteBuildingReq" //from client to server
	CMD_LOGGAMEEVENT    = "LogGameEvent"      //from client to server
)

const (
	//also MAGIC_BUILDING_MINE, MAGIC_BUILDING_MINE
	//
	//also MAGIC_DEFAULT_STAR, MAGIC_DEFAULT_PLANET, MAGIC_DEFAULT_ASTEROID,
	//also MAGIC_DEFAULT_SATELLITE, MAGIC_DEFAULT_WARP
	ShipAN ="MAGIC_ship"
	OtherShipAN = "MAGIC_othership"
	PredictorAN = "MAGIC_predictor"
	NaviMarkerAN = "MAGIC_navimarker"
	TrailAN = "MAGIC_trail"
    ScannerAN  = "MAGIC_scanner"
    ThrustArrowAN = "MAGIC_thurstarrow"
	RulerHAN = "MAGIC_rulerh"
	RulerVAN = "MAGIC_rulerv"
	DefaultBackgroundAN = "MAGIC_DEFAULT_background"
	EngiBackgroundAN = "MAGIC_engibackground"
	CompassAN = "MAGIC_compass"
	Frame9AN = "front9"
)

var (
	StartDateTime = time.Date(2018, 01, 01, 01, 01, 01, 01, time.Local)
)
