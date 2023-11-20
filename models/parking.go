package models

import (
	"fmt"
	"image/color"
	"math/rand"
	"sync"
	"time"

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

func SpotNumber() []int {
	rand.Seed(time.Now().UnixNano())

	numbers := make([]int, 20)
	for i := 0; i < 20; i++ {
		numbers[i] = i + 1
	}

	rand.Shuffle(len(numbers), func(i, j int) {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	})

	return numbers
}

var (
	spotNumber = SpotNumber()

	spotList = []*Spot{
		NewSpot(450, 390, 500, 420, 1, spotNumber[0]),
		NewSpot(150, 270, 200, 300, 4, spotNumber[1]),
		NewSpot(150, 390, 200, 420, 4, spotNumber[2]),
		NewSpot(450, 210, 500, 240, 1, spotNumber[3]),
		NewSpot(50, 270, 100, 300, 5, spotNumber[4]),
		NewSpot(350, 270, 400, 300, 2, spotNumber[5]),
		NewSpot(350, 330, 400, 360, 2, spotNumber[6]),
		NewSpot(350, 390, 400, 420, 2, spotNumber[7]),
		NewSpot(250, 210, 300, 240, 3, spotNumber[8]),
		NewSpot(250, 270, 300, 300, 3, spotNumber[9]),
		NewSpot(50, 390, 100, 420, 5, spotNumber[10]),
		NewSpot(250, 390, 300, 420, 3, spotNumber[11]),
		NewSpot(150, 210, 200, 240, 4, spotNumber[12]),
		NewSpot(450, 270, 500, 300, 1, spotNumber[13]),
		NewSpot(150, 330, 200, 360, 4, spotNumber[14]),
		NewSpot(450, 330, 500, 360, 1, spotNumber[15]),
		NewSpot(350, 210, 400, 240, 2, spotNumber[16]),
		NewSpot(50, 210, 100, 240, 5, spotNumber[17]),
		NewSpot(50, 330, 100, 360, 5, spotNumber[18]),
		NewSpot(250, 330, 300, 360, 3, spotNumber[19]),
	}
)

func drawGround(ctx *scene.Context, spots []*Spot) {

	fmt.Print(spotNumber)

	ground := floatgeom.NewRect2(0, 0, 640, 480)
	entities.New(ctx, entities.WithRect(ground), entities.WithColor(color.RGBA{85, 85, 85, 255}))

	grass := floatgeom.NewRect2(0, 0, 500, 170)
	grassMin := floatgeom.NewRect2(620, 0, 640, 170)
	entities.New(ctx, entities.WithRect(grass), entities.WithColor(color.RGBA{84, 137, 51, 255}))
	entities.New(ctx, entities.WithRect(grassMin), entities.WithColor(color.RGBA{84, 137, 51, 255}))

	parkingEnter := floatgeom.NewRect2(510, 40, 520, 60)
	entities.New(ctx, entities.WithRect(parkingEnter), entities.WithColor(color.RGBA{250, 250, 250, 1}))

	parkingLeave := floatgeom.NewRect2(600, 40, 610, 60)
	entities.New(ctx, entities.WithRect(parkingLeave), entities.WithColor(color.RGBA{250, 250, 250, 1}))

	for _, spot := range spots {
		entities.New(ctx, entities.WithRect(*spot.GetArea()), entities.WithColor(color.RGBA{255, 255, 255, 255}))
	}
}

func (p *Parking) ManageParkingSpot(spot *Spot) {
	p.mu.Lock()
	defer p.mu.Unlock()

	spot.SetIsEmpty(true)
	p.availableCond.Signal()
}

func NewParking(ctx *scene.Context) *Parking {

	spots := spotList

	drawGround(ctx, spots)
	queue := NewCarWaiting()
	p := &Parking{
		spots:     spots,
		queueCars: queue,
	}
	p.availableCond = sync.NewCond(&p.mu)
	return p
}

func (p *Parking) GetSpotAvailable() *Spot {
	p.mu.Lock()
	defer p.mu.Unlock()

	for {
		for _, spot := range p.spots {
			if spot.GetIsEmpty() {
				spot.SetIsEmpty(false)
				return spot
			}
		}
		p.availableCond.Wait()
	}
}

func NewCarWaiting() *CarQueue {
	return &CarQueue{
		cars: make([]Car, 0),
	}
}
