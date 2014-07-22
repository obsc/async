package async

import (
	"reflect"
)

type Async struct {
	ret  []reflect.Value
	done bool
	wait chan bool
}

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

func Unit(v ...interface{}) *Async {
	def := &Async{make([]reflect.Value, len(v)), true, make(chan bool)}
	for i := range v {
		def.ret[i] = reflect.ValueOf(v[i])
	}
	close(def.wait)
	return def
}

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

func (def *Async) IsDone() bool {
	return def.done
}
