package main

import (
	"github.com/Shnifer/magellan/graph"
	"golang.org/x/image/font"
)

const (
	face_cap = "caption"
)

var fonts map[string]font.Face

func init() {
	fonts = make(map[string]font.Face)
	face, err := graph.GetFace(fontPath+"phantom.ttf", 20)
	if err != nil {
		panic(err)
	}
	fonts[face_cap] = face
}
