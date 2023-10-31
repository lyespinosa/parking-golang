package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func productor(buffer chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 30; i++ {
		fmt.Println("+ Prodocido: ", i)
		buffer <- i
		time.Sleep(time.Second)
	}
	close(buffer)
}

func almacenador(buffer <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	almacen := make([]int, 0, 20)

	for num := range buffer {
		if len(almacen) < 20 {
			almacen = append(almacen, num)
			fmt.Println("Almacenado", almacen)
		} else {
			index := rand.Intn(20)
			fmt.Printf("Sacando: %d, del lugar: %d\n", almacen[index], index)
			almacen[index] = num
			fmt.Println("Almacen actualizado", almacen)
		}
		time.Sleep(2 * time.Second)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	var wg sync.WaitGroup
	buffer := make(chan int, 5)

	wg.Add(1)
	go productor(buffer, &wg)

	wg.Add(1)
	go almacenador(buffer, &wg)

	wg.Wait()
}
