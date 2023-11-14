package models

import (
	"math/rand"
	"time"
)

func GetRandomNumber(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	number := rand.Intn(max) + min
	return number
}
