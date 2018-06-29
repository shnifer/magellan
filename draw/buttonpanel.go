package draw

import (
	"github.com/Shnifer/magellan/v2"
	"github.com/Shnifer/magellan/graph"
	"image/color"
	"golang.org/x/image/font"
	"image"
)

type ButtonOpts struct{
	Tex graph.Tex
	Clr color.Color

	Caption string
	Face font.Face
	CapClr color.Color

	//will be returned on events
	Tags string
}

type ButtonsPanelOpts struct {
	PivotP v2.V2
	PivotV v2.V2

	ButtonLayer int
	CaptionLayer int
	ButtonSize v2.V2

	ButtonSpace float64
	BorderSpace float64

	SlideV v2.V2
	SlideT float64
}

type button struct{
	sprite *graph.Sprite
	caption *graph.Text
	tags string
}

type ButtonsPanel struct{
	opts ButtonsPanelOpts
	//top-left
	position v2.V2
	//size of panel
	size v2.V2

	slideK float64

	//active -- we want it to show and ready
	active bool
	//ready -- it is slided and now you may click it
	ready bool

	dirty bool

	buttons []button
}

func NewButtonsPanel(opts ButtonsPanelOpts) *ButtonsPanel{
	if opts.SlideT == 0{
		opts.SlideT = 0
	}

	res := &ButtonsPanel{
		opts: opts,

		dirty: true,

		buttons: []button{},
	}

	return res
}

func (bp *ButtonsPanel) SetActive(active bool) {
	bp.active = active
}

func (bp *ButtonsPanel) Update(dt float64) {
	if bp.dirty{
		bp.dirty = false
		bp.recalc()
	}

	slideD:=dt/bp.opts.SlideT
	if bp.active{
		if bp.slideK<1 {
			bp.slideK += slideD
			if bp.slideK > 1 {
				bp.slideK = 1
			}
			bp.dirty = true
		}
	} else {
		if bp.slideK>0 {
			bp.slideK -= slideD
			if bp.slideK < 0 {
				bp.slideK = 0
			}
			bp.dirty = true
		}
	}

	bp.ready = bp.slideK==1.0
}

func (bp *ButtonsPanel) Req() *graph.DrawQueue{
	if bp.dirty{
		bp.dirty = false
		bp.recalc()
	}

	R:=graph.NewDrawQueue()
	for _,button:=range bp.buttons{
		R.Add(button.sprite, bp.opts.ButtonLayer)
		R.Add(button.caption, bp.opts.CaptionLayer)
	}
	return R
}

func (bp *ButtonsPanel) recalc(){
	butC:=float64(len(bp.buttons))
	spaceC:=butC-1
	if spaceC<0{
		spaceC=0
	}

	bp.size = v2.V2{
		X: bp.opts.ButtonSize.X + 2*bp.opts.BorderSpace,
		Y: butC*bp.opts.ButtonSize.Y + spaceC*bp.opts.BorderSpace + 2*bp.opts.BorderSpace,
	}

	pivotP:=v2.V2{
		X: bp.opts.PivotV.X*bp.size.X,
		Y: bp.opts.PivotV.Y*bp.size.Y,
	}

	bp.position=bp.opts.PivotP.Sub(pivotP)

	butPos:= bp.position.AddMul(v2.V2{X:1, Y:1}, bp.opts.BorderSpace)
	for i, button:=range bp.buttons{
		if i>0{
			butPos.DoAddMul(v2.V2{X:0, Y:1}, bp.opts.ButtonSpace+bp.opts.ButtonSize.Y)
		}

		if bp.ready {
			button.sprite.SetAlpha(1)
		} else {
			button.sprite.SetAlpha(0.5)
		}

		button.sprite.SetPos(butPos)
		button.caption.SetPosPivot(butPos.AddMul(bp.opts.ButtonSize,0.5), graph.Center())
	}
}

//remove old buttons and create new
func (bp *ButtonsPanel) ClearButtons(){
	bp.buttons = []button{}
	bp.dirty = true
}

func (bp *ButtonsPanel) AddButton(opts ButtonOpts){
	sprite :=graph.NewSpriteHUD(opts.Tex)
	sprite.SetSize(bp.opts.ButtonSize.X, bp.opts.ButtonSize.Y)
	sprite.SetColor(opts.Clr)
	sprite.SetPivot(graph.TopLeft())

	caption:=graph.NewText(opts.Caption, opts.Face, opts.CapClr)

	bp.buttons = append(bp.buttons, button{
		sprite: sprite,
		caption: caption,
		tags: opts.Tags,
	})

	bp.dirty = true
}

func (bp *ButtonsPanel) ProcMouse(x,y int) (tags string, ok bool){
	if !bp.active {
		return "", false
	}

	butPos:= bp.position.AddMul(v2.V2{X:1, Y:1}, bp.opts.BorderSpace)
	sx,sy:=int(bp.opts.ButtonSize.X), int(bp.opts.ButtonSize.Y)
	for i, button:=range bp.buttons{
		if i>0{
			butPos.DoAddMul(v2.V2{X:0, Y:1}, bp.opts.ButtonSpace+bp.opts.ButtonSize.Y)
		}
		bx,by:=int(butPos.X), int(butPos.Y)
		r:=image.Rect(bx,by,bx+sx,by+sy)
		if image.Pt(x,y).In(r) {
			return button.tags, true
		}
	}

	return "", false
}