package main

import (
	. "github.com/Shnifer/magellan/commons"
	"math"
	"math/rand"
	"strconv"
)

var newEl, eEl bool
var starCount int

func createPlanets(stat WarpStat, points map[string]*GalaxyPoint, pref string, planets []Planet) {
	newEl = len(stat.GasList)+len(stat.MineralList) > 0
	eEl = len(stat.EGas)+len(stat.EMetals) > 0
	starCount = stat.StarCount

	var parentID string
	var minR float64
	if stat.StarCount == 1 {
		parentID = pref + "S"
		minR = Opts.PlanetMinR
	} else {
		parentID = pref + "sv"
		if stat.StarCount == 2 {
			minR = math.Max(points[pref+"S1"].Orbit, points[pref+"S2"].Orbit) + Opts.PlanetMinR
		} else {
			minR = 2*math.Max(points[pref+"S1"].Orbit, points[pref+"sv2"].Orbit) + Opts.PlanetMinR
		}
	}

	//0 or 1
	var moveLastHardOnGas int
	if stat.HardPlanetsCount > 0 && len(planets) > stat.HardPlanetsCount {
		if rand.Intn(100) < Opts.LastHardOnGasPercent {
			moveLastHardOnGas = 1
		}
	}

	//num of moved
	moveHardOnHard := -1
	if stat.HardPlanetsCount-moveLastHardOnGas > 1 {
		if rand.Intn(100) < Opts.MoveHardOnHardPercent {
			moveHardOnHard = rand.Intn(stat.HardPlanetsCount - moveLastHardOnGas - 1)
		}
	}

	var asteroidBelted bool

	n := 0
	k := math.Max(1, KDev(Opts.StepDev))
	dist := minR * k
	period := Opts.ClosePeriod * k
	for n < len(planets) {
		//asteroid belts
		if asteroidBelted {
			if rand.Intn(100) < Opts.MoreBeltsPercent {
				addBelt(points, parentID, dist, period, pref+strconv.Itoa(n)+"-")
				nextDistPeriod(&dist, &period)
				continue
			}
		} else {
			if rand.Intn(100) < Opts.FirstBeltPercent {
				addBelt(points, parentID, dist, period, pref+strconv.Itoa(n)+"-")
				nextDistPeriod(&dist, &period)
				asteroidBelted = true
				continue
			}
		}

		//move hard on hard
		if moveHardOnHard == n {
			addPlanetWithPlanet(points, parentID, dist, period, planets[n], planets[n+1], pref)
			n += 2
			nextDistPeriod(&dist, &period)
			continue
		}

		if moveLastHardOnGas == 1 && n == stat.HardPlanetsCount-1 {
			addPlanetWithPlanet(points, parentID, dist, period, planets[n], planets[n+1], pref)
			n += 2
			nextDistPeriod(&dist, &period)
			continue
		}

		addPlanet(points, parentID, dist, period, planets[n], pref)
		nextDistPeriod(&dist, &period)
		n++
	}
}

func addPlanetWithPlanet(points map[string]*GalaxyPoint, parent string, dist, period float64,
	smallplanet, mainplanet Planet, pref string) {

	mainid := addPlanet(points, parent, dist, period, mainplanet, pref)
	dist, period = satellite(dist, period)
	addPlanet(points, mainid, dist, period, smallplanet, pref)
}

func addPlanet(points map[string]*GalaxyPoint, parent string, dist, period float64, planet Planet, pref string) string {
	t := GPT_HARDPLANET
	if planet.IsGas {
		t = GPT_GASPLANET
	}
	var size, g, r10 float64
	if planet.IsGas {
		k := KDev(Opts.GasDev)
		size = Opts.GasSize * k
		g = Opts.GasG * k
		r10 = Opts.GasR10 * k
	} else {
		k := KDev(Opts.HardDev)
		size = Opts.HardSize * k
		g = Opts.HardG * k
		r10 = Opts.HardR10 * k
	}

	k := KDev(Opts.PlanetOrbitDev)
	id := pref + planet.ID
	points[id] = pOpts{
		t:        t,
		parent:   parent,
		orbit:    dist * k,
		period:   period * k,
		shps:     planet.Spheres,
		minerals: planet.Minerals,
		size:     size,
		maxG:     g,
		r10:      r10,
		phase:    0,
	}.gp()

	dist, period = satellite(dist, period)

	if planet.IsGas {
		if rand.Intn(100) < Opts.GasBeltPercent {
			addBelt(points, id, dist, period, id+"-")
		}
	}

	return id
}

func addBelt(points map[string]*GalaxyPoint, parent string, dist, period float64, pref string) {
	sphs := asteroidSphs()
	count := int(float64(Opts.BeltCount) * KDev(Opts.BeltCountDev))
	for i := 0; i < count; i++ {
		k := KDev(Opts.BeltOrbitDev)
		sk := KDev(Opts.BeltSizeDev)
		points[pref+strconv.Itoa(i)] = pOpts{
			t:      GPT_ASTEROID,
			parent: parent,
			phase:  360 * rand.Float64(),
			orbit:  dist * k,
			period: period * k,
			size:   Opts.AsteroidSize * sk,
			r10:    Opts.AsteroidR10 * sk,
			maxG:   Opts.AsteroidG * sk,
			shps:   sphs,
		}.gp()
	}
}

func asteroidSphs() [15]int {

	res := [15]int{}
	if !newEl {
		res[PEDOMETALS] = EARTH
	} else if eEl {
		res[PEDOMETALS] = EARTHANDNEW
	} else {
		res[PEDOMETALS] = NEW
	}

	switch rand.Intn(10) {
	case 0:
		res[WATER] = HARD
	case 1:
		res[WATER] = WAS
	}

	if res[WATER] > NONE {
		if !newEl {
			res[MIXTURES] = EARTH
		} else if starCount == 1 {
			switch rand.Intn(2) {
			case 0:
				res[MIXTURES] = EARTHANDNEW
			case 1:
				res[MIXTURES] = NEW
			}
		} else {
			switch rand.Intn(3) {
			case 0:
				res[MIXTURES] = EARTHANDNEW
			case 1:
				res[MIXTURES] = NEW
			case 2:
				res[MIXTURES] = RADICAL
			}
		}
	}

	return res
}

func nextDistPeriod(dist, period *float64) {
	k := KDev(Opts.StepDev)
	if k < 1 {
		k = 1
	}
	if *dist < Opts.DistStepLevel {
		*dist = *dist * ((Opts.DistStep-1)*k + 1)
		*period = *period * ((Opts.PeriodStep-1)*k + 1)
	} else {
		*dist = *dist * ((Opts.DistFarStep-1)*k + 1)
		*period = *period * ((Opts.PeriodFarStep-1)*k + 1)
	}
}

func satellite(dist, period float64) (sdist, speriod float64) {
	sdist = dist - (dist / Opts.DistStep)
	sdist /= Opts.SatelliteOrbitPart

	speriod = period / dist * sdist
	return
}
