package draw

import (
	"github.com/Shnifer/magellan/graph"
	. "github.com/Shnifer/magellan/log"
	"github.com/Shnifer/magellan/static"
	"golang.org/x/image/font"
)

const (
	Face_cap   = "caption"
	Face_stats = "stats"
	Face_list  = "list"
)

var Fonts map[string]font.Face

func fontLoader(filename string) ([]byte, error) {
	return static.Load("fonts", filename)
}

func InitFonts() {
	Fonts = make(map[string]font.Face)

	face, err := graph.GetFace("phantom.ttf", 20, fontLoader)
	if err != nil {
		Log(LVL_ERROR, err)
	}
	Fonts[Face_cap] = face

	face, err = graph.GetFace("interdim.ttf", 16, fontLoader)
	if err != nil {
		Log(LVL_ERROR, err)
	}
	Fonts[Face_stats] = face

	face, err = graph.GetFace("phantom.ttf", 12, fontLoader)
	if err != nil {
		Log(LVL_ERROR, err)
	}
	Fonts[Face_list] = face
}
