package graph

import (
	"github.com/hajimehoshi/ebiten"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
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

func (t Tex) Size() (w, h int) {
	return t.sw, t.sh
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

func GetTex(filename string, smoothFilter bool, sw, sh int, count int,
	loader func(filename string) (io.Reader, error)) (Tex, error) {

	filter := ebiten.FilterDefault
	if smoothFilter {
		filter = ebiten.FilterLinear
	}

	buf, err := loader(filename)
	if err != nil {
		return Tex{}, err
	}

	img, _, err := image.Decode(buf)
	if err != nil {
		return Tex{}, err
	}
	img2, err := ebiten.NewImageFromImage(img, filter)
	if err != nil {
		return Tex{}, err
	}

	t := TexFromImage(img2, filter, sw, sh, count, filename)
	if err != nil {
		return Tex{}, err
	}

	return t, nil
}

func CheckTexCache(cacheName string) (Tex, bool) {
	t, ok := texCache[cacheName]
	return t, ok
}

func StoreTexCache(cacheName string, tex Tex) {
	texCache[cacheName] = tex
}

//cache it manually
func SlidingTex(source Tex) (result Tex) {
	result = source
	w, h := source.image.Size()
	addW := source.sw
	newImage, _ := ebiten.NewImage(w+addW, h, source.filter)
	op := &ebiten.DrawImageOptions{}
	newImage.DrawImage(source.image, op)
	rect := image.Rect(0, 0, addW, h)
	op.SourceRect = &rect
	op.GeoM.Translate(float64(w-1), 0)
	newImage.DrawImage(source.image, op)
	result.image, _ = ebiten.NewImageFromImage(newImage, source.filter)
	return result
}

//cache it manually
func RoundTex(source Tex) (result Tex) {
	result = source
	w, h := source.sw, source.sh
	mask := NewSprite(CircleTex(), NoCam)
	mask.SetPivot(TopLeft())
	mask.SetSize(float64(w), float64(h))

	newImage, _ := ebiten.NewImage(w, h, source.filter)
	newImage.DrawImage(source.image, &ebiten.DrawImageOptions{})

	img, op := mask.ImageOp()
	op.CompositeMode = ebiten.CompositeModeDestinationIn
	newImage.DrawImage(img, op)
	result.image, _ = ebiten.NewImageFromImage(newImage, source.filter)
	return result
}

func ClearCache() {
	texCache = make(map[string]Tex, len(texCache))
}
