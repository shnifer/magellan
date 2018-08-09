package main

import (
	. "github.com/Shnifer/magellan/commons"
	. "github.com/Shnifer/magellan/draw"
	"github.com/Shnifer/magellan/graph"
	. "github.com/Shnifer/magellan/log"
	"github.com/Shnifer/magellan/ranma"
	"github.com/Shnifer/magellan/v2"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"log"
	"time"
)

type engiScene struct {
	shipID string

	ranma      *ranma.Ranma
	background *graph.Sprite

	systemsMonitor *systemsMonitor

	q *graph.DrawQueue

	tick <-chan time.Time

	wormOut string

	local localCounters
}

func newEngiScene() *engiScene {
	back := NewAtlasSpriteHUD(EngiBackgroundAN)
	back.SetSize(float64(WinW), float64(WinH))
	back.SetPivot(graph.TopLeft())

	//ranma:= ranma.NewRanma(DEFVAL.RanmaAddr, DEFVAL.DropOnRepair, DEFVAL.RanmaTimeoutMs, DEFVAL.RanmaHistoryDepth)
	ranma := &ranma.Ranma{}
	return &engiScene{
		ranma:          ranma,
		background:     back,
		systemsMonitor: newSystemsMonitor(),
		q:              graph.NewDrawQueue(),
		tick:           time.Tick(time.Second),
	}
}

func (s *engiScene) Init() {
	defer LogFunc("engiScene.Init")()

	if s.shipID == Data.ShipID {
		return
	}
	s.shipID = Data.ShipID

	s.local = initLocal()
	initMedi(Data.ShipID)

	for sysN := 0; sysN < SysCount; sysN++ {
		if s.ranma.GetIn(sysN) != Data.EngiData.InV[sysN] {
			s.ranma.SetIn(sysN, Data.EngiData.InV[sysN])
		}
	}

	s.wormOut = ""
}

func (s *engiScene) Update(dt float64) {
	defer LogFunc("engiScene.Update")()

	Data.Galaxy.Update(Data.PilotData.SessionTime)

	x, y := ebiten.CursorPosition()
	mouse := v2.V2{X: float64(x), Y: float64(y)}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		sysN, ok := s.systemsMonitor.mouseOverSystem(mouse)
		if ok {
			s.showSystemInfo(sysN)
		}
	}

	Data.EngiData.Emissions = CalculateEmissions(Data.Galaxy, Data.PilotData.Ship.Pos)
	Data.EngiData.BSPDegrade = CalculateBSPDegrade(s.ranma)
	CalculateCounters(dt)
	s.CalculateLocalCounters()

	select {
	case <-s.tick:
		s.procTick()
	default:
	}

	s.checkForWormHole()

	s.systemsMonitor.update(dt, s.ranma)
}

func (s *engiScene) Draw(image *ebiten.Image) {
	defer LogFunc("engiScene.Draw")()
	Q := s.q
	Q.Clear()

	Q.Append(s.systemsMonitor)

	Q.Run(image)
}

func (s *engiScene) OnCommand(command string) {
}

func (*engiScene) Destroy() {
}

func (s *engiScene) showSystemInfo(n int) {
	log.Println("show system info #", n)
}

func (s *engiScene) procTick() {
	s.checkDamage()
	s.checkMedicine()
}

func (s *engiScene) checkForWormHole() {
	if Data.EngiData.Emissions[EMI_WORMHOLE] == 0 {
		return
	}

	target, err := GetWormHoleTarget(Data.State.GalaxyID)
	if err != nil {
		Log(LVL_ERROR, err)
		return
	}

	if target == WarmHoleYouDIE {
		Log(LVL_INFO, "Die by wormhole")
		return
	}

	//to other system
	state := Data.State
	state.StateID = STATE_cosmo
	state.GalaxyID = target
	Client.RequestNewState(state.Encode(), false)
}
