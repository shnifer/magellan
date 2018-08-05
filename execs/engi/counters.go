package main

func CalculateCounters(dt float64) {
	Data.EngiData.Counters.Calories = calcHeat(dt)
}

func calcHeat(dt float64) float64 {
	pool := Data.EngiData.Counters.Calories
	pool += Data.PilotData.HeatProduction * dt
	pool -= Data.SP.Shields.Heat_sink * dt
	return pool
}
