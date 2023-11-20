package main

import (
	"parking/scenes"

	"github.com/oakmound/oak/v4"
)

func main() {
	scenes.NewParkingScene().Start()
	oak.Init("mainScene")
}
