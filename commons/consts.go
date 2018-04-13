package commons

import "time"

const (
	ROLE_Pilot = "Pilot"
	ROLE_Engi  = "Engi"
	ROLE_Navi  = "Navi"
	ROLE_Cargo = "Cargo"
)

const (
	STATE_login = "login"
	STATE_cosmo = "cosmo"
	STATE_warp  = "star"
)

const (
	START_Galaxy_ID = "solar"
)

const (
	CMD_STATECHANGRFAIL = "FailedStateChange"
)

var (
	StartDateTime = time.Date(2018, 01, 01, 01, 01, 01, 01, time.Local)
)
