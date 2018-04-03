package main

import (
	"github.com/Shnifer/magellan/graph"
	"github.com/hajimehoshi/ebiten"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"log"
	"github.com/hajimehoshi/ebiten/inpututil"
	"time"
	. "github.com/Shnifer/magellan/commons"
)

type LoginScene struct {
	face     font.Face
	question *graph.Text
	text     *graph.Text

	lastErrTime time.Time
	errorMsg *graph.Text

	inputText string
}

func NewLoginScene() *LoginScene {
	const questionText = "Enter login ID:"
	const errorText = "Wrong ID!"

	face:=fonts[face_cap]

	question := graph.NewText(questionText, face, colornames.Yellowgreen)
	question.SetPosPivot(graph.ScrP(0.5, 0.3), graph.Center())

	errorMsg := graph.NewText(errorText, face, colornames.Indianred)
	errorMsg.SetPosPivot(graph.ScrP(0.5, 0.7), graph.Center())
	return &LoginScene{
		face: face,
		question: question,
		errorMsg: errorMsg,
	}
}

func (p *LoginScene) Init() {
	log.Println("Loginscene inited")
	p.inputText = ""
	p.lastErrTime = time.Time{}
}

func (p *LoginScene) Update(float64) {
	var changed bool

	input := ebiten.InputChars()
	if len(input)>0 {
		changed = true
		p.inputText = string(append([]rune(p.inputText), input...))
	}


	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		if len(p.inputText) > 0 {
			changed = true
			p.inputText = p.inputText[0:len(p.inputText)-1]
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		p.tryToStartFly()
	}

	if changed || p.text == nil{
		p.text = graph.NewText(p.inputText, p.face, colornames.White)
		p.text.SetPosPivot(graph.ScrP(0.5, 0.5), graph.Center())
	}
}

func (p *LoginScene) Draw(image *ebiten.Image) {
	const ErrorShowtime = time.Second * 2

	p.question.Draw(image)
	p.text.Draw(image)

	errTime:=time.Since(p.lastErrTime)
	if errTime<ErrorShowtime {
		if int(errTime.Seconds()*4)%2 == 0 {
			p.errorMsg.Draw(image)
		}
	}
}

func (p *LoginScene) Destroy() {
}

func (p *LoginScene) tryToStartFly() {
	state :=State{
		Special:  STATE_cosmo,
		ShipID:   p.inputText,
		GalaxyID: START_Galaxy_ID,
	}.Encode()

	err:=Client.RequestNewState(state)
	if err!=nil{
		log.Println(err)
		p.lastErrTime = time.Now()
	}
}