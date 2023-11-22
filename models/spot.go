package models

import (
	"fmt"

	"github.com/oakmound/oak/v4/alg/floatgeom"
)

var (
	column    = 0
	iteration = 0
)

type Spot struct {
	isEmpty       bool
	number        int
	area          *floatgeom.Rect2
	routeEntering *[]Route
	routeLeaving  *[]Route
}

func NewSpot(x, y, x2, y2 float64) *Spot {

	column = (iteration % 5) + 1
	iteration++
	fmt.Println(column)
	routeEntering := getRouteEntering(x, y, column)
	routeLeaving := getRouteLeaving()
	area := floatgeom.NewRect2(x, y, x2, y2)

	return &Spot{
		area:          &area,
		routeEntering: routeEntering,
		routeLeaving:  routeLeaving,
		number:        iteration,
		isEmpty:       true,
	}
}

func (p *Spot) GetArea() *floatgeom.Rect2 {
	return p.area
}

func (p *Spot) GetNumber() int {
	return p.number
}

func (p *Spot) GetRouteEntering() *[]Route {
	return p.routeEntering
}

func (p *Spot) GetRouteLeaving() *[]Route {
	return p.routeLeaving
}

func (p *Spot) GetIsEmpty() bool {

	return p.isEmpty
}

func (p *Spot) SetIsEmpty(isEmpty bool) {
	p.isEmpty = isEmpty
}

func getRouteEntering(x, y float64, column int) *[]Route {
	var directions []Route

	if column == 1 {
		directions = append(directions, *newRoute("left", 550))
	} else if column == 2 {
		directions = append(directions, *newRoute("left", 400))
	} else if column == 3 {
		directions = append(directions, *newRoute("left", 300))
	} else if column == 4 {
		directions = append(directions, *newRoute("left", 200))
	} else if column == 5 {
		directions = append(directions, *newRoute("left", 100))
	}
	directions = append(directions, *newRoute("down", y+3))
	directions = append(directions, *newRoute("left", x+3))

	return &directions
}

func getRouteLeaving() *[]Route {
	var directions []Route

	directions = append(directions, *newRoute("down", 440))
	directions = append(directions, *newRoute("right", 570))
	directions = append(directions, *newRoute("up", 185))

	return &directions
}
