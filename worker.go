package Worker

import (
	"math/rand"
	"sync"
)

type Work struct {
	f       func(any)
	total   int
	items   map[any]bool
	inQueue []any
	cond    *sync.Cond
	todo    int
}

// init the worker
func (w *Work) init() {
	if w.items == nil {
		w.items = make(map[any]bool)
	}

	if w.cond == nil {
		w.cond = sync.NewCond(&sync.Mutex{})
	}
}

// Add adds a new item to the work queue
func (w *Work) Add(item any) {
	w.init()
	w.cond.L.Lock()

	if !w.items[item] {
		w.items[item] = true
		w.inQueue = append(w.inQueue, item)

		if w.todo > 0 {
			w.cond.Signal()
		}
	}
	w.cond.L.Unlock()
}

func (w *Work) Run(n int, f func(item any)) {
	if n < 1 {
		panic("n < 1")
	}

	if w.total >= 1 {
		panic("already called")
	}

	w.total = n
	w.f = f

	for i := 0; i < n-1; i++ {
		go w.worker()
	}

	w.worker()
}

func (w *Work) worker() {
	for {
		w.cond.L.Lock()

		for len(w.inQueue) == 0 {
			w.todo++

			//
			if w.todo == w.total {
				w.cond.Broadcast()
				w.cond.L.Unlock()
				return
			}

			w.cond.Wait()
			w.todo--
		}

		cur := rand.Intn(len(w.inQueue))

		item := w.inQueue[cur]
		// cut
		w.inQueue[cur] = w.inQueue[len(w.inQueue)-1]
		w.inQueue = w.inQueue[:len(w.inQueue)-1]

		w.cond.L.Unlock()
		// run
		w.f(item)
	}
}
