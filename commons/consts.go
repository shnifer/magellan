package commons

import "time"

const ShipSize = 20

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
	//also MAGIC_MARK_BUILDING_BEACON, MAGIC_MARK_BUILDING_BLACKBOX
	//
	//also MAGIC_DEFAULT_STAR, MAGIC_DEFAULT_HARDPLANET, MAGIC_DEFAULT_GASPLANET,
	// MAGIC_DEFAULT_ASTEROID, MAGIC_DEFAULT_WARP
	//
	//also MAGIC_MARK_STAR, MAGIC_DEFAULT_HARDPLANET, MAGIC_DEFAULT_GASPLANET,
	// MAGIC_MARK_ASTEROID, MAGIC_MARK_WARP
	ShipAN              = "MAGIC_ship"
	MARKtheEarthAN      = "MAGIC_MARK_theearth"
	MARKtheMagellanAN   = "MAGIC_MARK_themagellan"
	MARKGLOWAN          = "MAGIC_MARK_GLOW"
	MARKShipAN          = "MAGIC_MARK_ship"
	OtherShipAN         = "MAGIC_othership"
	MARKOtherShipAN     = "MAGIC_MARK_othership"
	PredictorAN         = "MAGIC_predictor"
	NaviMarkerAN        = "MAGIC_navimarker"
	TrailAN             = "MAGIC_trail"
	ScannerAN           = "MAGIC_scanner"
	ThrustArrowAN       = "MAGIC_thurstarrow"
	RulerHAN            = "MAGIC_rulerh"
	RulerVAN            = "MAGIC_rulerv"
	DefaultBackgroundAN = "MAGIC_DEFAULT_background"
	EngiBackgroundAN    = "MAGIC_engibackground"
	CompassAN           = "MAGIC_compass"
	Frame9AN            = "front9"
	ButtonAN            = "MAGIC_button"
	TextPanelAN         = "MAGIC_textpanel"
	WayPointAN          = "MAGIC_waypoint"
	WayArrowAN          = "MAGIC_wayarrow"
	WarpInnerAN         = "MAGIC_warpinner"
	WarpOuterAN         = "MAGIC_warpouter"
	WarpGreenAN         = "MAGIC_warpgreen"
)

var (
	StartDateTime = time.Date(2018, 01, 01, 01, 01, 01, 01, time.Local)
)
