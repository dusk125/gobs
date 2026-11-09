package gobs

// As of Go 1.27, noescape only changes code generation when the same function
// is also marked nocallback (cmd/cgo/out.go), and nocallback is unsafe here:
// the Go log handler (log.go) and Go signal handlers mean almost any libobs
// call can call back into Go, which panics under nocallback.
// #cgo pkg-config: libobs
// #include <obs/obs.h>
import "C"
