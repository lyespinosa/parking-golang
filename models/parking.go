package models

import (
	"image/color"
	"math/rand"
	"sync"
	"time"

	"github.com/oakmound/oak/v4/alg/floatgeom"
	"github.com/oakmound/oak/v4/entities"
	"github.com/oakmound/oak/v4/scene"
)

type Stack struct {
	cars []Car
}

type Parking struct {
	condition *sync.Cond
	spots     []*Spot
	stack     *Stack
	mu        sync.Mutex
}

var (
	spotList = []*Spot{
		NewSpot(450, 330, 500, 360),
		NewSpot(350, 330, 400, 360),
		NewSpot(250, 210, 300, 240),
		NewSpot(150, 330, 200, 360),
		NewSpot(50, 330, 100, 360),
		NewSpot(450, 270, 500, 300),
		NewSpot(350, 390, 400, 420),
		NewSpot(250, 270, 300, 300),
		NewSpot(150, 390, 200, 420),
		NewSpot(50, 390, 100, 420),
		NewSpot(450, 210, 500, 240),
		NewSpot(350, 270, 400, 300),
		NewSpot(250, 390, 300, 420),
		NewSpot(150, 270, 200, 300),
		NewSpot(50, 270, 100, 300),
		NewSpot(450, 390, 500, 420),
		NewSpot(350, 210, 400, 240),
		NewSpot(250, 330, 300, 360),
		NewSpot(150, 210, 200, 240),
		NewSpot(50, 210, 100, 240),
	}
)

func NewParking(ctx *scene.Context) *Parking {
	spots := setSpotsEntities()
	drawGround(ctx, spots)
	stack := NewCarStacked()
	p := &Parking{
		spots: spots,
		stack: stack,
	}
	p.condition = sync.NewCond(&p.mu)
	return p
}

func setSpotsEntities() []*Spot {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(spotList), func(i, j int) {
		spotList[i], spotList[j] = spotList[j], spotList[i]
	})

	return spotList
}

func drawGround(ctx *scene.Context, spots []*Spot) {

	ground := floatgeom.NewRect2(0, 0, 640, 480)
	grass := floatgeom.NewRect2(0, 0, 500, 170)
	grassMin := floatgeom.NewRect2(620, 0, 640, 170)
	line1 := floatgeom.NewRect2(555, 30, 565, 50)
	line2 := floatgeom.NewRect2(555, 90, 565, 110)

	entities.New(ctx, entities.WithRect(ground), entities.WithColor(color.RGBA{85, 85, 85, 255}))
	entities.New(ctx, entities.WithRect(grass), entities.WithColor(color.RGBA{84, 137, 51, 255}))
	entities.New(ctx, entities.WithRect(grassMin), entities.WithColor(color.RGBA{84, 137, 51, 255}))
	entities.New(ctx, entities.WithRect(line1), entities.WithColor(color.RGBA{255, 255, 255, 255}))
	entities.New(ctx, entities.WithRect(line2), entities.WithColor(color.RGBA{255, 255, 255, 255}))

	for _, spot := range spots {
		entities.New(ctx, entities.WithRect(*spot.GetArea()), entities.WithColor(color.RGBA{255, 255, 255, 255}))
	}
}

func (p *Parking) SpotParking(spot *Spot) {
	p.mu.Lock()
	defer p.mu.Unlock()
	spot.SetIsEmpty(true)
	p.condition.Signal()
}

func NewCarStacked() *Stack {
	return &Stack{
		cars: make([]Car, 0),
	}
}

func (p *Parking) GetEmptySpot() *Spot {
	p.mu.Lock()
	defer p.mu.Unlock()

	for {
		spot := p.findEmptySpot()
		if spot != nil {
			spot.SetIsEmpty(false)
			return spot
		}
		p.condition.Wait()
	}
}

func (p *Parking) findEmptySpot() *Spot {
	for _, spot := range p.spots {
		if spot.GetIsEmpty() {
			return spot
		}
	}
	return nil
}
