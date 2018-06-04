package main

import (
	"fmt"
	. "github.com/Shnifer/magellan/commons"
	. "github.com/Shnifer/magellan/draw"
	"github.com/Shnifer/magellan/graph"
	. "github.com/Shnifer/magellan/log"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"golang.org/x/image/colornames"
	"reflect"
	"strconv"
)

type engiScene struct {
	caption    *graph.Text
	background *graph.Sprite
}

func newEngiScene() *engiScene {
	caption := graph.NewText("Engi scene", Fonts[Face_cap], colornames.Aliceblue)
	caption.SetPosPivot(graph.ScrP(0.8, 0.1), graph.TopLeft())

	back := NewAtlasSpriteHUD("engibackground")
	back.SetSize(float64(WinW), float64(WinH))
	back.SetPivot(graph.TopLeft())

	return &engiScene{
		caption:    caption,
		background: back,
	}
}

func (*engiScene) Init() {
	defer LogFunc("engiScene.Init")()
}

func (scene *engiScene) Update(dt float64) {
	defer LogFunc("engiScene.Update")()

	emissions := CalculateEmissions(Data.Galaxy, Data.PilotData.Ship.Pos)
	heat := emissions[EMISSION_HEAT]

	prod := Data.PilotData.HeatProduction
	cool := Data.SP.Thrust_heat_sink

	Data.EngiData.HeatCumulated += (prod + heat - cool) * dt
	if Data.EngiData.HeatCumulated < 0 {
		Data.EngiData.HeatCumulated = 0
	}

	maxHeat := Data.SP.Thrust_heat_capacity
	if Data.EngiData.HeatCumulated > maxHeat {
		Data.EngiData.DmgCumulated[0] += dt
	}

	procDamage()

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		repairAll()
	}
}

func (scene *engiScene) Draw(image *ebiten.Image) {
	defer LogFunc("engiScene.Draw")()

	names, vals := degradeMsg()
	degrade := graph.NewText(names, Fonts[Face_list], colornames.Mediumvioletred)
	degrade.SetPosPivot(graph.ScrP(0.1, 0.1), graph.TopLeft())

	w, _ := degrade.GetSize()
	p := graph.ScrP(0.1, 0.1)
	p.X += float64(w) + 20
	stats := graph.NewText(vals, Fonts[Face_list], colornames.Palevioletred)
	stats.SetPosPivot(p, graph.TopLeft())

	heatMsg := fmt.Sprintf("%.1f/%.1f", Data.EngiData.HeatCumulated, Data.SP.Thrust_heat_capacity)
	heat := graph.NewText(heatMsg, Fonts[Face_stats], colornames.Forestgreen)
	heat.SetPosPivot(graph.ScrP(0.7, 0.9), graph.TopLeft())

	Q := graph.NewDrawQueue()
	Q.Add(scene.background, graph.Z_STAT_BACKGROUND)
	Q.Add(scene.caption, graph.Z_STAT_HUD)
	Q.Add(degrade, graph.Z_HUD)
	Q.Add(stats, graph.Z_HUD)
	Q.Add(heat, graph.Z_HUD)

	Q.Run(image)
}

func (scene *engiScene) OnCommand(command string) {
}

func (*engiScene) Destroy() {
}

func degradeMsg() (names, vals string) {
	v := reflect.ValueOf(Data.EngiData.BSPDegrade).Elem()
	t := v.Type()
	fc := t.NumField()
	for i := 0; i < fc; i++ {
		x := v.Field(i).Float()
		if x == 0 {
			continue
		}
		names += t.Field(i).Name + ":\n"
		vals += strconv.Itoa(int((1-x)*100)) + "%\n"
	}
	return names, vals
}

func procDamage() {
	if Data.EngiData.DmgCumulated[0] > 1 {
		Data.EngiData.DmgCumulated[0] -= 1
		Add1(&Data.EngiData.BSPDegrade.Thrust, 0.1)
		Add1(&Data.EngiData.BSPDegrade.Thrust_acc, 0.1)
		Add1(&Data.EngiData.BSPDegrade.Thrust_slow, 0.1)
	}
}

func repairAll() {
	Data.EngiData.BSPDegrade = &BSP{}
	Data.EngiData.HeatCumulated = 0
}
