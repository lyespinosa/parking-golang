package models

import "sync"

type CarHandler struct {
	Cars  []*Car
	Mutex sync.Mutex
}

func NewCarHandler() *CarHandler {
	return &CarHandler{
		Cars: make([]*Car, 0),
	}
}

func (carHandle *CarHandler) Add(car *Car) {
	carHandle.Mutex.Lock()
	defer carHandle.Mutex.Unlock()
	carHandle.Cars = append(carHandle.Cars, car)
}

func (carHandle *CarHandler) Remove(car *Car) {
	carHandle.Mutex.Lock()
	defer carHandle.Mutex.Unlock()
	for i, c := range carHandle.Cars {
		if c == car {
			carHandle.Cars = append(carHandle.Cars[:i], carHandle.Cars[i+1:]...)
			break
		}
	}
}

func (carHandle *CarHandler) GetCars() []*Car {
	carHandle.Mutex.Lock()
	defer carHandle.Mutex.Unlock()
	return carHandle.Cars
}
