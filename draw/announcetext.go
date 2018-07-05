package draw

import (
	"github.com/Shnifer/magellan/graph"
	"github.com/Shnifer/magellan/v2"
	"golang.org/x/image/font"
	"image/color"
)

type announceMsg struct {
	str string
	clr color.Color
	dur float64
}

type AnnounceText struct {
	pos   v2.V2
	pivot v2.V2
	face  font.Face
	layer int

	t float64

	text *graph.Text

	q []announceMsg
}

func NewAnnounceText(pos v2.V2, pivot v2.V2, face font.Face, layer int) *AnnounceText {
	res := &AnnounceText{
		pos:   pos,
		pivot: pivot,
		face:  face,
		layer: layer,
		q:     make([]announceMsg, 0),
	}
	return res
}

func (at *AnnounceText) AddMsg(msg string, clr color.Color, dur float64) {
	at.q = append(at.q, announceMsg{
		str: msg,
		clr: clr,
		dur: dur,
	})
}

func (at *AnnounceText) Update(dt float64) {
	if len(at.q) == 0 {
		at.text = nil
		return
	}
	msg := at.q[0]
	at.text = graph.NewText(msg.str, at.face, msg.clr)
	at.text.SetPosPivot(at.pos, at.pivot)

	at.t += dt
	if at.t >= at.q[0].dur {
		at.t -= at.q[0].dur
		at.q = at.q[1:]
	}
}

func (at *AnnounceText) Req(Q *graph.DrawQueue) {
	if at.text != nil {
		Q.Add(at.text, at.layer)
	}
}
