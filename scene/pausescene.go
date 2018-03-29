package scene

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/Shnifer/magellan/graph"
	"github.com/Shnifer/magellan/network"
	"image/color"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
)

type PauseScene struct{
	reason network.PauseReason
	caption *graph.Text

	face font.Face
	centered bool

	getReason func()network.PauseReason
}

func NewPauseScene(face font.Face, getReason func()network.PauseReason) *PauseScene{
	return &PauseScene{
		face:face,
		getReason: getReason,
	}
}

func (p *PauseScene) Init() {
}

func (p *PauseScene) Update(float64) {
	reason :=p.getReason()
	if reason !=p.reason {
		p.reason = reason

		//recalc caption
		var str string
		var color color.Color
		switch {
		case reason.PingLost:
			str="PING LOST!"
			color = colornames.Red
		case !reason.IsFull:
			str="other DISCONNECTED"
			color = colornames.Yellow
		case reason.WantState!=reason.CurState:
			str="Loading new data..."
			color = colornames.Yellowgreen
		case !reason.IsCoherent	:
			str = "waiting other loading..."
			color = colornames.Green
		}
		p.caption = graph.NewText(str,p.face,color)
		p.centered = false
	}
}

func (p *PauseScene) Draw(image *ebiten.Image) {
	if !p.centered {
		w,h:=image.Size()
		p.caption.SetPosPivot(graph.Point{float64(w/2),float64(h/2)}, graph.Center())
		p.centered = true
	}
	p.caption.Draw(image)
}

func (p *PauseScene) Destroy() {
}

