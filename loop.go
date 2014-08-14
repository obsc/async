package async

import (
	"reflect"
)

// Type: M a
//
// Represents a looping async monad
type Loop struct {
	Async
	stop chan bool
}

func Repeat(f interface{}) *Loop {
	newLoop := &Loop{Async{nil, false, make(chan bool)}, make(chan bool)}
	go func() {
		for {
			select {
			case <-newLoop.stop:
				return
			default:
				ffun := reflect.ValueOf(f)
				newLoop.completeOne(ffun.Call(nil))
			}
		}
	}()
	return newLoop
}

func (self *Loop) Stop() {
	self.stop <- true
}
