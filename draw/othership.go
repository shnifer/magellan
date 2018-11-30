package draw

import (
	"github.com/shnifer/magellan/commons"
	"github.com/shnifer/magellan/graph"
	"github.com/shnifer/magellan/v2"
	"image/color"
)

type OtherShip struct {
	markSprite *graph.Sprite
	sprite     *graph.Sprite
	capText    *graph.Text
	camParams  graph.CamParams
	rb         *commons.RBFollower
}

func NewOtherShip(params graph.CamParams, caption string, elastic float64) *OtherShip {
	sprite := NewAtlasSprite(commons.OtherShipAN, params)
	sprite.SetSize(30, 30)

	markParams := params
	markParams.DenyScale = true
	markParams.DenyAngle = true
	markSprite := NewAtlasSprite(commons.MARKOtherShipAN, markParams)
	sprite.SetSize(30, 30)

	capText := graph.NewText(caption, Fonts[Face_list], color.White)

	rb := commons.NewRBFollower(elastic)

	return &OtherShip{
		sprite:     sprite,
		markSprite: markSprite,
		capText:    capText,
		camParams:  params,
		rb:         rb,
	}
}

func (s *OtherShip) SetRB(rb commons.RBData) {
	s.rb.MoveTo(rb)
}

func (s *OtherShip) Update(dt float64) {
	s.rb.Update(dt)
	ship := s.rb.RB()
	pos := ship.Pos
	s.sprite.SetPosAng(pos, ship.Ang)
	s.markSprite.SetPos(pos)
	if s.camParams.Cam != nil {
		pos = s.camParams.Cam.Apply(pos)
	}
	pos.DoAddMul(v2.V2{X: 0, Y: -30}, graph.GS())
	s.capText.SetPosPivot(pos, graph.Center())
}

func (s *OtherShip) Req(Q *graph.DrawQueue) {
	markAlpha, spriteAlpha := MarkAlpha(commons.ShipSize, s.camParams.Cam)
	if markAlpha > 0 && s.markSprite != nil {
		Q.Add(s.markSprite, graph.Z_ABOVE_OBJECT)
	}

	if spriteAlpha > 0 && s.sprite != nil {
		Q.Add(s.sprite, graph.Z_ABOVE_OBJECT)
	}
	Q.Add(s.capText, graph.Z_HUD)
}
