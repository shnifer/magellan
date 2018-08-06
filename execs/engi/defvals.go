package main

import (
	"bytes"
	"encoding/json"
	"github.com/Shnifer/magellan/commons"
	"io/ioutil"
	"log"
)

const DefValPath = "./"

type tDefVals struct {
	Port           string
	Timeout        int
	PingPeriod     int
	Room           string
	Role           string
	FullScreen     bool
	WinW, WinH     int
	HalfResolution bool
	VSync          bool

	DoProf    bool
	DebugPort string

	//in ms
	LogLogTimeoutMs int
	LogRetryMinMs   int
	LogRetryMaxMs   int
	LogIP           string
	LogHostName     string

	SpriteSizeW,
	SpriteSizeH int

	RanmaAddr         string
	DropOnRepair      bool
	RanmaTimeoutMs    int
	RanmaHistoryDepth int

	RanmaMaxDegradePercent float64
	EmissionDegradePercent float64

	GravityConst float64
}

var DEFVAL tDefVals

func setDefDef() {
	DEFVAL = tDefVals{
		Port:            "http://localhost:8000",
		Timeout:         100,
		PingPeriod:      100,
		Room:            "room101",
		Role:            commons.ROLE_Engi,
		WinW:            1024,
		WinH:            768,
		LogLogTimeoutMs: 1000,
		LogRetryMinMs:   10,
		LogRetryMaxMs:   60000,
	}
}

func init() {
	setDefDef()

	exfn := DefValPath + "example_ini_" + roleName + ".json"
	exbuf, err := json.Marshal(DEFVAL)
	identbuf := bytes.Buffer{}
	json.Indent(&identbuf, exbuf, "", "    ")
	if err := ioutil.WriteFile(exfn, identbuf.Bytes(), 0); err != nil {
		log.Println("can't even write ", exfn)
	}

	fn := DefValPath + "ini_" + roleName + ".json"

	buf, err := ioutil.ReadFile(fn)
	if err != nil {
		log.Println("cant read ", fn, "using default")
		return
	}
	json.Unmarshal(buf, &DEFVAL)
}
