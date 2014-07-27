package async

// M a -> M b -> M (a, b)
func (self *Async) And(other *Async) *Async {
	newAsync := &Async{nil, false, make(chan bool)}
	go func() {
		for _ = range self.wait {
		}
		for _ = range other.wait {
		}

		newAsync.ret = concat(self.ret, other.ret)
		newAsync.done = true
		close(newAsync.wait)
	}()
	return newAsync
}

// M a -> M a -> M a
func (self *Async) Or(other *Async) *Async {
	newAsync := &Async{nil, false, make(chan bool)}
	waitFor := func(a *Async) {
		defer recover()
		for _ = range a.wait {
		}

		close(newAsync.wait)
		newAsync.ret = a.ret
		newAsync.done = true
	}
	go waitFor(self)
	go waitFor(other)
	return newAsync
}
