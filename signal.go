package gobs

/*
#include <stdlib.h>
#include <obs/callback/signal.h>

extern void signal_cb(void *, calldata_t*);
*/
import "C"

import (
	"runtime/cgo"
	"unsafe"
)

type Callback func(data any, cd CallData)

type Signal interface {
	isSignal()
}

type SignalHandler struct {
	c *C.signal_handler_t
}

func (sh SignalHandler) Valid() bool {
	return sh.c != nil
}

type internal_calldata struct {
	callback Callback
	data     any
}

//export signal_cb
func signal_cb(data unsafe.Pointer, cd *C.calldata_t) {
	h := *(*cgo.Handle)(data)
	v := h.Value().(*internal_calldata)

	v.callback(v.data, CallData{c: cd})
}

func (sh SignalHandler) Connect(name string, callback Callback, data any) SignalConnection {
	cname := C.CString(name)
	handle := cgo.NewHandle(&internal_calldata{callback: callback, data: data})
	C.signal_handler_connect(sh.c, cname, (*[0]byte)(C.signal_cb), unsafe.Pointer(&handle))
	return SignalConnection{
		handler: sh,
		name:    cname,
		handle:  &handle,
	}
}

type SignalConnection struct {
	handler SignalHandler
	name    *C.char
	handle  *cgo.Handle
}

func (sc SignalConnection) Disconnect() {
	// #cgo noescape signal_handler_disconnect
	// #cgo nocallback signal_handler_disconnect
	C.signal_handler_disconnect(sc.handler.c, sc.name, (*[0]byte)(C.signal_cb), unsafe.Pointer(sc.handle))
	(*sc.handle).Delete()
	C.free(unsafe.Pointer(sc.name))
}
