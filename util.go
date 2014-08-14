package async

import (
	"reflect"
)

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

func withNewAsync(fs ...func() []reflect.Value) *Async {
	newAsync := &Async{nil, false, make(chan bool)}
	for i := range fs {
		go newAsync.completeAll(fs[i]())
	}
	return newAsync
}

func (self *Async) completeAll(ret []reflect.Value) {
	self.completeOne(ret)
	self.done = true
}

func (self *Async) completeOne(ret []reflect.Value) {
	defer recover()
	close(self.wait)
	self.ret = ret
}

func (self *Async) get() []reflect.Value {
	self.Wait()
	return self.ret
}
