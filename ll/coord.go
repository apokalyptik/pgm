package ll

import (
	"context"
	"fmt"
	"math"
	"math/rand"

	"github.com/pogodevorg/pgoapi-go/api"

	"googlemaps.github.io/maps"
)

var GeoCodeFailure = fmt.Errorf("Unable to geocode requested location into a lat/lon pair")

type Coord struct {
	maps.LatLng
	Elevation float64
}

func (c *Coord) jitter(byUpTo int) *Coord {
	var jitterMeters = float64(rand.Intn(byUpTo)) / 1000
	var bearing = float64(rand.Intn(359))
	var jitterToEarthsRadius = jitterMeters / earthsRadius

	lat1 := radians(c.Lat)
	lng1 := radians(c.Lng)
	lat2 := math.Asin(
		(math.Sin(lat1) * math.Cos(jitterToEarthsRadius)) +
			(math.Cos(lat1) * math.Sin(jitterToEarthsRadius) * math.Cos(bearing)),
	)
	lng2 := lng1 + math.Atan2(
		(math.Sin(bearing)*math.Sin(jitterToEarthsRadius)*math.Cos(lat1)),
		(math.Cos(jitterToEarthsRadius)-(math.Sin(lat1)*math.Sin(lat2))),
	)
	return &Coord{LatLng: maps.LatLng{Lat: degrees(lat2), Lng: degrees(lng2)}}
}

func (c *Coord) Location() *api.Location {
	return &api.Location{
		Lat:      c.Lat,
		Lon:      c.Lng,
		Alt:      c.Elevation,
		Accuracy: 3.0,
	}
}

func (c *Coord) Next(d float64) *Coord {
	bearing := radians(d)
	lat1 := radians(c.Lat)
	lng1 := radians(c.Lng)
	lat2 := math.Asin(
		(math.Sin(lat1) * math.Cos(dr)) +
			(math.Cos(lat1) * math.Sin(dr) * math.Cos(bearing)),
	)
	lng2 := lng1 + math.Atan2(
		(math.Sin(bearing)*math.Sin(dr)*math.Cos(lat1)),
		(math.Cos(dr)-(math.Sin(lat1)*math.Sin(lat2))),
	)
	return &Coord{LatLng: maps.LatLng{Lat: degrees(lat2), Lng: degrees(lng2)}}
}

func (c *Coord) Grid(steps int) Grid {
	var rval = Grid{}
	rval = append(rval, c)
	for step := 1; step <= steps; step++ {
		for side := 0; side < 6; side++ {
			for brick := 0; brick < (step - 1); brick++ {
				from := rval[len(rval)-1]
				if side == 0 && brick == 0 {
					from = from.Next(240)
				}
				rval = append(rval, from.Next((float64(side) * 60)))
			}
		}
	}
	return rval
}

func NewCoord(location string, jitter int) (*Coord, error) {
	geos, err := Geo.Geocode(context.Background(), &maps.GeocodingRequest{
		Address: location,
	})
	if err != nil {
		return nil, err
	}
	if len(geos) < 1 {
		return nil, GeoCodeFailure
	}
	start := &Coord{LatLng: geos[0].Geometry.Location}
	if jitter > 0 {
		start = start.jitter(jitter)
	}
	return start, nil
}
