package lfstack_test

import (
	"testing"

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

func FuncTestRoutine() {

}

func TestLfStack(t *testing.T) {
	s := lfstack.Stack{}
}
