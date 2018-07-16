package main

import (
	"bytes"
	"encoding/json"
	. "github.com/Shnifer/magellan/log"
	"io/ioutil"
)

const DefValPath = "./"

var roleName = "missioncenter"

type tDefVals struct {
	FullScreen     bool
	WinW, WinH     int
	HalfResolution bool
	LowQ           bool
	VSync bool

	DebugControl bool
	DoProf       bool

	//in ms
	LogTimeoutMs  int
	LogRetryMinMs int
	LogRetryMaxMs int
	LogIP         string
	LogHostName   string

	//node
	FlightDiskPath    string
	NamesDiskPath     string
	NodeName          string
	DiskRefreshPeriod int

	//storage exchanger
	GameExchPort     string
	GameExchAddrs    []string
	GameExchPeriodMs int
}

var DEFVAL tDefVals

func setDefDef() {
	DEFVAL = tDefVals{
		WinW:              1024,
		WinH:              768,
		LogTimeoutMs:      1000,
		LogRetryMinMs:     10,
		LogRetryMaxMs:     60000,
		DiskRefreshPeriod: 30,
		NodeName:          "MissionControl",
		LogHostName:       "DummyLogHost",
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
