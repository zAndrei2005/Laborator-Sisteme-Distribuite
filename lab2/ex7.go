package main

import (
	"fmt"
	"time"
)

type Punct struct {
	X, Y int
}

func robotSender(ch chan Punct, rep chan string) {
	for i := 1; i <= 3; i++ {
		ch <- Punct{X: i, Y: i * 2}
		time.Sleep(100 * time.Millisecond) 
	}
	close(ch)            
	rep <- "Sender oprit"
}

func robotReceiver(ch chan Punct, rep chan string) {
	for coord := range ch {
		fmt.Printf("Coordonate citite: X=%d Y=%d\n", coord.X, coord.Y)
	}
	rep <- "Receiver oprit (toate pozițiile citite)"
}

func main() {
	chCoordonate := make(chan Punct)
	raportSender := make(chan string)
	raportReceiver := make(chan string)
	
	go robotSender(chCoordonate, raportSender)
	go robotReceiver(chCoordonate, raportReceiver)
	
	fmt.Println("Raport:", <-raportSender)
	fmt.Println("Raport:", <-raportReceiver)
}