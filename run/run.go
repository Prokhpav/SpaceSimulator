package run

import "C"
import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"math"
	"math/rand"
	"time"
)

const (
	G = 10 // Gravity constant. Can change speed of simulation

	NumOfStars = 1500
	MinR       = 500   // Minimal random distance to zero position
	MaxR       = 2000  // Maximal random distance to zero position
	MinSpeed   = -100  // Minimal random speed
	MaxSpeed   = 100   // Maximal random speed
	MinMass    = 2000  // Minimal star mass
	MaxMass    = 20000 // Maximal star mass

	MaxDist = 20000

	MaxCloudR = 100

	CameraMovementSpeed = 5
)

var (
	MinStars = 0 // Minimal number of stars

	CloudRB = 4 * MaxCloudR / ColorOfRadiusKeys[2]
	CloudRA = CloudRB / -ColorOfRadiusKeys[2]

	Focused = -1

	SumMass = 0.

	MeteorOrbitMul = 1.
)

func GenerateNewMeteor(MinR, MaxR, Mass float64) {
	angle := rand.Float64()
	R := rand.Float64()
	if R < 0.5 && rand.Float64() < 1-R*2 {
		R += 0.5
	}
	R = (R*(MaxR-MinR) + MinR) * MeteorOrbitMul
	NewStarVec(
		pixel.V(math.Cos(angle*2*math.Pi), math.Sin(angle*2*math.Pi)).Scaled(R),
		randPolarAngle(math.Sqrt(G*SumMass/R)+MinSpeed, math.Sqrt(G*SumMass/R)+MaxSpeed, angle+0.24, angle+0.26).Scaled(0.185),
		Mass)
}

func CheckCenter() {
	Pos := pixel.V(0, 0)
	V := pixel.V(0, 0)
	R := 0.
	for _, S := range Stars {
		Pos = Pos.Add(S.Pos.Scaled(S.Mass))
		V = V.Add(S.Speed.Scaled(S.Mass))
		R += S.Mass
	}
	Pos = Pos.Scaled(1. / R)
	V = V.Scaled(1. / R)
	for i := 0; i < len(Stars); i++ {
		Stars[i].Pos = Stars[i].Pos.Sub(Pos)
		Stars[i].Speed = Stars[i].Speed.Sub(V)
	}
	SumMass = R
}

func Run() {
	var x float64
	var R float64
	var Pos pixel.Vec

	rand.Seed(time.Now().UnixNano())

	NewStar(0, 0, 0, 0, 1000000000) // Generating a new star in centre

	CheckCenter()

	//Stars[0].Radius = math.Pow(Stars[0].Mass, 1./3) / 10
	for i := 0; i < NumOfStars; i++ {
		GenerateNewMeteor(MinR, MaxR, randInterval(MinMass, MaxMass))
		//angle := randInterval(0, 1)
		//R := rand.Float64()
		//if R < 0.5 && rand.Float64() < 1-R*2 {
		//	R += 0.5
		//}
		//R = R*(MaxR-MinR) + MinR
		//NewStarVec(
		//	pixel.V(math.Cos(angle*2*math.Pi), math.Sin(angle*2*math.Pi)).Scaled(R),
		//	randPolarAngle(math.Sqrt(G*Stars[0].Mass/R), math.Sqrt(G*Stars[0].Mass/R), angle+0.24, angle+0.26).Scaled(0.185).Add(randPolar(MinSpeed, MaxSpeed)),
		//randPolar(MinSpeed, MaxSpeed),
		//randInterval(MinMass, MaxMass))
	}

	CheckCenter()

	TickTimer := time.Tick(time.Second / 120)
	imd := imdraw.New(nil)

	_win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Title:  "Snake",
		Bounds: pixel.R(0, 0, float64(WinSizeW), float64(WinSizeH)),
	})
	if err != nil {
		panic(err)
	}
	win := Screen{Window: _win, Zoom: 8}
	win.ZoomPos = &pixel.Vec{X: -WinSizeW * win.Zoom / 2, Y: -WinSizeH * win.Zoom / 2}

	Paused := true
	Tracers := false

	for !win.Closed() {
		dt := 1. / FPS

		if win.MouseScroll().Y > 0 {
			win.ChangeZoom(win.MousePosition(), win.Zoom/1.1)
		} else if win.MouseScroll().Y < 0 {
			win.ChangeZoom(win.MousePosition(), win.Zoom*1.1)
		}
		if win.JustPressed(pixelgl.KeySpace) {
			Paused = !Paused
		}
		if win.JustPressed(pixelgl.KeyT) {
			Tracers = !Tracers
		}
		if win.JustPressed(pixelgl.KeyF) {
			if Focused != -1 {
				Focused = -1
			} else {
				Pos = win.MousePosition().Scaled(win.Zoom).Add(*win.ZoomPos)

				MinR := Dist(Pos, Stars[0].Pos)
				MinI := 0
				for i := 1; i < len(Stars); i++ {
					R := Dist(Pos, Stars[i].Pos)
					if R < MinR {
						MinR = R
						MinI = i
					}
				}
				Focused = MinI
			}
		}
		if win.JustPressed(pixelgl.Key0) {
			Focused = -2
		}
		if win.Pressed(pixelgl.KeyLeft) {
			win.ZoomPos.X -= CameraMovementSpeed * win.Zoom
		}
		if win.Pressed(pixelgl.KeyRight) {
			win.ZoomPos.X += CameraMovementSpeed * win.Zoom
		}
		if win.Pressed(pixelgl.KeyUp) {
			win.ZoomPos.Y += CameraMovementSpeed * win.Zoom
		}
		if win.Pressed(pixelgl.KeyDown) {
			win.ZoomPos.Y -= CameraMovementSpeed * win.Zoom
		}
		if win.JustPressed(pixelgl.KeyC) {
			CheckCenter()
		}
		if win.JustPressed(pixelgl.Key1) {
			MinStars = MinStars * 10 / 11
		}
		if win.JustPressed(pixelgl.Key2) {
			MinStars = MinStars * 11 / 10
		}
		if win.JustPressed(pixelgl.Key3) {
			MinStars = int(Max(Min(float64(MinStars*2), 2500), 1))
		}
		if win.JustPressed(pixelgl.Key4) {
			MinStars = 0
		}
		if win.JustPressed(pixelgl.Key5) {
			MeteorOrbitMul = Max(1, MeteorOrbitMul-1)
		}
		if win.JustPressed(pixelgl.Key6) {
			MeteorOrbitMul += 1
		}

		//////////////////////////////////////////////////////////////// calculation

		if !Paused {
			Pos = pixel.V(0, 0)
			R := 0.
			for _, S := range Stars {
				Pos = Pos.Add(S.Pos.Scaled(S.Mass))
				R += S.Mass
			}
			Pos = Pos.Scaled(1. / R)

			if Focused == -2 {
				win.ZoomPos.X = Pos.X - WinSizeW*win.Zoom/2
				win.ZoomPos.Y = Pos.Y - WinSizeH*win.Zoom/2
			}

			Gdt2 := G * dt * dt
			for i := 0; i < len(Stars); i++ {
				for j := i + 1; j < len(Stars); j++ {
					R = Dist(Stars[i].Pos, Stars[j].Pos)
					if R >= (Stars[i].Radius+Stars[j].Radius)*0.9 {
						x = Gdt2 / (R * R * R)
						if x*Stars[j].Mass > 0.000001 {
							Stars[i].Speed = Stars[i].Speed.Add(Stars[j].Pos.Sub(Stars[i].Pos).Scaled(x * Stars[j].Mass))
						}
						if x*Stars[i].Mass > 0.000001 {
							Stars[j].Speed = Stars[j].Speed.Add(Stars[i].Pos.Sub(Stars[j].Pos).Scaled(x * Stars[i].Mass))
						}
					} else {
						Collides = append(Collides, [2]int{i, j})
					}
				}
				Stars[i].Pos = Stars[i].Pos.Add(Stars[i].Speed.Scaled(dt))
				if Dist(Pos, Stars[i].Pos) > MaxDist {
					Stars[i].Pos = Stars[i].Pos.Sub(Pos).Scaled(0.99).Add(Pos)
					Stars[i].Speed = Stars[i].Speed.Scaled(-0.9)
				}
			}
		}
		for i := 0; i < len(Collides); i++ {
			Collide(Collides[i][0], Collides[i][1])
		}
		Collides = [][2]int{}

		for i := len(DelStars) - 1; i >= 0; i-- {
			DelStar(DelStars[i])
		}
		DelStars = []int{}

		for i := StarsLen; i < MinStars; i++ {
			GenerateNewMeteor(MinR, MaxR, randInterval(MinMass, MaxMass))
		}

		//////////////////////////////////////////////////////////////// drawing

		if Focused >= 0 {
			Pos = Stars[Focused].Pos.Sub(WinSize.Scaled(win.Zoom / 2))
			win.ZoomPos.X = Pos.X
			win.ZoomPos.Y = Pos.Y
		}

		imd.Clear()

		if Tracers {
			imd.Color = pixel.RGBA{A: dt}
			imd.Push(pixel.V(0, 0), pixel.V(WinSizeW, WinSizeH))
			imd.Rectangle(0)
		} else {
			win.Clear(pixel.RGB(0, 0, 0))
		}

		for _, S := range Stars {
			Pos = win.GetPos(S.Pos)
			R = Max(S.Radius/win.Zoom, 0.7)
			if Pos.X > -R && Pos.X < WinSizeW+R && Pos.Y > -R && Pos.Y < WinSizeH+R {
				imd.Color = S.Color
				imd.Push(win.GetPos(S.Pos))
				imd.Circle(R, 0)
			}
			if Tracers {
				if S.Radius < ColorOfRadiusKeys[2] {
					R = Max((CloudRA*S.Radius*S.Radius+(CloudRB+1)*S.Radius)/win.Zoom, 2)
					if Pos.X > -R && Pos.X < WinSizeW+R && Pos.Y > -R && Pos.Y < WinSizeH+R {
						Color := S.Color.Scaled(S.Radius / ColorOfRadiusKeys[2])
						Color = Color.Scaled(0.1)
						Color.A *= dt
						imd.Color = Color
						imd.Push(Pos)
						imd.Circle(R, 0)
					}
				} else if S.Radius >= ColorOfRadiusKeys[4] {
					R *= 4 - randInterval2(1, 3)
					if Pos.X > -R && Pos.X < WinSizeW+R && Pos.Y > -R && Pos.Y < WinSizeH+R {
						Color := S.Color.Scaled(S.Radius / ColorOfRadiusKeys[4])
						Color = Color.Scaled(0.01)
						Color.A *= dt
						imd.Color = Color
						imd.Push(Pos)
						imd.Circle(R, 0)
					}
				}
			}

		}

		//CheckCenter()

		imd.Draw(win)
		win.Update()

		select {
		case <-TickTimer:
		}
	}
}
