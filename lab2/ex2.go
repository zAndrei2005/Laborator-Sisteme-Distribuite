package main

import "fmt"

func primestemesaj(ch chan string, done chan bool) {
	mesaj := <-ch 
	fmt.Print(mesaj)
	done <- true
}

func main() {
	ch := make(chan string)

	done := make(chan bool)
	go primestemesaj(ch, done)
	ch <- "yo"
	<-done
}