package Worker

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestWork(t *testing.T) {
	var w Work

	const N = 1000000
	n := int32(0)
	w.Add(N)
	w.Run(100, func(x interface{}) {
		atomic.AddInt32(&n, 1)
		i := x.(int)
		if i >= 2 {
			w.Add(i - 1)
			w.Add(i - 2)
		}
		w.Add(i >> 1)
		w.Add((i >> 1) ^ 1)
	})
	if n != N+1 {
		t.Fatalf("ran %d items, expected %d", n, N+1)
	}
}

func TestWorkParallel(t *testing.T) {
	for tries := 0; tries < 10; tries++ {
		var w Work
		const N = 100
		for i := 0; i < N; i++ {
			w.Add(i)
		}
		start := time.Now()
		var n int32
		w.Run(N, func(x interface{}) {
			time.Sleep(1 * time.Millisecond)
			atomic.AddInt32(&n, +1)
		})
		if n != N {
			t.Fatalf("Work.Run did not do all the work")
		}
		if time.Since(start) < N/2*time.Millisecond {
			return
		}
	}
	t.Fatalf("Work.Run does not seem to be parallel")
}

