package main

import "fmt"

func sendcom(ch chan string) {
	ch <- "salut"
}

func recvcom(chin chan string, chout chan string) { 
	mesaj := <-chin
	chout <- mesaj
}

func main() {
	chcom := make(chan string)
	chrez := make(chan string)

	go sendcom(chcom)
	go recvcom(chcom, chrez)

	rez := <-chrez 
	fmt.Print(rez)
}

