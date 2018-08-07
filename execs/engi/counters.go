package main

import (
	. "github.com/Shnifer/magellan/commons"
)

func CalculateCounters(dt float64) {
	Data.EngiData.Counters.Calories = calcHeat(dt)
}

func calcHeat(dt float64) float64 {
	pool := Data.EngiData.Counters.Calories
	pool += Data.PilotData.HeatProduction * DEFVAL.HeatProdTemperatureK * dt
	pool += Data.EngiData.Emissions[EMI_DMG_HEAT] * DEFVAL.EmiTemperatureK *
		(100 - Data.SP.Shields.Heat_reflection) / 100 * dt
	pool -= Data.SP.Shields.Heat_sink * dt
	if pool < 0 {
		pool = 0
	}
	return pool
}

type localCounters struct {
	temperature float64
}

func (s *engiScene) CalculateLocalCounters() {
	s.local.temperature = DEFVAL.NormTemperature +
		DEFVAL.PoolTemperatureK*Data.EngiData.Counters.Calories*Data.SP.Lss.Thermal_def
}

func initLocal() localCounters {
	return localCounters{
		temperature: DEFVAL.NormTemperature,
	}
}
