package beehive

import (
	"log"

	"github.com/apokalyptik/pgm/ll"
	"github.com/pogodevorg/pgoapi-go/api"
)

var stepQueue = make(chan *api.Location)

func mindGridQueue(g ll.Grid) {
	for {
		for id, loc := range g {
			stepQueue <- loc.Location()
			log.Printf("queued location scan for #%d (%f,%f,%f)", id, loc.Lat, loc.Lng, loc.Elevation)
		}
	}
}
