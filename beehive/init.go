package beehive

import "github.com/apokalyptik/pgm/ll"

func Start(coord *ll.Coord, steps int, accounts [][]string) {
	doSimpleBeehiveFrom(coord, steps, accounts)
}
