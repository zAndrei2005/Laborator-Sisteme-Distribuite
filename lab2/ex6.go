package main

import (
	"fmt"
	"math/rand"
	"time"
)

func simulareRobot(nume string, finish chan string) {
	durata := time.Duration(rand.Intn(500)) * time.Millisecond
	finish <- nume
}

func main() {
	chR1 := make(chan string)
	chR2 := make(chan string)

	go simulareRobot("Robot1", chR1)
	go simulareRobot("Robot2", chR2)

	select {
	case castigator := <- chR1:
		fmt.Print("%s a ajuns primul la finish!\n", castigator)
	case castigator := <- chR2:
		fmt.Print("%s a ajuns primul la finish!\n", castigator)
	}
}