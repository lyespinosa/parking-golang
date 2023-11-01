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

func elementoVive(num int, almacen *[]int, espacioDisponible chan struct{}, mu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	tiempo := time.Duration(rand.Intn(3)+1) * time.Second
	time.Sleep(tiempo)

	mu.Lock()
	for i, val := range *almacen {
		if val == num {
			*almacen = append((*almacen)[:i], (*almacen)[i+1:]...)
			fmt.Printf("Elemento eliminado: %d. Almacen actualizado: %v, ocupados %v \n", num, *almacen, len(*almacen))
			espacioDisponible <- struct{}{}
			break
		}
	}
	mu.Unlock()
}

func almacenador(buffer <-chan int, almacen *[]int, espacioDisponible chan struct{}, mu *sync.Mutex, elementosEspera *[]int, wg *sync.WaitGroup) {
	for {
		select {
		case num, ok := <-buffer:
			if !ok {
				return
			}
			mu.Lock()
			if len(*almacen) < cap(*almacen) {
				*almacen = append(*almacen, num)
				fmt.Println("Almacenado el ", num, ": ", *almacen, "Ocupados:", len(*almacen))
				wg.Add(1)
				go elementoVive(num, almacen, espacioDisponible, mu, wg)
			} else {
				*elementosEspera = append(*elementosEspera, num)
			}
			mu.Unlock()
		case <-espacioDisponible:
			mu.Lock()
			if len(*elementosEspera) > 0 {
				num := (*elementosEspera)[0]
				*elementosEspera = (*elementosEspera)[1:]
				*almacen = append(*almacen, num)
				fmt.Println("Almacenado desde espera el ", num, " :", *almacen, "Ocupados:", len(*almacen))
				wg.Add(1)
				go elementoVive(num, almacen, espacioDisponible, mu, wg)
			}
			mu.Unlock()
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	var wg sync.WaitGroup
	mu := &sync.Mutex{}
	buffer := make(chan int, 5)
	almacen := make([]int, 0, 20)
	espacioDisponible := make(chan struct{}, cap(almacen))
	elementosEspera := make([]int, 0)

	wg.Add(1)
	go productor(buffer, &wg)

	go almacenador(buffer, &almacen, espacioDisponible, mu, &elementosEspera, &wg)

	wg.Wait()
}
