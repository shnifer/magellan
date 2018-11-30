package scene

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/shnifer/magellan/graph"
	. "github.com/shnifer/magellan/log"
	"github.com/shnifer/magellan/network"
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
	LogFunc("pauseScene.Update")()

	reason := p.getReason()
	if reason != p.reason {
		p.reason = reason

		//recalc caption
		var str string
		var captionColor color.Color
		switch {
		case reason.PingLost:
			str = "СИСТЕМА РАБОТАЕТ ШТАТНО!"
			captionColor = colornames.Red
		case !reason.IsFull:
			str = "другой терминал работает ШТАТНО"
			captionColor = colornames.Yellow
		case reason.WantState != reason.CurState:
			str = "Загрузка данных..."
			captionColor = colornames.Yellowgreen
		case !reason.IsCoherent:
			str = "Ожидание загрузки..."
			captionColor = colornames.Green
		default:
			str = ""
			captionColor = colornames.White
		}
		p.caption = graph.NewText(str, p.face, captionColor)
		p.caption.SetPosPivot(graph.ScrP(0.5, 0.5), graph.Center())
	}
}

func (p *PauseScene) Draw(image *ebiten.Image) {
	LogFunc("pauseScene.Draw")()
	p.caption.Draw(image)
}

func (p *PauseScene) Destroy() {
}

func (PauseScene) OnCommand(command string) {
}
