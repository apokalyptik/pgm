package feed

import (
	"log"
	"sync"
	"time"

	"github.com/mmcloughlin/geohash"
)

type SpawnPoints struct {
	sync.RWMutex
	Points     map[string][]Point
	LastUpdate time.Time
}

func (s *SpawnPoints) Exists(hash string, lat, lon float64) bool {
	s.RLock()
	defer s.RUnlock()
	if _, ok := s.Points[hash]; !ok {
		return false
	}
	for _, sp := range s.Points[hash] {
		if sp.Lat == lat && sp.Lon == lon {
			return true
		}
	}
	return false
}

func (s *SpawnPoints) Add(lat, lon float64, id string) {
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
			log.Printf("Added new spawn point and geohash: %s: %f,%f", hash, lat, lon)
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
			switch c.ID {
			case id:
				break
			case "":
				s.Lock()
				if c.ID == "" {
					c.ID = id
					s.Unlock()
					log.Printf("Updated spawn point %s/%f,%f with ID: %s", hash, lat, lon, id)
				} else {
					if c.ID != id {
						s.Unlock()
						log.Printf("Received conflicting ID for existing spawn %s/%s/%f,%f got: %s", hash, c.ID, lat, lon, id)
					}
				}
			default:
				log.Printf("Received conflicting ID for existing spawn %s/%s/%f,%f got: %s", hash, c.ID, lat, lon, id)
			}
			return
		}
	}
	s.RUnlock()
	s.Lock()
	// check again because of race, however, data is only appended
	// so we only have to check anything after lastIndex
	if len(s.Points[hash]) >= lastIndex {
		for i := lastIndex + 1; i < len(s.Points[hash]); i++ {
			c := s.Points[hash][i]
			if c.Lat == lat && c.Lon == lon {
				switch c.ID {
				case id:
					break
				case "":
					c.ID = id
					log.Printf("Updated spawn point %s/%f,%f with ID: %s", hash, lat, lon, id)
				default:
					log.Printf("Received conflicting ID for existing spawn %s/%s/%f,%f got: %s", hash, c.ID, lat, lon, id)
				}
				s.Unlock()
				// no-op
				return
			}
		}
	}
	s.Points[hash] = append(s.Points[hash], Point{Lat: lat, Lon: lon, ID: id})
	s.LastUpdate = time.Now()
	s.Unlock()
	log.Printf("Added new spawn point to geohash: %s: %f,%f", hash, lat, lon)
}
