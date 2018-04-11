package scene

import (
	"github.com/Shnifer/magellan/graph"
	"github.com/hajimehoshi/ebiten"
	"golang.org/x/image/font"
	"image/color"
)

type captionScene struct {
	caption *graph.Text

	face      font.Face
	textColor color.Color
	text      string

	getText func() string
}

func NewCaptionSceneString(face font.Face, textColor color.Color, text string) *captionScene {
	return &captionScene{
		face:      face,
		text:      text,
		textColor: textColor,
		getText:   nil,
	}
}

func NewCaptionSceneFunc(face font.Face, textColor color.Color, getText func() string) *captionScene {
	return &captionScene{
		face:      face,
		textColor: textColor,
		getText:   getText,
	}
}

func (s *captionScene) Init() {
}

func (s *captionScene) Update(float64) {
	var needRedraw bool
	if s.getText != nil {
		newText := s.getText()
		if newText != s.text {
			s.text = newText
			needRedraw = true
		}
	}

	if s.caption == nil {
		needRedraw = true
	}

	if !needRedraw {
		return
	}

	s.caption = graph.NewText(s.text, s.face, s.textColor)
	s.caption.SetPosPivot(graph.ScrP(0.5, 0.5), graph.Center())
}

func (s *captionScene) Draw(image *ebiten.Image) {
	s.caption.Draw(image)
}

func (s *captionScene) Destroy() {
}

func (*captionScene) OnCommand(command string) {
}
