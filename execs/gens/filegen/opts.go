package main

type Options struct {
	StarANCount        int
	SizeMassDevPercent float64
	OrbitDevPercent    float64

	SingleStar struct {
		R10  float64
		Size float64
		MaxG float64
	}

	DoubleStar struct {
		R10    float64
		Size   float64
		MaxG   float64
		Radius float64
		Period float64
	}

	TripleStar struct {
		R10    float64
		Size   float64
		MaxG   float64
		Radius float64
		Period float64
		Pair   struct {
			R10    float64
			Size   float64
			MaxG   float64
			Radius float64
			Period float64
		}
	}

	PlanetMinR float64
	ClosePeriod float64
	DistStep float64
	PeriodStep float64

	LastHardOnGasPercent int
	MoveHardOnHardPercent int

	FirstBeltPercent int
	MoreBeltsPercent int
	BeltOrbitDev     float64
	BeltSizeDev      float64

	AsteroidSize float64
	AsteroidR10 float64
	AsteroidG float64

	BeltCount int
	BeltCountDev float64

	PlanetOrbitDev float64
	GasDev float64
	GasSize float64
	GasG float64
	GasR10 float64
	HardDev float64
	HardSize float64
	HardG float64
	HardR10 float64

	SatelliteOrbitPart float64

	GasBeltPercent int
}

var Opts Options
