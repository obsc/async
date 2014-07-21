package async

import (
	"reflect"
)

type deffered struct {
	ret  []reflect.Value
	done bool
	wait chan bool
}

func Deffered(f interface{}) *deffered {
	def := &deffered{nil, false, make(chan bool)}
	go func() {
		ffun := reflect.ValueOf(f)
		def.ret = ffun.Call(nil)
		def.done = true
		close(def.wait)
	}()
	return def
}

func Init(v ...interface{}) *deffered {
	def := &deffered{make([]reflect.Value, len(v)), true, make(chan bool)}
	for i := range v {
		def.ret[i] = reflect.ValueOf(v[i])
	}
	close(def.wait)
	return def
}

func (def *deffered) IsDone() bool {
	return def.done
}

func (def *deffered) Bind(f interface{}) *deffered {
	newdef := &deffered{nil, false, make(chan bool)}
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
