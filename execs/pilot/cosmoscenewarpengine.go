package main

import (
	"fmt"
	. "github.com/shnifer/magellan/commons"
	"github.com/shnifer/magellan/draw"
	"github.com/shnifer/magellan/graph"
	"github.com/shnifer/magellan/input"
	"github.com/shnifer/magellan/v2"
	"golang.org/x/image/colornames"
	"image/color"
)

//in sec
const TimeToWarp = 3

type cosmoSceneWarpEngine struct {
	gravityAcc v2.V2

	wasReset bool
	toWarpT  float64
	fired    bool
}

func newCosmoSceneWarpEngine() *cosmoSceneWarpEngine {
	return &cosmoSceneWarpEngine{}
}

func (h *cosmoSceneWarpEngine) update(dt float64) {
	v := input.WarpLevel("warpspeed")
	if v <= 0 {
		h.wasReset = true
		h.fired = false
		h.toWarpT = 0
	}
	if v > 0 && h.wasReset {
		h.toWarpT += dt
	}
	if h.toWarpT > TimeToWarp && !h.fired {
		h.fired = true
		toWarp()
	}
}

func toWarp() {
	if Data.EngiData.Counters.Fuel < Data.SP.Warp_engine.Warp_enter_consumption {
		return
	}

	Data.EngiData.Counters.Fuel -= Data.SP.Warp_engine.Warp_enter_consumption

	state := Data.State
	state.StateID = STATE_warp
	state.GalaxyID = WARP_Galaxy_ID
	Client.RequestNewState(state.Encode(), false)
}

func (h *cosmoSceneWarpEngine) Req(Q *graph.DrawQueue) {
	basePoint := graph.ScrP(0.8, 0.2)

	courseMsg := fmt.Sprintf("Course: %.1f", 360-Data.PilotData.Ship.Vel.Dir())
	courseText := graph.NewText(courseMsg, draw.Fonts[draw.Face_mono], color.White)
	courseText.SetPosPivot(basePoint, graph.TopLeft())
	Q.Add(courseText, graph.Z_STAT_HUD+2)

	_, h_int := courseText.GetSize()
	interV := v2.V2{X: 0, Y: float64(h_int) * 1.4}

	velMsg := fmt.Sprintf("Velocity: %f", Data.PilotData.Ship.Vel.Len())
	velText := graph.NewText(velMsg, draw.Fonts[draw.Face_mono], color.White)
	velText.SetPosPivot(basePoint.AddMul(interV, 1), graph.TopLeft())
	Q.Add(velText, graph.Z_STAT_HUD+2)

	gravAcc := h.gravityAcc.Len() * 100
	gravityMsg := fmt.Sprintf("Gravity: %.1f%%", gravAcc)
	var gravityColor color.Color
	switch {
	case gravAcc > 25:
		gravityColor = colornames.Red
	case gravAcc > 10:
		gravityColor = colornames.Yellow
	case gravAcc > 2:
		gravityColor = colornames.Yellowgreen
	default:
		gravityColor = colornames.Lightgreen
	}

	gravText := graph.NewText(gravityMsg, draw.Fonts[draw.Face_mono], gravityColor)
	gravText.SetPosPivot(basePoint.AddMul(interV, 2), graph.TopLeft())
	Q.Add(gravText, graph.Z_STAT_HUD+2)

	warpMsg := ""
	var warpColor color.Color
	warpColor = color.White
	switch {
	case !h.wasReset:
		warpMsg = "RESET to zero"
		warpColor = colornames.Red
	case h.fired:
		warpMsg = "FIRE!!!"
	default:
		warpMsg = fmt.Sprintf("warping: %.0f%%", h.toWarpT/TimeToWarp*100)
	}
	warpText := graph.NewText(warpMsg, draw.Fonts[draw.Face_mono], warpColor)
	warpText.SetPosPivot(basePoint.AddMul(interV, 3), graph.TopLeft())
	Q.Add(warpText, graph.Z_STAT_HUD+2)
}
