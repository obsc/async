package async

import (
	"reflect"
)

// M a -> M b -> M (a, b)
func (self *Async) And(other *Async) *Async {
	newAsync := &Async{nil, false, make(chan bool)}
	go func() {
		self.Wait()
		other.Wait()

		newAsync.complete(concat(self.ret, other.ret))
	}()
	return newAsync
}

// M a -> M a -> M a
func (self *Async) Or(other *Async) *Async {
	newAsync := &Async{nil, false, make(chan bool)}
	waitFor := func(a *Async) {
		a.Wait()

		newAsync.complete(a.ret)
	}
	go waitFor(self)
	go waitFor(other)
	return newAsync
}

// unit -> M unit
func Always() *Async {
	return withNewAsync(func(newAsync *Async) []reflect.Value {
		return nil
	})
}

// unit -> M unit
func Never() *Async {
	return &Async{nil, false, make(chan bool)}
}
