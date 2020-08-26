// run

// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"runtime"
)

type T struct {
	io.Closer
}

func f1() {
	// The 5 here and below depends on the number of internal runtime frames
	// that sit between a deferred function called during panic and
	// the original frame. If that changes, this test will start failing and
	// the number here will need to be updated.
	defer checkLine(5)
	var t *T
	var c io.Closer = t
	c.Close()
}

func f2() {
	defer checkLine(5)
	var t T
	var c io.Closer = t
	c.Close()
}

func main() {
	f1()
	f2()
}

func checkLine(n int) {
	if err := recover(); err == nil {
		panic("did not panic")
	}
	var file string
	var line int
	for i := 1; i <= n; i++ {
		_, file, line, _ = runtime.Caller(i)
		if file != "<autogenerated>" || line != 1 {
			continue
		}
		return
	}
	panic(fmt.Sprintf("expected <autogenerated>:1 have %s:%d", file, line))
}
