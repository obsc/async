package async

import (
	"reflect"
)

// M a -> M b -> M (a, b)
func (self *Async) And(other *Async) *Async {
	return withNewAsync(func() []reflect.Value {
		return concat(self.get(), other.get())
	})
}

// [](M a) -> M ([]a)
func All(asyncs ...*Async) *Async {
	return withNewAsync(func() []reflect.Value {
		rets := make([][]reflect.Value, len(asyncs))
		for i := range asyncs {
			rets[i] = asyncs[i].get()
		}
		return concat(rets...)
	})
}

// M a -> M a -> M a
func (self *Async) Or(other *Async) *Async {
	return withNewAsync(self.get, other.get)
}

// [](M a) -> M ([]a)
func Any(asyncs ...*Async) *Async {
	gets := make([]func() []reflect.Value, len(asyncs))
	for i := range asyncs {
		gets[i] = asyncs[i].get
	}
	return withNewAsync(gets...)
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
