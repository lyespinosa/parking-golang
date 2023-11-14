package models

import (
	"image/color"

	"github.com/oakmound/oak/v4/entities"
)

var groundColor = color.RGBA{52, 52, 52, 255}

type Parking struct {
	entity *entities.Entity
}
