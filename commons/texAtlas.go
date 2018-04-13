package commons

import (
	"bytes"
	"encoding/json"
	"github.com/Shnifer/magellan/graph"
	"github.com/hajimehoshi/ebiten"
	"golang.org/x/image/colornames"
	"image/color"
	"io/ioutil"
)

type TexAtlasRec struct {
	FileName string
	Sx, Sy   int
	Color    color.RGBA
}
type TexAtlas map[string]TexAtlasRec

const atlasFN = "atlas.json"

var texPath string
var atlas TexAtlas

func InitTexAtlas(newTexPath string) {
	texPath = newTexPath
	saveAtlasExample(texPath + "example_" + atlasFN)
	data, err := ioutil.ReadFile(texPath + atlasFN)
	if err != nil {
		atlas = make(TexAtlas)
		panic("Can't find tex atlas file " + texPath + atlasFN)
	}
	err = json.Unmarshal(data, &atlas)
	if err != nil {
		panic(err)
	}
}

func GetAtlasTexColor(name string) (graph.Tex, color.RGBA) {
	rec, ok := atlas[name]
	if !ok {
		panic("GetAtlasTex: unknown name " + name)
	}

	tex, err := graph.GetTex(texPath+rec.FileName, ebiten.FilterLinear, rec.Sx, rec.Sy)
	if err != nil {
		panic(err)
	}
	return tex, rec.Color
}

func GetAtlasTex(name string) graph.Tex {
	tex, _ := GetAtlasTexColor(name)
	return tex
}

func saveAtlasExample(fn string) {
	exAtlas := make(map[string]TexAtlasRec)
	exAtlas["name"] = TexAtlasRec{
		FileName: "filename.png",
		Sx:       0,
		Sy:       0,
		Color:    colornames.White,
	}
	buf, err := json.Marshal(exAtlas)
	if err != nil {
		panic(err)
	}
	identbuf := bytes.Buffer{}
	json.Indent(&identbuf, buf, "", "  ")
	err = ioutil.WriteFile(fn, identbuf.Bytes(), 0)
	if err != nil {
		panic("can't write texture atlas example " + err.Error())
	}
}
