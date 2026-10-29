package main

import (
	"fmt"
	"math"
	"sync"
)

type VotingNode struct {
	id       int
	x        float64
	inChans  []<-chan float64
	outChans []chan<- float64
}

func (n *VotingNode) Step(wg *sync.WaitGroup) {
	defer wg.Done()

	for _, out := range n.outChans {
		out <- n.x
	}

	sumNeighbors := 0.0
	for _, in := range n.inChans {
		sumNeighbors += <-in
	}

	di := float64(len(n.inChans))
	n.x = (n.x + sumNeighbors) / (di + 1.0)
}

func main() {
	adjList := map[int][]int{
		0: {1, 4},
		1: {0, 2},
		2: {1, 3},
		3: {2, 4},
		4: {3, 0},
	}

	initialVotes := map[int]float64{
		0: 1.0,
		1: 0.0,
		2: 1.0,
		3: 1.0,
		4: 0.0,
	}

	nodes := make(map[int]*VotingNode)
	for id, vote := range initialVotes {
		nodes[id] = &VotingNode{
			id: id,
			x:  vote,
		}
	}

	for u, neighbors := range adjList {
		for _, v := range neighbors {
			ch := make(chan float64, 1)
			nodes[u].outChans = append(nodes[u].outChans, ch)
			nodes[v].inChans = append(nodes[v].inChans, ch)
		}
	}

	const eps = 1e-6
	iteration := 0

	for {
		minVal, maxVal := math.Inf(1), math.Inf(-1)
		for _, node := range nodes {
			if node.x < minVal {
				minVal = node.x
			}
			if node.x > maxVal {
				maxVal = node.x
			}
		}
		if maxVal-minVal < eps {
			break
		}

		iteration++

		var wg sync.WaitGroup
		for _, node := range nodes {
			wg.Add(1)
			go node.Step(&wg)
		}
		wg.Wait()
	}

	fmt.Printf("Convergenta atinsa in %d iteratii (eps = %g):\n", iteration, eps)
	for i := 0; i < len(nodes); i++ {
		fmt.Printf("  Nod %d -> stare finala x_%d = %.6f\n", i, i, nodes[i].x)
	}

	consensus := nodes[0].x
	fmt.Printf("Valoarea de consens obtinuta: %.6f\n", consensus)

	if consensus > 0.5 {
		fmt.Println("=> REZULTAT VOT MAJORITAR: 1 (deoarece media > 1/2)")
	} else if consensus < 0.5 {
		fmt.Println("=> REZULTAT VOT MAJORITAR: 0 (deoarece media < 1/2)")
	} else {
		fmt.Println("=> REZULTAT VOT MAJORITAR: Egalitate (media == 1/2)")
	}
}