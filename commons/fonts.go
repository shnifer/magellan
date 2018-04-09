package commons

import (
	"github.com/Shnifer/magellan/graph"
	"golang.org/x/image/font"
)

const (
	Face_cap = "caption"
)

const fontPath = "res/fonts/"

var Fonts map[string]font.Face

func init() {
	Fonts = make(map[string]font.Face)
	face, err := graph.GetFace(fontPath+"phantom.ttf", 20)
	if err != nil {
		panic(err)
	}
	Fonts[Face_cap] = face
}
