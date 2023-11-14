package models

type Route struct {
	route string
	spot  float64
}

func newRoute(route string, spot float64) *Route {
	return &Route{
		route: route,
		spot:  spot,
	}
}
