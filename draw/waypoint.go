package draw

import (
	"github.com/Shnifer/magellan/graph"
	"github.com/Shnifer/magellan/v2"
	"github.com/Shnifer/magellan/commons"
	"image/color"
)

type WayPoint struct{
	active bool
	ship v2.V2
	point v2.V2

	showPoint bool

	pointSprite *graph.Sprite
	arrowSprite *graph.Sprite


	cam *graph.Camera
}

func NewWayPoint(cam *graph.Camera, clr color.Color, showPoint bool) *WayPoint{
	point:=NewAtlasSprite(commons.WayPointAN, cam.Deny())
	point.SetSize(30,30)
	arrow:=NewAtlasSprite(commons.WayArrowAN, graph.NoCam)
	arrow.SetSize(50,50)
	point.SetColor(clr)
	arrow.SetColor(clr)
	res:=&WayPoint{
		pointSprite: point,
		arrowSprite: arrow,
		cam: cam,
		showPoint: showPoint,
	}
	return res
}


func (wp *WayPoint) SetActive(active bool){
	wp.active = active
}

func (wp *WayPoint) SetShipPoint(ship, point v2.V2) {
	wp.ship = ship
	wp.point = point
}

func (wp *WayPoint) Req(Q *graph.DrawQueue) {
	if !wp.active{
		return
	}

	if wp.showPoint {
		wp.pointSprite.SetPos(wp.point)
		Q.Add(wp.pointSprite, graph.Z_ABOVE_OBJECT)
	}

	radius := graph.ScrP(0, 0.3).Y
	shipScr:=wp.cam.Apply(wp.ship)
	pointScr:=wp.cam.Apply(wp.point)
	v:=pointScr.Sub(shipScr)
	if v.Len()<radius{
		return
	}
	wp.arrowSprite.SetPos(shipScr.AddMul(v.Normed(),radius))
	wp.arrowSprite.SetAng(v.Dir()+180)
	Q.Add(wp.arrowSprite, graph.Z_HUD)
}