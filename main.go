package main

import (
	"image/color"
	"math/rand"
	"sync"
	"time"

	"github.com/oakmound/oak/v4"
	"github.com/oakmound/oak/v4/event"
	"github.com/oakmound/oak/v4/render"
	"github.com/oakmound/oak/v4/scene"
)

const (
	espacios      = 20
	anchoEspacio  = 30
	altoEspacio   = 30
	margenEspacio = 5
	velocidad     = 2 // Píxeles por frame
)

type Dato struct {
	render.Sprite
	destX, destY float64
	almacenado   bool
	mu           sync.Mutex
}

func NuevoDato(x, y, destX, destY float64) *Dato {
	cubo := render.NewColorBox(anchoEspacio-10, altoEspacio-10, color.RGBA{100, 100, 250, 255})
	cubo.SetPos(x, y)
	return &Dato{
		Sprite: cubo,
		destX:  destX,
		destY:  destY,
	}
}

// Mover actualiza la posición del Dato hacia su destino.
func (d *Dato) Mover() bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.almacenado {
		return false
	}

	x, y := d.GetPos()
	if x < d.destX {
		x += velocidad
	}
	if y < d.destY {
		y += velocidad
	}

	// Si el Dato está lo suficientemente cerca de su destino, lo marca como almacenado
	if int(x) >= int(d.destX)-velocidad && int(y) >= int(d.destY)-velocidad {
		d.almacenado = true
		d.SetPos(d.destX, d.destY)
		return true
	}

	d.SetPos(x, y)
	return false
}

func main() {
	rand.Seed(time.Now().UnixNano())

	oak.AddScene("estacionamiento", scene.Scene{
		Start: func(ctx *scene.Context) {
			espaciosRender := make([]*render.ColorBox, espacios)

			// Crear espacios de estacionamiento
			for i := 0; i < espacios; i++ {
				x := margenEspacio + (i%5)*(anchoEspacio+margenEspacio)
				y := margenEspacio + (i/5)*(altoEspacio+margenEspacio)

				espacio := render.NewColorBox(anchoEspacio, altoEspacio, color.RGBA{150, 150, 150, 255})
				espacio.SetPos(float64(x), float64(y))
				ctx.DrawStack.Draw(espacio, 0)
				espaciosRender[i] = espacio
			}

			// Crear y mover un cubo (dato) hacia un espacio de estacionamiento
			dato := NuevoDato(100, 100, float64(margenEspacio), float64(margenEspacio))
			ctx.DrawStack.Draw(dato.Sprite, 1)

			// Mover el Dato en cada tick
			event.GlobalBind(ctx.CID, int(event.Tick), func(_ event.CID, _ interface{}) int {
				if almacenado := dato.Mover(); almacenado {
					// Marcar el espacio de estacionamiento como ocupado
					i := int(dato.destY/margenEspacio)*5 + int(dato.destX/margenEspacio)
					espaciosRender[i].SetFillColor(color.RGBA{100, 250, 100, 255})
				}
				return 0
			})
		},
	})

	// Iniciar el juego con la escena del estacionamiento
	oak.Init("estacionamiento")
}
