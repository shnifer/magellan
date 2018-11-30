package main

import (
	"bytes"
	"encoding/json"
	"github.com/shnifer/magellan/commons"
	"golang.org/x/image/colornames"
	"log"
)

func main() {
	gp := commons.GalaxyPoint{
		InnerColor: colornames.Aliceblue,
		Color:      colornames.Orange,
	}
	dat, _ := json.Marshal(gp)

	dat = bytes.Replace(dat, []byte(`{"R":0,"G":0,"B":0,"A":0}`), []byte("{}"), -1)

	log.Println(string(dat))
	var rest commons.GalaxyPoint
	json.Unmarshal(dat, &rest)
	log.Println(rest)
}
