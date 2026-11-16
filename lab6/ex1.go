package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Info struct {
	ID    int
	Value int
}

type Memory map[Info]struct{}

type Node struct {
	ID        int
	Neighbors []int
	M         Memory
	Crashed   bool
}

func diametru(adj [][]int) int {
	n := len(adj)
	diam := 0
	for src := 0; src < n; src++ {
		dist := make([]int, n)
		for i := range dist {
			dist[i] = -1
		}
		dist[src] = 0
		queue := []int{src}

		for len(queue) > 0 {
			u := queue[0]
			queue = queue[1:]
			for _, v := range adj[u] {
				if dist[v] == -1 {
					dist[v] = dist[u] + 1
					queue = append(queue, v)
				}
			}
		}

		for _, d := range dist {
			if d == -1 {
				return -1 
			}
			if d > diam {
				diam = d
			}
		}
	}
	return diam
}

func copyMemory(src Memory) Memory {
	dst := make(Memory)
	for k := range src {
		dst[k] = struct{}{}
	}
	return dst
}

func main() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	scanner := bufio.NewScanner(os.Stdin)

	var n, s int
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		n, _ = strconv.Atoi(fields[0])
		s, _ = strconv.Atoi(fields[1])
		break
	}

	tokens := make([]int, n)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		for i := 0; i < n; i++ {
			tokens[i], _ = strconv.Atoi(fields[i])
		}
		break
	}

	adj := make([][]int, n)
	for i := 0; i < n; i++ {
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}
			fields := strings.Fields(line)
			for _, f := range fields {
				neighbor, _ := strconv.Atoi(f)
				adj[i] = append(adj[i], neighbor)
			}
			break
		}
	}

	diam := diametru(adj)
	if diam == -1 {
		fmt.Println("Graful nu este tare conex!")
		return
	}

	nodes := make([]Node, n)
	for i := 0; i < n; i++ {
		mem := make(Memory)
		mem[Info{ID: i, Value: tokens[i]}] = struct{}{}
		nodes[i] = Node{
			ID:        i,
			Neighbors: adj[i],
			M:         mem,
			Crashed:   false,
		}
	}

	crashesCount := 0
	t := 0
	maxIter := (s + 1) * diam

	for {
		inbox := make([]Memory, n)
		for i := 0; i < n; i++ {
			inbox[i] = make(Memory)
		}

		for i := 0; i < n; i++ {
			if nodes[i].Crashed {
				continue 
			}
			for _, neighbor := range nodes[i].Neighbors {
				if !nodes[neighbor].Crashed {
					for info := range nodes[i].M {
						inbox[neighbor][info] = struct{}{}
					}
				}
			}
		}

		if crashesCount < s && rng.Float64() < 0.5 {
			var activeNodes []int
			for i := 0; i < n; i++ {
				if !nodes[i].Crashed {
					activeNodes = append(activeNodes, i)
				}
			}
			if len(activeNodes) > 1 {
				victim := activeNodes[rng.Intn(len(activeNodes))]
				nodes[victim].Crashed = true
				crashesCount++
			}
		}

		for i := 0; i < n; i++ {
			if nodes[i].Crashed {
				continue
			}
			nextM := copyMemory(nodes[i].M)
			for info := range inbox[i] {
				nextM[info] = struct{}{}
			}
			nodes[i].M = nextM
		}

		if t > maxIter {
			break
		}

		t++
	}

	fmt.Printf("Diametrul grafului: %d\n", diam)
	fmt.Printf("Iteratia de oprire: %d\n", t)

	var correctNode *Node
	for i := 0; i < n; i++ {
		if !nodes[i].Crashed {
			correctNode = &nodes[i]
			break
		}
	}

	if correctNode != nil {
		var result []Info
		for info := range correctNode.M {
			result = append(result, info)
		}
		sort.Slice(result, func(i, j int) bool {
			return result[i].ID < result[j].ID
		})

		for i, info := range result {
			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Printf("(%d, %d)", info.ID, info.Value)
		}
		fmt.Println()
	}
}