package main

import (
	"encoding/json"
	"github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/v2"
	"image/color"
	"io/ioutil"
	"strconv"
	"strings"
)

type fileData struct {
	ID, Parent  string
	Diameter    float64
	Distance    float64
	Mass        float64
	OrbitPeriod float64
	Color       struct{ R, G, B byte }
	Count       int
	//начальный угол, если объект 1
	StartAng float64
	//отклонения от базовых значений в процентах, если объектов много
	RadMassDev     float64
	PeriodOrbitDev float64
	TexName        string
}

const DEFType = "planet"
const K_OrbitPeriod = 1.0
const K_Radius = 1.0

func main() {
	buf, err := ioutil.ReadFile("galaxyPredata.json")
	if err != nil {
		panic(err)
	}

	var inData []fileData
	err = json.Unmarshal(buf, &inData)
	if err != nil {
		panic(err)
	}

	var outData commons.Galaxy
	maxOrbit := 0.0

	for _, v := range inData {
		if v.Count == 1 {
			if v.Distance > maxOrbit {
				maxOrbit = v.Distance
			}
			gp := createGP(v)
			outData.Points = append(outData.Points, gp)
		} else {
			for i := 0; i < v.Count; i++ {
				w := v
				w.ID = v.ID + "-" + strconv.Itoa(i)

				kPeriodOrbit := commons.KDev(v.PeriodOrbitDev)
				w.OrbitPeriod *= kPeriodOrbit
				w.Distance *= kPeriodOrbit

				kRadMass := commons.KDev(v.RadMassDev)
				w.Diameter *= kRadMass
				w.Mass *= kRadMass

				gp := createGP(w)
				outData.Points = append(outData.Points, gp)
			}
		}
	}

	outData.SpawnDistance = maxOrbit * 1.1

	buf, err = json.Marshal(outData)
	if err != nil {
		panic(err)
	}

	ioutil.WriteFile("galaxy_solar.json", buf, 0)
}

func createGP(v fileData) commons.GalaxyPoint {
	objType := DEFType
	if v.TexName != "" {
		s := strings.Split(v.TexName, ".")
		objType = s[0]
	}

	color := color.RGBA{R: v.Color.R, G: v.Color.G, B: v.Color.B, A: 255}

	okr := func(x float64) float64 {
		const sgn = 10
		return float64(int(x*sgn)) / sgn
	}

	gp := commons.GalaxyPoint{
		ID:       v.ID,
		ParentID: v.Parent,
		Pos:      v2.V2{},
		Orbit:    okr(v.Distance * K_Radius),
		Period:   okr(v.OrbitPeriod * K_OrbitPeriod),
		Type:     objType,
		Size:     okr(v.Diameter / 2),
		Color:    color,
		Mass:     okr(v.Mass),
		ScanData: v.ID,
	}
	return gp
}
