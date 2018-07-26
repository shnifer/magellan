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
}

var Opts Options
