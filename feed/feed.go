package feed

import (
	"log"

	"github.com/mmcloughlin/geohash"
	protos "github.com/pogodevorg/POGOProtos-go"
)

type Point struct {
	ID  string
	Lat float64
	Lon float64
}

type Coord struct {
	Lat float64
	Lon float64
}

type Handler struct {
	sp *SpawnPoints
	gp *GymPoints
	pp *PokestopPoints
}

func (h *Handler) Push(entry interface{}) {
	switch e := entry.(type) {
	default:
		// NOOP: Will not report type
	case *protos.GetMapObjectsResponse:
		for _, cell := range e.GetMapCells() {
			for _, point := range cell.GetSpawnPoints() {
				h.sp.Add(point.Latitude, point.Longitude, "")
			}
			for _, pokemon := range cell.GetWildPokemons() {
				h.sp.Add(pokemon.Latitude, pokemon.Longitude, pokemon.SpawnPointId)
				hash := geohash.Encode(pokemon.Latitude, pokemon.Longitude)
				log.Printf(
					"Found spawn at %s / %s / %f,%f: %s (%d lm:). Hides in %dms",
					pokemon.SpawnPointId,
					hash,
					pokemon.Latitude,
					pokemon.Longitude,
					pokemon.PokemonData.PokemonId.String(),
					pokemon.EncounterId,
					pokemon.LastModifiedTimestampMs,
					pokemon.TimeTillHiddenMs,
				)
			}
			for _, fort := range cell.GetForts() {
				switch fort.Type.String() {
				case "CHECKPOINT":
					go h.pp.Add(fort.Latitude, fort.Longitude, fort.Id)
				case "GYM":
					go h.gp.Add(fort.Latitude, fort.Longitude, fort.Id)
				}
			}
		}
	}
}

func New() *Handler {
	return &Handler{
		sp: &SpawnPoints{
			Points: map[string][]Point{},
		},
		pp: &PokestopPoints{
			Points: map[string][]Point{},
		},
		gp: &GymPoints{
			Points: map[string][]Point{},
		},
	}
}
