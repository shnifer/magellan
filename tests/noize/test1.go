package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/ojrac/opensimplex-go"
	"golang.org/x/image/colornames"
	"image/color"
	"image/png"
	"os"
	"time"
)

const winSize = 400

func noise(x, y, t, frequency, lacunarity, gain float64, octaves int) float64 {
	const bigDelta = 100
	var sum float64
	amplitude := 1.0
	for i := 0; i < octaves; i++ {
		sum += generator.Eval3(x*frequency+float64(i)*bigDelta, y*frequency, t) * amplitude
		frequency *= lacunarity
		amplitude *= gain
	}
	return sum
}

func lerp(color1, color2 color.RGBA, x, min, max float64) color.RGBA {
	if max <= min {
		return color1
	}
	k1 := (x - min) / (max - min)
	k2 := 1 - k1
	var res color.RGBA
	res.R = uint8(float64(color1.R)*k2 + float64(color2.R)*k1)
	res.G = uint8(float64(color1.G)*k2 + float64(color2.G)*k1)
	res.B = uint8(float64(color1.B)*k2 + float64(color2.B)*k1)
	res.A = uint8(float64(color1.A)*k2 + float64(color2.A)*k1)
	return res
}

var generator *opensimplex.Noise

var frequency, lacunarity, gain float64
var octaves int
var level1, level2 float64

var img *ebiten.Image

var startT time.Time

func calcPixels(frequency, lacunarity, gain float64, octaves int, t float64) []byte {
	f := make([]float64, winSize*winSize)
	var fMin, fMax float64
	for y := 0; y < winSize; y++ {
		for x := 0; x < winSize; x++ {
			fx, fy := float64(x-winSize/2), float64(y-winSize/2)
			v := noise(fx, fy, t, frequency, lacunarity, gain, octaves)
			if v < fMin {
				fMin = v
			}
			if v > fMax {
				fMax = v
			}
			f[x+y*winSize] = v
		}
	}

	for i, v := range f {
		f[i] = (v - fMin) / (fMax - fMin)
	}

	color1 := colornames.Deepskyblue
	color2 := colornames.Forestgreen

	p := make([]byte, 4*winSize*winSize)
	for y := 0; y < winSize; y++ {
		for x := 0; x < winSize; x++ {
			f := f[x+y*winSize]

			var clr color.RGBA
			switch {
			case f <= level1:
				clr = lerp(colornames.Black, color1, f, 0, level1)
			case f <= level2:
				clr = lerp(color1, color2, f, level1, level2)
			case f <= 1:
				clr = lerp(color2, colornames.White, f, level2, 1)
			}
			p[4*x+4*winSize*y+0] = clr.R
			p[4*x+4*winSize*y+1] = clr.G
			p[4*x+4*winSize*y+2] = clr.B
			p[4*x+4*winSize*y+3] = 255
		}
	}
	return p
}

func saveSpriteSheet() {
	const rows = 6
	const cols = 6

	f, err := os.Create("noisesheet.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	sheet, _ := ebiten.NewImage(cols*winSize, rows*winSize, ebiten.FilterDefault)

	all := make([]byte, cols*rows*winSize*winSize*4)
	for col := 0; col < cols; col++ {
		for row := 0; row < rows; row++ {
			t := float64(col+row*cols)/20
			p := calcPixels(frequency, lacunarity, gain, octaves, t)
			for y := 0; y < winSize; y++ {
				for x := 0; x < winSize; x++ {
					ax := col*winSize + x
					ay := row*winSize + y
					for i := 0; i < 4; i++ {
						all[i+ax*4+ay*(cols*winSize*4)] = p[i+x*4+y*winSize*4]
					}
				}
			}
		}
	}
	sheet.ReplacePixels(all)

	png.Encode(f, sheet)
}

func update(window *ebiten.Image) error {
	procInput()

	if ebiten.IsRunningSlowly() {
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		saveSpriteSheet()
	}

	t := time.Since(startT).Seconds()

	p := calcPixels(frequency, lacunarity, gain, octaves, t)

	err := img.ReplacePixels(p)
	if err != nil {
		panic(err)
	}

	window.DrawImage(img, &ebiten.DrawImageOptions{})

	ebitenutil.DebugPrint(window, fmt.Sprintf(
		"FPS: %v\n[A-D] frequency: %.2f\n[S-W] lacunarity: %.2f\n[Q-E] gain: %.2f\n"+
			"[1-2] octaves: %v\n[Z-X] level1:%.2f\n[C-V] level2:%.2f",
		ebiten.CurrentFPS(), frequency, lacunarity, gain, octaves, level1, level2))

	return nil
}

func main() {
	frequency = 0.03
	lacunarity = 2
	gain = 0.5
	octaves = 3
	level1 = 0.4
	level2 = 0.65

	generator = opensimplex.New()

	startT = time.Now()

	img, _ = ebiten.NewImage(winSize, winSize, ebiten.FilterDefault)

	ebiten.Run(update, winSize, winSize, 1, "noize")
}

func procInput() {
	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		frequency *= 1.1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		frequency /= 1.1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		lacunarity *= 1.1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		lacunarity /= 1.1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		gain *= 1.1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		gain /= 1.1
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		level1 -= 0.05
		if level1 < 0 {
			level1 = 0
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyX) {
		level1 += 0.05
		if level1 > level2 {
			level1 = level2
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		level2 -= 0.05
		if level2 < level1 {
			level2 = level1
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyV) {
		level2 += 0.05
		if level2 > 1 {
			level2 = 1
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		if octaves > 1 {
			octaves--
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		octaves++
	}
}
