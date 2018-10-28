package main

import (
	"fmt"
	"github.com/ValeryPiashchynski/Worker"
	"io/ioutil"
	"net/http"
	"runtime"
)

func main() {
	var w Worker.Work

	sites := []string{"http://google.com", "http://amazon.com", "http://spiralscout.com", "http://0xdev.me"}

	for _, v := range sites {
		w.Add(v)
	}

	fmt.Println(runtime.NumGoroutine())

	w.Run(10, f)

	fmt.Println(runtime.NumGoroutine())
}

func f(x interface{}) {
	str := x.(string)

	r, err := http.Get(str)
	if err != nil {
		panic(err)
	}

	_, err = ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	//fmt.Print(string(b))

}
