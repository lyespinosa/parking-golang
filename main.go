package main

import (
	"fmt"
	"image/color"
	"sync"
	"time"

	"github.com/oakmound/oak/v4"
	"github.com/oakmound/oak/v4/alg/floatgeom"
	"github.com/oakmound/oak/v4/entities"
	"github.com/oakmound/oak/v4/scene"
	"gonum.org/v1/gonum/stat/distuv"
)

const parkingSpaces = 20

func main() {
	oak.AddScene("parkingLot", scene.Scene{Start: func(ctx *scene.Context) {

		for i := 0; i < parkingSpaces/2; i++ {
			entities.New(ctx,
				entities.WithRect(floatgeom.NewRect2WH(250+float64(i*40), 150, 25, 45)),
				entities.WithColor(color.RGBA{128, 128, 128, 255}),
			)
			entities.New(ctx,
				entities.WithRect(floatgeom.NewRect2WH(250+float64(i*40), 250, 25, 45)),
				entities.WithColor(color.RGBA{128, 128, 128, 255}),
			)
		}

		entities.New(ctx,
			entities.WithRect(floatgeom.NewRect2WH(180, 180, 25, 80)),
			entities.WithColor(color.RGBA{255, 0, 0, 255}),
		)

	}})
	oak.Init("parkingLot")
}

const lambda = 2
const numCars = 20

func entrada() {
	var wg sync.WaitGroup
	mutex := make(chan struct{}, 1)
	poisson := distuv.Poisson{Lambda: lambda, Src: nil}
	for i := 1; i <= numCars; i++ {
		wg.Add(1)
		waitTime := time.Duration(poisson.Rand()) * time.Second
		time.Sleep(waitTime)
		go func(carID int) {
			defer wg.Done()

			fmt.Printf("Car %d spanwed.\n", carID)
			mutex <- struct{}{} // Bloquea la entrada

			fmt.Printf("Car %d entering to enter parking.\n", carID)
			time.Sleep(4 * time.Second)
			fmt.Printf("Car %d leaving to enter parking.\n", carID)

			<-mutex
		}(i)
	}

	wg.Wait()
	fmt.Println("All cars have entered and exited the parking lot.")
}
