package models

import (
	"image/color"
	"sync"

	"github.com/oakmound/oak/v4/alg/floatgeom"
	"github.com/oakmound/oak/v4/entities"
	"github.com/oakmound/oak/v4/scene"
)

type CarQueue struct {
	cars []Car
}

type Parking struct {
	spots         []*Spot
	queueCars     *CarQueue
	mu            sync.Mutex
	availableCond *sync.Cond
}

func NewParking(ctx *scene.Context) *Parking {

	spots := []*Spot{

		NewSpot(450, 210, 500, 240, 1, 1),
		NewSpot(450, 270, 500, 300, 1, 2),
		NewSpot(450, 330, 500, 360, 1, 3),
		NewSpot(450, 390, 500, 420, 1, 4),
		NewSpot(350, 210, 400, 240, 2, 5),
		NewSpot(350, 270, 400, 300, 2, 6),
		NewSpot(350, 330, 400, 360, 2, 7),
		NewSpot(350, 390, 400, 420, 2, 8),
		NewSpot(250, 210, 300, 240, 3, 9),
		NewSpot(250, 270, 300, 300, 3, 10),
		NewSpot(250, 330, 300, 360, 3, 11),
		NewSpot(250, 390, 300, 420, 3, 12),
		NewSpot(150, 210, 200, 240, 4, 13),
		NewSpot(150, 270, 200, 300, 4, 14),
		NewSpot(150, 330, 200, 360, 4, 15),
		NewSpot(150, 390, 200, 420, 4, 16),
		NewSpot(50, 210, 100, 240, 5, 17),
		NewSpot(50, 270, 100, 300, 5, 18),
		NewSpot(50, 330, 100, 360, 5, 19),
		NewSpot(50, 390, 100, 420, 5, 20),
	}

	setUpScene(ctx, spots)
	queue := NewCarQueue()
	p := &Parking{
		spots:     spots,
		queueCars: queue,
	}
	p.availableCond = sync.NewCond(&p.mu)
	return p
}

func (p *Parking) GetSpots() []*Spot {
	return p.spots
}

func (p *Parking) GetSpotAvailable() *Spot {
	p.mu.Lock()
	defer p.mu.Unlock()

	for {
		for _, spot := range p.spots {
			if spot.GetIsAvailable() {
				spot.SetIsAvailable(false)
				return spot
			}
		}
		p.availableCond.Wait()
	}
}

func (p *Parking) ReleaseParkingSpot(spot *Spot) {
	p.mu.Lock()
	defer p.mu.Unlock()

	spot.SetIsAvailable(true)
	p.availableCond.Signal()
}

func (p *Parking) GetQueueCars() *CarQueue {
	return p.queueCars
}

func NewCarQueue() *CarQueue {
	return &CarQueue{
		cars: make([]Car, 0),
	}
}

func setUpScene(ctx *scene.Context, spots []*Spot) {

	parkingArea := floatgeom.NewRect2(10, 170, 610, 470)
	entities.New(ctx, entities.WithRect(parkingArea), entities.WithColor(color.RGBA{100, 100, 100, 1}))

	parkingDoor := floatgeom.NewRect2(530, 140, 600, 160)
	entities.New(ctx, entities.WithRect(parkingDoor), entities.WithColor(color.RGBA{250, 250, 250, 1}))

	for _, spot := range spots {
		entities.New(ctx, entities.WithRect(*spot.GetArea()), entities.WithColor(color.RGBA{255, 255, 255, 255}))
	}
}
