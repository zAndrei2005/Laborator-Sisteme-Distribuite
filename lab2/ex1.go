package main

import "fmt"

func trimitemesaj(ch chan string){
	ch <- "salut" 
}

func main(){
	ch := make(chan string)
	go trimitemesaj(ch)
	mesaj := <-ch
	fmt.Print(mesaj)
}