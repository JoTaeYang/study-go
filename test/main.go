package main

import (
	"log"
	"sync/atomic"
	"time"

	"github.com/JoTaeYang/study-go/pkg/lockfree/lfstack"
	"github.com/JoTaeYang/study-go/pkg/lockfree/queue"
)

const (
	GO_CNT     = 4
	ARRAY_CNT  = 2
	POOL_COUNT = GO_CNT * ARRAY_CNT
)

type Player struct {
	Hp  int32
	Att int32
	Def int32
}

var s *lfstack.Stack[*Player]
var q *queue.Queue[*Player]

func FuncTestRoutineQ() {
	pList := make([]*Player, ARRAY_CNT, ARRAY_CNT)

	var i int
	for {
		for i = 0; i < ARRAY_CNT; i++ {
			q.Dequeue(&pList[i])

			if pList[i] == nil {
				log.Println("Player is Nil")
				panic("Player is Nil")
			}

			if pList[i].Hp != 0x0000555 {
				log.Println("Player Hp Is Not Init Data Equal")
				panic("Player Hp Is Not Init Data Equal")
			}

			if pList[i].Att != 0 {
				log.Println("Player Att Is Not Init Data Equal")
				panic("Player Att Is Not Init Data Equal")
			}
		}
		time.Sleep(0)

		for i = 0; i < ARRAY_CNT; i++ {
			atomic.AddInt32(&pList[i].Hp, 1)
			atomic.AddInt32(&pList[i].Att, 1)
		}

		time.Sleep(0)
		for i = 0; i < ARRAY_CNT; i++ {
			if pList[i].Hp != 0x0000556 {
				log.Println("Player Hp Is Not Init Data Equal")
				panic("Player Hp Is Not Init Data Equal")
			}

			if pList[i].Att != 1 {
				log.Println("Player Att Is Not Init Data Equal")
				panic("Player Att Is Not Init Data Equal")
			}
		}

		time.Sleep(0)
		for i = 0; i < ARRAY_CNT; i++ {
			atomic.AddInt32(&pList[i].Hp, -1)
			atomic.AddInt32(&pList[i].Att, -1)
		}

		time.Sleep(0)
		for i = 0; i < ARRAY_CNT; i++ {
			if pList[i].Hp != 0x0000555 {
				log.Println("Player Hp Is Not Init Data Equal")
				panic("Player Hp Is Not Init Data Equal")
			}

			if pList[i].Att != 0 {
				log.Println("Player Att Is Not Init Data Equal")
				panic("Player Att Is Not Init Data Equal")
			}

			q.Enqueue(pList[i])
		}
	}
}

func FuncTestRoutine() {
	pList := make([]*Player, ARRAY_CNT, ARRAY_CNT)

	var check bool
	var i int
	for {
		for i = 0; i < ARRAY_CNT; i++ {
			pList[i], check = s.Pop()
			if !check {
				log.Println("Stack Is Empty")
				panic("Stack Is Empty")
			}

			if pList[i] == nil {
				log.Println("Player is Nil")
				panic("Player is Nil")
			}

			if pList[i].Hp != 0x0000555 {
				log.Println("Player Hp Is Not Init Data Equal")
				panic("Player Hp Is Not Init Data Equal")
			}

			if pList[i].Att != 0 {
				log.Println("Player Att Is Not Init Data Equal")
				panic("Player Att Is Not Init Data Equal")
			}
		}
		time.Sleep(0)

		for i = 0; i < ARRAY_CNT; i++ {
			atomic.AddInt32(&pList[i].Hp, 1)
			atomic.AddInt32(&pList[i].Att, 1)
		}

		time.Sleep(0)
		for i = 0; i < ARRAY_CNT; i++ {
			if pList[i].Hp != 0x0000556 {
				log.Println("Player Hp Is Not Init Data Equal")
				panic("Player Hp Is Not Init Data Equal")
			}

			if pList[i].Att != 1 {
				log.Println("Player Att Is Not Init Data Equal")
				panic("Player Att Is Not Init Data Equal")
			}
		}

		time.Sleep(0)
		for i = 0; i < ARRAY_CNT; i++ {
			atomic.AddInt32(&pList[i].Hp, -1)
			atomic.AddInt32(&pList[i].Att, -1)
		}

		time.Sleep(0)
		for i = 0; i < ARRAY_CNT; i++ {
			if pList[i].Hp != 0x0000555 {
				log.Println("Player Hp Is Not Init Data Equal")
				panic("Player Hp Is Not Init Data Equal")
			}

			if pList[i].Att != 0 {
				log.Println("Player Att Is Not Init Data Equal")
				panic("Player Att Is Not Init Data Equal")
			}

			s.Push(pList[i])
		}
	}
}
func main() {
	s = lfstack.NewStack[*Player]()
	q = queue.NewQueue[*Player]()
	pList := make([]*Player, POOL_COUNT, POOL_COUNT)
	for i := 0; i < POOL_COUNT; i++ {
		pList[i] = &Player{
			Hp:  0x0000555,
			Att: 0,
			Def: 5,
		}
	}

	for i := 0; i < POOL_COUNT; i++ {
		q.Enqueue(pList[i])
	}

	player := &Player{}
	q.Dequeue(&player)

	for i := 0; i < GO_CNT; i++ {
		go FuncTestRoutineQ()
	}

	for {
		time.Sleep(time.Second)
		if s.GetCount() > POOL_COUNT {
			panic("pool count over")
		}

		log.Println("lock free test ing ..")
	}
	// s = lfstack.NewStack[*Player]()

	// pList := make([]*Player, POOL_COUNT, POOL_COUNT)
	// for i := 0; i < POOL_COUNT; i++ {
	// 	pList[i] = &Player{
	// 		Hp:  0x0000555,
	// 		Att: 0,
	// 		Def: 5,
	// 	}
	// }

	// for i := 0; i < POOL_COUNT; i++ {
	// 	s.Push(pList[i])
	// }

	// for i := 0; i < GO_CNT; i++ {
	// 	go FuncTestRoutine()
	// }

	// for {
	// 	time.Sleep(time.Second)
	// 	if s.GetCount() > POOL_COUNT {
	// 		panic("pool count over")
	// 	}

	// 	log.Println("lock free test ing ..")
	// }
}
