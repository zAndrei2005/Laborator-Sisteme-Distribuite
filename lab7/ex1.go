package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

type Message struct {
	SenderID int
	Round    int
	Value    int
}

type Node struct {
	id          int
	x           int
	s           int
	n           int
	isByzantine bool
	inbox       chan Message
	majority    int
	mult        int
}

func (node *Node) BcastState(peers []*Node) {
	for _, peer := range peers {
		peer.inbox <- Message{
			SenderID: node.id,
			Round:    1,
			Value:    node.x,
		}
	}
}

func (node *Node) MaliciousBcastState(peers []*Node) {
	for _, peer := range peers {
		val := rand.Intn(2)
		peer.inbox <- Message{
			SenderID: node.id,
			Round:    1,
			Value:    val,
		}
	}
}

func (node *Node) BcastTie(peers []*Node) {
	for _, peer := range peers {
		peer.inbox <- Message{
			SenderID: node.id,
			Round:    2,
			Value:    node.majority,
		}
	}
}

func (node *Node) MaliciousBcastTie(peers []*Node) {
	for _, peer := range peers {
		peer.inbox <- Message{
			SenderID: node.id,
			Round:    2,
			Value:    rand.Intn(2),
		}
	}
}

func readInitialStates(n int) ([]int, error) {
	reader := bufio.NewReader(os.Stdin)
	states := make([]int, 0, n)

	for len(states) < n {
		b, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}
		if b == '0' {
			states = append(states, 0)
		} else if b == '1' {
			states = append(states, 1)
		}
	}

	return states, nil
}

func setByzantineRoles(nodes []*Node, s int) {
	for i := range nodes {
		nodes[i].isByzantine = i < s
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	var n, s int
	if _, err := fmt.Fscan(os.Stdin, &n, &s); err != nil {
		return
	}

	states, err := readInitialStates(n)
	if err != nil || len(states) < n {
		return
	}

	nodes := make([]*Node, n)
	for i := 0; i < n; i++ {
		nodes[i] = &Node{
			id:    i,
			x:     states[i],
			s:     s,
			n:     n,
			inbox: make(chan Message, 2*n),
		}
	}

	setByzantineRoles(nodes, s)

	for t := 0; ; t++ {
		var wg sync.WaitGroup

		wg.Add(n)
		for i := 0; i < n; i++ {
			go func(node *Node) {
				defer wg.Done()
				if node.isByzantine {
					node.MaliciousBcastState(nodes)
				} else {
					node.BcastState(nodes)
				}
			}(nodes[i])
		}
		wg.Wait()

		wg.Add(n)
		for i := 0; i < n; i++ {
			go func(node *Node) {
				defer wg.Done()
				ones := 0
				for j := 0; j < n; j++ {
					msg := <-node.inbox
					if msg.Round == 1 && msg.Value == 1 {
						ones++
					}
				}
				if ones >= n/2+1 {
					node.majority = 1
					node.mult = ones
				} else {
					node.majority = 0
					node.mult = n - ones
				}
			}(nodes[i])
		}
		wg.Wait()

		leader := nodes[t%n]
		if leader.isByzantine {
			leader.MaliciousBcastTie(nodes)
		} else {
			leader.BcastTie(nodes)
		}

		wg.Add(n)
		for i := 0; i < n; i++ {
			go func(node *Node) {
				defer wg.Done()
				tieMsg := <-node.inbox
				tie := tieMsg.Value

				if !node.isByzantine {
					if float64(node.mult) > float64(n)/2.0+float64(s) {
						node.x = node.majority
					} else {
						node.x = tie
					}
				}
			}(nodes[i])
		}
		wg.Wait()

		if t > s+1 {
			break
		}
	}

	for _, node := range nodes {
		if !node.isByzantine {
			fmt.Println(node.x)
			return
		}
	}
}