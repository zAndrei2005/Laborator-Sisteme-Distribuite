package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type Clock struct{ offset time.Duration }

func (c *Clock) Now() time.Time { return time.Now().Add(c.offset) }

func estimate(t1, t2, t3, t4 time.Time) (theta, delta time.Duration) {
	theta = (t2.Sub(t1) + t3.Sub(t4)) / 2
	delta = t4.Sub(t1) - t3.Sub(t2)
	return theta, delta
}

func measure(c_clock, s_clock *Clock) (theta, delta time.Duration) {
	t1 := c_clock.Now()
	time.Sleep(3*time.Millisecond + time.Duration(rand.Intn(10))*time.Millisecond)

	t2 := s_clock.Now()
	time.Sleep(5*time.Millisecond + time.Duration(rand.Intn(20))*time.Millisecond)

	t3 := s_clock.Now()

	time.Sleep(3*time.Millisecond + time.Duration(rand.Intn(10))*time.Millisecond)
	t4 := c_clock.Now()

	return estimate(t1, t2, t3, t4)
}

func absDuration(d time.Duration) time.Duration {
	if d < 0 {
		return -d
	}
	return d
}

func main() {
	initialOffset := 200 * time.Millisecond

	// a
	clientClock := &Clock{offset: initialOffset}
	serverClock := &Clock{offset: 0} 

	thetaA, deltaA := measure(clientClock, serverClock)
	newOffsetA := clientClock.offset + thetaA

	fmt.Printf("Theta estimat: %v si Delta: %v\n", thetaA, deltaA)
	fmt.Printf("Offset dupa corectie: %v (A scazut? %v)\n\n",
		newOffsetA, absDuration(newOffsetA) < absDuration(initialOffset))

	// b
	minDelta := time.Duration(math.MaxInt64)
	minTheta := time.Duration(math.MaxInt64)
	var bestThetaForMinDelta time.Duration

	for i := 1; i <= 10; i++ {
		cClock := &Clock{offset: initialOffset}
		sClock := &Clock{offset: 0}

		theta, delta := measure(cClock, sClock)
		correctedOffset := cClock.offset + theta

		fmt.Printf("Run %2d -> theta: %10v | delta: %10v | offset ramas: %10v\n",
			i, theta, delta, correctedOffset)

		if delta < minDelta {
			minDelta = delta
			bestThetaForMinDelta = theta
		}
		if absDuration(theta) < absDuration(minTheta) {
			minTheta = theta
		}
	}

	fmt.Println("\nComparatie rezultate b vs a")
	fmt.Printf("Delta minim obtinut in 10 rulari: %v (fata de %v la pct. a)\n", minDelta, deltaA)
	fmt.Printf("Theta minim: %v\n", minTheta)
	fmt.Printf("Offsetul corectat folosind incercarea cu delta minim: %v (fata de %v la pct. a)\n",
		initialOffset+bestThetaForMinDelta, newOffsetA)
}