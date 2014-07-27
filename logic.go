package async

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

// M a -> M a -> M a
func (def *Async) Or(other *Async) *Async {
	newdef := &Async{nil, false, make(chan bool)}
	waitFor := func(d *Async) {
		defer recover()
		for _ = range d.wait {
		}

		close(newdef.wait)
		newdef.ret = d.ret
		newdef.done = true
	}
	go waitFor(def)
	go waitFor(other)
	return newdef
}
