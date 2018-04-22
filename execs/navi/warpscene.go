package main

import (
	"fmt"
	. "github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/graph"
	"github.com/Shnifer/magellan/input"
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
}

func newWarpScene() *warpScene {
	caption := graph.NewText("Navi warp scene", Fonts[Face_cap], colornames.Aliceblue)
	caption.SetPosPivot(graph.ScrP(0.1, 0.1), graph.TopLeft())

	cam := graph.NewCamera()
	cam.Center = graph.ScrP(0.5, 0.5)
	cam.Recalc()

	ship := graph.NewSprite(GetAtlasTex("ship"), cam, false, false)
	ship.SetSize(100, 100)
	ship.SetAlpha(0.5)

	sonarSector := graph.NewSector(cam, false, false)
	sonarSector.SetColor(colornames.Indigo)

	return &warpScene{
		caption: caption,
		ship:    ship,
		cam:     cam,
		face:    Fonts[Face_stats],
		sonar:   sonarSector,
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

	Data.NaviData.SonarDir += turn * dt * Data.BSP.Sonar_rotate_speed
	Data.NaviData.SonarRange += rang * dt * Data.BSP.Sonar_range_change
	Data.NaviData.SonarRange = Clamp(Data.NaviData.SonarRange, Data.BSP.Sonar_range_min, Data.BSP.Sonar_range_max)
	Data.NaviData.SonarWide += wide * dt * Data.BSP.Sonar_angle_change
	Data.NaviData.SonarWide = Clamp(Data.NaviData.SonarWide, Data.BSP.Sonar_angle_min, Data.BSP.Sonar_angle_max)

	//PilotData Rigid Body emulation
	Data.PilotData.Ship = Data.PilotData.Ship.Extrapolate(dt)

	s.sonar.SetRadius(Data.NaviData.SonarRange)
	s.sonar.SetAngles(
		Data.NaviData.SonarDir-Data.NaviData.SonarWide/2,
		Data.NaviData.SonarDir+Data.NaviData.SonarWide/2)
}

func (s *warpScene) Draw(image *ebiten.Image) {
	defer LogFunc("cosmoScene.Draw")()

	s.sonar.Draw(image)

	s.ship.SetAng(Data.PilotData.Ship.Ang)
	s.ship.Draw(image)

	msg := fmt.Sprintf("DIRECTION: %.f\nRANGE: %.f\nWIDE: %.1f",
		Data.NaviData.SonarDir, Data.NaviData.SonarRange, Data.NaviData.SonarWide)
	stats := graph.NewText(msg, s.face, colornames.Palegoldenrod)
	stats.SetPosPivot(graph.ScrP(0.6, 0.1), graph.TopLeft())
	stats.Draw(image)
	s.caption.Draw(image)
}

func (s *warpScene) procMouseClick(scrPos v2.V2) {
}

func (s *warpScene) OnCommand(command string) {
}

func (*warpScene) Destroy() {
}
