package main

import (
	"container/heap"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/Gordy96/evt-sim/internal"
)

type MessageKind uint64

type Message struct {
	ID          string
	Src         string
	Dst         string
	Kind        MessageKind
	TimestampNS time.Time
	Params      map[string]interface{}
}

func (m *Message) Priority() int64 {
	return m.TimestampNS.UnixNano()
}

type Node struct {
	ID        string
	Params    map[string]interface{}
	Init      func(self *Node)
	OnMessage func(self *Node, msg *Message, timestamp time.Time)
}

func main() {
	pq := internal.PriorityQueue[*Message]{}
	heap.Init(&pq)

	makeSender := func() func(self *Node, msg *Message, timestamp time.Time) {
		return func(self *Node, msg *Message, timestamp time.Time) {
			j, _ := json.Marshal(msg)
			fmt.Printf("[%s] node '%s' received message: %s\n", time.Now().Format(time.RFC3339Nano), self.ID, j)
			itries, ok := self.Params["retries"]
			var tries int
			if !ok {
				self.Params["retries"] = 3
				tries = 3
			} else {
				tries = itries.(int)
			}

			if tries > 0 {
				i, _ := strconv.ParseInt(msg.ID, 10, 64)
				heap.Push(&pq, &Message{
					ID:          fmt.Sprintf("%d", i+1),
					Kind:        msg.Kind,
					Src:         self.ID,
					Dst:         msg.Src,
					TimestampNS: timestamp.Add(100 * time.Millisecond),
				})
				tries--
				self.Params["retries"] = tries
			}
		}
	}

	nodes := []Node{
		{
			ID: "first",
			Init: func(self *Node) {
				heap.Push(&pq, &Message{
					ID:          "1",
					Kind:        0,
					Src:         self.ID,
					Dst:         "second",
					TimestampNS: time.Now(),
					Params: map[string]interface{}{
						"foo": "bar",
					},
				})
			},
			OnMessage: makeSender(),
			Params:    make(map[string]interface{}),
		}, {
			ID:        "second",
			Init:      func(self *Node) {},
			OnMessage: makeSender(),
			Params:    make(map[string]interface{}),
		},
	}

	for _, node := range nodes {
		if node.Init != nil {
			node.Init(&node)
		}
	}

	for pq.Len() > 0 {
		msg := heap.Pop(&pq).(*Message)
		for _, node := range nodes {
			if node.ID == msg.Dst {
				node.OnMessage(&node, msg, msg.TimestampNS)
			}
		}
	}
}
