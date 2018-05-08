package main

import (
	"bytes"
	"encoding/json"
	"github.com/Shnifer/magellan/commons"
	"io/ioutil"
	"log"
)

type tDefVals struct {
	Port string
	//in ms, 0 - default
	Timeout    int
	Room       string
	Role       string
	FullScreen bool
	WinW, WinH int

	DoProf          bool
	CpuProfFileName string
	MemProfFileName string

	GravityConst     float64
	WarpGravityConst float64
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
	}
}

func init() {
	setDefDef()

	exfn := resPath + "example_defdata.json"
	exbuf, _ := json.Marshal(DEFVAL)
	identbuf := bytes.Buffer{}
	json.Indent(&identbuf, exbuf, "", "    ")
	if err := ioutil.WriteFile(exfn, identbuf.Bytes(), 0); err != nil {
		log.Println("can't even write ", exfn)
	}

	fn := resPath + "defdata.json"

	buf, err := ioutil.ReadFile(fn)
	if err != nil {
		log.Println("cant read ", fn, "using default")
		return
	}
	json.Unmarshal(buf, &DEFVAL)
}
