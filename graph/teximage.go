package graph

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

type Tex struct {
	image  *ebiten.Image
	filter ebiten.Filter
	//0 for solid image
	sw, sh     int
	cols, rows int
}

func newTex(filename string, filter ebiten.Filter, sw, sh int) (Tex, error) {
	img, _, err := ebitenutil.NewImageFromFile(filename, filter)
	if err != nil {
		return Tex{}, err
	}
	return TexFromImage(img, filter, sw, sh), nil
}

func TexFromImage(image *ebiten.Image, filter ebiten.Filter, sw, sh int) Tex {
	if image == nil {
		return Tex{}
	}
	w, h := image.Size()
	if sw == 0 || sh == 0 {
		return Tex{
			image:  image,
			sw:     w,
			sh:     h,
			cols:   1,
			rows:   1,
			filter: filter,
		}
	}
	return Tex{
		image:  image,
		sw:     sw,
		sh:     sh,
		cols:   w / sw,
		rows:   h / sh,
		filter: filter,
	}
}

var texCache map[string]Tex

func init() {
	texCache = make(map[string]Tex)
}

func GetTex(filename string, filter ebiten.Filter, sw, sh int) (Tex, error) {
	if Tex, ok := texCache[filename]; ok {
		return Tex, nil
	}
	t, err := newTex(filename, filter, sw, sh)
	if err != nil {
		return Tex{}, err
	}
	texCache[filename] = t
	return t, nil
}
