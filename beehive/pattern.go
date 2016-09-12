package beehive

import (
	"fmt"
	"log"
	"strings"

	"github.com/apokalyptik/pgm/ll"
)

func doSimpleBeehiveFrom(start *ll.Coord, steps int, accounts [][]string) {
	grid := start.Grid(steps)
	log.Println("Built un-elevated", steps, "step,", len(grid), "coordinate, grid from starting position of", start.Lat, start.Lng)

	// Step 3: gather elevation data for all the cells
	if err := grid.Elevate(); err != nil {
		log.Fatalf("Unable to get elevation data for the grid: %s", err.Error())
	}
	log.Println("Successfully elevated the grid")

	var gridOutParts = []string{}
	for _, coord := range grid {
		gridOutParts = append(gridOutParts, fmt.Sprintf("(%f,%f,%f)", coord.Lat, coord.Lng, coord.Elevation))
	}
	log.Println("Complete (lat,lon,elevation) grid:", strings.Join(gridOutParts, ","))

	go mindGridQueue(grid)
	for id, account := range accounts {
		go mindWorker(id, account[0], account[1], account[2])
	}

	<-make(chan struct{})
}
