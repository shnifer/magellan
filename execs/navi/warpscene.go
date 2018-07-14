package main

import (
	"fmt"
	. "github.com/Shnifer/magellan/commons"
	. "github.com/Shnifer/magellan/draw"
	"github.com/Shnifer/magellan/graph"
	"github.com/Shnifer/magellan/input"
	. "github.com/Shnifer/magellan/log"
	"github.com/Shnifer/magellan/v2"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
)

type warpScene struct {
	ship    *graph.Sprite
	caption *graph.Text
	cam     *graph.Camera

	sonar *graph.Sector

	face font.Face

	q *graph.DrawQueue
}

func newWarpScene() *warpScene {
	caption := graph.NewText("Navi warp scene", Fonts[Face_cap], colornames.Aliceblue)
	caption.SetPosPivot(graph.ScrP(0.1, 0.1), graph.TopLeft())

	cam := graph.NewCamera()
	cam.Center = graph.ScrP(0.5, 0.5)
	cam.Recalc()

	ship := graph.NewSprite(GetAtlasTex(ShipAN), cam.Phys())
	ship.SetSize(100, 100)
	ship.SetAlpha(0.5)

	sonarSector := graph.NewSector(cam.Phys())
	sonarSector.SetColor(colornames.Indigo)

	return &warpScene{
		caption: caption,
		ship:    ship,
		cam:     cam,
		face:    Fonts[Face_stats],
		sonar:   sonarSector,
		q:       graph.NewDrawQueue(),
	}
}

func (s *warpScene) Init() {
	defer LogFunc("cosmoScene.Init")()

}

func (s *warpScene) Update(dt float64) {
	defer LogFunc("cosmoScene.Update")()
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mousex, mousey := ebiten.CursorPosition()
		s.procMouseClick(v2.V2{X: float64(mousex), Y: float64(mousey)})
	}

	turn := input.GetF("sonarturn")
	rang := input.GetF("sonarrange")
	wide := input.GetF("sonarwide")

	Data.NaviData.SonarDir += turn * dt * Data.SP.Radar.Rotate_Speed
	Data.NaviData.SonarRange += rang * dt * Data.SP.Radar.Range_Change
	Data.NaviData.SonarRange = Clamp(Data.NaviData.SonarRange, 0, Data.SP.Radar.Range_Max)
	Data.NaviData.SonarWide += wide * dt * Data.SP.Radar.Angle_Change
	Data.NaviData.SonarWide = Clamp(Data.NaviData.SonarWide, Data.SP.Radar.Angle_Min, Data.SP.Radar.Angle_Max)

	//PilotData Rigid Body emulation
	Data.PilotData.Ship = Data.PilotData.Ship.Extrapolate(dt)

	s.sonar.SetRadius(Data.NaviData.SonarRange)
	s.sonar.SetAngles(
		Data.NaviData.SonarDir-Data.NaviData.SonarWide/2,
		Data.NaviData.SonarDir+Data.NaviData.SonarWide/2)

	s.ship.SetAng(Data.PilotData.Ship.Ang)
}

func (s *warpScene) Draw(image *ebiten.Image) {
	defer LogFunc("cosmoScene.Draw")()

	Q := s.q
	Q.Clear()

	Q.Add(s.sonar, graph.Z_UNDER_OBJECT)

	Q.Add(s.ship, graph.Z_UNDER_OBJECT)

	msg := fmt.Sprintf("DIRECTION: %.f\nRANGE: %.f\nWIDE: %.1f",
		Data.NaviData.SonarDir, Data.NaviData.SonarRange, Data.NaviData.SonarWide)
	stats := graph.NewText(msg, s.face, colornames.Palegoldenrod)
	stats.SetPosPivot(graph.ScrP(0.6, 0.1), graph.TopLeft())
	Q.Add(stats, graph.Z_HUD)
	Q.Add(s.caption, graph.Z_STAT_HUD)

	Q.Run(image)
}

func (s *warpScene) procMouseClick(scrPos v2.V2) {
}

func (s *warpScene) OnCommand(command string) {
}

func (*warpScene) Destroy() {
}
