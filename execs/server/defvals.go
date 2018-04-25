package main

import (
	"bytes"
	"encoding/json"
	. "github.com/Shnifer/magellan/commons"
	"io/ioutil"
	"log"
)

const DefValPath = "res/server/"

type tDefVals struct {
	Port        string
	NeededRoles []string

	StartWarpSpeed         float64
	SolarStartLocationName string
}

var DEFVAL tDefVals

func setDefDef() {
	DEFVAL = tDefVals{
		Port: "8000",
		NeededRoles: []string{
			ROLE_Pilot,
			ROLE_Navi,
			ROLE_Engi,
			ROLE_Cargo},
		SolarStartLocationName: "magellan",
	}
}

func init() {
	setDefDef()

	exfn := DefValPath + "example_defdata.json"
	exbuf, err := json.Marshal(DEFVAL)
	identbuf := bytes.Buffer{}
	json.Indent(&identbuf, exbuf, "", "    ")
	if err := ioutil.WriteFile(exfn, identbuf.Bytes(), 0); err != nil {
		log.Println("can't even write ", exfn)
	}

	fn := DefValPath + "defdata.json"

	buf, err := ioutil.ReadFile(DefValPath + "defdata.json")
	if err != nil {
		log.Println("cant read ", fn, "using default")
		return
	}
	json.Unmarshal(buf, &DEFVAL)
}
