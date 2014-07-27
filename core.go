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
	newAsync := &Async{nil, false, make(chan bool)}
	go func() {
		ffun := reflect.ValueOf(f)
		newAsync.ret = ffun.Call(nil)
		newAsync.done = true
		close(newAsync.wait)
	}()
	return newAsync
}

// a -> M a
func Return(vs ...interface{}) *Async {
	newAsync := &Async{make([]reflect.Value, len(vs)), true, make(chan bool)}
	for i := range vs {
		newAsync.ret[i] = reflect.ValueOf(vs[i])
	}
	newAsync.done = true
	close(newAsync.wait)
	return newAsync
}

// M a -> (a -> b) -> M b
func (self *Async) Fmap(f interface{}) *Async {
	newAsync := &Async{nil, false, make(chan bool)}
	go func() {
		for _ = range self.wait {
		}

		ffun := reflect.ValueOf(f)
		ftyp := reflect.TypeOf(f)
		newAsync.ret = ffun.Call(self.ret[0:ftyp.NumIn()])
		newAsync.done = true
		close(newAsync.wait)
	}()
	return newAsync
}

// M a -> (a -> M b) -> M b
func (self *Async) Bind(f interface{}) *Async {
	newAsync := &Async{nil, false, make(chan bool)}
	go func() {
		for _ = range self.wait {
		}

		ffun := reflect.ValueOf(f)
		ftyp := reflect.TypeOf(f)
		next := ffun.Call(self.ret[0:ftyp.NumIn()])[0].Interface().(*Async)

		for _ = range next.wait {
		}

		newAsync.ret = next.ret
		newAsync.done = true
		close(newAsync.wait)
	}()
	return newAsync
}

// M (M a) -> M a
func (self *Async) Join(f interface{}) *Async {
	newAsync := &Async{nil, false, make(chan bool)}
	go func() {
		for _ = range self.wait {
		}
		selfRet := self.ret[0].Interface().(*Async)

		for _ = range selfRet.wait {
		}

		newAsync.ret = selfRet.ret
		newAsync.done = true
		close(newAsync.wait)
	}()
	return newAsync
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

// unit -> M a
func Never() *Async {
	return &Async{nil, false, make(chan bool)}
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
