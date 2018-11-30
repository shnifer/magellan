package main

//COPYPASTE: in Pilot

import (
	. "github.com/shnifer/magellan/commons"
	. "github.com/shnifer/magellan/draw"
	"github.com/shnifer/magellan/graph"
	"github.com/shnifer/magellan/v2"
	"golang.org/x/image/colornames"
)

type predictors struct {
	show            bool
	predictorZero   *TrackPredictor
	predictorThrust *TrackPredictor
}

func (p *predictors) init(cam *graph.Camera) {
	gps := NewGravityPredictorSource(Data.Galaxy, 0.1, 600)

	predictorSprite := NewAtlasSprite(PredictorAN, cam.Deny())
	predictorSprite.SetSize(20, 20)
	opts := TrackPredictorOpts{
		Cam:      cam,
		Sprite:   predictorSprite,
		Clr:      colornames.Palevioletred,
		Layer:    graph.Z_ABOVE_OBJECT + 1,
		GPS:      gps,
		UpdT:     DEFVAL.CosmoPredictorUpdT,
		NumInSec: DEFVAL.CosmoPredictorNumInSec,
		GravEach: DEFVAL.CosmoPredictorGravEach,
		TrackLen: DEFVAL.CosmoPredictorTrackLen,
		DrawMaxP: DEFVAL.CosmoPredictorDrawMaxP,
	}

	p.predictorThrust = NewTrackPredictor(opts)

	predictor2Sprite := NewAtlasSprite(PredictorAN, cam.Deny())
	predictor2Sprite.SetSize(15, 15)
	predictor2Sprite.SetColor(colornames.Darkgray)
	opts.Sprite = predictor2Sprite
	opts.Clr = colornames.Cadetblue
	p.predictorZero = NewTrackPredictor(opts)

	p.show = true
}

func (p predictors) setParams(ss float64, ship RBData) {
	p.predictorThrust.SetAccelSessionTimeShipPos(Data.PilotData.ThrustVector, ss, ship)
	p.predictorZero.SetAccelSessionTimeShipPos(v2.ZV, ss, ship)
}

func (p predictors) Req(Q *graph.DrawQueue) {
	if p.show {
		Q.Append(p.predictorZero)
		Q.Append(p.predictorThrust)
	}
}
