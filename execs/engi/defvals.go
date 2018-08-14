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

	NormTemperature      float64
	PoolTemperatureK     float64
	EmiTemperatureK      float64
	HeatProdTemperatureK float64
	CO2StepK             float64
	TankRadiBase         float64
	TankRadiK            float64
	OutRadiK             float64
	InRadiK              float64
	MaxAirLose           float64
	AirLoseQuot          float64
	RepairQuot           float64
	NormPressure         float64
	MinHole              float64
	MaxHole              float64
	HoleAZK              float64
	BrakeChanceK         float64
	OverheatAZK          float64
	RadiAZK              float64
	HitsToCounter        float64
	RadiCockPitK         float64
	MedRadiLevel         float64
	HardGDmg             float64
	HardGDmgRepeats      int
	MediumGDmg           float64
	MediumGDmgRepeats    int

	MediLowCounterS       int
	MediMidTotalS         int
	MediMidNeededS        int
	MediHittedDropPeriodS int
	MediOpts              MediOpts

	AliceAddr  string
	AlicePath  string
	AlicePass  string
	AliceLogin string

	RequestHyBoostListAddr string
	RequestHyBoostUseAddr  string
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
