package main

import "fmt"

func proc(done chan bool) {
	done <- true
}

func main() {
	done := make(chan bool)

	go proc(done)

	<- done
	fmt.Print("gata")
}