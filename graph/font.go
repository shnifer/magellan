package graph

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"math"
)

type faceSign struct {
	filename string
	size     float64
}

var faceCache map[faceSign]font.Face

func init() {
	faceCache = make(map[faceSign]font.Face)
}

func newFace(b []byte, size float64) (font.Face, error) {
	f, err := truetype.Parse(b)
	if err != nil {
		return nil, err
	}
	tto := &truetype.Options{
		Size: size,
	}
	face := truetype.NewFace(f, tto)
	return face, nil
}

func GetFace(fileName string, size float64, loader func(filename string) ([]byte, error)) (font.Face, error) {
	size = math.Floor(size * globalScale)
	sign := faceSign{
		filename: fileName,
		size:     size,
	}
	if face, ok := faceCache[sign]; ok {
		return face, nil
	}
	data, err := loader(fileName)
	if err != nil {
		return nil, err
	}
	face, err := newFace(data, size)
	if err != nil {
		return nil, err
	}
	faceCache[sign] = face
	return face, nil
}
