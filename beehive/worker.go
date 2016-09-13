package beehive

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/apokalyptik/pgm/encryption"
	"github.com/pogodevorg/pgoapi-go/api"
	"github.com/pogodevorg/pgoapi-go/auth"
)

var Feed api.Feed

func doWork(id int, kind, username, password string) {
	var session *api.Session

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in doWork for worker", id, r)
		}
	}()

	provider, err := auth.NewProvider(kind, username, password)
	if err != nil {
		log.Printf("Worker #%d provider error: %s", id, err.Error())
		return
	}

	tick := time.Tick(10 * time.Second)
	for {
		time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
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

func mindWorker(id int, kind, username, password string) {
	for {
		for i := 0; i < 5; i++ {
			before := time.Now()
			doWork(id, kind, username, password)
			if time.Now().Sub(before) > (5 * time.Minute) {
				break
			}
			if i == 4 {
				log.Println("Worker", id, "seems to be having a hard time, or the service is down. Abandoning.")
				return
			}
			time.Sleep(time.Second * time.Duration((30*i)+30))
		}
	}
	log.Println("Giving up on worker", id, kind, username, password)
}
