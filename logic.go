package async

import (
	"reflect"
)

// Type: M a -> M b -> M (a, b)
//
// And returns an Async struct that becomes determined when both self
// and other become determined. It contains return values of self and
// other concatenated together.
func (self *Async) And(other *Async) *Async {
	return withNewAsync(func() []reflect.Value {
		return concat(self.get(), other.get())
	})
}

// Type: [](M a) -> M ([]a)
//
// All returns an Async struct that becomes determined when all of
// asyncs becomes determined. It contains the return value of all the
// structs concatenated together.
func All(asyncs ...*Async) *Async {
	return withNewAsync(func() []reflect.Value {
		rets := make([][]reflect.Value, len(asyncs))
		for i := range asyncs {
			rets[i] = asyncs[i].get()
		}
		return concat(rets...)
	})
}

// Type: M a -> M a -> M a
//
// Or returns an Async struct that becomes determined whenever the first
// of either self or other becomes determined. It contains the return
// value of whichever one is determined first.
func (self *Async) Or(other *Async) *Async {
	return withNewAsync(self.get, other.get)
}

// Type: [](M a) -> M ([]a)
//
// Any returns an Async struct that becomes determined whenever any of
// asyncs becomes determined. It contains the return value of whichever
// one is determined first.
func Any(asyncs ...*Async) *Async {
	gets := make([]func() []reflect.Value, len(asyncs))
	for i := range asyncs {
		gets[i] = asyncs[i].get
	}
	return withNewAsync(gets...)
}

// Type: unit -> M unit
//
// Always represents an Async struct that is determined as soon as it is created.
func Always() *Async {
	return withNewAsync(func() []reflect.Value {
		return nil
	})
}

// Type: unit -> M unit
//
// Never represents an Async struct that is never determined.
func Never() *Async {
	return &Async{nil, false, make(chan bool)}
}
