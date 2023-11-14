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

func (cm *CarHandler) Add(car *Car) {
	cm.Mutex.Lock()
	defer cm.Mutex.Unlock()
	cm.Cars = append(cm.Cars, car)
}

func (cm *CarHandler) Remove(car *Car) {
	cm.Mutex.Lock()
	defer cm.Mutex.Unlock()
	for i, c := range cm.Cars {
		if c == car {
			cm.Cars = append(cm.Cars[:i], cm.Cars[i+1:]...)
			break
		}
	}
}

func (cm *CarHandler) GetCars() []*Car {
	cm.Mutex.Lock()
	defer cm.Mutex.Unlock()
	return cm.Cars
}
