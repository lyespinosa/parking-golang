package main

import (
	"parking/scenes"

	"github.com/oakmound/oak/v4"
)

func main() {
	parkingScene := scenes.NewParkingScene()
	parkingScene.Generate()
	oak.Init("mainScene")
}
