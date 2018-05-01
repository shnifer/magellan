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
	//size of single sprite in sheet
	//0 for solid image
	sw, sh     int
	cols, rows int
	//count of sprites in sprite sheet
	count int
	name  string
}

func newTex(filename string, filter ebiten.Filter, sw, sh int, count int) (Tex, error) {
	img, _, err := ebitenutil.NewImageFromFile(filename, filter)
	if err != nil {
		return Tex{}, err
	}
	return TexFromImage(img, filter, sw, sh, count, filename), nil
}

func TexFromImage(image *ebiten.Image, filter ebiten.Filter, sw, sh int, count int, name string) Tex {
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
			count:  1,
			filter: filter,
		}
	}
	cols, rows := w/sw, h/sh
	if count == 0 {
		count = cols * rows
	}
	if count > cols*rows {
		panic("TexFromImage: count>cols*rows")
	}

	return Tex{
		image:  image,
		sw:     sw,
		sh:     sh,
		cols:   cols,
		rows:   rows,
		count:  count,
		filter: filter,
		name:   name,
	}
}

var texCache map[string]Tex

func init() {
	texCache = make(map[string]Tex)
}

func GetTex(filename string, smoothFilter bool, sw, sh int, count int) (Tex, error) {
	var cacheKey string
	if smoothFilter {
		cacheKey = "0" + filename
	} else {
		cacheKey = "1" + filename
	}
	if Tex, ok := texCache[cacheKey]; ok {
		return Tex, nil
	}

	filter := ebiten.FilterDefault
	if smoothFilter {
		filter = ebiten.FilterLinear
	}

	t, err := newTex(filename, filter, sw, sh, count)
	if err != nil {
		return Tex{}, err
	}
	texCache[cacheKey] = t
	return t, nil
}
