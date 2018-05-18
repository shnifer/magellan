package main

import (
	"math"
)

func FtoB(f float64) byte {
	switch {
	case f > 1:
		return 255
	case f < 0:
		return 0
	default:
		return byte(f * 255)
	}
}

//tabTan[0] = tan(-89.5 grad)
//tabTan[179] = tan(89.5 grad)

//медленнее почемуто

//var tabTan sort.IntSlice
//
//func init()  {
//	tabTan=make(sort.IntSlice,180)
//	for i:=0; i<180; i++{
//		tabTan[i] = int(math.Tan((float64(i)-90+0.5)*math.Pi/180)*1000)
//	}
//	fmt.Println(tabTan)
//}

//integer angle degree [0,360)
func getDeg(x, y int) (deg int) {
	const Pi = math.Pi
	if x == 0 {
		deg = 90
		if y < 0 {
			deg *= -1
		}
	} else {
		t := float64(y) / float64(x)
		ang := math.Atan(t)
		deg = int(ang / Pi * 180)
	}
	if x < 0 {
		deg += 180
	}
	if deg < 0 {
		deg += 360
	}
	return deg
}

//normalized absoluted [0;180]
func defDeg(deg1, deg2 int) (def int) {
	def = deg1 - deg2
	for def >= 180 {
		def -= 360
	}
	for def < (-180) {
		def += 360
	}
	if def < 0 {
		def *= -1
	}
	return def
}

const CoreRadius = 0.25
const CoreHalo = 0.40

//rad - часть от радиуса галактики
//swirlang - угол поворота при удалении на RadiusGalaxy
func Whirl(deg, midDeg, widthDeg int, rad float64, swirl float64) float64 {
	k := 1.0
	if rad < CoreHalo {
		k = rad / CoreHalo
		widthDeg = int(float64(widthDeg) / k)
	}

	midDeg += int(rad * swirl)

	def := defDeg(deg, midDeg)
	if def > widthDeg {
		return 0
	}
	return k * (1 - float64(def)/float64(widthDeg))
}

func Core(rad float64) float64 {
	//уменьшение яркости ядра при удалении на CoreRadius
	const CoreRadiusDeplete = 0.5
	if rad > CoreHalo {
		return 0
	}
	if rad < CoreRadius {
		return 1 - rad/CoreRadius*CoreRadiusDeplete
	}
	if (CoreHalo - CoreRadius) <= 0 {
		return 0
	}
	return (1 - (rad-CoreRadius)/(CoreHalo-CoreRadius)) * (1 - CoreRadiusDeplete)
}

func max(vs ...float64) (res float64) {
	for _, v := range vs {
		if v > res {
			res = v
		}
	}
	return res
}

//takes x,y in GalaxyUnits, return Density of star population
//0 -can't, 255 - max
func Dens(x, y int) byte {
	rad := math.Sqrt(float64(x*x+y*y)) / GalaxyRadius
	radk := 1 - rad
	//angle beams
	deg := getDeg(x, y)
	r1 := Whirl(deg, 45, 45, rad, 270)
	r2 := Whirl(deg, -45, 45, rad, 270)
	r3 := Whirl(deg, 135, 45, rad, 270)
	r4 := Whirl(deg, -135, 45, rad, 270)
	c := Core(rad) * 0.8
	fon := 0.1

	r := max(r1, r2, r3, r4, c, fon)

	return FtoB(radk * r)
}
