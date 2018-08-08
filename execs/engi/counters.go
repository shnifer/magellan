package main

import (
	. "github.com/Shnifer/magellan/commons"
	"math"
)

/*type EngiCounters struct {
Fuel     float64 `json:"f,omitempty"`
HoleSize float64 `json:"h,omitempty"`
Pressure float64 `json:"p,omitempty"`
Air      float64 `json:"a,omitempty"`
Calories float64 `json:"t,omitempty"`
CO2      float64 `json:"co2,omitempty"`
FlightTime float64 `json:"ft,omitempty"`}*/

func CalculateCounters(dt float64) {
	Data.EngiData.Counters.FlightTime += dt
	Data.EngiData.Counters.Calories = calcHeat(dt)
	Data.EngiData.Counters.CO2 = calcCO2(dt)
	Data.EngiData.Counters.Fuel = calcFuel(dt)
	calcHolePressureAir(dt)
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

func calcCO2(dt float64) float64 {
	const NormalCo2 = 0.04
	target := (700 - Data.BSP.Lss.Co2_level) / 33.3
	target *= e(EMI_CO2)
	N := Data.EngiData.BSPDegrade.Lss.Co2_level
	target = (target-NormalCo2)/8*N + NormalCo2
	current := Data.EngiData.Counters.CO2

	delta := target - current
	return current + delta*dt*DEFVAL.CO2StepK
}

func calcFuel(dt float64) float64 {
	if Data.State.GalaxyID != WARP_Galaxy_ID {
		return Data.EngiData.Counters.Fuel
	}

	return Data.EngiData.Counters.Fuel - Data.PilotData.Distortion*Data.SP.Warp_engine.Consumption*dt
}

func calcHolePressureAir(dt float64) {
	hs := Data.EngiData.Counters.HoleSize
	if hs <= 0 {
		return
	}
	p := Data.EngiData.Counters.Pressure
	AirLose := math.Min(DEFVAL.MaxAirLose, p*hs*DEFVAL.AirLoseQuot)
	Data.EngiData.Counters.Pressure -= AirLose * dt
	hs -= Data.SP.Shields.Mechanical_def * DEFVAL.RepairQuot * dt
	if hs <= 0 {
		hs = 0
	}
	Data.EngiData.Counters.HoleSize = hs

	if Data.EngiData.Counters.Air > 0 {
		AirRestore := math.Min(DEFVAL.NormPressure-p, Data.SP.Lss.Air_prepare_speed) * dt
		Data.EngiData.Counters.Pressure += AirRestore
		Data.EngiData.Counters.Air -= AirRestore
		if Data.EngiData.Counters.Air < 0 {
			Data.EngiData.Counters.Air = 0
		}
	}
}

type localCounters struct {
	temperature float64
	radiation   float64
}

func (s *engiScene) CalculateLocalCounters() {
	s.local.temperature = DEFVAL.NormTemperature +
		DEFVAL.PoolTemperatureK*Data.EngiData.Counters.Calories*
			(100-Data.SP.Lss.Thermal_def)/100

	outRadi := Data.EngiData.Emissions[EMI_DMG_RADI] * (100 - Data.SP.Shields.Radiation_def) / 100 * DEFVAL.OutRadiK
	tankRadi := DEFVAL.TankRadiBase + DEFVAL.TankRadiK*Data.EngiData.Counters.Fuel
	pointUndef := (240 - 0.4*Data.BSP.Fuel_tank.Radiation_def) / 1600
	totalUndef := Data.EngiData.BSPDegrade.Fuel_tank.Radiation_def * pointUndef
	if totalUndef > 1 {
		totalUndef = 1
	}
	inRadi := tankRadi * totalUndef * DEFVAL.InRadiK

	s.local.radiation = outRadi + inRadi
}

func initLocal() localCounters {
	return localCounters{
		temperature: DEFVAL.NormTemperature,
	}
}
