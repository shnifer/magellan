package main

import (
	"fmt"
	"github.com/shnifer/magellan/commons"
	"github.com/shnifer/magellan/draw"
	"github.com/shnifer/magellan/graph"
	"golang.org/x/image/colornames"
	"strings"
)

func (s *engiScene) systemInfo(n int) *graph.Text {
	var msg string
	moni := s.systemsMonitor
	a := func(s string) {
		msg += s + "\n"
	}
	af := func(name string, v float64) {
		if v == 1 {
			a(name + ": РАБОТАЕТ ШТАТНО")
			return
		}
		a(fmt.Sprintf("%v: %.1f", name, v*100))
	}
	var sysNames = [8]string{"Маршевый двигатель", "Гипер двигатель", "Маневровые двигатели", "Радар",
		"Сканеры и посадка", "Топливный бак", "Системы жизнеобеспечения", "Защитные системы"}

	a(sysNames[n])
	if moni.isDmgd[n] {
		a("ПОВРЕЖДЕНА")
	}
	if moni.isEmid[n] {
		a("ВНЕШНЕЕ ВОЗДЕЙСТВИЕ")
	}
	if moni.isBoosted[n] {
		a("ПОД БУСТОМ")
	}
	AZ := Data.EngiData.AZ[n] * getBoostPow(n)
	af("Ресурс защиты", AZ/100)
	a("")

	switch n {
	case commons.SYS_MARCH:
		deg := Data.EngiData.BSPDegrade.March_engine
		af("Мощность", deg.Thrust_max)
		af("Разгон", deg.Thrust_acc)
		af("Сброс", deg.Thrust_slow)
		a("")
		if Data.BSP.BSPParams.March_engine.Reverse_max > 0 {
			af("Реверс", deg.Reverse_max)
			af("Разгон", deg.Thrust_slow)
			af("Сброс", deg.Thrust_slow)
			a("")
		}
		af("Тепловыделение", deg.Heat_prod)
	case commons.SYS_SHUNTER:
		deg := Data.EngiData.BSPDegrade.Shunter
		af("Поворот", deg.Turn_max)
		af("Разгон", deg.Turn_acc)
		af("Сброс", deg.Turn_slow)
		a("")
		af("Тепловыделение", deg.Heat_prod)
	case commons.SYS_WARP:
		deg := Data.EngiData.BSPDegrade.Warp_engine
		af("Макс. искривление", deg.Distort_max)
		af("Разгон", deg.Distort_acc)
		af("Сброс", deg.Distort_slow)
		af("Поворот", deg.Turn_speed)
		a("")
		af("Потребление топлива", deg.Consumption)
		af("Потребление на поворот", deg.Turn_consumption)
		af("Цена входа в гипер", deg.Warp_enter_consumption)
	case commons.SYS_SHIELD:
		deg := Data.EngiData.BSPDegrade.Shields
		af("Механическая защита", deg.Mechanical_def)
		af("Радиационная защита", deg.Radiation_def)
		a("")
		af("Отражатель излучения", deg.Heat_reflection)
		af("Ёмкость накопителя", deg.Heat_capacity)
		af("Охлаждение", deg.Heat_sink)
	case commons.SYS_RADAR:
		deg := Data.EngiData.BSPDegrade.Radar
		af("Макс. дальность", deg.Range_Max)
		af("Изм. дальности", deg.Range_Max)
		af("Скорость поворота", deg.Rotate_Speed)
		a("")
		af("Макс. угол", deg.Angle_Max)
		af("Мин. угол", deg.Angle_Min)
		af("Изменение угла", deg.Angle_Change)
	case commons.SYS_SCANNER:
		deg := Data.EngiData.BSPDegrade.Scanner
		af("Дальность сканирования", deg.ScanRange)
		af("Скорость сканирования", deg.ScanSpeed)
		a("")
		af("Дальность сброса", deg.DropRange)
		af("Скорость сброса", deg.DropSpeed)
	case commons.SYS_FUEL:
		deg := Data.EngiData.BSPDegrade.Fuel_tank
		af("Радиационная защита", deg.Radiation_def)
		af("Защита топлива", deg.Fuel_Protection)
	case commons.SYS_LSS:
		deg := Data.EngiData.BSPDegrade.Lss
		af("Термическая защита", deg.Thermal_def)
		af("Регенератор CO2", deg.Co2_level)
		af("Поддержание давления", deg.Air_prepare_speed)
	}

	msg = strings.TrimSuffix(msg, "\n")
	res := graph.NewText(msg, draw.Fonts[draw.Face_cap], colornames.White)
	res.SetPosPivot(graph.ScrP(0.25, 0.25), graph.TopLeft())
	return res
}
