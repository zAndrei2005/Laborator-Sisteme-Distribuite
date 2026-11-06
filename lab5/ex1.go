package main

import (
	"fmt"
	"math"
	"os"
)

type Point struct {
	X float64
	Y float64
}

func distance(a, b Point) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return math.Sqrt(dx*dx + dy*dy)
}

func compressRobots(robots []Point, eps float64) []Point {
	var compressed []Point
	for _, r := range robots {
		duplicate := false
		for _, c := range compressed {
			if distance(r, c) < eps {
				duplicate = true
				break
			}
		}
		if !duplicate {
			compressed = append(compressed, r)
		}
	}
	return compressed
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Eroare la deschiderea fisierului:", err)
		return
	}
	defer file.Close()

	var n int
	if _, err := fmt.Fscan(file, &n); err != nil {
		fmt.Println("Eroare la citirea numarului de roboti:", err)
		return
	}

	robots := make([]Point, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(file, &robots[i].X, &robots[i].Y)
	}

	const (
		commRadius = 15.0 
		eps        = 1e-4 
		maxIter    = 1000 
	)

	fmt.Printf("Start: %d roboti activi | Raza de comunicare: %.2f\n", len(robots), commRadius)

	robots = compressRobots(robots, eps)

	iter := 0
	for iter < maxIter && len(robots) > 1 {
		iter++
		newRobots := make([]Point, len(robots))
		moved := false

		for i := 0; i < len(robots); i++ {
			sumX := robots[i].X
			sumY := robots[i].Y
			count := 1 

			for j := 0; j < len(robots); j++ {
				if i != j && distance(robots[i], robots[j]) <= commRadius {
					sumX += robots[j].X
					sumY += robots[j].Y
					count++
				}
			}

			newRobots[i] = Point{
				X: sumX / float64(count),
				Y: sumY / float64(count),
			}

			if distance(robots[i], newRobots[i]) > 1e-9 {
				moved = true
			}
		}

		robots = newRobots

		robots = compressRobots(robots, eps)

		fmt.Printf("Iteratia %d: %d roboti ramasi -> %v\n", iter, len(robots), robots)

		if !moved {
			fmt.Println("Stare stationara atinsa (robotii nu mai au vecini in raza de comunicare).")
			break
		}
	}
	
	if len(robots) == 1 {
		fmt.Printf("Rendezvous reusit in %d iteratii! Punctul de intalnire: (%.4f, %.4f)\n",
			iter, robots[0].X, robots[0].Y)
	} else {
		fmt.Printf("Algoritmul s-a oprit dupa %d iteratii cu %d clustere separate:\n", iter, len(robots))
		for idx, r := range robots {
			fmt.Printf("  Cluster %d: (%.4f, %.4f)\n", idx+1, r.X, r.Y)
		}
	}
}