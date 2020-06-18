package run

import (
	"github.com/faiface/pixel"
	"math"
)

type Star struct {
	Pos    pixel.Vec
	Speed  pixel.Vec
	Mass   float64
	Radius float64
	Color  pixel.RGBA
}

var (
	Stars    []*Star
	StarsLen = 0

	DelStars []int
	Collides [][2]int

	ColorOfRadiusKeys   = [6]float64{0, 7, 10, 20, 90, 100}
	ColorOfRadiusValues = [6]pixel.RGBA{pixel.RGB(0, 0, 0), pixel.RGB(1, 0.7, 0.2), pixel.RGB(0.6, 0.6, 0.), pixel.RGB(0.3, 1, 0.8), pixel.RGB(1, 0.2, 0.2), pixel.RGB(1, 1, 0.5)}
)

func (S *Star) CheckRadius() {
	S.Radius = math.Pow(S.Mass, 1./3) / 5
}

func (S *Star) CheckColor() {
	for i := 1; i < len(ColorOfRadiusKeys); i++ {
		if S.Radius < ColorOfRadiusKeys[i] {
			x := (S.Radius - ColorOfRadiusKeys[i-1]) / (ColorOfRadiusKeys[i] - ColorOfRadiusKeys[i-1])
			S.Color = ColorOfRadiusValues[i-1].Scaled(1 - x).Add(ColorOfRadiusValues[i].Scaled(x))
			return
		}
	}
	S.Color = ColorOfRadiusValues[len(ColorOfRadiusValues)-1]
}

func NewStarVec(Pos pixel.Vec, Speed pixel.Vec, Mass float64) {
	S := Star{
		Pos:   Pos,
		Speed: Speed,
		Mass:  Mass,
	}
	S.CheckRadius()
	S.CheckColor()
	Stars = append(Stars, &S)
	StarsLen++
}

func NewStar(PosX, PosY, SpeedX, SpeedY, Mass float64) {
	NewStarVec(pixel.V(PosX, PosY), pixel.V(SpeedX, SpeedY), Mass)
}

func DelStar(i int) {
	if Focused >= i {
		Focused--
	}
	Stars = append(Stars[:i], Stars[i+1:]...)
	StarsLen--
}

func Collide(j, i int) {

	if Focused == j {
		Focused = i
	}

	//if i == 0 {
	//	GenerateNewMeteor(MinR, MaxR, Stars[j].Mass)
	//	DelStars = append(DelStars, j)
	//	return
	//} else if j == 0 {
	//	if Focused == i {
	//		Focused = 0
	//	}
	//	GenerateNewMeteor(MinR, MaxR, Stars[i].Mass)
	//	DelStars = append(DelStars, i)
	//	return
	//}
	for _, k := range DelStars {
		if i == k || j == k {
			return
		}
	}

	Stars[i].Pos = Stars[i].Pos.Scaled(Stars[i].Mass).Add(Stars[j].Pos.Scaled(Stars[j].Mass)).Scaled(1 / (Stars[i].Mass + Stars[j].Mass))
	Stars[i].Speed = Stars[i].Speed.Scaled(Stars[i].Mass).Add(Stars[j].Speed.Scaled(Stars[j].Mass)).Scaled(1 / (Stars[i].Mass + Stars[j].Mass))
	Stars[i].Mass += Stars[j].Mass
	Stars[i].CheckRadius()
	Stars[i].CheckColor()

	Stars[j].Mass = 0
	Stars[j].Pos = pixel.V(0, 0)

	DelStars = append(DelStars, j)
}
