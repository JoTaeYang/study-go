package lfstack_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/JoTaeYang/study-go/pkg/lockfree/lfstack"
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

		for i = 0; i < ARRAY_CNT; i ++ {
			atomic.AddInt32(&pList[i].Hp, 1)
			atomic.AddInt32(&pList[i].Att, 1)
		}
	}
}

func TestLfStack(t *testing.T) {
	s = &lfstack.Stack[*Player]{}

	pList := make([]*Player, POOL_COUNT, POOL_COUNT)
	for i := 0; i < POOL_COUNT; i++ {
		pList[i] = &Player{
			Hp:  0x0000555,
			Att: 0,
			Def: 5,
		}
	}

	for i := 0; i < POOL_COUNT; i++ {
		s.Push(pList[i])
	}


	for i := 0; i < GO_CNT; i ++ {
		go 
	}
}
