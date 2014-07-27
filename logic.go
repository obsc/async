package async

import (
	"reflect"
)

// M a -> M b -> M (a, b)
func (self *Async) And(other *Async) *Async {
	return withNewAsync(func() []reflect.Value {
		self.Wait()
		other.Wait()

		return concat(self.ret, other.ret)
	})
}

// M a -> M a -> M a
func (self *Async) Or(other *Async) *Async {
	waitFor := func(a *Async) func() []reflect.Value {
		return func() []reflect.Value {
			a.Wait()

			return a.ret
		}
	}

	return withNewAsync(waitFor(self), waitFor(other))
}

// unit -> M unit
func Always() *Async {
	return withNewAsync(func() []reflect.Value {
		return nil
	})
}

// unit -> M unit
func Never() *Async {
	return &Async{nil, false, make(chan bool)}
}
