package debug

import (
	"fmt"

	protos "github.com/pogodevorg/POGOProtos-go"
)

type Feed struct{}

func (f *Feed) Push(entry interface{}) {
	switch e := entry.(type) {
	default:
		// NOOP: Will not report type
	case *protos.GetMapObjectsResponse:
		cells := e.GetMapCells()
		for _, cell := range cells {
			pokemons := cell.GetWildPokemons()
			if len(pokemons) > 0 {
				fmt.Println(pokemons)
			}
			forts := cell.GetForts()
			if len(forts) > 0 {
				fmt.Println(forts)
			}
		}
	}
}
