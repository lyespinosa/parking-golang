package models

import (
	"github.com/oakmound/oak/v4/alg/floatgeom"
)

type Spot struct {
	area          *floatgeom.Rect2
	routeEntering *[]Route
	routeLeaving  *[]Route
	number        int
	isAvailable   bool
}

func NewSpot(x, y, x2, y2 float64, column, number int) *Spot {
	routeEntering := getRouteEntering(x, y, column)
	routeLeaving := getRouteLeaving()
	area := floatgeom.NewRect2(x, y, x2, y2)

	return &Spot{
		area:          &area,
		routeEntering: routeEntering,
		routeLeaving:  routeLeaving,
		number:        number,
		isAvailable:   true,
	}
}

func getRouteEntering(x, y float64, column int) *[]Route {
	var directions []Route

	if column == 1 {
		directions = append(directions, *newRoute("left", 445))
	} else if column == 2 {
		directions = append(directions, *newRoute("left", 355))
	} else if column == 3 {
		directions = append(directions, *newRoute("left", 265))
	} else if column == 4 {
		directions = append(directions, *newRoute("left", 175))
	} else if column == 5 {
		directions = append(directions, *newRoute("left", 85))
	}

	directions = append(directions, *newRoute("down", y+5))
	directions = append(directions, *newRoute("left", x+5))

	return &directions
}

func getRouteLeaving() *[]Route {
	var directions []Route

	directions = append(directions, *newRoute("down", 380))
	directions = append(directions, *newRoute("right", 475))
	directions = append(directions, *newRoute("up", 185))

	return &directions
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

func (p *Spot) GetIsAvailable() bool {

	return p.isAvailable
}

func (p *Spot) SetIsAvailable(isAvailable bool) {
	p.isAvailable = isAvailable
}
