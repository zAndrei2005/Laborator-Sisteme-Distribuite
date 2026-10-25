package main

import "fmt"

func trimitemesaj(ch chan string) {
	ch <- "sal"
	ch <- "ut"
}

func main() {
	ch := make(chan string)

	go trimitemesaj(ch, ch2)
	mesaj := <-ch 
	mesaj2 := <- ch
	fmt.Print(mesaj, mesaj2)
}