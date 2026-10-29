package main

import (
	"fmt"
	"sync"
)

type GraphNode struct {
	id       int
	maxID    int
	isLeader bool
	inChans  []<-chan int
	outChans []chan<- int
}

func (n *GraphNode) Run(diameter int, wg *sync.WaitGroup) {
	defer wg.Done()

	for round := 1; round <= diameter; round++ {
		for _, out := range n.outChans {
			out <- n.maxID
		}

		nextMax := n.maxID
		for _, in := range n.inChans {
			receivedID := <-in
			if receivedID > nextMax {
				nextMax = receivedID
			}
		}
		n.maxID = nextMax
	}

	if n.maxID == n.id {
		n.isLeader = true
	}
}

func main() {
	adjList := map[int][]int{
		10: {25},
		25: {80, 45},
		80: {45, 10},
		45: {10},
	}
	diameter := 3 

	nodes := make(map[int]*GraphNode)
	for id := range adjList {
		nodes[id] = &GraphNode{
			id:    id,
			maxID: id, 
		}
	}

	for u, neighbors := range adjList {
		for _, v := range neighbors {
			ch := make(chan int, 1) 
			nodes[u].outChans = append(nodes[u].outChans, ch)
			nodes[v].inChans = append(nodes[v].inChans, ch)
		}
	}

	var wg sync.WaitGroup
	for _, node := range nodes {
		wg.Add(1)
		go node.Run(diameter, &wg)
	}

	wg.Wait()

	for _, id := range []int{10, 25, 45, 80} {
		node := nodes[id]
		status := "Non-Lider"
		if node.isLeader {
			status = "LIDER"
		}
		fmt.Printf("Nod ID = %2d ; max_id = %2d -> %s\n", node.id, node.maxID, status)
	}
}