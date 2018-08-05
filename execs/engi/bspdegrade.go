package main

import (
	. "github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/ranma"
)

func CalculateBSPDegrade(ranma *ranma.Ranma) (res BSPDegrade) {
	var b uint16
	k := func(x uint16) float64 {
		var n uint16
		for i := 0; i < 16; i++ {
			n += x & 1
			x = x >> 1
		}
		return float64(n)
	}
	f := func(v uint16, mask uint16) float64 {
		x := v & mask
		return k(x) / k(mask) * DEFVAL.RanmaMaxDegradePercent
	}
	d := func(v uint16, mask uint16) float64 {
		x := f(v, mask)
		if x > 1 {
			return 0
		} else {
			return 1 - x
		}
	}
	u := func(v uint16, mask uint16) float64 {
		return 1 + f(v, mask)
	}
	e := func(emi string) float64 {
		return 1 + Data.EngiData.Emissions[emi]*(DEFVAL.EmissionDegradePercent/100)
	}

	b = ranma.GetOut(SYS_MARCH)
	res.March_engine.Thrust_max = d(b, 36411)
	res.March_engine.Thrust_acc = d(b, 18863)
	res.March_engine.Thrust_slow = d(b, 9590)
	res.March_engine.Reverse_max = d(b, 29125)
	res.March_engine.Reverse_acc = d(b, 48721)
	res.March_engine.Reverse_slow = d(b, 56008)
	res.March_engine.Heat_prod = u(b, 5822) * e(EMI_ENGINE_HEAT)

	b = ranma.GetOut(SYS_WARP)
	res.Warp_engine.Distort_max = d(b, 36411)
	res.Warp_engine.Warp_enter_consumption = u(b, 18863)
	res.Warp_engine.Distort_acc = d(b, 9590)
	res.Warp_engine.Distort_slow = d(b, 29125)
	res.Warp_engine.Consumption = u(b, 48721)
	res.Warp_engine.Turn_speed = d(b, 56008)
	res.Warp_engine.Turn_consumption = u(b, 5822)

	b = ranma.GetOut(SYS_SHUNTER)
	res.Shunter.Turn_max = d(b, 36411)
	res.Shunter.Turn_acc = d(b, 18863)
	res.Shunter.Turn_slow = d(b, 9590)
	res.Shunter.Strafe_max = d(b, 29125)
	res.Shunter.Strafe_acc = d(b, 48721)
	res.Shunter.Strafe_slow = d(b, 56008)
	res.Shunter.Heat_prod = u(b, 5822) * e(EMI_ENGINE_HEAT)

	b = ranma.GetOut(SYS_RADAR)
	res.Radar.Range_Max = d(b, 36411) / e(EMI_RADAR_COSMOS) / e(EMI_RADAR_WARP)
	res.Radar.Angle_Min = u(b, 18863)
	res.Radar.Angle_Max = d(b, 9590)
	checkAngles(Data.BSP.Radar.Angle_Min, Data.BSP.Radar.Angle_Max, &res.Radar.Angle_Min, &res.Radar.Angle_Max)
	ak := e(EMI_RADAR_ANG_UP) / e(EMI_RADAR_ANG_DOWN)
	res.Radar.Angle_Min *= ak
	res.Radar.Angle_Max *= ak

	res.Radar.Angle_Change = d(b, 29125)
	res.Radar.Range_Change = d(b, 56008)
	res.Radar.Rotate_Speed = d(b, 5822)

	b = ranma.GetOut(SYS_SCANNER)
	res.Scanner.ScanRange = d(b, 36410) / e(EMI_SCAN_RADIUS)
	res.Scanner.ScanSpeed = d(b, 18862) / e(EMI_SCAN_SPEED)
	res.Scanner.DropRange = d(b, 9590) / e(EMI_DROP_RADIUS)
	res.Scanner.DropSpeed = d(b, 4830) / e(EMI_DROP_SPEED)

	b = ranma.GetOut(SYS_FUEL)
	res.Fuel_tank.Fuel_Protection = d(b, 8095)
	res.Fuel_tank.Radiation_def = k(b & 58355)

	b = ranma.GetOut(SYS_LSS)
	res.Lss.Thermal_def = d(b, 36411)
	//fixme:check this!
	res.Lss.Co2_level = d(b, 18862) / e(EMI_CO2)
	res.Lss.Air_prepare_speed = d(b, 9590)
	res.Lss.Lightness = d(b, 4831)

	b = ranma.GetOut(SYS_SHIELD)
	res.Shields.Radiation_def = d(b, 36411) / e(EMI_DEF_RADI)
	res.Shields.Disinfect_level = d(b, 18863)
	res.Shields.Mechanical_def = d(b, 9590) / e(EMI_DEF_MECH)
	res.Shields.Heat_reflection = d(b, 48721) / e(EMI_DEF_HEAT)
	res.Shields.Heat_capacity = d(b, 56008)
	res.Shields.Heat_sink = d(b, 5822)

	return res
}

func checkAngles(bMin, bMax float64, dMin, dMax *float64) {
	min := bMin * (*dMin)
	max := bMax * (*dMax)
	if min <= max {
		return
	}
	v := Clamp((min+max)/2, bMin, bMax)
	*dMin = v / bMin
	*dMax = v / bMax
}
