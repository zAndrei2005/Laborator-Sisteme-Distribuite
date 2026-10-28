package main

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"
)

type LamportClock struct{ counter int }

func (l *LamportClock) Tick() {
	l.counter++
}

func (l *LamportClock) Send() {
	l.counter++
}

func (l *LamportClock) Receive(ts int) {
	if ts > l.counter {
		l.counter = ts
	}
	l.counter++
}

type Message struct{ sender, ts int }

type Event struct {
	process          int
	internal_counter int
	kind             string 
}

func less(a, b Event) bool {
	if a.internal_counter != b.internal_counter {
		return a.internal_counter < b.internal_counter
	}
	return a.process < b.process
}

func runProcess(id int, out chan<- Message, in <-chan Message, eventsChan chan<- Event, wg *sync.WaitGroup) {
	defer wg.Done()
	clock := &LamportClock{counter: 0}

	for runda := 0; runda < 2; runda++ {
		randomTicks := rand.Intn(3) + 1
		for i := 0; i < randomTicks; i++ {
			clock.Tick()
		}
		eventsChan <- Event{
			process:          id,
			internal_counter: clock.counter,
			kind:             "internal",
		}

		clock.Send()
		msgOut := Message{sender: id, ts: clock.counter}
		eventsChan <- Event{
			process:          id,
			internal_counter: clock.counter,
			kind:             "send",
		}
		out <- msgOut

		if rand.Intn(2) == 1 {
			clock.Tick()
			eventsChan <- Event{
				process:          id,
				internal_counter: clock.counter,
				kind:             "internal",
			}
		}

		msgIn := <-in
		clock.Receive(msgIn.ts)
		eventsChan <- Event{
			process:          id,
			internal_counter: clock.counter,
			kind:             "recv",
		}
	}
}

func main() {

	ch12 := make(chan Message, 5)
	ch23 := make(chan Message, 5)
	ch31 := make(chan Message, 5)

	eventsChan := make(chan Event, 100)

	var wg sync.WaitGroup
	wg.Add(3)

	go runProcess(1, ch12, ch31, eventsChan, &wg) 
	go runProcess(2, ch23, ch12, eventsChan, &wg) 
	go runProcess(3, ch31, ch23, eventsChan, &wg) 

	wg.Wait()
	close(eventsChan)

	var events []Event
	for ev := range eventsChan {
		events = append(events, ev)
	}

	sort.Slice(events, func(i, j int) bool {
		return less(events[i], events[j])
	})

	fmt.Println("Evenimentele ordonate crescator:")
	for _, ev := range events {
		fmt.Printf("Proces P%d; Ceas Lamport: %2d; Tip: %-8s\n",
			ev.process, ev.internal_counter, ev.kind)
	}
}