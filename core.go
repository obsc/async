package async

import (
	"reflect"
)

type deferred struct {
	ret  []reflect.Value
	done bool
	wait chan bool
}

func Deferred(f interface{}) *deferred {
	def := &deferred{nil, false, make(chan bool)}
	go func() {
		ffun := reflect.ValueOf(f)
		def.ret = ffun.Call(nil)
		def.done = true
		close(def.wait)
	}()
	return def
}

func Return(v ...interface{}) *deferred {
	def := &deferred{make([]reflect.Value, len(v)), true, make(chan bool)}
	for i := range v {
		def.ret[i] = reflect.ValueOf(v[i])
	}
	close(def.wait)
	return def
}

func (def *deferred) IsDone() bool {
	return def.done
}

func (def *deferred) Fmap(f interface{}) *deferred {
	newdef := &deferred{nil, false, make(chan bool)}
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

func (def *deferred) Bind(f interface{}) *deferred {
	newdef := &deferred{nil, false, make(chan bool)}
	go func() {
		for _ = range def.wait {
		}

		ffun := reflect.ValueOf(f)
		ftyp := reflect.TypeOf(f)
		next := ffun.Call(def.ret[0:ftyp.NumIn()])[0].Interface().(*deferred)

		for _ = range next.wait {
		}

		newdef.ret = next.ret
		newdef.done = true
		close(newdef.wait)
	}()
	return newdef
}

func (def *deferred) Join(f interface{}) *deferred {
	newdef := &deferred{nil, false, make(chan bool)}
	go func() {
		for _ = range def.wait {
		}
		defret := def.ret[0].Interface().(*deferred)

		for _ = range defret.wait {
		}

		newdef.ret = defret.ret
		newdef.done = true
		close(newdef.wait)
	}()
	return newdef
}
