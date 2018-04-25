package commons

import (
	"bytes"
	"encoding/json"
	"github.com/Shnifer/magellan/graph"
	"github.com/hajimehoshi/ebiten"
	"io/ioutil"
)

type TexAtlasRec struct {
	FileName string
	Sx, Sy   int
	Count    int
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

func GetAtlasTex(name string) graph.Tex {
	rec, ok := atlas[name]
	if !ok {
		panic("GetAtlasTex: unknown name " + name)
	}

	tex, err := graph.GetTex(texPath+rec.FileName, ebiten.FilterLinear, rec.Sx, rec.Sy, rec.Count)
	if err != nil {
		panic(err)
	}
	return tex
}

func NewAtlasSprite(atlasName string, cam *graph.Camera, denyCamScale, denyCamAngle bool) *graph.Sprite {
	return graph.NewSprite(GetAtlasTex(atlasName), cam, denyCamScale, denyCamAngle)
}

func NewAtlasSpriteHUD(atlasName string) *graph.Sprite {
	return graph.NewSpriteHUD(GetAtlasTex(atlasName))
}

func saveAtlasExample(fn string) {
	exAtlas := make(map[string]TexAtlasRec)
	exAtlas["name"] = TexAtlasRec{
		FileName: "filename.png",
		Sx:       0,
		Sy:       0,
		Count:    1,
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
