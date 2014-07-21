package main

import (
	"fmt"
	"github.com/obsc/async"
	"time"
)

func timer() {
	// fmt.Println("starting timer")
	time.Sleep(3 * time.Second)
	// fmt.Println("ending timer")
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
	// async.Deffered(timer)
	async.Deffered(a).Bind(b).Bind(a)
	async.Init(5, 10).Bind(b).Bind(a)
	time.Sleep(10 * time.Second)
	fmt.Println("done")
}
