package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"bytes"
)

const DefValPath = "res/pilot/"

type tDefVals struct {
	Port        string
	Room string
	FullScreen bool
	WinW, WinH int

}

var DEFVAL tDefVals

func setDefDef() {
	DEFVAL = tDefVals{
		Port: "http://localhost:8000",
		Room: "room101",
		WinW:1024,
		WinH:768,
	}
}

func init() {
	setDefDef()

	exfn := DefValPath + "example_defdata.json"
	exbuf, err := json.Marshal(DEFVAL)
	identbuf:=bytes.Buffer{}
	json.Indent(&identbuf, exbuf,"","    ")
	if err := ioutil.WriteFile(exfn, identbuf.Bytes(), 0); err != nil {
		log.Println("can't even write ", exfn)
	}

	fn := DefValPath + "defdata.json"

	buf, err := ioutil.ReadFile(fn)
	if err != nil {
		log.Println("cant read ", fn, "using default")
		return
	}
	json.Unmarshal(buf, &DEFVAL)
}
