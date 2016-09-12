package ll

import "math"

const (
	distanceToNextPoint = 0.113364 // km
	earthsRadius        = 6378.14  // km
	dr                  = distanceToNextPoint / earthsRadius
)

func radians(d float64) float64 {
	return (d * (math.Pi / 180))
}

func degrees(r float64) float64 {
	return (r * (180 / math.Pi))
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
