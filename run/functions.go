package run

import (
	"github.com/faiface/pixel"
	"math"
	"math/rand"
)

func randInterval(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}

func randInterval2(min, max float64) float64 {
	R := rand.Float64()
	if R < 0.5 && rand.Float64() < 1-R*2 {
		R += 0.5
	}
	return R*(max-min) + min
}

//func randPolar(MinR, MaxR float64) pixel.Vec {
//	angle := rand.Float64() * 2 * math.Pi
//	return pixel.V(math.Cos(angle), math.Sin(angle)).Scaled(randInterval(MinR, MaxR))
//}

func randPolarAngle(MinR, MaxR, MinAngle, MaxAngle float64) pixel.Vec {
	angle := randInterval(MinAngle, MaxAngle) * 2 * math.Pi
	return pixel.V(math.Cos(angle), math.Sin(angle)).Scaled(randInterval(MinR, MaxR))
}

//func randPolar2(MinR, MaxR float64) pixel.Vec {
//	angle := rand.Float64() * 2 * math.Pi
//	R := rand.Float64()
//	if R < 0.5 && rand.Float64() < 1-R*2 {
//		R += 0.5
//	}
//	return pixel.V(math.Cos(angle), math.Sin(angle)).Scaled(R*(MaxR-MinR) + MinR)
//}
//
//func randPolar3(MinR, MaxR float64) pixel.Vec {
//	angle := rand.Float64() * 2 * math.Pi
//	R := rand.Float64()
//	if R < 0.5 && rand.Float64() < 1-R*2 {
//		R += 0.5
//	}
//	R = R*R*(3-2*R)*(MaxR-MinR) + MinR
//	return pixel.V(math.Cos(angle), math.Sin(angle)).Scaled(R)
//}

func Dist(u, v pixel.Vec) float64 {
	dx, dy := u.X-v.X, u.Y-v.Y
	return math.Sqrt(dx*dx + dy*dy)
}

func Max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func Min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
