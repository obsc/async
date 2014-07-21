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
		ff := reflect.ValueOf(f)
		def.ret = ff.Call(nil)
		def.done = true
		close(def.wait)
	}()
	return def
}

func (def *deffered) IsDone() bool {
	return def.done
}
