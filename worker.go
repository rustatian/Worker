package Worker

import (
	"math/rand"
	"sync"
)

type Work struct {
	f       func(interface{})
	running int
	sync.Mutex
	items     map[interface{}]bool
	inQueue   []interface{}
	cond      sync.Cond
	waitToRun int
}

func (w *Work) initialize() {
	if w.items == nil {
		w.items = make(map[interface{}]bool)
	}
	//w.cond = *sync.NewCond(w)
}

func (w *Work) Add(item interface{}) {
	w.Lock()
	w.initialize()

	if !w.items[item] {
		w.items[item] = true
		w.inQueue = append(w.inQueue, item)

		if w.waitToRun > 0 {
			w.cond.Signal()
		}
	}
	w.Unlock()
}

func (w *Work) Run(n int, f func(item interface{})) {
	if n < 1 {
		panic("n < 1")
	}

	if w.running >= 1 {
		panic("already called")
	}

	w.running = n
	w.f = f
	w.cond.L = w

	for i:=0; i < n-1; i ++ {
		go w.worker()
	}

	w.worker()
}

func (w *Work) worker() {
	for {
		w.Lock()
		for len(w.inQueue) == 0 {
			w.waitToRun ++

			//
			if w.waitToRun == w.running {
				w.cond.Broadcast()
				w.Unlock()
				return
			}

			w.cond.Wait()
			w.waitToRun--
		}

		cur := rand.Intn(len(w.inQueue))

		item := w.inQueue[cur]
		// cut
		w.inQueue[cur] = w.inQueue[len(w.inQueue)-1]
		w.inQueue = w.inQueue[:len(w.inQueue)-1]
		w.Unlock()

		// run
		w.f(item)
	}
}
