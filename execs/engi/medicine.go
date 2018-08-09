package main

const (
	MC_Radi = iota
	MC_Temp
	MC_CO2
	MC_Air
	MC_Hit
	MC_RadiTemp
	MC_RadiCO2
	MC_RadiAir
)

type MediOpts struct {
	Radi, Temp, CO2, Air, Hit, RadiTemp, RadiCO2, RadiAir CounterOpts
}
type CounterOpts struct {
	Levels [3]float64
	BioInf [3][7]int
	Nucleo [3][7]int
}

var counters [8]*mediCounters

func initMedi(shipId string) {
	o := DEFVAL.MediOpts
	counters[MC_Radi] = newCounter(o.Radi)
	counters[MC_Temp] = newCounter(o.Temp)
	counters[MC_CO2] = newCounter(o.CO2)
	counters[MC_Air] = newCounter(o.Air)
	counters[MC_Hit] = newCounter(o.Hit)
	counters[MC_RadiTemp] = newCounter(o.RadiTemp)
	counters[MC_RadiCO2] = newCounter(o.RadiCO2)
	counters[MC_RadiAir] = newCounter(o.RadiAir)
}

func (s *engiScene) checkMedicine() {
	radiInCockpit := s.local.radiationSum * DEFVAL.RadiCockPitK
	lostAir := DEFVAL.NormPressure - Data.EngiData.Counters.Air

	counters[MC_Radi].AddValue(radiInCockpit)
	counters[MC_Temp].AddValue(s.local.temperature)
	counters[MC_CO2].AddValue(Data.EngiData.Counters.CO2)
	counters[MC_Air].AddValue(lostAir)
	counters[MC_Hit].AddValue(Data.EngiData.Counters.Hitted)
	if radiInCockpit > DEFVAL.MedRadiLevel {
		counters[MC_RadiTemp].AddValue(Data.EngiData.Counters.CO2)
		counters[MC_RadiCO2].AddValue(lostAir)
		counters[MC_RadiAir].AddValue(Data.EngiData.Counters.Hitted)
	}
}

func (s *engiScene) procPhysMedicine(emiLvl float64) {
	//todo: some other way
	Data.EngiData.Counters.Hitted += emiLvl * DEFVAL.HitsToCounter
}
