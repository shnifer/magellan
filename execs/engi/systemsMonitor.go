package main

import (
	"fmt"
	. "github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/draw"
	. "github.com/Shnifer/magellan/graph"
	"github.com/Shnifer/magellan/ranma"
	"github.com/Shnifer/magellan/v2"
	"golang.org/x/image/colornames"
	"image/color"
	"math"
	"strconv"
	"time"
)

type systemsMonitor struct {
	params CamParams

	sprites map[string]*Sprite
	sysNvs  [SysCount]v2.V2

	isEmid [8]bool
	isDmgd [8]bool
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

func (s *systemsMonitor) update(dt float64, ranma *ranma.Ranma) {
	e := Data.EngiData.Emissions
	s.isEmid[0] = 0 < e[EMI_VEL_UP]+e[EMI_VEL_DOWN]+e[EMI_ACCEL]+e[EMI_REVERSE]+e[EMI_ENGINE_HEAT]
	s.isEmid[1] = 0 < e[EMI_DIST_UP]+e[EMI_DIST_DOWN]+e[EMI_WARP_TURN]
	s.isEmid[2] = 0 < e[EMI_TURN]+e[EMI_STRAFE]
	s.isEmid[3] = 0 < e[EMI_RADAR_COSMOS]+e[EMI_RADAR_WARP]+e[EMI_RADAR_ANG_DOWN]+e[EMI_RADAR_ANG_UP]
	s.isEmid[4] = 0 < e[EMI_SCAN_RADIUS]+e[EMI_SCAN_SPEED]+e[EMI_DROP_RADIUS]+e[EMI_DROP_SPEED]
	s.isEmid[5] = 0 < e[EMI_FUEL]
	s.isEmid[6] = 0 < e[EMI_CO2]
	s.isEmid[7] = 0 < e[EMI_DEF_HEAT]+e[EMI_DEF_RADI]+e[EMI_DEF_MECH]

	for i := 0; i < SysCount; i++ {
		s.isDmgd[i] = ranma.GetOut(i) > 0
	}
}

func (s *systemsMonitor) Req(Q *DrawQueue) {
	Q.Add(s.sprites["all"], Z_GAME_OBJECT-1)
	Q.Add(s.sprites["upp"], Z_GAME_OBJECT+2)
	for i := 0; i < SysCount; i++ {
		sprite := s.sprites[strconv.Itoa(i)]
		var clr color.Color
		if s.isDmgd[i] {
			clr = colornames.Darkgray
		} else {
			clr = sysColor(Data.EngiData.AZ[i] / 100)
		}
		sprite.SetColor(clr)
		if s.isEmid[i] {
			t := float64(time.Now().Nanosecond()) / 1000000000
			a := math.Sin(t*2*math.Pi)*0.5 + 0.5
			sprite.SetAlpha(a)
		} else {
			sprite.SetAlpha(1)
		}
		Q.Add(sprite, Z_GAME_OBJECT)

		if i == 5 {
			sprite := s.sprites["fuelc"]
			sprite.SetColor(clr)
			Q.Add(sprite, Z_GAME_OBJECT)
		}

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
		} else if pos.Sub(s.sysNvs[i]).Len() < 200 {
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
	a("fuelc")
	a("tempc")
	a("airc")

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

func sysColor(v float64) color.Color {
	k := Clamp(v, 0, 1)
	if k > 0.5 {
		k := (k - 0.5) * 2
		return color.RGBA{
			R: uint8(255 * (1 - k)),
			G: 255,
			B: 0,
			A: 255,
		}
	} else {
		k := k * 2
		return color.RGBA{
			R: 255,
			G: uint8(255 * k),
			B: 0,
			A: 255,
		}
	}
}

func (s *systemsMonitor) s(pref string, n int) *Sprite {
	return s.sprites[pref+strconv.Itoa(n)]
}

func (s *systemsMonitor) drawAir(Q *DrawQueue) {
	v := Clamp(Data.EngiData.Counters.Air/Data.BSP.Lss.Air_volume, 0, 1)
	n := int(math.Round(v * 10))
	n *= 10
	Q.Add(s.s("a", n), Z_GAME_OBJECT+1)
	sprite := s.sprites["airc"]
	sprite.SetColor(sysColor(v))
	Q.Add(sprite, Z_GAME_OBJECT)
}

func (s *systemsMonitor) drawTemp(Q *DrawQueue) {
	v := Clamp(Data.EngiData.Counters.Calories/Data.BSP.Shields.Heat_capacity, 0, 2)
	n := int(math.Round(v * 10))
	if n > 10 {
		n = 11
	}
	n *= 10
	Q.Add(s.s("t", n), Z_GAME_OBJECT+1)
	sprite := s.sprites["tempc"]
	sprite.SetColor(sysColor(1 - v))
	Q.Add(sprite, Z_GAME_OBJECT)
}

func (s *systemsMonitor) drawFuel(Q *DrawQueue) {
	v := Clamp(Data.EngiData.Counters.Fuel/Data.BSP.Fuel_tank.Fuel_volume, 0, 1)
	n := int(math.Round(v * 20))
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
