package main

import (
	. "github.com/Shnifer/magellan/commons"
	. "github.com/Shnifer/magellan/log"
	"math"
	"math/rand"
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
	target := (700 - Data.BSP.Lss.Co2_level*getBoostPow(SYS_LSS)) / 33.3
	target *= e(EMI_CO2)
	N := Data.EngiData.BSPDegrade.Lss.Co2_level
	target = (target-NormalCo2)/8*N + NormalCo2

	//kill after 90min
	if Data.EngiData.Counters.FlightTime > DEFVAL.KillCO2Mins*60 {
		target = 20
	}

	current := Data.EngiData.Counters.CO2

	delta := target - current
	return current + delta*dt*DEFVAL.CO2StepK
}

func calcFuel(dt float64) float64 {
	if Data.State.GalaxyID != WARP_Galaxy_ID {
		return Data.EngiData.Counters.Fuel
	}

	return Data.EngiData.Counters.Fuel -
		Data.PilotData.Distortion*Data.SP.Warp_engine.Consumption*dt -
		Data.PilotData.DistTurn*Data.SP.Warp_engine.Turn_consumption*dt
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
	temperature  float64
	radiationSum float64
}

func (s *engiScene) CalculateLocalCounters() {
	s.local.temperature = DEFVAL.NormTemperature +
		DEFVAL.PoolTemperatureK*Data.EngiData.Counters.Calories*
			(100-Data.SP.Lss.Thermal_def)/100

	outRadi := Data.EngiData.Emissions[EMI_DMG_RADI] * (100 - Data.SP.Shields.Radiation_def) / 100 * DEFVAL.OutRadiK
	tankRadi := DEFVAL.TankRadiBase + DEFVAL.TankRadiK*Data.EngiData.Counters.Fuel
	pointUndef := (240 - 0.4*Data.BSP.Fuel_tank.Radiation_def*getBoostPow(SYS_FUEL)) / 1600
	totalUndef := Data.EngiData.BSPDegrade.Fuel_tank.Radiation_def * pointUndef
	if totalUndef > 1 {
		totalUndef = 1
	}
	inRadi := tankRadi * totalUndef * DEFVAL.InRadiK

	s.local.radiationSum = outRadi + inRadi
}

func initLocal() localCounters {
	return localCounters{
		temperature: DEFVAL.NormTemperature,
	}
}

func (s *engiScene) checkDamage() {
	if Data.EngiData.Emissions[EMI_DMG_MECH] > 0 {
		s.procPhys()
	}
	if s.local.radiationSum > 0 {
		s.procRadiation(s.local.radiationSum)
	}
	overHeat := Data.EngiData.Counters.Calories - Data.SP.Shields.Heat_capacity
	if overHeat > 0 {
		s.procOverHeat(overHeat)
	}
}

func (s *engiScene) procPhys() {
	hMin, hMax := DEFVAL.MinHole, DEFVAL.MaxHole
	makeHole := Data.EngiData.Emissions[EMI_DMG_MECH] *
		(rand.Float64()*(hMax-hMin) + hMin)
	Data.EngiData.Counters.HoleSize += makeHole
	if Data.EngiData.Counters.HoleSize < hMin {
		Data.EngiData.Counters.HoleSize = hMin
	} else if Data.EngiData.Counters.HoleSize > hMax {
		Data.EngiData.Counters.HoleSize = hMax
	}

	var countSys int
	var doMedical bool
	switch rand.Intn(8) {
	case 1, 2:
		countSys = 1
	case 3, 4:
		countSys = 1
		doMedical = true
	case 5:
		countSys = 2
	case 6:
		countSys = 2
		doMedical = true
	case 7:
		doMedical = true
	case 0:
		countSys = 3
	}

	Log(LVL_WARN, "Phys damage ", countSys, "x", Data.EngiData.Counters.HoleSize*DEFVAL.HoleAZK)
	s.doAZDamage(countSys, makeHole*DEFVAL.HoleAZK)
	if doMedical {
		s.procPhysMedicine(Data.EngiData.Emissions[EMI_DMG_MECH])
	}
}

func (s *engiScene) procOverHeat(overheat float64) {
	var countSys int
	switch rand.Intn(8) {
	case 0:
	case 1, 2, 3, 4:
		countSys = 1
	case 5, 6:
		countSys = 2
	case 7:
		countSys = 3
	}
	Log(LVL_WARN, "Heat damage")
	s.doAZDamage(countSys, overheat*DEFVAL.OverheatAZK)
}

func (s *engiScene) procRadiation(radiation float64) {
	var countSys int
	switch rand.Intn(8) {
	case 0:
	case 1, 2, 3, 4:
		countSys = 1
	case 5, 6:
		countSys = 2
	case 7:
		countSys = 3
	}
	Log(LVL_WARN, "Radi damage")
	s.doAZDamage(countSys, radiation*DEFVAL.RadiAZK)
}

func (s *engiScene) doAZDamage(repeats int, dmg float64) {
	if dmg == 0 || repeats == 0 {
		return
	}

	brakeSystems := make(map[int]struct{})

	for i := 0; i < repeats; i++ {
		sysN := getRandomSysByValue()
		dmg = math.Min(dmg, Data.EngiData.AZ[sysN])
		var brakeChance float64
		if Data.EngiData.AZ[sysN] == 0 {
			brakeChance = 1
		} else {
			brakeChance = dmg / Data.EngiData.AZ[sysN] * DEFVAL.BrakeChanceK
		}
		Data.EngiData.AZ[sysN] -= dmg
		if rand.Float64() < brakeChance {
			brakeSystems[sysN] = struct{}{}
		}
	}

	for sysN := range brakeSystems {
		Log(LVL_WARN, "we brake sys#", sysN)
		v := uint16(rand.Intn(65536))
		s.ranma.SetIn(sysN, v)
	}
}

func (s *engiScene) doTargetAZDamage(sysN int, dmg float64) {
	if dmg == 0 || sysN > SysCount || sysN < 0 {
		return
	}

	dmg = math.Min(dmg, Data.EngiData.AZ[sysN])
	var brakeChance float64
	if Data.EngiData.AZ[sysN] == 0 {
		brakeChance = 1
	} else {
		brakeChance = dmg / Data.EngiData.AZ[sysN] * DEFVAL.BrakeChanceK
	}
	Data.EngiData.AZ[sysN] -= dmg
	if rand.Float64() < brakeChance {
		Log(LVL_WARN, "we brake sys#", sysN)
		v := uint16(rand.Intn(65536))
		s.ranma.SetIn(sysN, v)
	}
}

func getRandomSysByValue() int {
	sum := Data.BSP.March_engine.Volume +
		Data.BSP.Shunter.Volume +
		Data.BSP.Warp_engine.Volume +
		Data.BSP.Shields.Volume +
		Data.BSP.Radar.Volume +
		Data.BSP.Scanner.Volume +
		Data.BSP.Lss.Volume +
		Data.BSP.Fuel_tank.Volume

	if sum == 0 {
		Log(LVL_ERROR, "Zero ship .Volumes!")
		return 0
	}
	v := rand.Float64() * sum
	if v < Data.BSP.March_engine.Volume {
		return SYS_MARCH
	} else {
		v -= Data.BSP.March_engine.Volume
	}
	if v < Data.BSP.Shunter.Volume {
		return SYS_SHUNTER
	} else {
		v -= Data.BSP.Shunter.Volume
	}
	if v < Data.BSP.Warp_engine.Volume {
		return SYS_WARP
	} else {
		v -= Data.BSP.Warp_engine.Volume
	}
	if v < Data.BSP.Shields.Volume {
		return SYS_SHIELD
	} else {
		v -= Data.BSP.Shields.Volume
	}
	if v < Data.BSP.Radar.Volume {
		return SYS_RADAR
	} else {
		v -= Data.BSP.Radar.Volume
	}
	if v < Data.BSP.Scanner.Volume {
		return SYS_SCANNER
	} else {
		v -= Data.BSP.Scanner.Volume
	}
	if v < Data.BSP.Lss.Volume {
		return SYS_LSS
	} else {
		v -= Data.BSP.Lss.Volume
	}
	return SYS_FUEL
}
