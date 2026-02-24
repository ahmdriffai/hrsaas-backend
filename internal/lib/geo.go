package lib

import "math"

const earthRadiusMeter = 6371000 // radius bumi dalam meter

func DistanceMeter(lat1, lng1, lat2, lng2 float64) float64 {
	// konversi derajat ke radian
	toRad := func(deg float64) float64 {
		return deg * math.Pi / 180
	}

	lat1Rad := toRad(lat1)
	lng1Rad := toRad(lng1)
	lat2Rad := toRad(lat2)
	lng2Rad := toRad(lng2)

	dLat := lat2Rad - lat1Rad
	dLng := lng2Rad - lng1Rad

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dLng/2)*math.Sin(dLng/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusMeter * c
}
