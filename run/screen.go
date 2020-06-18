package run

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	WinSizeW = 1900
	WinSizeH = 1000

	FPS = 30
)

var WinSize = pixel.V(WinSizeW, WinSizeH)

type Screen struct {
	*pixelgl.Window
	ZoomPos *pixel.Vec
	Zoom    float64
}

func (Scr *Screen) GetPos(Pos pixel.Vec) pixel.Vec {
	return Pos.Sub(*Scr.ZoomPos).Scaled(1 / Scr.Zoom)
}

func (Scr *Screen) ChangeZoom(FocusedPoint pixel.Vec, NewZoom float64) {
	Pos := Scr.ZoomPos.Add(FocusedPoint.Scaled(Scr.Zoom - NewZoom))
	Scr.ZoomPos.X = Pos.X
	Scr.ZoomPos.Y = Pos.Y
	Scr.Zoom = NewZoom
}
