package beehive

import (
	"context"
	"log"
	"time"

	"github.com/apokalyptik/pgm/encryption"
	"github.com/pogodevorg/pgoapi-go/api"
	"github.com/pogodevorg/pgoapi-go/auth"
)

var Feed api.Feed

func mindWorker(id int, kind, username, password string) {
	var session *api.Session

	provider, err := auth.NewProvider(kind, username, password)
	if err != nil {
		log.Printf("Worker #%d provider error: %s", id, err.Error())
		return
	}

	tick := time.Tick(10 * time.Second)
	for {
		select {
		case next := <-stepQueue:
			if session == nil {
				session = api.NewSession(
					provider,
					&api.Location{
						Lat:      next.Lat,
						Lon:      next.Lon,
						Alt:      next.Alt,
						Accuracy: 3.0,
					},
					Feed,
					&encryption.FemotCrypto{},
					false,
				)

				if err := session.Init(context.Background()); err != nil {
					log.Printf("Worker #%d initialization error: %s", id, err.Error())
					return
				}
				log.Printf("Worker #%d Started at %f,%f,%f", id, next.Lat, next.Lon, next.Alt)
			} else {
				session.MoveTo(next)
				log.Printf("Worker #%d moved to %f,%f,%f", id, next.Lat, next.Lon, next.Alt)
			}
			time.Sleep(time.Second)
			session.GetPlayerMap(context.Background())
			<-tick
		}
	}
}
