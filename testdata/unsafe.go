package main

import (
	"unsafe"
)


func useUnsafe() uintptr {
	var i int;
	p := unsafe.Pointer(uintptr(i) + 0);

	return uintptr(p);
}
