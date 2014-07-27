/*
Package async provides a golang representation of the async/promise/deferred
monad along with useful methods to operate on these monads.

Additional information
    http://en.wikipedia.org/wiki/Futures_and_promises
    http://en.wikipedia.org/wiki/Monad_%28functional_programming%29
    http://en.wikipedia.org/wiki/Monad_%28category_theory%29

Installation
    go get github.com/obsc/async

The Async struct represents a value or slice of values that become determined
when some asynchronous operation complete.

An Async struct can either be "determined" or "undetermined" at any point
in time. However, it can only be "determined" once, and will henceforth
remain determined forever.

Unless otherwise specified, every library call will return immediately
and not block any other operations. Because of this, it is recommended that
any operation you write that returns an Async struct, do so immediately,
allowing the program to continue to other tasks.

While many of the methods can be written purely in terms of Return and Bind,
they are still implemented separately to minimize the overheard due to
the creation of more goroutines and more channels.

This package is loosely based off of Jane Street's OCaml library: Async
    https://ocaml.janestreet.com/ocaml-core/111.17.00/doc/async/#Std
*/
package async
