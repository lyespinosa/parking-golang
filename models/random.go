package models

import (
	"math/rand"
	"time"
)

func RandomSleep(value int) {
	rand.Seed(time.Now().UnixNano())
	number := rand.Intn(value*10) + value
	time.Sleep(time.Millisecond * time.Duration(number))
}
