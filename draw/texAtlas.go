package draw

import (
	"bytes"
	"encoding/json"
	"github.com/Shnifer/magellan/graph"
	"github.com/Shnifer/magellan/static"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"strconv"
)

type TexAtlasRec struct {
	FileName string
	Sx, Sy   int
	Count    int
	Smooth   bool
}
type TexAtlas map[string]TexAtlasRec

const atlasFN = "atlas.json"

var atlas TexAtlas

func InitTexAtlas() {
	saveAtlasExample("example_" + atlasFN)
	data, err := static.Load("textures", atlasFN)
	if err != nil {
		panic("Can't find tex atlas file " + atlasFN)
	}
	atlas = make(TexAtlas)
	err = json.Unmarshal(data, &atlas)
	if err != nil {
		panic(err)
	}
}

func atlasLoader(filename string) (io.Reader, error) {
	return static.Read("textures", filename)
}

//return tex and error
func getAtlasTex(name string) (graph.Tex, error) {
	rec, ok := atlas[name]
	if !ok {
		return graph.Tex{}, errors.New("Not found atlas")
	}

	tex, err := graph.GetTex(rec.FileName, rec.Smooth, rec.Sx, rec.Sy, rec.Count, atlasLoader)
	if err != nil {
		return graph.Tex{}, err
	}
	return tex, nil
}

//return tex, panic on error
func GetAtlasTex(name string) graph.Tex {
	tex, err := getAtlasTex(name)
	if err != nil {
		panic(err)
	}
	return tex
}

func NewAtlasSprite(atlasName string, params graph.CamParams) *graph.Sprite {
	return graph.NewSprite(GetAtlasTex(atlasName), params)
}

func NewAtlasSpriteHUD(atlasName string) *graph.Sprite {
	return graph.NewSpriteHUD(GetAtlasTex(atlasName))
}

func NewAtlasFrame9HUD(atlasName string, w, h int) *graph.Frame9HUD {
	var sprites [9]*graph.Sprite
	for i := 0; i < 9; i++ {
		tex, err := getAtlasTex(atlasName + strconv.Itoa(i))
		if err != nil {
			continue
		}
		sprites[i] = graph.NewSpriteHUD(tex)

	}
	return graph.NewFrame9(sprites, float64(w), float64(h))
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
