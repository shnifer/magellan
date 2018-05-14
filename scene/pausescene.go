package scene

import (
	"github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/graph"
	"github.com/Shnifer/magellan/network"
	"github.com/hajimehoshi/ebiten"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"image/color"
)

type PauseScene struct {
	reason  network.PauseReason
	caption *graph.Text

	face font.Face

	getReason func() network.PauseReason
}

func NewPauseScene(face font.Face, getReason func() network.PauseReason) *PauseScene {
	return &PauseScene{
		face:      face,
		getReason: getReason,
		caption:   graph.NewText("", face, color.White),
	}
}

func (p *PauseScene) Init() {
}

func (p *PauseScene) Update(float64) {
	commons.LogFunc("pauseScene.Update")()

	reason := p.getReason()
	if reason != p.reason {
		p.reason = reason

		//recalc caption
		var str string
		var captionColor color.Color
		switch {
		case reason.PingLost:
			str = "PING LOST!"
			captionColor = colornames.Red
		case !reason.IsFull:
			str = "other DISCONNECTED"
			captionColor = colornames.Yellow
		case reason.WantState != reason.CurState:
			str = "Loading new data..."
			captionColor = colornames.Yellowgreen
		case !reason.IsCoherent:
			str = "waiting other loading..."
			captionColor = colornames.Green
		}
		p.caption = graph.NewText(str, p.face, captionColor)
		p.caption.SetPosPivot(graph.ScrP(0.5, 0.5), graph.Center())
	}
}

func (p *PauseScene) Draw(image *ebiten.Image) {
	commons.LogFunc("pauseScene.Draw")()

	p.caption.Draw(image)
}

func (p *PauseScene) Destroy() {
}

func (PauseScene) OnCommand(command string) {
}
