package draw

import (
	"github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/graph"
	"github.com/Shnifer/magellan/v2"
	"image/color"
)

type OtherShip struct {
	sprite    *graph.Sprite
	capText   *graph.Text
	camParams graph.CamParams
	rb        *commons.RBFollower
}

func NewOtherShip(params graph.CamParams, caption string, elastic float64) *OtherShip {
	sprite := NewAtlasSprite(commons.OtherShipAN, params)
	sprite.SetSize(30, 30)

	capText := graph.NewText(caption, Fonts[Face_list], color.White)

	rb := commons.NewRBFollower(elastic)

	return &OtherShip{
		sprite:    sprite,
		capText:   capText,
		camParams: params,
		rb:        rb,
	}
}

func (s *OtherShip) SetRB(rb commons.RBData) {
	s.rb.MoveTo(rb)
}

func (s *OtherShip) Update(dt float64) {
	s.rb.Update(dt)
	ship := s.rb.RB()
	s.sprite.SetPosAng(ship.Pos, ship.Ang)
	pos := ship.Pos
	if s.camParams.Cam != nil {
		pos = s.camParams.Cam.Apply(pos)
	}
	pos.DoAddMul(v2.V2{X: 0, Y: -30}, graph.GS())
	s.capText.SetPosPivot(pos, graph.Center())
}

func (s *OtherShip) Req() *graph.DrawQueue {
	R := graph.NewDrawQueue()
	R.Add(s.sprite, graph.Z_ABOVE_OBJECT)
	R.Add(s.capText, graph.Z_HUD)
	return R
}
