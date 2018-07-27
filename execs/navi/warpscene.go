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
	"math"
	"strings"
)

type warpScene struct {
	ship    *graph.Sprite
	caption *graph.Text
	cam     *graph.Camera

	sonar *graph.Sector

	sonarHUD   *SonarHUD

	sonarText string

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

	sonarSize:=0.8*float64(WinH)

	return &warpScene{
		caption: caption,
		ship:    ship,
		cam:     cam,
		sonar:   sonarSector,
		sonarHUD: NewSonarHUD(graph.ScrP(0.5, 0.5), sonarSize, graph.NoCam, graph.Z_HUD),
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

	var activeSigs[]Signature
	activeSigs, s.sonarText = sonarSigs()
	s.sonarHUD.ActiveSignatures(activeSigs)
	s.sonarHUD.Update(dt)
}

func (s *warpScene) Draw(image *ebiten.Image) {
	defer LogFunc("cosmoScene.Draw")()

	Q := s.q
	Q.Clear()

	Q.Add(s.sonar, graph.Z_UNDER_OBJECT)

	Q.Add(s.ship, graph.Z_UNDER_OBJECT)

	msg := fmt.Sprintf("DIRECTION: %.f\nRANGE: %.f\nWIDE: %.1f",
		Data.NaviData.SonarDir, Data.NaviData.SonarRange, Data.NaviData.SonarWide)
	stats := graph.NewText(msg, Fonts[Face_stats], colornames.Palegoldenrod)
	stats.SetPosPivot(graph.ScrP(0.6, 0.1), graph.TopLeft())
	Q.Add(stats, graph.Z_HUD)
	Q.Add(s.caption, graph.Z_STAT_HUD)

	t:=graph.NewText(s.sonarText, Fonts[Face_cap], colornames.White)
	t.SetPosPivot(graph.ScrP(0.5,0.5),graph.Center())
	Q.Add(t, graph.Z_STAT_HUD)

	Q.Append(s.sonarHUD)

	Q.Run(image)
}

func sonarSigs() ([]Signature, string){
	res:=make([]Signature,0)
	var text string

	ship:=Data.PilotData.Ship.Pos
	range2:=Data.NaviData.SonarRange * Data.NaviData.SonarRange
	for _,p:=range Data.Galaxy.Ordered{
		v:=p.Pos.Sub(ship)
		if v.LenSqr()>range2{
			continue
		}
		angD:=math.Abs(AngDiff(Data.NaviData.SonarDir,v.Dir()))
		if angD>Data.NaviData.SonarWide/2{
			continue
		}
		res = append(res, p.Signatures...)

		for _,s:=range p.BlackBoxes{
			if s!=""{
				text  = text  + "ЧЯ: "+s+"\n"
			}
		}
		for _,s:=range p.Beacons{
			if s!=""{
				text  = text +"Маяк: "+s+"\n"
			}
		}
	}
	strings.TrimRight(text,"\n")
	return res, text
}

func (s *warpScene) procMouseClick(scrPos v2.V2) {
}

func (s *warpScene) OnCommand(command string) {
}

func (*warpScene) Destroy() {
}

func AngDiff(a,b float64) float64 {
	angle:=a-b
	for angle < -180 {
		angle += 360
	}
	for angle >= 180 {
		angle -= 360
	}
	return angle
}