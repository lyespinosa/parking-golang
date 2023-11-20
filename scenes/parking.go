package scenes

import (
	"parking/models"
	"sync"
	"time"

	"github.com/oakmound/oak/v4"
	"github.com/oakmound/oak/v4/event"
	"github.com/oakmound/oak/v4/scene"
)

type MainScene struct {
}

func NewParkingScene() *MainScene {
	return &MainScene{}
}

func (ps *MainScene) Start() {
	firstTime := true
	manager := models.NewCarHandler()
	mutex := sync.Mutex{}

	oak.AddScene("mainScene", scene.Scene{
		Start: func(ctx *scene.Context) {
			parking := models.NewParking(ctx)

			event.GlobalBind(ctx, event.Enter, func(enterPayload event.EnterPayload) event.Response {
				if !firstTime {
					return 0
				}
				firstTime = false

				for {
					car := models.NewCar(ctx)
					go models.CarBehaviour(car, manager, parking, &mutex)
					time.Sleep(time.Millisecond * time.Duration(models.GetRandomNumber(1000, 2000)))
				}

			})
		},
	})

	oak.Init("mainScene")
}
