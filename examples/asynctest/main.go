package main

import (
	"fmt"
	"github.com/obsc/async"
	"time"
)

func timer() {
	fmt.Println("starting timer")
	time.Sleep(3 * time.Second)
	fmt.Println("ending timer")
}

func a() (m int, n int) {
	fmt.Println("in a")
	timer()
	m = 10
	n = 20
	return
}

func b(x int, y int) {
	fmt.Println("in b")
	timer()
	fmt.Println(x + y)
}

func main() {
	async.Deferred(timer)
	async.Deferred(a).Fmap(b).Fmap(a)
	async.Unit(5, 10).Fmap(b).Fmap(a)
	async.Deferred(a).
		Fmap(b).
		Fmap(a)
	time.Sleep(10 * time.Second)
	fmt.Println("done")
}
