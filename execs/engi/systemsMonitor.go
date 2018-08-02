package main

import (
	"fmt"
	"github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/draw"
	. "github.com/Shnifer/magellan/graph"
	"github.com/Shnifer/magellan/v2"
	"image/color"
	"math"
	"strconv"
)

type systemsMonitor struct {
	params CamParams

	sprites map[string]*Sprite
	sysNvs  [SysCount]v2.V2
}

func newSystemsMonitor() *systemsMonitor {

	cam := NewCamera()
	cam.DenyGlobalScale = true
	cam.Center = ScrP(0.5, 0.5)

	scale1 := CalcGlobalScale(int(float64(WinH) / float64(DEFVAL.SpriteSizeH) * 1000))
	scale2 := CalcGlobalScale(int(float64(WinW) / float64(DEFVAL.SpriteSizeW) * 1000))
	scale := math.Min(scale1, scale2)

	cam.Scale = scale
	cam.Recalc()

	param := cam.Phys()

	sprites := getSprites(param)
	sysNvs := getSysNvs(sprites["all"], param)

	res := systemsMonitor{
		params:  param,
		sprites: sprites,
		sysNvs:  sysNvs,
	}

	return &res
}

func (s *systemsMonitor) Req(Q *DrawQueue) {
	Q.Add(s.sprites["all"], Z_GAME_OBJECT-1)
	Q.Add(s.sprites["upp"], Z_GAME_OBJECT+2)
	for i := 0; i < SysCount; i++ {
		sprite := s.sprites[strconv.Itoa(i)]
		sprite.SetColor(sysColor(Data.EngiData.AZ[i]))
		Q.Add(sprite, Z_GAME_OBJECT)

		s.writeSysText(Q, i)
	}
	s.drawTemp(Q)
	s.drawFuel(Q)
	s.drawAir(Q)
}

func (s *systemsMonitor) mouseOverSystem(pos v2.V2) (sysN int, isOver bool) {

	for i := 0; i < SysCount; i++ {
		if s.sprites[strconv.Itoa(i)].IsOver(pos, true) {
			return i, true
		}
	}

	return 0, false
}

func getSprites(params CamParams) map[string]*Sprite {
	res := make(map[string]*Sprite)
	a := func(id string) {
		res[id] = draw.NewAtlasSprite("engi_"+id, params)
	}

	a("all")
	a("0")
	a("1")
	a("2")
	a("3")
	a("4")
	a("5")
	a("6")
	a("7")
	a("upp")

	for i := 0; i <= 100; i += 5 {
		a("f" + strconv.Itoa(i))
	}

	for i := 0; i <= 100; i += 10 {
		a("a" + strconv.Itoa(i))
	}

	for i := 0; i <= 110; i += 10 {
		a("t" + strconv.Itoa(i))
	}

	return res
}

func sysColor(az float64) color.Color {
	k := commons.Clamp(az/100, 0, 1)
	return color.RGBA{
		R: uint8(255 * k),
		G: uint8(255 * (1 - k)),
		B: 0,
		A: 255,
	}
}

func (s *systemsMonitor) s(pref string, n int) *Sprite {
	return s.sprites[pref+strconv.Itoa(n)]
}

func (s *systemsMonitor) drawAir(Q *DrawQueue) {
	n := int(math.Round(commons.Clamp(Data.EngiData.Air/Data.BSP.Lss.Air_volume, 0, 1) * 10))
	n *= 10
	Q.Add(s.s("a", n), Z_GAME_OBJECT+1)
}

func (s *systemsMonitor) drawTemp(Q *DrawQueue) {
	n := int(math.Round(commons.Clamp(Data.EngiData.Calories/Data.BSP.Shields.Heat_capacity, 0, 2) * 10))
	if n > 10 {
		n = 11
	}
	n *= 10
	Q.Add(s.s("t", n), Z_GAME_OBJECT+1)
}

func (s *systemsMonitor) drawFuel(Q *DrawQueue) {
	n := int(math.Round(commons.Clamp(Data.EngiData.Fuel/Data.BSP.Fuel_tank.Fuel_volume, 0, 1) * 20))
	n *= 5
	Q.Add(s.s("f", n), Z_GAME_OBJECT+1)
}

func getSysNvs(sprite *Sprite, params CamParams) (res [SysCount]v2.V2) {
	var x1, x2, x3, x4 float64 = 90, 300, 700, 910
	var y1, y2 float64 = 50, 600
	xi := [SysCount]float64{x4, x3, x4, x1, x2, x2, x1, x3}
	yi := [SysCount]float64{y2, y1, y1, y2, y2, y1, y1, y2}

	_, op := sprite.ImageOp()
	G := op.GeoM
	for i := 0; i < SysCount; i++ {
		x, y := G.Apply(xi[i], yi[i])
		res[i] = v2.V2{X: x, Y: y}
	}
	return res
}

func (s *systemsMonitor) writeSysText(Q *DrawQueue, n int) {
	var sysNames = [8]string{"Маршевый\nдвигатель", "Гипер\nдвигатель", "Маневровые\nдвигатели", "Радар",
		"Сканеры\nи посадка", "Топливный\nбак", "Системы\nжизнеобеспечения", "Защитные\nсистемы"}
	az := Data.EngiData.AZ[n]
	msg := fmt.Sprintf("%v\n %.1f%%", sysNames[n], az)
	clr := sysColor(az)
	t := NewText(msg, draw.Fonts[draw.Face_list], clr)
	t.SetPosPivot(s.sysNvs[n], Center())
	Q.Add(t, Z_GAME_OBJECT+1)
}
