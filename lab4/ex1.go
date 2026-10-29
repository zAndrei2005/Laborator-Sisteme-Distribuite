package main

import (
	"fmt"
	"sync"
)

type RingNode struct {
	id 	int
	isLeader bool
	in	<-chan int
	out chan<- int

}

func (n *RingNode) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	n.out <- n.id

	for receivedID := range n.in {
		if receivedID == -1 {
			n.out <- -1
			return
		}

		if receivedID > n.id {
			n.out <- receivedID
		} else if receivedID == n.id {
			n.isLeader = true
			n.out <- -1
			return
		}
	}
}

func main() {
	ids := []int{14, 7, 42, 3, 99, 25}
	numNodes := len(ids)

	chans := make([]chan int, numNodes)
	for i := 0; i < numNodes; i++ {
		chans[i] = make(chan int, numNodes)
	}

	nodes := make([]*RingNode, numNodes)
	var wg sync.WaitGroup

	for i := 0; i < numNodes; i++ {
		leftChan := chans[(i-1+numNodes)%numNodes]
		rightChan := chans[i]
		nodes[i] = &RingNode{
			id:  ids[i],
			in:  leftChan,
			out: rightChan,
		}
		wg.Add(1)
		go nodes[i].Run(&wg)
	}

	wg.Wait()

	for i, node := range nodes {
		status := "Non-Lider"
		if node.isLeader {
			status = "LIDER"
		}
		fmt.Printf("Nod pozitia %d (ID = %2d) -> %s\n", i, node.id, status)
	}

}