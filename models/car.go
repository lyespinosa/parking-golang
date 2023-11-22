package models

import (
	"image/color"
	"sync"
	"time"

	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/render/mod"

	"github.com/oakmound/oak/v4/alg/floatgeom"
	"github.com/oakmound/oak/v4/entities"
	"github.com/oakmound/oak/v4/scene"
)

const (
	spawn   = 185.00
	despawn = 145.00
	speed   = 8
)

type Car struct {
	area   floatgeom.Rect2
	entity *entities.Entity
	mu     sync.Mutex
}

func NewCar(ctx *scene.Context) *Car {
	area := floatgeom.NewRect2(530, -40, 560, 0)
	spriteHpath := "assets/red-car.png"
	spriteVpath := "assets/red-car-h.png"
	spriteH, _ := render.LoadSprite(spriteHpath)
	spriteV, _ := render.LoadSprite(spriteVpath)

	carR := render.NewSwitch("Down", map[string]render.Modifiable{
		"Left":  spriteH.Copy().Modify(mod.FlipX),
		"Right": spriteH,
		"Down":  spriteV,
	})
	entity := entities.New(ctx, entities.WithRect(area), entities.WithColor(color.RGBA{255, 0, 0, 255}), entities.WithRenderable(carR), entities.WithDrawLayers([]int{1, 2}))

	return &Car{
		area:   area,
		entity: entity,
	}
}

func (c *Car) RunCar(mutex *sync.Mutex, manager *CarHandler, parking *Parking) {
	manager.Add(c)
	c.Stacked(manager)
	spotAvailable := parking.GetEmptySpot()
	mutex.Lock()
	c.EnterParking(manager)
	mutex.Unlock()
	c.Enter(spotAvailable, manager)
	RandomSleep(10000)
	c.LeaveSpot(manager)
	parking.SpotParking(spotAvailable)
	c.Leave(spotAvailable, manager)
	mutex.Lock()
	c.FinishEnterParking(manager)
	mutex.Unlock()
	c.GoAway(manager)
	c.Remove()
	manager.Remove(c)
}

func wait() {
	time.Sleep(speed * time.Millisecond)
}

func (c *Car) X() float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.entity.X()
}

func (c *Car) Y() float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.entity.Y()
}

func (c *Car) collides(direction string, cars []*Car) bool {
	distance := 25.0
	for _, car := range cars {
		switch direction {
		case "left":
			if c.X() > car.X() && c.X()-car.X() < distance && c.Y() == car.Y() {
				return true
			}
		case "right":
			if c.X() < car.X() && car.X()-c.X() < distance && c.Y() == car.Y() {
				return true
			}
		case "up":
			if c.Y() > car.Y() && c.Y()-car.Y() < distance && c.X() == car.X() {
				return true
			}
		case "down":
			if c.Y() < car.Y() && car.Y()-c.Y() < distance && c.X() == car.X() {
				return true
			}
		}
	}
	return false
}

func (c *Car) MovetoY(dy float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entity.ShiftY(dy)
}

func (c *Car) MovetoX(dx float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entity.ShiftX(dx)
}

func (c *Car) Stacked(manager *CarHandler) {

	for c.Y() < 145 {
		if !c.collides("down", manager.GetCars()) {
			c.MovetoY(2)
			wait()
		}
	}

}

func (c *Car) EnterParking(manager *CarHandler) {
	for c.Y() < spawn {
		if !c.collides("down", manager.GetCars()) {
			c.MovetoY(1)
			wait()
		}
	}
}

func (c *Car) FinishEnterParking(manager *CarHandler) {
	for c.Y() > despawn {
		if !c.collides("up", manager.GetCars()) {
			c.MovetoY(-1)
			wait()
		}
	}
}

func (c *Car) Enter(spot *Spot, manager *CarHandler) {
	for index := 0; index < len(*spot.GetRouteEntering()); index++ {
		routes := *spot.GetRouteEntering()
		switch routes[index].route {
		case "right":
			c.entity.Renderable.(*render.Switch).Set("Right")
			for c.X() < routes[index].spot {
				if !c.collides("right", manager.GetCars()) {
					c.MovetoX(1)
					wait()
				}
			}
		case "left":
			c.entity.Renderable.(*render.Switch).Set("Left")
			for c.X() > routes[index].spot {
				if !c.collides("left", manager.GetCars()) {
					c.MovetoX(-1)
					wait()
				}
			}
		case "down":
			c.entity.Renderable.(*render.Switch).Set("Down")
			for c.Y() < routes[index].spot {
				if !c.collides("down", manager.GetCars()) {
					c.MovetoY(1)
					wait()
				}
			}
		case "up":
			c.entity.Renderable.(*render.Switch).Set("Down")
			for c.Y() > routes[index].spot {
				if !c.collides("up", manager.GetCars()) {
					c.MovetoY(-1)
					wait()
				}
			}
		}
	}
}

func (c *Car) Remove() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entity.Destroy()
}

func (c *Car) Leave(spot *Spot, manager *CarHandler) {
	for index := 0; index < len(*spot.GetRouteLeaving()); index++ {
		routes := *spot.GetRouteLeaving()
		switch routes[index].route {
		case "right":
			c.entity.Renderable.(*render.Switch).Set("Right")
			for c.X() < routes[index].spot {
				if !c.collides("right", manager.GetCars()) {
					c.MovetoX(1)
					wait()
				}
			}
		case "left":
			c.entity.Renderable.(*render.Switch).Set("Left")
			for c.X() > routes[index].spot {
				if !c.collides("left", manager.GetCars()) {
					c.MovetoX(-1)
					wait()
				}
			}
		case "up":
			c.entity.Renderable.(*render.Switch).Set("Down")
			for c.Y() > routes[index].spot {
				if !c.collides("up", manager.GetCars()) {
					c.MovetoY(-1)
					wait()
				}
			}
		case "down":
			c.entity.Renderable.(*render.Switch).Set("Down")
			for c.Y() < routes[index].spot {
				if !c.collides("down", manager.GetCars()) {
					c.MovetoY(1)
					wait()
				}
			}
		}
	}
}

func (c *Car) LeaveSpot(manager *CarHandler) {
	spotX := c.X()
	for c.X() > spotX-50 {
		if !c.collides("left", manager.GetCars()) {
			c.MovetoX(-1)
			time.Sleep(speed * time.Millisecond)
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func (c *Car) GoAway(manager *CarHandler) {
	for c.Y() > -20 {
		if !c.collides("up", manager.GetCars()) {
			c.MovetoY(-1)
			time.Sleep(speed * time.Millisecond)
		}
	}
}
