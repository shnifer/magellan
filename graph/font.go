package graph

import (
	"golang.org/x/image/font"
	"io/ioutil"
	"github.com/golang/freetype/truetype"

)

type faceSign struct{
	filename string
	size float64
}

var faceCache map[faceSign]font.Face

func init(){
	faceCache = make(map[faceSign]font.Face)
}

func newFace(fileName string, size float64) (font.Face, error) {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	f, err := truetype.Parse(b)
	if err != nil {
		return nil, err
	}
	tto:=&truetype.Options{
		Size: size,
	}
	face := truetype.NewFace(f, tto)
	return face,nil
}

func GetFace(fileName string, size float64) (font.Face, error) {
	sign:=faceSign{
		filename:fileName,
		size:size,
	}
	if face, ok:=faceCache[sign]; ok{
		return face,nil
	}
	face,err:=newFace(fileName,size)
	if err!=nil{
		return nil,err
	}
	faceCache[sign] = face
	return face,nil
}