package main

import . "github.com/Shnifer/magellan/commons"

func createStars(stat WarpStat, points map[string]*GalaxyPoint, pref string) {
	switch stat.StarCount {
	case 1:
		log("just a star")
		points[pref+"S"] = pOpts{
			t:    GPT_STAR,
			r10:  Opts.SingleStar.R10,
			size: Opts.SingleStar.Size,
			maxG: Opts.SingleStar.MaxG,
		}.gp()
	case 2:
		points[pref+"sv"] = &GalaxyPoint{IsVirtual: true}

		kOrbitPeriod := KDev(Opts.OrbitDevPercent)
		r := Opts.DoubleStar.Radius * kOrbitPeriod
		period := Opts.DoubleStar.Period * kOrbitPeriod

		kr := KDev(Opts.OrbitDevPercent)

		log("2 stars, orbits: ", r*kr, ",", r/kr)
		points[pref+"S1"] = pOpts{
			t:      GPT_STAR,
			parent: pref + "sv",
			orbit:  r * kr,
			period: period,
			phase:  0,
			r10:    Opts.DoubleStar.R10 / kr,
			size:   Opts.DoubleStar.Size / kr,
			maxG:   Opts.DoubleStar.MaxG / kr,
		}.gp()
		points[pref+"S2"] = pOpts{
			t:      GPT_STAR,
			parent: pref + "sv",
			orbit:  r / kr,
			period: period,
			phase:  180,
			r10:    Opts.DoubleStar.R10 * kr,
			size:   Opts.DoubleStar.Size * kr,
			maxG:   Opts.DoubleStar.MaxG * kr,
		}.gp()
	case 3:
		points[pref+"sv"] = &GalaxyPoint{IsVirtual: true}

		kOrbitPeriod := KDev(Opts.OrbitDevPercent)
		r := Opts.TripleStar.Radius * kOrbitPeriod
		period := Opts.TripleStar.Period * kOrbitPeriod

		kr := KDev(Opts.OrbitDevPercent)

		log("3 stars, orbits main pair: ", r*kr, ",", r/kr)
		points[pref+"S1"] = pOpts{
			t:      GPT_STAR,
			parent: pref + "sv",
			orbit:  r * kr,
			period: period,
			phase:  0,
			r10:    Opts.DoubleStar.R10 / kr,
			size:   Opts.DoubleStar.Size / kr,
			maxG:   Opts.DoubleStar.MaxG / kr,
		}.gp()
		points[pref+"sv2"] = &GalaxyPoint{
			IsVirtual: true,
			ParentID:  pref + "sv",
			Orbit:     r / kr,
			Period:    period,
			AngPhase:  180,
		}

		kOrbitPeriod = KDev(Opts.OrbitDevPercent)
		r = Opts.TripleStar.Pair.Radius * kOrbitPeriod
		period = Opts.TripleStar.Pair.Period * kOrbitPeriod

		log("orbits in pair: ", r*kr, ",", r/kr)
		kr = KDev(Opts.OrbitDevPercent)
		points[pref+"S2"] = pOpts{
			t:      GPT_STAR,
			parent: pref + "sv2",
			orbit:  r * kr,
			period: period,
			phase:  0,
			r10:    Opts.TripleStar.Pair.R10 / kr,
			size:   Opts.TripleStar.Pair.Size / kr,
			maxG:   Opts.TripleStar.Pair.MaxG / kr,
		}.gp()
		points[pref+"S3"] = pOpts{
			t:      GPT_STAR,
			parent: pref + "sv2",
			orbit:  r / kr,
			period: period,
			phase:  180,
			r10:    Opts.TripleStar.Pair.R10 * kr,
			size:   Opts.TripleStar.Pair.Size * kr,
			maxG:   Opts.TripleStar.Pair.MaxG * kr,
		}.gp()
	}
}
