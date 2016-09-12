package ll

import (
	"context"
	"fmt"

	"googlemaps.github.io/maps"
)

var Geo *maps.Client

var ErrElevationRequestFailed = fmt.Errorf("Failed to get elevation data")

type Grid []*Coord

func (g Grid) Steps() int {
	l := len(g) - 1
	if l == 0 {
		return 1
	}

	s := 2
	n := 6
	for {
		if n == l {
			return s
		}
		s++
		n = n + (6 * (s - 1))
	}
}

func (g Grid) LatLng() []maps.LatLng {
	var rval []maps.LatLng
	for _, coord := range g {
		rval = append(rval, coord.LatLng)
	}
	return rval
}

func (g Grid) Elevate() error {
	req := new(maps.ElevationRequest)
	req.Locations = g.LatLng()
	eles, err := Geo.Elevation(context.Background(), req)
	if err != nil {
		return err
	}
	if len(eles) < 1 {
		return ErrElevationRequestFailed
	}
	for _, rsp := range eles {
		for _, coord := range g {
			if coord.Elevation != 0 {
				continue
			}
			if toFixed(coord.Lng, 4) != toFixed(rsp.Location.Lng, 4) {
				continue
			}
			if toFixed(coord.Lat, 4) != toFixed(rsp.Location.Lat, 4) {
				continue
			}
			coord.Elevation = rsp.Elevation
		}
	}
	return nil
}
