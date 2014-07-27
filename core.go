package async

import (
	"reflect"
)

// M a
type Async struct {
	ret  []reflect.Value
	done bool
	wait chan bool
}

// (unit -> a) -> M a
func Deferred(f interface{}) *Async {
	return withNewAsync(func() []reflect.Value {
		ffun := reflect.ValueOf(f)
		return ffun.Call(nil)
	})
}

// a -> M a
func Return(vs ...interface{}) *Async {
	return withNewAsync(func() []reflect.Value {
		ret := make([]reflect.Value, len(vs))
		for i := range vs {
			ret[i] = reflect.ValueOf(vs[i])
		}
		return ret
	})
}

// M a -> (a -> b) -> M b
func (self *Async) Fmap(f interface{}) *Async {
	return withNewAsync(func() []reflect.Value {
		self.Wait()

		ffun := reflect.ValueOf(f)
		ftyp := reflect.TypeOf(f)
		return ffun.Call(self.ret[0:ftyp.NumIn()])
	})
}

// M a -> (a -> M b) -> M b
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

// M (M a) -> M a
func (self *Async) Join(f interface{}) *Async {
	return withNewAsync(func() []reflect.Value {
		self.Wait()

		selfRet := self.ret[0].Interface().(*Async)
		selfRet.Wait()

		return selfRet.ret
	})
}

// M a -> bool
func (self *Async) IsDone() bool {
	return self.done
}

// M a -> unit
func (self *Async) Wait() {
	for _ = range self.wait {
	}
}

// (a -> b) -> (M a -> M b)
func Lift(f interface{}) func(*Async) *Async {
	return func(self *Async) *Async {
		return self.Fmap(f)
	}
}

// ((a, b) -> c) -> M a -> M b -> M c
func Lift2(f interface{}) func(*Async) func(*Async) *Async {
	return func(self *Async) func(*Async) *Async {
		return func(other *Async) *Async {
			return self.And(other).Fmap(f)
		}
	}
}
