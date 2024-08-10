package lock

import (
	"sync"
	"testing"
)

func BenchmarkStack(b *testing.B) {

	b.Run("lock", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			var wg sync.WaitGroup
			s := NewStack[int]()

			for i := 0; i < 1; i++ {
				s.Push(1)
			}

			Push := func() {
				for i := 0; i < 1; i++ {
					s.Push(1)
				}
				wg.Done()
			}

			Pop := func() {
				for i := 0; i < 1; i++ {
					s.Pop()
				}
				wg.Done()
			}

			wg.Add(5)
			go Push()
			go Push()

			go Pop()
			go Pop()
			go Pop()

			wg.Wait()
		}
	})
}
