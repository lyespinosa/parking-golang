package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func productor(buffer chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 100; i++ {
		fmt.Println("Producido:", i)
		buffer <- i
		time.Sleep(time.Millisecond * 100)
	}
	close(buffer)
}

func elementoVive(num int, almacen *[]int, espacioDisponible chan<- bool, mu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	tiempo := time.Duration(rand.Intn(3)+1) * time.Second
	time.Sleep(tiempo)

	mu.Lock()
	// Elimina el elemento después de que ha vivido su tiempo aleatorio.
	for i, val := range *almacen {
		if val == num {
			*almacen = append((*almacen)[:i], (*almacen)[i+1:]...)
			fmt.Printf("Elemento eliminado: %d. Almacen actualizado: %v, ocupados %v \n", num, *almacen, len(*almacen)) // Editado
			if len(*almacen) < cap(*almacen) {
				// Si hay espacio disponible, envía una señal al canal.
				espacioDisponible <- true // Editado
			}
			break
		}
	}
	mu.Unlock()
}

func almacenador(buffer <-chan int, almacen *[]int, espacioDisponible chan<- bool, mu *sync.Mutex, wg *sync.WaitGroup) {
	for num := range buffer {
		mu.Lock()
		if len(*almacen) < cap(*almacen) {
			*almacen = append(*almacen, num)
			fmt.Println("Almacenado el ", num, ": ", *almacen, "Ocupados:", len(*almacen))
			wg.Add(1)
			go elementoVive(num, almacen, espacioDisponible, mu, wg)
		} else {
			// Si el almacen está lleno, envía una señal al canal.
			espacioDisponible <- true // Editado
		}
		mu.Unlock()
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	var wg sync.WaitGroup
	mu := &sync.Mutex{}
	buffer := make(chan int, 5)          // Canal para los elementos producidos.
	almacen := make([]int, 0, 20)        // Almacen con capacidad para 20 elementos.
	espacioDisponible := make(chan bool) // Canal para señalizar cuando hay espacio disponible.

	wg.Add(1)
	go productor(buffer, &wg)

	// Inicia el almacenador en su propia goroutine.
	go almacenador(buffer, &almacen, espacioDisponible, mu, &wg)

	// Gestiona el espacio disponible en el almacen.
	go func() {
		for range espacioDisponible {
			mu.Lock()
			if len(almacen) < cap(almacen) && len(buffer) > 0 {
				num := <-buffer // Toma el siguiente elemento en espera.
				almacen = append(almacen, num)
				fmt.Println("Almacenado desde espera:", &almacen)
				wg.Add(1)
				go elementoVive(num, &almacen, espacioDisponible, mu, &wg)
			}
			mu.Unlock()
		}
	}()

	wg.Wait() // Espera a que todos los elementos hayan sido procesados.
}
