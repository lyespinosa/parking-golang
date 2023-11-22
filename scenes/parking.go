package scenes

import (
	"parking/models"
	"sync"

	"github.com/oakmound/oak/v4"
	"github.com/oakmound/oak/v4/event"
	"github.com/oakmound/oak/v4/scene"
)

var (
	mutex sync.Mutex
)

type MainScene struct {
}

func NewParkingScene() *MainScene {
	return &MainScene{}
}

func runParking(ctx *scene.Context, handler *models.CarHandler, parking *models.Parking) {
	event.GlobalBind(ctx, event.Enter, func(enterPayload event.EnterPayload) event.Response {
		for {
			car := models.NewCar(ctx)
			go car.RunCar(&mutex, handler, parking)
			models.RandomSleep(300)
		}
	})
}

func (ps *MainScene) Start() {
	oak.AddScene("mainScene", scene.Scene{
		Start: func(ctx *scene.Context) {
			handler := models.NewCarHandler()
			parking := models.NewParking(ctx)
			go runParking(ctx, handler, parking)
		},
	})
	oak.Init("mainScene")
}
