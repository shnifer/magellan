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
	Port       string
	Room       string
	Role       string
	FullScreen bool
	WinW, WinH int

	DoProf bool
}

var DEFVAL tDefVals

func setDefDef() {
	DEFVAL = tDefVals{
		Port: "http://localhost:8000",
		Room: "room101",
		Role: commons.ROLE_Engi,
		WinW: 1024,
		WinH: 768,
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
