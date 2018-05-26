package commons

import "time"

const (
	ROLE_Pilot = "Pilot"
	ROLE_Engi  = "Engi"
	ROLE_Navi  = "Navi"
)

const (
	STATE_login = "login"
	STATE_cosmo = "cosmo"
	STATE_warp  = "star"
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
)

var (
	StartDateTime = time.Date(2018, 01, 01, 01, 01, 01, 01, time.Local)
)
