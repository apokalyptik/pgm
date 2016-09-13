package feed

import (
	"log"
	"sync"
	"time"

	"github.com/mmcloughlin/geohash"
)

type GymPoints struct {
	sync.RWMutex
	Points     map[string][]Point
	LastUpdate time.Time
}

func (s *GymPoints) Add(lat, lon float64, id string) {
	hash := geohash.Encode(lat, lon)
	s.RLock()
	if _, ok := s.Points[hash]; !ok {
		s.RUnlock()
		// Check again after upgrading because of race
		s.Lock()
		if _, ok := s.Points[hash]; !ok {
			s.Points[hash] = []Point{Point{Lat: lat, Lon: lon, ID: id}}
			s.LastUpdate = time.Now()
			s.Unlock()
			log.Printf("Added new gym and geohash: %s: %s: %f,%f", hash, id, lat, lon)
			return
		}
		s.Unlock()
		s.RLock()
	}
	// If we made it here then s.Points[hash] MUST exist
	lastIndex := 0
	for i, c := range s.Points[hash] {
		lastIndex = i
		if c.Lat == lat && c.Lon == lon {
			s.RUnlock()
			// no-op
			return
		}
	}
	s.RUnlock()
	s.Lock()
	// check again because of race, however, data is only appended
	// so we only have to check anything after lastIndex
	if len(s.Points[hash]) >= lastIndex {
		for i := lastIndex + 1; i < len(s.Points[hash]); i++ {
			if s.Points[hash][i].Lat == lat && s.Points[hash][i].Lon == lon {
				s.Unlock()
				// no-op
				return
			}
		}
	}
	s.Points[hash] = append(s.Points[hash], Point{Lat: lat, Lon: lon, ID: id})
	s.LastUpdate = time.Now()
	s.Unlock()
	log.Printf("Added new gym point to geohash: %s: %s: %f,%f", hash, id, lat, lon)
}
