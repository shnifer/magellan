package main

import (
	. "github.com/Shnifer/magellan/commons"
	. "github.com/Shnifer/magellan/draw"
	"github.com/Shnifer/magellan/graph"
	. "github.com/Shnifer/magellan/log"
	"github.com/Shnifer/magellan/v2"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"strconv"
	"time"
)

type LoginScene struct {
	face     font.Face
	question *graph.Text
	text     *graph.Text

	lastErrTime time.Time
	errorMsg    *graph.Text

	inputText string

	backT float64
	dir   int

	sonarHUD   *SonarHUD
	activeSigs [10]bool
}

func NewLoginScene() *LoginScene {
	const questionText = "Enter login ID (\"firefly\")"
	const errorText = "Wrong ID!"

	face := Fonts[Face_cap]

	question := graph.NewText(questionText, face, colornames.Yellowgreen)
	question.SetPosPivot(graph.ScrP(0.5, 0.3), graph.Center())

	errorMsg := graph.NewText(errorText, face, colornames.Indianred)
	errorMsg.SetPosPivot(graph.ScrP(0.5, 0.7), graph.Center())

	cam := graph.NewCamera()
	cam.Center = graph.ScrP(0.5, 0.5)
	cam.Scale = cam.Center.Y * 0.8
	cam.Recalc()

	return &LoginScene{
		face:     face,
		question: question,
		errorMsg: errorMsg,
		dir:      1,
		sonarHUD: NewSonarHUD(graph.ScrP(0.3, 0.7), 500, graph.NoCam, graph.Z_HUD),
	}
}

func (p *LoginScene) Init() {
	defer LogFunc("LoginScene.Init")()
	p.inputText = ""
	p.text = nil
	p.lastErrTime = time.Time{}
}

func (p *LoginScene) Update(dt float64) {
	defer LogFunc("LoginScene.Update")()
	var changed bool

	input := ebiten.InputChars()
	if len(input) > 0 {
		changed = true
		p.inputText = string(append([]rune(p.inputText), input...))
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		if len(p.inputText) > 0 {
			changed = true
			p.inputText = p.inputText[0 : len(p.inputText)-1]
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		p.tryToStartFly()
	}

	if changed || p.text == nil {
		p.text = graph.NewText(p.inputText, p.face, colornames.White)
		p.text.SetPosPivot(graph.ScrP(0.5, 0.5), graph.Center())
	}

	for i := 0; i < 10; i++ {
		if inpututil.IsKeyJustPressed(ebiten.Key0 + ebiten.Key(i)) {
			p.activeSigs[i] = !p.activeSigs[i]
		}
	}

	activeSigs := make([]Signature, 0, 10)
	for i, active := range p.activeSigs {
		if active {
			activeSigs = append(activeSigs, Signature{
				Dev:      v2.ZV,
				TypeName: strconv.Itoa(i)})
		}
	}
	p.sonarHUD.ActiveSignatures(activeSigs)
	p.sonarHUD.Update(dt)
}

func (p *LoginScene) Draw(image *ebiten.Image) {
	defer LogFunc("LoginScene.Draw")()

	const ErrorShowtime = time.Second * 2

	Q := graph.NewDrawQueue()

	Q.Add(p.question, graph.Z_HUD)
	Q.Add(p.text, graph.Z_HUD)

	errTime := time.Since(p.lastErrTime)
	if errTime < ErrorShowtime {
		if int(errTime.Seconds()*8)%2 == 0 {
			Q.Add(p.errorMsg, graph.Z_HUD)
		}
	}

	Q.Append(p.sonarHUD)

	Q.Run(image)
}

func (p *LoginScene) Destroy() {
}

func (p *LoginScene) tryToStartFly() {
	defer LogFunc("LoginScene.tryToStartFly")()

	state := State{
		StateID:  STATE_cosmo,
		ShipID:   p.inputText,
		GalaxyID: START_Galaxy_ID,
	}.Encode()
	Client.RequestNewState(state, false)
}

func (p *LoginScene) OnCommand(command string) {
	defer LogFunc("LoginScene.OnCommand")()
	if command == CMD_STATECHANGEFAIL {
		p.lastErrTime = time.Now()
	}
}
