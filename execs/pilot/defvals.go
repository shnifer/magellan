package main

import (
	"bytes"
	"encoding/json"
	"github.com/Shnifer/magellan/commons"
	. "github.com/Shnifer/magellan/log"
	"io/ioutil"
)

const DefValPath = "./"

type tDefVals struct {
	Port string
	//in ms, 0 - default
	Timeout    int
	PingPeriod int
	Room       string
	Role       string

	FullScreen     bool
	WinW, WinH     int
	HalfResolution bool
	LowQ           bool
	VSync          bool

	DebugControl bool
	DoProf       bool
	DebugPort    string

	GravityConst      float64
	WarpGravityConst  float64
	DT                float64
	MinDistortion     float64
	VelDistWarpK      float64
	WarpGravThreshold float64
	WarpGravPowN      float64

	//predictors
	CosmoPredictorUpdT     float64
	CosmoPredictorNumInSec int
	CosmoPredictorGravEach int
	CosmoPredictorTrackLen int
	CosmoPredictorDrawMaxP int
	WarpPredictorUpdT      float64
	WarpPredictorNumInSec  int
	WarpPredictorTrackLen  int
	WarpPredictorDrawMaxP  int

	//in ms
	LogTimeoutMs  int
	LogRetryMinMs int
	LogRetryMaxMs int
	LogIP         string
	LogHostName   string

	//inms
	OtherShipElastic int

	DefWarpCamScale  float64
	DefCosmoCamScale float64
}

var DEFVAL tDefVals

func setDefDef() {
	DEFVAL = tDefVals{
		Port:             "http://localhost:8000",
		Room:             "room101",
		Role:             commons.ROLE_Pilot,
		WinW:             1024,
		WinH:             768,
		GravityConst:     100,
		WarpGravityConst: 100,
		WarpGravPowN:     3,
		DT:               0.001,
		LogTimeoutMs:     1000,
		LogRetryMinMs:    10,
		LogRetryMaxMs:    60000,
		OtherShipElastic: 400,
	}
}

func init() {
	setDefDef()

	exfn := DefValPath + "example_ini_" + roleName + ".json"
	exbuf, err := json.Marshal(DEFVAL)
	identbuf := bytes.Buffer{}
	json.Indent(&identbuf, exbuf, "", "    ")
	if err := ioutil.WriteFile(exfn, identbuf.Bytes(), 0); err != nil {
		Log(LVL_WARN, "can't even write ", exfn)
	}

	fn := DefValPath + "ini_" + roleName + ".json"

	buf, err := ioutil.ReadFile(fn)
	if err != nil {
		Log(LVL_WARN, "cant read ", fn, "using default")
		return
	}
	json.Unmarshal(buf, &DEFVAL)
}
