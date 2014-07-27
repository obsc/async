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
