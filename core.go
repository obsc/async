package async

import (
	"reflect"
)

// Type: M a
//
// Represents the async monad
type Async struct {
	ret  []reflect.Value
	done bool
	wait chan bool
}

// Type: (unit -> a) -> M a
//
// Deferred executes a function asynchronously and returns an Async
// struct representing the return values.
func Deferred(f interface{}) *Async {
	return withNewAsync(func() []reflect.Value {
		ffun := reflect.ValueOf(f)
		return ffun.Call(nil)
	})
}

// Type: a -> M a
//
// Return takes a set of inputs and returns an already determined Async
// struct wrapping those inputs
func Return(vs ...interface{}) *Async {
	return withNewAsync(func() []reflect.Value {
		ret := make([]reflect.Value, len(vs))
		for i := range vs {
			ret[i] = reflect.ValueOf(vs[i])
		}
		return ret
	})
}

// Type: M a -> (a -> b) -> M b
//
// Fmapping a function to self guarantees that function will execute on the
// return values of self as soon self is determined. It returns an Async
// struct representing the return values of the called function.
//
// If Fmap was called after self is already deteremined, the function will
// begin to execute right away.
//
// Note that while this represents a functor, the input types are reversed.
// This is done to make it a method rather than a function as Go does not
// allow creation of new operators.
func (self *Async) Fmap(f interface{}) *Async {
	return withNewAsync(func() []reflect.Value {
		self.Wait()

		ffun := reflect.ValueOf(f)
		ftyp := reflect.TypeOf(f)
		return ffun.Call(self.ret[0:ftyp.NumIn()])
	})
}

// Type: M a -> (a -> M b) -> M b
//
// Bind behaves the same way as Fmap except it operates on functions
// that return Async structs.
func (self *Async) Bind(f interface{}) *Async {
	return withNewAsync(func() []reflect.Value {
		self.Wait()

		ffun := reflect.ValueOf(f)
		ftyp := reflect.TypeOf(f)
		next := ffun.Call(self.ret[0:ftyp.NumIn()])[0].Interface().(*Async)
		next.Wait()

		return next.ret
	})
}

// Type: M (M a) -> M a
//
// Join is used to unify multiple layers of Async structs.
func (self *Async) Join(f interface{}) *Async {
	return withNewAsync(func() []reflect.Value {
		self.Wait()

		selfRet := self.ret[0].Interface().(*Async)
		selfRet.Wait()

		return selfRet.ret
	})
}

// Type: M a -> bool
//
// IsDone returns whether or not self has becomed determined yet.
func (self *Async) IsDone() bool {
	return self.done
}

// Type: M a -> unit
//
// Wait is a blocking function that waits until self has become determined.
func (self *Async) Wait() {
	for _ = range self.wait {
	}
}

// Type: (a -> b) -> (M a -> M b)
//
// Lift is another representation of Fmap with the types reversed.
//
// Lift can take any unary function and return a function that operates on
// Async structs.
func Lift(f interface{}) func(*Async) *Async {
	return func(self *Async) *Async {
		return self.Fmap(f)
	}
}

// Type: ((a, b) -> c) -> M a -> M b -> M c
//
// Lift2 behaves the same as lift except with binary functions.
//
// The parameters to the function are curried as well.
func Lift2(f interface{}) func(*Async) func(*Async) *Async {
	return func(self *Async) func(*Async) *Async {
		return func(other *Async) *Async {
			return self.And(other).Fmap(f)
		}
	}
}
