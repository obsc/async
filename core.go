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

func concat(rets ...[]reflect.Value) []reflect.Value {
	lens := 0
	for _, v := range rets {
		lens += len(v)
	}
	ret := make([]reflect.Value, lens)
	curIndex := 0
	for _, reti := range rets {
		copy(ret[curIndex:curIndex+len(reti)], reti)
		curIndex += len(reti)
	}
	return ret
}

// (unit -> a) -> M a
func Deferred(f interface{}) *Async {
	def := &Async{nil, false, make(chan bool)}
	go func() {
		ffun := reflect.ValueOf(f)
		def.ret = ffun.Call(nil)
		def.done = true
		close(def.wait)
	}()
	return def
}

// a -> M a
func Return(v ...interface{}) *Async {
	def := &Async{make([]reflect.Value, len(v)), true, make(chan bool)}
	for i := range v {
		def.ret[i] = reflect.ValueOf(v[i])
	}
	close(def.wait)
	return def
}

// M a -> (a -> b) -> M b
func (def *Async) Fmap(f interface{}) *Async {
	newdef := &Async{nil, false, make(chan bool)}
	go func() {
		for _ = range def.wait {
		}

		ffun := reflect.ValueOf(f)
		ftyp := reflect.TypeOf(f)
		newdef.ret = ffun.Call(def.ret[0:ftyp.NumIn()])
		newdef.done = true
		close(newdef.wait)
	}()
	return newdef
}

// M a -> (a -> M b) -> M b
func (def *Async) Bind(f interface{}) *Async {
	newdef := &Async{nil, false, make(chan bool)}
	go func() {
		for _ = range def.wait {
		}

		ffun := reflect.ValueOf(f)
		ftyp := reflect.TypeOf(f)
		next := ffun.Call(def.ret[0:ftyp.NumIn()])[0].Interface().(*Async)

		for _ = range next.wait {
		}

		newdef.ret = next.ret
		newdef.done = true
		close(newdef.wait)
	}()
	return newdef
}

// M (M a) -> M a
func (def *Async) Join(f interface{}) *Async {
	newdef := &Async{nil, false, make(chan bool)}
	go func() {
		for _ = range def.wait {
		}
		defret := def.ret[0].Interface().(*Async)

		for _ = range defret.wait {
		}

		newdef.ret = defret.ret
		newdef.done = true
		close(newdef.wait)
	}()
	return newdef
}

// M a -> bool
func (def *Async) IsDone() bool {
	return def.done
}

// M a -> unit
func (def *Async) Wait() {
	for _ = range def.wait {
	}
}

// M a -> (unit -> b) -> M b
func (def *Async) Next(f interface{}) *Async {
	newdef := &Async{nil, false, make(chan bool)}
	go func() {
		for _ = range def.wait {
		}

		ffun := reflect.ValueOf(f)
		newdef.ret = ffun.Call(nil)
		newdef.done = true
		close(newdef.wait)
	}()
	return newdef
}

// M a -> M b -> M (a, b)
func (def *Async) And(other *Async) *Async {
	newdef := &Async{nil, false, make(chan bool)}
	go func() {
		for _ = range def.wait {
		}
		for _ = range other.wait {
		}

		newdef.ret = concat(def.ret, other.ret)
		newdef.done = true
		close(newdef.wait)
	}()
	return newdef
}

// (a -> b) -> (M a -> M b)
func Lift(f interface{}) func(*Async) *Async {
	return func(def *Async) *Async {
		return def.Fmap(f)
	}
}

// ((a, b) -> c) -> M a -> M b -> M c
func Lift2(f interface{}) func(*Async) func(*Async) *Async {
	return func(def *Async) func(*Async) *Async {
		return func(other *Async) *Async {
			return def.And(other).Fmap(f)
		}
	}
}
