package main

import (
	"fmt"
	"github.com/obsc/async"
	"time"
)

func loop(n int, f func(i int)) {
	for i := 0; i < n; i++ {
		f(i)
	}
}

func delay(wait int) func(string, int) (string, int) {
	return func(prefix string, start int) (string, int) {
		time.Sleep(time.Duration(wait) * time.Second)
		return prefix, start
	}
}

func step(prefix string, v int, wait int) {
	fmt.Printf("%v: %v\n", prefix, v)
	time.Sleep(time.Duration(wait) * time.Second)
}

func cycle(interval int, num int) func(string, int) (string, int) {
	return func(prefix string, start int) (string, int) {
		loop(num, func(i int) {
			step(prefix, start+i*interval, interval)
		})
		return prefix, start + num*interval
	}
}

func main() {
	async.Return("Actual", 0).
		Fmap(cycle(3, 5)).
		Fmap(cycle(3, 5))
	async.Return("Actual", 1).
		Fmap(delay(1)).
		Fmap(cycle(3, 5)).
		Fmap(cycle(3, 5))
	async.Return("Actual", 2).
		Fmap(delay(2)).
		Fmap(cycle(3, 5)).
		Fmap(cycle(3, 5))
	async.Return("Expected", 0).Fmap(cycle(1, 30)).Wait()
	fmt.Println("done")
}
